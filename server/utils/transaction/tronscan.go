package transaction

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type (
	AccountInfo struct {
		TRXBalance   float64
		USDTBalance  float64
		TotalTxCount int64
	}
	TxStats struct {
		FirstTime    int64
		TodayIn      float64
		TodayOut     float64
		YesterdayIn  float64
		YesterdayOut float64
	}
)

func GetAccountInfo(address string) (*AccountInfo, error) {

	client := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("https://apilist.tronscan.org/api/account?address=%s", address)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var acc struct {
		Balance               int64 `json:"balance"`
		TotalTransactionCount int64 `json:"totalTransactionCount"`
		Trc20TokenBalances    []struct {
			TokenAbbr string `json:"tokenAbbr"`
			Balance   string `json:"balance"`
		} `json:"trc20token_balances"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&acc); err != nil {
		return nil, err
	}

	info := &AccountInfo{
		TRXBalance:   float64(acc.Balance) / 1e6,
		TotalTxCount: acc.TotalTransactionCount,
	}

	// 找 USDT
	for _, t := range acc.Trc20TokenBalances {
		if t.TokenAbbr == "USDT" {
			val, _ := strconv.ParseFloat(t.Balance, 64)
			info.USDTBalance = val / 1e6
			break
		}
	}

	return info, nil
}

// 获取统计数据
func GetTxStats(address string) (*TxStats, error) {

	client := &http.Client{Timeout: 10 * time.Second}

	now := time.Now()
	loc := now.Location()

	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).UnixMilli()
	yesterdayStart := todayStart - 86400000

	var (
		start = 0
		limit = 200
		stats = &TxStats{}
	)

	for {
		url := fmt.Sprintf(
			"https://apilist.tronscan.org/api/transaction?address=%s&limit=%d&start=%d&sort=-timestamp",
			address, limit, start,
		)

		resp, err := client.Get(url)
		if err != nil {
			return nil, err
		}

		var tx struct {
			Data []struct {
				BlockTimestamp int64  `json:"block_timestamp"`
				OwnerAddress   string `json:"ownerAddress"`
				ToAddress      string `json:"toAddress"`
				Amount         int64  `json:"amount"`
			} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&tx); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		if len(tx.Data) == 0 {
			break
		}

		for _, t := range tx.Data {

			amount := float64(t.Amount) / 1e6

			// 激活时间
			if stats.FirstTime == 0 || t.BlockTimestamp < stats.FirstTime {
				stats.FirstTime = t.BlockTimestamp
			}

			// 今日
			if t.BlockTimestamp >= todayStart {
				if t.ToAddress == address {
					stats.TodayIn += amount
				} else if t.OwnerAddress == address {
					stats.TodayOut += amount
				}
			}

			// 昨日
			if t.BlockTimestamp >= yesterdayStart && t.BlockTimestamp < todayStart {
				if t.ToAddress == address {
					stats.YesterdayIn += amount
				} else if t.OwnerAddress == address {
					stats.YesterdayOut += amount
				}
			}
		}

		// 优化停止
		last := tx.Data[len(tx.Data)-1]
		if last.BlockTimestamp < yesterdayStart {
			break
		}

		start += limit
	}

	return stats, nil
}
