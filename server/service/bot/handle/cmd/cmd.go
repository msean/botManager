package cmd

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/service/cache"
	"github.com/msean/botmanager/server/utils/bot_handler"
	"go.uber.org/zap"
)

func Handle(update tgbotapi.Update, token string, botID int64) (err error) {
	var text string

	switch {
	case update.CallbackQuery != nil:
		text = update.CallbackQuery.Data
	case update.Message != nil:
		text = update.Message.Text
	}

	if text == "" {
		return
	}

	cacheObjects := cache.NewBotCmdCacheList(int(botID))
	_, err = cache.CacheGetItem(cacheObjects)
	if err != nil {
		global.GVA_LOG.Error("fetch ban content failed", zap.Int("botID", int(botID)), zap.Error(err))
		return
	}

	var cfg cache.BotCmdCache
	var inCfg bool
	var cmd string

	// 真实命令 → cfg
	cmdCfgMapper := make(map[string]cache.BotCmdCache)

	// 用户输入 → bindCmd
	triggerMapper := make(map[string]string)

	// 构建映射
	for _, c := range cacheObjects.Objects {

		// 1. 主命令 → config
		cmdCfgMapper[c.Cmd] = c

		// 2. 按钮映射
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

	global.GVA_LOG.Debug("BotMsgHandlerSvc handleCmd", zap.Any("any", cmdCfgMapper), zap.Any("triggerMapper", triggerMapper), zap.Any("cmd", cmd), zap.Any("inCfg", inCfg), zap.Any("cfg", cfg))

	switch cmd {
	case "/start":
		StartHandlerfunc(update, token, cfg)
	default:
		if inCfg {
			SendCfgMessage(update, token, cfg, 2)
		}
	}
	ProcessBindCommand(update, token, botID, cmd)
	return
}

func ProcessBindCommand(update tgbotapi.Update, token string, botID int64, bindCmd string) {
	switch bindCmd {
	case "/publishAd":
	case "/showPrice":
	default:
	}
}

func SendCfgMessage(update tgbotapi.Update, token string, cfg cache.BotCmdCache, buttonType int) error {
	var markup interface{} // 最终传给 HandleTexWithMarup

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
