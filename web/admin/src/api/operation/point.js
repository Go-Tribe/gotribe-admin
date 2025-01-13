import request from '@/utils/request'

// 列表接口
export function getPointList(params) {
  return request({
    url: '/api/point',
    method: 'get',
    params
  })
}

// 创建接口
export function createPoint(data) {
  return request({
    url: '/api/point',
    method: 'post',
    data
  })
}
