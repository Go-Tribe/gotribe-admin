import { memo, type ReactNode } from 'react'
import { SubTitle } from '@/components/sub-title'
import { DataTableCard } from '@/components/data-table/data-table-card'
import { cn } from '@/lib/utils'

export interface ListPageLayoutProps {
  /** 页面标题 */
  title: string
  /** 页面描述 */
  description?: string
  /** 标题栏右侧操作区（如新建按钮） */
  actions?: ReactNode
  /** 表格区域（通常包含 DataTable） */
  children: ReactNode
  /** 筛选区域内容（放在表格上方） */
  filterContent?: ReactNode
  /** 弹窗区域（FormDialog, ConfirmDialog 等） */
  dialogs?: ReactNode
  /** 自定义类名 */
  className?: string
  /** 卡片内边距 */
  cardPadding?: 'none' | 'sm' | 'md' | 'lg'
  /** 卡片外边距 */
  cardMargin?: 'none' | 'normal'
}

/**
 * 列表页面统一布局组件
 * 
 * 封装所有列表页的共同结构：
 * - 标题区（SubTitle）
 * - 卡片容器（DataTableCard）
 * - 筛选区
 * - 表格区
 * - 弹窗区
 * 
 * 使用 React.memo 优化性能，避免不必要的重渲染
 * 
 * @example
 * ```tsx
 * function AdminPage() {
 *   return (
 *     <ListPageLayout
 *       title={t('system.admin.title')}
 *       description={t('system.admin.description')}
 *       actions={<Button onClick={handleCreate}>新建</Button>}
 *       filterContent={<Input placeholder="搜索..." />}
 *       dialogs={<><FormDialog /><ConfirmDialog /></>}
 *     >
 *       <DataTable table={table} columns={columns} />
 *     </ListPageLayout>
 *   )
 * }
 * ```
 */
export const ListPageLayout = memo(function ListPageLayout({
  title,
  description,
  actions,
  children,
  filterContent,
  dialogs,
  className,
  cardPadding = 'md',
  cardMargin = 'normal',
}: ListPageLayoutProps) {
  return (
    <div className={cn('space-y-4', className)}>
      {/* 标题区 */}
      <div className="flex items-center justify-between px-4 pt-4">
        <SubTitle title={title} description={description}>
          {actions}
        </SubTitle>
      </div>

      {/* 表格卡片区 */}
      <DataTableCard padding={cardPadding} margin={cardMargin}>
        {filterContent && <div className="mb-4">{filterContent}</div>}
        {children}
      </DataTableCard>

      {/* 弹窗区 */}
      {dialogs}
    </div>
  )
})

export default ListPageLayout
