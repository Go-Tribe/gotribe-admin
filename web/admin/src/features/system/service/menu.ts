import { request } from '@/service'
import type {
  Menu,
  MenuTreeResponse,
  MenuListResponse,
  CreateMenuParams,
  BatchDeleteMenuParams,
  UserMenuResponse,
  UserMenuTreeResponse,
} from '../types/menu'

export async function getMenuTree(): Promise<MenuTreeResponse> {
  return request.get<MenuTreeResponse>('/api/menu/tree')
}

export async function getMenus(): Promise<MenuListResponse> {
  return request.get<MenuListResponse>('/api/menu/list')
}

export async function createMenu(
  params: CreateMenuParams,
): Promise<{ success: boolean }> {
  return request.post<{ success: boolean }>('/api/menu/create', params)
}

export async function updateMenuById(
  id: number,
  params: Partial<Menu>,
): Promise<{ success: boolean }> {
  return request.patch<{ success: boolean }>(`/api/menu/update/${id}`, params)
}

export async function batchDeleteMenuByIds(
  params: BatchDeleteMenuParams,
): Promise<{ success: boolean }> {
  return request.delete<{ success: boolean }>('/api/menu/delete/batch', {
    data: params,
  })
}

export async function getUserMenusByUserId(
  userId: number,
): Promise<UserMenuResponse> {
  return request.get<UserMenuResponse>(`/api/menu/access/list/${userId}`)
}

export async function getUserMenuTreeByUserId(
  userId: number,
): Promise<UserMenuTreeResponse> {
  return request.get<UserMenuTreeResponse>(`/api/menu/access/tree/${userId}`)
}
