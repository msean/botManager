package ledger2

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger2"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

type LedgerSummaryHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	ShouldPermissionAwareWithOutAdmin
}

func (l *LedgerSummaryHandler) Match(botModel bot.Bot, update botapi.Update) bool {

	// ✅ 导出按钮
	if update.CallbackQuery != nil {

		data := update.CallbackQuery.Data

		if strings.HasPrefix(data, "export_excel_") {

			date := strings.TrimPrefix(data, "export_excel_")

			l.botModel = botModel
			l.chatGroupID = update.CallbackQuery.Message.Chat.ID

			go l.handleExport(date)

			return true
		}
	}

	// ✅ 下课
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)

	if text != "下课" {
		return false
	}

	l.botModel = botModel
	l.chatGroupID = update.Message.Chat.ID

	return true
}

func (l *LedgerSummaryHandler) Handle() error {

	today := time.Now().Format("2006-01-02")

	list, err := l.getListByDate(today)
	if err != nil {
		return err
	}

	if len(list) == 0 {
		return l.reply("📊 今日暂无账单")
	}

	groupMap, totalIn, totalOut, totalBalanceAdjust := l.buildStat(list)

	var builder strings.Builder
	builder.WriteString("📊 今日账单统计\n\n")

	for k, v := range groupMap {

		balance := v.in - v.out + v.balance

		builder.WriteString(fmt.Sprintf(
			"%s\n收入：%.2f\n收入笔数：%d\n支出：%.2f\n支出笔数：%d\n余额：%.2f\n\n",
			k, v.in, v.inCount, v.out, v.outCount, balance,
		))
	}

	totalBalance := totalIn - totalOut + totalBalanceAdjust

	builder.WriteString("——————\n")
	builder.WriteString(fmt.Sprintf(
		"汇总：\n当日收入：%.2f\n当日支出：%.2f\n当日余额：%.2f",
		totalIn, totalOut, totalBalance,
	))

	// 关闭会话
	_ = global.GVA_MYSQL.Model(&ledger2.LedgerSession{}).
		Where("bot_id = ? AND chat_group_id = ? AND work_date = ?",
			l.botModel.BotID, l.chatGroupID, today).
		Update("is_active", 0).Error

	return l.reply(builder.String())
}

func (l *LedgerSummaryHandler) getListByDate(date string) ([]ledger2.Ledger, error) {

	var list []ledger2.Ledger

	err := global.GVA_MYSQL.Where(
		"bot_id = ? AND chat_group_id = ? AND DATE(created_at) = ?",
		l.botModel.BotID,
		l.chatGroupID,
		date,
	).Find(&list).Error

	return list, err
}

type stat struct {
	in       float64
	out      float64
	balance  float64
	inCount  int
	outCount int
}

func (l *LedgerSummaryHandler) buildStat(list []ledger2.Ledger) (map[string]*stat, float64, float64, float64) {

	groupMap := make(map[string]*stat)
	userAccountMap := make(map[string]string)

	var totalIn, totalOut, totalBalanceAdjust float64

	for _, v := range list {
		if v.ActionType != 3 {
			userAccountMap[v.UserName] = v.Account
		}
	}

	for _, v := range list {

		account := v.Account

		if v.ActionType == 3 {
			if acc, ok := userAccountMap[v.UserName]; ok {
				account = acc
			} else {
				account = "余额" // ✅ 关键修复
			}
		}

		key := v.UserName + "+" + account

		if _, ok := groupMap[key]; !ok {
			groupMap[key] = &stat{}
		}

		switch v.ActionType {
		case 1:
			groupMap[key].in += v.Amount
			groupMap[key].inCount++
			totalIn += v.Amount
		case 2:
			groupMap[key].out += v.Amount
			groupMap[key].outCount++
			totalOut += v.Amount
		case 3:
			groupMap[key].balance += v.Amount
			totalBalanceAdjust += v.Amount
		}
	}

	return groupMap, totalIn, totalOut, totalBalanceAdjust
}

func (l *LedgerSummaryHandler) reply(text string) error {

	botSender, err := botapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return err
	}

	today := time.Now().Format("2006-01-02")

	msg := botapi.NewMessage(l.chatGroupID, text)

	msg.ReplyMarkup = botapi.NewInlineKeyboardMarkup(
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonData(
				"📥 导出Excel（"+today+"）",
				"export_excel_"+today,
			),
		),
	)

	_, err = botSender.Send(msg)
	return err
}

func (l *LedgerSummaryHandler) handleExport(date string) {

	list, err := l.getListByDate(date)
	if err != nil || len(list) == 0 {
		l.reply("该日期暂无数据")
		return
	}

	filePath := fmt.Sprintf("ledger_%s.csv", date)

	err = l.generateCSV(list, filePath, date)
	if err != nil {
		l.reply("导出失败")
		return
	}

	botSender, _ := botapi.NewBotAPI(l.botModel.Token)

	doc := botapi.NewDocument(l.chatGroupID, botapi.FilePath(filePath))
	botSender.Send(doc)

	_ = os.Remove(filePath)
}

func (l *LedgerSummaryHandler) generateCSV(list []ledger2.Ledger, filePath string, date string) error {

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, _ = file.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"日期", "账户", "收入", "收入笔数", "支出", "支出笔数", "余额"}
	writer.Write(header)

	type stat struct {
		in       float64
		out      float64
		balance  float64
		inCount  int
		outCount int
	}

	groupMap := make(map[string]*stat)
	userAccountMap := make(map[string]string)

	for _, v := range list {
		if v.ActionType != 3 {
			userAccountMap[v.UserName] = v.Account
		}
	}

	for _, v := range list {

		account := v.Account

		if v.ActionType == 3 {
			if acc, ok := userAccountMap[v.UserName]; ok {
				account = acc
			} else {
				account = "余额" // ✅ 修复点
			}
		}

		key := v.UserName + "+" + account

		if _, ok := groupMap[key]; !ok {
			groupMap[key] = &stat{}
		}

		switch v.ActionType {
		case 1:
			groupMap[key].in += v.Amount
			groupMap[key].inCount++
		case 2:
			groupMap[key].out += v.Amount
			groupMap[key].outCount++
		case 3:
			groupMap[key].balance += v.Amount
		}
	}

	var totalIn, totalOut, totalBalanceAdjust float64

	for k, v := range groupMap {

		balance := v.in - v.out + v.balance

		row := []string{
			date,
			k,
			fmt.Sprintf("%.2f", v.in),
			strconv.Itoa(v.inCount),
			fmt.Sprintf("%.2f", v.out),
			strconv.Itoa(v.outCount),
			fmt.Sprintf("%.2f", balance),
		}

		writer.Write(row)

		totalIn += v.in
		totalOut += v.out
		totalBalanceAdjust += v.balance
	}

	totalBalance := totalIn - totalOut + totalBalanceAdjust

	writer.Write([]string{
		"合计",
		"",
		fmt.Sprintf("%.2f", totalIn),
		"",
		fmt.Sprintf("%.2f", totalOut),
		"",
		fmt.Sprintf("%.2f", totalBalance),
	})

	return nil
}
