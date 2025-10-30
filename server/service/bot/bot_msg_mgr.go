
package bot

import (
	"context"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
    botReq "github.com/msean/botmanager/server/model/bot/request"
)

type BotMsgMgrService struct {}
// CreateBotMsgMgr 创建机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (bot_msg_mgrService *BotMsgMgrService) CreateBotMsgMgr(ctx context.Context, bot_msg_mgr *bot.BotMsgMgr) (err error) {
	err = global.GVA_DB.Create(bot_msg_mgr).Error
	return err
}

// DeleteBotMsgMgr 删除机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (bot_msg_mgrService *BotMsgMgrService)DeleteBotMsgMgr(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&bot.BotMsgMgr{},"id = ?",ID).Error
	return err
}

// DeleteBotMsgMgrByIds 批量删除机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (bot_msg_mgrService *BotMsgMgrService)DeleteBotMsgMgrByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]bot.BotMsgMgr{},"id in ?",IDs).Error
	return err
}

// UpdateBotMsgMgr 更新机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (bot_msg_mgrService *BotMsgMgrService)UpdateBotMsgMgr(ctx context.Context, bot_msg_mgr bot.BotMsgMgr) (err error) {
	err = global.GVA_DB.Model(&bot.BotMsgMgr{}).Where("id = ?",bot_msg_mgr.ID).Updates(&bot_msg_mgr).Error
	return err
}

// GetBotMsgMgr 根据ID获取机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (bot_msg_mgrService *BotMsgMgrService)GetBotMsgMgr(ctx context.Context, ID string) (bot_msg_mgr bot.BotMsgMgr, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&bot_msg_mgr).Error
	return
}
// GetBotMsgMgrInfoList 分页获取机器人消息管理记录
// Author [yourname](https://github.com/yourname)
func (bot_msg_mgrService *BotMsgMgrService)GetBotMsgMgrInfoList(ctx context.Context, info botReq.BotMsgMgrSearch) (list []bot.BotMsgMgr, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&bot.BotMsgMgr{})
    var bot_msg_mgrs []bot.BotMsgMgr
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.Ban_content != nil && *info.Ban_content != "" {
        db = db.Where("ban_content LIKE ?", "%"+ *info.Ban_content+"%")
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&bot_msg_mgrs).Error
	return  bot_msg_mgrs, total, err
}
func (bot_msg_mgrService *BotMsgMgrService)GetBotMsgMgrPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
