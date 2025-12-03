package recharge

import (
	
	"github.com/msean/botmanager/server/global"
    "github.com/msean/botmanager/server/model/common/response"
    "github.com/msean/botmanager/server/model/recharge"
    rechargeReq "github.com/msean/botmanager/server/model/recharge/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type AdPublishRecordApi struct {}



// CreateAdPublishRecord 创建广告发布记录
// @Tags AdPublishRecord
// @Summary 创建广告发布记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body recharge.AdPublishRecord true "创建广告发布记录"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /adPublishRecord/createAdPublishRecord [post]
func (adPublishRecordApi *AdPublishRecordApi) CreateAdPublishRecord(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var adPublishRecord recharge.AdPublishRecord
	err := c.ShouldBindJSON(&adPublishRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = adPublishRecordService.CreateAdPublishRecord(ctx,&adPublishRecord)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteAdPublishRecord 删除广告发布记录
// @Tags AdPublishRecord
// @Summary 删除广告发布记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body recharge.AdPublishRecord true "删除广告发布记录"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /adPublishRecord/deleteAdPublishRecord [delete]
func (adPublishRecordApi *AdPublishRecordApi) DeleteAdPublishRecord(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := adPublishRecordService.DeleteAdPublishRecord(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteAdPublishRecordByIds 批量删除广告发布记录
// @Tags AdPublishRecord
// @Summary 批量删除广告发布记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /adPublishRecord/deleteAdPublishRecordByIds [delete]
func (adPublishRecordApi *AdPublishRecordApi) DeleteAdPublishRecordByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := adPublishRecordService.DeleteAdPublishRecordByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateAdPublishRecord 更新广告发布记录
// @Tags AdPublishRecord
// @Summary 更新广告发布记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body recharge.AdPublishRecord true "更新广告发布记录"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /adPublishRecord/updateAdPublishRecord [put]
func (adPublishRecordApi *AdPublishRecordApi) UpdateAdPublishRecord(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var adPublishRecord recharge.AdPublishRecord
	err := c.ShouldBindJSON(&adPublishRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = adPublishRecordService.UpdateAdPublishRecord(ctx,adPublishRecord)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindAdPublishRecord 用id查询广告发布记录
// @Tags AdPublishRecord
// @Summary 用id查询广告发布记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询广告发布记录"
// @Success 200 {object} response.Response{data=recharge.AdPublishRecord,msg=string} "查询成功"
// @Router /adPublishRecord/findAdPublishRecord [get]
func (adPublishRecordApi *AdPublishRecordApi) FindAdPublishRecord(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	readPublishRecord, err := adPublishRecordService.GetAdPublishRecord(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(readPublishRecord, c)
}
// GetAdPublishRecordList 分页获取广告发布记录列表
// @Tags AdPublishRecord
// @Summary 分页获取广告发布记录列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query rechargeReq.AdPublishRecordSearch true "分页获取广告发布记录列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /adPublishRecord/getAdPublishRecordList [get]
func (adPublishRecordApi *AdPublishRecordApi) GetAdPublishRecordList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo rechargeReq.AdPublishRecordSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := adPublishRecordService.GetAdPublishRecordInfoList(ctx,pageInfo)
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

// GetAdPublishRecordPublic 不需要鉴权的广告发布记录接口
// @Tags AdPublishRecord
// @Summary 不需要鉴权的广告发布记录接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /adPublishRecord/getAdPublishRecordPublic [get]
func (adPublishRecordApi *AdPublishRecordApi) GetAdPublishRecordPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    adPublishRecordService.GetAdPublishRecordPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的广告发布记录接口信息",
    }, "获取成功", c)
}
