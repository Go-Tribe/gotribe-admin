import { cn } from '@/lib/utils'

export interface DataTableCardProps {
  children: React.ReactNode
  className?: string
  padding?: 'none' | 'sm' | 'md' | 'lg'
  margin?: 'none' | 'normal'
}

/**
 * 表格卡片容器组件
 * 
 * 统一表格外层容器样式，消除重复代码
 * 
 * @example
 * ```tsx
 * <DataTableCard>
 *   <DataTable table={table} columns={columns} />
 * </DataTableCard>
 * ```
 */
export function DataTableCard({
  children,
  className,
  padding = 'md',
  margin = 'normal',
}: DataTableCardProps) {
  const paddingClasses = {
    none: '',
    sm: 'p-4',
    md: 'p-6',
    lg: 'p-8',
  }

  const marginClasses = {
    none: '',
    normal: 'mx-4',
  }

  return (
    <div
      className={cn(
        'rounded-md border bg-card',
        paddingClasses[padding],
        marginClasses[margin],
        className
      )}
    >
      {children}
    </div>
  )
}
