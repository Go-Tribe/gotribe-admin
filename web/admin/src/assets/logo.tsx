import { cn } from '@/lib/utils'

export function Logo({ className, ...props }: { className?: string, props?: React.HTMLAttributes<HTMLImageElement> }) {
  return (
    <img className={cn('size-6', className)} src="https://avatars.githubusercontent.com/u/106083123?s=400&u=985e6fbd108c676fe8a72743d675a047ae4785ba&v=4" alt="Go-Gribe Logo" {...props} />
  )
}
