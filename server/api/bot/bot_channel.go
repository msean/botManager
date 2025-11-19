package bot

import (
	
	"github.com/msean/botmanager/server/global"
    "github.com/msean/botmanager/server/model/common/response"
    "github.com/msean/botmanager/server/model/bot"
    botReq "github.com/msean/botmanager/server/model/bot/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type BotChannelApi struct {}



// CreateBotChannel 创建机器人渠道
// @Tags BotChannel
// @Summary 创建机器人渠道
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotChannel true "创建机器人渠道"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /botChannel/createBotChannel [post]
func (botChannelApi *BotChannelApi) CreateBotChannel(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var botChannel bot.BotChannel
	err := c.ShouldBindJSON(&botChannel)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = botChannelService.CreateBotChannel(ctx,&botChannel)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteBotChannel 删除机器人渠道
// @Tags BotChannel
// @Summary 删除机器人渠道
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotChannel true "删除机器人渠道"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /botChannel/deleteBotChannel [delete]
func (botChannelApi *BotChannelApi) DeleteBotChannel(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := botChannelService.DeleteBotChannel(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBotChannelByIds 批量删除机器人渠道
// @Tags BotChannel
// @Summary 批量删除机器人渠道
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /botChannel/deleteBotChannelByIds [delete]
func (botChannelApi *BotChannelApi) DeleteBotChannelByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := botChannelService.DeleteBotChannelByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBotChannel 更新机器人渠道
// @Tags BotChannel
// @Summary 更新机器人渠道
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotChannel true "更新机器人渠道"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /botChannel/updateBotChannel [put]
func (botChannelApi *BotChannelApi) UpdateBotChannel(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var botChannel bot.BotChannel
	err := c.ShouldBindJSON(&botChannel)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = botChannelService.UpdateBotChannel(ctx,botChannel)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBotChannel 用id查询机器人渠道
// @Tags BotChannel
// @Summary 用id查询机器人渠道
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询机器人渠道"
// @Success 200 {object} response.Response{data=bot.BotChannel,msg=string} "查询成功"
// @Router /botChannel/findBotChannel [get]
func (botChannelApi *BotChannelApi) FindBotChannel(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	rebotChannel, err := botChannelService.GetBotChannel(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(rebotChannel, c)
}
// GetBotChannelList 分页获取机器人渠道列表
// @Tags BotChannel
// @Summary 分页获取机器人渠道列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotChannelSearch true "分页获取机器人渠道列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /botChannel/getBotChannelList [get]
func (botChannelApi *BotChannelApi) GetBotChannelList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo botReq.BotChannelSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := botChannelService.GetBotChannelInfoList(ctx,pageInfo)
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

// GetBotChannelPublic 不需要鉴权的机器人渠道接口
// @Tags BotChannel
// @Summary 不需要鉴权的机器人渠道接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /botChannel/getBotChannelPublic [get]
func (botChannelApi *BotChannelApi) GetBotChannelPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    botChannelService.GetBotChannelPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的机器人渠道接口信息",
    }, "获取成功", c)
}
