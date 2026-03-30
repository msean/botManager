package msgmass

import (
	"strings"

	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

type HelpHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	ShouldPermissionAwareWithOutAdmin
}

func (h *HelpHandler) Match(botModel bot.Bot, update botapi.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)

	if text != "help" && text != "帮助" {
		return false
	}

	h.botModel = botModel
	h.chatGroupID = update.Message.Chat.ID

	return true
}

func (h *HelpHandler) Handle() error {

	text := `
群发机器人帮助菜单
1、增加操作人
命令：增加操作人 xx  
示例：增加操作人 bot  
说明：xx是群用户，将 bot 设置为可操作账单人（仅管理员）

2、艾特设置 
命令：艾特设置 xx1,xx2
示例：
艾特设置 bot1
艾特设置 bot1,bot2 
说明：多个艾特用户用,分割

3、查看群艾特用户
命令：艾特查看 
`
	return h.reply(text)
}

func (h *HelpHandler) reply(text string) error {
	botSender, err := botapi.NewBotAPI(h.botModel.Token)
	if err != nil {
		return err
	}

	msg := botapi.NewMessage(h.chatGroupID, text)
	_, err = botSender.Send(msg)
	return err
}
