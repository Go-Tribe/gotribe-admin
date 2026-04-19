import type { ReactNode } from 'react'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Form } from '@/components/ui/form'
import type { UseFormReturn, FieldValues } from 'react-hook-form'
import { useI18n } from '@/context/i18n-provider'
import { cn } from '@/lib/utils'

/** Section卡片配置 */
export interface FormSection {
  /** 分组标题 */
  title: string
  /** 分组描述 */
  description?: string
  /** 分组内容 */
  children: ReactNode
  /** 自定义类名 */
  className?: string
}

/**
 * Section卡片组件
 * 用于将表单字段分组展示，带有圆角边框和背景色
 */
export function FormSectionCard({
  title,
  description,
  children,
  className,
}: FormSection) {
  return (
    <section className={cn('py-2', className)}>
      <div className='mb-4 space-y-1'>
        <h3 className='text-sm font-semibold tracking-tight text-foreground'>{title}</h3>
        {description && <p className='text-sm text-muted-foreground'>{description}</p>}
      </div>
      <div className='space-y-4'>{children}</div>
    </section>
  )
}

/**
 * 信息摘要组件
 * 用于展示表单字段分组概览
 */
export function FormSummary({
  label,
  tags,
}: {
  label: string
  tags: string[]
}) {
  return (
    <div className='rounded-2xl border border-border/60 bg-card/80 p-4 shadow-sm'>
      <p className='text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground'>
        {label}
      </p>
      <div className='mt-2 flex flex-wrap gap-2 text-sm text-muted-foreground'>
        {tags.map((tag, index) => (
          <span key={index} className='rounded-full bg-muted px-3 py-1'>
            {tag}
          </span>
        ))}
      </div>
    </div>
  )
}

export interface FormDialogProps<TFormData extends FieldValues> {
  /** Dialog open state */
  open: boolean
  /** Callback when open state changes */
  onOpenChange: (open: boolean) => void
  /** Form instance from react-hook-form */
  form: UseFormReturn<TFormData>
  /** Dialog title */
  title: string
  /** Dialog description */
  description?: string
  /** Form content (fields) */
  children: ReactNode
  /** Submit handler */
  onSubmit: (data: TFormData) => void
  /** Loading state */
  isLoading?: boolean
  /** Submit button text */
  submitText?: string
  /** Cancel button text */
  cancelText?: string
  /** Whether the form is in edit mode */
  isEdit?: boolean
  /** Edit mode title (optional, defaults to title) */
  editTitle?: string
  /** Maximum width of the dialog */
  maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
}

const maxWidthMap = {
  sm: 'sm:max-w-sm',
  md: 'sm:max-w-md',
  lg: 'sm:max-w-lg',
  xl: 'sm:max-w-xl',
  full: 'sm:max-w-full',
}

/**
 * 通用表单对话框组件
 * 
 * 使用示例:
 * ```tsx
 * const form = useForm<FormData>({
 *   resolver: zodResolver(formSchema),
 *   defaultValues: {...},
 * })
 * 
 * <FormDialog
 *   open={dialogOpen}
 *   onOpenChange={setDialogOpen}
 *   form={form}
 *   title="Create User"
 *   editTitle="Edit User"
 *   isEdit={!!editingItem}
 *   onSubmit={handleSubmit}
 *   isLoading={createMutation.isPending || updateMutation.isPending}
 * >
 *   <FormField
 *     control={form.control}
 *     name="username"
 *     render={({ field }) => (
 *       <FormItem>
 *         <FormLabel>Username</FormLabel>
 *         <FormControl>
 *           <Input {...field} />
 *         </FormControl>
 *         <FormMessage />
 *       </FormItem>
 *     )}
 *   />
 * </FormDialog>
 * ```
 */
export function FormDialog<TFormData extends FieldValues>({
  open,
  onOpenChange,
  form,
  title,
  description,
  children,
  onSubmit,
  isLoading = false,
  submitText,
  cancelText,
  isEdit = false,
  editTitle,
  maxWidth = 'md',
}: FormDialogProps<TFormData>) {
  const { t } = useI18n()

  const displayTitle = isEdit && editTitle ? editTitle : title
  const displaySubmitText = submitText || (isEdit 
    ? (t('components.confirmDialog.save') || 'Save')
    : (t('components.confirmDialog.create') || 'Create')
  )
  const displayCancelText = cancelText || t('components.confirmDialog.cancel') || 'Cancel'

  return (
    <Dialog
      open={open}
      onOpenChange={(nextOpen) => {
        if (isLoading) return
        onOpenChange(nextOpen)
      }}
    >
      <DialogContent className={cn(maxWidthMap[maxWidth], 'max-h-[90vh] flex flex-col')}>
        <DialogHeader className='shrink-0'>
          <DialogTitle>{displayTitle}</DialogTitle>
          {description && <DialogDescription>{description}</DialogDescription>}
        </DialogHeader>
        
        <Form {...form}>
          <form 
            onSubmit={form.handleSubmit(onSubmit)} 
            className='flex-1 overflow-y-auto pr-2 space-y-5 min-h-0'
          >
            {children}
          </form>
        </Form>
        
        <DialogFooter className='shrink-0 gap-2 pt-4 border-t mt-4'>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isLoading}
          >
            {displayCancelText}
          </Button>
          <Button type="submit" disabled={isLoading} onClick={form.handleSubmit(onSubmit)}>
            {isLoading ? (
              <span className="flex items-center gap-2">
                <span className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                {isEdit
                  ? t('components.confirmDialog.updateLoading') || displaySubmitText
                  : t('components.confirmDialog.createLoading') || displaySubmitText}
              </span>
            ) : (
              displaySubmitText
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
