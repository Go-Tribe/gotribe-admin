import { useQuery } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import { Menu, X } from 'lucide-react'
import { SidebarMenu, SidebarMenuItem, useSidebar } from '@/components/ui/sidebar'
import { Button } from '@/components/ui/button'
import { getConfig } from '@/features/system/service/config'
import { Logo } from '@/assets/logo'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'

export function SidebarBrand() {
  const { setOpenMobile, toggleSidebar, state, isMobile } = useSidebar()
  const isVertical = state === 'collapsed'
  const { t } = useTranslation()
  const { data: configData } = useQuery({
    queryKey: ['systemConfig'],
    queryFn: () => getConfig(),
    staleTime: 5 * 60 * 1000,
  })
  const systemConfig = configData?.systemConfig
  const title = systemConfig?.title?.trim() || t('components.layout.appTitle.title')
  const logoUrl = systemConfig?.logo?.trim()

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <div
          className={cn(
            'flex w-full gap-2 py-2',
            isVertical && !isMobile ? 'flex-col-reverse' : 'flex-row',
            isMobile ? 'items-start' : 'item-center'
          )}
        >
          <Link
            to='/'
            onClick={() => setOpenMobile(false)}
            className={cn('flex w-full shrink-0 items-center flex-1 gap-2',
            )}
          >
            <div
              className={cn(
                'flex aspect-square w-full items-center justify-center overflow-hidden rounded-md',
                isMobile ? 'size-12' : 'size-8'
              )}
            >
              {logoUrl ? (
                <img src={logoUrl} alt='' className='size-full object-contain' />
              ) : (
                <Logo className='size-4 shrink-0' />
              )}
            </div>
            {!isVertical && (
              <span className='truncate font-semibold text-sm'>{title}</span>
            )}
          </Link>
          <Button
            data-sidebar='trigger'
            data-slot='sidebar-trigger'
            variant='ghost'
            size='icon'
            className='aspect-square size-8 shrink-0 max-md:scale-125'
            onClick={(e) => {
              e.preventDefault()
              toggleSidebar()
            }}
          >
            <X className='md:hidden' />
            <Menu className='max-md:hidden' />
            <span className='sr-only'>{t('components.layout.appTitle.toggleSidebar')}</span>
          </Button>
        </div>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
