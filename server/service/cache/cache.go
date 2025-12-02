package cache

import (
	"fmt"
)

func AdWaitCacheKey(botID, userID int64) string {
	return fmt.Sprintf("bot_manager:wait_state:%d_%d", botID, userID)
}

// 用户输入的广告内容
func AdDraftCacheKey(botID, userID, updateID int64) string {
	return fmt.Sprintf("bot_manager:ad_draft:%d_%d_%d", botID, userID, updateID)
}

func RechargeTryCountKey(botID int64, userID int64) string {
	return fmt.Sprintf("bot_manager:recharge:try:%d:%d", botID, userID)
}
