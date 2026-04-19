import { request } from '@/service'
import type { User, UserListResponse } from '../types/user'

function mapUser(raw: Omit<User, 'userID'>): User {
  return {
    ...raw,
    userID: String(raw.id),
  }
}

// 获取用户详情（单条）
export async function getUserDetail(id: number): Promise<User> {
  const data = await request.get<{ user: User }>(`/api/user/${id}`)
  return mapUser(data.user)
}

// 获取用户列表（支持分页和筛选）
export const getUserList = async (params?: {
  current?: number
  userID?: string
  projectId?: number
  pageNum?: number
  pageSize?: number
}) => {
  return request
    .get<UserListResponse>('/api/user', { params })
    .then((data) => ({
      ...data,
      users: (data.users ?? []).map((user) => mapUser(user)),
    }))
}

// 创建用户
export async function createUser(
  params: Partial<User>,
): Promise<{ success: boolean }> {
  return request.post<{ success: boolean }>('/api/user', params)
}

// 更新用户
export async function updateUser(
  params: Partial<User>,
): Promise<{ success: boolean }> {
  return request.patch<{ success: boolean }>(
    `/api/user/${params.id}`,
    params,
  )
}
