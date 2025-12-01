package recharge

import (
	
	"github.com/msean/botmanager/server/global"
    "github.com/msean/botmanager/server/model/common/response"
    "github.com/msean/botmanager/server/model/recharge"
    rechargeReq "github.com/msean/botmanager/server/model/recharge/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type UserWalletApi struct {}



// CreateUserWallet 创建用户钱包
// @Tags UserWallet
// @Summary 创建用户钱包
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body recharge.UserWallet true "创建用户钱包"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /userWallet/createUserWallet [post]
func (userWalletApi *UserWalletApi) CreateUserWallet(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var userWallet recharge.UserWallet
	err := c.ShouldBindJSON(&userWallet)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userWalletService.CreateUserWallet(ctx,&userWallet)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteUserWallet 删除用户钱包
// @Tags UserWallet
// @Summary 删除用户钱包
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body recharge.UserWallet true "删除用户钱包"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /userWallet/deleteUserWallet [delete]
func (userWalletApi *UserWalletApi) DeleteUserWallet(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := userWalletService.DeleteUserWallet(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteUserWalletByIds 批量删除用户钱包
// @Tags UserWallet
// @Summary 批量删除用户钱包
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /userWallet/deleteUserWalletByIds [delete]
func (userWalletApi *UserWalletApi) DeleteUserWalletByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := userWalletService.DeleteUserWalletByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateUserWallet 更新用户钱包
// @Tags UserWallet
// @Summary 更新用户钱包
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body recharge.UserWallet true "更新用户钱包"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /userWallet/updateUserWallet [put]
func (userWalletApi *UserWalletApi) UpdateUserWallet(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var userWallet recharge.UserWallet
	err := c.ShouldBindJSON(&userWallet)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userWalletService.UpdateUserWallet(ctx,userWallet)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindUserWallet 用id查询用户钱包
// @Tags UserWallet
// @Summary 用id查询用户钱包
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询用户钱包"
// @Success 200 {object} response.Response{data=recharge.UserWallet,msg=string} "查询成功"
// @Router /userWallet/findUserWallet [get]
func (userWalletApi *UserWalletApi) FindUserWallet(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	reuserWallet, err := userWalletService.GetUserWallet(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(reuserWallet, c)
}
// GetUserWalletList 分页获取用户钱包列表
// @Tags UserWallet
// @Summary 分页获取用户钱包列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query rechargeReq.UserWalletSearch true "分页获取用户钱包列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /userWallet/getUserWalletList [get]
func (userWalletApi *UserWalletApi) GetUserWalletList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo rechargeReq.UserWalletSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := userWalletService.GetUserWalletInfoList(ctx,pageInfo)
	if err != nil {
	    global.GVA_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败:" + err.Error(), c)
        return
    }
    response.OkWithDetailed(response.PageResult{
        List:     list,
        Total:    total,
        Page:     pageInfo.Page,
        PageSize: pageInfo.PageSize,
    }, "获取成功", c)
}

// GetUserWalletPublic 不需要鉴权的用户钱包接口
// @Tags UserWallet
// @Summary 不需要鉴权的用户钱包接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /userWallet/getUserWalletPublic [get]
func (userWalletApi *UserWalletApi) GetUserWalletPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    userWalletService.GetUserWalletPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的用户钱包接口信息",
    }, "获取成功", c)
}
