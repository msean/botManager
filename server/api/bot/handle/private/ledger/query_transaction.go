package ledger

import (
	"fmt"
	"regexp"

	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"github.com/msean/botmanager/server/utils/transaction/trongrid"
	"github.com/msean/botmanager/server/utils/transaction/tronscan"
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

	// ===== 1. 获取余额（tronscan）=====
	info, err := tronscan.GetAccountInfo(address)
	if err != nil {
		return "", fmt.Errorf("获取账户信息失败: %v", err)
	}

	// ===== 2. 获取今日 + 昨日统计（trongrid）=====
	todayStat, yesterdayStat, err := trongrid.CalcTodayYesterday(address)
	if err != nil {
		return "", fmt.Errorf("获取交易统计失败: %v, 请稍后再尝试", err)
	}

	if todayStat == nil {
		msg := fmt.Sprintf(
			"您的地址：%s\n\n"+
				"USDT余额：%.2f\n"+
				"TRX余额：%.2f\n\n"+
				"今日：\n"+
				"收入：%.2f USDT\n"+
				"支出：%.2f USDT\n\n"+
				"昨日：\n"+
				"收入：%.2f USDT\n"+
				"支出：%.2f USDT\n\n"+
				"🕛 今日0点余额：%.2f USDT",
			address,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
		)
		return msg, nil
	}

	// ===== 3. 计算今日0点余额 =====
	zeroUSDT := tronscan.CalcUSDTZeroBalance(
		info.USDTBalance,
		todayStat.In,
		todayStat.Out,
	)

	// ===== 4. 构建消息 =====
	msg := fmt.Sprintf(
		"您的地址：%s\n\n"+
			"USDT余额：%.2f\n"+
			"TRX余额：%.2f\n\n"+
			"今日：\n"+
			"收入：%.2f USDT\n"+
			"支出：%.2f USDT\n\n"+
			"昨日：\n"+
			"收入：%.2f USDT\n"+
			"支出：%.2f USDT\n\n"+
			"🕛 今日0点余额：%.2f USDT",
		address,
		info.USDTBalance,
		info.TRXBalance,
		todayStat.In,
		todayStat.Out,
		yesterdayStat.In,
		yesterdayStat.Out,
		zeroUSDT,
	)

	return msg, nil
}
