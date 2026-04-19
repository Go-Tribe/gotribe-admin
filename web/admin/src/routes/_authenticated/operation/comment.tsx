import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const OperationComment = lazy(() => 
  import('@/features/operation/comment').then(m => ({ 
    default: m.OperationComment 
  }))
)

export const Route = createFileRoute('/_authenticated/operation/comment')({
  component: () => (
    <Suspense fallback={<TablePageSkeleton />}>
      <OperationComment />
    </Suspense>
  ),
})
