package bot

import (
	
	"github.com/msean/botmanager/server/global"
    "github.com/msean/botmanager/server/model/common/response"
    "github.com/msean/botmanager/server/model/bot"
    botReq "github.com/msean/botmanager/server/model/bot/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type BotApi struct {}



// CreateBot 创建机器人
// @Tags Bot
// @Summary 创建机器人
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.Bot true "创建机器人"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /bot_mgr/createBot [post]
func (bot_mgrApi *BotApi) CreateBot(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var bot_mgr bot.Bot
	err := c.ShouldBindJSON(&bot_mgr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = bot_mgrService.CreateBot(ctx,&bot_mgr)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteBot 删除机器人
// @Tags Bot
// @Summary 删除机器人
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.Bot true "删除机器人"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /bot_mgr/deleteBot [delete]
func (bot_mgrApi *BotApi) DeleteBot(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := bot_mgrService.DeleteBot(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBotByIds 批量删除机器人
// @Tags Bot
// @Summary 批量删除机器人
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /bot_mgr/deleteBotByIds [delete]
func (bot_mgrApi *BotApi) DeleteBotByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := bot_mgrService.DeleteBotByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBot 更新机器人
// @Tags Bot
// @Summary 更新机器人
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.Bot true "更新机器人"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /bot_mgr/updateBot [put]
func (bot_mgrApi *BotApi) UpdateBot(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var bot_mgr bot.Bot
	err := c.ShouldBindJSON(&bot_mgr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = bot_mgrService.UpdateBot(ctx,bot_mgr)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBot 用id查询机器人
// @Tags Bot
// @Summary 用id查询机器人
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询机器人"
// @Success 200 {object} response.Response{data=bot.Bot,msg=string} "查询成功"
// @Router /bot_mgr/findBot [get]
func (bot_mgrApi *BotApi) FindBot(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	rebot_mgr, err := bot_mgrService.GetBot(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(rebot_mgr, c)
}
// GetBotList 分页获取机器人列表
// @Tags Bot
// @Summary 分页获取机器人列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotSearch true "分页获取机器人列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /bot_mgr/getBotList [get]
func (bot_mgrApi *BotApi) GetBotList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo botReq.BotSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := bot_mgrService.GetBotInfoList(ctx,pageInfo)
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

// GetBotPublic 不需要鉴权的机器人接口
// @Tags Bot
// @Summary 不需要鉴权的机器人接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /bot_mgr/getBotPublic [get]
func (bot_mgrApi *BotApi) GetBotPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    bot_mgrService.GetBotPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的机器人接口信息",
    }, "获取成功", c)
}
