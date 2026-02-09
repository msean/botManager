package msgmass

import (
	"strings"

	"go.uber.org/zap"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

type AddOprUserHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	userID      string
	NotPermissionAwareWithAdmin
}

func (l *AddOprUserHandler) Match(botModel bot.Bot, update botapi.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)

	if !strings.HasPrefix(text, "增加操作人") {
		return false
	}

	parts := strings.Fields(text)
	if len(parts) != 2 {
		return false
	}

	l.userID = parts[1]
	l.botModel = botModel
	l.chatGroupID = update.Message.Chat.ID

	return true
}

func (l *AddOprUserHandler) Handle() error {
	var permission bot.BotMassMsgPermission

	err := global.GVA_MYSQL.Where(
		"bot_id = ? AND chat_group_id = ?",
		l.botModel.BotID,
		l.chatGroupID,
	).First(&permission).Error

	if err != nil {
		global.GVA_LOG.Error("AddOprUserHandler Handle", zap.Any("bot", l.botModel), zap.Error(err))
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		return err
		// 不存在就创建
		// permission = ledger.LedgerPermission{
		// 	BotID:       l.botModel.BotID,
		// 	ChatGroupID: l.chatGroupID,
		// 	OprUsers:    l.userID,
		// }
		// if err = global.GVA_MYSQL.Create(&permission).Error; err != nil {
		// 	return err
		// 	// }
		// } else {
		// 	return err
		// }
	} else {
		// 已存在，追加 opr user（去重）
		if !hasOprUser(permission.OprUsers, l.userID) {
			if permission.OprUsers == "" {
				permission.OprUsers = l.userID
			} else {
				permission.OprUsers = permission.OprUsers + "," + l.userID
			}
			if err = global.GVA_MYSQL.Model(&permission).
				Update("opr_users", permission.OprUsers).Error; err != nil {
				return err
			}
		}
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

func hasOprUser(oprUsers string, userID string) bool {
	if oprUsers == "" || userID == "" {
		return false
	}

	users := strings.Split(oprUsers, ",")
	for _, u := range users {
		if strings.TrimSpace(u) == userID {
			return true
		}
	}
	return false
}
