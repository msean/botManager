package cmd

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/utils/bot_handler"
)

func BalanceShowHandle(update tgbotapi.Update, token string, botID int64) (err error) {
	wallet, getWalletErr := dao.RechargeDao.GetUserWallet(global.GVA_MYSQL, botID, bot_handler.GetChatUserID(update), "")
	if getWalletErr != nil {
		err = getWalletErr
		return
	}
	botApi, _ := bot_handler.NewBot(token)
	botApi.SendTextMessage(bot_handler.GetChatID(update), fmt.Sprintf("当前余额为：%.3fUSDT", wallet.Balance))
	return
}
