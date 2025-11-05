package bot

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
	"github.com/msean/botmanager/server/model/common/response"
	"go.uber.org/zap"
)

type BotBanGroupMemApi struct{}

// CreateBotBanGroupMem 创建封禁成员设置
// @Tags BotBanGroupMem
// @Summary 创建封禁成员设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotBanGroupMem true "创建封禁成员设置"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /botBanGroupMem/createBotBanGroupMem [post]
func (botBanGroupMemApi *BotBanGroupMemApi) CreateBotBanGroupMem(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var botBanGroupMem bot.BotBanGroupMem
	err := c.ShouldBindJSON(&botBanGroupMem)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = botBanGroupMemService.CreateBotBanGroupMem(ctx, &botBanGroupMem)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteBotBanGroupMem 删除封禁成员设置
// @Tags BotBanGroupMem
// @Summary 删除封禁成员设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotBanGroupMem true "删除封禁成员设置"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /botBanGroupMem/deleteBotBanGroupMem [delete]
func (botBanGroupMemApi *BotBanGroupMemApi) DeleteBotBanGroupMem(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := botBanGroupMemService.DeleteBotBanGroupMem(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBotBanGroupMemByIds 批量删除封禁成员设置
// @Tags BotBanGroupMem
// @Summary 批量删除封禁成员设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /botBanGroupMem/deleteBotBanGroupMemByIds [delete]
func (botBanGroupMemApi *BotBanGroupMemApi) DeleteBotBanGroupMemByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := botBanGroupMemService.DeleteBotBanGroupMemByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBotBanGroupMem 更新封禁成员设置
// @Tags BotBanGroupMem
// @Summary 更新封禁成员设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotBanGroupMem true "更新封禁成员设置"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /botBanGroupMem/updateBotBanGroupMem [put]
func (botBanGroupMemApi *BotBanGroupMemApi) UpdateBotBanGroupMem(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var botBanGroupMem bot.BotBanGroupMem
	err := c.ShouldBindJSON(&botBanGroupMem)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = botBanGroupMemService.UpdateBotBanGroupMem(ctx, botBanGroupMem)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBotBanGroupMem 用id查询封禁成员设置
// @Tags BotBanGroupMem
// @Summary 用id查询封禁成员设置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询封禁成员设置"
// @Success 200 {object} response.Response{data=bot.BotBanGroupMem,msg=string} "查询成功"
// @Router /botBanGroupMem/findBotBanGroupMem [get]
func (botBanGroupMemApi *BotBanGroupMemApi) FindBotBanGroupMem(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	rebotBanGroupMem, err := botBanGroupMemService.GetBotBanGroupMem(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(rebotBanGroupMem, c)
}

// GetBotBanGroupMemList 分页获取封禁成员设置列表
// @Tags BotBanGroupMem
// @Summary 分页获取封禁成员设置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotBanGroupMemSearch true "分页获取封禁成员设置列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /botBanGroupMem/getBotBanGroupMemList [get]
func (botBanGroupMemApi *BotBanGroupMemApi) GetBotBanGroupMemList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo botReq.BotBanGroupMemSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := botBanGroupMemService.GetBotBanGroupMemInfoList(ctx, pageInfo)
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

// GetBotBanGroupMemPublic 不需要鉴权的封禁成员设置接口
// @Tags BotBanGroupMem
// @Summary 不需要鉴权的封禁成员设置接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /botBanGroupMem/getBotBanGroupMemPublic [get]
func (botBanGroupMemApi *BotBanGroupMemApi) GetBotBanGroupMemPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	botBanGroupMemService.GetBotBanGroupMemPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的封禁成员设置接口信息",
	}, "获取成功", c)
}
