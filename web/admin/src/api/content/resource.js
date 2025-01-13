import request from '@/utils/request'

// 获取接口列表
export function getResourceList(params) {
  return request({
    url: '/api/resource',
    method: 'get',
    params
  })
}

// 上传接口
export function uploadResource(data) {
  return request({
    url: '/api/resource/upload',
    method: 'post',
    data
  })
}

// 更新接口
export function updateResource(Id, data) {
  return request({
    url: '/api/resource/' + Id,
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

// 删除接口
export function deleteResource(data) {
  return request({
    url: '/api/resource',
    method: 'delete',
    data
  })
}
