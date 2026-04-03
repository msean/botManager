package ledger

import (
	"fmt"
	"regexp"
	"time"

	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
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

	info, err := tronscan.GetAccountInfo(address)
	if err != nil {
		return "", fmt.Errorf("获取账户信息失败: %v", err)
	}

	now := time.Now()

	// 今日
	todayStart, todayEnd := tronscan.GetDayRange(now)
	todayStat, err := tronscan.CalcUSDTStat(address, todayStart, todayEnd)
	if err != nil {
		return "", fmt.Errorf("获取今日数据失败: %v", err)
	}

	// 昨日
	yesterday := now.AddDate(0, 0, -1)
	yStart, yEnd := tronscan.GetDayRange(yesterday)
	yStat, err := tronscan.CalcUSDTStat(address, yStart, yEnd)
	if err != nil {
		return "", fmt.Errorf("获取昨日数据失败: %v", err)
	}

	// 0点USDT余额
	zeroUSDT := tronscan.CalcUSDTZeroBalance(
		info.USDTBalance,
		todayStat.In,
		todayStat.Out,
	)

	// ===== 构建返回文本 =====
	msg := fmt.Sprintf(
		"📊 地址统计\n\n"+
			"📍 地址：%s\n\n"+
			"💰 当前余额：\n"+
			"USDT：%.2f\n"+
			"TRX：%.2f\n\n"+
			"📅 今日：\n"+
			"收入：%.2f USDT\n"+
			"支出：%.2f USDT\n\n"+
			"📆 昨日：\n"+
			"收入：%.2f USDT\n"+
			"支出：%.2f USDT\n\n"+
			"🕛 今日0点余额：%.2f USDT"+
			address,
		info.USDTBalance,
		info.TRXBalance,
		todayStat.In,
		todayStat.Out,
		yStat.In,
		yStat.Out,
		zeroUSDT,
	)

	return msg, nil
}
