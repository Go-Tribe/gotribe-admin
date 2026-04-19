import { request } from '@/service'
import type { SystemConfig, ConfigResponse } from '../types/config'

// 获取系统配置（拦截器已解包，直接返回 data）
export const getConfig = async (): Promise<ConfigResponse> => {
  return request.get<ConfigResponse>('/api/base/config')
}

// 更新系统配置
export async function updateConfig(
  params: Partial<SystemConfig>,
): Promise<{ success: boolean }> {
  return request.patch<{ success: boolean }>('/api/system', params)
}
