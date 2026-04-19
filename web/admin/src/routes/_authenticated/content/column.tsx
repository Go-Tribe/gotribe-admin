import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const ContentColumn = lazy(() => 
  import('@/features/content/column').then(m => ({ 
    default: m.ContentColumn 
  }))
)

export const Route = createFileRoute('/_authenticated/content/column')({
  component: () => (
    <Suspense fallback={<TablePageSkeleton />}>
      <ContentColumn />
    </Suspense>
  ),
})
