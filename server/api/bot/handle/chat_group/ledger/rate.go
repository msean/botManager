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

	rate, err := getTop10USDTToCNY()
	if err != nil {
		return r.reply("❌ 获取汇率失败，请稍后再试")
	}

	text := fmt.Sprintf(
		"💱 当前汇率（欧易）\n\nUSDT ≈ %.2f CNY\n\n⏰ 更新时间：%s",
		rate,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return r.reply(text)
}

func getTop10USDTToCNY() ([]float64, error) {
	url := "https://www.okx.com/v3/c2c/tradingOrders/books?quoteCurrency=cny&baseCurrency=usdt&side=sell"

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 结构体（只解析需要的字段）
	var result struct {
		Data []struct {
			Price string `json:"price"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no data")
	}

	// ✅ 取前10档
	limit := 10
	if len(result.Data) < 10 {
		limit = len(result.Data)
	}

	prices := make([]float64, 0, limit)

	for i := 0; i < limit; i++ {
		var price float64
		fmt.Sscanf(result.Data[i].Price, "%f", &price)
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
