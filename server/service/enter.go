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
	"github.com/msean/botmanager/server/service/ledger"
	"github.com/msean/botmanager/server/service/listen"
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
	UsageServiceGroup    ledger.ServiceGroup
	ListenServiceGroup   listen.ServiceGroup
}

func Init() {
	botSrv.InitBotTaskManager()
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range // 每一分钟去检查过期订单
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
			time.Sleep(20 * time.Second)
		}
	}()
}
func CheckExpiredOrders() {
	deadline := time.Now().Add(-constant.OrderMatchAgo * time.Minute)
	err := global.GVA_MYSQL.Model(&recharge.UserRechargeRecord{}).Where("status = ?", constant.AdRechargeCreate).Where("created_at <= ?", deadline).Updates(map[string]any{"status": constant.AdRechargeTimeout, "updated_at": time.Now()}).Error
	if err != nil {
		global.GVA_LOG.Error("订单超时更新失败", zap.Error(err))
	}
}
func ReconcileAccounts() {
	db := global.GVA_MYSQL
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
	var channels []bot.BotChannel
	if err = global.GVA_MYSQL.Where("bot_id = ?", botID).Find(&channels).Error; err != nil {
		global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int64("botID", botID), zap.Error(err))
		return
	}
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
	trxResp.Data = append(trxResp.Data, utils.TronResponseData{TransactionID: "mock_tx_1001", BlockTimestamp: 1764746290000, From: "TEST_FROM_ADDRESS", To: "TKBDsYcVgvBMFi2qmhf88JDaMPYkqH8x2E", Type: "Transfer", Value: "10004000", TokenInfo: struct {
		Symbol   string `json:"symbol"`
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		Name     string `json:"name"`
	}{Symbol: "USDT", Address: "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", Decimals: 6, Name: "Tether USD"}})
	var orders []recharge.UserRechargeRecord
	err = global.GVA_MYSQL.Where("bot_id = ? AND status = 1", botID).Find(&orders).Error
	if err != nil {
		global.GVA_LOG.Error("查询订单失败", zap.Error(err))
		return
	}
	for _, order := range orders {
		for _, trx := range trxResp.Data {
			var match bool
			if match = rechargeSrv.MatchTransaction(paymentAddr, order, trx); match {
				if err := global.GVA_MYSQL.Model(&recharge.UserRechargeRecord{}).Where("id = ? AND status = 1", order.ID).Updates(map[string]interface{}{"status": constant.AdRechargePaid, "tx_id": trx.TransactionID, "updated_at": time.Now()}).Error; err != nil {
					global.GVA_LOG.Error("更新订单失败", zap.Error(err))
				}
				var balance float64
				if balance, err = dao.RechargeDao.AddBalance(global.GVA_MYSQL, botID, order.UserID, order.Price); err != nil {
					global.GVA_LOG.Error("reconcileAccount AddBalance", zap.Int64("botID", botID), zap.Int64("userID", order.UserID), zap.Float64("price", order.Price), zap.Error(err))
					continue
				}
				draftKey := cache.AdDraftConfirmCacheKey(botID, order.UserID)
				var publishContent string
				publishContent, err = global.GVA_REDIS.Get(context.Background(), draftKey).Result()
				if err != nil || publishContent == "" {
					global.GVA_LOG.Error("reconcileAccount DeleteMsg", zap.Int64("botID", botID), zap.Int64("userID", order.UserID), zap.Int("msgID", order.MsgID), zap.Error(err))
					continue
				} else {
					var medias []bot_handler.MediaItem
					if err = json.Unmarshal([]byte(publishContent), &medias); err != nil {
						global.GVA_LOG.Error("botHandle HandleAdConfirm", zap.Int("botID", int(botID)), zap.Any("val", publishContent), zap.Error(err))
						continue
					}
					var cnf cache.RechargeCnfObj
					var has bool
					if cnf, has, err = cache.NewRechargeCnfListCache(botID).WherePublishTimes(order.PublishTimes); !has || err != nil {
						global.GVA_LOG.Error("HandleAdConfirm RechargeCnfListCache", zap.Int64("botID", botID), zap.String("userName", order.UserName), zap.Int64("userID", order.UserID), zap.Bool("has", has), zap.Error(err))
						continue
					}
					if balance < order.Price {
						continue
					}
					hook := func(channels []cache.BotChannelCache) error {
						go func() {
							var channelIDList []int64
							for _, channel := range channels {
								channelIDList = append(channelIDList, channel.ChannelID)
							}
							err := dao.RechargeDao.CreatePublishRecords(global.GVA_MYSQL, recharge.AdPublishRecord{BotID: botID, PublishTimes: 1, UserID: order.UserID, UserName: order.UserName, Price: cnf.Price, Content: order.PublishContent}, channelIDList)
							if err != nil {
								global.GVA_LOG.Error("保存发布记录失败", zap.Error(err))
							}
						}()
						return nil
					}
					if err = botSrv.NewBotHandlerSvc(botID).PublishAd2Channel(*botHandler, order.ChatID, medias, hook); err != nil {
						global.GVA_LOG.Error("botHandle PublishAd2Channel", zap.Int("botID", int(botID)), zap.Any("order", order.ID), zap.Error(err))
						continue
					}
					if err = global.GVA_REDIS.Del(context.Background(), draftKey).Err(); err != nil {
						global.GVA_LOG.Error("botHandle global.GVA_REDIS.Del", zap.String("draftKey", draftKey), zap.Any("order", order.ID), zap.Error(err))
						continue
					}
					if _, err = dao.RechargeDao.ReduceBalance(global.GVA_MYSQL, botID, order.UserID, cnf.Price); err != nil {
						global.GVA_LOG.Error("HandleAdConfirm ReduceBalance", zap.Int64("botID", botID), zap.Int64("userID", order.UserID), zap.Any("price", cnf.Price), zap.Error(err))
						continue
					}
				}
			}
		}
	}
	return
}
