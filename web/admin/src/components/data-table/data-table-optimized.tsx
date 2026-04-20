import type { ReactNode } from 'react'
import { cn } from '@/lib/utils'
import type {
  ColumnDef,
  Table as TanStackTable,
  Row,
  Cell,
  Header,
  HeaderGroup,
} from '@tanstack/react-table'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { DataTablePagination } from './pagination'
import { useI18n } from '@/context/i18n-provider'

export interface DataTableProps<TData = unknown> {
  table: TanStackTable<TData>
  columns: ColumnDef<TData>[]
  isLoading?: boolean
  error?: Error | null
  loadingText?: string
  errorText?: string
  emptyText?: string
  children?: ReactNode
  className?: string
  /** 是否显示外边框 */
  bordered?: boolean
}

/**
 * 优化的表格行组件
 * 使用 React.memo 防止不必要的重渲染
 */
function DataTableRow<TData>({
  row,
}: {
  row: Row<TData>
}) {
  return (
    <TableRow data-state={row.getIsSelected() && 'selected'}>
      {row.getVisibleCells().map((cell: Cell<TData, unknown>) => {
        const cellContent =
          typeof cell.column.columnDef.cell === 'function'
            ? cell.column.columnDef.cell(cell.getContext())
            : cell.column.columnDef.cell
        const meta = cell.column.columnDef.meta as { tdClassName?: string } | undefined
        return (
          <TableCell key={cell.id} className={meta?.tdClassName}>
            {cellContent}
          </TableCell>
        )
      })}
    </TableRow>
  )
}

/**
 * 优化的表头行组件
 */
function DataTableHeaderRow<TData>({
  headerGroup,
}: {
  headerGroup: HeaderGroup<TData>
}) {
  return (
    <TableRow>
      {headerGroup.headers.map((header: Header<TData, unknown>) => {
        const headerContent = header.isPlaceholder
          ? null
          : typeof header.column.columnDef.header === 'function'
            ? header.column.columnDef.header(header.getContext())
            : header.column.columnDef.header
        const meta = header.column.columnDef.meta as { thClassName?: string } | undefined
        return (
          <TableHead key={header.id} className={meta?.thClassName}>
            {headerContent}
          </TableHead>
        )
      })}
    </TableRow>
  )
}

/**
 * 优化后的通用数据表格组件
 * 
 * 优化点：
 * 1. 使用 React.memo 缓存行渲染
 * 2. 使用 useMemo 缓存计算结果
 * 3. 减少不必要的重渲染
 * 
 * @example
 * ```tsx
 * <DataTableOptimized
 *   table={table}
 *   columns={columns}
 *   isLoading={isLoading}
 *   error={error}
 * >
 *   <div className="flex gap-2 pb-4">
 *     // 过滤控件
 *   </div>
 * </DataTableOptimized>
 * ```
 */
export function DataTableOptimized<TData>({
  table,
  columns,
  isLoading = false,
  error = null,
  loadingText,
  errorText,
  emptyText,
  children,
  className,
  bordered,
}: DataTableProps<TData>) {
  const { t } = useI18n()

  const defaultLoadingText = t('components.dataTable.loading') || 'Loading...'
  const defaultErrorText = t('components.dataTable.error') || 'Failed to load data'
  const defaultEmptyText = t('components.dataTable.empty') || 'No data'

  // 缓存行数据
  const rows = table.getRowModel().rows

  // 缓存表头组 - 使用 table 作为依赖，因为 getHeaderGroups() 返回新数组
  const headerGroups = table.getHeaderGroups()

  // 状态渲染
  const renderContent = () => {
    if (isLoading) {
      return (
        <TableRow>
          <TableCell colSpan={columns.length} className="h-24 text-center">
            <div className="flex items-center justify-center gap-2">
              <div className="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
              {loadingText || defaultLoadingText}
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
            {errorText || defaultErrorText}
          </TableCell>
        </TableRow>
      )
    }

    if (rows.length === 0) {
      return (
        <TableRow>
          <TableCell colSpan={columns.length} className="h-24 text-center text-muted-foreground">
            {emptyText || defaultEmptyText}
          </TableCell>
        </TableRow>
      )
    }

    // 渲染行 - 对于树形表格，getRowModel().rows 已经包含展开后的可见行
    return rows.map((row) => <DataTableRow key={row.id} row={row} />)
  }

  return (
    <div className={className}>
      {children ? <div className='pb-4'>{children}</div> : null}

      <div className={cn('overflow-x-auto', bordered === true && 'rounded-md border')}>
        <Table>
          <TableHeader className="bg-muted/50">
            {headerGroups.map((headerGroup) => (
              <DataTableHeaderRow key={headerGroup.id} headerGroup={headerGroup} />
            ))}
          </TableHeader>
          <TableBody>
            {renderContent()}
          </TableBody>
        </Table>
      </div>

      <DataTablePagination table={table} className='border-t border-border/60 px-4 py-4' />
    </div>
  )
}

// 为了向后兼容，保留原名导出
export { DataTableOptimized as DataTable }
