import { create } from 'zustand'

const ACCESS_TOKEN_KEY = 'accessToken'
const USER_INFO_KEY = 'user'

/**
 * 认证用户信息接口
 * 基于 Admin 类型，但允许扩展其他字段
 * 使用 Record 类型提供更好的类型安全性和灵活性
 */
export interface AuthUser extends Record<string, unknown> {
  id?: number
  username?: string
  mobile?: string
  avatar?: string
  nickname?: string
  introduction?: string
  status?: number
  creator?: string
  roleIds?: number[]
}

/**
 * 认证状态接口
 * 管理用户认证信息和 token
 */
export interface AuthState {
  auth: {
    /** 当前用户信息 */
    user: AuthUser | null
    /** 设置用户信息 */
    setUser: (user: AuthUser) => void
    /** 访问令牌 */
    accessToken: string
    /** 设置访问令牌 */
    setAccessToken: (accessToken: string) => void
    /** 重置访问令牌 */
    resetAccessToken: () => void
    /** 重置所有认证信息 */
    reset: () => void
  }
}

/**
 * localStorage 工具函数
 * 提供安全的 token 和用户信息存储/读取
 */

/**
 * 从 localStorage 获取访问令牌
 * @returns 访问令牌，如果不存在或出错则返回空字符串
 */
function getAccessTokenFromStorage(): string {
  if (typeof window === 'undefined') return ''
  try {
    const token = localStorage.getItem(ACCESS_TOKEN_KEY)
    return token || ''
  } catch {
    return ''
  }
}

/**
 * 将访问令牌保存到 localStorage
 * @param accessToken 访问令牌
 */
function setAccessTokenToStorage(accessToken: string): void {
  if (typeof window === 'undefined') return
  try {
    if (accessToken) {
      localStorage.setItem(ACCESS_TOKEN_KEY, accessToken)
    } else {
      localStorage.removeItem(ACCESS_TOKEN_KEY)
    }
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Failed to save accessToken to localStorage:', error)
  }
}

/**
 * 从 localStorage 获取用户信息
 * @returns 用户信息对象，如果不存在或出错则返回 null
 */
function getUserInfoFromStorage(): AuthUser | null {
  if (typeof window === 'undefined') return null
  try {
    const userInfo = localStorage.getItem(USER_INFO_KEY)
    return userInfo ? JSON.parse(userInfo) : null
  } catch {
    return null
  }
}

/**
 * 将用户信息保存到 localStorage
 * @param user 用户信息对象
 */
function setUserInfoToStorage(user: AuthUser | null): void {
  if (typeof window === 'undefined') return
  try {
    if (user) {
      localStorage.setItem(USER_INFO_KEY, JSON.stringify(user))
    } else {
      localStorage.removeItem(USER_INFO_KEY)
    }
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Failed to save user info to localStorage:', error)
  }
}

/**
 * 认证状态管理 Store
 * 使用 zustand 管理全局认证状态，包括用户信息和访问令牌
 * 自动同步到 localStorage，确保刷新页面后状态保持
 */
export const useAuthStore = create<AuthState>()((set) => {
  // 从 localStorage 读取初始 token 和用户信息
  const initToken = getAccessTokenFromStorage()
  const initUser = getUserInfoFromStorage()

  return {
    auth: {
      user: initUser,
      setUser: (user) =>
        set((state) => {
          setUserInfoToStorage(user)
          return { ...state, auth: { ...state.auth, user } }
        }),
      accessToken: initToken,
      setAccessToken: (accessToken) =>
        set((state) => {
          setAccessTokenToStorage(accessToken)
          return { ...state, auth: { ...state.auth, accessToken } }
        }),
      resetAccessToken: () =>
        set((state) => {
          setAccessTokenToStorage('')
          return { ...state, auth: { ...state.auth, accessToken: '' } }
        }),
      reset: () =>
        set((state) => {
          setAccessTokenToStorage('')
          setUserInfoToStorage(null)
          return {
            ...state,
            auth: { ...state.auth, user: null, accessToken: '' },
          }
        }),
    },
  }
})

// ==========================================
// 优化的 Selector Hooks - 避免不必要的重渲染
// ==========================================

/**
 * 获取当前登录用户信息
 * 仅当 user 对象变化时才会触发重渲染
 * 
 * @example
 * const user = useAuthUser()
 * // user 对象的引用稳定，可以安全用于依赖数组
 */
export function useAuthUser(): AuthUser | null {
  return useAuthStore(state => state.auth.user)
}

/**
 * 获取访问令牌
 * 仅当 token 变化时才会触发重渲染
 * 
 * @example
 * const token = useAccessToken()
 */
export function useAccessToken(): string {
  return useAuthStore(state => state.auth.accessToken)
}

/**
 * 检查用户是否已登录
 * 返回布尔值，性能最优
 * 
 * @example
 * const isAuthenticated = useIsAuthenticated()
 */
export function useIsAuthenticated(): boolean {
  return useAuthStore(state => !!state.auth.accessToken)
}

/**
 * 获取用户信息选择器（选择特定字段）
 * 
 * @example
 * // 只订阅 id 和 username
 * const { id, username } = useAuthUserSelect(user => ({ 
 *   id: user?.id, 
 *   username: user?.username 
 * }))
 */
export function useAuthUserSelect<T>(selector: (user: AuthUser | null) => T): T {
  return useAuthStore((state) => selector(state.auth.user))
}

/**
 * 获取 setUser 方法
 * 引用稳定，可用于依赖数组
 * 
 * @example
 * const setUser = useSetAuthUser()
 */
export function useSetAuthUser(): (user: AuthUser) => void {
  return useAuthStore(state => state.auth.setUser)
}

/**
 * 获取 setAccessToken 方法
 * 引用稳定，可用于依赖数组
 * 
 * @example
 * const setToken = useSetAccessToken()
 */
export function useSetAccessToken(): (accessToken: string) => void {
  return useAuthStore(state => state.auth.setAccessToken)
}

/**
 * 获取 reset 方法（登出使用）
 * 引用稳定，可用于依赖数组
 * 
 * @example
 * const logout = useLogout()
 * <button onClick={logout}>退出登录</button>
 */
export function useLogout(): () => void {
  return useAuthStore(state => state.auth.reset)
}

/**
 * 获取完整的 auth 对象（向后兼容，但不推荐在新代码中使用）
 * 会订阅所有 auth 字段的变化
 * 
 * @deprecated 建议使用上面的细粒度 hooks
 */
export function useAuthState(): AuthState['auth'] {
  return useAuthStore(state => state.auth)
}
