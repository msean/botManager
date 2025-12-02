package recharge

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/bot"
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
	baseStr := fmt.Sprintf("%.3f", base)
	if !strings.HasSuffix(baseStr, ".000") {
		return base
	}
	rand2 := rand.Intn(100)
	randomDecimal := float64(rand2) / 1000.0
	newPrice := base + randomDecimal
	return utils.FloatReserve(newPrice, 3)
}

func (pay *Pay) Recharge(token string, userID int64, chatID int64, updateID int64, amount float64) {

	price := pay.RandomPrice(amount)

	paymentAddr, err := pay.GetPaymentAddr()
	if err != nil {
		global.GVA_LOG.Error("GetPaymentAddress failed", zap.Error(err))
		return
	}

	record := recharge.UserRechargeRecord{
		BotID:       pay.botID,
		UserID:      userID,
		UpdateID:    int64(updateID),
		ChatID:      chatID,
		Price:       price,
		Status:      1, // 创建
		PaymentAddr: paymentAddr,
	}

	if err := global.GVA_DB.Create(&record).Error; err != nil {
		global.GVA_LOG.Error("create recharge record failed", zap.Error(err))
		return
	}

	createdAt := record.CreatedAt.Format("2006-01-02 15:04:05")
	// 给 Telegram 的消息内容
	msg := FormatRechargeMessage(
		record.ID,
		fmt.Sprintf("%.3f", record.Price),
		record.PaymentAddr,
		createdAt,
		constant.OrderLeftPaid,
	)

	// 发送给 Telegram
	msgConfig := tgbotapi.NewMessage(chatID, msg)
	msgConfig.ParseMode = "MarkdownV2"

	botApi, _ := bot_handler.NewBot(token)
	botApi.SendMarkDownMessage(chatID, msg)
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

func CheckExpiredOrders() {
	deadline := time.Now().Add(-constant.OrderMatchAgo * time.Minute)

	// 批量更新
	err := global.GVA_DB.Model(&recharge.UserRechargeRecord{}).
		Where("status = ?", constant.AdRechargeCreate).
		Where("created_at <= ?", deadline).
		Updates(map[string]any{
			"status":     constant.AdRechargeTimeout,
			"updated_at": time.Now(),
		}).Error

	if err != nil {
		global.GVA_LOG.Error("订单超时更新失败", zap.Error(err))
	}
}

func ReconcileAccounts() {
	db := global.GVA_DB

	bots, err := dao.BotDao.All(db)
	if err != nil {
		global.GVA_LOG.Error("获取机器人失败", zap.Error(err))
		return
	}

	for _, botModel := range bots {
		reconcileAccount(botModel)
	}
}

func reconcileAccount(botModel bot.Bot) (err error) {
	botID := botModel.BotID
	token := botModel.Token
	// 获取机器人所有的channel
	var channels []bot.BotChannel
	if err = global.GVA_DB.Where("bot_id = ?", botID).Find(&channels).Error; err != nil {
		global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int64("botID", botID), zap.Error(err))
		return
	}

	var botHandler *bot_handler.Bot
	if botHandler, err = bot_handler.NewBot(token); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm NewBot", zap.Int64("botID", botID), zap.Error(err))
		return
	}
	_ = botHandler
	// 1. 读取收款地址
	key := fmt.Sprintf("payment:%d", botID)
	paymentSysCnf, err := cache.LoadSyscnf(key, false, "")
	if err != nil || paymentSysCnf.Value == "" {
		global.GVA_LOG.Error("机器人未设置收款地址", zap.Int64("BotID", botID))
		return
	}
	paymentAddr := paymentSysCnf.Value

	var buttons any
	var hasPublishCfg bool
	cmdCfg := cache.NewBotCmdCache(int64(botID), constant.BotReplyCnfPublish2Channel, constant.BotReplyCnfType)
	if hasPublishCfg, err = cache.CacheGetItem(cmdCfg); err != nil {
		global.GVA_LOG.Error("handleBot", zap.Int("botID", int(botID)), zap.Error(err))
		return
	}

	if hasPublishCfg {
		buttons = bot_handler.ParseContentFromCfg(*cmdCfg, constant.ButtonTypeInline)
		global.GVA_LOG.Debug("handleBot", zap.Any("buttons", buttons))
	}

	trxResp, err := utils.FetchTransactions(paymentAddr, 20)
	if err != nil || !trxResp.Success {
		global.GVA_LOG.Error("获取链上交易失败", zap.Error(err))
		return
	}
	// trxResp.Data = append(trxResp.Data, utils.TronResponseData{
	// 	TransactionID:  "mock_tx_1001",
	// 	BlockTimestamp: 1764518680000, // 2025-11-30 22:40 北京时间
	// 	From:           "TEST_FROM_ADDRESS",
	// 	To:             "TKBDsYcVgvBMFi2qmhf88JDaMPYkqH8x2E",
	// 	Type:           "Transfer",
	// 	Value:          "10076000", // 10.085 * 1e6
	// 	TokenInfo: struct {
	// 		Symbol   string `json:"symbol"`
	// 		Address  string `json:"address"`
	// 		Decimals int    `json:"decimals"`
	// 		Name     string `json:"name"`
	// 	}{
	// 		Symbol:   "USDT",
	// 		Address:  "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
	// 		Decimals: 6,
	// 		Name:     "Tether USD",
	// 	},
	// })
	var orders []recharge.UserRechargeRecord
	err = global.GVA_DB.Where("bot_id = ? AND status = 1", botID).Find(&orders).Error
	if err != nil {
		global.GVA_LOG.Error("查询订单失败", zap.Error(err))
		return
	}

	for _, order := range orders {
		global.GVA_LOG.Info("reconcileAccount", zap.Uint("orderID", order.ID), zap.Any("order", order))
		var medias []bot_handler.MediaItem
		if err = json.Unmarshal([]byte(order.PublishContent), &medias); err != nil {
			global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int("botID", int(botID)), zap.Any("val", order.PublishContent), zap.Error(err))
			continue
		}
		for _, trx := range trxResp.Data {
			var match bool
			// 更新
			if match = matchTransaction(paymentAddr, order, trx); match {
				if err := global.GVA_DB.Model(&recharge.UserRechargeRecord{}).
					Where("id = ? AND status = 1", order.ID).
					Updates(map[string]interface{}{
						"status":     constant.AdRechargePaid,
						"tx_id":      trx.TransactionID,
						"updated_at": time.Now(),
					}).Error; err != nil {
					global.GVA_LOG.Error("更新订单失败", zap.Error(err))
				}

				for _, ch := range channels {
					// global.GVA_LOG.Error("HandleAdConfirm TgSend", zap.Int64("channelID", ch.ChannelID), zap.Error(err))
					if _, err = botHandler.TgSend(ch.ChannelID, medias, buttons); err != nil {
						global.GVA_LOG.Error("HandleAdConfirm TgSend", zap.Int64("channelID", ch.ChannelID), zap.Error(err))
					}
				}
				for i := range medias {
					if medias[i].Type == "text" {
						medias[i].Text += "\n\n✅ 发布成功"
						break
					}
				}
				if _, err = botHandler.TgSend(order.ChatID, medias, nil); err != nil {
					global.GVA_LOG.Error("HandleAdConfirm NewBot", zap.Int64("botID", botID), zap.Any("medias", medias), zap.Any("buttons", buttons), zap.Error(err))
				}
			}
		}
	}

	return
}

func matchTransaction(paymentAddr string, order recharge.UserRechargeRecord, trx utils.TronResponseData) (match bool) {
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
