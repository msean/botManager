package cmd

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils/bot_handler"
	"go.uber.org/zap"
)

var (
	startCmd          = "/start"        // 开始按钮
	AdPublishCmd      = "/publishAd"    // 广告发布
	AdCancelCmd       = "/AdCancel"     // 取消广告发布
	AdRcvContentCmd   = "/AdRcvContent" // 收到广告内容
	AdConfirmCmd      = "/AdConfirm"    // 确认发布
	NoticeRechargeCmd = "/NoticeRecharge"
)

var (
	waitAdContentState = "waiting_publish_content" // 用户状态
)

var (
	waitAdContentExpire = 30 * time.Minute // 等待用户输入广告超时时间
	confirmAdExpire     = 30 * time.Minute // 确认广告超时有效时间
)

func getChatUserID(update tgbotapi.Update) (userId int64) {
	switch {
	case update.Message != nil:
		// 如果是用户发送的消息
		userId = int64(update.Message.From.ID)
	case update.CallbackQuery != nil:
		// 如果是用户点击了按钮
		userId = int64(update.CallbackQuery.From.ID)
	default:
		// 其他情况
		userId = 0
	}
	return
}

func Handle(update tgbotapi.Update, token string, botID int64) (err error) {
	var text string

	switch {
	case update.CallbackQuery != nil:
		text = update.CallbackQuery.Data
	case update.Message != nil:
		text = update.Message.Text
	}

	if update.CallbackQuery != nil {
		if err = HandleCallback(update.CallbackQuery, token, botID); err != nil {
			global.GVA_LOG.Error("Handle HandleCallback", zap.Int64("botID", botID), zap.Error(err))
			return
		}
	}

	cmds := cache.NewBotCmdCacheList(botID)
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

	cmdCfg := cache.NewBotCmdCache(botID, cmd, constant.BotReplyCmdType)
	if _, err = cache.CacheGetItem(cmdCfg); err != nil {
		global.GVA_LOG.Error("botHandle GetBotCmdCache", zap.Int("botID", int(botID)), zap.Error(err))
		return
	}

	global.GVA_LOG.Debug("BotMsgHandlerSvc handleCmd", zap.Any("cmd", cmd), zap.Any("any", cmdCfgMapper), zap.Any("triggerMapper", triggerMapper), zap.Any("cmd", cmd), zap.Any("inCfg", inCfg), zap.Any("cfg", cfg))
	switch cmd {
	case startCmd:
		StartHandlerfunc(update, token, *cmdCfg)
	default:
		if inCfg {
			switch cmd {
			case AdPublishCmd:
				if canPublich := PublishAdCheckHandle(update, token, botID); canPublich {
					SendCfgMessage(update, token, *cmdCfg, constant.ButtonTypeInline)
				}
			default:
				SendCfgMessage(update, token, *cmdCfg, constant.ButtonTypeInline)
			}

		}
	}
	global.GVA_LOG.Debug("BotMsgHandlerSvc ProcessBindCommand", zap.Any("cmd", cmd))
	ProcessBindCommand(update, token, botID, cmd)
	return
}

func ProcessBindCommand(update tgbotapi.Update, token string, botID int64, cmd string) {
	switch cmd {
	case AdPublishCmd: // 点击发布广告
		PublishAdHandle(update, token, botID)
	case AdRcvContentCmd: // 用户输入广告内容
		ReceiveAdContentHandle(update, token, botID)
	default:
	}
}

func ParseContentFromCfg(cfg cache.BotCmdCache, buttonType int) (markup any) {
	switch buttonType {
	case constant.ButtonTypeKeyBoard: // 普通键盘（ReplyKeyboard）
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

	case constant.ButtonTypeInline: // 内联键盘（InlineKeyboard）
		var rows [][]struct {
			Name    string `json:"name"`
			BindCmd string `json:"bindCmd"`
		}
		_ = json.Unmarshal([]byte(cfg.CmdButtons), &rows)

		inlineRows := make([][]tgbotapi.InlineKeyboardButton, 0)
		for _, row := range rows {
			btnRow := make([]tgbotapi.InlineKeyboardButton, 0)
			for _, b := range row {
				var btn tgbotapi.InlineKeyboardButton
				if strings.HasPrefix(b.BindCmd, "http://") || strings.HasPrefix(b.BindCmd, "https://") {
					btn = tgbotapi.NewInlineKeyboardButtonURL(b.Name, b.BindCmd)
				} else {
					btn = tgbotapi.NewInlineKeyboardButtonData(b.Name, b.BindCmd)
				}

				btnRow = append(btnRow, btn)
			}
			inlineRows = append(inlineRows, btnRow)
		}
		markup = tgbotapi.NewInlineKeyboardMarkup(inlineRows...)
	}
	return
}

func SendCfgMessage(update tgbotapi.Update, token string, cfg cache.BotCmdCache, buttonType int) error {
	markup := ParseContentFromCfg(cfg, buttonType)
	chatID := update.Message.Chat.ID // 获取聊天 ID
	global.GVA_LOG.Debug("BotMsgHandlerSvc send", zap.Any("any", cfg.Content), zap.Any("markup", markup))
	return bot_handler.HandleTexWithMarup(chatID, token, cfg.Content, markup)
}

func WaitCmd(update tgbotapi.Update, botID int64) string {

	state, _ := global.GVA_REDIS.Get(context.Background(), cache.AdWaitCacheKey(botID, getChatUserID(update))).Result()

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

	global.GVA_LOG.Debug("BotMsgHandlerSvc CallbackQuery", zap.Any("data", cb.Data))

	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return nil
	}

	cmd := parts[0]
	updateID, _ := strconv.Atoi(parts[1])

	global.GVA_LOG.Debug("BotMsgHandlerSvc CallbackQuery", zap.Any("updateID", updateID), zap.Any("cmd", cmd))

	switch cmd {
	case AdConfirmCmd:
		userName := cb.From.UserName
		if userName == "" {
			userName = cb.From.FirstName + " " + cb.From.LastName
		}
		return HandleAdConfirm(chatID, userID, userName, int64(updateID), token, botID, cb.Message.MessageID, 1)

	case AdCancelCmd:
		return HandleAdCancel(chatID, userID, updateID, token, botID, cb.Message.MessageID)
	}

	return nil
}
