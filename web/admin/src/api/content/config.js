import request from '@/utils/request'

// 获取接口列表
export function getConfigList(params) {
  return request({
    url: '/api/config',
    method: 'get',
    params
  })
}

// 创建接口
export function createConfig(data) {
  return request({
    url: '/api/config',
    method: 'post',
    data
  })
}

// 更新接口
export function updateConfig(data) {
  return request({
    url: '/api/config/' + data.configID,
    method: 'patch',
    data
  })
}

// 获取详情接口
export function getConfigDetail(Id) {
  return request({
    url: '/api/config/' + Id,
    method: 'get'
  })
}

// 批量删除接口
export function batchDeleteConfig(data) {
  return request({
    url: '/api/config',
    method: 'delete',
    data
  })
}
