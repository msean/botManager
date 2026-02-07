package tg_auto_helper

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/common/response"
	"github.com/msean/botmanager/server/model/tg_auto_helper"
	tgAutoHelperReq "github.com/msean/botmanager/server/model/tg_auto_helper/request"
	"go.uber.org/zap"
)

type TgUserApi struct{}

// ================= CRUD =================

// CreateTgUser
func (api *TgUserApi) CreateTgUser(c *gin.Context) {
	var user tg_auto_helper.TgUser
	if err := c.ShouldBindJSON(&user); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := tgUserService.Create(&user); err != nil {
		global.GVA_LOG.Error("创建失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdateTgUser
func (api *TgUserApi) UpdateTgUser(c *gin.Context) {
	var user tg_auto_helper.TgUser
	if err := c.ShouldBindJSON(&user); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := tgUserService.UpdateTgUser(c.Request.Context(), user); err != nil {
		global.GVA_LOG.Error("更新失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteTgUser
func (api *TgUserApi) DeleteTgUser(c *gin.Context) {
	ID := c.Query("ID")
	if err := tgUserService.DeleteTgUser(c.Request.Context(), ID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetTgUserList
func (api *TgUserApi) GetTgUserList(c *gin.Context) {
	var pageInfo tgAutoHelperReq.TgUserSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := tgUserService.GetTgUserInfoList(c.Request.Context(), pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// ================= Telegram 登录 =================

// SendCode
func (api *TgUserApi) SendCode(c *gin.Context) {
	ID := c.PostForm("id")
	user, err := tgUserService.GetTgUser(c.Request.Context(), ID)
	if err != nil {
		response.FailWithMessage("用户不存在", c)
		return
	}

	if err := tgUserService.SendCode(c.Request.Context(), &user); err != nil {
		global.GVA_LOG.Error("发送验证码失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("验证码已发送", c)
}

// VerifyCode
func (api *TgUserApi) VerifyCode(c *gin.Context) {
	ID := c.PostForm("id")
	code := c.PostForm("code")

	user, err := tgUserService.GetTgUser(c.Request.Context(), ID)
	if err != nil {
		response.FailWithMessage("用户不存在", c)
		return
	}

	if err := tgUserService.VerifyCode(c.Request.Context(), &user, code); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("验证码验证成功", c)
}

// VerifyPassword
func (api *TgUserApi) VerifyPassword(c *gin.Context) {
	ID := c.PostForm("id")
	password := c.PostForm("password")

	user, err := tgUserService.GetTgUser(c.Request.Context(), ID)
	if err != nil {
		response.FailWithMessage("用户不存在", c)
		return
	}

	if err := tgUserService.VerifyPassword(c.Request.Context(), &user, password); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("登录完成", c)
}
