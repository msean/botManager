package collect

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/davecgh/go-spew/spew"
	"github.com/gotd/td/tg"
	"github.com/msean/botmanager/server/model/tg_auto_helper"
)

type CollectSoSo struct{}

func (c *CollectSoSo) AfterSearch(
	ctx context.Context,
	api *tg.Client,
	bot *tg.InputPeerUser,
	msg *tg.Message,
) error {
	if msg.ReplyMarkup == nil {
		return nil
	}

	markup := msg.ReplyMarkup.(*tg.ReplyInlineMarkup)

	for _, row := range markup.Rows {

		for _, btn := range row.Buttons {

			b, ok := btn.(*tg.KeyboardButtonCallback)
			if !ok {
				continue
			}

			if strings.Contains(b.Text, "👥") || strings.Contains(b.Text, "群组") {
				_, err := api.MessagesGetBotCallbackAnswer(
					ctx,
					&tg.MessagesGetBotCallbackAnswerRequest{
						Peer:  bot,
						MsgID: msg.ID,
						Data:  b.Data,
					},
				)
				return err
			}
		}
	}
	return nil
}

func utf16Len(s string) int {
	return len(utf16.Encode([]rune(s)))
}

func (c *CollectSoSo) ParseGroups(msg *tg.Message, taskID uint) []tg_auto_helper.CollectGroupInfo {

	var result []tg_auto_helper.CollectGroupInfo

	text := msg.Message

	reg := regexp.MustCompile(`👥\s*(.*?)\s*([0-9.]+)\s*K`)

	type LineInfo struct {
		Title   string
		Members int
		Start   int // UTF16 offset
		End     int // UTF16 offset
	}

	var linesInfo []LineInfo

	// =========================
	// ✅ 用 UTF-16 计算 offset
	// =========================
	offset := 0

	for _, line := range strings.Split(text, "\n") {

		lineLen := utf16Len(line)

		if strings.HasPrefix(line, "👥") {

			match := reg.FindStringSubmatch(line)
			if len(match) >= 3 {

				title := strings.TrimSpace(match[1])

				f, _ := strconv.ParseFloat(match[2], 64)
				members := int(f * 1000)

				linesInfo = append(linesInfo, LineInfo{
					Title:   title,
					Members: members,
					Start:   offset,
					End:     offset + lineLen,
				})
				spew.Dump(">>>>>>>>title", title, members, offset, offset+lineLen)
			}
		}

		// ⚠️ 换行也算一个 UTF16
		offset += lineLen + 1
	}

	// =========================
	// 2️⃣ 用 entity offset 匹配
	// =========================
	for _, e := range msg.Entities {

		entity, ok := e.(*tg.MessageEntityTextURL)
		if !ok {
			continue
		}

		url := entity.URL

		if !isValidGroupLink(url) {
			continue
		}

		for _, line := range linesInfo {

			if entity.Offset >= line.Start && entity.Offset <= line.End {

				result = append(result, tg_auto_helper.CollectGroupInfo{
					Title:   line.Title,
					Link:    url,
					Members: line.Members,
					TaskID:  taskID,
				})
				spew.Dump(">>>>>>>>title", line.Title, url)
				break
			}
		}
	}

	return result
}

func isValidGroupLink(url string) bool {

	if !strings.Contains(url, "t.me/") {
		return false
	}

	// ❌ soso跳转
	if strings.Contains(url, "t.me/soso") {
		return false
	}

	// ❌ 分享
	if strings.Contains(url, "share/url") {
		return false
	}

	// ❌ start参数（广告）
	if strings.Contains(url, "start=") {
		return false
	}

	// ❌ sponsor 常见（你这个坑）
	if strings.Contains(url, "laicai") {
		return false
	}

	return true
}

func (c *CollectSoSo) FindNextButton(
	msg *tg.Message,
) *tg.KeyboardButtonCallback {

	if msg.ReplyMarkup == nil {
		return nil
	}

	markup := msg.ReplyMarkup.(*tg.ReplyInlineMarkup)

	for _, row := range markup.Rows {

		for _, btn := range row.Buttons {

			b, ok := btn.(*tg.KeyboardButtonCallback)

			if !ok {
				continue
			}

			if strings.Contains(b.Text, "下一页") {

				return b
			}
		}
	}

	return nil
}

func (c *CollectSoSo) ParsePage(
	msg *tg.Message,
) (int, int) {

	reg := regexp.MustCompile(`当前第(\d+)页.*共(\d+)页`)

	match := reg.FindStringSubmatch(msg.Message)

	if len(match) < 3 {

		return 1, 1
	}

	current, _ := strconv.Atoi(match[1])
	total, _ := strconv.Atoi(match[2])

	return current, total
}
