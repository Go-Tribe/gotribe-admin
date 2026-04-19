import { request } from '@/service'
import type { Role } from '../types/admin'
import type { Menu } from '../types/menu'
import type { ApiListResponse } from '../types/role'

export async function getRoleMenus(roleId: number): Promise<{ menus: Menu[] }> {
  return request.get<{ menus: Menu[] }>(`/api/role/menus/get/${roleId}`)
}

export async function updateRoleMenus(
  roleId: number,
  data: { menuIds: number[] },
): Promise<{ success: boolean }> {
  return request.patch<{ success: boolean }>(`/api/role/menus/update/${roleId}`, data)
}

export async function getRoleApis(roleId: number): Promise<{ apis: { id: number }[] }> {
  return request.get<{ apis: { id: number }[] }>(`/api/role/apis/get/${roleId}`)
}

export async function updateRoleApis(
  roleId: number,
  data: { apiIds: number[] },
): Promise<{ success: boolean }> {
  return request.patch<{ success: boolean }>(`/api/role/apis/update/${roleId}`, data)
}

export async function deleteRole(id: number): Promise<{ success: boolean }> {
  return request.delete<{ success: boolean }>('/api/role/delete/batch', {
    data: { roleIds: [id] },
  })
}

export async function batchDeleteRoles(roleIds: number[]): Promise<{ success: boolean }> {
  return request.delete<{ success: boolean }>('/api/role/delete/batch', {
    data: { roleIds },
  })
}

export const getRoleList = async (params?: {
  current?: number
  name?: string
  keyword?: string
  status?: string
  pageNum?: number
  pageSize?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}) => {
  return request.get<{ roles: Role[]; total: number }>('/api/role/list', { params })
}

export async function getApis(): Promise<ApiListResponse> {
  return request.get<ApiListResponse>('/api/api/tree')
}

export async function createRole(
  params: Partial<Role>,
): Promise<{ success: boolean }> {
  return request.post<{ success: boolean }>('/api/role/create', params)
}

export async function updateRole(
  params: Partial<Role>,
): Promise<{ success: boolean }> {
  return request.patch<{ success: boolean }>(`/api/role/update/${params.id}`, params)
}
