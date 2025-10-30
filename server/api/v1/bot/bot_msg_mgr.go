package bot

import (
	
	"github.com/msean/botmanager/server/global"
    "github.com/msean/botmanager/server/model/common/response"
    "github.com/msean/botmanager/server/model/bot"
    botReq "github.com/msean/botmanager/server/model/bot/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type BotMsgMgrApi struct {}



// CreateBotMsgMgr 创建机器人消息管理
// @Tags BotMsgMgr
// @Summary 创建机器人消息管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotMsgMgr true "创建机器人消息管理"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /bot_msg_mgr/createBotMsgMgr [post]
func (bot_msg_mgrApi *BotMsgMgrApi) CreateBotMsgMgr(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var bot_msg_mgr bot.BotMsgMgr
	err := c.ShouldBindJSON(&bot_msg_mgr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = bot_msg_mgrService.CreateBotMsgMgr(ctx,&bot_msg_mgr)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteBotMsgMgr 删除机器人消息管理
// @Tags BotMsgMgr
// @Summary 删除机器人消息管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotMsgMgr true "删除机器人消息管理"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /bot_msg_mgr/deleteBotMsgMgr [delete]
func (bot_msg_mgrApi *BotMsgMgrApi) DeleteBotMsgMgr(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := bot_msg_mgrService.DeleteBotMsgMgr(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBotMsgMgrByIds 批量删除机器人消息管理
// @Tags BotMsgMgr
// @Summary 批量删除机器人消息管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /bot_msg_mgr/deleteBotMsgMgrByIds [delete]
func (bot_msg_mgrApi *BotMsgMgrApi) DeleteBotMsgMgrByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := bot_msg_mgrService.DeleteBotMsgMgrByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBotMsgMgr 更新机器人消息管理
// @Tags BotMsgMgr
// @Summary 更新机器人消息管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotMsgMgr true "更新机器人消息管理"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /bot_msg_mgr/updateBotMsgMgr [put]
func (bot_msg_mgrApi *BotMsgMgrApi) UpdateBotMsgMgr(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var bot_msg_mgr bot.BotMsgMgr
	err := c.ShouldBindJSON(&bot_msg_mgr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = bot_msg_mgrService.UpdateBotMsgMgr(ctx,bot_msg_mgr)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBotMsgMgr 用id查询机器人消息管理
// @Tags BotMsgMgr
// @Summary 用id查询机器人消息管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询机器人消息管理"
// @Success 200 {object} response.Response{data=bot.BotMsgMgr,msg=string} "查询成功"
// @Router /bot_msg_mgr/findBotMsgMgr [get]
func (bot_msg_mgrApi *BotMsgMgrApi) FindBotMsgMgr(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	rebot_msg_mgr, err := bot_msg_mgrService.GetBotMsgMgr(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(rebot_msg_mgr, c)
}
// GetBotMsgMgrList 分页获取机器人消息管理列表
// @Tags BotMsgMgr
// @Summary 分页获取机器人消息管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotMsgMgrSearch true "分页获取机器人消息管理列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /bot_msg_mgr/getBotMsgMgrList [get]
func (bot_msg_mgrApi *BotMsgMgrApi) GetBotMsgMgrList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo botReq.BotMsgMgrSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := bot_msg_mgrService.GetBotMsgMgrInfoList(ctx,pageInfo)
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

// GetBotMsgMgrPublic 不需要鉴权的机器人消息管理接口
// @Tags BotMsgMgr
// @Summary 不需要鉴权的机器人消息管理接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /bot_msg_mgr/getBotMsgMgrPublic [get]
func (bot_msg_mgrApi *BotMsgMgrApi) GetBotMsgMgrPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    bot_msg_mgrService.GetBotMsgMgrPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的机器人消息管理接口信息",
    }, "获取成功", c)
}
