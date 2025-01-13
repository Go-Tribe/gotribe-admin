import request from '@/utils/request'

// 获取接口列表
export function getCommentList(params) {
  return request({
    url: '/api/comment',
    method: 'get',
    params
  })
}

// 更新接口
export function updateComment(data) {
  return request({
    url: '/api/comment/' + data.commentID,
    method: 'patch',
    data
  })
}
