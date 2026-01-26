package private

import (
	"fmt"

	"github.com/msean/botmanager/server/dao"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/utils/bot_handler"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

func BalanceShowHandle(update botapi.Update, token string, botID int64) (err error) {
	wallet, getWalletErr := dao.RechargeDao.GetUserWallet(global.GVA_MYSQL, botID, bot_handler.GetChatUserID(update), "")
	if getWalletErr != nil {
		err = getWalletErr
		return
	}
	botApi, _ := bot_handler.NewBot(token)
	botApi.SendTextMessage(bot_handler.GetChatID(update), fmt.Sprintf("当前余额为：%.3fUSDT", wallet.Balance))
	return
}
