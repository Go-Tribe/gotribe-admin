import { useState } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from '@tanstack/react-router'
import { Loader2, LogIn } from 'lucide-react'
import { toast } from 'sonner'
import { useSetAccessToken, useSetAuthUser } from '@/stores/auth-store'
import { useI18n } from '@/context/i18n-provider'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { PasswordInput } from '@/components/password-input'
import { login, type LoginResult } from '../service'
import { getCurrentUser } from '@/service/user'

interface UserAuthFormProps extends React.HTMLAttributes<HTMLFormElement> {
  redirectTo?: string
}

export function UserAuthForm({
  className,
  redirectTo,
  ...props
}: UserAuthFormProps) {
  const [isLoading, setIsLoading] = useState(false)
  const navigate = useNavigate()
  const setAccessToken = useSetAccessToken()
  const setAuthUser = useSetAuthUser()
  const { t } = useI18n()

  // Create form schema with translations
  const formSchema = z.object({
    username: z.string().min(1, t('features.auth.validation.usernameRequired')),
    password: z.string().min(1, t('features.auth.validation.passwordRequired')),
  })

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: '',
      password: '',
    },
  })

  function onSubmit(data: z.infer<typeof formSchema>) {
    setIsLoading(true)

    toast.promise(
      login({ username: data.username, password: data.password }).then(async (result: LoginResult) => {
        // 先存储 accessToken 到 localStorage
        setAccessToken(result.token)
        // 调用接口获取用户信息
        let displayName = data.username
        try {
          const userInfoResult = await getCurrentUser()
          if (userInfoResult && userInfoResult.admin) {
            const admin = userInfoResult.admin as { nickname?: string; username?: string }
            // 优先展示昵称，其次用户名
            displayName = admin.nickname?.trim() || admin.username || data.username
            setAuthUser({
              ...userInfoResult.admin,
            })
          } else {
            setAuthUser({})
          }
        } catch (error) {
          setAuthUser({})
          // eslint-disable-next-line no-console
          console.error('Failed to get user info:', error)
        }

        // 重定向到目标页面或默认到仪表板
        const targetPath = redirectTo || '/'
        navigate({ to: targetPath, replace: true })

        return t('features.auth.signIn.welcomeBack', { username: displayName })
      }),
      {
        loading: t('features.auth.signIn.loggingIn'),
        success: (message) => {
          setIsLoading(false)
          return message
        },
        error: (error) => {
          setIsLoading(false)
          return error?.message || t('features.auth.signIn.error')
        },
      }
    )
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className={cn('grid gap-3', className)}
        {...props}
      >
        <FormField
          control={form.control}
          name='username'
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t('features.auth.signIn.account')}</FormLabel>
              <FormControl>
                <Input placeholder={t('features.auth.signIn.accountPlaceholder')} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name='password'
          render={({ field }) => (
            <FormItem className='relative'>
              <FormLabel>{t('features.auth.signIn.password')}</FormLabel>
<FormControl>
                  <PasswordInput placeholder={t('features.auth.signIn.passwordPlaceholder')} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button className='mt-2' disabled={isLoading}>
          {isLoading ? <Loader2 className='animate-spin' /> : <LogIn />}
          {t('features.auth.signIn.login')}
        </Button>
        <p className='text-center text-xs text-muted-foreground'>
          {redirectTo ? t('features.auth.signIn.redirectHint') : t('features.auth.signIn.securityHint')}
        </p>
      </form>
    </Form>
  )
}
