package ledger

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/common/response"
	"github.com/msean/botmanager/server/model/ledger"
	ledgerReq "github.com/msean/botmanager/server/model/ledger/request"
	"go.uber.org/zap"
)

type LedgerAccountGroupApi struct{}

// CreateLedgerAccountGroup 创建记账账号组
// @Tags LedgerAccountGroup
// @Summary 创建记账账号组
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body ledger.LedgerAccountGroup true "创建记账账号组"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /ledgerAccountGroup/createLedgerAccountGroup [post]
func (ledgerAccountGroupApi *LedgerAccountGroupApi) CreateLedgerAccountGroup(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var ledgerAccountGroup ledger.LedgerAccountGroup
	err := c.ShouldBindJSON(&ledgerAccountGroup)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = ledgerAccountGroupService.CreateLedgerAccountGroup(ctx, &ledgerAccountGroup)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteLedgerAccountGroup 删除记账账号组
// @Tags LedgerAccountGroup
// @Summary 删除记账账号组
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body ledger.LedgerAccountGroup true "删除记账账号组"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /ledgerAccountGroup/deleteLedgerAccountGroup [delete]
func (ledgerAccountGroupApi *LedgerAccountGroupApi) DeleteLedgerAccountGroup(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := ledgerAccountGroupService.DeleteLedgerAccountGroup(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteLedgerAccountGroupByIds 批量删除记账账号组
// @Tags LedgerAccountGroup
// @Summary 批量删除记账账号组
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /ledgerAccountGroup/deleteLedgerAccountGroupByIds [delete]
func (ledgerAccountGroupApi *LedgerAccountGroupApi) DeleteLedgerAccountGroupByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := ledgerAccountGroupService.DeleteLedgerAccountGroupByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateLedgerAccountGroup 更新记账账号组
// @Tags LedgerAccountGroup
// @Summary 更新记账账号组
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body ledger.LedgerAccountGroup true "更新记账账号组"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /ledgerAccountGroup/updateLedgerAccountGroup [put]
func (ledgerAccountGroupApi *LedgerAccountGroupApi) UpdateLedgerAccountGroup(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var ledgerAccountGroup ledger.LedgerAccountGroup
	err := c.ShouldBindJSON(&ledgerAccountGroup)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = ledgerAccountGroupService.UpdateLedgerAccountGroup(ctx, ledgerAccountGroup)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindLedgerAccountGroup 用id查询记账账号组
// @Tags LedgerAccountGroup
// @Summary 用id查询记账账号组
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询记账账号组"
// @Success 200 {object} response.Response{data=ledger.LedgerAccountGroup,msg=string} "查询成功"
// @Router /ledgerAccountGroup/findLedgerAccountGroup [get]
func (ledgerAccountGroupApi *LedgerAccountGroupApi) FindLedgerAccountGroup(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	reledgerAccountGroup, err := ledgerAccountGroupService.GetLedgerAccountGroup(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(reledgerAccountGroup, c)
}

// GetLedgerAccountGroupList 分页获取记账账号组列表
// @Tags LedgerAccountGroup
// @Summary 分页获取记账账号组列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query ledgerReq.LedgerAccountGroupSearch true "分页获取记账账号组列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /ledgerAccountGroup/getLedgerAccountGroupList [get]
func (ledgerAccountGroupApi *LedgerAccountGroupApi) GetLedgerAccountGroupList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo ledgerReq.LedgerAccountGroupSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := ledgerAccountGroupService.GetLedgerAccountGroupInfoList(ctx, pageInfo)
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

// GetLedgerAccountGroupPublic 不需要鉴权的记账账号组接口
// @Tags LedgerAccountGroup
// @Summary 不需要鉴权的记账账号组接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /ledgerAccountGroup/getLedgerAccountGroupPublic [get]
func (ledgerAccountGroupApi *LedgerAccountGroupApi) GetLedgerAccountGroupPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	ledgerAccountGroupService.GetLedgerAccountGroupPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的记账账号组接口信息",
	}, "获取成功", c)
}
