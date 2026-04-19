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
import { OtpForm } from './components/otp-form'

export function Otp() {
  const { t } = useI18n()

  return (
    <AuthLayout>
      <Card className='w-full border shadow-sm'>
        <CardHeader className='pb-4 pt-6'>
          <AuthBrand />
          <div className='space-y-1.5 border-t pt-5'>
            <CardTitle className='text-base font-medium'>
              {t('features.auth.otp.title')}
            </CardTitle>
            <CardDescription>
              {t('features.auth.otp.description')}
            </CardDescription>
          </div>
        </CardHeader>
        <CardContent>
          <OtpForm />
        </CardContent>
        <CardFooter className='border-t pt-6'>
          <p className='w-full text-center text-sm text-muted-foreground'>
            {t('features.auth.otp.footerText')}{' '}
            <Link
              to='/sign-in'
              className='underline underline-offset-4 hover:text-primary'
            >
              {t('features.auth.otp.resendLink')}
            </Link>
            .
          </p>
        </CardFooter>
      </Card>
    </AuthLayout>
  )
}
