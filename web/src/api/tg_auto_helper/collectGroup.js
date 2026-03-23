import service from '@/utils/request'

export const createCollectGroupTask = (data) => {
  return service({
    url: '/collect_group/create',
    method: 'post',
    data
  })
}


export const deleteCollectGroupTask = (params) => {
  return service({
    url: '/collect_group/delete',
    method: 'delete',
    params
  })
}


export const listCollectGroupTask = (params) => {
  return service({
    url: '/collect_group/list',
    method: 'get',
    params
  })
}


export const listCollectGroupInfo = (params) => {
  return service({
    url: '/collect_group/list_collect_group_info',
    method: 'get',
    params
  })
}