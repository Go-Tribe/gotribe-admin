import request from '@/utils/request'

// 获取接口列表
export function getArticleList(params) {
  return request({
    url: '/api/post',
    method: 'get',
    params
  })
}

// 创建接口
export function createArticle(data) {
  return request({
    url: '/api/post',
    method: 'post',
    data
  })
}

// 更新接口
export function updateArticle(data) {
  return request({
    url: '/api/post/' + data.postID,
    method: 'patch',
    data
  })
}

// 获取详情接口
export function getArticleDetail(Id) {
  return request({
    url: '/api/post/' + Id,
    method: 'get'
  })
}

// 批量删除接口
export function batchDeleteArticle(data) {
  return request({
    url: '/api/post',
    method: 'delete',
    data
  })
}

// 发布接口
export function pushArticle(data) {
  return request({
    url: '/api/post/' + data.postID,
    method: 'put'
  })
}
