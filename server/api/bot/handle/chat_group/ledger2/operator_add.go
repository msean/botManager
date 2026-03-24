package ledger2

import (
	"errors"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger2"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

type AddOprUserHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	userList    []string
	OnlyAdminAware
}

func (l *AddOprUserHandler) Match(botModel bot.Bot, update botapi.Update) bool {

	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)

	if !strings.HasPrefix(text, "设置操作人") {
		return false
	}

	users := parseAtUsers(text)
	global.GVA_LOG.Debug("AddOprUserHandler", zap.Any("users", users))
	if len(users) == 0 {
		return false
	}

	l.userList = users
	l.botModel = botModel
	l.chatGroupID = update.Message.Chat.ID

	return true
}

func (l *AddOprUserHandler) Handle() error {
	var permission ledger2.LedgerPermission

	err := global.GVA_MYSQL.Where(
		"bot_id = ? AND chat_group_id = ?",
		l.botModel.BotID,
		l.chatGroupID,
	).First(&permission).Error

	// ✅ 如果是“没找到”，就创建
	if errors.Is(err, gorm.ErrRecordNotFound) {
		permission = ledger2.LedgerPermission{
			BotID:       l.botModel.BotID,
			ChatGroupID: l.chatGroupID,
			OprUsers:    "",
		}
	} else if err != nil {
		// ❗其他错误才返回
		global.GVA_LOG.Error("AddOprUserHandler Handle",
			zap.Any("bot", l.botModel),
			zap.Error(err),
		)
		return err
	}

	// ===== 原有逻辑 =====
	existing := splitUsers(permission.OprUsers)

	userMap := make(map[string]struct{})

	for _, u := range existing {
		userMap[u] = struct{}{}
	}
	for _, u := range l.userList {
		userMap[u] = struct{}{}
	}

	var finalUsers []string
	for u := range userMap {
		finalUsers = append(finalUsers, u)
	}

	permission.OprUsers = strings.Join(finalUsers, ",")

	// ✅ 关键：用 Save（有ID就更新，没ID就插入）
	if err = global.GVA_MYSQL.Save(&permission).Error; err != nil {
		return err
	}

	return l.reply("操作成功")
}

func (l *AddOprUserHandler) reply(text string) error {
	botSender, err := botapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return err
	}
	msg := botapi.NewMessage(l.chatGroupID, text)
	_, err = botSender.Send(msg)
	return err
}

func parseAtUsers(text string) []string {

	parts := strings.Fields(text)

	var users []string

	for _, p := range parts {

		if strings.HasPrefix(p, "@") {
			u := strings.TrimPrefix(p, "@")

			u = strings.TrimSpace(u)
			u = strings.Trim(u, ",， ")

			if u != "" {
				users = append(users, u)
			}
		}
	}

	return users
}

func splitUsers(s string) []string {
	if s == "" {
		return []string{}
	}
	arr := strings.Split(s, ",")

	var res []string
	for _, v := range arr {
		v = strings.TrimSpace(v)
		if v != "" {
			res = append(res, v)
		}
	}
	return res
}
