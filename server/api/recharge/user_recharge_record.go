package recharge

import (
	
	"github.com/msean/botmanager/server/global"
    "github.com/msean/botmanager/server/model/common/response"
    "github.com/msean/botmanager/server/model/recharge"
    rechargeReq "github.com/msean/botmanager/server/model/recharge/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type UserRechargeRecordApi struct {}



// CreateUserRechargeRecord 创建用户充值记录
// @Tags UserRechargeRecord
// @Summary 创建用户充值记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body recharge.UserRechargeRecord true "创建用户充值记录"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /userRechargeRecord/createUserRechargeRecord [post]
func (userRechargeRecordApi *UserRechargeRecordApi) CreateUserRechargeRecord(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var userRechargeRecord recharge.UserRechargeRecord
	err := c.ShouldBindJSON(&userRechargeRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userRechargeRecordService.CreateUserRechargeRecord(ctx,&userRechargeRecord)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteUserRechargeRecord 删除用户充值记录
// @Tags UserRechargeRecord
// @Summary 删除用户充值记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body recharge.UserRechargeRecord true "删除用户充值记录"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /userRechargeRecord/deleteUserRechargeRecord [delete]
func (userRechargeRecordApi *UserRechargeRecordApi) DeleteUserRechargeRecord(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := userRechargeRecordService.DeleteUserRechargeRecord(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteUserRechargeRecordByIds 批量删除用户充值记录
// @Tags UserRechargeRecord
// @Summary 批量删除用户充值记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /userRechargeRecord/deleteUserRechargeRecordByIds [delete]
func (userRechargeRecordApi *UserRechargeRecordApi) DeleteUserRechargeRecordByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := userRechargeRecordService.DeleteUserRechargeRecordByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateUserRechargeRecord 更新用户充值记录
// @Tags UserRechargeRecord
// @Summary 更新用户充值记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body recharge.UserRechargeRecord true "更新用户充值记录"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /userRechargeRecord/updateUserRechargeRecord [put]
func (userRechargeRecordApi *UserRechargeRecordApi) UpdateUserRechargeRecord(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var userRechargeRecord recharge.UserRechargeRecord
	err := c.ShouldBindJSON(&userRechargeRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userRechargeRecordService.UpdateUserRechargeRecord(ctx,userRechargeRecord)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindUserRechargeRecord 用id查询用户充值记录
// @Tags UserRechargeRecord
// @Summary 用id查询用户充值记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询用户充值记录"
// @Success 200 {object} response.Response{data=recharge.UserRechargeRecord,msg=string} "查询成功"
// @Router /userRechargeRecord/findUserRechargeRecord [get]
func (userRechargeRecordApi *UserRechargeRecordApi) FindUserRechargeRecord(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	reuserRechargeRecord, err := userRechargeRecordService.GetUserRechargeRecord(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(reuserRechargeRecord, c)
}
// GetUserRechargeRecordList 分页获取用户充值记录列表
// @Tags UserRechargeRecord
// @Summary 分页获取用户充值记录列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query rechargeReq.UserRechargeRecordSearch true "分页获取用户充值记录列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /userRechargeRecord/getUserRechargeRecordList [get]
func (userRechargeRecordApi *UserRechargeRecordApi) GetUserRechargeRecordList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo rechargeReq.UserRechargeRecordSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := userRechargeRecordService.GetUserRechargeRecordInfoList(ctx,pageInfo)
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

// GetUserRechargeRecordPublic 不需要鉴权的用户充值记录接口
// @Tags UserRechargeRecord
// @Summary 不需要鉴权的用户充值记录接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /userRechargeRecord/getUserRechargeRecordPublic [get]
func (userRechargeRecordApi *UserRechargeRecordApi) GetUserRechargeRecordPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    userRechargeRecordService.GetUserRechargeRecordPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的用户充值记录接口信息",
    }, "获取成功", c)
}
