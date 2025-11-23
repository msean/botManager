package recharge

import (
	"github.com/msean/botmanager/server/middleware"
	"github.com/gin-gonic/gin"
)

type UserRechargeRecordRouter struct {}

// InitUserRechargeRecordRouter 初始化 用户充值记录 路由信息
func (s *UserRechargeRecordRouter) InitUserRechargeRecordRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	userRechargeRecordRouter := Router.Group("userRechargeRecord").Use(middleware.OperationRecord())
	userRechargeRecordRouterWithoutRecord := Router.Group("userRechargeRecord")
	userRechargeRecordRouterWithoutAuth := PublicRouter.Group("userRechargeRecord")
	{
		userRechargeRecordRouter.POST("createUserRechargeRecord", userRechargeRecordApi.CreateUserRechargeRecord)   // 新建用户充值记录
		userRechargeRecordRouter.DELETE("deleteUserRechargeRecord", userRechargeRecordApi.DeleteUserRechargeRecord) // 删除用户充值记录
		userRechargeRecordRouter.DELETE("deleteUserRechargeRecordByIds", userRechargeRecordApi.DeleteUserRechargeRecordByIds) // 批量删除用户充值记录
		userRechargeRecordRouter.PUT("updateUserRechargeRecord", userRechargeRecordApi.UpdateUserRechargeRecord)    // 更新用户充值记录
	}
	{
		userRechargeRecordRouterWithoutRecord.GET("findUserRechargeRecord", userRechargeRecordApi.FindUserRechargeRecord)        // 根据ID获取用户充值记录
		userRechargeRecordRouterWithoutRecord.GET("getUserRechargeRecordList", userRechargeRecordApi.GetUserRechargeRecordList)  // 获取用户充值记录列表
	}
	{
	    userRechargeRecordRouterWithoutAuth.GET("getUserRechargeRecordPublic", userRechargeRecordApi.GetUserRechargeRecordPublic)  // 用户充值记录开放接口
	}
}
