package trongrid

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	limiterOnce sync.Once
	trxLimiter  *rate.Limiter
	clientOnce  sync.Once
	httpClient  *http.Client
)

func getTronLimiter() *rate.Limiter {
	limiterOnce.Do(func() {
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

func doRequest(url string, result interface{}) error {
	limiter := getTronLimiter()
	client := getHTTPClient()

	var lastErr error

	for i := 0; i < 3; i++ { // ✅ 重试3次
		if err := limiter.Wait(context.Background()); err != nil {
			return err
		}

		resp, err := client.Get(url)
		if err != nil {
			lastErr = err
			time.Sleep(time.Millisecond * 300)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != 200 {
			lastErr = errors.New("http error")
			time.Sleep(time.Millisecond * 300)
			continue
		}

		if err := json.Unmarshal(body, result); err != nil {
			lastErr = err
			continue
		}

		return nil
	}

	return lastErr
}
