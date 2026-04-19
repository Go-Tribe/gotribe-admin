import { request } from '@/service'
import type {
  AdListParams,
  AdListResponse,
  AdCreateParams,
  AdUpdateParams,
} from './types'

/**
 * 获取广告列表
 * GET /api/ad?pageNum=1&pageSize=3&sceneID=&title&status
 */
export async function getAdList(
  params?: AdListParams
): Promise<AdListResponse> {
  const requestParams: Record<string, string | number | undefined> = {
    pageNum: params?.pageNum ?? 1,
    pageSize: params?.pageSize ?? 10,
  }

  // 只有当参数有值时才添加到请求参数中
  if (params?.sceneID != null && params.sceneID !== '') {
    requestParams.sceneID = params.sceneID
  }
  if (params?.title != null && params.title !== '') {
    requestParams.title = params.title
  }
  if (params?.status != null && params.status !== '') {
    requestParams.status = params.status
  }
  if (params?.sortBy) {
    requestParams.sortBy = params.sortBy
  }
  if (params?.sortOrder) {
    requestParams.sortOrder = params.sortOrder
  }

  const data = await request.get<AdListResponse>('/api/ad', {
    params: requestParams,
  })
  return data as AdListResponse
}

/**
 * 创建广告
 * POST /api/ad?project_id
 */
export async function createAd(
  data: AdCreateParams
): Promise<unknown> {
  const requestData: Record<string, unknown> = {
    title: data.title,
    description: data.description,
    sceneID: data.sceneID,
    url: data.url,
    urlType: data.urlType,
    image: data.image,
    sort: data.sort,
    status: data.status,
    ext: data.ext,
  }
  // 只有当 video 有值时才添加
  if (data.video) {
    requestData.video = data.video
  }
  const result = await request.post('/api/ad', requestData)
  return result ?? {}
}

/**
 * 更新广告
 * PATCH /api/ad/:adID
 */
export async function updateAd(
  adID: string,
  data: AdUpdateParams
): Promise<void> {
  const id = adID.trim()
  await request.patch(`/api/ad/${id}`, data)
}

/**
 * 删除广告
 * DELETE /api/ad  body: { adsIds: string }
 */
export async function deleteAd(adIds: string): Promise<void> {
  await request.delete('/api/ad', {
    data: { ids: [Number(adIds.trim())] },
  })
}
