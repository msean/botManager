package ledger

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
1、增加操作人
命令：增加操作人 xx  
示例：增加操作人 bot  
说明：xx是群用户，将 bot 设置为可操作账单人（仅管理员）

2、清空今日账单 
命令：清空今日账单 
说明：删除今天所有账单记录

3、费率设置
命令：设置费率 xx  
示例：设置费率 2.3 
说明：xx是要设置的费率 

4、查看费率  
命令：查看费率

5、入款操作
命令：+xx1 xx2
示例：
+800  
+800 张三  
说明：xx1是入款金额， xx2是入款账户姓名

6、下发操作
命令：下发xx1 xx2
示例：
下发1000  
下发1000 张三  
说明：xx1是下发金额， xx2是下发账户姓名

7、查看汇率
命令：汇率/hl
说明：来源于okex
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
