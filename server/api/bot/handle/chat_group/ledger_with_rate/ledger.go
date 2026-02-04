package rate_ledger

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger"
	"github.com/msean/botmanager/server/utils"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LedgerHandler struct {
	model       ledger.Ledger
	botModel    bot.Bot
	chatGroupID int64
	msg         botapi.Update
	ShouldPermissionAwareWithOutAdmin
}

/* ===================== Match ===================== */

func (l *LedgerHandler) Match(botModel bot.Bot, update botapi.Update) (match bool) {
	msg := update.Message
	input := strings.TrimSpace(msg.Text)

	l.botModel = botModel
	l.chatGroupID = msg.Chat.ID
	l.msg = update

	rawInput := input

	var (
		actionType int
		body       string
	)

	// 操作类型
	if strings.HasPrefix(input, "+") {
		actionType = 1
		body = strings.TrimSpace(input[1:])
	} else if strings.HasPrefix(input, "下发") {
		actionType = 2
		body = strings.TrimSpace(string([]rune(input)[2:]))
	} else {
		return
	}

	if body == "" {
		return
	}

	parts := strings.Fields(body)
	if len(parts) == 0 {
		return
	}

	amountExpr := parts[0]
	remark := ""
	if len(parts) > 1 {
		remark = strings.TrimSpace(body[len(amountExpr):])
	}

	parseResult, err := parseAmountExpr(amountExpr)
	if err != nil {
		global.GVA_LOG.Error("Ledger parse amount failed",
			zap.String("expr", amountExpr),
			zap.Error(err),
		)
		return
	}

	l.model = ledger.Ledger{
		OprUserID:        msg.From.ID,
		OprUserFirstName: msg.From.FirstName,
		OprUserLastName:  msg.From.LastName,
		OprUserNickname:  msg.From.UserName,

		ActionType: actionType,
		Amount:     parseResult.Amount,
		Remark:     remark,

		BotID:       l.botModel.BotID,
		ChatGroupID: msg.Chat.ID,
		MessageID:   int64(msg.MessageID),
		RawInput:    rawInput,
	}

	// 入款时支持 /汇率
	if actionType == 1 && parseResult.CurrencyRate > 0 {
		l.model.CurrentCurrencyFeeRate = parseResult.CurrencyRate
	}

	return true
}

/* ===================== Handle ===================== */

func (l *LedgerHandler) Handle() (err error) {

	// 入款：优先使用输入的汇率
	if l.model.ActionType == 1 && l.model.CurrentCurrencyFeeRate > 0 {
		l.model.AmountWithFee = utils.FloatReserve(
			l.model.Amount/l.model.CurrentCurrencyFeeRate, 2)
	} else {
		l.model.CurrentFeeRate = l.confModel.CurrentFeeRate
		l.model.CurrentCurrencyFeeRate = l.confModel.CurrentCurrencyFeeRate
		l.model.AmountWithFee = utils.FloatReserve(
			(100-l.confModel.CurrentFeeRate)*l.model.Amount/100, 2)
	}

	if err = l.Create(); err != nil {
		global.GVA_LOG.Error("Ledger Create", zap.Error(err))
		return
	}

	var reply, url string
	if reply, url, err = l.BuildReply(); err != nil {
		return
	}

	botSender, err := botapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return
	}

	msg := botapi.NewMessage(l.chatGroupID, reply)
	msg.ReplyMarkup = botapi.NewInlineKeyboardMarkup(
		botapi.NewInlineKeyboardRow(
			botapi.NewInlineKeyboardButtonURL("📊 点击查看完整账单", url),
		),
	)
	_, err = botSender.Send(msg)
	return
}

/* ===================== DB ===================== */

func (l *LedgerHandler) Create() error {
	return global.GVA_MYSQL.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&l.model).Error
	})
}

/* ===================== BuildReply ===================== */

func (l *LedgerHandler) BuildReply() (content string, url string, err error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var records []ledger.Ledger
	err = global.GVA_MYSQL.
		Where("bot_id=?", l.botModel.BotID).
		Where("chat_group_id = ?", l.chatGroupID).
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Where("deleted_at IS NULL").
		Order("id asc").
		Find(&records).Error
	if err != nil || len(records) == 0 {
		return
	}

	var (
		incomeLines        []string
		payoutLines        []string
		totalIncome        float64
		totalPayout        float64
		totalAmountWithFee float64
	)

	for _, r := range records {
		var line string

		if r.ActionType == 1 {
			// 入款展示：金额 / 汇率 = 实际
			if r.CurrentCurrencyFeeRate > 0 {
				line = fmt.Sprintf(
					"%s %.2f  / %.2f = %.2f",
					r.CreatedAt.Format("15:04:05"),
					r.Amount,
					r.CurrentCurrencyFeeRate,
					r.AmountWithFee,
				)
			} else {
				line = fmt.Sprintf(
					"%s %.2f %s",
					r.CreatedAt.Format("15:04:05"),
					r.Amount,
					strings.TrimSpace(r.Remark),
				)
			}

			incomeLines = append(incomeLines, line)
			totalIncome += r.Amount
			totalAmountWithFee += r.AmountWithFee
		} else {
			line = fmt.Sprintf(
				"%s %.2f %s",
				r.CreatedAt.Format("15:04:05"),
				r.Amount,
				strings.TrimSpace(r.Remark),
			)

			payoutLines = append(payoutLines, line)
			totalPayout += r.Amount
		}
	}

	url = fmt.Sprintf(
		"%s/#/ledger/full?bot_id=%d&chat_group_id=%d&idmin=%d&idmax=%d",
		global.GVA_CONFIG.System.Domain,
		l.botModel.BotID,
		l.chatGroupID,
		records[0].ID,
		records[len(records)-1].ID,
	)

	return fmt.Sprintf(
		`账单（%s）

入款（%d笔）：
%s

下发（%d笔）：
%s

当前费率：%.2f

总入款：%.2f
应下发：%.2f
总下发：%.2f
未下发：%.2f`,
		startOfDay.Format("2006-01-02"),
		len(incomeLines), strings.Join(incomeLines, "\n"),
		len(payoutLines), strings.Join(payoutLines, "\n"),
		l.confModel.CurrentFeeRate,
		totalIncome,
		totalAmountWithFee,
		totalPayout,
		totalAmountWithFee-totalPayout,
	), url, nil
}

/* ===================== 工具函数 ===================== */

type AmountParseResult struct {
	Amount       float64
	CurrencyRate float64 // 0 表示未传
}

// 支持：
// 10000
// 10000*7.8
// 10000/7.01
func parseAmountExpr(expr string) (*AmountParseResult, error) {
	expr = strings.TrimSpace(expr)

	// 金额 / 汇率
	if strings.Contains(expr, "/") {
		parts := strings.Split(expr, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid amount/rate expr: %s", expr)
		}

		amount, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return nil, err
		}

		rate, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return nil, err
		}

		return &AmountParseResult{
			Amount:       amount,
			CurrencyRate: rate,
		}, nil
	}

	// 金额 * 倍数
	if strings.Contains(expr, "*") {
		parts := strings.Split(expr, "*")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid amount expr: %s", expr)
		}

		a, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return nil, err
		}
		b, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return nil, err
		}

		return &AmountParseResult{
			Amount: a * b,
		}, nil
	}

	amount, err := strconv.ParseFloat(expr, 64)
	if err != nil {
		return nil, err
	}

	return &AmountParseResult{
		Amount: amount,
	}, nil
}
