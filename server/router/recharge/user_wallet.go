package recharge

import (
	"github.com/msean/botmanager/server/middleware"
	"github.com/gin-gonic/gin"
)

type UserWalletRouter struct {}

// InitUserWalletRouter 初始化 用户钱包 路由信息
func (s *UserWalletRouter) InitUserWalletRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	userWalletRouter := Router.Group("userWallet").Use(middleware.OperationRecord())
	userWalletRouterWithoutRecord := Router.Group("userWallet")
	userWalletRouterWithoutAuth := PublicRouter.Group("userWallet")
	{
		userWalletRouter.POST("createUserWallet", userWalletApi.CreateUserWallet)   // 新建用户钱包
		userWalletRouter.DELETE("deleteUserWallet", userWalletApi.DeleteUserWallet) // 删除用户钱包
		userWalletRouter.DELETE("deleteUserWalletByIds", userWalletApi.DeleteUserWalletByIds) // 批量删除用户钱包
		userWalletRouter.PUT("updateUserWallet", userWalletApi.UpdateUserWallet)    // 更新用户钱包
	}
	{
		userWalletRouterWithoutRecord.GET("findUserWallet", userWalletApi.FindUserWallet)        // 根据ID获取用户钱包
		userWalletRouterWithoutRecord.GET("getUserWalletList", userWalletApi.GetUserWalletList)  // 获取用户钱包列表
	}
	{
	    userWalletRouterWithoutAuth.GET("getUserWalletPublic", userWalletApi.GetUserWalletPublic)  // 用户钱包开放接口
	}
}
