import { Toaster as Sonner, ToasterProps } from 'sonner'
import { useTheme } from '@/context/theme-provider'

export function Toaster({ ...props }: ToasterProps) {
  const { theme = 'system' } = useTheme()

  return (
    <Sonner
      theme={theme as ToasterProps['theme']}
      className='toaster group [&_div[data-content]]:w-full'
      toastOptions={{
        classNames: {
          toast:
            'rounded-2xl border border-border/60 bg-background/95 shadow-xl backdrop-blur-sm',
          title: 'text-sm font-semibold tracking-tight',
          description: 'text-sm text-muted-foreground',
          actionButton: 'rounded-xl',
          cancelButton: 'rounded-xl',
          success: 'border-emerald-500/20',
          error: 'border-destructive/20',
          warning: 'border-amber-500/20',
          info: 'border-primary/20',
        },
      }}
      style={
        {
          '--normal-bg': 'var(--popover)',
          '--normal-text': 'var(--popover-foreground)',
          '--normal-border': 'var(--border)',
        } as React.CSSProperties
      }
      {...props}
    />
  )
}
