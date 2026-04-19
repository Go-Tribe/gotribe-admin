import { request } from '@/service'
import type {
  SceneCreateParams,
  SceneListParams,
  SceneListResponse,
  SceneUpdateParams,
} from './types'

/**
 * 获取广告场景列表
 * GET /api/ad/scene?pageNum=1&pageSize=10&projectID=
 */
export async function getSceneList(
  params?: SceneListParams
): Promise<SceneListResponse> {
  const requestParams: Record<string, string | number | undefined> = {
    pageNum: params?.pageNum ?? 1,
    pageSize: params?.pageSize ?? 10,
  }
  if (params?.projectID != null && params.projectID !== '') {
    requestParams.projectID = params.projectID
  }
  const data = await request.get<SceneListResponse>('/api/ad/scene', {
    params: requestParams,
  })
  return data as SceneListResponse
}

/**
 * 新建广告场景
 * POST /api/ad/scene?project_id=  body: title, description, projectID
 */
export async function createScene(
  data: SceneCreateParams
): Promise<unknown> {
  const result = await request.post('/api/ad/scene', data)
  return result ?? {}
}

/**
 * 更新广告场景
 * PATCH /api/ad/scene/:adSceneID  body: { title, description }
 */
export async function updateScene(
  adSceneID: string,
  data: SceneUpdateParams
): Promise<void> {
  const id = adSceneID.trim()
  await request.patch(`/api/ad/scene/${id}`, data)
}

/**
 * 删除广告场景
 * DELETE /api/ad/scene  body: { adScenesIds: string }
 */
export async function deleteScene(adSceneIds: string): Promise<void> {
  await request.delete('/api/ad/scene', {
    data: { ids: [Number(adSceneIds.trim())] },
  })
}
