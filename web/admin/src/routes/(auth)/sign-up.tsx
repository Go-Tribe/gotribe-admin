import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { PageSkeleton } from '@/components/page-skeleton'

const SignUp = lazy(() => 
  import('@/features/auth/sign-up').then(m => ({ 
    default: m.SignUp 
  }))
)

export const Route = createFileRoute('/(auth)/sign-up')({
  component: () => (
    <Suspense fallback={<PageSkeleton />}>
      <SignUp />
    </Suspense>
  ),
})
