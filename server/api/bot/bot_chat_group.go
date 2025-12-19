package bot

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	"github.com/msean/botmanager/server/model/bot/request"
	botReq "github.com/msean/botmanager/server/model/bot/request"
	"github.com/msean/botmanager/server/model/common/response"
	botService "github.com/msean/botmanager/server/service/bot"
	"go.uber.org/zap"
)

type BotChatGroupApi struct{}

// CreateBotChatGroup 创建机器人群组列表
// @Tags BotChatGroup
// @Summary 创建机器人群组列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotChatGroup true "创建机器人群组列表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /botChatGroup/createBotChatGroup [post]
func (botChatGroupApi *BotChatGroupApi) CreateBotChatGroup(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var botChatGroup bot.BotChatGroup
	err := c.ShouldBindJSON(&botChatGroup)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = botChatGroupService.CreateBotChatGroup(ctx, &botChatGroup)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteBotChatGroup 删除机器人群组列表
// @Tags BotChatGroup
// @Summary 删除机器人群组列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotChatGroup true "删除机器人群组列表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /botChatGroup/deleteBotChatGroup [delete]
func (botChatGroupApi *BotChatGroupApi) DeleteBotChatGroup(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := botChatGroupService.DeleteBotChatGroup(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBotChatGroupByIds 批量删除机器人群组列表
// @Tags BotChatGroup
// @Summary 批量删除机器人群组列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /botChatGroup/deleteBotChatGroupByIds [delete]
func (botChatGroupApi *BotChatGroupApi) DeleteBotChatGroupByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := botChatGroupService.DeleteBotChatGroupByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBotChatGroup 更新机器人群组列表
// @Tags BotChatGroup
// @Summary 更新机器人群组列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotChatGroup true "更新机器人群组列表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /botChatGroup/updateBotChatGroup [put]
func (botChatGroupApi *BotChatGroupApi) UpdateBotChatGroup(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var botChatGroup bot.BotChatGroup
	err := c.ShouldBindJSON(&botChatGroup)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = botChatGroupService.UpdateBotChatGroup(ctx, botChatGroup)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBotChatGroup 用id查询机器人群组列表
// @Tags BotChatGroup
// @Summary 用id查询机器人群组列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询机器人群组列表"
// @Success 200 {object} response.Response{data=bot.BotChatGroup,msg=string} "查询成功"
// @Router /botChatGroup/findBotChatGroup [get]
func (botChatGroupApi *BotChatGroupApi) FindBotChatGroup(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	rebotChatGroup, err := botChatGroupService.GetBotChatGroup(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(rebotChatGroup, c)
}

// GetBotChatGroupList 分页获取机器人群组列表列表
// @Tags BotChatGroup
// @Summary 分页获取机器人群组列表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotChatGroupSearch true "分页获取机器人群组列表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /botChatGroup/getBotChatGroupList [get]
func (botChatGroupApi *BotChatGroupApi) GetBotChatGroupList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo botReq.BotChatGroupSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := botChatGroupService.GetBotChatGroupInfoList(ctx, pageInfo)
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

// GetBotChatGroupPublic 不需要鉴权的机器人群组列表接口
// @Tags BotChatGroup
// @Summary 不需要鉴权的机器人群组列表接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /botChatGroup/getBotChatGroupPublic [get]
func (botChatGroupApi *BotChatGroupApi) GetBotChatGroupPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	botChatGroupService.GetBotChatGroupPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的机器人群组列表接口信息",
	}, "获取成功", c)
}

func (api *BotChatGroupApi) ChatHistory(c *gin.Context) {
	var req request.ChatMessageQuery
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	layout := "2006-01-02 15:04:05"
	if req.SrcStartTime != "" {
		t, err := time.ParseInLocation(layout, req.SrcStartTime, time.Local)
		if err != nil {
			response.FailWithMessage("startTime 格式错误", c)
			return
		}
		req.StartTime = &t
	}

	if req.SrcEndTime != "" {
		t, err := time.ParseInLocation(layout, req.SrcEndTime, time.Local)
		if err != nil {
			response.FailWithMessage("endTime 格式错误", c)
			return
		}
		req.EndTime = &t
	}

	if req.Limit <= 0 || req.Limit > 200 {
		req.Limit = 20
	}

	svc := botService.NewBotChatHistorySvc(req.BotID, req.ChatGroupID)
	list, hasMore, err := svc.QueryMessages(req)
	if err != nil {
		global.GVA_LOG.Error("ChatHistory QueryMessages", zap.Error(err))
	}

	response.OkWithDetailed(gin.H{
		"list":    list,
		"hasMore": hasMore,
	}, "ok", c)
}
