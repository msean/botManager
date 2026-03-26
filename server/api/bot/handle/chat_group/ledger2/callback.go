package ledger2

import (
	"fmt"
	"os"
	"strings"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/ledger2"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

func handleCallback(botModel bot.Bot, update botapi.Update) error {

	if update.CallbackQuery == nil {
		return nil
	}

	data := update.CallbackQuery.Data

	if !strings.HasPrefix(data, "export_excel_") {
		return nil
	}

	date := strings.TrimPrefix(data, "export_excel_")

	chatID := update.CallbackQuery.Message.Chat.ID

	var list []ledger2.Ledger
	err := global.GVA_MYSQL.Where(
		"bot_id = ? AND chat_group_id = ? AND DATE(created_at) = ?",
		botModel.BotID,
		chatID,
		date,
	).Find(&list).Error

	if err != nil {
		return err
	}

	if len(list) == 0 {
		botSender, _ := botapi.NewBotAPI(botModel.Token)
		msg := botapi.NewMessage(chatID, "📊 "+date+" 暂无数据")
		_, _ = botSender.Send(msg)
		return nil
	}

	filePath := fmt.Sprintf("ledger_%s.csv", date)

	handler := LedgerSummaryHandler{
		botModel:    botModel,
		chatGroupID: chatID,
	}

	err = handler.generateCSV(list, filePath, date)
	if err != nil {
		return err
	}

	botSender, err := botapi.NewBotAPI(botModel.Token)
	if err != nil {
		return err
	}

	doc := botapi.NewDocument(chatID, botapi.FilePath(filePath))
	_, err = botSender.Send(doc)
	if err != nil {
		return err
	}

	_ = os.Remove(filePath)

	return nil
}
