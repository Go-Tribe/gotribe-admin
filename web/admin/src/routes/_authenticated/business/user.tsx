import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const BusinessUser = lazy(() => 
  import('@/features/business/user').then(m => ({ 
    default: m.BusinessUser 
  }))
)

export const Route = createFileRoute('/_authenticated/business/user')({
  component: () => (
    <Suspense fallback={<TablePageSkeleton />}>
      <BusinessUser />
    </Suspense>
  ),
})
