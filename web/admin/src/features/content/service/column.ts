import { request } from '@/service'
import type {
  Column,
  ColumnListParams,
  ColumnListResponse,
  ColumnCreateParams,
  ColumnUpdateParams,
} from '../types/column'

/**
 * 获取专栏列表
 * GET /api/column?projectID=&title=&pageNum=1&pageSize=3
 */
export async function getColumnList(
  params?: ColumnListParams
): Promise<ColumnListResponse> {
  const data = await request.get<ColumnListResponse>('/api/column', {
    params: {
      projectID: params?.projectID ?? '',
      title: params?.title ?? '',
      pageNum: params?.pageNum ?? 1,
      pageSize: params?.pageSize ?? 10,
    },
  })
  return data as ColumnListResponse
}

/**
 * 获取专栏详情（编辑时用于回显）
 * GET /api/column/:columnID
 * 响应: { code, message, data: { column: { columnID, title, info, projectID, icon, description, createdAt } } }
 */
export async function getColumn(columnID: string): Promise<Column> {
  const id = columnID.trim()
  const data = await request.get<{ column: Column }>(`/api/column/${id}`)
  return (data as { column: Column }).column
}

/**
 * 新增专栏
 * POST /api/column?projectID=xxx  body: title, description, info, projectID, icon
 */
export async function createColumn(
  data: ColumnCreateParams
): Promise<{ columnID?: string }> {
  const result = await request.post<{ columnID?: string }>('/api/column', data)
  return (result ?? {}) as { columnID?: string }
}

/**
 * 更新专栏
 * PATCH /api/column/:columnID
 */
export async function updateColumn(
  columnID: string,
  data: ColumnUpdateParams
): Promise<void> {
  const id = columnID.trim()
  await request.patch(`/api/column/${id}`, data)
}

/**
 * 删除专栏
 * DELETE /api/column  body: { columnIds: string }
 */
export async function deleteColumn(columnIds: string): Promise<void> {
  await request.delete('/api/column', {
    data: { ids: [Number(columnIds.trim())] },
  })
}
