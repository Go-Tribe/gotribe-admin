import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const ContentResource = lazy(() => 
  import('@/features/content/resource').then(m => ({ 
    default: m.ContentResource 
  }))
)

export const Route = createFileRoute('/_authenticated/content/resource')({
  component: () => (
    <Suspense fallback={<TablePageSkeleton />}>
      <ContentResource />
    </Suspense>
  ),
})
