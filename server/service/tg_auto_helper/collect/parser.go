package collect

import (
	"context"

	"github.com/gotd/td/tg"
	"github.com/msean/botmanager/server/model/tg_auto_helper"
)

type CollectParser interface {

	// 发送关键词后的处理（例如点击群组按钮）
	AfterSearch(
		ctx context.Context,
		api *tg.Client,
		bot *tg.InputPeerUser,
		msg *tg.Message,
	) error

	// 解析群
	ParseGroups(
		msg *tg.Message,
		taskID uint,
	) []tg_auto_helper.CollectGroupInfo

	// 查找下一页按钮
	FindNextButton(
		msg *tg.Message,
	) *tg.KeyboardButtonCallback

	// 解析页码
	ParsePage(
		msg *tg.Message,
	) (int, int)
}

func GetParser(name string) CollectParser {

	switch name {

	case "sosoNewBot":
		return &CollectSoSo{}

	default:
		return &CollectSoSo{}
	}
}
