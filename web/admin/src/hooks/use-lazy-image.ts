import { useState, useEffect, useRef, useCallback } from 'react'

interface UseLazyImageOptions {
  /** 根元素的 margin，用于提前加载 */
  rootMargin?: string
  /** 阈值（0-1），元素可见比例达到时触发加载 */
  threshold?: number
  /** 默认占位图片 */
  placeholder?: string
  /** 加载失败时的 fallback 图片 */
  fallback?: string
  /** 是否只加载一次 */
  once?: boolean
}

interface UseLazyImageReturn {
  /** 图片引用，需要绑定到 img 元素 */
  ref: React.RefObject<HTMLImageElement | null>
  /** 当前图片 src */
  src: string | undefined
  /** 是否正在加载 */
  isLoading: boolean
  /** 是否加载完成 */
  isLoaded: boolean
  /** 是否加载失败 */
  isError: boolean
  /** 手动重新加载 */
  reload: () => void
}

/**
 * 图片懒加载 Hook
 * 使用 Intersection Observer 实现
 * 
 * @example
 * function LazyImage({ src, alt }: { src: string; alt: string }) {
 *   const { ref, src: lazySrc, isLoading } = useLazyImage({ 
 *     rootMargin: '50px',
 *     placeholder: '/placeholder.png'
 *   })
 *   
 *   return (
 *     <div className="relative">
 *       {isLoading && <Skeleton className="absolute inset-0" />}
 *       <img ref={ref} src={lazySrc} alt={alt} className={isLoading ? 'opacity-0' : 'opacity-100'} />
 *     </div>
 *   )
 * }
 */
export function useLazyImage(
  src: string | undefined,
  options: UseLazyImageOptions = {}
): UseLazyImageReturn {
  const {
    rootMargin = '50px',
    threshold = 0,
    placeholder,
    fallback,
    once = true,
  } = options

  const imgRef = useRef<HTMLImageElement>(null)
  const [isInView, setIsInView] = useState(false)
  const [isLoaded, setIsLoaded] = useState(false)
  const [isError, setIsError] = useState(false)
  const [loadedSrc, setLoadedSrc] = useState<string | undefined>(undefined)

  // 计算当前应该显示的图片 src
  const currentSrc = loadedSrc || (src && isInView ? undefined : placeholder)

  // Intersection Observer 检测元素是否在视口内
  useEffect(() => {
    const element = imgRef.current
    if (!element || !src) return

    // 如果已经加载过且只需要加载一次，直接返回
    if (once && isInView) return

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setIsInView(true)
            if (once) {
              observer.unobserve(element)
            }
          } else if (!once) {
            setIsInView(false)
          }
        })
      },
      { rootMargin, threshold }
    )

    observer.observe(element)

    return () => observer.disconnect()
  }, [src, rootMargin, threshold, once, isInView])

  // 加载图片
  useEffect(() => {
    if (!isInView || !src) {
      // 使用 requestAnimationFrame 延迟状态重置，避免同步 setState
      const rafId = requestAnimationFrame(() => {
        setLoadedSrc(undefined)
        setIsLoaded(false)
        setIsError(false)
      })
      return () => cancelAnimationFrame(rafId)
    }

    let cancelled = false
    const img = new Image()
    
    const handleLoad = () => {
      if (cancelled) return
      setLoadedSrc(src)
      setIsLoaded(true)
      setIsError(false)
    }

    const handleError = () => {
      if (cancelled) return
      setIsError(true)
      if (fallback) {
        setLoadedSrc(fallback)
      }
      setIsLoaded(true)
    }

    img.addEventListener('load', handleLoad)
    img.addEventListener('error', handleError)
    img.src = src

    // 如果图片已经缓存，onload 可能不会触发
    if (img.complete) {
      handleLoad()
    }

    return () => {
      cancelled = true
      img.removeEventListener('load', handleLoad)
      img.removeEventListener('error', handleError)
    }
  }, [isInView, src, fallback])

  const reload = useCallback(() => {
    if (src) {
      setIsError(false)
      setIsLoaded(false)
      const img = new Image()
      const handleLoad = () => {
        setLoadedSrc(src)
        setIsLoaded(true)
      }
      const handleError = () => {
        setIsError(true)
        if (fallback) {
          setLoadedSrc(fallback)
        }
      }
      img.addEventListener('load', handleLoad)
      img.addEventListener('error', handleError)
      img.src = src + '?t=' + Date.now() // 添加时间戳避免缓存
    }
  }, [src, fallback])

  return {
    ref: imgRef,
    src: currentSrc,
    isLoading: isInView && !isLoaded,
    isLoaded,
    isError,
    reload,
  }
}

/**
 * 预加载图片
 * 
 * @example
 * preloadImage('/image.png').then(() => console.log('loaded'))
 */
export function preloadImage(src: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.onload = () => resolve()
    img.onerror = reject
    img.src = src
  })
}

/**
 * 预加载多张图片
 * 
 * @example
 * preloadImages(['/1.png', '/2.png']).then(() => console.log('all loaded'))
 */
export function preloadImages(srcs: string[]): Promise<void> {
  return Promise.all(srcs.map(preloadImage)).then(() => undefined)
}
