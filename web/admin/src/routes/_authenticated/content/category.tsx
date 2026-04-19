import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const ContentCategory = lazy(() => 
  import('@/features/content/category').then(m => ({ 
    default: m.ContentCategory 
  }))
)

export const Route = createFileRoute('/_authenticated/content/category')({
  component: () => (
    <Suspense fallback={<TablePageSkeleton />}>
      <ContentCategory />
    </Suspense>
  ),
})
