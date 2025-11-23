package recharge

import (
	
	"github.com/msean/botmanager/server/global"
    "github.com/msean/botmanager/server/model/common/response"
    "github.com/msean/botmanager/server/model/recharge"
    rechargeReq "github.com/msean/botmanager/server/model/recharge/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type RechargeConfigApi struct {}



// CreateRechargeConfig 创建充值配置
// @Tags RechargeConfig
// @Summary 创建充值配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body recharge.RechargeConfig true "创建充值配置"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /rechargeConfig/createRechargeConfig [post]
func (rechargeConfigApi *RechargeConfigApi) CreateRechargeConfig(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var rechargeConfig recharge.RechargeConfig
	err := c.ShouldBindJSON(&rechargeConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = rechargeConfigService.CreateRechargeConfig(ctx,&rechargeConfig)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteRechargeConfig 删除充值配置
// @Tags RechargeConfig
// @Summary 删除充值配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body recharge.RechargeConfig true "删除充值配置"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /rechargeConfig/deleteRechargeConfig [delete]
func (rechargeConfigApi *RechargeConfigApi) DeleteRechargeConfig(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := rechargeConfigService.DeleteRechargeConfig(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteRechargeConfigByIds 批量删除充值配置
// @Tags RechargeConfig
// @Summary 批量删除充值配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /rechargeConfig/deleteRechargeConfigByIds [delete]
func (rechargeConfigApi *RechargeConfigApi) DeleteRechargeConfigByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := rechargeConfigService.DeleteRechargeConfigByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateRechargeConfig 更新充值配置
// @Tags RechargeConfig
// @Summary 更新充值配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body recharge.RechargeConfig true "更新充值配置"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /rechargeConfig/updateRechargeConfig [put]
func (rechargeConfigApi *RechargeConfigApi) UpdateRechargeConfig(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var rechargeConfig recharge.RechargeConfig
	err := c.ShouldBindJSON(&rechargeConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = rechargeConfigService.UpdateRechargeConfig(ctx,rechargeConfig)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindRechargeConfig 用id查询充值配置
// @Tags RechargeConfig
// @Summary 用id查询充值配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询充值配置"
// @Success 200 {object} response.Response{data=recharge.RechargeConfig,msg=string} "查询成功"
// @Router /rechargeConfig/findRechargeConfig [get]
func (rechargeConfigApi *RechargeConfigApi) FindRechargeConfig(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	rerechargeConfig, err := rechargeConfigService.GetRechargeConfig(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(rerechargeConfig, c)
}
// GetRechargeConfigList 分页获取充值配置列表
// @Tags RechargeConfig
// @Summary 分页获取充值配置列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query rechargeReq.RechargeConfigSearch true "分页获取充值配置列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /rechargeConfig/getRechargeConfigList [get]
func (rechargeConfigApi *RechargeConfigApi) GetRechargeConfigList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo rechargeReq.RechargeConfigSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := rechargeConfigService.GetRechargeConfigInfoList(ctx,pageInfo)
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

// GetRechargeConfigPublic 不需要鉴权的充值配置接口
// @Tags RechargeConfig
// @Summary 不需要鉴权的充值配置接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /rechargeConfig/getRechargeConfigPublic [get]
func (rechargeConfigApi *RechargeConfigApi) GetRechargeConfigPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    rechargeConfigService.GetRechargeConfigPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的充值配置接口信息",
    }, "获取成功", c)
}
