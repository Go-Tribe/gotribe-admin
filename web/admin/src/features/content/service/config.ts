import { request } from '@/service'
import type {
  Config,
  ConfigListParams,
  ConfigListResponse,
  ConfigCreateParams,
  ConfigUpdateParams,
} from '../types/config'

/**
 * 获取配置列表
 * GET /api/config?configID=&title=&pageNum=1&pageSize=3&projectID=&type=
 */
export async function getConfigList(
  params?: ConfigListParams
): Promise<ConfigListResponse> {
  const data = await request.get<ConfigListResponse>('/api/config', {
    params: {
      configID: params?.configID ?? '',
      title: params?.title ?? '',
      pageNum: params?.pageNum ?? 1,
      pageSize: params?.pageSize ?? 10,
      projectID: params?.projectID,
      type: params?.type,
    },
  })
  return data as ConfigListResponse
}

/**
 * 新增配置
 * POST /api/config?project_id=  body: title, description, info, projectID, alias, type, mdContent
 */
export async function createConfig(
  data: ConfigCreateParams
): Promise<{ configID?: string }> {
  const result = await request.post<{ configID?: string }>('/api/config', data)
  return (result ?? {}) as { configID?: string }
}

/**
 * 获取配置详情
 * GET /api/config/:configID
 */
export async function getConfig(configID: string): Promise<Config> {
  const id = configID.trim()
  const data = await request.get<{ config: Config }>(`/api/config/${id}`)
  return (data as { config: Config }).config
}

/**
 * 更新配置
 * PATCH /api/config/:configID  body: title, description, projectID, info, mdContent
 */
export async function updateConfig(
  configID: string,
  data: ConfigUpdateParams
): Promise<void> {
  const id = configID.trim()
  await request.patch(`/api/config/${id}`, data)
}

/**
 * 删除配置
 * DELETE /api/config  body: { configIds: string }
 */
export async function deleteConfig(configIds: string): Promise<void> {
  await request.delete('/api/config', {
    data: { configIds: configIds.trim() },
  })
}
