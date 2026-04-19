import { request } from '@/service'

export interface LoginResult {
  [key: string]: unknown
}

export interface CurrentUserResult {
  admin: { [key: string]: unknown }
}

/** 获取当前用户 GET /api/admin/info；拦截器已解包，直接返回 data */
export async function getCurrentUser(): Promise<CurrentUserResult> {
  return request.get<CurrentUserResult>('/api/admin/info')
}

/** 退出登录 POST /api/base/logout */
export async function outLogin() {
  return request.post<LoginResult>('/api/base/logout')
}
