package bot

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"go.uber.org/zap"
	"golang.org/x/net/html"
)

type (
	TaskManager struct {
		mu    sync.Mutex
		tasks map[uint]*TaskRunner
	}
)

var BotTaskManager *TaskManager

type TaskRunner struct {
	Task     *bot.BotTask
	StopChan chan struct{}
}

func (tm *TaskManager) StartTask(task *bot.BotTask) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 如果已经有任务在跑，先停掉
	if runner, ok := tm.tasks[task.ID]; ok {
		close(runner.StopChan)
		delete(tm.tasks, task.ID)
	}

	runner := &TaskRunner{
		Task:     task,
		StopChan: make(chan struct{}),
	}
	tm.tasks[task.ID] = runner

	go runner.Run()
	global.GVA_LOG.Info("TaskManager StartTask", zap.Any("task", runner))
}

func (tm *TaskManager) StopTask(taskID uint) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if runner, ok := tm.tasks[taskID]; ok {
		close(runner.StopChan)
		delete(tm.tasks, taskID)
	}
	global.GVA_LOG.Info("TaskManager StopTask", zap.Any("task", taskID))
}

func (tm *TaskManager) ReloadTask(task *bot.BotTask) {
	if task.Status == 1 {
		tm.StopTask(task.ID)
		tm.StartTask(task)
	} else {
		tm.StopTask(task.ID)
	}
}

func (tm *TaskManager) StopAll() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, runner := range tm.tasks {
		close(runner.StopChan)
	}
	tm.tasks = make(map[uint]*TaskRunner)
	global.GVA_LOG.Info("TaskManager StopAll")
}

func (tr *TaskRunner) Run() {
	for {
		select {
		case <-tr.StopChan:
			return
		default:
			now := time.Now()

			// 如果任务设置了 StopTime 且已到期，停止任务
			if !tr.Task.StopTime.IsZero() && now.After(tr.Task.StopTime) {
				tr.Task.Status = 2
				global.GVA_DB.Model(tr.Task).Update("status", 2)
				global.GVA_LOG.Info("TaskRunner StopTask due to StopTime", zap.Int("taskID", int(tr.Task.ID)))
				return
			}

			// 更新 NextSendTime，保证它在未来
			for !tr.Task.NextSendTime.After(now) {
				tr.Task.NextSendTime = tr.Task.NextSendTime.Add(time.Duration(tr.Task.SendInterval) * time.Minute)
				// 如果下次发送时间已经超过 StopTime，则任务结束
				if !tr.Task.StopTime.IsZero() && tr.Task.NextSendTime.After(tr.Task.StopTime) {
					tr.Task.Status = 2
					global.GVA_DB.Model(tr.Task).Update("status", 2)
					global.GVA_LOG.Info("TaskRunner StopTask due to StopTime (next send)", zap.Int("taskID", int(tr.Task.ID)))
					return
				}
			}

			global.GVA_LOG.Info("TaskRunner Run", zap.Int("taskID", int(tr.Task.ID)), zap.Time("nextSendTime", tr.Task.NextSendTime))

			wait := time.Until(tr.Task.NextSendTime)
			timer := time.NewTimer(wait)

			select {
			case <-tr.StopChan:
				timer.Stop()
				return
			case <-timer.C:
				err := SendTelegramMessage(int64(tr.Task.GroupID), tr.Task)
				if err != nil {
					global.GVA_LOG.Error("TaskRunner Run", zap.Int("taskID", int(tr.Task.ID)), zap.Error(err))
					continue
				} else {
					global.GVA_LOG.Info("TaskRunner Run", zap.Int("taskID", int(tr.Task.ID)), zap.Error(err))
				}

				now = time.Now()
				tr.Task.PreSendTime = &now
				tr.Task.NextSendTime = now.Add(time.Duration(tr.Task.SendInterval) * time.Minute)

				// 如果下次发送时间已经超过 StopTime，则任务结束
				if !tr.Task.StopTime.IsZero() && tr.Task.NextSendTime.After(tr.Task.StopTime) {
					tr.Task.Status = 2
					global.GVA_DB.Model(tr.Task).Updates(map[string]interface{}{
						"pre_send_time":  tr.Task.PreSendTime,
						"next_send_time": tr.Task.NextSendTime,
						"status":         2,
					})
					global.GVA_LOG.Info("TaskRunner StopTask due to StopTime after send", zap.Int("taskID", int(tr.Task.ID)))
					return
				}

				global.GVA_DB.Model(tr.Task).Updates(map[string]interface{}{
					"pre_send_time":  tr.Task.PreSendTime,
					"next_send_time": tr.Task.NextSendTime,
				})
			}
		}
	}
}

