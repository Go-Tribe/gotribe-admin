import { useMemo, useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { useTranslation } from 'react-i18next'
import { Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { FeedbackState, LoadingState } from '@/components/feedback-state'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { PasswordInput } from '@/components/password-input'
import { ResourceUpload } from '@/components/resource-upload'
import { Separator } from '@/components/ui/separator'
import { useAuthUser, useSetAuthUser, useAuthStore } from '@/stores/auth-store'
import type { Admin } from '@/shared/types'
import { getAdminInfo, changePassword, updateAdmin } from '@/shared/api'

type ChangePasswordValues = {
  old_password: string
  new_password: string
  confirm_password: string
}

/** 从 admin 或 auth.user 得到规范的 number[]，保证更新时 roleIds 传值正确 */
function normalizeRoleIds(
  admin: Admin & { roles?: { id?: number }[] },
  authUser: { roleIds?: number[] } | null
): number[] {
  if (Array.isArray(admin.roleIds)) {
    const ids = admin.roleIds.map((id) => Number(id)).filter((n) => !Number.isNaN(n))
    if (ids.length > 0) return ids
  }
  if (Array.isArray(admin.roles)) {
    const ids = admin.roles.map((r) => Number(r?.id)).filter((n) => !Number.isNaN(n))
    if (ids.length > 0) return ids
  }
  if (authUser?.roleIds && Array.isArray(authUser.roleIds)) {
    const ids = authUser.roleIds.map((id) => Number(id)).filter((n) => !Number.isNaN(n))
    if (ids.length > 0) return ids
  }
  return []
}

export function ProfileForm() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const user = useAuthUser()
  const setAuthUser = useSetAuthUser()
  const resetAuth = useAuthStore(state => state.auth.reset)
  const [avatarDialogOpen, setAvatarDialogOpen] = useState(false)

  const { mutate: updateAvatar } = useMutation({
    mutationFn: async (payload: { admin: Admin & { roles?: { id?: number }[] }; avatar: string }) => {
      const { admin, avatar } = payload
      const status = admin.status === 1 || admin.status === 2 ? admin.status : 1
      const roleIds = normalizeRoleIds(admin, user)
      return updateAdmin({ ...admin, avatar, status, roleIds })
    },
    onSuccess: (_, { avatar: newAvatar }) => {
      if (user) setAuthUser({ ...user, avatar: newAvatar })
      queryClient.invalidateQueries({ queryKey: ['adminInfo'] })
      setAvatarDialogOpen(false)
      toast.success(t('features.settings.profile.avatarUpdateSuccess'))
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : t('features.settings.profile.avatarUpdateError'))
    },
  })

  const changePasswordSchema = useMemo(
    () =>
      z
        .object({
          old_password: z.string().min(1, t('features.settings.validation.oldPasswordRequired')),
          new_password: z
            .string()
            .min(6, t('features.settings.validation.newPasswordMin'))
            .max(20, t('features.settings.validation.newPasswordMax')),
          confirm_password: z.string().min(1, t('features.settings.validation.confirmPasswordRequired')),
        })
        .refine((data) => data.new_password === data.confirm_password, {
          message: t('features.settings.validation.passwordMismatch'),
          path: ['confirm_password'],
        }),
    [t]
  )

  const { data: admin, isLoading, error } = useQuery({
    queryKey: ['adminInfo'],
    queryFn: () => getAdminInfo({ page_num: 1, page_size: 3 }),
  })

  const passwordForm = useForm<ChangePasswordValues>({
    resolver: zodResolver(changePasswordSchema),
    defaultValues: {
      old_password: '',
      new_password: '',
      confirm_password: '',
    },
  })

  const { mutate: submitChangePassword, isPending: isChangingPassword } = useMutation({
    mutationFn: (values: ChangePasswordValues) =>
      changePassword({
        OldPassword: values.old_password,
        NewPassword: values.new_password,
      }),
    onSuccess: () => {
      toast.success(t('features.settings.changePassword.success'))
      passwordForm.reset()
      queryClient.invalidateQueries({ queryKey: ['adminInfo'] })
      resetAuth()
      navigate({ to: '/sign-in', search: { redirect: '/personal-center' }, replace: true })
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : t('features.settings.changePassword.error'))
    },
  })

  if (isLoading) {
    return (
      <LoadingState
        title={t('features.settings.profile.title')}
        description={t('features.settings.profile.desc')}
        className='py-12'
      />
    )
  }

  if (error) {
    return (
      <FeedbackState
        title={t('features.settings.profile.title')}
        description={t('features.settings.profile.loadError')}
        tone='danger'
        className='py-12'
      />
    )
  }

  if (!admin) {
    return (
      <FeedbackState
        title={t('features.settings.profile.title')}
        description={t('features.settings.profile.noData')}
        tone='muted'
        className='py-12'
      />
    )
  }

  const initials = (admin.nickname || admin.username || '?').slice(0, 2).toUpperCase()

  const fieldClass = 'space-y-4 max-w-sm'

  return (
    <div className='space-y-6'>
      {/* 个人资料：与下方表单统一的左对齐、标签在上的纵向布局 */}
      <div className={fieldClass}>
        <div className='mb-4'>
          <button
            type='button'
            className='rounded-lg outline-none ring-ring focus-visible:ring-2'
            onClick={() => setAvatarDialogOpen(true)}
            aria-label={t('features.settings.profile.changeAvatar')}
          >
            <Avatar className='h-20 w-20 rounded-lg cursor-pointer'>
              <AvatarImage src={admin.avatar} alt={admin.nickname || admin.username} />
              <AvatarFallback className='rounded-lg text-lg'>{initials}</AvatarFallback>
            </Avatar>
          </button>
        </div>
        <div className='grid gap-4'>
          <div>
            <p className='text-sm font-medium text-muted-foreground'>{t('features.settings.profile.nickname')}</p>
            <p className='text-base font-medium mt-1'>{admin.nickname || '—'}</p>
          </div>
          <div>
            <p className='text-sm font-medium text-muted-foreground'>{t('features.settings.profile.username')}</p>
            <p className='text-base font-medium mt-1'>{admin.username || '—'}</p>
          </div>
          <div>
            <p className='text-sm font-medium text-muted-foreground'>{t('features.settings.profile.mobile')}</p>
            <p className='text-base font-medium mt-1'>{admin.mobile || '—'}</p>
          </div>
          <div>
            <p className='text-sm font-medium text-muted-foreground'>{t('features.settings.profile.introduction')}</p>
            <p className='text-base text-muted-foreground mt-1 whitespace-pre-wrap'>{admin.introduction || '—'}</p>
          </div>
        </div>
      </div>

      <Separator className='my-8' />

      {/* 修改密码：与个人资料同宽、同标签样式与行距 */}
      <div className={fieldClass}>
        <h4 className='text-sm font-medium mb-4'>{t('features.settings.changePassword.title')}</h4>
        <Form {...passwordForm}>
          <form
            onSubmit={passwordForm.handleSubmit((values) => submitChangePassword(values))}
            className='grid gap-4'
          >
            <FormField
              control={passwordForm.control}
              name='old_password'
              render={({ field }) => (
                <FormItem className='gap-1'>
                  <FormLabel className='text-sm font-medium text-muted-foreground'>{t('features.settings.changePassword.oldPassword')}</FormLabel>
                  <FormControl>
                    <PasswordInput placeholder={t('features.settings.changePassword.oldPasswordPlaceholder')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={passwordForm.control}
              name='new_password'
              render={({ field }) => (
                <FormItem className='gap-1'>
                  <FormLabel className='text-sm font-medium text-muted-foreground'>{t('features.settings.changePassword.newPassword')}</FormLabel>
                  <FormControl>
                    <PasswordInput placeholder={t('features.settings.changePassword.newPasswordPlaceholder')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={passwordForm.control}
              name='confirm_password'
              render={({ field }) => (
                <FormItem className='gap-1'>
                  <FormLabel className='text-sm font-medium text-muted-foreground'>{t('features.settings.changePassword.confirmPassword')}</FormLabel>
                  <FormControl>
                    <PasswordInput placeholder={t('features.settings.changePassword.confirmPasswordPlaceholder')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type='submit' disabled={isChangingPassword}>
              {isChangingPassword && <Loader2 className='mr-2 h-4 w-4 animate-spin' />}
              {t('features.settings.changePassword.submit')}
            </Button>
          </form>
        </Form>
      </div>

      <ResourceUpload
        open={avatarDialogOpen}
        onOpenChange={setAvatarDialogOpen}
        onSelect={(resource) => updateAvatar({ admin, avatar: resource.url })}
        type={1}
        title={t('features.settings.profile.changeAvatar')}
        description={t('features.settings.profile.changeAvatarDesc')}
      />
    </div>
  )
}
