package chat_group

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LedgerHandler struct {
	model       ledger.Ledger
	botModel    bot.Bot
	chatGroupID int64
	confModel   cache.LedgerPermissionCache
}

func (l *LedgerHandler) Match(botModel bot.Bot, update tgbotapi.Update) (match bool) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	msg := update.Message
	input := strings.TrimSpace(msg.Text)

	l.botModel = botModel
	l.chatGroupID = msg.Chat.ID

	rawInput := input

	var (
		actionType int
		body       string
	)

	// 判断操作类型
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

	// body 示例：
	// "8000 aa"
	// "-8000    aa"
	// "1000*7.8 aa"
	parts := strings.Fields(body)
	if len(parts) == 0 {
		return
	}

	amountExpr := parts[0]
	remark := ""
	if len(parts) > 1 {
		remark = strings.TrimSpace(body[len(amountExpr):])
	}

	amount, err := parseAmount(amountExpr)
	if err != nil {
		global.GVA_LOG.Error("Ledger parse amount failed",
			zap.String("expr", amountExpr),
			zap.Error(err),
		)
		return
	}

	global.GVA_LOG.Debug("LedgerHandler", zap.Any("l.confModel.CurrentFeeRate", l.confModel.CurrentFeeRate), zap.Any("AmountWithFee", utils.FloatReserve((100-l.confModel.CurrentFeeRate)*amount/100, 2)))
	l.model = ledger.Ledger{
		OprUserID:        msg.From.ID,
		OprUserFirstName: msg.From.FirstName,
		OprUserLastName:  msg.From.LastName,
		OprUserNickname:  msg.From.UserName,

		CurrentFeeRate: l.confModel.CurrentFeeRate,
		ActionType:     actionType,
		Amount:         amount,
		AmountWithFee:  utils.FloatReserve((100-l.confModel.CurrentFeeRate)*amount/100, 2),
		Remark:         remark,

		BotID:       l.botModel.BotID,
		ChatGroupID: msg.Chat.ID,
		MessageID:   int64(msg.MessageID),
		RawInput:    rawInput,
	}

	return true
}

func (l *LedgerHandler) HasPerMission() (permit bool, err error) {
	ledgerPermission := cache.NewLedgerPermissionCache(l.botModel.BotID, l.chatGroupID)
	var has bool
	if has, err = cache.CacheGetItem(ledgerPermission); err != nil {
		global.GVA_LOG.Error("Ledger HasPerMission", zap.Error(err))
		return
	}
	if !has {
		return
	}
	if !ledgerPermission.HasUserPermission(l.model.OprUserID, l.model.OprUserNickname) {
		return
	}
	permit = true
	l.confModel = *ledgerPermission
	return
}

func (l *LedgerHandler) Handle() (err error) {
	var permit bool
	if permit, err = l.HasPerMission(); err != nil || !permit {
		return
	}

	if err = l.Create(); err != nil {
		global.GVA_LOG.Error("Ledger Create", zap.Error(err))
		return
	}

	var reply, url string
	if reply, url, err = l.BuildReply(global.GVA_MYSQL); err != nil {
		return
	}

	botSender, err := tgbotapi.NewBotAPI(l.botModel.Token)
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(l.chatGroupID, reply)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("📊 点击查看完整账单", url),
		),
	)
	_, err = botSender.Send(msg)
	return
}

func (l *LedgerHandler) Create() error {
	return global.GVA_MYSQL.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&l.model).Error
	})
}

func (l *LedgerHandler) BuildReply(db *gorm.DB) (content string, url string, err error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var records []ledger.Ledger
	err = db.
		Where("chat_group_id = ?", l.chatGroupID).
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
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
		// var showRemark string
		// if r.ActionType == 1 {
		// 	showRemark = r.Remark
		// } else {
		// 	showRemark = r.OprUserNickname
		// }

		line := fmt.Sprintf(
			"%s %.2f %s",
			r.CreatedAt.Format("15:04:05"),
			r.Amount,
			strings.TrimSpace(r.Remark),
		)

		if r.ActionType == 1 {
			incomeLines = append(incomeLines, line)
			totalIncome += r.Amount
			totalAmountWithFee += r.AmountWithFee
		} else {
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

// ----------------- 工具函数 -----------------

// 支持：
// 8000
// -8000
// 1000*7.8
func parseAmount(expr string) (float64, error) {
	expr = strings.TrimSpace(expr)

	if !strings.Contains(expr, "*") {
		return strconv.ParseFloat(expr, 64)
	}

	parts := strings.Split(expr, "*")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid amount expr: %s", expr)
	}

	a, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, err
	}
	b, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, err
	}

	return a * b, nil
}
