import type { ReactNode } from 'react'
import type {
  ColumnDef,
  Table as TanStackTable,
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
import { cn } from '@/lib/utils'

export interface TreeDataTableProps<TData> {
  table: TanStackTable<TData>
  columns: ColumnDef<TData>[]
  isLoading?: boolean
  error?: Error | null
  loadingText?: string
  errorText?: string
  emptyText?: string
  children?: ReactNode
  className?: string
  bordered?: boolean
}

/**
 * 树形数据表格组件
 * 
 * 特点：
 * - 支持展开/折叠行
 * - 不使用 memo 优化，确保展开状态变化时正确渲染
 * - 适用于菜单、分类等树形结构数据
 * 
 * @example
 * ```tsx
 * <TreeDataTable
 *   table={table}
 *   columns={columns}
 *   isLoading={isLoading}
 *   error={error}
 * >
 *   <div>筛选区</div>
 * </TreeDataTable>
 * ```
 */
export function TreeDataTable<TData>({
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
}: TreeDataTableProps<TData>) {
  const { t } = useI18n()

  const defaultLoadingText = loadingText || t('components.dataTable.loading') || 'Loading...'
  const defaultErrorText = errorText || t('components.dataTable.error') || 'Failed to load data'
  const defaultEmptyText = emptyText || t('components.dataTable.empty') || 'No data'

  // 获取可见行（包含展开后的子行）
  const rows = table.getRowModel().rows
  const headerGroups = table.getHeaderGroups()

  // 渲染行 - getRowModel().rows 已经包含展开后的所有可见行
  const renderRows = () => {
    return rows.map((row) => (
      <TableRow key={row.id} data-state={row.getIsSelected() && 'selected'}>
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
    ))
  }

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

    if (rows.length === 0) {
      return (
        <TableRow>
          <TableCell colSpan={columns.length} className="h-24 text-center text-muted-foreground">
            {defaultEmptyText}
          </TableCell>
        </TableRow>
      )
    }

    return renderRows()
  }

  return (
    <div className={className}>
      {children ? <div className='pb-4'>{children}</div> : null}

      <div className={cn('overflow-x-auto', bordered === true && 'rounded-md border')}>
        <Table>
          <TableHeader className="bg-muted/50">
            {headerGroups.map((headerGroup: HeaderGroup<TData>) => (
              <TableRow key={headerGroup.id}>
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

export default TreeDataTable
