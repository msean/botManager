package collect

import (
	"context"
	"sync"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/tg_auto_helper"
	"go.uber.org/zap"
)

var TaskManager *taskManager

type taskManager struct {
	mu    sync.Mutex
	tasks map[uint]context.CancelFunc
}

func (m *taskManager) Add(taskID uint, cancel context.CancelFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tasks[taskID] = cancel
}

func (m *taskManager) Stop(taskID uint) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cancel, ok := m.tasks[taskID]; ok {
		cancel()
		delete(m.tasks, taskID)
	}
}

func (m *taskManager) Remove(taskID uint) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tasks, taskID)
}

func InitCollectTaskManager() {
	TaskManager = &taskManager{
		tasks: make(map[uint]context.CancelFunc),
	}
	go RecoverRunningTasks()
}

func RecoverRunningTasks() {
	time.Sleep(10 * time.Second)
	var tasks []tg_auto_helper.CollectGroupTask

	err := global.GVA_MYSQL.
		Where("status = ?", 1).
		Find(&tasks).Error

	if err != nil {
		global.GVA_LOG.Error("recover task error", zap.Error(err))
		return
	}

	global.GVA_LOG.Debug("RecoverRunningTasks", zap.Any("task", tasks))
	for _, task := range tasks {
		t := CollectGroupTask{
			CollectGroupTask: task,
		}
		go func(tt CollectGroupTask) {
			tt.Run()
		}(t)
	}
}