func extractImgsAndText(htmlStr string) (imgs []string, textWithoutImgs string) {
	// 找所有 <img ... src="..."> 并取 src
	imgRe := regexp.MustCompile(`(?i)<img[^>]+src=["']([^"']+)["'][^>]*>`)
	matches := imgRe.FindAllStringSubmatch(htmlStr, -1)
	for _, m := range matches {
		if len(m) >= 2 {
			imgs = append(imgs, m[1])
		}
	}
	// 去掉所有 <img> 标签，保留其余 HTML
	textWithoutImgs = imgRe.ReplaceAllString(htmlStr, "")
	// 去掉空的 <p></p> 等多余空白（可选）
	// 将 HTML 实体转回正常字符（比如 &gt; -> >）
	textWithoutImgs = html.UnescapeString(textWithoutImgs)
	textWithoutImgs = strings.TrimSpace(textWithoutImgs)
	return
}

// 辅助：创建 InlineKeyboardMarkup（如果没有按钮则返回 nil）
func buildMarkupFromExtrend(raw json.RawMessage) *tgbotapi.InlineKeyboardMarkup {
	if len(raw) == 0 {
		return nil
	}
	var btns []bot.ButtonItem
	if err := json.Unmarshal(raw, &btns); err != nil || len(btns) == 0 {
		return nil
	}
	var row []tgbotapi.InlineKeyboardButton
	for _, b := range btns {
		// 只创建 URL 按钮（和你之前逻辑保持一致）
		row = append(row, tgbotapi.NewInlineKeyboardButtonURL(b.Name, b.URL))
	}
	m := tgbotapi.NewInlineKeyboardMarkup(row)
	return &m
}

func CleanHTMLForTelegram(htmlStr string) string {
	if htmlStr == "" {
		return ""
	}

	htmlStr = regexp.MustCompile(`(?i)<p[^>]*>`).ReplaceAllString(htmlStr, "")
	htmlStr = regexp.MustCompile(`(?i)</p>`).ReplaceAllString(htmlStr, "\n")
	htmlStr = regexp.MustCompile(`(?i)<div[^>]*>`).ReplaceAllString(htmlStr, "")
	htmlStr = regexp.MustCompile(`(?i)</div>`).ReplaceAllString(htmlStr, "\n")
	htmlStr = regexp.MustCompile(`(?i)<br\s*/?>`).ReplaceAllString(htmlStr, "\n")

	allowed := map[string]bool{
		"b": true, "i": true, "strong": true, "em": true,
		"u": true, "s": true, "strike": true, "del": true,
		"a": true, "code": true, "pre": true,
	}

	tagRe := regexp.MustCompile(`(?i)</?([a-z0-9]+)[^>]*>`)
	htmlStr = tagRe.ReplaceAllStringFunc(htmlStr, func(tag string) string {
		m := tagRe.FindStringSubmatch(tag)
		if len(m) < 2 {
			return ""
		}
		name := strings.ToLower(m[1])
		if allowed[name] {
			return tag
		}
		return ""
	})

	htmlStr = html.UnescapeString(htmlStr)

	htmlStr = regexp.MustCompile(`\n{2,}`).ReplaceAllString(htmlStr, "\n\n")
	htmlStr = strings.TrimSpace(htmlStr)
	return htmlStr
}

