package cmd

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/utils/bot_handler"
)

func BalanceShowHandle(update tgbotapi.Update, token string, botID int64) (err error) {
	wallet, getWalletErr := dao.RechargeDao.GetUserWallet(global.GVA_DB, botID, getChatUserID(update), "")
	if getWalletErr != nil {
		err = getWalletErr
		return
	}
	botApi, _ := bot_handler.NewBot(token)
	botApi.SendTextMessage(getChatID(update), fmt.Sprintf("您的当前余额是：%.3f", wallet.Balance))
	return
}
