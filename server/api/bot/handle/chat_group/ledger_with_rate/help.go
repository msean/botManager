package rate_ledger

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

func (l *HelpHandler) Match(botModel bot.Bot, update botapi.Update) (match bool) {
	text := strings.TrimSpace(update.Message.Text)

	if text != "/help" {
		return false
	}

	return false
}

// 功能：增加操作人 格式：增加操作人xx(可以是用户名，也可以是用户id)
// 功能：清空今日账单 格式：清空今日账单
// 功能：设置汇率 格式：设置汇率 0.1
// 功能：查看汇率 格式：查看汇率
// 功能：设置费率 格式：设置费率 0.1
// 功能：查看费率 格式：查看费率
// 功能：入款    支持格式：1、+1000 张三 2、+1000*7.0 张三  3、+10000/7.0 张三 张三作为入款人备注，可以不用输入
// 功能：下发 支持格式：下发 1000 张三  下发 张三作为下发人备注，可以不用输入
func (l *HelpHandler) Handle() (err error) {

	helpText := `
📌 账单机器人使用说明

【基础功能】

增加操作人
格式：增加操作人xx(可以是用户名，也可以是用户id)

清空今日账单
格式：清空今日账单

设置汇率
格式：设置汇率 7.01

查看汇率
格式：查看汇率

设置费率
格式：设置费率 0.1

查看费率
格式：查看费率

【账单操作】
入款（张三为备注，可不填）
支持格式：1、+1000 张三 2、+1000*7.0 张三 3、+10000/7.01 张三

下发（张三为备注，可不填）
格式：下发 1000 张三
`

	botSender, err := botapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return err
	}

	msg := botapi.NewMessage(l.chatGroupID, strings.TrimSpace(helpText))
	_, err = botSender.Send(msg)
	return err
}
