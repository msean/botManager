package collect

import (
	"context"
	"fmt"
	"log"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
	"github.com/msean/botmanager/server/model/tg_auto_helper"
	"gorm.io/gorm"
)

type TgCollector struct {
	Client *telegram.Client
	DB     *gorm.DB
}

// ================= 自定义 Handler（关键） =================

type MyUpdateHandler struct {
	t *TgCollector
}

func (h *MyUpdateHandler) Handle(ctx context.Context, u tg.UpdatesClass) error {

	switch upd := u.(type) {

	case *tg.Updates:
		return h.t.handleUpdates(ctx, upd.Updates, upd.Users)

	case *tg.UpdatesCombined:
		return h.t.handleUpdates(ctx, upd.Updates, upd.Users)
	}

	return nil
}

// ================= Run =================

func (t *TgCollector) Run(ctx context.Context) error {

	handler := &MyUpdateHandler{t: t}

	updateManager := updates.New(updates.Config{
		Handler: handler, // ✅ 必须是接口实现
	})

	return t.Client.Run(ctx, func(ctx context.Context) error {

		self, err := t.Client.Self(ctx)
		if err != nil {
			return err
		}

		return updateManager.Run(
			ctx,
			t.Client.API(),
			self.ID,
			updates.AuthOptions{},
		)
	})
}

// ================= 核心处理 =================

func (t *TgCollector) handleUpdates(
	ctx context.Context,
	updateList []tg.UpdateClass,
	users []tg.UserClass,
) error {

	userMap := make(map[int64]*tg.User)

	for _, u := range users {
		if user, ok := u.(*tg.User); ok {
			userMap[user.ID] = user
		}
	}

	for _, u := range updateList {

		switch update := u.(type) {

		case *tg.UpdateNewMessage:

			msg, ok := update.Message.(*tg.Message)
			if !ok {
				continue
			}

			groupID := getGroupID(msg)
			if groupID == 0 {
				continue
			}

			userID, username, nickname := getUserInfo(msg, userMap)
			if userID == 0 {
				continue
			}

			fmt.Println("采集用户:", userID, username, nickname)

			t.saveUser(groupID, userID, username, nickname)
		}
	}

	return nil
}

func getGroupID(msg *tg.Message) int64 {

	switch peer := msg.PeerID.(type) {

	case *tg.PeerChat:
		return peer.ChatID

	case *tg.PeerChannel:
		return peer.ChannelID
	}

	return 0
}

func getUserInfo(msg *tg.Message, userMap map[int64]*tg.User) (int64, string, string) {

	userPeer, ok := msg.FromID.(*tg.PeerUser)
	if !ok {
		return 0, "", ""
	}

	userID := userPeer.UserID

	user, ok := userMap[userID]
	if !ok {
		return userID, "", ""
	}

	return userID, user.Username, user.FirstName + user.LastName
}

func (t *TgCollector) saveUser(groupID, userID int64, username, nickname string) {

	user := tg_auto_helper.CollectUser{
		GroupID:  groupID,
		UserID:   userID,
		Username: username,
		NickName: nickname,
	}

	err := t.DB.
		Where("group_id = ? AND user_id = ?", groupID, userID).
		FirstOrCreate(&user).Error

	if err != nil {
		log.Println("保存用户失败:", err)
	}
}