func extractVideosFromHTML(content string) []string {
	var urls []string
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return urls
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "source" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					urls = append(urls, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return urls
}

func SendTelegramMessage(chatID int64, task *bot.BotTask) error {
	var botModel bot.Bot
	var has bool
	var err error
	if botModel, has, err = dao.BotDao.FromBotID(global.GVA_DB, int(task.BotID)); !has || err != nil {
		if !has {
			err = fmt.Errorf("bot %d not found", task.BotID)
		}
		return err
	}

	botAPI, err := tgbotapi.NewBotAPI(botModel.Token)
	if err != nil {
		return err
	}

	// 处理按钮（二维数组）
	var markup *tgbotapi.InlineKeyboardMarkup
	if len(task.ExtrendButton) > 0 {
		var btns [][]bot.ButtonItem
		if err := json.Unmarshal(task.ExtrendButton, &btns); err == nil && len(btns) > 0 {
			var rows [][]tgbotapi.InlineKeyboardButton
			for _, row := range btns {
				var telegramRow []tgbotapi.InlineKeyboardButton
				for _, b := range row {
					telegramRow = append(telegramRow, tgbotapi.NewInlineKeyboardButtonURL(b.Name, b.URL))
				}
				rows = append(rows, telegramRow)
			}
			m := tgbotapi.NewInlineKeyboardMarkup(rows...)
			markup = &m
		}
	}

	switch task.TaskSendType {
	case 0: // 仅按钮
		if markup != nil {
			msg := tgbotapi.NewMessage(chatID, "")
			msg.ReplyMarkup = markup
			_, err = botAPI.Send(msg)
		}
		return err

	case 1: // 富文本
		imgs, text := extractImgsAndText(task.Content)
		videos := extractVideosFromHTML(task.Content)

		caption := CleanHTMLForTelegram(text)
		if len(caption) > 1024 {
			caption = caption[:1020] + "..."
		}

		// 先发送图片
		for i, img := range imgs {
			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(img))
			if i == 0 {
				photo.Caption = caption
				photo.ParseMode = tgbotapi.ModeHTML
				if markup != nil {
					photo.ReplyMarkup = markup
				}
			}
			if _, err := botAPI.Send(photo); err != nil {
				return err
			}
		}

		// 发送视频
		for i, vid := range videos {
			video := tgbotapi.NewVideo(chatID, tgbotapi.FileURL(vid))
			if i == 0 && len(imgs) == 0 { // 如果前面没有图片，第一条视频附加 caption
				video.Caption = caption
				video.ParseMode = tgbotapi.ModeHTML
				if markup != nil {
					video.ReplyMarkup = markup
				}
			} else if i == 0 && markup != nil {
				video.ReplyMarkup = markup
			}
			if _, err := botAPI.Send(video); err != nil {
				return err
			}
		}

		// 如果没有图片和视频，则发送纯文本
		if len(imgs) == 0 && len(videos) == 0 && caption != "" {
			msg := tgbotapi.NewMessage(chatID, caption)
			msg.ParseMode = tgbotapi.ModeHTML
			if markup != nil {
				msg.ReplyMarkup = markup
			}
			_, err = botAPI.Send(msg)
		}
		return err

	case 2: // 文本
		msg := tgbotapi.NewMessage(chatID, task.Content)
		if markup != nil {
			msg.ReplyMarkup = markup
		}
		_, err = botAPI.Send(msg)
		return err

	case 3: // 图片
		var urls []string
		if err := json.Unmarshal([]byte(task.Content), &urls); err != nil {
			return err
		}
		for i, url := range urls {
			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(url))
			if i == 0 && markup != nil {
				photo.ReplyMarkup = markup
			}
			if _, err := botAPI.Send(photo); err != nil {
				return err
			}
		}
		return nil

	case 4: // 视频
		var urls []string
		if err := json.Unmarshal([]byte(task.Content), &urls); err != nil {
			return err
		}
		for i, url := range urls {
			video := tgbotapi.NewVideo(chatID, tgbotapi.FileURL(url))
			if i == 0 && markup != nil {
				video.ReplyMarkup = markup
			}
			if _, err := botAPI.Send(video); err != nil {
				return err
			}
		}
		return nil

	default:
		return fmt.Errorf("不支持的 TaskSendType: %d", task.TaskSendType)
	}
}

func InitBotTaskManager() {
	BotTaskManager = &TaskManager{
		tasks: make(map[uint]*TaskRunner),
	}

	var tasks []bot.BotTask
	// 只加载未结束且 status=1 的任务
	err := global.GVA_DB.
		Where("status = ? AND (stop_time IS NULL OR stop_time > ?)", 1, time.Now()).
		Find(&tasks).Error
	if err != nil {
		global.GVA_LOG.Error("加载任务失败", zap.Error(err))
		return
	}

	for _, task := range tasks {
		t := task
		BotTaskManager.StartTask(&t)
	}

	global.GVA_LOG.Info("InitBotTaskManager", zap.Int("count", len(tasks)))
}
