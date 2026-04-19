import { request } from '@/service'
import type { CommentListParams, CommentListResponse } from '../types/comment'

/**
 * 获取评论列表
 * GET /api/comment?pageNum=1&pageSize=10&status=&nickname=&projectID=
 */
export async function getCommentList(
  params?: CommentListParams
): Promise<CommentListResponse> {
  const requestParams: Record<string, string | number | undefined> = {
    pageNum: params?.pageNum ?? 1,
    pageSize: params?.pageSize ?? 10,
  }

  if (params?.status != null && params.status !== '') {
    requestParams.status = params.status
  }
  if (params?.nickname != null && params.nickname !== '') {
    requestParams.nickname = params.nickname
  }
  if (params?.projectID != null && params.projectID !== '') {
    requestParams.projectID = params.projectID
  }
  if (params?.sortBy) {
    requestParams.sortBy = params.sortBy
  }
  if (params?.sortOrder) {
    requestParams.sortOrder = params.sortOrder
  }

  const data = await request.get<CommentListResponse>('/api/comment', {
    params: requestParams,
  })
  return data as CommentListResponse
}

/**
 * 审核评论（通过）
 * PATCH /api/comment/:commentID  body: {}
 */
export async function approveComment(commentID: string): Promise<void> {
  const id = commentID.trim()
  await request.patch(`/api/comment/${id}`, {})
}

/**
 * 评论设为不通过
 * PATCH /api/comment/:commentID  body: { status: 1 }
 */
export async function rejectComment(commentID: string): Promise<void> {
  const id = commentID.trim()
  await request.patch(`/api/comment/${id}`, { status: 1 })
}
