import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const BusinessProject = lazy(() => 
  import('@/features/business/project').then(m => ({ 
    default: m.BusinessProject 
  }))
)

export const Route = createFileRoute('/_authenticated/business/project')({
  component: () => (
    <Suspense fallback={<TablePageSkeleton />}>
      <BusinessProject />
    </Suspense>
  ),
})
