import request from '@/utils/request'

export function getSpecList(params) {
  return request({
    url: '/api/product/spec',
    method: 'get',
    params
  })
}

export function createSpec(data) {
  return request({
    url: '/api/product/spec',
    method: 'post',
    data
  })
}

export function updateSpec(Id, data) {
  return request({
    url: '/api/product/spec/' + Id,
    method: 'patch',
    data
  })
}

export function batchDeleteSpec(data) {
  return request({
    url: '/api/product/spec',
    method: 'delete',
    data
  })
}

export function getSpecItemList(params) {
  return request({
    url: '/api/product/spec/item',
    method: 'get',
    params
  })
}

export function createSpecItem(data) {
  return request({
    url: '/api/product/spec/item',
    method: 'post',
    data
  })
}

export function updateSpecItem(Id, data) {
  return request({
    url: '/api/product/spec/item/' + Id,
    method: 'patch',
    data
  })
}

export function batchDeleteSpecItem(data) {
  return request({
    url: '/api/product/spec/item',
    method: 'delete',
    data
  })
}

