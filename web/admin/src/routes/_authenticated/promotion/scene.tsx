import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const OperationScene = lazy(() => 
  import('@/features/promotion/scene').then(m => ({ 
    default: m.OperationScene 
  }))
)

export const Route = createFileRoute('/_authenticated/promotion/scene')({
  component: () => (
    <Suspense fallback={<TablePageSkeleton />}>
      <OperationScene />
    </Suspense>
  ),
})
