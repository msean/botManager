package bot

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/system"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils/bot_handler"
	"go.uber.org/zap"
)

type BotMsgHandlerSvc struct{}

func (svc *BotMsgHandlerSvc) Handle(c *gin.Context, botID int, body []byte) (err error) {
	var tgMsg tgbotapi.Update
	if err = json.Unmarshal(body, &tgMsg); err != nil {
		global.GVA_LOG.Error("invalid telegram tgMsg", zap.Error(err))
		return
	}

	botModel, has, err := dao.BotDao.FromBotID(global.GVA_DB, botID)
	if err != nil || !has {
		global.GVA_LOG.Error("bot not found", zap.Int("botID", botID), zap.Error(err))
		return
	}

	if tgMsg.MyChatMember != nil {
		svc.SyncChatGroup(botModel, tgMsg)
		return nil
	}

	svc.SyncChatGroup(botModel, tgMsg)

	// 只要消息是转发的都需要禁止
	if tgMsg.Message.ForwardFrom != nil || tgMsg.Message.ForwardFromChat != nil {
		svc.BanUser(botModel, tgMsg, global.BanTypeForword)
		return nil
	}

	// 普通消息
	if tgMsg.Message == nil {
		return nil
	}

	var find bool
	if tgMsg.Message.Text != "" {
		if find, err = svc.CheckBanContent(botModel, tgMsg); err != nil || find {
			return
		}
	}

	if find, err = svc.CheckGroupMem(botModel, tgMsg); err != nil || find {
		return
	}

	return nil
}

func (svc *BotMsgHandlerSvc) BanUser(botModel bot.Bot, tgMsg tgbotapi.Update, _type int) (err error) {
	// 获取封禁时长
	var param system.SysParams
	param, err = dao.SysParamsDao.FromKey(global.GVA_DB, global.UserBanDuritonKey, global.DefaultUserBanDuriton)
	if err != nil {
		global.GVA_LOG.Error("get ban duration failed", zap.Error(err))
		return
	}
	durationMinutes, _ := strconv.Atoi(param.Value)

	// 发送api 封禁用户
	chatID := tgMsg.Message.Chat.ID
	messageID := tgMsg.Message.MessageID

	botHandler := bot_handler.NewBot(botModel.Token)
	var banErr error
	until := time.Now().Add(time.Duration(durationMinutes) * time.Minute).Unix()
	if banErr = botHandler.BanUser(tgMsg.Message.Chat.ID, tgMsg.Message.From.ID, until); err != nil {
		global.GVA_LOG.Error("ban user failed", zap.Error(err))
	} else {
		global.GVA_LOG.Info("ban user success",
			zap.Int64("chatID", tgMsg.Message.Chat.ID),
			zap.Int64("user_id", tgMsg.Message.Chat.ID),
			zap.Int64("util", until),
		)
	}

	if chatID != 0 && messageID != 0 && banErr != nil {
		if deleteErr := botHandler.DeleteMsg(chatID, messageID); deleteErr != nil {
			global.GVA_LOG.Error("delete msg error",
				zap.Int64("chatID", tgMsg.Message.Chat.ID),
				zap.Int64("user_id", tgMsg.Message.Chat.ID),
				zap.Int64("util", until),
				zap.Error(deleteErr),
			)
		}
	}

	remark := ""
	if banErr != nil {
		remark = banErr.Error()
	}
	msg := ""
	if _type == global.BanTypeWord {
		msg = tgMsg.Message.Text
	}
	record := bot.BanRecord{
		BotID:       botModel.BotID,
		UserID:      tgMsg.Message.From.ID,
		UserName:    tgMsg.Message.From.UserName,
		ChatID:      tgMsg.Message.Chat.ID,
		ChatName:    tgMsg.Message.Chat.Title,
		BanDuration: int64(durationMinutes),
		Remark:      remark,
		BanType:     _type,
		FullName:    fmt.Sprintf("%s%s", tgMsg.Message.From.FirstName, tgMsg.Message.From.LastName),
		Msg:         msg,
	}
	if err := global.GVA_DB.Create(&record).Error; err != nil {
		global.GVA_LOG.Error("failed to insert BanRecord", zap.Any("record", record), zap.Error(err))
	}
	return
}

