package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TronResponse struct {
	Data []struct {
		TransactionID  string `json:"transaction_id"`
		BlockTimestamp int64  `json:"block_timestamp"`
		From           string `json:"from"`
		To             string `json:"to"`
		Type           string `json:"type"`
		Value          string `json:"value"`
		TokenInfo      struct {
			Symbol   string `json:"symbol"`
			Address  string `json:"address"`
			Decimals int    `json:"decimals"`
			Name     string `json:"name"`
		} `json:"token_info"`
	} `json:"data"`
	Success bool `json:"success"`
}

func FetchTransactions(account string, limit int) (*TronResponse, error) {
	contract := "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" // USDT 合约

	url := fmt.Sprintf(
		"https://api.trongrid.io/v1/accounts/%s/transactions/trc20?limit=%d&contract_address=%s",
		account, limit, contract,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result TronResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("json error: %v", err)
	}

	return &result, nil
}

func ParseAmount(value string, decimals int) float64 {
	var v float64
	fmt.Sscanf(value, "%f", &v)
	for i := 0; i < decimals; i++ {
		v /= 10
	}
	return v
}
