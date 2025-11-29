package recharge

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

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
		botID        int64 // 机器人ID
		publishTimes int   // 发布次数
	}
)

func NewPay(botID int64, publishTimes int) *Pay {
	return &Pay{botID: botID, publishTimes: publishTimes}
}

func (pay *Pay) RandomPrice() (float64, error) {
	rechargeCnfList := cache.NewRechargeCnfListCache(pay.botID)
	if _, err := cache.CacheGetItem(rechargeCnfList); err != nil {
		return 0, err
	}

	var base float64
	found := false

	for _, object := range rechargeCnfList.Objects {
		if object.PublishTimes == pay.publishTimes {
			base = object.Price
			found = true
			break
		}
	}

	if !found {
		return 0, fmt.Errorf("no config found for publishTimes=%d", pay.publishTimes)
	}

	second := rand.Intn(10)
	third := rand.Intn(10)

	randomDecimal := float64(second*10+third) / 1000.0

	return utils.FloatReserve(float64(base+randomDecimal), 3), nil
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
	deadline := time.Now().Add(constant.OrderMatchAgo * time.Minute)

	// 批量更新
	err := global.GVA_DB.Model(&recharge.UserRechargeRecord{}).
		Where("status = ?", constant.AdRechargeCreate).
		Where("created_at >= ?", deadline).
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
					if _, err = botHandler.TgSend(token, ch.ChannelID, medias, buttons); err != nil {
						global.GVA_LOG.Error("HandleAdConfirm TgSend", zap.Int64("channelID", ch.ChannelID), zap.Error(err))
					}
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