func (svc *BotMsgHandlerSvc) CheckBanContent(botModel bot.Bot, tgMsg tgbotapi.Update) (find bool, err error) {
	// var banContents []bot.BotBanContent
	// banContents, err = dao.BotDao.ListBotBannerContentByID(global.GVA_DB, botModel.BotID)
	// if err != nil {
	// 	global.GVA_LOG.Error("fetch ban content failed", zap.Int("botID", botModel.BotID), zap.Error(err))
	// 	return
	// }
	if tgMsg.Message == nil {
		return
	}

	var botBanContentCache []cache.BotBanContentCache
	if _, err = cache.CacheGet(cache.BotBanContentCache{}.TableName(), cache.BotBanContentPk(botModel.BotID), &botBanContentCache, cache.LoadFromDBList); err != nil {
		global.GVA_LOG.Error("fetch ban content failed", zap.Int("botID", botModel.BotID), zap.Error(err))
		return
	}

	messageText := tgMsg.Message.Text
	global.GVA_LOG.Info("msg", zap.String("bot", botModel.Name), zap.String("msg", tgMsg.Message.Text))

	for _, rule := range botBanContentCache {
		if strings.Contains(messageText, rule.BanContent) {
			global.GVA_LOG.Info("found banned word",
				zap.String("word", rule.BanContent),
				zap.String("user", tgMsg.Message.From.UserName),
				zap.String("msg", messageText),
			)

			svc.BanUser(botModel, tgMsg, global.BanTypeWord)

			find = true
			return
		}
	}
	return
}

func (svc *BotMsgHandlerSvc) CheckGroupMem(botModel bot.Bot, tgMsg tgbotapi.Update) (found bool, err error) {
	if tgMsg.Message == nil {
		return
	}

	chatID := tgMsg.Message.Chat.ID
	user := tgMsg.Message.From

	var botChatGroupBanMemList []cache.BotChatGroupBanMemCache
	if _, err = cache.CacheGet(cache.BotChatGroupBanMemCache{}.TableName(), cache.BotChatGroupMemPk(botModel.BotID, int(chatID)), &botChatGroupBanMemList, cache.LoadFromDBList); err != nil {
		global.GVA_LOG.Error("fetch ban content failed", zap.Int("botID", botModel.BotID), zap.Error(err))
		return
	}

	fullName := fmt.Sprintf("%s%s", user.FirstName, user.LastName)

	for _, ban := range botChatGroupBanMemList {
		banStr := strings.ToLower(strings.TrimSpace(ban.BanMemContent))
		if banStr == "" {
			continue
		}

		if strings.Contains(fullName, banStr) {
			global.GVA_LOG.Info("found banned member",
				zap.String("ban_word", ban.BanMemContent),
				zap.String("member", user.UserName),
				zap.Int64("chatID", chatID),
			)

			svc.BanUser(botModel, tgMsg, global.BanTypeMem)

			found = true
			return
		}
	}
	return
}

func (svc *BotMsgHandlerSvc) SyncChatGroup(botModel bot.Bot, tgMsg tgbotapi.Update) {
	if tgMsg.Message == nil {
		return
	}
	chatID := tgMsg.Message.Chat.ID
	chatName := tgMsg.Message.Chat.Title

	if chatID == 0 || chatName == "" {
		return
	}

	// 构造缓存主键
	pkPairs := cache.BotChatGroupPk(botModel.BotID, int(chatID))

	// 尝试从缓存或数据库读取
	var chatGroupCache cache.BotChatGroupCache
	has, err := cache.CacheGet(chatGroupCache.TableName(), pkPairs, &chatGroupCache, cache.LoadFromDBGet)
	if err != nil {
		global.GVA_LOG.Error("CacheGet failed", zap.Int64("chatID", chatID), zap.Error(err))
		return
	}

	// 如果缓存或数据库不存在，创建新记录
	if !has {
		newGroup := bot.BotChatGroup{
			BotID:         botModel.BotID,
			ChatGroupID:   chatID,
			ChatGroupName: chatName,
		}
		if createErr := global.GVA_DB.Create(&newGroup).Error; createErr != nil {
			global.GVA_LOG.Error("failed to create new chat group",
				zap.Int64("chatID", chatID),
				zap.String("chatName", chatName),
				zap.Error(createErr),
			)
			return
		}
		global.GVA_LOG.Info("new chat group added",
			zap.Int64("chatID", chatID),
			zap.String("chatName", chatName),
		)
		return
	}

	// 如果数据库/缓存存在，但名称不同，则更新数据库和缓存
	if chatGroupCache.ChatGroupName != chatName {
		if err := global.GVA_DB.Model(&bot.BotChatGroup{}).
			Where("bot_id = ? AND chat_group_id = ?", botModel.BotID, chatID).
			Update("chat_group_name", chatName).Error; err != nil {
			global.GVA_LOG.Error("failed to update chat group name",
				zap.Int64("chatID", chatID),
				zap.String("newName", chatName),
				zap.Error(err),
			)
		} else {
			global.GVA_LOG.Info("chat group name updated",
				zap.Int64("chatID", chatID),
				zap.String("newName", chatName),
			)
		}

		if err = cache.CacheDelete(chatGroupCache.TableName(), pkPairs); err != nil {
			global.GVA_LOG.Error("failed to update chat group name",
				zap.Int64("chatID", chatID),
				zap.String("newName", chatName),
				zap.Error(err),
			)
		}
	}
}
