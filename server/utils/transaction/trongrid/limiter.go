package trongrid

import (
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
