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

			for !tr.Task.NextSendTime.After(now) {
				tr.Task.NextSendTime = tr.Task.NextSendTime.Add(time.Duration(tr.Task.SendInterval) * time.Minute)
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

		caption := text
		if len(caption) > 1024 {
			caption = caption[:1020] + "..."
		}

		var markup *tgbotapi.InlineKeyboardMarkup
		if len(task.ExtrendButton) > 0 {
			markup = buildMarkupFromExtrend(task.ExtrendButton)
		}

		if len(imgs) > 0 {
			first := imgs[0]

			resp, err := http.Get(first)
			if err != nil {
				global.GVA_LOG.Error("download image failed", zap.Error(err), zap.String("url", first))
				return err
			}
			defer resp.Body.Close()

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				global.GVA_LOG.Error("read downloaded image failed", zap.Error(err), zap.String("url", first))
				return err
			}

			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{
				Name:  filepath.Base(first),
				Bytes: data,
			})
			photo.Caption = caption
			photo.ParseMode = tgbotapi.ModeHTML
			if markup != nil {
				photo.ReplyMarkup = markup
			}

			if _, err = botAPI.Send(photo); err != nil {
				global.GVA_LOG.Error("send photo with caption failed", zap.Error(err), zap.String("url", first))
				return err
			}
			return nil
		}

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
	err := global.GVA_DB.Where("status = ?", 1).Find(&tasks).Error
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
