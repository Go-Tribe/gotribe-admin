import { type QueryClient } from '@tanstack/react-query'
import { createRootRouteWithContext, Outlet, redirect } from '@tanstack/react-router'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'
import { AxiosError } from 'axios'
import { Toaster } from '@/components/ui/sonner'
import { NavigationProgress } from '@/components/navigation-progress'
import { GeneralError } from '@/features/errors/general-error'
import { NotFoundError } from '@/features/errors/not-found-error'
import { useAuthStore } from '@/stores/auth-store'
import { getCurrentUser } from '@/service/user'

/**
 * 需要排除的认证相关路由
 * 这些路由不需要认证即可访问
 */
const AUTH_ROUTES = ['/sign-in', '/sign-up', '/otp']

/**
 * 根路由配置
 * 提供全局布局、错误处理和认证守卫
 */
export const Route = createRootRouteWithContext<{
  queryClient: QueryClient
}>()({
  /**
   * 路由加载前的认证检查
   * 在进入任何需要认证的路由前，验证 token 并获取用户信息
   */
  beforeLoad: async ({ location }) => {
    const { pathname } = location
    const { auth } = useAuthStore.getState()

    // 如果是认证相关页面，直接放行
    if (AUTH_ROUTES.some(route => pathname.startsWith(route))) {
      return
    }

    // 检查是否有 token
    if (!auth.accessToken) {
      // 没有 token，跳转到登录页，并保存当前路径用于登录后重定向
      throw redirect({
        to: '/sign-in',
        search: {
          redirect: pathname,
        },
      })
    }

    // 每次刷新页面都请求用户信息（除了登录/注册页面）
    // 这样可以确保用户信息是最新的，并且验证 token 是否仍然有效
    try {
      const result = await getCurrentUser()
      // 根据实际 API 返回的数据结构更新用户信息
      if (result && result.admin) {
        const userInfo = result.admin

        // 设置用户信息，不包含 accessToken（accessToken 单独存储）
        auth.setUser({
          ...userInfo,
        })
      }
    } catch (error: unknown) {
      // 获取用户信息失败
      // 检查是否是 401 错误（token 过期或无效）
      const isUnauthorized =
        error instanceof AxiosError && error.response?.status === 401

      if (isUnauthorized) {
        // 401 错误：token 已过期或无效
        // 注意：axios 拦截器已经清除了 token 并触发了事件，这里只需要跳转
        throw redirect({
          to: '/sign-in',
          search: {
            redirect: pathname,
          },
        })
      }

      // 其他错误（如网络错误、超时等）
      // 如果有 token 但没有用户信息，可能是刚登录成功但获取用户信息失败
      // 这种情况下允许继续访问，让页面组件自己处理用户信息获取
      // 这样可以避免因临时网络问题导致登录后无法跳转
      if (!auth.user && !auth.accessToken) {
        // 既没有用户信息也没有 token，跳转到登录页
        auth.reset()
        throw redirect({
          to: '/sign-in',
          search: {
            redirect: pathname,
          },
        })
      }
      // 有 token 但获取用户信息失败（可能是网络问题），允许继续访问
      // 页面组件可以自己处理用户信息的获取和显示
    }
  },
  /**
   * 根路由组件
   * 提供全局 UI 元素：导航进度条、Toast 通知、开发工具
   */
  component: () => {
    return (
      <>
        <NavigationProgress />
        <Outlet />
        <Toaster duration={5000} position="bottom-right" />
        {import.meta.env.MODE === 'development' && (
          <>
            <ReactQueryDevtools buttonPosition='bottom-left' />
            <TanStackRouterDevtools position='bottom-right' />
          </>
        )}
      </>
    )
  },
  /** 404 错误组件 */
  notFoundComponent: NotFoundError,
  /** 通用错误组件 */
  errorComponent: GeneralError,
})
