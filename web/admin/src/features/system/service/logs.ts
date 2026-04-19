import { request } from '@/service'
import type { OperationLog } from '../types/log'

export const getLogList = async (params?: {
  current?: number
  username?: string
  ip?: string
  path?: string
  status?: string
  pageNum?: number
  pageSize?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}) => {
  return request.get<{ logs: OperationLog[]; total: number }>(
    '/api/log/operation/list',
    { params }
  )
}

export async function deleteLog(id: number): Promise<{ success: boolean }> {
  return request.delete<{ success: boolean }>(
    '/api/log/operation/delete/batch',
    { data: { operationLogIds: [id] } }
  )
}
