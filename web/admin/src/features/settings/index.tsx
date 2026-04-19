import { Outlet } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { UserCog } from 'lucide-react'
import { Separator } from '@/components/ui/separator'
// import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
// import { ProfileDropdown } from '@/components/profile-dropdown'
import { SidebarNav } from './components/sidebar-nav'
import { SubTitle } from '@/components/sub-title'

export function Settings() {
  const { t } = useTranslation()
  const sidebarNavItems = [
    {
      title: t('features.settings.sidebar.profile'),
      href: '/personal-center',
      icon: <UserCog size={18} />,
    },
  ]

  return (
    <>
      {/* ===== Top Heading ===== */}
      {/* <Header>
        <div className='ms-auto flex items-center space-x-4'>
          <ProfileDropdown />
        </div>
      </Header> */}

      <Main fixed>
        <SubTitle title={t('features.settings.title')} description={t('features.settings.description')} />
        <Separator className='my-4 lg:my-6' />
        <div className='flex flex-1 flex-col space-y-2 overflow-hidden md:space-y-2 lg:flex-row lg:space-y-0 lg:space-x-12'>
          <aside className='top-0 lg:sticky lg:w-1/5'>
            <SidebarNav items={sidebarNavItems} />
          </aside>
          <div className='flex w-full overflow-y-hidden p-1'>
            <Outlet />
          </div>
        </div>
      </Main>
    </>
  )
}
