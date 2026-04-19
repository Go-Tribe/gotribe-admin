import { request } from '@/service'
export interface LoginParams {
  username: string
  password: string
}
export interface LoginResult {
  token: string
  expires: string
}

/** 登录接口 POST /api/base/login；返回 data 即 LoginResult */
export async function login(body: LoginParams): Promise<LoginResult> {
  return request.post<LoginResult>('/api/base/login', body)
}