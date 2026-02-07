package bot

import (
	
	"github.com/msean/botmanager/server/global"
    "github.com/msean/botmanager/server/model/common/response"
    "github.com/msean/botmanager/server/model/bot"
    botReq "github.com/msean/botmanager/server/model/bot/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type BotMassMsgRecordApi struct {}



// CreateBotMassMsgRecord 创建群发历史记录
// @Tags BotMassMsgRecord
// @Summary 创建群发历史记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotMassMsgRecord true "创建群发历史记录"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /botMassMsgRecord/createBotMassMsgRecord [post]
func (botMassMsgRecordApi *BotMassMsgRecordApi) CreateBotMassMsgRecord(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var botMassMsgRecord bot.BotMassMsgRecord
	err := c.ShouldBindJSON(&botMassMsgRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = botMassMsgRecordService.CreateBotMassMsgRecord(ctx,&botMassMsgRecord)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteBotMassMsgRecord 删除群发历史记录
// @Tags BotMassMsgRecord
// @Summary 删除群发历史记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotMassMsgRecord true "删除群发历史记录"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /botMassMsgRecord/deleteBotMassMsgRecord [delete]
func (botMassMsgRecordApi *BotMassMsgRecordApi) DeleteBotMassMsgRecord(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := botMassMsgRecordService.DeleteBotMassMsgRecord(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteBotMassMsgRecordByIds 批量删除群发历史记录
// @Tags BotMassMsgRecord
// @Summary 批量删除群发历史记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /botMassMsgRecord/deleteBotMassMsgRecordByIds [delete]
func (botMassMsgRecordApi *BotMassMsgRecordApi) DeleteBotMassMsgRecordByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := botMassMsgRecordService.DeleteBotMassMsgRecordByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateBotMassMsgRecord 更新群发历史记录
// @Tags BotMassMsgRecord
// @Summary 更新群发历史记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body bot.BotMassMsgRecord true "更新群发历史记录"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /botMassMsgRecord/updateBotMassMsgRecord [put]
func (botMassMsgRecordApi *BotMassMsgRecordApi) UpdateBotMassMsgRecord(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var botMassMsgRecord bot.BotMassMsgRecord
	err := c.ShouldBindJSON(&botMassMsgRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = botMassMsgRecordService.UpdateBotMassMsgRecord(ctx,botMassMsgRecord)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindBotMassMsgRecord 用id查询群发历史记录
// @Tags BotMassMsgRecord
// @Summary 用id查询群发历史记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询群发历史记录"
// @Success 200 {object} response.Response{data=bot.BotMassMsgRecord,msg=string} "查询成功"
// @Router /botMassMsgRecord/findBotMassMsgRecord [get]
func (botMassMsgRecordApi *BotMassMsgRecordApi) FindBotMassMsgRecord(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	rebotMassMsgRecord, err := botMassMsgRecordService.GetBotMassMsgRecord(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(rebotMassMsgRecord, c)
}
// GetBotMassMsgRecordList 分页获取群发历史记录列表
// @Tags BotMassMsgRecord
// @Summary 分页获取群发历史记录列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotMassMsgRecordSearch true "分页获取群发历史记录列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /botMassMsgRecord/getBotMassMsgRecordList [get]
func (botMassMsgRecordApi *BotMassMsgRecordApi) GetBotMassMsgRecordList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo botReq.BotMassMsgRecordSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := botMassMsgRecordService.GetBotMassMsgRecordInfoList(ctx,pageInfo)
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

// GetBotMassMsgRecordPublic 不需要鉴权的群发历史记录接口
// @Tags BotMassMsgRecord
// @Summary 不需要鉴权的群发历史记录接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /botMassMsgRecord/getBotMassMsgRecordPublic [get]
func (botMassMsgRecordApi *BotMassMsgRecordApi) GetBotMassMsgRecordPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    botMassMsgRecordService.GetBotMassMsgRecordPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的群发历史记录接口信息",
    }, "获取成功", c)
}
