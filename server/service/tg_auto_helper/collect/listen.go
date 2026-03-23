package collect

import (
	"context"
	"fmt"
	"log"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/msean/botmanager/server/model/tg_auto_helper"
	"gorm.io/gorm"
)

// ================= 核心服务 =================

type TgCollector struct {
	Client *telegram.Client
	DB     *gorm.DB
}

// ================= 启动监听 =================

func (t *TgCollector) Run(ctx context.Context) error {

	updateManager := updates.New(updates.Config{
		Handler: updates.HandlerFunc(func(ctx context.Context, u tg.UpdatesClass) error {

			switch update := u.(type) {

			case *tg.UpdateNewMessage:

				msg, ok := update.Message.(*tg.Message)
				if !ok {
					return nil
				}

				// 只处理群 / 超级群
				groupID := getGroupID(msg)
				if groupID == 0 {
					return nil
				}

				// 获取发送人
				userID, username, nickname := getUserInfo(msg)
				if userID == 0 {
					return nil
				}

				fmt.Println("采集用户:", userID, username, nickname)

				// 入库
				t.saveUser(groupID, userID, username, nickname)

			}

			return nil
		}),
	})

	return t.Client.Run(ctx, func(ctx context.Context) error {
		return updateManager.Run(ctx, t.Client.API(), t.Client.Self())
	})
}

// ================= 获取群ID =================

func getGroupID(msg *tg.Message) int64 {

	switch peer := msg.PeerID.(type) {

	case *tg.PeerChat:
		return peer.ChatID

	case *tg.PeerChannel:
		return peer.ChannelID

	default:
		return 0
	}
}

// ================= 获取用户信息 =================

func getUserInfo(msg *tg.Message) (int64, string, string) {

	userPeer, ok := msg.FromID.(*tg.PeerUser)
	if !ok {
		return 0, "", ""
	}

	userID := userPeer.UserID

	// 从 Users 里找详细信息
	for _, u := range msg.Entities.Users() {
		user, ok := u.(*tg.User)
		if !ok {
			continue
		}

		if user.ID == userID {
			return userID, user.Username, user.FirstName + user.LastName
		}
	}

	return userID, "", ""
}

// ================= 保存用户 =================

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
