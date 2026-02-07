package tg_auto_helper

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/middleware"
)

type TgUserRouter struct{}

// export const sendCode = (data) => post('/tgUser/sendCode', data)
// export const verifyCode = (data) => post('/tgUser/verifyCode', data)
// export const verifyPassword = (data) => post('/tgUser/verifyPassword', data)

// InitTgUserRouter 初始化 telegram用户管理 路由信息
func (s *TgUserRouter) InitTgUserRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	tgUserRouter := Router.Group("tgUser").Use(middleware.OperationRecord())
	tgUserRouterWithoutRecord := Router.Group("tgUser")
	{
		tgUserRouter.POST("createTgUser", tgUserApi.CreateTgUser)
		tgUserRouter.DELETE("deleteTgUser", tgUserApi.DeleteTgUser)
		tgUserRouter.PUT("updateTgUser", tgUserApi.UpdateTgUser)
		tgUserRouter.POST("sendCode", tgUserApi.SendCode)
		tgUserRouter.POST("verifyCode", tgUserApi.VerifyCode)
		tgUserRouter.POST("verifyPassword", tgUserApi.VerifyPassword)
	}
	{
		tgUserRouterWithoutRecord.GET("getTgUserList", tgUserApi.GetTgUserList) // 获取telegram用户管理列表
	}
}
