import { request } from '@/service'
import type { Category, CategoryParams, BatchDeleteParams, CategoryListResponse } from '../types/category'

/**
 * 获取分类树
 * @returns 分类树形数据
 */
export function getCategoryTree() {
  return request.get<CategoryListResponse>('/api/category/tree')
}

/**
 * 创建分类
 * @param data 分类数据（接口要求 body: title, description, icon, path, sort, hidden, parentId）
 * @param options projectID 可选，作为 query 传递
 */
export function createCategory(
  data: CategoryParams,
  options?: { projectID?: string },
) {
  const params = new URLSearchParams()
  if (options?.projectID) params.set('projectID', options.projectID)
  const query = params.toString()
  const url = query ? `/api/category?${query}` : '/api/category'
  const body = {
    title: data.title,
    slug: data.slug,
    description: data.description ?? '',
    icon: data.icon ?? '',
    path: data.path ?? '', // 链接，接口字段为 path
    sort: data.sort ?? 1,
    status: data.status ?? 1,
    hidden: data.hidden ?? 1,
    parentID: data.parentId ?? 0,
  }
  return request.post<Category>(url, body)
}

/**
 * 更新分类（接口 PATCH /api/category/:id，body: title, slug, description, path, sort, status, hidden, parentID）
 * @param id 分类 ID（路径参数）
 * @param data 更新数据
 */
export function updateCategory(id: number, data: Partial<CategoryParams>) {
  const body = {
    title: data.title,
    slug: data.slug,
    description: data.description ?? '',
    icon: data.icon ?? '',
    path: data.path ?? '',
    sort: data.sort ?? 1,
    status: data.status ?? 1,
    hidden: data.hidden ?? 1,
    parentID: data.parentId ?? 0,
  }
  return request.patch<Category>(`/api/category/${id}`, body)
}

/**
 * 删除分类（接口要求 body: { ids: number[] }）
 * @param data 包含要删除的分类 ID 数组
 * @returns 删除结果
 */
export function batchDeleteCategory(data: BatchDeleteParams) {
  return request.delete<void>('/api/category', {
    data: { ids: data.ids },
  })
}
