package dao

type sysParamsDao struct{}

func newsysParamsDao() *sysParamsDao {
	return &sysParamsDao{}
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
