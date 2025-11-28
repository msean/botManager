package recharge

import (
	"fmt"
	"math/rand"

	"github.com/msean/botmanager/server/service/cache"
)

type (
	Pay struct {
		botID        int64 // 机器人ID
		publishTimes int   // 发布次数
	}
)

func NewPay(botID int64, publishTimes int) *Pay {
	return &Pay{botID: botID, publishTimes: publishTimes}
}

func (pay *Pay) RandomPrice() (float64, error) {
	rechargeCnfList := cache.NewRechargeCnfListCache(pay.botID)
	if _, err := cache.CacheGetItem(rechargeCnfList); err != nil {
		return 0, err
	}

	var base float64
	found := false

	for _, object := range rechargeCnfList.Objects {
		if object.PublishTimes == pay.publishTimes {
			base = object.Price
			found = true
			break
		}
	}

	if !found {
		return 0, fmt.Errorf("no config found for publishTimes=%d", pay.publishTimes)
	}

	second := rand.Intn(10)
	third := rand.Intn(10)

	randomDecimal := float64(second*10+third) / 1000.0

	return base + randomDecimal, nil
}
