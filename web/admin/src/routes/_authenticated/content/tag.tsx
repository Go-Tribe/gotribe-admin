import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const ContentTag = lazy(() => 
  import('@/features/content/tag').then(m => ({ 
    default: m.ContentTag 
  }))
)

export const Route = createFileRoute('/_authenticated/content/tag')({
  component: () => (
    <Suspense fallback={<TablePageSkeleton />}>
      <ContentTag />
    </Suspense>
  ),
})
