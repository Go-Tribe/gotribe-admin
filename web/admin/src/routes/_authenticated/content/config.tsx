import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const ContentConfig = lazy(() => 
  import('@/features/content/config').then(m => ({ 
    default: m.ContentConfig 
  }))
)

export const Route = createFileRoute('/_authenticated/content/config')({
  component: () => (
    <Suspense fallback={<TablePageSkeleton />}>
      <ContentConfig />
    </Suspense>
  ),
})
