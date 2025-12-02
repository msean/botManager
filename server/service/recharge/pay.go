package recharge

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/recharge"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils"
	"github.com/msean/botmanager/server/utils/bot_handler"
	"go.uber.org/zap"
)

type (
	Pay struct {
		botID int64 // 机器人ID
	}
)

func NewPay(botID int64) *Pay {
	return &Pay{botID: botID}
}

func (pay *Pay) RandomPrice(base float64) float64 {
	rand2 := rand.Intn(100)
	randomDecimal := float64(rand2) / 1000.0
	newPrice := base + randomDecimal
	return utils.FloatReserve(newPrice, 3)
}

func (pay *Pay) Recharge(token string, userID int64, chatID int64, msgID int, amount float64) (err error) {

	price := pay.RandomPrice(amount)

	var paymentAddr string
	paymentAddr, err = pay.GetPaymentAddr()
	if err != nil {
		global.GVA_LOG.Error("GetPaymentAddress failed", zap.Error(err))
		return
	}

	record := recharge.UserRechargeRecord{
		BotID:       pay.botID,
		UserID:      userID,
		MsgID:       msgID,
		ChatID:      chatID,
		Price:       price,
		Status:      1, // 创建
		PaymentAddr: paymentAddr,
	}

	if err = global.GVA_DB.Create(&record).Error; err != nil {
		global.GVA_LOG.Error("create recharge record failed", zap.Error(err))
		return
	}

	createdAt := bot_handler.EscapeMarkdownV2(record.CreatedAt.Format("2006-01-02 15:04:05"))
	// 给 Telegram 的消息内容
	msg := FormatRechargeMessage(
		record.ID,
		fmt.Sprintf("%.3f", record.Price),
		record.PaymentAddr,
		createdAt,
		constant.OrderLeftPaid,
	)

	btnAmount := tgbotapi.NewInlineKeyboardButtonData(
		fmt.Sprintf("💰 请支付%.3f USDT", record.Price),
		fmt.Sprintf("recharge_amount:%d", record.ID),
	)

	// 按钮2：取消充值
	btnCancel := tgbotapi.NewInlineKeyboardButtonData(
		"❌ 取消充值",
		fmt.Sprintf("/rechargeCancel:%d", record.ID),
	)

	// 组装键盘
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{btnAmount},
		[]tgbotapi.InlineKeyboardButton{btnCancel},
	)

	botApi, _ := bot_handler.NewBot(token)
	return botApi.SendMarkDownMessage(chatID, msg, keyboard)
}

func FormatRechargeMessage(orderID uint, amount, paymentAddr string, createdAt string, leftPaidMinutes int) string {
	return fmt.Sprintf(
		"订单号：%d\n"+
			"转账金额：`%s` USDT （点击即可复制）\n"+
			"转账地址：`%s` （点击即可复制）\n"+
			"充值时间：%s\n\n"+
			"⚠️注意：\n"+
			"▫️注意小数点 %s 转错金额不能到账\n"+
			"▫️请在%d分钟完成付款，转错金额不能到账。\n\n"+
			"转账%d分钟后没到账及时联系",
		orderID,
		bot_handler.EscapeMarkdownV2CodeBlock(amount),
		bot_handler.EscapeMarkdownV2CodeBlock(paymentAddr),
		createdAt,
		bot_handler.EscapeMarkdownV2(amount),
		leftPaidMinutes,
		leftPaidMinutes,
	)
}

func (pay *Pay) GetPaymentAddr() (paymentAddr string, err error) {
	// 获取支付方式
	var paymentWaySysCnf *cache.SysCnfCache
	if paymentWaySysCnf, err = cache.LoadSyscnf(constant.SysCnfPaymentWayKey, true, constant.DefaultSysCnfPaymentWay); err != nil {
		return
	}

	switch paymentWaySysCnf.Value {
	case constant.DefaultSysCnfPaymentWay:
		key := fmt.Sprintf("payment:%d", pay.botID)
		var paymentSysCnf *cache.SysCnfCache
		if paymentSysCnf, err = cache.LoadSyscnf(key, false, ""); err != nil {
			return
		}
		paymentAddr = paymentSysCnf.Value
	default:
		err = fmt.Errorf("支付方式暂未开通")
	}
	return
}

func MatchTransaction(paymentAddr string, order recharge.UserRechargeRecord, trx utils.TronResponseData) (match bool) {
	if trx.To != paymentAddr {
		return
	}

	global.GVA_LOG.Info("reconcileAccount matchTransaction", zap.Uint("orderID", order.ID), zap.Any("trx", trx))
	amount := utils.ParseAmount(trx.Value, trx.TokenInfo.Decimals)

	global.GVA_LOG.Info("reconcileAccount matchTransaction", zap.Any("amount1", math.Abs(amount-order.Price) > 0.000001), zap.Any("amount", amount))
	if math.Abs(amount-order.Price) > 0.000001 {
		return
	}

	trxTime := time.UnixMilli(trx.BlockTimestamp)
	if trxTime.After(order.CreatedAt.Add(constant.OrderMatchAgo * time.Minute)) {
		return
	}
	match = true
	return
}
