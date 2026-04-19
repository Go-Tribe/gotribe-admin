import { Link } from '@tanstack/react-router'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { useI18n } from '@/context/i18n-provider'
import { AuthLayout } from '../auth-layout'
import { AuthBrand } from '../auth-brand'
import { SignUpForm } from './components/sign-up-form'

export function SignUp() {
  const { t } = useI18n()

  return (
    <AuthLayout>
      <Card className='w-full border shadow-sm'>
        <CardHeader className='pb-4 pt-6'>
          <AuthBrand />
          <div className='space-y-1.5 border-t pt-5'>
            <CardTitle className='text-base font-medium'>
              {t('features.auth.signUp.title')}
            </CardTitle>
            <CardDescription>
              {t('features.auth.signUp.description')}{' '}
              <Link
                to='/sign-in'
                className='underline underline-offset-4 hover:text-primary'
              >
                {t('features.auth.signUp.signInLink')}
              </Link>
            </CardDescription>
          </div>
        </CardHeader>
        <CardContent>
          <SignUpForm />
        </CardContent>
        <CardFooter className='border-t pt-6'>
          <p className='w-full text-center text-sm text-muted-foreground'>
            {t('features.auth.signUp.footerText')}{' '}
            <a
              href='/terms'
              className='underline underline-offset-4 hover:text-primary'
            >
              {t('features.auth.signUp.termsOfService')}
            </a>{' '}
            {t('features.auth.signUp.and')}{' '}
            <a
              href='/privacy'
              className='underline underline-offset-4 hover:text-primary'
            >
              {t('features.auth.signUp.privacyPolicy')}
            </a>
            .
          </p>
        </CardFooter>
      </Card>
    </AuthLayout>
  )
}
