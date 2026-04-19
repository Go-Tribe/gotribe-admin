import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider, createRouter } from '@tanstack/react-router'
import { authEventEmitter } from '@/lib/auth-events'
import { AppProviders } from '@/components/app-providers'
import { queryClient } from '@/lib/query-client'
// 环境变量验证（在应用启动时验证）
import './config/env'
// Generated Routes
import { routeTree } from './routeTree.gen'
// Styles
import './styles/index.css'
// i18n
import './i18n/config'

/**
 * 创建路由实例
 * 配置路由上下文和预加载策略
 */
const router = createRouter({
  routeTree,
  context: { queryClient },
  defaultPreload: 'intent',
  defaultPreloadStaleTime: 0,
  // 添加滚动恢复
  scrollRestoration: true,
})

// 注册路由类型，提供类型安全
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

/**
 * 订阅认证事件
 * 在 axios 拦截器中触发认证相关事件时，统一处理路由跳转
 * 这样可以避免在拦截器中直接使用 window.location，保持 SPA 路由状态
 */
authEventEmitter.subscribe((event) => {
  if (event.type === 'unauthorized') {
    // 未授权：跳转到登录页
    const redirect = event.redirect || router.history.location.href
    router.navigate({ to: '/sign-in', search: { redirect } })
  } else if (event.type === 'forbidden') {
    router.navigate({ to: '/403' })
  }
})

// 延迟加载非关键配置
if (typeof window !== 'undefined') {
  // 使用 requestIdleCallback 在浏览器空闲时加载
  const scheduleInit = () => {
    void import('./config/app')
      .then(({ updateAppConfigFromApi }) => updateAppConfigFromApi())
      .catch((error) => {
        // eslint-disable-next-line no-console
        console.warn('Failed to initialize app config:', error)
      })
  }

  if ('requestIdleCallback' in window) {
    window.requestIdleCallback(scheduleInit, { timeout: 3000 })
  } else {
    setTimeout(scheduleInit, 100)
  }
}

/**
 * 渲染应用
 * 在渲染前检查 root 元素是否已渲染，避免重复渲染
 */
const rootElement = document.getElementById('root')
if (!rootElement) {
  throw new Error('Root element not found')
}

// 使用 data 属性标记挂载状态，避免 innerHTML 检查不可靠（SSR 可能残留注释/空白）
if (!rootElement.hasAttribute('data-mounted')) {
  rootElement.setAttribute('data-mounted', 'true')
  const root = ReactDOM.createRoot(rootElement)
  root.render(
    <StrictMode>
      <AppProviders>
        <RouterProvider router={router} />
      </AppProviders>
    </StrictMode>
  )
}

// 注册 Service Worker（生产环境）
if (import.meta.env.PROD && 'serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker
      .register('/sw.js')
      .then((registration) => {
        // eslint-disable-next-line no-console
        console.log('SW registered:', registration.scope)
      })
      .catch((error) => {
        // eslint-disable-next-line no-console
        console.log('SW registration failed:', error)
      })
  })
}
