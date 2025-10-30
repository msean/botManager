import service from '@/utils/request'
// @Tags Bot
// @Summary 创建机器人
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Bot true "创建机器人"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /bot_mgr/createBot [post]
export const createBot = (data) => {
  return service({
    url: '/bot_mgr/createBot',
    method: 'post',
    data
  })
}

// @Tags Bot
// @Summary 删除机器人
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Bot true "删除机器人"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bot_mgr/deleteBot [delete]
export const deleteBot = (params) => {
  return service({
    url: '/bot_mgr/deleteBot',
    method: 'delete',
    params
  })
}

// @Tags Bot
// @Summary 批量删除机器人
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除机器人"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bot_mgr/deleteBot [delete]
export const deleteBotByIds = (params) => {
  return service({
    url: '/bot_mgr/deleteBotByIds',
    method: 'delete',
    params
  })
}

// @Tags Bot
// @Summary 更新机器人
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Bot true "更新机器人"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /bot_mgr/updateBot [put]
export const updateBot = (data) => {
  return service({
    url: '/bot_mgr/updateBot',
    method: 'put',
    data
  })
}

// @Tags Bot
// @Summary 用id查询机器人
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Bot true "用id查询机器人"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /bot_mgr/findBot [get]
export const findBot = (params) => {
  return service({
    url: '/bot_mgr/findBot',
    method: 'get',
    params
  })
}

// @Tags Bot
// @Summary 分页获取机器人列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取机器人列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bot_mgr/getBotList [get]
export const getBotList = (params) => {
  return service({
    url: '/bot_mgr/getBotList',
    method: 'get',
    params
  })
}

// @Tags Bot
// @Summary 不需要鉴权的机器人接口
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotSearch true "分页获取机器人列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /bot_mgr/getBotPublic [get]
export const getBotPublic = () => {
  return service({
    url: '/bot_mgr/getBotPublic',
    method: 'get',
  })
}
