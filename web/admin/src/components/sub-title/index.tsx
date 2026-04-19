import { useSidebar } from '@/components/ui/sidebar'
import { Button } from '@/components/ui/button'
import { Menu } from 'lucide-react'
export function SubTitle({ 
  title, 
  description,
  children 
}: { 
  title: string
  description?: string
  children?: React.ReactNode 
}) {
  const { toggleSidebar, isMobile } = useSidebar()
  return (
    <div className='flex flex-col w-full'>
      {isMobile && (
        <Button
          data-sidebar='trigger'
          data-slot='sidebar-trigger'
          variant='ghost'
          size='icon'
          className='aspect-square size-8 shrink-0 max-md:scale-125 mb-2'
          onClick={(e) => {
            e.preventDefault()
            toggleSidebar()
          }}
        >
          {isMobile ? <Menu /> : ''}
        </Button>
      )}
      <div className='flex items-center justify-between w-full'>
        <div className='flex flex-col gap-2'>
          <h2 className='text-2xl font-bold'>{title}</h2>
          <p className='text-sm text-muted-foreground'>{description}</p>
        </div>
        {children && <div>{children}</div>}
      </div>
    </div>

  )
}
