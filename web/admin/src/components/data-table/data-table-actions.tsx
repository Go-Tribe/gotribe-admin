import { memo, useCallback, useState } from 'react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Pencil1Icon, TrashIcon, EyeOpenIcon, DotsHorizontalIcon } from '@radix-ui/react-icons'
import { useI18n } from '@/context/i18n-provider'

export interface DataTableActionsProps {
  /** 编辑回调 */
  onEdit?: () => void
  /** 删除回调 */
  onDelete?: () => void
  /** 查看详情回调 */
  onView?: () => void
  /** 删除确认标题 */
  deleteConfirmTitle?: string
  /** 删除确认描述 */
  deleteConfirmDescription?: string
  /** 是否禁用编辑 */
  disabledEdit?: boolean
  /** 是否禁用删除 */
  disabledDelete?: boolean
  /** 是否禁用查看 */
  disabledView?: boolean
  /** 是否显示下拉菜单（否则直接显示按钮组） */
  useDropdown?: boolean
  /** 自定义类名 */
  className?: string
}

/**
 * 表格操作列组件
 * 
 * 统一表格行操作按钮的样式和行为，避免每个表格重复实现
 * 
 * 性能优化：使用 React.memo 避免不必要的重渲染
 * 
 * @example
 * ```tsx
 * // 基础用法（下拉菜单模式）
 * columns: [
 *   {
 *     id: 'actions',
 *     cell: ({ row }) => (
 *       <DataTableActions
 *         onEdit={() => handleEdit(row.original)}
 *         onDelete={() => handleDelete(row.original.id)}
 *         deleteConfirmTitle="确认删除此用户？"
 *       />
 *     ),
 *   },
 * ]
 * 
 * // 直接显示按钮组
 * <DataTableActions
 *   useDropdown={false}
 *   onEdit={handleEdit}
 *   onDelete={handleDelete}
 *   onView={handleView}
 * />
 * ```
 */
export const DataTableActions = memo(function DataTableActions({
  onEdit,
  onDelete,
  onView,
  deleteConfirmTitle,
  deleteConfirmDescription,
  disabledEdit = false,
  disabledDelete = false,
  disabledView = false,
  useDropdown = true,
  className,
}: DataTableActionsProps) {
  const { t } = useI18n()
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  const handleDeleteClick = useCallback(() => {
    if (deleteConfirmTitle) {
      setDeleteDialogOpen(true)
    } else {
      onDelete?.()
    }
  }, [deleteConfirmTitle, onDelete])

  const handleConfirmDelete = useCallback(() => {
    onDelete?.()
    setDeleteDialogOpen(false)
  }, [onDelete])

  const hasActions = onEdit || onDelete || onView

  if (!hasActions) {
    return null
  }

  // 下拉菜单模式（默认）
  if (useDropdown) {
    return (
      <>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">Open menu</span>
              <DotsHorizontalIcon className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            {onView && (
              <DropdownMenuItem
                onClick={onView}
                disabled={disabledView}
                className="cursor-pointer"
              >
                <EyeOpenIcon className="mr-2 h-4 w-4" />
                {t('common.view') || '查看'}
              </DropdownMenuItem>
            )}
            {onEdit && (
              <DropdownMenuItem
                onClick={onEdit}
                disabled={disabledEdit}
                className="cursor-pointer"
              >
                <Pencil1Icon className="mr-2 h-4 w-4" />
                {t('common.edit') || '编辑'}
              </DropdownMenuItem>
            )}
            {onDelete && (
              <DropdownMenuItem
                onClick={handleDeleteClick}
                disabled={disabledDelete}
                className="cursor-pointer text-destructive focus:text-destructive"
              >
                <TrashIcon className="mr-2 h-4 w-4" />
                {t('common.delete') || '删除'}
              </DropdownMenuItem>
            )}
          </DropdownMenuContent>
        </DropdownMenu>

        {/* 删除确认弹窗 */}
        <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{deleteConfirmTitle}</DialogTitle>
              {deleteConfirmDescription && (
                <DialogDescription>{deleteConfirmDescription}</DialogDescription>
              )}
            </DialogHeader>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => setDeleteDialogOpen(false)}
              >
                {t('common.cancel') || '取消'}
              </Button>
              <Button variant="destructive" onClick={handleConfirmDelete}>
                {t('common.confirm') || '确认'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </>
    )
  }

  // 直接显示按钮组模式
  return (
    <>
      <div className={cn('flex items-center gap-2', className)}>
        {onView && (
          <Button
            variant="outline"
            size="sm"
            className="h-8 border-border/60"
            onClick={onView}
            disabled={disabledView}
            title={t('common.view') || '查看'}
          >
            <EyeOpenIcon className="h-4 w-4" />
          </Button>
        )}
        {onEdit && (
          <Button
            variant="outline"
            size="sm"
            className="h-8 border-border/60"
            onClick={onEdit}
            disabled={disabledEdit}
            title={t('common.edit') || '编辑'}
          >
            <Pencil1Icon className="h-4 w-4" />
          </Button>
        )}
        {onDelete && (
          <Button
            variant="ghost"
            size="sm"
            className="h-8 text-destructive hover:text-destructive"
            onClick={handleDeleteClick}
            disabled={disabledDelete}
            title={t('common.delete') || '删除'}
          >
            <TrashIcon className="h-4 w-4" />
          </Button>
        )}
      </div>

      {/* 删除确认弹窗 */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{deleteConfirmTitle}</DialogTitle>
            {deleteConfirmDescription && (
              <DialogDescription>{deleteConfirmDescription}</DialogDescription>
            )}
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteDialogOpen(false)}
            >
              {t('common.cancel') || '取消'}
            </Button>
            <Button variant="destructive" onClick={handleConfirmDelete}>
              {t('common.confirm') || '确认'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
})

export default DataTableActions
