package bot

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/utils/bot_handler"
	"go.uber.org/zap"
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
				global.GVA_MYSQL.Model(tr.Task).Update("status", 2)
				global.GVA_LOG.Info("TaskRunner StopTask due to StopTime", zap.Int("taskID", int(tr.Task.ID)))
				return
			}

			// 更新 NextSendTime，保证它在未来
			for !tr.Task.NextSendTime.After(now) {
				tr.Task.NextSendTime = tr.Task.NextSendTime.Add(time.Duration(tr.Task.SendInterval) * time.Minute)
				// 如果下次发送时间已经超过 StopTime，则任务结束
				if !tr.Task.StopTime.IsZero() && tr.Task.NextSendTime.After(tr.Task.StopTime) {
					tr.Task.Status = 2
					global.GVA_MYSQL.Model(tr.Task).Update("status", 2)
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
					global.GVA_MYSQL.Model(tr.Task).Updates(map[string]interface{}{
						"pre_send_time":  tr.Task.PreSendTime,
						"next_send_time": tr.Task.NextSendTime,
						"status":         2,
					})
					global.GVA_LOG.Info("TaskRunner StopTask due to StopTime after send", zap.Int("taskID", int(tr.Task.ID)))
					return
				}

				global.GVA_MYSQL.Model(tr.Task).Updates(map[string]interface{}{
					"pre_send_time":  tr.Task.PreSendTime,
					"next_send_time": tr.Task.NextSendTime,
				})
			}
		}
	}
}

func SendTelegramMessage(chatID int64, task *bot.BotTask) error {
	// 1. 获取 bot
	botModel, has, err := dao.BotDao.FromBotID(global.GVA_MYSQL, int(task.BotID))
	if !has || err != nil {
		if !has {
			return fmt.Errorf("bot %d not found", task.BotID)
		}
		return err
	}

	botAPI, err := tgbotapi.NewBotAPI(botModel.Token)
	if err != nil {
		return err
	}

	var markup tgbotapi.InlineKeyboardMarkup
	var hasButton bool

	if len(task.ExtrendButton) > 0 {
		var btns [][]bot.ButtonItem
		if err := json.Unmarshal(task.ExtrendButton, &btns); err == nil {

			var rows [][]tgbotapi.InlineKeyboardButton
			for _, row := range btns {
				var r []tgbotapi.InlineKeyboardButton
				for _, b := range row {
					r = append(r, tgbotapi.NewInlineKeyboardButtonURL(b.Name, b.URL))
				}
				if len(r) > 0 {
					rows = append(rows, r)
				}
			}

			if len(rows) > 0 {
				markup = tgbotapi.NewInlineKeyboardMarkup(rows...)
				hasButton = true
			}
		}
	}

	spew.Dump(">>>>>>>>>>>>>>>>>markup", markup)
	// ==================== 分类型发送 ====================
	switch task.TaskSendType {

	case 0: // 仅按钮
		if hasButton {
			msg := tgbotapi.NewMessage(chatID, "请选择：")
			msg.ReplyMarkup = markup
			_, err = botAPI.Send(msg)
		}
		return err

	case 1: // 富文本
		return bot_handler.HandleTexWithMarup(chatID, botModel.Token, task.Content, markup)

	case 2: // 文本
		msg := tgbotapi.NewMessage(chatID, task.Content)
		if hasButton {
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
			if i == 0 && hasButton {
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
			if i == 0 && hasButton {
				video.ReplyMarkup = markup
			}
			if _, err := botAPI.Send(video); err != nil {
				return err
			}
		}
		return nil
	}

	return fmt.Errorf("不支持的 TaskSendType: %d", task.TaskSendType)
}

func InitBotTaskManager() {
	BotTaskManager = &TaskManager{
		tasks: make(map[uint]*TaskRunner),
	}
	var tasks []bot.BotTask
	err := global.GVA_MYSQL.
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
