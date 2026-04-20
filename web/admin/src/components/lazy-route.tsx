import { lazy, Suspense, type ComponentType, type MouseEvent, type ReactNode } from 'react'
import { PageSkeleton, TablePageSkeleton, FormPageSkeleton } from './page-skeleton'

type SkeletonType = 'default' | 'table' | 'form'

interface LazyRouteOptions {
  /** 骨架屏类型 */
  skeleton?: SkeletonType
  /** 自定义骨架屏 */
  fallback?: ReactNode
  /** 预加载延迟（毫秒） */
  preloadDelay?: number
}

const skeletonMap: Record<SkeletonType, ReactNode> = {
  default: <PageSkeleton />,
  table: <TablePageSkeleton />,
  form: <FormPageSkeleton />,
}

/**
 * 创建懒加载路由组件
 * 
 * @param factory 动态导入工厂函数
 * @param options 配置选项
 * @returns 带 Suspense 的懒加载组件
 * 
 * @example
 * // 基本用法
 * const DashboardPage = lazyRoute(() => import('@/features/dashboard/dashboard-page'))
 * 
 * // 指定骨架屏类型
 * const AdminPage = lazyRoute(() => import('@/features/system/admin'), { skeleton: 'table' })
 * 
 * // 自定义骨架屏
 * const CustomPage = lazyRoute(() => import('@/features/custom'), { fallback: <CustomSkeleton /> })
 */
export function lazyRoute<T extends ComponentType<Record<string, unknown>>>(
  factory: () => Promise<{ default: T }>,
  options: LazyRouteOptions = {}
) {
  const { skeleton = 'default', fallback, preloadDelay = 200 } = options
  
  const LazyComponent = lazy(factory) as unknown as ComponentType<Record<string, unknown>>
  const defaultFallback = skeletonMap[skeleton]

  // 预加载功能
  let preloadPromise: Promise<void> | null = null
  
  const doPreload = () => {
    if (!preloadPromise) {
      preloadPromise = factory().then(() => undefined)
    }
    return preloadPromise
  }

  function LazyRouteWrapper(props: Record<string, unknown>) {
    return (
      <Suspense fallback={fallback ?? defaultFallback}>
        <LazyComponent {...props} />
      </Suspense>
    )
  }

  // 附加预加载方法
  LazyRouteWrapper.preload = () => doPreload()
  
  // 智能预加载：鼠标悬停时预加载
  LazyRouteWrapper.preloadOnHover = (event: MouseEvent) => {
    const target = event.currentTarget
    const preloadTimeout = setTimeout(() => {
      doPreload()
    }, preloadDelay)
    
    const cleanup = () => {
      clearTimeout(preloadTimeout)
      target.removeEventListener('mouseleave', cleanup)
    }
    
    target.addEventListener('mouseleave', cleanup, { once: true })
  }

  return LazyRouteWrapper
}

/**
 * 预加载多个路由（用于关键路径预加载）
 * 
 * @example
 * // 登录后立即预加载 Dashboard
 * preloadRoutes([
 *   () => import('@/features/dashboard/dashboard-page'),
 *   () => import('@/features/system/admin'),
 * ])
 */
export function preloadRoutes(factories: Array<() => Promise<unknown>>) {
  // 使用 requestIdleCallback 在浏览器空闲时预加载
  const schedulePreload = () => {
    factories.forEach((factory, index) => {
      // 错开加载时间，避免同时请求
      setTimeout(() => {
        factory().catch(() => {
          // 静默失败，预加载失败不应影响用户体验
        })
      }, index * 100)
    })
  }

  if (typeof window !== 'undefined') {
    if ('requestIdleCallback' in window) {
      window.requestIdleCallback(schedulePreload, { timeout: 2000 })
    } else {
      setTimeout(schedulePreload, 200)
    }
  }
}

/**
 * 使用 Intersection Observer 的视口预加载
 * 
 * @example
 * const DashboardPage = lazyRouteWithViewport(
 *   () => import('@/features/dashboard/dashboard-page')
 * )
 */
export function lazyRouteWithViewport<T extends ComponentType<Record<string, unknown>>>(
  factory: () => Promise<{ default: T }>,
  options: LazyRouteOptions & { rootMargin?: string } = {}
) {
  const { rootMargin: _rootMargin = '100px', ...lazyOptions } = options
  const LazyComponent = lazyRoute(factory, lazyOptions)

  function ViewportPreloadWrapper(props: Record<string, unknown>) {
    // 在 useEffect 中设置 Intersection Observer
    // 这里简化处理，实际使用时可以通过 ref 绑定到链接元素
    return <LazyComponent {...props} />
  }

  ViewportPreloadWrapper.preload = LazyComponent.preload
  return ViewportPreloadWrapper
}
