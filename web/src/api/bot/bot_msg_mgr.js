import service from '@/utils/request'
// @Tags BotMsgMgr
// @Summary 创建机器人消息管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotMsgMgr true "创建机器人消息管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /bot_msg_mgr/createBotMsgMgr [post]
export const createBotMsgMgr = (data) => {
  return service({
    url: '/bot_msg_mgr/createBotMsgMgr',
    method: 'post',
    data
  })
}

// @Tags BotMsgMgr
// @Summary 删除机器人消息管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotMsgMgr true "删除机器人消息管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bot_msg_mgr/deleteBotMsgMgr [delete]
export const deleteBotMsgMgr = (params) => {
  return service({
    url: '/bot_msg_mgr/deleteBotMsgMgr',
    method: 'delete',
    params
  })
}

// @Tags BotMsgMgr
// @Summary 批量删除机器人消息管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除机器人消息管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bot_msg_mgr/deleteBotMsgMgr [delete]
export const deleteBotMsgMgrByIds = (params) => {
  return service({
    url: '/bot_msg_mgr/deleteBotMsgMgrByIds',
    method: 'delete',
    params
  })
}

// @Tags BotMsgMgr
// @Summary 更新机器人消息管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotMsgMgr true "更新机器人消息管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /bot_msg_mgr/updateBotMsgMgr [put]
export const updateBotMsgMgr = (data) => {
  return service({
    url: '/bot_msg_mgr/updateBotMsgMgr',
    method: 'put',
    data
  })
}

// @Tags BotMsgMgr
// @Summary 用id查询机器人消息管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BotMsgMgr true "用id查询机器人消息管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /bot_msg_mgr/findBotMsgMgr [get]
export const findBotMsgMgr = (params) => {
  return service({
    url: '/bot_msg_mgr/findBotMsgMgr',
    method: 'get',
    params
  })
}

// @Tags BotMsgMgr
// @Summary 分页获取机器人消息管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取机器人消息管理列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bot_msg_mgr/getBotMsgMgrList [get]
export const getBotMsgMgrList = (params) => {
  return service({
    url: '/bot_msg_mgr/getBotMsgMgrList',
    method: 'get',
    params
  })
}

// @Tags BotMsgMgr
// @Summary 不需要鉴权的机器人消息管理接口
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotMsgMgrSearch true "分页获取机器人消息管理列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /bot_msg_mgr/getBotMsgMgrPublic [get]
export const getBotMsgMgrPublic = () => {
  return service({
    url: '/bot_msg_mgr/getBotMsgMgrPublic',
    method: 'get',
  })
}
