package ledger

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/common/response"
	"github.com/msean/botmanager/server/model/ledger"
	ledgerReq "github.com/msean/botmanager/server/model/ledger/request"
	"go.uber.org/zap"
)

type LedgerApi struct{}

// CreateLedger 创建帐薄
// @Tags Ledger
// @Summary 创建帐薄
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body usage.Ledger true "创建帐薄"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /ledger/createLedger [post]
func (ledgerApi *LedgerApi) CreateLedger(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var ledger ledger.Ledger
	err := c.ShouldBindJSON(&ledger)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = ledgerService.CreateLedger(ctx, &ledger)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteLedger 删除帐薄
// @Tags Ledger
// @Summary 删除帐薄
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body usage.Ledger true "删除帐薄"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /ledger/deleteLedger [delete]
func (ledgerApi *LedgerApi) DeleteLedger(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := ledgerService.DeleteLedger(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteLedgerByIds 批量删除帐薄
// @Tags Ledger
// @Summary 批量删除帐薄
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /ledger/deleteLedgerByIds [delete]
func (ledgerApi *LedgerApi) DeleteLedgerByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := ledgerService.DeleteLedgerByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateLedger 更新帐薄
// @Tags Ledger
// @Summary 更新帐薄
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body usage.Ledger true "更新帐薄"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /ledger/updateLedger [put]
func (ledgerApi *LedgerApi) UpdateLedger(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var ledger ledger.Ledger
	err := c.ShouldBindJSON(&ledger)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = ledgerService.UpdateLedger(ctx, ledger)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindLedger 用id查询帐薄
// @Tags Ledger
// @Summary 用id查询帐薄
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询帐薄"
// @Success 200 {object} response.Response{data=usage.Ledger,msg=string} "查询成功"
// @Router /ledger/findLedger [get]
func (ledgerApi *LedgerApi) FindLedger(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	reledger, err := ledgerService.GetLedger(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(reledger, c)
}

// GetLedgerList 分页获取帐薄列表
// @Tags Ledger
// @Summary 分页获取帐薄列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query usageReq.LedgerSearch true "分页获取帐薄列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /ledger/getLedgerList [get]
func (ledgerApi *LedgerApi) GetLedgerList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo ledgerReq.LedgerSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := ledgerService.GetLedgerInfoList(ctx, pageInfo)
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

// GetLedgerPublic 不需要鉴权的帐薄接口
// @Tags Ledger
// @Summary 不需要鉴权的帐薄接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /ledger/getLedgerPublic [get]
func (ledgerApi *LedgerApi) GetLedgerPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	ledgerService.GetLedgerPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的帐薄接口信息",
	}, "获取成功", c)
}

func (ledgerApi *LedgerApi) Full(c *gin.Context) {
	var req struct {
		BotID       int64 `form:"bot_id" binding:"required"`
		ChatGroupID int64 `form:"chat_group_id" binding:"required"`
		IDMin       int64 `form:"idmin" binding:"required"`
		IDMax       int64 `form:"idmax" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
		return
	}

	var list []ledger.Ledger
	err := global.GVA_MYSQL.
		Where("bot_id = ?", req.BotID).
		Where("chat_group_id = ?", req.ChatGroupID).
		Where("id BETWEEN ? AND ?", req.IDMin, req.IDMax).
		Order("id asc").
		Find(&list).Error
	if err != nil {
		c.JSON(500, gin.H{"msg": "查询失败"})
		return
	}

	var incomeList, payoutList []gin.H
	var totalIncome, totalPayout float64

	for _, v := range list {
		row := gin.H{
			"time":      v.CreatedAt.Format("15:04:05"),
			"amount":    v.Amount,
			"remark":    "",
			"replyUser": "",
			"operator":  v.OprUserNickname,
			"afterNote": "",
		}

		if v.ActionType == 1 {
			incomeList = append(incomeList, row)
			totalIncome += v.Amount
		} else if v.ActionType == 2 {
			payoutList = append(payoutList, row)
			totalPayout += v.Amount
		}
	}

	c.JSON(200, gin.H{
		"income": gin.H{
			"list":  incomeList,
			"count": len(incomeList),
		},
		"payout": gin.H{
			"list":  payoutList,
			"count": len(payoutList),
		},
		"summary": gin.H{
			"totalIncome": totalIncome,
			"totalPayout": totalPayout,
			"unpaid":      totalIncome - totalPayout,
		},
	})
}
