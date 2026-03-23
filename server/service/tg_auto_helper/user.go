package tg_auto_helper

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"

	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/global/constant"
	"github.com/msean/botmanager/server/model/tg_auto_helper"
	tgAutoHelperReq "github.com/msean/botmanager/server/model/tg_auto_helper/request"
)

type TgUserService struct{}

func (s *TgUserService) Create(user *tg_auto_helper.TgUser) error {
	user.Status = constant.TgStatusInit
	user.SessionPath = fmt.Sprintf("%s/%s.session", global.GVA_CONFIG.Local.SessionPath, user.Phone)
	return global.GVA_MYSQL.Create(user).Error
}

// ================= 创建 Client =================

func (s *TgUserService) newClient(user *tg_auto_helper.TgUser) *telegram.Client {
	return telegram.NewClient(
		user.ApiId,
		user.ApiHash,
		telegram.Options{
			SessionStorage: &session.FileStorage{
				Path: user.SessionPath,
			},
		},
	)
}

// ================= 发送验证码 =================
func (s *TgUserService) SendCode(ctx context.Context, user *tg_auto_helper.TgUser) error {
	tgClient := s.newClient(user)

	authClient := auth.NewClient(
		tgClient.API(), // ✅ 高层 client -> 底层 *tg.Client
		rand.Reader,    // crypto/rand.Reader
		user.ApiId,
		user.ApiHash,
	)

	return tgClient.Run(ctx, func(ctx context.Context) error {
		resp, err := authClient.SendCode(ctx, user.Phone, auth.SendCodeOptions{})

		if err != nil {
			return err
		}
		sentCode := resp.(*tg.AuthSentCode)

		// gotd v1.9+ 已经不需要存 PhoneCodeHash
		user.PhoneCodeHash = sentCode.PhoneCodeHash
		user.Status = constant.TgStatusCodeSent
		return global.GVA_MYSQL.Save(user).Error
	})
}

// ================= 验证验证码 =================
func (s *TgUserService) VerifyCode(
	ctx context.Context,
	user *tg_auto_helper.TgUser,
	code string,
) error {

	client := s.newClient(user)

	authClient := auth.NewClient(
		client.API(),
		rand.Reader,
		user.ApiId,
		user.ApiHash,
	)

	return client.Run(ctx, func(ctx context.Context) error {

		_, err := authClient.SignIn(
			ctx,
			user.Phone,
			code,
			user.PhoneCodeHash,
		)

		if err != nil {
			// ===== 需要二步验证 =====
			if errors.Is(err, auth.ErrPasswordAuthNeeded) {

				if user.NextVerification == "" {
					user.Status = constant.TgStatusPasswordRequired
					global.GVA_MYSQL.Save(user)
					return errors.New("账号需要二步验证密码")
				}

				// 自动执行二步验证
				if _, err := authClient.Password(ctx, user.NextVerification); err != nil {
					return err
				}

				user.Status = constant.TgStatusReady
				return global.GVA_MYSQL.Save(user).Error
			}

			return err
		}

		user.Status = constant.TgStatusReady
		return global.GVA_MYSQL.Save(user).Error
	})
}

// ================= 二步验证 =================
// func (s *TgUserService) VerifyPassword(
// 	ctx context.Context,
// 	user *tg_auto_helper.TgUser,
// 	password string,
// ) error {

// 	client := s.newClient(user)
// 	authClient := auth.NewClient(
// 		client.API(), // 这里 API() 返回 *tg.Client
// 		rand.Reader,  // crypto/rand.Reader
// 		user.ApiId,   // appID
// 		user.ApiHash, // appHash
// 	)
// 	return client.Run(ctx, func(ctx context.Context) error {
// 		if _, err := authClient.Password(ctx, password); err != nil {
// 			return err
// 		}

// 		user.Status = constant.TgStatusReady
// 		return global.GVA_MYSQL.Save(user).Error
// 	})
// }

//
// ================= CRUD =================
//

func (s *TgUserService) DeleteTgUser(ctx context.Context, ID string) error {
	return global.GVA_MYSQL.Delete(&tg_auto_helper.TgUser{}, "id = ?", ID).Error
}

func (s *TgUserService) UpdateTgUser(ctx context.Context, user tg_auto_helper.TgUser) error {
	user.SessionPath = fmt.Sprintf("%s/%s.session", global.GVA_CONFIG.Local.SessionPath, user.Phone)

	return global.GVA_MYSQL.Model(&tg_auto_helper.TgUser{}).
		Where("id = ?", user.ID).
		Updates(map[string]any{
			"phone":        user.Phone,
			"api_id":       user.ApiId,
			"api_hash":     user.ApiHash,
			"session_path": user.SessionPath,
		}).Error
}

func (s *TgUserService) GetTgUser(ctx context.Context, ID int) (tg_auto_helper.TgUser, error) {
	var user tg_auto_helper.TgUser
	err := global.GVA_MYSQL.Where("id = ?", ID).First(&user).Error
	return user, err
}

func (s *TgUserService) GetTgUserInfoList(
	ctx context.Context,
	info tgAutoHelperReq.TgUserSearch,
) (list []tg_auto_helper.TgUser, total int64, err error) {

	db := global.GVA_MYSQL.Model(&tg_auto_helper.TgUser{})
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	if len(info.CreatedAtRange) == 2 {
		db = db.Where(
			"created_at BETWEEN ? AND ?",
			info.CreatedAtRange[0],
			info.CreatedAtRange[1],
		)
	}

	if err = db.Count(&total).Error; err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&list).Error
	return
}
