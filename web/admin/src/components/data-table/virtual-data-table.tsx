import { memo, useRef, useCallback, useState } from 'react'
import type { CellContext, ColumnDef } from '@tanstack/react-table'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

import { useI18n } from '@/context/i18n-provider'
import { cn } from '@/lib/utils'

export interface VirtualDataTableProps<TData> {
  /** 表格数据 */
  data: TData[]
  /** 列定义 */
  columns: ColumnDef<TData>[]
  /** 是否加载中 */
  isLoading?: boolean
  /** 错误对象 */
  error?: Error | null
  /** 行高度（像素） */
  rowHeight?: number
  /** 可视区域高度（像素） */
  containerHeight?: number
  /** 上下缓冲区行数 */
  overscan?: number
  /** 加载文本 */
  loadingText?: string
  /** 错误文本 */
  errorText?: string
  /** 空数据文本 */
  emptyText?: string
  /** 自定义类名 */
  className?: string
  /** 是否启用分页 */
  enablePagination?: boolean
  /** 每页条数 */
  pageSize?: number
  /** 当前页 */
  currentPage?: number
  /** 总条数 */
  total?: number
  /** 页码变化回调 */
  onPageChange?: (page: number) => void
  /** 每页条数变化回调 */
  onPageSizeChange?: (pageSize: number) => void
  /** 获取行唯一标识，用于稳定 key。默认尝试取 row.id / row._id */
  getRowId?: (row: TData) => string
}

/**
 * 虚拟滚动数据表格组件
 * 
 * 性能优化：大数据量时使用虚拟滚动，只渲染可视区域数据
 * 适用于 1000+ 行数据的场景
 * 
 * @example
 * ```tsx
 * <VirtualDataTable
 *   data={largeData} // 10000 条数据
 *   columns={columns}
 *   rowHeight={48}
 *   containerHeight={600}
 *   overscan={5}
 * />
 * ```
 */
