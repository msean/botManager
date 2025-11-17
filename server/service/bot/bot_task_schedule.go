package bot

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"go.uber.org/zap"
)

type (
	TaskManager struct {
		mu    sync.Mutex
		tasks map[uint]*TaskRunner
	}
	TgSendFunc func(int64, *bot.BotTask) error
)

var BotTaskManager *TaskManager

type TaskRunner struct {
	Task     *bot.BotTask
	StopChan chan struct{}
}

func (tm *TaskManager) StartTask(task *bot.BotTask, sendFunc TgSendFunc) {
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

	go runner.Run(sendFunc)
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

func (tm *TaskManager) ReloadTask(task *bot.BotTask, sendFunc TgSendFunc) {
	tm.StopTask(task.ID)
	tm.StartTask(task, sendFunc)
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

func (tr *TaskRunner) Run(sendFunc TgSendFunc) {
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
				err := sendFunc(tr.Task.ChatGroupID, tr.Task)
				if err != nil {
					global.GVA_LOG.Error("发送失败", zap.Error(err))
					continue
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

	// 3️⃣ 转义 HTML 实体
	htmlStr = html.UnescapeString(htmlStr)

	// 4️⃣ 合并多余换行
	htmlStr = regexp.MustCompile(`\n{2,}`).ReplaceAllString(htmlStr, "\n\n")
	htmlStr = strings.TrimSpace(htmlStr)
	return htmlStr
}

func SendTelegramMessage(chatID int64, task *bot.BotTask) (err error) {
	var botModel bot.Bot
	var has bool
	if botModel, has, err = dao.BotDao.FromBotID(global.GVA_DB, int(task.BotID)); !has || err != nil {
		if !has {
			err = fmt.Errorf("bot %d not found", task.BotID)
		}
		return
	}
	botAPI, err := tgbotapi.NewBotAPI(botModel.Token)
	if err != nil {
		return err
	}

	var markup *tgbotapi.InlineKeyboardMarkup
	if len(task.ExtrendButton) > 0 {
		var btns []bot.ButtonItem
		_ = json.Unmarshal(task.ExtrendButton, &btns)

		if len(btns) > 0 {
			var row []tgbotapi.InlineKeyboardButton
			for _, b := range btns {
				row = append(row, tgbotapi.NewInlineKeyboardButtonURL(b.Name, b.URL))
			}
			m := tgbotapi.NewInlineKeyboardMarkup(row)
			markup = &m
		}
	}

	switch task.TaskSendType {

	case 1:
		imgs, text := extractImgsAndText(task.Content)

		// 清理 HTML
		caption := CleanHTMLForTelegram(text)
		if len(caption) > 1024 {
			caption = caption[:1020] + "..."
		}

		// 发送单图 + caption
		if len(imgs) > 0 {
			first := imgs[0]
			resp, err := http.Get(first)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			data, _ := ioutil.ReadAll(resp.Body)

			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{
				Name:  filepath.Base(first),
				Bytes: data,
			})
			photo.Caption = caption
			photo.ParseMode = tgbotapi.ModeHTML
			if markup != nil {
				photo.ReplyMarkup = markup
			}
			_, err = botAPI.Send(photo)
			return err
		}

		// 没有图片则直接发送文本
		msg := tgbotapi.NewMessage(chatID, caption)
		msg.ParseMode = tgbotapi.ModeHTML
		if markup != nil {
			msg.ReplyMarkup = markup
		}
		_, err = botAPI.Send(msg)
		return err
	case 2: // 普通文本
		msg := tgbotapi.NewMessage(chatID, task.Content)
		if markup != nil {
			msg.ReplyMarkup = markup
		}
		_, err = botAPI.Send(msg)
		return err

	case 3: // 图片（多图）
		var urls []string
		_ = json.Unmarshal([]byte(task.Content), &urls)
		global.GVA_LOG.Info("SendTelegramMessage", zap.Any("urls", urls))

		for _, url := range urls {
			data, _ := ioutil.ReadFile(url)
			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{
				Name:  "file.jpg",
				Bytes: data,
			})
			photo = tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(url))
			if markup != nil {
				photo.ReplyMarkup = markup
			}
			_, err = botAPI.Send(photo)
			if err != nil {
				return err
			}
		}
		return nil

	case 4: // 视频
		var urls []string
		_ = json.Unmarshal([]byte(task.Content), &urls)

		for _, url := range urls {
			video := tgbotapi.NewVideo(chatID, tgbotapi.FileURL(url))
			if markup != nil {
				video.ReplyMarkup = markup
			}
			_, err = botAPI.Send(video)
			if err != nil {
				return err
			}
		}
		return nil

	default:
		return fmt.Errorf("不支持的 TaskSendType：%d", task.TaskSendType)
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
		BotTaskManager.StartTask(&t, SendTelegramMessage)
	}

	global.GVA_LOG.Info("InitBotTaskManager", zap.Int("count", len(tasks)))
}
