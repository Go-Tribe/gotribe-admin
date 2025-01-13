import request from '@/utils/request'

// 获取推广场景列表
export function getSceneList(params) {
  return request({
    url: '/api/ad/scene',
    method: 'get',
    params
  })
}

// 创建推广场景
export function createScene(data) {
  return request({
    url: '/api/ad/scene',
    method: 'post',
    data
  })
}

// 更新推广场景
export function updateScene(data) {
  return request({
    url: '/api/ad/scene/' + data.adSceneID,
    method: 'patch',
    data
  })
}

// 批量删除推广场景
export function batchDeleteScene(data) {
  return request({
    url: '/api/ad/scene',
    method: 'delete',
    data
  })
}

// 获取推广内容列表
export function getAdList(params) {
  return request({
    url: '/api/ad',
    method: 'get',
    params
  })
}

// 创建推广内容
export function createAd(data) {
  return request({
    url: '/api/ad',
    method: 'post',
    data
  })
}

// 更新推广内容
export function updateAd(data) {
  return request({
    url: '/api/ad/' + data.adID,
    method: 'patch',
    data
  })
}

// 批量删除推广内容
export function batchDeleteAd(data) {
  return request({
    url: '/api/ad',
    method: 'delete',
    data
  })
}
