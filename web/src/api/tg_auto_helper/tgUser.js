import service from '@/utils/request'
// @Tags TgUser
// @Summary 创建telegram用户管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.TgUser true "创建telegram用户管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /tgUser/createTgUser [post]
export const createTgUser = (data) => {
  return service({
    url: '/tgUser/createTgUser',
    method: 'post',
    data
  })
}

// @Tags TgUser
// @Summary 删除telegram用户管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.TgUser true "删除telegram用户管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /tgUser/deleteTgUser [delete]
export const deleteTgUser = (params) => {
  return service({
    url: '/tgUser/deleteTgUser',
    method: 'delete',
    params
  })
}

// @Tags TgUser
// @Summary 批量删除telegram用户管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除telegram用户管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /tgUser/deleteTgUser [delete]
export const deleteTgUserByIds = (params) => {
  return service({
    url: '/tgUser/deleteTgUserByIds',
    method: 'delete',
    params
  })
}

// @Tags TgUser
// @Summary 更新telegram用户管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.TgUser true "更新telegram用户管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /tgUser/updateTgUser [put]
export const updateTgUser = (data) => {
  return service({
    url: '/tgUser/updateTgUser',
    method: 'put',
    data
  })
}

// @Tags TgUser
// @Summary 用id查询telegram用户管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.TgUser true "用id查询telegram用户管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /tgUser/findTgUser [get]
export const findTgUser = (params) => {
  return service({
    url: '/tgUser/findTgUser',
    method: 'get',
    params
  })
}

// @Tags TgUser
// @Summary 分页获取telegram用户管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取telegram用户管理列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /tgUser/getTgUserList [get]
export const getTgUserList = (params) => {
  return service({
    url: '/tgUser/getTgUserList',
    method: 'get',
    params
  })
}

// @Tags TgUser
// @Summary 不需要鉴权的telegram用户管理接口
// @Accept application/json
// @Produce application/json
// @Param data query tg_auto_helperReq.TgUserSearch true "分页获取telegram用户管理列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /tgUser/getTgUserPublic [get]
export const getTgUserPublic = () => {
  return service({
    url: '/tgUser/getTgUserPublic',
    method: 'get',
  })
}


export const sendCode = (data) => {
  return service({
    url: '/tgUser/sendCode',
    method: 'post',
    data
  })
}


export const verifyCode = (data) => {
  return service({
    url: '/tgUser/verifyCode',
    method: 'post',
    data
  })
}

export const verifyPassword = (data) => {
  return service({
    url: '/tgUser/verifyPassword',
    method: 'post',
    data
  })
}
// export const sendCode = (data) => post('/tgUser/sendCode', data)
// export const verifyCode = (data) => post('/tgUser/verifyCode', data)
// export const verifyPassword = (data) => post('/tgUser/verifyPassword', data)
