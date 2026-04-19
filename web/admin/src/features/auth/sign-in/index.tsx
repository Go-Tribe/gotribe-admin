import { useSearch } from '@tanstack/react-router'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { LanguageSwitcher } from '@/components/language-switcher'
import { AuthLayout } from '../auth-layout'
import { AuthBrand } from '../auth-brand'
import { UserAuthForm } from './components/user-auth-form'

export function SignIn() {
  const { redirect } = useSearch({ from: '/(auth)/sign-in' })

  return (
    <AuthLayout>
      <div className='absolute right-4 top-4 z-10'>
        <LanguageSwitcher />
      </div>
      <Card className='w-full rounded-3xl border border-border/60 bg-background/90 shadow-2xl backdrop-blur-sm'>
        <CardHeader className='pb-4 pt-6'>
          <AuthBrand />
        </CardHeader>
        <CardContent className='pb-6 pt-0'>
          <UserAuthForm redirectTo={redirect} />
        </CardContent>
      </Card>
    </AuthLayout>
  )
}
