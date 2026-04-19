// {
//   "code": 200,
//   "message": "Get list successfully",
//   "data": {
//       "apis": [
//           {
//               "id": 126,
//               "createdAt": "2026-01-06T14:43:21.337972+08:00",
//               "updatedAt": "2026-01-06T14:43:21.337972+08:00",
//               "deletedAt": null,
//               "method": "GET",
//               "path": "/index/data",
//               "category": "Index",
//               "desc": "获取首页时间数据",
//               "creator": "系统"
//           },
//           {
//               "id": 125,
//               "createdAt": "2026-01-06T14:43:21.180331+08:00",
//               "updatedAt": "2026-01-06T14:43:21.180331+08:00",
//               "deletedAt": null,
//               "method": "GET",
//               "path": "/index",
//               "category": "Index",
//               "desc": "获取首页当日数据",
//               "creator": "系统"
//           },
//           {
//               "id": 124,
//               "createdAt": "2026-01-06T14:43:21.024844+08:00",
//               "updatedAt": "2026-01-06T14:43:21.024844+08:00",
//               "deletedAt": null,
//               "method": "PATCH",
//               "path": "/order/logistics/:orderID",
//               "category": "order",
//               "desc": "设置物流信息",
//               "creator": "系统"
//           },
//           {
//               "id": 123,
//               "createdAt": "2026-01-06T14:43:20.87868+08:00",
//               "updatedAt": "2026-01-06T14:43:20.87868+08:00",
//               "deletedAt": null,
//               "method": "GET",
//               "path": "/feedback",
//               "category": "feedback",
//               "desc": "获取反馈列表",
//               "creator": "系统"
//           },
//           {
//               "id": 122,
//               "createdAt": "2026-01-06T14:43:20.729218+08:00",
//               "updatedAt": "2026-01-06T14:43:20.729218+08:00",
//               "deletedAt": null,
//               "method": "PATCH",
//               "path": "/system",
//               "category": "system",
//               "desc": "更新配置",
//               "creator": "系统"
//           },
//           {
//               "id": 121,
//               "createdAt": "2026-01-06T14:43:20.578544+08:00",
//               "updatedAt": "2026-01-06T14:43:20.578544+08:00",
//               "deletedAt": null,
//               "method": "GET",
//               "path": "/base/config",
//               "category": "base",
//               "desc": "获取后台配置",
//               "creator": "系统"
//           },
//           {
//               "id": 120,
//               "createdAt": "2026-01-06T14:43:20.431437+08:00",
//               "updatedAt": "2026-01-06T14:43:20.431437+08:00",
//               "deletedAt": null,
//               "method": "DELETE",
//               "path": "/order",
//               "category": "order",
//               "desc": "删除订单",
//               "creator": "系统"
//           },
//           {
//               "id": 119,
//               "createdAt": "2026-01-06T14:43:20.282893+08:00",
//               "updatedAt": "2026-01-06T14:43:20.282893+08:00",
//               "deletedAt": null,
//               "method": "PATCH",
//               "path": "/order/:orderID",
//               "category": "order",
//               "desc": "修改订单信息",
//               "creator": "系统"
//           },
//           {
//               "id": 118,
//               "createdAt": "2026-01-06T14:43:20.136576+08:00",
//               "updatedAt": "2026-01-06T14:43:20.136576+08:00",
//               "deletedAt": null,
//               "method": "GET",
//               "path": "/order",
//               "category": "order",
//               "desc": "获取订单列表",
//               "creator": "系统"
//           },
//           {
//               "id": 117,
//               "createdAt": "2026-01-06T14:43:19.985384+08:00",
//               "updatedAt": "2026-01-06T14:43:19.985384+08:00",
//               "deletedAt": null,
//               "method": "GET",
//               "path": "/order/log/:orderID",
//               "category": "order",
//               "desc": "获取订单记录",
//               "creator": "系统"
//           }
//       ],
//       "total": 126
//   }
// }
import { request } from '@/service'
import type { Api } from '../types/api'

export const getApiList = async (params?: {
  current?: number
  path?: string
  category?: string
  method?: string
  creator?: string
  pageNum?: number
  pageSize?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}) => {
  return request.get<{ apis: Api[]; total: number }>('/api/api/list', { params })
}

export async function createApi(
  params: Partial<Api>,
): Promise<{ success: boolean }> {
  return request.post<{ success: boolean }>('/api/api/create', params)
}

export async function updateApi(
  params: Partial<Api>,
): Promise<{ success: boolean }> {
  return request.patch<{ success: boolean }>(`/api/api/update/${params.id}`, params)
}

export async function deleteApi(id: number): Promise<{ success: boolean }> {
  return request.delete<{ success: boolean }>('/api/api/delete/batch', {
    data: { apiIds: [id] },
  })
}
