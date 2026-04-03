package tronscan

import (
	"fmt"
	"strconv"
)

type accountResp struct {
	Balance               int64 `json:"balance"`
	TotalTransactionCount int64 `json:"totalTransactionCount"`
	TRC20TokenBalances    []struct {
		TokenAbbr string `json:"tokenAbbr"`
		Balance   string `json:"balance"`
	} `json:"trc20token_balances"`
}

type AccountInfo struct {
	TRXBalance   float64
	USDTBalance  float64
	TotalTxCount int64
}

func GetAccountInfo(address string) (*AccountInfo, error) {

	url := fmt.Sprintf("https://apilist.tronscan.org/api/account?address=%s", address)

	var resp accountResp
	if err := doRequest(url, &resp); err != nil {
		return nil, err
	}

	info := &AccountInfo{
		TRXBalance:   float64(resp.Balance) / 1e6,
		TotalTxCount: resp.TotalTransactionCount,
	}

	for _, t := range resp.TRC20TokenBalances {
		if t.TokenAbbr == "USDT" {
			val, _ := strconv.ParseFloat(t.Balance, 64)
			info.USDTBalance = val / 1e6
		}
	}

	return info, nil
}
