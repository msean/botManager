package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/recharge"
	botSrv "github.com/msean/botmanager/server/service/bot"
	"github.com/msean/botmanager/server/service/cache"
	rechargeSrv "github.com/msean/botmanager/server/service/recharge"
	"github.com/msean/botmanager/server/service/system"
	"github.com/msean/botmanager/server/utils"
	"github.com/msean/botmanager/server/utils/bot_handler"
	"go.uber.org/zap"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup   system.ServiceGroup
	BotServiceGroup      botSrv.ServiceGroup
	RechargeServiceGroup rechargeSrv.ServiceGroup
}

func Init() {
	botSrv.InitBotTaskManager()
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range // 每一分钟去检查过期订单
		// 每25s 检查收款记录
		ticker.C {
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("ReconcileAccounts panic:", r)
					}
					CheckExpiredOrders()
				}()
			}()
		}
	}()
	go func() {
		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("ReconcileAccounts panic:", r)
					}
				}()
				if !global.GVA_CONFIG.System.UnMatchPayment {
					ReconcileAccounts()
				}
			}()
			time.Sleep(15 * time.Second)
		}
	}()
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
	// 获取机器人所有的channel
	var channels []bot.BotChannel
	if err = global.GVA_DB.Where("bot_id = ?", botID).Find(&channels).Error; err != nil {
		global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int64("botID", botID), zap.Error(err))
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

	var botHandler *bot_handler.Bot
	if botHandler, err = bot_handler.NewBot(botModel.Token); err != nil {
		global.GVA_LOG.Error("HandleAdConfirm NewBot", zap.Int64("botID", botID), zap.Error(err))
		return
	}
	trxResp.Data = append(trxResp.Data, utils.TronResponseData{
		TransactionID:  "mock_tx_1001",
		BlockTimestamp: 1764662681000, // 2025-11-30 22:40 北京时间
		From:           "TEST_FROM_ADDRESS",
		To:             "TKBDsYcVgvBMFi2qmhf88JDaMPYkqH8x2E",
		Type:           "Transfer",
		Value:          "20064000",
		TokenInfo: struct {
			Symbol   string `json:"symbol"`
			Address  string `json:"address"`
			Decimals int    `json:"decimals"`
			Name     string `json:"name"`
		}{
			Symbol:   "USDT",
			Address:  "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
			Decimals: 6,
			Name:     "Tether USD",
		},
	})
	var orders []recharge.UserRechargeRecord
	err = global.GVA_DB.Where("bot_id = ? AND status = 1", botID).Find(&orders).Error
	if err != nil {
		global.GVA_LOG.Error("查询订单失败", zap.Error(err))
		return
	}

	for _, order := range orders {
		for _, trx := range trxResp.Data {
			var match bool
			// 更新
			if match = rechargeSrv.MatchTransaction(paymentAddr, order, trx); match {
				if err := global.GVA_DB.Model(&recharge.UserRechargeRecord{}).
					Where("id = ? AND status = 1", order.ID).
					Updates(map[string]interface{}{
						"status":     constant.AdRechargePaid,
						"tx_id":      trx.TransactionID,
						"updated_at": time.Now(),
					}).Error; err != nil {
					global.GVA_LOG.Error("更新订单失败", zap.Error(err))
				}
				// 加钱
				var balance float64
				if balance, err = dao.RechargeDao.AddBalance(global.GVA_DB, botID, order.UserID, order.Price); err != nil {
					global.GVA_LOG.Error("reconcileAccount AddBalance", zap.Int64("botID", botID), zap.Int64("userID", order.UserID), zap.Float64("price", order.Price), zap.Error(err))
					continue
				}

				// 假如有缓存的广告 发布吧
				draftKey := cache.AdDraftConfirmCacheKey(botID, order.UserID, order.UpdateID)
				var publishContent string
				publishContent, err = global.GVA_REDIS.Get(context.Background(), draftKey).Result()
				if err != nil || publishContent == "" {
					global.GVA_LOG.Error("reconcileAccount DeleteMsg", zap.Int64("botID", botID), zap.Int64("userID", order.UserID), zap.Int64("updateID", order.UpdateID), zap.Error(err))
					continue
				} else {
					var medias []bot_handler.MediaItem
					if err = json.Unmarshal([]byte(publishContent), &medias); err != nil {
						global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int("botID", int(botID)), zap.Any("val", publishContent), zap.Error(err))
						return
					}

					rechargeCnfList := cache.NewRechargeCnfListCache(botID)
					if _, err = cache.CacheGetItem(rechargeCnfList); err != nil {
						global.GVA_LOG.Error("HandleAdConfirm RechargeCnfListCache", zap.Int64("botID", botID), zap.Int64("chatID", order.ChatID), zap.Int64("userID", order.UserID), zap.Error(err))
						continue
					}
					cnf, has := rechargeCnfList.WherePublishTimes(1)
					if !has {
						global.GVA_LOG.Error("HandleAdConfirm RechargeCnfListCache", zap.Int64("botID", botID), zap.Int64("chatID", order.ChatID), zap.Int64("userID", order.UserID))
						continue
					}

					if balance < order.Price {
						continue
					}

					// 余额充足 立马 扣减余额
					if _, err = dao.RechargeDao.ReduceBalance(global.GVA_DB, botID, order.UserID, cnf.Price); err != nil {
						global.GVA_LOG.Error("HandleAdConfirm ReduceBalance", zap.Int64("botID", botID), zap.Int64("userID", order.UserID), zap.Any("price", cnf.Price), zap.Error(err))
						continue
					}
					// 发布到所有渠道
					if err = botSrv.NewBotHandlerSvc(botID).PublishAd2Channel(*botHandler, order.ChatID, medias); err != nil {
						global.GVA_LOG.Error("botHandle PublishAd2Channel", zap.Int("botID", int(botID)), zap.Any("order", order.ID), zap.Error(err))
						continue
					}
				}
			}
		}
	}

	return
}
