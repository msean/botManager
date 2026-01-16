package private

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/utils/bot_handler"
)

func getTelegramUser(update tgbotapi.Update) *tgbotapi.User {
	if update.Message != nil {
		return update.Message.From
	}
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From
	}
	return nil
}

func UserCenterHandler(update tgbotapi.Update, token string, botID int64) (err error) {
	user := getTelegramUser(update)
	if user == nil {
		return fmt.Errorf("无法获取 Telegram 用户信息")
	}

	userID := user.ID
	username := user.UserName
	fullName := strings.TrimSpace(user.FirstName + " " + user.LastName)

	// 没有 username 的兜底
	if username == "" {
		username = "未设置"
	}
	if fullName == "" {
		fullName = "未知"
	}

	// 查询钱包
	wallet, getWalletErr := dao.RechargeDao.GetUserWallet(
		global.GVA_MYSQL,
		botID,
		userID,
		"",
	)
	if getWalletErr != nil {
		return getWalletErr
	}

	botApi, _ := bot_handler.NewBot(token)

	text := fmt.Sprintf(
		"🆔 用户ID：%d\n"+
			"👤 姓名：%s\n"+
			"🔗 用户名：@%s\n"+
			"💰 USDT余额：%.3f",
		userID,
		fullName,
		username,
		wallet.Balance,
	)

	botApi.SendTextMessage(
		bot_handler.GetChatID(update),
		text,
	)

	return nil
}
