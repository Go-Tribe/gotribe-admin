import { AlertTriangle, Loader2 } from 'lucide-react'
import { cn } from '@/lib/utils'
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { useTranslation } from 'react-i18next'

type ConfirmDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  title: React.ReactNode
  disabled?: boolean
  desc: React.JSX.Element | string
  cancelBtnText?: string
  confirmText?: React.ReactNode
  destructive?: boolean
  handleConfirm: () => void
  isLoading?: boolean
  className?: string
  children?: React.ReactNode
}

export function ConfirmDialog(props: ConfirmDialogProps) {
  const {
    title,
    desc,
    children,
    className,
    confirmText,
    cancelBtnText,
    destructive,
    isLoading,
    disabled = false,
    handleConfirm,
    onOpenChange,
    ...actions
  } = props
  const { t } = useTranslation()
  return (
    <AlertDialog
      {...actions}
      open={actions.open}
      onOpenChange={(nextOpen) => {
        if (isLoading) return
        onOpenChange(nextOpen)
      }}
    >
      <AlertDialogContent className={cn(className && className)}>
        <AlertDialogHeader className='text-start'>
          <div className='mb-1 flex h-12 w-12 items-center justify-center rounded-2xl border border-destructive/20 bg-destructive/10 text-destructive shadow-sm'>
            <AlertTriangle className='h-5 w-5' />
          </div>
          <AlertDialogTitle>{title}</AlertDialogTitle>
          <AlertDialogDescription asChild>
            <div>{desc}</div>
          </AlertDialogDescription>
          {destructive ? (
            <div className='mt-3 rounded-2xl border border-destructive/15 bg-destructive/5 px-3 py-2 text-sm text-muted-foreground'>
              {t('components.confirmDialog.irreversibleHint')}
            </div>
          ) : null}
        </AlertDialogHeader>
        {children}
        <AlertDialogFooter>
          <AlertDialogCancel disabled={isLoading}>
            {cancelBtnText ?? t('components.confirmDialog.cancel')}
          </AlertDialogCancel>
          <Button
            variant={destructive ? 'destructive' : 'default'}
            onClick={handleConfirm}
            disabled={disabled || isLoading}
            className='min-w-24'
          >
            {isLoading ? (
              <span className='flex items-center gap-2'>
                <Loader2 className='h-4 w-4 animate-spin' />
                {t('components.confirmDialog.processing')}
              </span>
            ) : (
              confirmText ?? t('components.confirmDialog.continue')
            )}
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
