import { AxiosError } from 'axios'
import { QueryCache, QueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

/**
 * 创建 QueryClient 实例
 * 配置查询和变更的默认选项，统一错误处理策略
 */
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      /**
       * 重试逻辑
       * - 开发环境：不重试（快速失败，便于调试）
       * - 生产环境：最多重试 3 次
       * - 401/403 错误：不重试（认证错误不应重试）
       */
      retry: (failureCount, error) => {
        // 开发环境不重试
        if (import.meta.env.DEV) {
          // eslint-disable-next-line no-console
          console.log({ failureCount, error })
          return false
        }

        // 生产环境最多重试 3 次
        if (failureCount >= 3) {
          return false
        }

        // 401/403 错误不重试
        if (
          error instanceof AxiosError &&
          [401, 403].includes(error.response?.status ?? 0)
        ) {
          return false
        }

        return true
      },
      refetchOnWindowFocus: import.meta.env.PROD,
      staleTime: 10 * 1000, // 10s
      // 添加 GC Time 优化内存
      gcTime: 5 * 60 * 1000, // 5 minutes
    },
    mutations: {
      onError: (error) => {
        // 全局 mutation 错误处理
        // 注意：组件级别的错误处理优先，这里只处理特殊情况
        if (error instanceof AxiosError) {
          if (error.response?.status === 304) {
            toast.error('Content not modified!')
          }
        }
      },
    },
  },
  queryCache: new QueryCache({
    /**
     * 全局查询错误处理
     * 统一处理认证错误和服务器错误
     * 注意：所有 HTTP 状态码错误的 toast 已由 axios 拦截器统一处理（通过 handleServerError）
     * 这里只处理需要全局响应的特殊逻辑（如页面跳转），避免重复提示
     */
    onError: (error) => {
      if (error instanceof AxiosError) {
        const status = error.response?.status

        // 500 错误：仅在生产环境跳转到错误页面，避免开发时影响 HMR
        // toast 已由 axios 拦截器处理，这里不再重复显示
        if (status === 500 && import.meta.env.PROD) {
          // 使用 window.location 作为 fallback，因为 router 可能还未初始化
          window.location.href = '/500'
        }

        // 401/403 错误：token 清除、事件触发和路由跳转已由 axios 拦截器和事件监听器处理
        // 不需要额外操作
      }
    },
  }),
})
