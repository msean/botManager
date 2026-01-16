package bot

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
	"github.com/msean/botmanager/server/utils"
	"go.uber.org/zap"
)

type BotTaskService struct{}

// CreateBotTask 创建任务列表记录
// Author [yourname](https://github.com/yourname)
func (taskService *BotTaskService) CreateBotTask(ctx context.Context, task *bot.BotTask) (err error) {
	task.ID = 0
	task.PreSendTime = nil
	if task.NextSendTimeStr == "" {
		err = fmt.Errorf("下次发送时间不能为空")
		return
	}
	if task.StopTimeText == "" {
		err = fmt.Errorf("下次发送时间不能为空")
		return
	}
	layout := "2006-01-02 15:04:05"
	if task.NextSendTime, err = time.ParseInLocation(layout, task.NextSendTimeStr, time.Local); err != nil {
		return
	}
	if task.StopTime, err = time.ParseInLocation(layout, task.StopTimeText, time.Local); err != nil {
		return
	}
	if task.StopTime.Before(task.NextSendTime) || task.StopTime.Equal(task.NextSendTime) {
		err = fmt.Errorf("发送时间大于或者等于结束时间")
		return
	}
	if err = global.GVA_MYSQL.Create(task).Error; err != nil {
		return
	}

	BotTaskManager.StartTask(task)
	return err
}

// DeleteBotTask 删除任务列表记录
// Author [yourname](https://github.com/yourname)
func (taskService *BotTaskService) DeleteBotTask(ctx context.Context, ID string) (err error) {
	var idInt int
	if idInt, err = strconv.Atoi(ID); err != nil {
		return
	}
	if err = global.GVA_MYSQL.Delete(&bot.BotTask{}, "id = ?", ID).Error; err != nil {
		return
	}
	BotTaskManager.StopTask(uint(idInt))
	return err
}

// DeleteBotTaskByIds 批量删除任务列表记录
// Author [yourname](https://github.com/yourname)
func (taskService *BotTaskService) DeleteBotTaskByIds(ctx context.Context, IDs []string) (err error) {
	ids := utils.StringsToIntsIgnoreError(IDs)
	if err = global.GVA_MYSQL.Delete(&[]bot.BotTask{}, "id in ?", IDs).Error; err != nil {
		return
	}
	for _, id := range ids {
		BotTaskManager.StopTask(uint(id))
	}
	return err
}

// UpdateBotTask 更新任务列表记录
// Author [yourname](https://github.com/yourname)
func (taskService *BotTaskService) UpdateBotTask(ctx context.Context, task *bot.BotTask) (err error) {
	if task.NextSendTimeStr == "" {
		err = fmt.Errorf("下次发送时间不能为空")
		return
	}
	if task.StopTimeText == "" {
		err = fmt.Errorf("下次发送时间不能为空")
		return
	}
	layout := "2006-01-02 15:04:05"
	if task.NextSendTime, err = time.ParseInLocation(layout, task.NextSendTimeStr, time.Local); err != nil {
		return
	}
	if task.StopTime, err = time.ParseInLocation(layout, task.StopTimeText, time.Local); err != nil {
		return
	}
	if err = global.GVA_MYSQL.Model(&bot.BotTask{}).Where("id = ?", task.ID).Updates(&task).Error; err != nil {
		return
	}
	BotTaskManager.ReloadTask(task)
	return err
}

// GetBotTask 根据ID获取任务列表记录
// Author [yourname](https://github.com/yourname)
func (taskService *BotTaskService) GetBotTask(ctx context.Context, ID string) (task *bot.BotTask, err error) {
	if err = global.GVA_MYSQL.Where("id = ?", ID).First(&task).Error; err != nil {
		return
	}
	var botModel bot.Bot
	var has bool
	if botModel, has, err = dao.BotDao.FromBotID(global.GVA_MYSQL, int(task.BotID)); !has || err != nil {
		if !has {
			err = fmt.Errorf("没有找到改机器人")
		}
		global.GVA_LOG.Error("taskService GetBotTask", zap.Any("botID", task.BotID), zap.Bool("has", has), zap.Error(err))
		return
	}
	task.BotName = botModel.Name

	switch task.GroupType {
	case constant.GroupTypeChat:
		var botChatGroupModel bot.BotChatGroup
		if botChatGroupModel, has, err = dao.BotChatGroupDao.FromID(global.GVA_MYSQL, int(task.GroupID)); !has || err != nil {
			if !has {
				err = fmt.Errorf("没有找到改机器人")
			}
			global.GVA_LOG.Error("taskService GetBotTask", zap.Any("GroupID", task.GroupID), zap.Bool("has", has), zap.Error(err))
			return
		}
		task.GroupName = botChatGroupModel.ChatGroupName
	case constant.GroupTypeChannel:
		var botChannel bot.BotChannel
		if botChannel, has, err = dao.BotChannelDao.FromBotID(global.GVA_MYSQL, int(task.GroupID)); !has || err != nil {
			if !has {
				err = fmt.Errorf("没有找到改渠道")
			}
			global.GVA_LOG.Error("taskService GetBotTask", zap.Any("", task.GroupID), zap.Bool("has", has), zap.Error(err))
			return
		}
		task.GroupName = botChannel.ChannelName
	}

	task.NextSendTimeStr = task.NextSendTime.Format("2006-01-02 15:04:05")
	task.StopTimeText = task.StopTime.Format("2006-01-02 15:04:05")
	return
}

// GetBotTaskInfoList 分页获取任务列表记录
// Author [yourname](https://github.com/yourname)
func (taskService *BotTaskService) GetBotTaskInfoList(ctx context.Context, info botReq.BotTaskSearch) (list []*bot.BotTask, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_MYSQL.Model(&bot.BotTask{})
	var tasks []*bot.BotTask
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	if err = db.Find(&tasks).Error; err != nil {
		return
	}
	var botList []int64
	var chatGroupList []int64
	var channelList []int64
	for _, object := range tasks {
		botList = append(botList, object.BotID)
		if object.GroupType == constant.GroupTypeChat {
			chatGroupList = append(chatGroupList, object.GroupID)
		} else {
			channelList = append(channelList, object.GroupID)
		}
	}

	var botMapper map[int64]bot.Bot
	var chatGroupMapper map[int64]bot.BotChatGroup
	if botMapper, err = dao.BotDao.MappByIDList(global.GVA_MYSQL, botList); err != nil {
		return
	}
	var channelMapper map[int64]bot.BotChannel
	if botMapper, err = dao.BotDao.MappByIDList(global.GVA_MYSQL, botList); err != nil {
		return
	}

	if chatGroupMapper, err = dao.BotChatGroupDao.MappByChatGroupIDList(global.GVA_MYSQL, chatGroupList); err != nil {
		return
	}
	if channelMapper, err = dao.BotChannelDao.MappByChannelIDList(global.GVA_MYSQL, channelList); err != nil {
		return
	}

	for _, object := range tasks {
		object.BotName = botMapper[object.BotID].Name
		if object.GroupType == constant.GroupTypeChat {
			object.GroupName = chatGroupMapper[object.GroupID].ChatGroupName
		} else {
			object.GroupName = channelMapper[object.GroupID].ChannelName
		}
	}

	return tasks, total, err
}

func (taskService *BotTaskService) GetBotTaskPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
