package bot

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/bot"
	botReq "github.com/msean/botmanager/server/model/bot/request"
	"github.com/msean/botmanager/server/model/common/response"
	"go.uber.org/zap"
)

type BotCmdConfigApi struct{}

// CreateBotCmdConfig 创建机器人命令配置
// @Tags BotCmdConfig
// @Summary 创建机器人命令配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotCmdConfig true "创建机器人命令配置"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /botCmdConfig/createBotCmdConfig [post]
func (botCmdConfigApi *BotCmdConfigApi) CreateBotCmdConfig(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var botCmdConfig bot.BotCmdConfig
	err := c.ShouldBindJSON(&botCmdConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = botCmdConfigService.CreateBotCmdConfig(ctx, &botCmdConfig)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteBotCmdConfig 删除机器人命令配置
// @Tags BotCmdConfig
// @Summary 删除机器人命令配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotCmdConfig true "删除机器人命令配置"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /botCmdConfig/deleteBotCmdConfig [delete]
func (botCmdConfigApi *BotCmdConfigApi) DeleteBotCmdConfig(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := botCmdConfigService.DeleteBotCmdConfig(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}

	// cache.NewBotCmdCacheList()
	response.OkWithMessage("删除成功", c)
}

// DeleteBotCmdConfigByIds 批量删除机器人命令配置
// @Tags BotCmdConfig
// @Summary 批量删除机器人命令配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /botCmdConfig/deleteBotCmdConfigByIds [delete]
func (botCmdConfigApi *BotCmdConfigApi) DeleteBotCmdConfigByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := botCmdConfigService.DeleteBotCmdConfigByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBotCmdConfig 更新机器人命令配置
// @Tags BotCmdConfig
// @Summary 更新机器人命令配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotCmdConfig true "更新机器人命令配置"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /botCmdConfig/updateBotCmdConfig [put]
func (botCmdConfigApi *BotCmdConfigApi) UpdateBotCmdConfig(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var botCmdConfig bot.BotCmdConfig
	err := c.ShouldBindJSON(&botCmdConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = botCmdConfigService.UpdateBotCmdConfig(ctx, botCmdConfig)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBotCmdConfig 用id查询机器人命令配置
// @Tags BotCmdConfig
// @Summary 用id查询机器人命令配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询机器人命令配置"
// @Success 200 {object} response.Response{data=bot.BotCmdConfig,msg=string} "查询成功"
// @Router /botCmdConfig/findBotCmdConfig [get]
func (botCmdConfigApi *BotCmdConfigApi) FindBotCmdConfig(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	rebotCmdConfig, err := botCmdConfigService.GetBotCmdConfig(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(rebotCmdConfig, c)
}

// GetBotCmdConfigList 分页获取机器人命令配置列表
// @Tags BotCmdConfig
// @Summary 分页获取机器人命令配置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotCmdConfigSearch true "分页获取机器人命令配置列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /botCmdConfig/getBotCmdConfigList [get]
func (botCmdConfigApi *BotCmdConfigApi) GetBotCmdConfigList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo botReq.BotCmdConfigSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := botCmdConfigService.GetBotCmdConfigInfoList(ctx, pageInfo)
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

// GetBotCmdConfigPublic 不需要鉴权的机器人命令配置接口
// @Tags BotCmdConfig
// @Summary 不需要鉴权的机器人命令配置接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /botCmdConfig/getBotCmdConfigPublic [get]
func (botCmdConfigApi *BotCmdConfigApi) GetBotCmdConfigPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	botCmdConfigService.GetBotCmdConfigPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的机器人命令配置接口信息",
	}, "获取成功", c)
}
