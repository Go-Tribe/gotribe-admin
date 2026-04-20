import { useState, useEffect, type ComponentType, type ReactNode } from 'react'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { useI18n } from '@/context/i18n-provider'
import { cn } from '@/lib/utils'

export interface CodeSplitDialogProps {
  /** 是否打开 */
  open: boolean
  /** 打开状态变化回调 */
  onOpenChange: (open: boolean) => void
  /** 标题 */
  title: string
  /** 描述 */
  description?: string
  /** 懒加载的组件工厂函数 */
  contentComponent: () => Promise<{ default: ComponentType<Record<string, unknown>> }>
  /** 传递给内容组件的 props */
  contentProps?: Record<string, unknown>
  /** 加载占位 */
  fallback?: ReactNode
  /** 是否显示底部按钮 */
  showFooter?: boolean
  /** 确认按钮文本 */
  confirmText?: string
  /** 取消按钮文本 */
  cancelText?: string
  /** 确认回调 */
  onConfirm?: () => void | Promise<void>
  /** 是否加载中 */
  isLoading?: boolean
  /** 对话框最大宽度 */
  maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
  /** 自定义类名 */
  className?: string
}

const maxWidthMap = {
  sm: 'sm:max-w-sm',
  md: 'sm:max-w-md',
  lg: 'sm:max-w-lg',
  xl: 'sm:max-w-xl',
  full: 'sm:max-w-full',
}

/**
 * 代码分割对话框组件
 * 
 * 性能优化：弹窗内容懒加载，减少首屏 JS 体积
 * 适用于包含富文本编辑器等大型组件的弹窗
 * 
 * @example
 * ```tsx
 * // 懒加载 Slate 编辑器
 * <CodeSplitDialog
 *   open={open}
 *   onOpenChange={setOpen}
 *   title="编辑文章"
 *   contentComponent={() => import('./article-editor')}
 *   contentProps={{ articleId: 1 }}
 *   onConfirm={handleSave}
 * />
 * 
 * // 带自定义占位
 * <CodeSplitDialog
 *   open={open}
 *   onOpenChange={setOpen}
 *   title="编辑"
 *   contentComponent={() => import('./heavy-form')}
 *   fallback={<CustomSkeleton />}
 * />
 * ```
 */
export function CodeSplitDialog({
  open,
  onOpenChange,
  title,
  description,
  contentComponent,
  contentProps = {},
  fallback,
  showFooter = true,
  confirmText,
  cancelText,
  onConfirm,
  isLoading = false,
  maxWidth = 'lg',
  className,
}: CodeSplitDialogProps) {
  const { t } = useI18n()
  const [ContentComponent, setContentComponent] = useState<ComponentType<Record<string, unknown>> | null>(null)

  // 懒加载内容组件
  useEffect(() => {
    if (!open || ContentComponent) {
      return
    }

    let cancelled = false
    void contentComponent().then((module) => {
      if (!cancelled) {
        setContentComponent(() => module.default)
      }
    })

    return () => {
      cancelled = true
    }
  }, [open, ContentComponent, contentComponent])

  const defaultFallback = fallback || (
    <div className="space-y-4 py-4">
      <Skeleton className="h-8 w-full" />
      <Skeleton className="h-32 w-full" />
      <Skeleton className="h-8 w-full" />
    </div>
  )

  const displayConfirmText = confirmText || t('common.confirm') || '确认'
  const displayCancelText = cancelText || t('common.cancel') || '取消'

  return (
    <Dialog
      open={open}
      onOpenChange={(nextOpen) => {
        if (isLoading) return
        onOpenChange(nextOpen)
      }}
    >
      <DialogContent className={cn(maxWidthMap[maxWidth], className)}>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          {description && <DialogDescription>{description}</DialogDescription>}
        </DialogHeader>

        <div className="py-4">
          {ContentComponent ? (
            <ContentComponent {...contentProps} />
          ) : (
            defaultFallback
          )}
        </div>

        {showFooter && (
          <DialogFooter className="gap-2">
            <Button
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isLoading}
            >
              {displayCancelText}
            </Button>
            {onConfirm && (
              <Button onClick={onConfirm} disabled={isLoading}>
                {isLoading ? (
                  <span className="flex items-center gap-2">
                    <span className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                    {displayConfirmText}
                  </span>
                ) : (
                  displayConfirmText
                )}
              </Button>
            )}
          </DialogFooter>
        )}
      </DialogContent>
    </Dialog>
  )
}

/**
 * 简化版懒加载弹窗
 * 
 * 适用于只需要懒加载内容，不需要复杂配置的场景
 */
export interface LazyDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  component: () => Promise<{ default: ComponentType<Record<string, unknown>> }>
  componentProps?: Record<string, unknown>
  title?: string
  maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
  fallback?: ReactNode
}

export function LazyDialog({
  open,
  onOpenChange,
  component,
  componentProps = {},
  title,
  maxWidth = 'lg',
  fallback,
}: LazyDialogProps) {
  const [ContentComponent, setContentComponent] = useState<ComponentType<Record<string, unknown>> | null>(null)

  useEffect(() => {
    if (!open || ContentComponent) {
      return
    }

    let cancelled = false
    void component().then((module) => {
      if (!cancelled) {
        setContentComponent(() => module.default)
      }
    })

    return () => {
      cancelled = true
    }
  }, [open, ContentComponent, component])

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className={maxWidthMap[maxWidth]}>
        {title && (
          <DialogHeader>
            <DialogTitle>{title}</DialogTitle>
          </DialogHeader>
        )}
        {ContentComponent ? (
          <ContentComponent {...componentProps} />
        ) : (
          fallback || <DialogFallback />
        )}
      </DialogContent>
    </Dialog>
  )
}

function DialogFallback() {
  return (
    <div className="space-y-4 py-4">
      <Skeleton className="h-8 w-full" />
      <Skeleton className="h-32 w-full" />
      <Skeleton className="h-8 w-full" />
    </div>
  )
}

export default CodeSplitDialog
