import request from '@/utils/request'

export function getTypeList(params) {
  return request({
    url: '/api/product/type',
    method: 'get',
    params
  })
}

export function createType(data) {
  return request({
    url: '/api/product/type',
    method: 'post',
    data
  })
}

export function updateType(Id, data) {
  return request({
    url: '/api/product/type/' + Id,
    method: 'patch',
    data
  })
}

export function batchDeleteType(data) {
  return request({
    url: '/api/product/type',
    method: 'delete',
    data
  })
}

