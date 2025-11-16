package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
}

func (tm *TaskManager) StopTask(taskID uint) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if runner, ok := tm.tasks[taskID]; ok {
		close(runner.StopChan)
		delete(tm.tasks, taskID)
	}
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
}

func (tr *TaskRunner) Run(sendFunc TgSendFunc) {
	global.GVA_LOG.Info("TaskRunner", zap.Any("next_send_time", tr.Task.NextSendTime))
	for {
		select {
		case <-tr.StopChan:
			return

		case <-time.After(time.Until(tr.Task.NextSendTime)):
			fmt.Printf("开始执行任务 %d\n", tr.Task.ID)

			err := sendFunc(tr.Task.ChatGroupID, tr.Task)
			if err != nil {
				global.GVA_LOG.Error("SendTelegramMessage", zap.Any("tr.Task", tr.Task), zap.Error(err))
			}

			// 记录发送时间
			now := time.Now()
			global.GVA_DB.Model(tr.Task).Updates(map[string]interface{}{
				"pre_send_time":  now,
				"next_send_time": now.Add(time.Duration(tr.Task.SendInterval) * time.Minute),
			})

			tr.Task.NextSendTime = tr.Task.NextSendTime.Add(time.Duration(tr.Task.SendInterval) * time.Minute)
		}
	}
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
		msg := tgbotapi.NewMessage(chatID, task.Content)
		msg.ParseMode = "HTML"
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
}
