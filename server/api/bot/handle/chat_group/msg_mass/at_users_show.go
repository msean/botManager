package msgmass

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

type BotAtUserShowHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	NeedAdminAware
}

func (h *BotAtUserShowHandler) Match(botModel bot.Bot, update botapi.Update) (match bool) {
	text := strings.TrimSpace(update.Message.Text)

	// 精确匹配「艾特查看」
	if text != "艾特查看" {
		return false
	}

	h.botModel = botModel
	h.chatGroupID = update.Message.Chat.ID
	return true
}

func (h *BotAtUserShowHandler) Handle() (err error) {
	if h.botModel.BotID == 0 || h.chatGroupID == 0 {
		return errors.New("bot or chat group invalid")
	}

	var record bot.BotMsgMass

	err = global.GVA_MYSQL.Where(
		"bot_id = ? AND chat_group_id = ?",
		h.botModel.BotID,
		h.chatGroupID,
	).First(&record).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.reply("当前未设置艾特成员")
			return nil
		}
		return err
	}

	if strings.TrimSpace(record.Members) == "" {
		h.reply("当前未设置艾特成员")
		return nil
	}

	// 展示成员（支持你之前用的空格 / 逗号）
	members := strings.FieldsFunc(record.Members, func(r rune) bool {
		return r == ' ' || r == ',' || r == '，'
	})

	replyText := fmt.Sprintf(
		"当前艾特成员：\n%s",
		strings.Join(members, "\n"),
	)

	h.reply(replyText)
	return nil
}

func (h *BotAtUserShowHandler) reply(text string) (err error) {
	botSender, err := botapi.NewBotAPI(h.botModel.Token)
	if err != nil {
		return
	}
	msg := botapi.NewMessage(h.chatGroupID, text)
	_, err = botSender.Send(msg)
	return
}
