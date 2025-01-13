import request from '@/utils/request'

// 获取接口列表
export function getColumnList(params) {
  return request({
    url: '/api/column',
    method: 'get',
    params
  })
}

// 创建接口
export function createColumn(data) {
  return request({
    url: '/api/column',
    method: 'post',
    data
  })
}

// 更新接口
export function updateColumn(Id, data) {
  return request({
    url: '/api/column/' + Id,
    method: 'patch',
    data
  })
}

// 批量删除接口
export function batchDeleteColumn(data) {
  return request({
    url: '/api/column',
    method: 'delete',
    data
  })
}
