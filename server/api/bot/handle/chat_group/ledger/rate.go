package ledger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
)

type RateHandler struct {
	botModel    bot.Bot
	chatGroupID int64
	text        string
}

func (r *RateHandler) Match(botModel bot.Bot, update botapi.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)
	text = strings.ToUpper(text)

	if text != "Q" && text != "K" && text != "Z" && text != "V" {
		return false
	}

	r.botModel = botModel
	r.chatGroupID = update.Message.Chat.ID
	r.text = text

	return true
}

func (r *RateHandler) Handle() error {

	now := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("欧易 USDT 实时汇率 %s\n", now)

	switch r.text {
	case "Q":
		msg += r.getUSDTData("bank")
		msg += r.getUSDTData("aliPay")
		msg += r.getUSDTData("wxPay")
	case "K":
		msg += r.getUSDTData("bank")
	case "Z":
		msg += r.getUSDTData("aliPay")
	case "V":
		msg += r.getUSDTData("wxPay")
	}

	return r.reply(msg)
}

func (r *RateHandler) getUSDTData(payType string) string {

	typeMap := map[string]string{
		"bank":   "银行卡",
		"aliPay": "支付宝",
		"wxPay":  "微信",
	}

	typeCN := typeMap[payType]
	msg := fmt.Sprintf("\n类型： %s\n", typeCN)

	url := fmt.Sprintf("https://www.okx.com/v3/c2c/tradingOrders/books?quoteCurrency=CNY&baseCurrency=USDT&side=sell&paymentMethod=%s&page=1&rows=10", payType)

	client := &http.Client{Timeout: 10 * time.Second}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return msg + "请求失败\n\n"
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Code int `json:"code"`
		Data struct {
			Sell []struct {
				Price    string `json:"price"`
				NickName string `json:"nickName"`
			} `json:"sell"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return msg + "解析失败\n"
	}

	if result.Code != 0 || len(result.Data.Sell) == 0 {
		return msg + "暂无数据\n"
	}

	// 取前10
	list := result.Data.Sell
	limit := 10
	if len(list) < 10 {
		limit = len(list)
	}

	for i := 0; i < limit; i++ {
		item := list[i]
		msg += fmt.Sprintf("%s           %s\n", item.Price, item.NickName)
	}

	msg += "\n"
	return msg
}

func (r *RateHandler) reply(text string) error {
	botSender, err := botapi.NewBotAPI(r.botModel.Token)
	if err != nil {
		return err
	}

	msg := botapi.NewMessage(r.chatGroupID, text)
	_, err = botSender.Send(msg)
	return err
}
