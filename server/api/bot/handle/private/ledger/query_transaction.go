package ledger

import (
	"fmt"
	"regexp"
	"time"

	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"github.com/msean/botmanager/server/utils/transaction"
)

var tronAddressRegex = regexp.MustCompile(`^T[1-9A-HJ-NP-Za-km-z]{33}$`)

type TronAddressParser struct {
	botModel bot.Bot
	update   botapi.Update
	address  string
}

// Match
func (p *TronAddressParser) Match(botModel bot.Bot, update botapi.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := update.Message.Text

	if tronAddressRegex.MatchString(text) {
		p.botModel = botModel
		p.update = update
		p.address = text
		return true
	}

	return false
}

// Handle
func (p *TronAddressParser) Handle() error {

	chatID := p.update.Message.Chat.ID

	msg, err := BuildAddressReport(p.address)
	if err != nil {
		replyText(p.botModel.Token, chatID, "❌ 查询失败："+err.Error())
		return nil
	}

	replyText(p.botModel.Token, chatID, msg)
	return nil
}

func BuildAddressReport(address string) (string, error) {

	acc, err := transaction.GetAccountInfo(address)
	if err != nil {
		return "", err
	}

	stats, err := transaction.GetTxStats(address)
	if err != nil {
		return "", err
	}

	msg := fmt.Sprintf(`
地址：%s

💰 TRX余额：%.2f
💵 USDT余额：%.2f

📊 累计交易：%d
⏱ 激活时间：%s

📅 今日
收入：%.2f
支出：%.2f

📆 昨日
收入：%.2f
支出：%.2f
`,
		address,
		acc.TRXBalance,
		acc.USDTBalance,
		acc.TotalTxCount,
		time.UnixMilli(stats.FirstTime).Format("2006-01-02 15:04:05"),
		stats.TodayIn,
		stats.TodayOut,
		stats.YesterdayIn,
		stats.YesterdayOut,
	)

	return msg, nil
}
