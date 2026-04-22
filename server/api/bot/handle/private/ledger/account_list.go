package ledger

import (
	"fmt"
	"strings"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

type GroupListParser struct {
	botModel bot.Bot
	update   botapi.Update
}

// 匹配
func (p *GroupListParser) Match(botModel bot.Bot, update botapi.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	if strings.TrimSpace(update.Message.Text) == "账户分组列表" {
		p.botModel = botModel
		p.update = update
		return true
	}

	return false
}

// 处理
func (p *GroupListParser) Handle() error {
	chatID := p.update.Message.Chat.ID

	var list []ledger.LedgerAccountGroup
	err := global.GVA_MYSQL.Find(&list).Error
	if err != nil {
		replyText(p.botModel.Token, chatID, "❌ 查询失败："+err.Error())
		return nil
	}

	if len(list) == 0 {
		replyText(p.botModel.Token, chatID, "暂无分组数据")
		return nil
	}

	var builder strings.Builder

	for i, item := range list {

		groupID := item.ID
		title := ""
		accountGroup := ""

		if item.Title != "" {
			title = item.Title
		}
		if item.AccountGroup != "" {
			accountGroup = item.AccountGroup
		}

		accounts := strings.Split(accountGroup, ",")

		builder.WriteString(fmt.Sprintf("ID：%d\n", groupID))
		builder.WriteString(fmt.Sprintf("标题：%s\n", title))
		builder.WriteString("地址列表：\n")

		for _, acc := range accounts {
			acc = strings.TrimSpace(acc)
			if acc != "" {
				builder.WriteString(fmt.Sprintf("  %s\n", acc))
			}
		}

		if i != len(list)-1 && len(list) > 1 {
			builder.WriteString("---------------------\n")
		}
	}

	replyText(p.botModel.Token, chatID, builder.String())
	return nil
}
