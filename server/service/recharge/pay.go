package recharge

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/recharge"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils"
	"github.com/msean/botmanager/server/utils/bot_handler"
	bot_api "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"github.com/msean/botmanager/server/utils/transaction/trongrid"
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

func (pay *Pay) Recharge(token string, userID int64, chatID int64, msgID int, userName string, amount float64) (err error) {

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
		UserName:    userName,
	}

	if err = global.GVA_MYSQL.Create(&record).Error; err != nil {
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

	btnAmount := bot_api.NewInlineKeyboardButtonData(
		fmt.Sprintf("💰 请支付%.3fUSDT", record.Price),
		fmt.Sprintf("recharge_amount:%d", record.ID),
	)

	// 按钮2：取消充值
	btnCancel := bot_api.NewInlineKeyboardButtonData(
		"❌ 取消充值",
		fmt.Sprintf("/rechargeCancel:%d", record.ID),
	)

	// 组装键盘
	keyboard := bot_api.NewInlineKeyboardMarkup(
		[]bot_api.InlineKeyboardButton{btnAmount},
		[]bot_api.InlineKeyboardButton{btnCancel},
	)

	bot_api, _ := bot_handler.NewBot(token)
	return bot_api.SendMarkDownMessage(chatID, msg, keyboard)
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

func MatchTransaction(paymentAddr string, order recharge.UserRechargeRecord, trx trongrid.TronResponseData) (match bool) {
	if trx.To != paymentAddr {
		return
	}

	global.GVA_LOG.Info("reconcileAccount matchTransaction", zap.Uint("orderID", order.ID), zap.Any("trx", trx))
	amount := trongrid.ParseAmount(trx.Value, trx.TokenInfo.Decimals)

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
