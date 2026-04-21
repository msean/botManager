package ledger

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"github.com/msean/botmanager/server/utils/transaction/trongrid"
	"github.com/msean/botmanager/server/utils/transaction/tronscan"
)

type GroupStatParser struct {
	botModel bot.Bot
	update   botapi.Update
	groupID  uint
}

func (p *GroupStatParser) Match(botModel bot.Bot, update botapi.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)

	if strings.HasPrefix(text, "统计分组") {
		idStr := strings.TrimPrefix(text, "统计分组")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return false
		}

		p.botModel = botModel
		p.update = update
		p.groupID = uint(id)

		return true
	}

	return false
}

func (p *GroupStatParser) Handle() error {
	chatID := p.update.Message.Chat.ID

	var group ledger.LedgerAccountGroup
	err := global.GVA_MYSQL.First(&group, p.groupID).Error
	if err != nil {
		replyText(p.botModel.Token, chatID, "❌ 分组不存在")
		return nil
	}

	if group.AccountGroup == nil {
		replyText(p.botModel.Token, chatID, "❌ 分组没有账号")
		return nil
	}

	accounts := strings.Split(*group.AccountGroup, ",")

	replyText(p.botModel.Token, chatID, "🚀 开始统计，请稍候...")

	var totalUSDT float64
	var totalTRX float64
	var totalTodayIn float64
	var totalTodayOut float64
	var totalYesterdayIn float64
	var totalYesterdayOut float64
	var totalZeroUSDT float64

	for i, acc := range accounts {
		acc = strings.TrimSpace(acc)
		if acc == "" {
			continue
		}

		replyText(
			p.botModel.Token,
			chatID,
			fmt.Sprintf("⏳ 正在统计第 %d 个账户：%s", i+1, acc),
		)

		// ===== 1. 获取账户信息 =====
		info, err := tronscan.GetAccountInfo(acc)
		if err != nil {
			replyText(p.botModel.Token, chatID, "❌ "+acc+" 获取账户失败")
			continue
		}

		// ===== 2. 获取统计 =====
		todayStat, yesterdayStat, err := trongrid.CalcTodayYesterday(acc)
		if err != nil {
			replyText(p.botModel.Token, chatID, "❌ "+acc+" 统计失败")
			continue
		}

		if todayStat == nil {
			replyText(p.botModel.Token, chatID, "⚠️ "+acc+" 无交易数据")
			continue
		}

		// ===== 3. 今日0点余额 =====
		zeroUSDT := tronscan.CalcUSDTZeroBalance(
			info.USDTBalance,
			todayStat.In,
			todayStat.Out,
		)

		// ===== 4. 累加 =====
		totalUSDT += info.USDTBalance
		totalTRX += info.TRXBalance
		totalTodayIn += todayStat.In
		totalTodayOut += todayStat.Out
		totalYesterdayIn += yesterdayStat.In
		totalYesterdayOut += yesterdayStat.Out
		totalZeroUSDT += zeroUSDT

		// ===== 5. ⭐ 每个账号完整结果（核心）=====
		msg := fmt.Sprintf(
			"📍 地址：%s\n\n"+
				"💰 USDT余额：%.2f\n"+
				"💰 TRX余额：%.2f\n\n"+
				"📅 今日：\n收入：%.2f\n支出：%.2f\n\n"+
				"📅 昨日：\n收入：%.2f\n支出：%.2f\n\n"+
				"🕛 今日0点余额：%.2f",
			acc,
			info.USDTBalance,
			info.TRXBalance,
			todayStat.In,
			todayStat.Out,
			yesterdayStat.In,
			yesterdayStat.Out,
			zeroUSDT,
		)

		// 👉 直接输出该账号结果
		replyText(p.botModel.Token, chatID, msg)

		replyText(
			p.botModel.Token,
			chatID,
			fmt.Sprintf("✅ %s 统计完成", acc),
		)
	}

	// ===== 汇总 =====
	summary := fmt.Sprintf(
		"📊 分组汇总结果\n\n"+
			"💰 总USDT余额：%.2f\n"+
			"💰 总TRX余额：%.2f\n\n"+
			"📅 今日：\n收入：%.2f\n支出：%.2f\n\n"+
			"📅 昨日：\n收入：%.2f\n支出：%.2f\n\n"+
			"🕛 今日0点总余额：%.2f\n",
		totalUSDT,
		totalTRX,
		totalTodayIn,
		totalTodayOut,
		totalYesterdayIn,
		totalYesterdayOut,
		totalZeroUSDT,
	)

	replyText(p.botModel.Token, chatID, summary)

	return nil
}
