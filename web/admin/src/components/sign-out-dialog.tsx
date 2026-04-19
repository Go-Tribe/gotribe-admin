import { useNavigate, useLocation } from '@tanstack/react-router'
import { useLogout } from '@/stores/auth-store'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { useTranslation } from 'react-i18next'
import { outLogin } from '@/service/user'

interface SignOutDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function SignOutDialog({ open, onOpenChange }: SignOutDialogProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const logout = useLogout()
  const { t } = useTranslation()

  const handleSignOut = async () => {
    try {
      await outLogin()
    } catch {
      // Ignore logout API errors and still clear local auth state.
    } finally {
      logout()
      // Preserve current location for redirect after sign-in
      const currentPath = location.href
      navigate({
        to: '/sign-in',
        search: { redirect: currentPath },
        replace: true,
      })
    }
  }

  return (
    <ConfirmDialog
      open={open}
      onOpenChange={onOpenChange}
      title={t('components.layout.signOutDialog.title')}
      desc={t('components.layout.signOutDialog.description')}
      confirmText={t('components.layout.signOutDialog.confirmText')}
      destructive
      handleConfirm={handleSignOut}
      className='sm:max-w-sm'
    />
  )
}
