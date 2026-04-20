import { useState, useEffect, useRef, memo, useCallback } from 'react'
import { cn } from '@/lib/utils'
import { Skeleton } from '@/components/ui/skeleton'
import { ImageIcon, AlertCircle } from 'lucide-react'

export interface LazyImageProps {
  /** 图片地址 */
  src: string
  /** 替代文本 */
  alt?: string
  /** 宽度 */
  width?: number | string
  /** 高度 */
  height?: number | string
  /** 自定义类名 */
  className?: string
  /** 图片样式类名 */
  imageClassName?: string
  /** 占位组件 */
  placeholder?: React.ReactNode
  /** 加载失败显示 */
  fallback?: React.ReactNode
  /** 视口交叉阈值（0-1） */
  threshold?: number
  /** 根元素外边距（提前加载距离） */
  rootMargin?: string
  /** 图片加载完成回调 */
  onLoad?: () => void
  /** 图片加载失败回调 */
  onError?: (error: Error) => void
  /** 是否开启懒加载 */
  lazy?: boolean
  /** 填充模式 */
  objectFit?: 'cover' | 'contain' | 'fill' | 'none' | 'scale-down'
  /** 圆角 */
  rounded?: boolean | 'sm' | 'md' | 'lg' | 'full'
}

type ImageState = 'loading' | 'loaded' | 'error'

const roundedMap = {
  true: 'rounded-md',
  sm: 'rounded-sm',
  md: 'rounded-md',
  lg: 'rounded-lg',
  full: 'rounded-full',
}

/**
 * 懒加载图片组件
 * 
 * 性能优化：
 * 1. 使用 Intersection Observer 实现视口内懒加载
 * 2. 图片未进入视口时不加载
 * 3. 支持占位图和错误处理
 * 
 * @example
 * ```tsx
 * // 基础用法
 * <LazyImage src="https://example.com/image.jpg" width={200} height={150} />
 * 
 * // 自定义占位
 * <LazyImage
 *   src="large-image.jpg"
   width={400}
 *   height={300}
 *   placeholder={<CustomSkeleton />}
 *   threshold={0.1}
 *   rootMargin="100px"
 * />
 * 
 * // 头像（全圆角）
 * <LazyImage
 *   src={avatarUrl}
 *   width={40}
 *   height={40}
 *   rounded="full"
 *   objectFit="cover"
 * />
 * ```
 */
export const LazyImage = memo(function LazyImage({
  src,
  alt = '',
  width = 200,
  height = 150,
  className,
  imageClassName,
  placeholder,
  fallback,
  threshold = 0.1,
  rootMargin = '50px',
  onLoad,
  onError,
  lazy = true,
  objectFit = 'cover',
  rounded = true,
}: LazyImageProps) {
  const [state, setState] = useState<ImageState>('loading')
  const [isInView, setIsInView] = useState(!lazy)
  const imgRef = useRef<HTMLImageElement>(null)
  const containerRef = useRef<HTMLDivElement>(null)

  // 使用 ref 存储回调，避免父组件传递内联函数导致 useEffect 重复触发
  const onLoadRef = useRef(onLoad)
  const onErrorRef = useRef(onError)

  useEffect(() => {
    onLoadRef.current = onLoad
  }, [onLoad])

  useEffect(() => {
    onErrorRef.current = onError
  }, [onError])

  // Intersection Observer 监听
  useEffect(() => {
    if (!lazy || isInView) return

    const element = containerRef.current
    if (!element) return

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsInView(true)
          observer.disconnect()
        }
      },
      {
        threshold,
        rootMargin,
      }
    )

    observer.observe(element)

    return () => {
      observer.disconnect()
    }
  }, [lazy, isInView, threshold, rootMargin])

  // 加载图片
  useEffect(() => {
    if (!isInView || !src) return

    const img = new Image()
    img.src = src

    img.onload = () => {
      setState('loaded')
      onLoadRef.current?.()
    }

    img.onerror = () => {
      setState('error')
      onErrorRef.current?.(new Error(`Failed to load image: ${src}`))
    }
  }, [isInView, src])

  // 计算圆角类名
  const roundedClass = typeof rounded === 'boolean' 
    ? (rounded ? roundedMap.true : '') 
    : roundedMap[rounded]

  // 占位组件
  const defaultPlaceholder = (
    <Skeleton 
      className={cn(
        'flex items-center justify-center bg-muted',
        roundedClass
      )}
      style={{ width, height }}
    >
      <ImageIcon className="h-8 w-8 text-muted-foreground/50" />
    </Skeleton>
  )

  // 错误组件
  const defaultFallback = (
    <div 
      className={cn(
        'flex flex-col items-center justify-center bg-muted text-muted-foreground',
        roundedClass
      )}
      style={{ width, height }}
    >
      <AlertCircle className="h-8 w-8 mb-1" />
      <span className="text-xs">加载失败</span>
    </div>
  )

  return (
    <div 
      ref={containerRef}
      className={cn('relative overflow-hidden', roundedClass, className)}
      style={{ width, height }}
    >
      {state === 'loading' && (placeholder || defaultPlaceholder)}
      
      {state === 'error' && (fallback || defaultFallback)}
      
      {(state === 'loaded' || isInView) && (
        <img
          ref={imgRef}
          src={src}
          alt={alt}
          className={cn(
            'transition-opacity duration-300',
            state === 'loaded' ? 'opacity-100' : 'opacity-0',
            roundedClass,
            imageClassName
          )}
          style={{ 
            width, 
            height, 
            objectFit,
          }}
        />
      )}
    </div>
  )
})

/** 图片预览组件 */
export interface ImagePreviewProps extends Omit<LazyImageProps, 'onClick'> {
  /** 点击预览 */
  onPreview?: (src: string) => void
  /** 是否可预览 */
  previewable?: boolean
}

export const ImagePreview = memo(function ImagePreview({
  onPreview,
  previewable = true,
  className,
  ...props
}: ImagePreviewProps) {
  const handleClick = useCallback(() => {
    if (previewable && onPreview) {
      onPreview(props.src)
    }
  }, [previewable, onPreview, props.src])

  return (
    <div 
      className={cn(
        'relative cursor-pointer group',
        previewable && 'hover:opacity-90',
        className
      )}
      onClick={handleClick}
    >
      <LazyImage {...props} />
      {previewable && (
        <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity bg-black/30 rounded-md">
          <span className="text-white text-sm">查看</span>
        </div>
      )}
    </div>
  )
})

export default LazyImage
