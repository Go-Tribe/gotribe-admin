import { z } from 'zod'
import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { PageSkeleton } from '@/components/page-skeleton'

const SignIn = lazy(() => 
  import('@/features/auth/sign-in').then(m => ({ 
    default: m.SignIn 
  }))
)

const searchSchema = z.object({
  redirect: z.string().optional(),
})

export const Route = createFileRoute('/(auth)/sign-in')({
  component: () => (
    <Suspense fallback={<PageSkeleton />}>
      <SignIn />
    </Suspense>
  ),
  validateSearch: searchSchema,
})
