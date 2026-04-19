import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const PromotionAdvertising = lazy(() => 
  import('@/features/promotion/advertising').then(m => ({ 
    default: m.PromotionAdvertising 
  }))
)

export const Route = createFileRoute('/_authenticated/promotion/advertising')({
  component: () => (
    <Suspense fallback={<TablePageSkeleton />}>
      <PromotionAdvertising />
    </Suspense>
  ),
})
