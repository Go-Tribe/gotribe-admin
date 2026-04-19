import { useState, useEffect, useCallback, useRef } from 'react'

interface UseCachedFetchOptions<T> {
  /** 缓存时间（毫秒），默认 5 分钟 */
  cacheTime?: number
  /** 防抖延迟（毫秒），默认 300ms */
  debounceMs?: number
  /** 页面重新可见时是否重新获取，默认 true */
  refetchOnVisible?: boolean
  /** 重新获取的冷却时间（毫秒），默认 30 秒 */
  refetchCooldown?: number
  /** 初始数据 */
  initialData?: T
}

interface CachedData<T> {
  data: T
  timestamp: number
}

const cache = new Map<string, CachedData<unknown>>()
/** 最大缓存条目数，防止无限制增长 */
const MAX_CACHE_SIZE = 50

/** LRU 淘汰：当缓存超过上限时，删除最久未使用的条目 */
function enforceLRU(): void {
  if (cache.size <= MAX_CACHE_SIZE) return
  // Map 按插入顺序维护键，删除最早的条目
  const firstKey = cache.keys().next().value
  if (firstKey !== undefined) {
    cache.delete(firstKey)
  }
}

/**
 * 带缓存的数据获取 Hook
 *
 * 特性：
 * 1. 内存级缓存，减少重复请求
 * 2. LRU 淘汰机制，防止内存无限增长
 * 3. 防抖处理，避免频繁触发
 * 4. 页面可见性变化时智能重新获取
 * 5. 冷却时间控制，防止过度刷新
 * 
 * @param key 缓存键
 * @param fetchFn 数据获取函数
 * @param options 配置选项
 * 
 * @example
 * const { data, isLoading, error, refetch } = useCachedFetch(
 *   'menu',
 *   () => getMenuAccessTree(userId),
 *   { cacheTime: 60000 } // 缓存 1 分钟
 * )
 */
export function useCachedFetch<T>(
  key: string | null,
  fetchFn: () => Promise<T>,
  options: UseCachedFetchOptions<T> = {}
) {
  const {
    cacheTime = 5 * 60 * 1000, // 5 分钟
    debounceMs = 300,
    refetchOnVisible = true,
    refetchCooldown = 30 * 1000, // 30 秒
    initialData,
  } = options

  const [data, setData] = useState<T | undefined>(initialData)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<Error | null>(null)
  
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const lastFetchRef = useRef<number>(0)
  const isMountedRef = useRef(true)

  const fetchData = useCallback(async (force = false) => {
    if (!key) return

    // 检查缓存
    if (!force) {
      const cached = cache.get(key) as CachedData<T> | undefined
      if (cached && Date.now() - cached.timestamp < cacheTime) {
        setData(cached.data)
        // 更新缓存顺序以实现 LRU
        cache.delete(key)
        cache.set(key, cached)
        return
      }
    }

    // 防抖处理
    if (debounceRef.current) {
      clearTimeout(debounceRef.current)
    }

    debounceRef.current = setTimeout(async () => {
      setIsLoading(true)
      setError(null)

      try {
        const result = await fetchFn()
        
        if (isMountedRef.current) {
          setData(result)
          // 更新缓存
          cache.set(key, { data: result, timestamp: Date.now() })
          enforceLRU()
          lastFetchRef.current = Date.now()
        }
      } catch (err) {
        if (isMountedRef.current) {
          setError(err instanceof Error ? err : new Error(String(err)))
        }
      } finally {
        if (isMountedRef.current) {
          setIsLoading(false)
        }
      }
    }, debounceMs)
  }, [key, fetchFn, cacheTime, debounceMs])

  // 初始获取
  useEffect(() => {
    isMountedRef.current = true
    fetchData()

    return () => {
      isMountedRef.current = false
      if (debounceRef.current) {
        clearTimeout(debounceRef.current)
      }
    }
  }, [fetchData])

  // 监听页面可见性变化
  useEffect(() => {
    if (!refetchOnVisible || !key) return

    const handleVisibilityChange = () => {
      if (document.visibilityState === 'visible') {
        // 检查冷却时间
        const timeSinceLastFetch = Date.now() - lastFetchRef.current
        if (timeSinceLastFetch > refetchCooldown) {
          fetchData(true)
        }
      }
    }

    document.addEventListener('visibilitychange', handleVisibilityChange)
    return () => document.removeEventListener('visibilitychange', handleVisibilityChange)
  }, [fetchData, refetchOnVisible, refetchCooldown, key])

  // 手动刷新
  const refetch = useCallback(() => {
    fetchData(true)
  }, [fetchData])

  // 清除缓存
  const clearCache = useCallback(() => {
    if (key) {
      cache.delete(key)
    }
  }, [key])

  return {
    data,
    isLoading,
    error,
    refetch,
    clearCache,
  }
}

/**
 * 清除指定 key 的缓存
 */
export function clearFetchCache(key: string): void {
  cache.delete(key)
}

/**
 * 清除所有缓存
 */
export function clearAllFetchCache(): void {
  cache.clear()
}

/**
 * 预加载数据（用于预加载关键数据）
 */
export function prefetchData<T>(key: string, fetchFn: () => Promise<T>): Promise<void> {
  return fetchFn().then((data) => {
    cache.set(key, { data, timestamp: Date.now() })
    enforceLRU()
  }).catch(() => {
    // 预加载失败静默处理
  })
}
