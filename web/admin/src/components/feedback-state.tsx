import type { ReactNode } from 'react'
import { Loader2 } from 'lucide-react'
import { cn } from '@/lib/utils'

type FeedbackStateProps = {
  icon?: ReactNode
  title: ReactNode
  description?: ReactNode
  actions?: ReactNode
  className?: string
  tone?: 'default' | 'danger' | 'muted'
}

const toneClassMap: Record<NonNullable<FeedbackStateProps['tone']>, string> = {
  default: 'border-border/60 bg-muted/20',
  danger: 'border-destructive/20 bg-destructive/5',
  muted: 'border-border/40 bg-background/80',
}

export function FeedbackState({
  icon,
  title,
  description,
  actions,
  className,
  tone = 'default',
}: FeedbackStateProps) {
  return (
    <div
      className={cn(
        'flex flex-col items-center justify-center rounded-3xl border px-6 py-10 text-center shadow-sm',
        toneClassMap[tone],
        className
      )}
    >
      {icon ? (
        <div className='mb-4 flex h-14 w-14 items-center justify-center rounded-2xl border border-border/60 bg-background text-foreground shadow-sm'>
          {icon}
        </div>
      ) : null}
      <div className='space-y-1.5'>
        <h3 className='text-base font-semibold tracking-tight'>{title}</h3>
        {description ? (
          <p className='mx-auto max-w-md text-sm leading-6 text-muted-foreground'>
            {description}
          </p>
        ) : null}
      </div>
      {actions ? <div className='mt-5 flex flex-wrap items-center justify-center gap-3'>{actions}</div> : null}
    </div>
  )
}

export function LoadingState({
  title,
  description,
  className,
}: {
  title: ReactNode
  description?: ReactNode
  className?: string
}) {
  return (
    <FeedbackState
      className={className}
      title={title}
      description={description}
      icon={<Loader2 className='h-5 w-5 animate-spin text-primary' />}
      tone='muted'
    />
  )
}
