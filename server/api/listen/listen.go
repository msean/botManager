package listen

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/api/listen"
	"github.com/msean/botmanager/server/model/common/response"
)

type ListenApi struct{}

func (api *ListenApi) Choice(c *gin.Context) {
	list, err := listenService.Choice(c.Request.Context())
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
	return
}

func (api *ListenApi) Query(c *gin.Context) {
	var req listen.ListenQueryReq

	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 默认分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	list, total, err := listenService.Query(
		c.Request.Context(),
		req,
	)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(gin.H{
		"list":  list,
		"total": total,
	}, c)
}

func (api *ListenApi) Export(c *gin.Context) {
	var req listen.ListenQueryReq

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	req.IsExport = true

	file, err := listenService.Export(
		c.Request.Context(),
		req,
	)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(gin.H{
		"file": file,
	}, c)
}
