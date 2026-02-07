import service from '@/utils/request'
// @Tags BotMassMsgRecord
// @Summary 创建群发历史记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotMassMsgRecord true "创建群发历史记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /botMassMsgRecord/createBotMassMsgRecord [post]
export const createBotMassMsgRecord = (data) => {
  return service({
    url: '/botMassMsgRecord/createBotMassMsgRecord',
    method: 'post',
    data
  })
}

// @Tags BotMassMsgRecord
// @Summary 删除群发历史记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotMassMsgRecord true "删除群发历史记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /botMassMsgRecord/deleteBotMassMsgRecord [delete]
export const deleteBotMassMsgRecord = (params) => {
  return service({
    url: '/botMassMsgRecord/deleteBotMassMsgRecord',
    method: 'delete',
    params
  })
}

// @Tags BotMassMsgRecord
// @Summary 批量删除群发历史记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除群发历史记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /botMassMsgRecord/deleteBotMassMsgRecord [delete]
export const deleteBotMassMsgRecordByIds = (params) => {
  return service({
    url: '/botMassMsgRecord/deleteBotMassMsgRecordByIds',
    method: 'delete',
    params
  })
}

// @Tags BotMassMsgRecord
// @Summary 更新群发历史记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.BotMassMsgRecord true "更新群发历史记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /botMassMsgRecord/updateBotMassMsgRecord [put]
export const updateBotMassMsgRecord = (data) => {
  return service({
    url: '/botMassMsgRecord/updateBotMassMsgRecord',
    method: 'put',
    data
  })
}

// @Tags BotMassMsgRecord
// @Summary 用id查询群发历史记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.BotMassMsgRecord true "用id查询群发历史记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /botMassMsgRecord/findBotMassMsgRecord [get]
export const findBotMassMsgRecord = (params) => {
  return service({
    url: '/botMassMsgRecord/findBotMassMsgRecord',
    method: 'get',
    params
  })
}

// @Tags BotMassMsgRecord
// @Summary 分页获取群发历史记录列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取群发历史记录列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /botMassMsgRecord/getBotMassMsgRecordList [get]
export const getBotMassMsgRecordList = (params) => {
  return service({
    url: '/botMassMsgRecord/getBotMassMsgRecordList',
    method: 'get',
    params
  })
}

// @Tags BotMassMsgRecord
// @Summary 不需要鉴权的群发历史记录接口
// @Accept application/json
// @Produce application/json
// @Param data query botReq.BotMassMsgRecordSearch true "分页获取群发历史记录列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /botMassMsgRecord/getBotMassMsgRecordPublic [get]
export const getBotMassMsgRecordPublic = () => {
  return service({
    url: '/botMassMsgRecord/getBotMassMsgRecordPublic',
    method: 'get',
  })
}