export const VirtualDataTable = memo(function VirtualDataTable<TData>({
  data,
  columns,
  isLoading = false,
  error = null,
  rowHeight = 48,
  containerHeight = 600,
  overscan = 5,
  loadingText,
  errorText,
  emptyText,
  className,
  enablePagination = false,
  pageSize = 10,
  currentPage = 1,
  total,
  onPageChange,
  onPageSizeChange,
  getRowId,
}: VirtualDataTableProps<TData>) {
  const { t } = useI18n()
  const containerRef = useRef<HTMLDivElement>(null)
  const rafRef = useRef<number | null>(null)
  const [scrollTop, setScrollTop] = useState(0)

  // 计算可视区域
  const totalHeight = data.length * rowHeight
  const startIndex = Math.max(0, Math.floor(scrollTop / rowHeight) - overscan)
  const visibleCount = Math.ceil(containerHeight / rowHeight) + overscan * 2
  const endIndex = Math.min(data.length, startIndex + visibleCount)
  const visibleData = data.slice(startIndex, endIndex)
  const offsetY = startIndex * rowHeight

  // 滚动事件处理（使用 requestAnimationFrame 节流，避免 60fps+ 触发 setState）
  const handleScroll = useCallback((e: React.UIEvent<HTMLDivElement>) => {
    const newScrollTop = e.currentTarget.scrollTop
    if (rafRef.current) {
      cancelAnimationFrame(rafRef.current)
    }
    rafRef.current = requestAnimationFrame(() => {
      setScrollTop(newScrollTop)
    })
  }, [])

  // 默认文本
  const defaultLoadingText = loadingText || t('components.dataTable.loading') || 'Loading...'
  const defaultErrorText = errorText || t('components.dataTable.error') || 'Failed to load data'
  const defaultEmptyText = emptyText || t('components.dataTable.empty') || 'No data'

  // 状态渲染
  const renderContent = () => {
    if (isLoading) {
      return (
        <TableRow>
          <TableCell colSpan={columns.length} className="h-24 text-center">
            <div className="flex items-center justify-center gap-2">
              <div className="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
              {defaultLoadingText}
            </div>
          </TableCell>
        </TableRow>
      )
    }

    if (error) {
      return (
        <TableRow>
          <TableCell
            colSpan={columns.length}
            className="h-24 text-center text-destructive"
          >
            {defaultErrorText}
          </TableCell>
        </TableRow>
      )
    }

    if (data.length === 0) {
      return (
        <TableRow>
          <TableCell
            colSpan={columns.length}
            className="h-24 text-center text-muted-foreground"
          >
            {defaultEmptyText}
          </TableCell>
        </TableRow>
      )
    }

    // 虚拟渲染行
    return (
      <>
        {/* 占位撑开高度 */}
        <tr style={{ height: offsetY }} />
        {visibleData.map((row, index) => {
          const actualIndex = startIndex + index
          const rowKey = getRowId
            ? getRowId(row)
            : (row as Record<string, unknown>).id?.toString() ??
              (row as Record<string, unknown>)._id?.toString() ??
              `row-${actualIndex}`
          return (
            <TableRow
              key={rowKey}
              style={{ height: rowHeight }}
              data-index={actualIndex}
            >
              {columns.map((column, colIndex) => {
                const rowRecord = row as Record<string, unknown>
                const cellValue = column.id ? rowRecord[column.id] : undefined
                const cellContent = column.cell
                  ? typeof column.cell === 'function'
                    ? column.cell({
                        row: { original: row },
                        getValue: () => cellValue,
                      } as unknown as CellContext<TData, unknown>)
                    : column.cell
                  : cellValue
                const meta = column.meta as { tdClassName?: string } | undefined
                return (
                  <TableCell key={colIndex} className={cn(meta?.tdClassName, 'truncate')}>
                    {cellContent}
                  </TableCell>
                )
              })}
            </TableRow>
          )
        })}
        {/* 底部占位 */}
        <tr style={{ height: totalHeight - offsetY - visibleData.length * rowHeight }} />
      </>
    )
  }

  return (
    <div className={className}>
      <div
        ref={containerRef}
        onScroll={handleScroll}
        className="overflow-auto rounded-md border"
        style={{ height: containerHeight }}
      >
        <Table>
          <TableHeader className="sticky top-0 z-10 bg-muted/50">
            <TableRow style={{ height: rowHeight }}>
              {columns.map((column, index) => {
                const meta = column.meta as { thClassName?: string } | undefined
                return (
                  <TableHead key={index} className={meta?.thClassName}>
                    {column.header as React.ReactNode}
                  </TableHead>
                )
              })}
            </TableRow>
          </TableHeader>
          <TableBody>
            {renderContent()}
          </TableBody>
        </Table>
      </div>

      {enablePagination && (
        <div className="flex items-center justify-between border-t border-border/60 px-4 py-4">
          <div className="text-sm text-muted-foreground">
            {total !== undefined && `共 ${total} 条`}
          </div>
          <div className="flex items-center gap-2">
            <select
              value={pageSize}
              onChange={(e) => onPageSizeChange?.(Number(e.target.value))}
              className="h-8 rounded-md border bg-background px-2 text-sm"
            >
              {[10, 20, 50, 100].map((size) => (
                <option key={size} value={size}>
                  {size} 条/页
                </option>
              ))}
            </select>
            <div className="flex gap-1">
              <button
                onClick={() => onPageChange?.(currentPage - 1)}
                disabled={currentPage <= 1}
                className="h-8 w-8 rounded-md border text-sm disabled:opacity-50"
              >
                ←
              </button>
              <span className="flex h-8 items-center px-3 text-sm">
                第 {currentPage} 页
              </span>
              <button
                onClick={() => onPageChange?.(currentPage + 1)}
                disabled={total !== undefined && currentPage * pageSize >= total}
                className="h-8 w-8 rounded-md border text-sm disabled:opacity-50"
              >
                →
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
})

export default VirtualDataTable
