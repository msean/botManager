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
	"github.com/xuri/excelize/v2"
)

type LedgerSummaryHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	ShouldPermissionAwareWithOutAdmin
}

// ================== 匹配 ==================
func (l *LedgerSummaryHandler) Match(botModel bot.Bot, update botapi.Update) bool {

	// ✅ 处理按钮点击（导出）
	if update.CallbackQuery != nil {

		data := update.CallbackQuery.Data

		if strings.HasPrefix(data, "export_excel_") {

			date := strings.TrimPrefix(data, "export_excel_")

			l.botModel = botModel
			l.chatGroupID = update.CallbackQuery.Message.Chat.ID

			go l.handleExport(date) // 异步执行

			return true
		}
	}

	// ✅ 处理“下课”
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

// ================== 主处理 ==================
func (l *LedgerSummaryHandler) Handle() error {

	today := time.Now().Format("2006-01-02")

	list, err := l.getListByDate(today)
	if err != nil {
		return err
	}

	if len(list) == 0 {
		return l.reply("📊 今日暂无账单")
	}

	groupMap, totalIn, totalOut := l.buildStat(list)

	var builder strings.Builder
	builder.WriteString("📊 今日账单统计\n\n")

	for k, v := range groupMap {

		balance := v.in - v.out

		builder.WriteString(fmt.Sprintf(
			"%s\n收入：%.2f\n收入笔数：%d\n支出：%.2f\n支出笔数：%d\n余额：%.2f\n\n",
			k, v.in, v.inCount, v.out, v.outCount, balance,
		))
	}

	totalBalance := totalIn - totalOut

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

// ================== 查询 ==================
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

// ================== 统计 ==================
type stat struct {
	in       float64
	out      float64
	inCount  int
	outCount int
}

func (l *LedgerSummaryHandler) buildStat(list []ledger2.Ledger) (map[string]*stat, float64, float64) {

	groupMap := make(map[string]*stat)
	var totalIn, totalOut float64

	for _, v := range list {

		key := v.UserName + "+" + v.Account

		if _, ok := groupMap[key]; !ok {
			groupMap[key] = &stat{}
		}

		if v.ActionType == 1 {
			groupMap[key].in += v.Amount
			groupMap[key].inCount++
			totalIn += v.Amount
		} else {
			groupMap[key].out += v.Amount
			groupMap[key].outCount++
			totalOut += v.Amount
		}
	}

	return groupMap, totalIn, totalOut
}

// ================== 回复 ==================
func (l *LedgerSummaryHandler) reply(text string) error {

	botSender, err := botapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return err
	}

	today := time.Now().Format("2006-01-02")

	msg := botapi.NewMessage(l.chatGroupID, text)

	// ✅ 按钮
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

// ================== 导出处理 ==================
func (l *LedgerSummaryHandler) handleExport(date string) {

	list, err := l.getListByDate(date)
	if err != nil || len(list) == 0 {
		l.reply("该日期暂无数据")
		return
	}

	filePath := fmt.Sprintf("ledger_%s.xlsx", date)

	err = l.generateExcel(list, filePath, date)
	if err != nil {
		l.reply("导出失败")
		return
	}

	botSender, _ := botapi.NewBotAPI(l.botModel.Token)

	doc := botapi.NewDocument(l.chatGroupID, botapi.FilePath(filePath))
	botSender.Send(doc)

	_ = os.Remove(filePath)
}

// ================== Excel ==================
func (l *LedgerSummaryHandler) generateExcel(list []ledger2.Ledger, filePath string, date string) error {

	f := excelize.NewFile()
	sheet := "Sheet1"

	headers := []string{"日期", "账户", "收入", "收入笔数", "支出", "支出笔数", "余额"}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	groupMap, totalIn, totalOut := l.buildStat(list)

	row := 2

	for k, v := range groupMap {

		balance := v.in - v.out

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), date)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), k)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.in)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), v.inCount)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), v.out)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), v.outCount)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), balance)

		row++
	}

	totalBalance := totalIn - totalOut

	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "合计")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), totalIn)
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), totalOut)
	f.SetCellValue(sheet, fmt.Sprintf("G%d", row), totalBalance)

	return f.SaveAs(filePath)
}

func (l *LedgerSummaryHandler) generateCSV(list []ledger2.Ledger, filePath string, date string) error {

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 表头
	header := []string{"日期", "账户", "收入", "收入笔数", "支出", "支出笔数", "余额"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// ===== 统计 =====
	type stat struct {
		in       float64
		out      float64
		inCount  int
		outCount int
	}

	groupMap := make(map[string]*stat)

	for _, v := range list {
		key := v.UserName + "+" + v.Account

		if _, ok := groupMap[key]; !ok {
			groupMap[key] = &stat{}
		}

		if v.ActionType == 1 {
			groupMap[key].in += v.Amount
			groupMap[key].inCount++
		} else {
			groupMap[key].out += v.Amount
			groupMap[key].outCount++
		}
	}

	// ===== 写数据 =====
	var totalIn, totalOut float64

	for k, v := range groupMap {

		balance := v.in - v.out

		row := []string{
			date,
			k,
			fmt.Sprintf("%.2f", v.in),
			strconv.Itoa(v.inCount),
			fmt.Sprintf("%.2f", v.out),
			strconv.Itoa(v.outCount),
			fmt.Sprintf("%.2f", balance),
		}

		if err := writer.Write(row); err != nil {
			return err
		}

		totalIn += v.in
		totalOut += v.out
	}

	// ===== 合计 =====
	totalBalance := totalIn - totalOut

	totalRow := []string{
		"合计",
		"",
		fmt.Sprintf("%.2f", totalIn),
		"",
		fmt.Sprintf("%.2f", totalOut),
		"",
		fmt.Sprintf("%.2f", totalBalance),
	}

	if err := writer.Write(totalRow); err != nil {
		return err
	}

	return nil
}
