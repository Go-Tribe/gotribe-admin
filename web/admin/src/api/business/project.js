import request from '@/utils/request'

// 获取接口列表
export function getProjectList(params) {
  return request({
    url: '/api/project',
    method: 'get',
    params
  })
}

// 创建接口
export function createProject(data) {
  return request({
    url: '/api/project',
    method: 'post',
    data
  })
}

// 更新接口
export function updateProject(Id, data) {
  return request({
    url: '/api/project/' + Id,
    method: 'patch',
    data
  })
}

// 批量删除接口
export function batchDeleteProject(data) {
  return request({
    url: '/api/project',
    method: 'delete',
    data
  })
}
