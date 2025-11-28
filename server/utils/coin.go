package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	trxLimiterOnce sync.Once
	trxLimiter     *rate.Limiter
	clientOnce     sync.Once
	httpClient     *http.Client
)

func getTronLimiter() *rate.Limiter {
	trxLimiterOnce.Do(func() {
		trxLimiter = rate.NewLimiter(rate.Limit(1), 1)
	})
	return trxLimiter
}

func getHTTPClient() *http.Client {
	clientOnce.Do(func() {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	})
	return httpClient
}

type (
	TronResponseData struct {
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
	}

	TronResponse struct {
		Data    []TronResponseData `json:"data"`
		Success bool               `json:"success"`
	}
)

func FetchTransactions(account string, limit int) (*TronResponse, error) {

	contract := "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	url := fmt.Sprintf(
		"https://api.trongrid.io/v1/accounts/%s/transactions/trc20?limit=%d&contract_address=%s",
		account, limit, contract,
	)

	limiter := getTronLimiter()
	client := getHTTPClient()

	var lastErr error

	for attempt := 1; attempt <= 3; attempt++ {

		limiter.Wait(context.Background())

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("http error: %v (attempt %d/3)", err, attempt)
			time.Sleep(time.Duration(attempt) * 300 * time.Millisecond)
			continue
		}

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("status %d from API (attempt %d/3)", resp.StatusCode, attempt)
			time.Sleep(time.Duration(attempt) * 300 * time.Millisecond)
			continue
		}

		var result TronResponse
		if err := json.Unmarshal(body, &result); err != nil {
			lastErr = fmt.Errorf("json parse error: %v", err)
			continue
		}

		if !result.Success {
			lastErr = fmt.Errorf("api returned success=false (attempt %d/3)", attempt)
			time.Sleep(time.Duration(attempt) * 300 * time.Millisecond)
			continue
		}

		return &result, nil
	}

	return nil, lastErr
}

func ParseAmount(value string, decimals int) float64 {
	var v float64
	fmt.Sscanf(value, "%f", &v)
	for i := 0; i < decimals; i++ {
		v /= 10
	}
	return v
}
