import { Separator } from '@/components/ui/separator'

type ContentSectionProps = {
  title: string
  desc: string
  children: React.JSX.Element
}

export function ContentSection({ title, desc, children }: ContentSectionProps) {
  return (
    <div className='flex flex-1 flex-col'>
      <div className='flex-none rounded-3xl border border-border/60 bg-card/80 p-5 shadow-sm'>
        <p className='text-xs font-semibold uppercase tracking-[0.16em] text-muted-foreground'>Personal workspace</p>
        <h3 className='mt-2 text-2xl font-semibold tracking-tight'>{title}</h3>
        <p className='mt-1 text-sm text-muted-foreground'>{desc}</p>
      </div>
      <Separator className='my-5 flex-none' />
      <div className='faded-bottom h-full w-full overflow-y-auto scroll-smooth pe-4 pb-12'>
        <div className='-mx-1 px-1.5 lg:max-w-2xl'>{children}</div>
      </div>
    </div>
  )
}
