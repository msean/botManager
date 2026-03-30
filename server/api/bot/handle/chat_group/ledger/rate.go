package ledger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botapi "github.com/msean/botmanager/server/utils/bot_handler/bot_api"
	"go.uber.org/zap"
)

type RateHandler struct {
	botModel    bot.Bot
	chatGroupID int64
}

func (r *RateHandler) Match(botModel bot.Bot, update botapi.Update) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := strings.TrimSpace(update.Message.Text)

	if text != "hl" && text != "汇率" {
		return false
	}

	r.botModel = botModel
	r.chatGroupID = update.Message.Chat.ID

	return true
}

func (r *RateHandler) Handle() error {

	prices, err := getTop10USDTToCNY()
	if err != nil {
		global.GVA_LOG.Error("getTop10USDTToCNY", zap.Error(err))
		return r.reply("❌ 获取汇率失败")
	}

	// _, _ = getTop10USDTToCNY()

	text := "💱 USDT 场外汇率（OKX）\n\n"

	for i, p := range prices {
		text += fmt.Sprintf("第%d档：%.2f\n", i+1, p)
	}

	// text += fmt.Sprintf("\n📊 平均价：%.2f", avg)

	return r.reply(text)
}

func getTop10USDTToCNY() ([]float64, error) {
	url := "https://www.okx.com/v3/c2c/tradingOrders/books?quoteCurrency=cny&baseCurrency=usdt&side=sell&page=1&size=20"

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, _ := http.NewRequest("GET", url, nil)

	// ✅ 非常关键：伪装浏览器（否则返回空）
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Code int `json:"code"`
		Data struct {
			Sell []struct {
				Price string `json:"price"`
			} `json:"sell"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// 🔍 调试用（建议你先开）
	if len(result.Data.Sell) == 0 {
		return nil, fmt.Errorf("no sell data, raw: %s", string(body))
	}

	limit := 10
	if len(result.Data.Sell) < 10 {
		limit = len(result.Data.Sell)
	}

	prices := make([]float64, 0, limit)

	for i := 0; i < limit; i++ {
		price, _ := strconv.ParseFloat(result.Data.Sell[i].Price, 64)
		prices = append(prices, price)
	}

	return prices, nil
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
