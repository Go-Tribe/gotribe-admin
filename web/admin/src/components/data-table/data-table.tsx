import type { ReactNode } from 'react'
import type {
  ColumnDef,
  Table as TanStackTable,
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

export interface DataTableProps<TData> {
  table: TanStackTable<TData>
  columns: ColumnDef<TData>[]
  isLoading?: boolean
  error?: Error | null
  loadingText?: string
  errorText?: string
  emptyText?: string
  children?: ReactNode
  className?: string
}

/**
 * 通用数据表格组件
 * 
 * 使用示例:
 * ```tsx
 * <DataTable
 *   table={table}
 *   columns={columns}
 *   isLoading={isLoading}
 *   error={error}
 *   loadingText={t('loading')}
 *   errorText={t('loadError')}
 *   emptyText={t('noData')}
 * >
 *   <div className="flex gap-2 pb-4">
 *     // 过滤控件
 *   </div>
 * </DataTable>
 * ```
 */
export function DataTable<TData>({
  table,
  columns,
  isLoading = false,
  error = null,
  loadingText,
  errorText,
  emptyText,
  children,
  className,
}: DataTableProps<TData>) {
  const { t } = useI18n()

  const defaultLoadingText = t('components.dataTable.loading') || 'Loading...'
  const defaultErrorText = t('components.dataTable.error') || 'Failed to load data'
  const defaultEmptyText = t('components.dataTable.empty') || 'No data'

  return (
    <div className={className}>
      {children ? <div className='pb-4'>{children}</div> : null}

      <div className='overflow-x-auto'>
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => {
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
            {isLoading ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  {loadingText || defaultLoadingText}
                </TableCell>
              </TableRow>
            ) : error ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center text-destructive"
                >
                  {errorText || defaultErrorText}
                </TableCell>
              </TableRow>
            ) : table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && 'selected'}
                >
                  {row.getVisibleCells().map((cell) => {
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
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  {emptyText || defaultEmptyText}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <DataTablePagination table={table} className='border-t border-border/60 px-4 py-4' />
    </div>
  )
}
