import { request } from '@/service'
import type { Tag, TagListResponse } from '../types/tag'

// 获取标签列表（支持分页和筛选）
export const getTagList = async (params?: {
  current?: number
  title?: string
  description?: string
  pageNum?: number
  pageSize?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}) => {
  const query = {
    current: params?.current,
    title: params?.title,
    description: params?.description,
    pageNum: params?.pageNum,
    pageSize: params?.pageSize,
    sortBy: params?.sortBy,
    sortOrder: params?.sortOrder,
  }
  return request.get<TagListResponse>('/api/tag', { params: query })
}

// 创建标签
export async function createTag(
  params: Partial<Tag>,
): Promise<{ success: boolean }> {
  return request.post<{ success: boolean }>('/api/tag', params)
}

// 更新标签
export async function updateTag(
  id: number,
  params: Partial<Tag>,
): Promise<{ success: boolean }> {
  return request.patch<{ success: boolean }>(`/api/tag/${id}`, params)
}

// 删除标签
export async function deleteTag(id: number): Promise<{ success: boolean }> {
  return request.delete<{ success: boolean }>('/api/tag', {
    data: { ids: [id] },
  })
}
