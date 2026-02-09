package msgmass

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

type BotMassMsgHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	members     string
	NeedAdminAware
}

func (h *BotMassMsgHandler) Match(botModel bot.Bot, update botapi.Update) (match bool) {
	text := strings.TrimSpace(update.Message.Text)

	// 必须以「艾特设置」开头
	if !strings.HasPrefix(text, "艾特设置") {
		return false
	}

	parts := strings.Fields(text)
	if len(parts) < 2 {
		return false
	}

	// 获取成员列表，空格分隔
	h.members = strings.Join(parts[1:], ",")
	h.botModel = botModel
	h.chatGroupID = update.Message.Chat.ID

	return true
}

func (h *BotMassMsgHandler) Handle() (err error) {
	if h.botModel.BotID == 0 || h.chatGroupID == 0 {
		return errors.New("bot or chat group invalid")
	}

	var record bot.BotMsgMass

	// 查询是否存在
	err = global.GVA_MYSQL.Where(
		"bot_id = ? AND chat_group_id = ?",
		h.botModel.BotID,
		h.chatGroupID,
	).First(&record).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在就创建
			record = bot.BotMsgMass{
				BotID:       h.botModel.BotID,
				ChatGroupID: h.chatGroupID,
				Members:     h.members,
			}
			if err = global.GVA_MYSQL.Create(&record).Error; err != nil {
				return
			}
			h.reply("设置成功")
			return nil
		}
		return err
	}

	// 已存在就更新 Members
	if err = global.GVA_MYSQL.Model(&record).Updates(map[string]interface{}{
		"members": h.members,
	}).Error; err != nil {
		return
	}

	h.reply("设置成功")
	return nil
}

func (h *BotMassMsgHandler) reply(text string) (err error) {
	botSender, err := botapi.NewBotAPI(h.botModel.Token)
	if err != nil {
		return
	}
	msg := botapi.NewMessage(h.chatGroupID, text)
	_, err = botSender.Send(msg)
	return
}
