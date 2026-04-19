import { useState, useMemo, useCallback } from 'react'
import type {
  ColumnFiltersState,
  PaginationState,
  SortingState,
  VisibilityState,
} from '@tanstack/react-table'

export interface UseDataTableOptions {
  defaultPageSize?: number
  defaultPageIndex?: number
}

export interface UseDataTableReturn {
  // 状态
  columnFilters: ColumnFiltersState
  setColumnFilters: React.Dispatch<React.SetStateAction<ColumnFiltersState>>
  pagination: PaginationState
  setPagination: React.Dispatch<React.SetStateAction<PaginationState>>
  sorting: SortingState
  setSorting: React.Dispatch<React.SetStateAction<SortingState>>
  columnVisibility: VisibilityState
  setColumnVisibility: React.Dispatch<React.SetStateAction<VisibilityState>>
  
  // 计算属性
  pageNum: number
  pageSize: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  
  // 工具函数
  getFilterValue: (columnId: string) => string
  getSortingParams: () => { sortBy?: string; sortOrder?: 'asc' | 'desc' }
  resetFilters: () => void
}

/**
 * 统一的数据表格状态管理 Hook
 * 
 * 使用示例:
 * ```typescript
 * const {
 *   columnFilters,
 *   setColumnFilters,
 *   pagination,
 *   setPagination,
 *   sorting,
 *   setSorting,
 *   pageNum,
 *   getFilterValue,
 * } = useDataTable()
 * 
 * // 构建查询参数
 * const queryParams = useMemo(() => ({
 *   pageNum,
 *   pageSize: pagination.pageSize,
 *   username: getFilterValue('username'),
 * }), [pagination, columnFilters])
 * ```
 */
export function useDataTable(options: UseDataTableOptions = {}): UseDataTableReturn {
  const { defaultPageSize = 10, defaultPageIndex = 0 } = options

  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([])
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: defaultPageIndex,
    pageSize: defaultPageSize,
  })
  const [sorting, setSorting] = useState<SortingState>([])
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({})

  // 计算页码（API 使用 1-based）
  const pageNum = useMemo(() => pagination.pageIndex + 1, [pagination.pageIndex])
  const pageSize = pagination.pageSize
  const sortBy = sorting[0]?.id
  const sortOrder = sorting[0]?.desc ? 'desc' : sorting[0] ? 'asc' : undefined

  /**
   * 从 columnFilters 获取指定列的过滤值
   */
  const getFilterValue = useCallback((columnId: string): string => {
    const filter = columnFilters.find((f) => f.id === columnId)
    if (filter?.value && Array.isArray(filter.value)) {
      return filter.value[0] || ''
    }
    return (filter?.value as string) || ''
  }, [columnFilters])

  /**
   * 重置所有过滤条件
   */
  const resetFilters = useCallback(() => {
    setColumnFilters([])
    setSorting([])
    setPagination(prev => ({
      ...prev,
      pageIndex: 0,
    }))
  }, [])

  const getSortingParams = useCallback(() => {
    if (!sorting[0]?.id) {
      return {}
    }

    return {
      sortBy: sorting[0].id,
      sortOrder: sorting[0].desc ? 'desc' as const : 'asc' as const,
    }
  }, [sorting])

  return {
    // 状态
    columnFilters,
    setColumnFilters,
    pagination,
    setPagination,
    sorting,
    setSorting,
    columnVisibility,
    setColumnVisibility,
    
    // 计算属性
    pageNum,
    pageSize,
    sortBy,
    sortOrder,

    // 工具函数
    getFilterValue,
    getSortingParams,
    resetFilters,
  }
}
