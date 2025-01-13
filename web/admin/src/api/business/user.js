import request from '@/utils/request'

// 获取接口列表
export function getUserList(params) {
  return request({
    url: '/api/user',
    method: 'get',
    params
  })
}

// 创建接口
export function createUser(data) {
  return request({
    url: '/api/user',
    method: 'post',
    data
  })
}

// 更新接口
export function updateUser(Id, data) {
  return request({
    url: '/api/user/' + Id,
    method: 'patch',
    data
  })
}
