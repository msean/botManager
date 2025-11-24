package cache

import (
	"fmt"
)

func AdWaitCacheKey(botID, userID int64) string {
	return fmt.Sprintf("bot_manager:wait_state:%d_%d", botID, userID)
}

func AdDraftCacheKey(botID, userID, updateID int64) string {
	return fmt.Sprintf("bot_manager:ad_draft:%d_%d_%d", botID, userID, updateID)
}
