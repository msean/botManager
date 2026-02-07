package bot

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
	"github.com/msean/botmanager/server/model/common/response"
	"go.uber.org/zap"
)

type BotMsgMassApi struct{}

// CreateBotMsgMass 创建机器人群发
// @Tags BotMsgMass
// @Summary 创建机器人群发
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotMsgMass true "创建机器人群发"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /botMsgMass/createBotMsgMass [post]
func (botMsgMassApi *BotMsgMassApi) CreateBotMsgMass(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var botMsgMass bot.BotMsgMass
	err := c.ShouldBindJSON(&botMsgMass)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = botMsgMassService.CreateBotMsgMass(ctx, &botMsgMass)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteBotMsgMass 删除机器人群发
// @Tags BotMsgMass
// @Summary 删除机器人群发
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotMsgMass true "删除机器人群发"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /botMsgMass/deleteBotMsgMass [delete]
func (botMsgMassApi *BotMsgMassApi) DeleteBotMsgMass(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := botMsgMassService.DeleteBotMsgMass(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBotMsgMassByIds 批量删除机器人群发
// @Tags BotMsgMass
// @Summary 批量删除机器人群发
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /botMsgMass/deleteBotMsgMassByIds [delete]
func (botMsgMassApi *BotMsgMassApi) DeleteBotMsgMassByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := botMsgMassService.DeleteBotMsgMassByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBotMsgMass 更新机器人群发
// @Tags BotMsgMass
// @Summary 更新机器人群发
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotMsgMass true "更新机器人群发"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /botMsgMass/updateBotMsgMass [put]
func (botMsgMassApi *BotMsgMassApi) UpdateBotMsgMass(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var botMsgMass bot.BotMsgMass
	err := c.ShouldBindJSON(&botMsgMass)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = botMsgMassService.UpdateBotMsgMass(ctx, botMsgMass)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBotMsgMass 用id查询机器人群发
// @Tags BotMsgMass
// @Summary 用id查询机器人群发
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询机器人群发"
// @Success 200 {object} response.Response{data=bot.BotMsgMass,msg=string} "查询成功"
// @Router /botMsgMass/findBotMsgMass [get]
func (botMsgMassApi *BotMsgMassApi) FindBotMsgMass(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	rebotMsgMass, err := botMsgMassService.GetBotMsgMass(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(rebotMsgMass, c)
}

// GetBotMsgMassList 分页获取机器人群发列表
// @Tags BotMsgMass
// @Summary 分页获取机器人群发列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotMsgMassSearch true "分页获取机器人群发列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /botMsgMass/getBotMsgMassList [get]
func (botMsgMassApi *BotMsgMassApi) GetBotMsgMassList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo botReq.BotMsgMassSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := botMsgMassService.GetBotMsgMassInfoList(ctx, pageInfo)
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

// GetBotMsgMassPublic 不需要鉴权的机器人群发接口
// @Tags BotMsgMass
// @Summary 不需要鉴权的机器人群发接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /botMsgMass/getBotMsgMassPublic [get]
func (botMsgMassApi *BotMsgMassApi) GetBotMsgMassPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	botMsgMassService.GetBotMsgMassPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的机器人群发接口信息",
	}, "获取成功", c)
}

func (botMsgMassApi *BotMsgMassApi) SendBotMsgMass(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var req botReq.BotMsgMassSend
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = botMsgMassService.SendBotMsgMass(ctx, req)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
