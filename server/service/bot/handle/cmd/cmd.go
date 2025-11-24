package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/recharge"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils/bot_handler"
	"go.uber.org/zap"
)

var (
	startCmd        = "/start"           // 开始按钮
	AdPublishCmd    = "/publishAd"       // 广告发布
	AdCancelCmd     = "/cancelPublishAd" // 取消广告发布
	AdRcvContentCmd = "/AdRcvContent"    // 收到广告内容
	AdConfirmCmd    = "/AdConfirm"       // 确认发布
)

var (
	waitAdContentState = "waiting_publish_content" // 用户状态
)

var (
	waitAdContentExpire = 30 * time.Minute // 等待用户输入广告超时时间
	confirmAdExpire     = 30 * time.Minute // 确认广告超时有效时间
)

func Handle(update tgbotapi.Update, token string, botID int64) (err error) {
	var text string

	switch {
	case update.CallbackQuery != nil:
		text = update.CallbackQuery.Data
	case update.Message != nil:
		text = update.Message.Text
	}

	if update.CallbackQuery != nil {
		return HandleCallback(update.CallbackQuery, token, botID)
	}

	if text == "" {
		return
	}

	cmds := cache.NewBotCmdCacheList(int(botID))
	if _, err = cache.CacheGetItem(cmds); err != nil {
		global.GVA_LOG.Error("botHandle GetBotCmdCacheList", zap.Int("botID", int(botID)), zap.Error(err))
		return
	}

	var cfg cache.BotCmdCacheWithNoContent
	var inCfg bool
	var cmd string

	cmdCfgMapper := make(map[string]cache.BotCmdCacheWithNoContent)
	triggerMapper := make(map[string]string)

	// 构建映射
	for _, c := range cmds.Objects {
		cmdCfgMapper[c.Cmd] = c
		if len(c.CmdButtons) > 0 {
			var buttons [][]struct {
				Name    string `json:"name"`
				BindCmd string `json:"bindCmd"`
			}
			_ = json.Unmarshal([]byte(c.CmdButtons), &buttons)
			for _, row := range buttons {
				for _, b := range row {
					// 输入按钮名 → 执行按钮绑定命令
					triggerMapper[b.Name] = b.BindCmd
				}
			}
		}
	}

	// 找到对应配置
	var ok bool
	if cfg, ok = cmdCfgMapper[text]; ok {
		inCfg = true
		cmd = text
	} else {
		bindCmd := triggerMapper[text]
		if bindCmd != "" {
			if cfg, ok = cmdCfgMapper[bindCmd]; ok {
				inCfg = true
				cmd = bindCmd
			}
		}
	}

	if cmd == "" {
		if cmd = WaitCmd(update, botID); cmd == "" {
			return
		}
	}

	cmdCfg := cache.NewBotCmdCache(botID, cmd)
	if _, err = cache.CacheGetItem(cmdCfg); err != nil {
		global.GVA_LOG.Error("botHandle GetBotCmdCache", zap.Int("botID", int(botID)), zap.Error(err))
		return
	}

	global.GVA_LOG.Debug("BotMsgHandlerSvc handleCmd", zap.Any("any", cmdCfgMapper), zap.Any("triggerMapper", triggerMapper), zap.Any("cmd", cmd), zap.Any("inCfg", inCfg), zap.Any("cfg", cfg))
	switch cmd {
	case startCmd:
		StartHandlerfunc(update, token, *cmdCfg)
	default:
		if inCfg {
			SendCfgMessage(update, token, *cmdCfg, 2)
		}
	}
	ProcessBindCommand(update, token, botID, cmd)
	return
}

func ProcessBindCommand(update tgbotapi.Update, token string, botID int64, bindCmd string) {
	switch bindCmd {
	case AdPublishCmd: // 点击发布广告
		PublishAdHandle(update, token, botID)
	case AdRcvContentCmd: // 用户输入广告内容
		ReceiveAdContentHandle(update, token, botID)
	default:
	}
}

