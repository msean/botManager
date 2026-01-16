package ledger

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/common/response"
	"github.com/msean/botmanager/server/model/ledger"
	ledgerReq "github.com/msean/botmanager/server/model/ledger/request"
	"go.uber.org/zap"
)

type LedgerPermissionApi struct{}

// CreateLedgerPermission 创建帐薄权限管理
// @Tags LedgerPermission
// @Summary 创建帐薄权限管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body usage.LedgerPermission true "创建帐薄权限管理"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /ledgerPermission/createLedgerPermission [post]
func (ledgerPermissionApi *LedgerPermissionApi) CreateLedgerPermission(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var ledgerPermission ledger.LedgerPermission
	err := c.ShouldBindJSON(&ledgerPermission)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = ledgerPermissionService.CreateLedgerPermission(ctx, &ledgerPermission)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteLedgerPermission 删除帐薄权限管理
// @Tags LedgerPermission
// @Summary 删除帐薄权限管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body usage.LedgerPermission true "删除帐薄权限管理"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /ledgerPermission/deleteLedgerPermission [delete]
func (ledgerPermissionApi *LedgerPermissionApi) DeleteLedgerPermission(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := ledgerPermissionService.DeleteLedgerPermission(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteLedgerPermissionByIds 批量删除帐薄权限管理
// @Tags LedgerPermission
// @Summary 批量删除帐薄权限管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /ledgerPermission/deleteLedgerPermissionByIds [delete]
func (ledgerPermissionApi *LedgerPermissionApi) DeleteLedgerPermissionByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := ledgerPermissionService.DeleteLedgerPermissionByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateLedgerPermission 更新帐薄权限管理
// @Tags LedgerPermission
// @Summary 更新帐薄权限管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body usage.LedgerPermission true "更新帐薄权限管理"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /ledgerPermission/updateLedgerPermission [put]
func (ledgerPermissionApi *LedgerPermissionApi) UpdateLedgerPermission(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var ledgerPermission ledger.LedgerPermission
	err := c.ShouldBindJSON(&ledgerPermission)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = ledgerPermissionService.UpdateLedgerPermission(ctx, ledgerPermission)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindLedgerPermission 用id查询帐薄权限管理
// @Tags LedgerPermission
// @Summary 用id查询帐薄权限管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询帐薄权限管理"
// @Success 200 {object} response.Response{data=usage.LedgerPermission,msg=string} "查询成功"
// @Router /ledgerPermission/findLedgerPermission [get]
func (ledgerPermissionApi *LedgerPermissionApi) FindLedgerPermission(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	reledgerPermission, err := ledgerPermissionService.GetLedgerPermission(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}

	response.OkWithData(reledgerPermission, c)
}

// GetLedgerPermissionList 分页获取帐薄权限管理列表
// @Tags LedgerPermission
// @Summary 分页获取帐薄权限管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query usageReq.LedgerPermissionSearch true "分页获取帐薄权限管理列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /ledgerPermission/getLedgerPermissionList [get]
func (ledgerPermissionApi *LedgerPermissionApi) GetLedgerPermissionList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo ledgerReq.LedgerPermissionSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := ledgerPermissionService.GetLedgerPermissionInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetLedgerPermissionPublic 不需要鉴权的帐薄权限管理接口
// @Tags LedgerPermission
// @Summary 不需要鉴权的帐薄权限管理接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /ledgerPermission/getLedgerPermissionPublic [get]
func (ledgerPermissionApi *LedgerPermissionApi) GetLedgerPermissionPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	ledgerPermissionService.GetLedgerPermissionPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的帐薄权限管理接口信息",
	}, "获取成功", c)
}
