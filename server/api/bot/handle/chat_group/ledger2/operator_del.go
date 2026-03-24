package ledger2

import (
	"strings"

	"go.uber.org/zap"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger2"
	"github.com/msean/botmanager/server/service/cache"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

// 删除操作人
type DelOprUserHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	userList    []string
	OnlyAdminAware
}

func (l *DelOprUserHandler) Match(botModel bot.Bot, update botapi.Update) bool {

	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)

	if !strings.HasPrefix(text, "删除操作人") {
		return false
	}

	users := parseAtUsers(text)
	if len(users) == 0 {
		return false
	}

	l.userList = users
	l.botModel = botModel
	l.chatGroupID = update.Message.Chat.ID

	return true
}

func (l *DelOprUserHandler) Handle() error {

	var permission ledger2.LedgerPermission

	err := global.GVA_MYSQL.Where(
		"bot_id = ? AND chat_group_id = ?",
		l.botModel.BotID,
		l.chatGroupID,
	).First(&permission).Error

	if err != nil {
		global.GVA_LOG.Error("DelOprUserHandler Handle",
			zap.Any("bot", l.botModel),
			zap.Error(err),
		)
		return err
	}

	// 👉 原有用户
	existing := splitUsers(permission.OprUsers)

	delMap := make(map[string]struct{})
	for _, u := range l.userList {
		delMap[u] = struct{}{}
	}

	var finalUsers []string
	for _, u := range existing {
		if _, ok := delMap[u]; !ok {
			finalUsers = append(finalUsers, u)
		}
	}

	permission.OprUsers = strings.Join(finalUsers, ",")

	if err = global.GVA_MYSQL.Model(&permission).
		Update("opr_users", permission.OprUsers).Error; err != nil {
		return err
	}

	_ = cache.NewLedgerPermissionCache(
		l.botModel.BotID,
		l.chatGroupID,
	).Release()

	return l.reply("删除成功")
}

func (l *DelOprUserHandler) reply(text string) error {
	botSender, err := botapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return err
	}
	msg := botapi.NewMessage(l.chatGroupID, text)
	_, err = botSender.Send(msg)
	return err
}