func SendCfgMessage(update tgbotapi.Update, token string, cfg cache.BotCmdCache, buttonType int) error {
	var markup any // 最终传给 HandleTexWithMarup

	switch buttonType {
	case 1: // 普通键盘（ReplyKeyboard）
		var keyboard [][]tgbotapi.KeyboardButton

		if len(cfg.CmdButtons) > 0 {
			var buttons [][]struct {
				Name    string `json:"name"`
				BindCmd string `json:"bindCmd"`
			}
			_ = json.Unmarshal([]byte(cfg.CmdButtons), &buttons)

			for _, row := range buttons {
				kbRow := []tgbotapi.KeyboardButton{}
				for _, btn := range row {
					kbRow = append(kbRow, tgbotapi.NewKeyboardButton(btn.Name))
				}
				keyboard = append(keyboard, kbRow)
			}

			// 创建 ReplyKeyboardMarkup
			replyKeyboard := tgbotapi.ReplyKeyboardMarkup{
				Keyboard:        keyboard,
				ResizeKeyboard:  true,
				OneTimeKeyboard: false,
			}
			markup = replyKeyboard
		}

	case 2: // 内联键盘（InlineKeyboard）
		var rows [][]struct {
			Name    string `json:"name"`
			BindCmd string `json:"bindCmd"`
		}
		_ = json.Unmarshal([]byte(cfg.CmdButtons), &rows)

		inlineRows := make([][]tgbotapi.InlineKeyboardButton, 0)
		for _, row := range rows {
			btnRow := make([]tgbotapi.InlineKeyboardButton, 0)
			for _, b := range row {
				btnRow = append(btnRow, tgbotapi.NewInlineKeyboardButtonData(b.Name, b.BindCmd))
			}
			inlineRows = append(inlineRows, btnRow)
		}
		markup = tgbotapi.NewInlineKeyboardMarkup(inlineRows...)
	}

	chatID := update.Message.Chat.ID // 获取聊天 ID

	global.GVA_LOG.Debug("BotMsgHandlerSvc send", zap.Any("any", cfg.Content), zap.Any("markup", markup))
	return bot_handler.HandleTexWithMarup(chatID, token, cfg.Content, markup)
}

func WaitCmd(update tgbotapi.Update, botID int64) string {

	var userID int64

	cacheKey := fmt.Sprintf("bot:%d:user:%d:state", botID, userID)
	state, _ := global.GVA_REDIS.Get(context.Background(), cacheKey).Result()

	switch state {
	case waitAdContentState:
		return AdRcvContentCmd
	}
	return ""
}

func HandleCallback(cb *tgbotapi.CallbackQuery, token string, botID int64) error {
	chatID := cb.Message.Chat.ID
	userID := cb.From.ID
	data := cb.Data

	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return nil
	}

	cmd := parts[0]
	updateID, _ := strconv.Atoi(parts[1])

	switch cmd {
	case "AdConfirm":
		return HandleAdConfirm(chatID, userID, updateID, token, botID)

	case "AdCancel":
		return HandleAdCancel(chatID, userID, updateID, token, botID, cb.Message.MessageID)
	}

	return nil
}

func HandleAdCancel(chatID int64, userID int64, updateID int, token string, botID int64, msgID int) error {
	ctx := context.Background()

	draftKey := fmt.Sprintf("bot:%d:user:%d:ad_draft:%d", botID, userID, updateID)

	global.GVA_REDIS.Del(ctx, draftKey)

	bot, _ := tgbotapi.NewBotAPI(token)

	del := tgbotapi.NewDeleteMessage(chatID, msgID)
	bot.Send(del)

	bot.Send(tgbotapi.NewMessage(chatID, "❌ 已取消发布。"))

	return nil
}

func HandleAdConfirm(chatID int64, userID int64, updateID int, token string, botID int64) error {
	ctx := context.Background()
	draftKey := fmt.Sprintf("bot:%d:user:%d:ad_draft:%d", botID, userID, updateID)

	val, err := global.GVA_REDIS.Get(ctx, draftKey).Result()
	if err != nil || val == "" {
		bot_handler.SendTextMessage(chatID, token, "❌ 此发布请求已过期，请重新发送内容。")
		return nil
	}

	// 写入订单
	rec := recharge.UserRechargeRecord{
		BotID:           botID,
		PublishTimes:    1,
		StartTime:       time.Now(),
		PublishInterval: 30,
		PublishContent:  val,
		Status:          1, // 创建
	}
	global.GVA_DB.Create(&rec)

	global.GVA_REDIS.Del(ctx, draftKey)

	bot_handler.SendTextMessage(chatID, token, "✅ 广告订单创建成功，请前往后台完成支付。")
	return nil
}
