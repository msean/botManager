package tg_auto_helper

import (
	"errors"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/tg_auto_helper"
	"github.com/msean/botmanager/server/model/tg_auto_helper/request"
	tgAutoHelperReq "github.com/msean/botmanager/server/model/tg_auto_helper/request"
	"github.com/msean/botmanager/server/service/tg_auto_helper/collect"
	"gorm.io/gorm"
)

type CollectGroupTaskService struct{}

func (s *CollectGroupTaskService) Create(
	task *tg_auto_helper.CollectGroupTask,
) error {

	var tgUser tg_auto_helper.TgUser

	err := global.GVA_MYSQL.Transaction(func(tx *gorm.DB) error {

		err := tx.
			Where("status = ?", 3). // 已登录
			Where("id NOT IN (?)",
				tx.Model(&tg_auto_helper.CollectGroupTask{}).
					Select("tg_user_id").
					Where("status = ?", 1),
			).
			First(&tgUser).Error

		if err != nil {
			return errors.New("没有可用的TG账号执行任务")
		}

		task.TgUserID = int(tgUser.ID)
		task.Status = 1
		task.CurrentPage = 1
		// 默认搜索soso
		if task.SourceBotName == "" {
			task.SourceBotName = "soso"
		}

		if err := tx.Create(task).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	go func() {

		t := collect.CollectGroupTask{
			CollectGroupTask: *task,
		}

		t.Run()
	}()

	return nil
}

func (s *CollectGroupTaskService) Delete(ID string) error {

	var task tg_auto_helper.CollectGroupTask

	if err := global.GVA_MYSQL.Where("id = ?", ID).First(&task).Error; err != nil {
		return err
	}

	// 停止任务
	collect.TaskManager.Stop(task.ID)

	// 删除任务
	return global.GVA_MYSQL.Delete(&tg_auto_helper.CollectGroupTask{}, "id = ?", ID).Error
}

// func (s *CollectGroupTaskService) UpdateCollectGroupTask(
// 	ctx context.Context,
// 	task tg_auto_helper.CollectGroupTask,
// ) error {

// 	return global.GVA_MYSQL.Model(&tg_auto_helper.CollectGroupTask{}).
// 		Where("id = ?", task.ID).
// 		Updates(map[string]any{
// 			"search_text":   task.SearchText,
// 			"members_lower": task.MemberLower,
// 			"status":        task.Status,
// 		}).Error
// }

// func (s *CollectGroupTaskService) GetCollectGroupTask(
// 	ctx context.Context,
// 	ID int,
// ) (task tg_auto_helper.CollectGroupTask, err error) {

// 	err = global.GVA_MYSQL.
// 		Where("id = ?", ID).
// 		First(&task).Error

// 	return
// }

func (s *CollectGroupTaskService) List(
	info tgAutoHelperReq.CollectGroupTaskSearch,
) (list []tg_auto_helper.CollectGroupTask, total int64, err error) {

	db := global.GVA_MYSQL.Model(&tg_auto_helper.CollectGroupTask{})

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	if len(info.CreatedAtRange) == 2 {
		db = db.Where(
			"created_at BETWEEN ? AND ?",
			info.CreatedAtRange[0],
			info.CreatedAtRange[1],
		)
	}

	if info.Status != 0 {
		db = db.Where("status = ?", info.Status)
	}

	if info.SearchText != "" {
		db = db.Where("search_text LIKE ?", "%"+info.SearchText+"%")
	}

	if err = db.Count(&total).Error; err != nil {
		return
	}

	err = db.
		Limit(limit).
		Offset(offset).
		Order("id desc").
		Find(&list).Error

	return
}

// func (s *CollectGroupTaskService) StopCollectGroupTask(
// 	ctx context.Context,
// 	ID int,
// ) error {

//		return global.GVA_MYSQL.Model(&tg_auto_helper.CollectGroupTask{}).
//			Where("id = ?", ID).
//			Update("status", 3).Error
//	}
func (s *CollectGroupTaskService) ListCollectGroupInfo(info request.CollectGroupInfoSearch) (list []tg_auto_helper.CollectGroupInfo, total int64, err error) {

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_MYSQL.Model(&tg_auto_helper.CollectGroupInfo{})

	if info.TaskID != 0 {
		db = db.Where("task_id = ?", info.TaskID)
	}

	if info.SearchText != "" {
		db = db.Where("title LIKE ?", "%"+info.SearchText+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.
		Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&list).Error

	return
}
