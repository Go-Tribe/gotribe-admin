import request from '@/utils/request'

// 获取接口列表
export function getTagList(params) {
  return request({
    url: '/api/tag',
    method: 'get',
    params
  })
}

// 创建接口
export function createTag(data) {
  return request({
    url: '/api/tag',
    method: 'post',
    data
  })
}

// 更新接口
export function updateTag(Id, data) {
  return request({
    url: '/api/tag/' + Id,
    method: 'patch',
    data
  })
}

// 批量删除接口
export function batchDeleteTag(data) {
  return request({
    url: '/api/tag',
    method: 'delete',
    data
  })
}
