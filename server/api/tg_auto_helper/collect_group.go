package tg_auto_helper

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/global"
	"github.com/msean/botmanager/server/model/common/response"
	"github.com/msean/botmanager/server/model/tg_auto_helper"
	tgAutoHelperReq "github.com/msean/botmanager/server/model/tg_auto_helper/request"
	"go.uber.org/zap"
)

type CollectGroupApi struct{}

// CreateTgUser
func (api *CollectGroupApi) Create(c *gin.Context) {
	var object tg_auto_helper.CollectGroupTask
	if err := c.ShouldBindJSON(&object); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := collectGroupSvc.Create(&object); err != nil {
		global.GVA_LOG.Error("创建失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (api *CollectGroupApi) Delete(c *gin.Context) {
	ID := c.Query("ID")
	if err := collectGroupSvc.Delete(ID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (api *CollectGroupApi) List(c *gin.Context) {
	var pageInfo tgAutoHelperReq.CollectGroupTaskSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := collectGroupSvc.List(pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (api *CollectGroupApi) ListCollectGroupInfo(c *gin.Context) {

	var search tgAutoHelperReq.CollectGroupInfoSearch

	_ = c.ShouldBindQuery(&search)

	// 默认值
	if search.Page == 0 {
		search.Page = 1
	}
	if search.PageSize == 0 {
		search.PageSize = 10
	}

	list, total, err := collectGroupSvc.ListCollectGroupInfo(search)

	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}

	response.OkWithData(gin.H{
		"list":  list,
		"total": total,
	}, c)
}
