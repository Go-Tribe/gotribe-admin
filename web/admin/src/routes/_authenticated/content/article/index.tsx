import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const ContentArticle = lazy(() => 
  import('@/features/content/article').then(m => ({ 
    default: m.ContentArticle 
  }))
)

export const Route = createFileRoute('/_authenticated/content/article/')({
  component: () => (
    <Suspense fallback={<TablePageSkeleton />}>
      <ContentArticle />
    </Suspense>
  ),
})
