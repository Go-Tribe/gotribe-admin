import request from '@/utils/request'

export function getProductList(params) {
  return request({
    url: '/api/product',
    method: 'get',
    params
  })
}

export function createProduct(data) {
  return request({
    url: '/api/product',
    method: 'post',
    data
  })
}

export function updateProduct(data) {
  return request({
    url: '/api/product/' + data.productID,
    method: 'patch',
    data
  })
}

export function batchDeleteProduct(data) {
  return request({
    url: '/api/product',
    method: 'delete',
    data
  })
}

export function getProductDetail(Id) {
  return request({
    url: '/api/product/' + Id,
    method: 'get'
  })
}

export function getSpecDetail(Id) {
  return request({
    url: '/api/product/spec/info/' + Id,
    method: 'get'
  })
}

