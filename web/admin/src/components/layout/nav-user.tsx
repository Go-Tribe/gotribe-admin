import { useState } from 'react'
import { Link } from '@tanstack/react-router'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  BadgeCheck,
  ChevronsUpDown,
  Languages,
  LogOut,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import useDialogState from '@/hooks/use-dialog-state'
import { useI18n } from '@/context/i18n-provider'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { ResourceUpload } from '@/components/resource-upload'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar'
import { SignOutDialog } from '@/components/sign-out-dialog'
import { ThemeSwitch } from '@/components/theme-switch'
import { ConfigDrawer } from '@/components/config-drawer'
import { type AuthUser, useSetAuthUser } from '@/stores/auth-store'
import { getAdminInfo, updateAdmin } from '@/features/system/service'

type NavUserProps = AuthUser | null

export function NavUser({ user }: { user: NavUserProps }) {
  const { isMobile } = useSidebar()
  const [open, setOpen] = useDialogState()
  const [avatarDialogOpen, setAvatarDialogOpen] = useState(false)
  const { language, setLanguage } = useI18n()
  const { t } = useTranslation()
  const setAuthUser = useSetAuthUser()
  const queryClient = useQueryClient()

  const { mutate: updateAvatar } = useMutation({
    mutationFn: async (avatar: string) => {
      const admin = await getAdminInfo()
      if (!admin) throw new Error(t('features.settings.profile.noData'))
      const status = admin.status === 1 || admin.status === 2 ? admin.status : 1
      return updateAdmin({ ...admin, avatar, status })
    },
    onSuccess: (_, newAvatar) => {
      if (user) setAuthUser({ ...user, avatar: newAvatar })
      queryClient.invalidateQueries({ queryKey: ['adminInfo'] })
      setAvatarDialogOpen(false)
      toast.success(t('components.layout.navUser.avatarUpdateSuccess'))
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : '')
    },
  })

  const onAvatarClick = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (user?.id != null) setAvatarDialogOpen(true)
  }

  return (
    <>
      <SidebarMenu>
        <SidebarMenuItem>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <SidebarMenuButton
                size='lg'
                asChild
                className='data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground'
              >
                <div>
                  <button
                    type='button'
                    className='rounded-lg outline-none ring-sidebar-ring focus-visible:ring-2'
                    onClick={onAvatarClick}
                    aria-label={t('components.layout.navUser.changeAvatar')}
                  >
                    <Avatar className='h-8 w-8 rounded-lg'>
                      <AvatarImage src={user?.avatar as string} alt={(user?.nickname || user?.username) as string} />
                      <AvatarFallback className='rounded-lg'>SN</AvatarFallback>
                    </Avatar>
                  </button>
                  <div className='grid flex-1 text-start text-sm leading-tight'>
                    <span className='truncate font-semibold'>{(user?.nickname || user?.username) as string}</span>
                    <span className='truncate text-xs'>{user?.email as string}</span>
                  </div>
                  <ChevronsUpDown className='ms-auto size-4' />
                </div>
              </SidebarMenuButton>
            </DropdownMenuTrigger>
            <DropdownMenuContent
              className='w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg'
              side={isMobile ? 'bottom' : 'right'}
              align='end'
              sideOffset={4}
            >
              <DropdownMenuLabel className='p-0 font-normal'>
                <div className='flex items-center gap-2 px-1 py-1.5 text-start text-sm'>
                  <Avatar className='h-8 w-8 rounded-full'>
                    <AvatarImage src={user?.avatar as string} alt={(user?.nickname || user?.username) as string} />
                    <AvatarFallback className='rounded-lg'>{user?.username as string}</AvatarFallback>
                  </Avatar>
                  <div className='grid flex-1 text-start text-sm leading-tight'>
                    <span className='truncate font-semibold'>{(user?.nickname || user?.username) as string}</span>
                    <span className='truncate text-xs'>{user?.email as string}</span>
                  </div>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuGroup>
                <DropdownMenuItem asChild>
                  <Link to='/personal-center'>
                    <BadgeCheck />
                    {t('components.layout.navUser.personalCenter')}
                  </Link>
                </DropdownMenuItem>
              </DropdownMenuGroup>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => setLanguage(language === 'zh' ? 'en' : 'zh')}
              >
                <Languages />
                {language === 'zh' ? t('components.layout.navUser.english') : t('components.layout.navUser.chinese')}
              </DropdownMenuItem>
              <DropdownMenuItem
                onSelect={(e) => e.preventDefault()}
                className='relative flex gap-2'
              >
                <ConfigDrawer asMenuItem />
                {t('components.layout.navUser.config')}
              </DropdownMenuItem>
              <ThemeSwitch asMenuItem label={t('components.layout.navUser.theme')} />
              <DropdownMenuSeparator />
              <DropdownMenuItem
                variant='destructive'
                onClick={() => setOpen(true)}
              >
                <LogOut />
                {t('components.layout.navUser.signOut')}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </SidebarMenuItem>
      </SidebarMenu>

      <SignOutDialog open={!!open} onOpenChange={setOpen} />
      <ResourceUpload
        open={avatarDialogOpen}
        onOpenChange={setAvatarDialogOpen}
        onSelect={(resource) => updateAvatar(resource.url)}
        type={1}
        title={t('components.layout.navUser.changeAvatar')}
        description={t('components.layout.navUser.changeAvatarDesc')}
      />
    </>
  )
}
