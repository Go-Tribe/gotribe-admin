import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const OperationPoint = lazy(() => 
  import('@/features/operation/point').then(m => ({ 
    default: m.OperationPoint 
  }))
)

export const Route = createFileRoute('/_authenticated/operation/point')({
  component: () => (
    <Suspense fallback={<TablePageSkeleton />}>
      <OperationPoint />
    </Suspense>
  ),
})
