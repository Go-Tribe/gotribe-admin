import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { PageSkeleton } from '@/components/page-skeleton'

const Otp = lazy(() => 
  import('@/features/auth/otp').then(m => ({ 
    default: m.Otp 
  }))
)

export const Route = createFileRoute('/(auth)/otp')({
  component: () => (
    <Suspense fallback={<PageSkeleton />}>
      <Otp />
    </Suspense>
  ),
})
