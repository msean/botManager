package ledger2

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
普通记账机器人帮助菜单
1、设置操作人
命令：设置操作人 @xx   
示例：
设置操作人 @bot @bot2
说明：xx是群用户，将 bot、bot2 设置为可操作账单人（仅管理员）

2、删除操作人 
命令：删除操作人 @xx
示例：
删除操作人 @bot @bot2
说明：xx是群用户，将 bot、bot2 设置为可操作账单人（仅管理员）

3、查看操作人
命令：删除操作人  

4、开始/重新开始记账  
命令：开始

5、结束当天记账
命令：下课

6、清除当天账单重新记账
命令：上课
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
