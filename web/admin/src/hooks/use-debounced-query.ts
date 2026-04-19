import { useState, useEffect, useRef, useCallback } from 'react'
import {
  useQuery,
  useQueryClient,
  type UseQueryOptions,
  type QueryKey,
} from '@tanstack/react-query'

export interface UseDebouncedQueryOptions<TData = unknown>
  extends Omit<UseQueryOptions<TData>, 'queryKey' | 'queryFn'> {
  /** 查询 Key */
  queryKey: QueryKey
  /** 查询函数 */
  queryFn: () => Promise<TData>
  /** 防抖延迟（毫秒） */
  debounceMs?: number
  /** 是否启用防抖 */
  enabledDebounce?: boolean
}

export interface UseDebouncedQueryReturn<TData = unknown> {
  /** 查询数据 */
  data: TData | undefined
  /** 是否加载中 */
  isLoading: boolean
  /** 是否获取中 */
  isFetching: boolean
  /** 是否错误 */
  isError: boolean
  /** 错误对象 */
  error: Error | null
  /** 手动刷新 */
  refetch: () => void
  /** 取消防抖 */
  cancelDebounce: () => void
}

/**
 * 防抖查询 Hook
 * 
 * 性能优化：防抖处理频繁变化的查询参数，减少服务器请求
 * 
 * @example
 * ```tsx
 * // 基础用法
 * const { data, isLoading } = useDebouncedQuery({
 *   queryKey: ['users', filters],
 *   queryFn: () => fetchUsers(filters),
 *   debounceMs: 300,
 * })
 * 
 * // 表格筛选场景
 * const { data, isLoading } = useDebouncedQuery({
 *   queryKey: ['users', searchValue, page, pageSize],
 *   queryFn: () => getUserList({ search: searchValue, page, pageSize }),
 *   debounceMs: 300,
 * })
 * ```
 */
export function useDebouncedQuery<TData = unknown>({
  queryKey,
  queryFn,
  debounceMs = 300,
  enabledDebounce = true,
  ...options
}: UseDebouncedQueryOptions<TData>): UseDebouncedQueryReturn<TData> {
  const [debouncedKey, setDebouncedKey] = useState(queryKey)
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const queryClient = useQueryClient()
  // 使用 ref 存储上一次的 queryKey 引用，避免频繁 JSON.stringify
  const prevKeyRef = useRef(queryKey)

  // 防抖更新 queryKey
  useEffect(() => {
    if (!enabledDebounce) {
      // 只在值变化时更新，避免不必要的重新渲染
      // eslint-disable-next-line react-hooks/set-state-in-effect -- 防抖逻辑需要在 effect 中设置 state
      setDebouncedKey((prev) => {
        // 使用 ref 比较同一引用，若引用不同再用深度比较
        if (prevKeyRef.current === queryKey) return prev
        prevKeyRef.current = queryKey
        // 快速路径：基本类型或同一引用已处理，复杂对象做深度比较
        const prevStr = JSON.stringify(prev)
        const nextStr = JSON.stringify(queryKey)
        if (prevStr === nextStr) return prev
        return queryKey
      })
      return
    }

    // 清除之前的定时器
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
    }

    // 设置新的定时器
    timeoutRef.current = setTimeout(() => {
      prevKeyRef.current = queryKey
      setDebouncedKey(queryKey)
    }, debounceMs)

    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }
    }
  }, [queryKey, debounceMs, enabledDebounce])

  // 执行查询
  const queryResult = useQuery({
    queryKey: debouncedKey,
    queryFn,
    ...options,
  })

  // 手动刷新
  const refetch = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: debouncedKey })
  }, [queryClient, debouncedKey])

  // 取消防抖
  const cancelDebounce = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
      setDebouncedKey(queryKey)
    }
  }, [queryKey])

  return {
    data: queryResult.data,
    isLoading: queryResult.isLoading,
    isFetching: queryResult.isFetching,
    isError: queryResult.isError,
    error: queryResult.error,
    refetch,
    cancelDebounce,
  }
}

/**
 * 防抖搜索 Hook
 * 
 * 简化版的 useDebouncedQuery，专为搜索场景设计
 */
export interface UseDebouncedSearchOptions<TData = unknown> {
  /** 搜索关键词 */
  keyword: string
  /** 搜索函数 */
  searchFn: (keyword: string) => Promise<TData>
  /** 防抖延迟 */
  debounceMs?: number
  /** 最小关键词长度 */
  minLength?: number
}

export function useDebouncedSearch<TData = unknown>({
  keyword,
  searchFn,
  debounceMs = 300,
  minLength = 0,
}: UseDebouncedSearchOptions<TData>) {
  const shouldSearch = keyword.length >= minLength

  const { data, isLoading, isError, error } = useDebouncedQuery({
    queryKey: ['search', keyword],
    queryFn: () => searchFn(keyword),
    debounceMs,
    enabled: shouldSearch,
  })

  return {
    data,
    isLoading: shouldSearch && isLoading,
    isError,
    error,
    isEmpty: keyword.length > 0 && !isLoading && !data,
  }
}

/**
 * 防抖突变 Hook
 * 
 * 用于防抖提交表单等场景
 */
export function useDebouncedMutation<TData = unknown, TVariables = void>(
  mutationFn: (variables: TVariables) => Promise<TData>,
  debounceMs = 300
) {
  const timeoutRef = useRef<NodeJS.Timeout | null>(null)
  const [isPending, setIsPending] = useState(false)

  const mutate = useCallback(
    (variables: TVariables) => {
      return new Promise<TData>((resolve, reject) => {
        // 清除之前的定时器
        if (timeoutRef.current) {
          clearTimeout(timeoutRef.current)
        }

        // 设置新的定时器
        timeoutRef.current = setTimeout(async () => {
          setIsPending(true)
          try {
            const result = await mutationFn(variables)
            resolve(result)
          } catch (error) {
            reject(error)
          } finally {
            setIsPending(false)
          }
        }, debounceMs)
      })
    },
    [mutationFn, debounceMs]
  )

  const cancel = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
      setIsPending(false)
    }
  }, [])

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }
    }
  }, [])

  return { mutate, isPending, cancel }
}

export default useDebouncedQuery
