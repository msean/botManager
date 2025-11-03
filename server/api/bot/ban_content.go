package bot

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
	"github.com/msean/botmanager/server/model/common/response"
	"go.uber.org/zap"
)

type BotBanContentApi struct{}

// CreateBotBanContent 创建机器人消息管理
// @Tags BotBanContent
// @Summary 创建机器人消息管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotBanContent true "创建机器人消息管理"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /bot_msg_mgr/createBotBanContent [post]
func (api *BotBanContentApi) CreateBotBanContent(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var botBanContent bot.BotBanContent
	err := c.ShouldBindJSON(&botBanContent)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = BotBanContentService.CreateBotBanContent(ctx, &botBanContent)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteBotBanContent 删除机器人消息管理
// @Tags BotBanContent
// @Summary 删除机器人消息管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotBanContent true "删除机器人消息管理"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /bot_msg_mgr/deleteBotBanContent [delete]
func (api *BotBanContentApi) DeleteBotBanContent(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := BotBanContentService.DeleteBotBanContent(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBotBanContentByIds 批量删除机器人消息管理
// @Tags BotBanContent
// @Summary 批量删除机器人消息管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /bot_msg_mgr/deleteBotBanContentByIds [delete]
func (api *BotBanContentApi) DeleteBotBanContentByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := BotBanContentService.DeleteBotBanContentByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBotBanContent 更新机器人消息管理
// @Tags BotBanContent
// @Summary 更新机器人消息管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotBanContent true "更新机器人消息管理"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /bot_msg_mgr/updateBotBanContent [put]
func (api *BotBanContentApi) UpdateBotBanContent(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var bot_msg_mgr bot.BotBanContent
	err := c.ShouldBindJSON(&bot_msg_mgr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = BotBanContentService.UpdateBotBanContent(ctx, bot_msg_mgr)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBotBanContent 用id查询机器人消息管理
// @Tags BotBanContent
// @Summary 用id查询机器人消息管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询机器人消息管理"
// @Success 200 {object} response.Response{data=bot.BotBanContent,msg=string} "查询成功"
// @Router /bot_msg_mgr/findBotBanContent [get]
func (api *BotBanContentApi) FindBotBanContent(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	rebot_msg_mgr, err := BotBanContentService.GetBotBanContent(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(rebot_msg_mgr, c)
}

// GetBotBanContentList 分页获取机器人消息管理列表
// @Tags BotBanContent
// @Summary 分页获取机器人消息管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotBanContentSearch true "分页获取机器人消息管理列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /bot_msg_mgr/getBotBanContentList [get]
func (api *BotBanContentApi) GetBotBanContentList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo botReq.BotBanContentSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := BotBanContentService.GetBotBanContentInfoList(ctx, pageInfo)
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

// GetBotBanContentPublic 不需要鉴权的机器人消息管理接口
// @Tags BotBanContent
// @Summary 不需要鉴权的机器人消息管理接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /bot_msg_mgr/getBotBanContentPublic [get]
func (api *BotBanContentApi) GetBotBanContentPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	BotBanContentService.GetBotBanContentPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的机器人消息管理接口信息",
	}, "获取成功", c)
}
