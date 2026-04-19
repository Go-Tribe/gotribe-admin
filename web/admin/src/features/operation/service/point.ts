import { request } from '@/service'
import { type PointListResponse, type PointListParams, type PointCreateParams } from '../types/point'

/**
 * 获取积分列表
 * GET /api/point?pageNum=1&pageSize=3&userID=&nickname=&projectID=
 */
export async function getPointList(
  params?: PointListParams
): Promise<PointListResponse> {
  const requestParams: Record<string, string | number | undefined> = {
    pageNum: params?.pageNum ?? 1,
    pageSize: params?.pageSize ?? 10,
  }

  // 只有当参数有值时才添加到请求参数中
  if (params?.userID != null && params.userID !== '') {
    requestParams.userID = params.userID
  }
  if (params?.nickname != null && params.nickname !== '') {
    requestParams.nickname = params.nickname
  }
  if (params?.projectID != null && params.projectID !== '') {
    requestParams.projectID = params.projectID
  }

  const data = await request.get<PointListResponse>('/api/point', {
    params: requestParams,
  })
  return data as PointListResponse
}

/**
 * 创建积分
 * POST /api/point?project_id
 */
export async function createPoint(
  data: PointCreateParams,
): Promise<unknown> {
  const result = await request.post('/api/point', data)
  return result ?? {}
}
