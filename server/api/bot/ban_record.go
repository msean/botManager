package bot

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
	"github.com/msean/botmanager/server/model/common/response"
	"go.uber.org/zap"
)

type BanRecordApi struct{}

// CreateBanRecord 创建封禁记录
// @Tags BanRecord
// @Summary 创建封禁记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BanRecord true "创建封禁记录"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /banRecord/createBanRecord [post]
func (banRecordApi *BanRecordApi) CreateBanRecord(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var banRecord bot.BanRecord
	err := c.ShouldBindJSON(&banRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = banRecordService.CreateBanRecord(ctx, &banRecord)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteBanRecord 删除封禁记录
// @Tags BanRecord
// @Summary 删除封禁记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BanRecord true "删除封禁记录"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /banRecord/deleteBanRecord [delete]
func (banRecordApi *BanRecordApi) DeleteBanRecord(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := banRecordService.DeleteBanRecord(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBanRecordByIds 批量删除封禁记录
// @Tags BanRecord
// @Summary 批量删除封禁记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /banRecord/deleteBanRecordByIds [delete]
func (banRecordApi *BanRecordApi) DeleteBanRecordByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := banRecordService.DeleteBanRecordByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBanRecord 更新封禁记录
// @Tags BanRecord
// @Summary 更新封禁记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BanRecord true "更新封禁记录"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /banRecord/updateBanRecord [put]
func (banRecordApi *BanRecordApi) UpdateBanRecord(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var banRecord bot.BanRecord
	err := c.ShouldBindJSON(&banRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = banRecordService.UpdateBanRecord(ctx, banRecord)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBanRecord 用id查询封禁记录
// @Tags BanRecord
// @Summary 用id查询封禁记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询封禁记录"
// @Success 200 {object} response.Response{data=bot.BanRecord,msg=string} "查询成功"
// @Router /banRecord/findBanRecord [get]
func (banRecordApi *BanRecordApi) FindBanRecord(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	rebanRecord, err := banRecordService.GetBanRecord(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(rebanRecord, c)
}

// GetBanRecordList 分页获取封禁记录列表
// @Tags BanRecord
// @Summary 分页获取封禁记录列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BanRecordSearch true "分页获取封禁记录列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /banRecord/getBanRecordList [get]
func (banRecordApi *BanRecordApi) GetBanRecordList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo botReq.BanRecordSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := banRecordService.GetBanRecordInfoList(ctx, pageInfo)
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

// GetBanRecordPublic 不需要鉴权的封禁记录接口
// @Tags BanRecord
// @Summary 不需要鉴权的封禁记录接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /banRecord/getBanRecordPublic [get]
func (banRecordApi *BanRecordApi) GetBanRecordPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	banRecordService.GetBanRecordPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的封禁记录接口信息",
	}, "获取成功", c)
}
