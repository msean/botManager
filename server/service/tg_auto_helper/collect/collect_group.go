package collect

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/tg_auto_helper"
)

type CollectGroupTask struct {
	tg_auto_helper.CollectGroupTask
}

func (t *CollectGroupTask) Run() error {

	var user tg_auto_helper.TgUser

	if err := global.GVA_MYSQL.
		Where("id = ?", t.TgUserID).
		First(&user).Error; err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	TaskManager.Add(t.ID, cancel)
	defer TaskManager.Remove(t.ID)

	client := telegram.NewClient(
		user.ApiId,
		user.ApiHash,
		telegram.Options{
			SessionStorage: &session.FileStorage{
				Path: user.SessionPath,
			},
		},
	)

	err := client.Run(ctx, func(ctx context.Context) error {

		api := client.API()

		bot, err := t.resolveSosoBot(ctx, api)
		if err != nil {
			return err
		}

		// 发送关键词
		if err := t.sendKeyword(ctx, api, bot); err != nil {
			return err
		}

		t.randomSleep(2, 4)

		return t.collectLoop(ctx, api, bot)
	})

	if err != nil {
		if err != nil {
			global.GVA_LOG.Error("task run error", zap.Error(err))
			global.GVA_MYSQL.Model(&tg_auto_helper.CollectGroupTask{}).
				Where("id = ?", t.ID).Updates(map[string]any{"status": 3, "remark": err.Error()})
		}
	}
	return err
}

func (t *CollectGroupTask) resolveSosoBot(
	ctx context.Context,
	api *tg.Client,
) (*tg.InputPeerUser, error) {

	res, err := api.ContactsResolveUsername(ctx, &tg.ContactsResolveUsernameRequest{
		Username: t.SourceBotName,
	})
	if err != nil {
		return nil, err
	}

	user := res.Users[0].(*tg.User)

	return &tg.InputPeerUser{
		UserID:     user.ID,
		AccessHash: user.AccessHash,
	}, nil
}

func (t *CollectGroupTask) sendKeyword(
	ctx context.Context,
	api *tg.Client,
	bot *tg.InputPeerUser,
) error {

	randomID := func() int64 {
		b := make([]byte, 8)
		_, _ = rand.Read(b)
		var id int64
		for _, v := range b {
			id = id<<8 | int64(v)
		}
		return id
	}

	_, err := api.MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
		Peer:     bot,
		Message:  t.SearchText,
		RandomID: randomID(),
	})

	return err
}

func (t *CollectGroupTask) collectLoop(
	ctx context.Context,
	api *tg.Client,
	bot *tg.InputPeerUser,
) error {

	parser := GetParser(t.SourceBotName)

	msg, err := t.waitBotMessage(ctx, api, bot)
	if err != nil {
		return err
	}

	// 点击群组按钮
	parser.AfterSearch(ctx, api, bot, msg)

	t.randomSleep(2, 4)

	for {

		msg, err = t.waitBotMessage(ctx, api, bot)
		if err != nil {
			return err
		}

		spew.Dump(">>>>>>>>>>>>msg2", msg)

		if t.isRiskControl(msg) {
			return errors.New("触发风控，请稍后再试")
		}

		current, total := parser.ParsePage(msg)

		groups := parser.ParseGroups(msg, t.ID)

		var filtered []tg_auto_helper.CollectGroupInfo
		for _, g := range groups {
			if g.Members >= t.MemberLower {
				filtered = append(filtered, g)
			}
		}

		t.saveGroups(filtered, current, total)

		spew.Dump("当前页:", current, "总页:", total)

		if current >= total {
			break
		}

		next := parser.FindNextButton(msg)
		if next == nil {
			break
		}

		err = t.clickNext(ctx, api, bot, msg.ID, next.Data)
		if err != nil {
			return err
		}

		// ✅ 防风控：随机延迟
		t.randomSleep(3, 6)
	}

	t.finishTask()
	return nil
}

func (t *CollectGroupTask) waitBotMessage(
	ctx context.Context,
	api *tg.Client,
	bot *tg.InputPeerUser,
) (*tg.Message, error) {

	for {
		history, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
			Peer:  bot,
			Limit: 1,
		})
		if err != nil {
			return nil, err
		}

		switch h := history.(type) {
		case *tg.MessagesMessages:
			if len(h.Messages) > 0 {
				return h.Messages[0].(*tg.Message), nil
			}
		case *tg.MessagesMessagesSlice:
			if len(h.Messages) > 0 {
				return h.Messages[0].(*tg.Message), nil
			}
		case *tg.MessagesChannelMessages:
			if len(h.Messages) > 0 {
				return h.Messages[0].(*tg.Message), nil
			}
		}

		time.Sleep(time.Second)
	}
}

//////////////////////////////////////////////////////
// ✅ 风控检测
//////////////////////////////////////////////////////

func (t *CollectGroupTask) isRiskControl(msg *tg.Message) bool {

	text := msg.Message

	keywords := []string{
		// "验证",
		// "请输入",
		"点击正确答案",
		// "机器人",
		"验证身份",
	}

	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}

	return false
}

func (t *CollectGroupTask) randomSleep(min, max int) {

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	sleep := int(n.Int64()) + min

	time.Sleep(time.Duration(sleep) * time.Second)
}

func (t *CollectGroupTask) saveGroups(
	groups []tg_auto_helper.CollectGroupInfo,
	currentPage int,
	totalPage int,
) {
	err := global.GVA_MYSQL.Model(&tg_auto_helper.CollectGroupTask{}).
		Where("id = ?", t.ID).
		Updates(map[string]interface{}{
			"current_page": currentPage,
			"total_page":   totalPage,
		}).Error

	if err != nil {
		global.GVA_LOG.Error("更新任务分页失败", zap.Error(err))
	}

	for _, g := range groups {

		var saveErr error

		saveErr = global.GVA_MYSQL.
			Where("link = ?", g.Link).
			Assign(map[string]interface{}{
				"title":   g.Title,
				"members": g.Members,
			}).
			FirstOrCreate(&g).Error

		if saveErr != nil {
			global.GVA_LOG.Error("CollectGroupTask saveGroups",
				zap.Any("group", g),
				zap.Error(saveErr),
			)
		}
	}
}

func (t *CollectGroupTask) clickNext(
	ctx context.Context,
	api *tg.Client,
	bot *tg.InputPeerUser,
	msgID int,
	data []byte,
) error {

	_, err := api.MessagesGetBotCallbackAnswer(
		ctx,
		&tg.MessagesGetBotCallbackAnswerRequest{
			Peer:  bot,
			MsgID: msgID,
			Data:  data,
		})

	return err
}

func (t *CollectGroupTask) finishTask() {
	global.GVA_MYSQL.Model(&tg_auto_helper.CollectGroupTask{}).
		Where("id = ?", t.ID).
		Update("status", 2)
}
