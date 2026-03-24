package ledger2

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger2"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

// 查看操作人
type ListOprUserHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	msg         botapi.Update
	OnlyAdminAware
}

func (l *ListOprUserHandler) Match(botModel bot.Bot, update botapi.Update) bool {

	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)

	if text != "查看操作人" {
		return false
	}

	l.botModel = botModel
	l.chatGroupID = update.Message.Chat.ID
	l.msg = update

	return true
}

func (l *ListOprUserHandler) Handle() error {

	var permission ledger2.LedgerPermission

	err := global.GVA_MYSQL.Where(
		"bot_id = ? AND chat_group_id = ?",
		l.botModel.BotID,
		l.chatGroupID,
	).First(&permission).Error

	if err != nil {
		global.GVA_LOG.Error("ListOprUserHandler Handle",
			zap.Any("bot", l.botModel),
			zap.Error(err),
		)
		return l.reply("获取失败")
	}

	users := splitUsers(permission.OprUsers)

	if len(users) == 0 {
		return l.reply("当前没有设置操作人")
	}

	content := fmt.Sprintf(
		"允许操作人（%d人）：\n%s",
		len(users),
		strings.Join(users, "\n"),
	)

	return l.reply(content)
}

func (l *ListOprUserHandler) reply(text string) error {
	botSender, err := botapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return err
	}

	msg := botapi.NewMessage(l.chatGroupID, text)
	msg.ReplyToMessageID = l.msg.Message.MessageID

	_, err = botSender.Send(msg)
	return err
}
