package dao

import (
	"github.com/msean/botmanager/server/model/system"
	"gorm.io/gorm"
)

type sysDao struct{}

func newsysDao() *sysDao {
	return &sysDao{}
}

// func (dao *sysParamsDao) FromKey(db *gorm.DB, key string, defaultValue string) (param system.SysParams, err error) {
// 	err = db.Where("`key` = ?", key).First(&param).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			// 没找到则创建新记录
// 			param = system.SysParams{
// 				Key:   key,
// 				Value: defaultValue,
// 				Name:  key,
// 				Desc:  "系统自动初始化默认参数",
// 			}
// 			if createErr := db.Create(&param).Error; createErr != nil {
// 				err = createErr
// 				return
// 			}
// 			err = nil
// 		}
// 	}
// 	return
// }

func (dao *sysDao) NameMapperFromIDList(db *gorm.DB, userID []int64) (mapper map[int64]string, err error) {
	var models []system.SysUser
	mapper = make(map[int64]string)
	if err = db.Find(&models, "id in (?)", userID).Error; err != nil {
		return
	}
	for _, model := range models {
		mapper[int64(model.ID)] = model.Username
	}
	return
}

func (dao *sysDao) AllUser(db *gorm.DB) (models []system.SysUser, err error) {
	err = db.Find(&models).Error
	return
}
