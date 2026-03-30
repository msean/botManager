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
📘 记账系统使用说明

👤 操作人管理
————————————
增加操作人 xx  
👉 示例：增加操作人 bot  
👉 说明：将 bot 设置为可操作账单人（仅管理员）

🧹 账单管理
————————————
清空今日账单  
👉 说明：删除今天所有账单记录（谨慎使用）

📊 费率设置
————————————
设置费率 xx  
👉 示例：设置费率 2.3  

查看费率  
👉 查看当前费率

💰 入款（收入）
————————————
+金额 [姓名]  
👉 示例：
+800  
+800 张三  
👉 说明：姓名可选，不填默认自己

💸 下发（支出）
————————————
下发金额 [姓名]  
👉 示例：
下发1000  
下发1000 张三  
👉 说明：姓名可选
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
