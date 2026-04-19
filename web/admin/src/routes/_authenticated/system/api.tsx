import z from 'zod'
import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const SystemApi = lazy(() => 
  import('@/features/system/api').then(m => ({ 
    default: m.SystemApi 
  }))
)

const appsSearchSchema = z.object({
  type: z
    .enum(['all', 'connected', 'notConnected'])
    .optional()
    .catch(undefined),
  filter: z.string().optional().catch(''),
  sort: z.enum(['asc', 'desc']).optional().catch(undefined),
})

export const Route = createFileRoute('/_authenticated/system/api')({
  validateSearch: appsSearchSchema,
  component: () => (
    <Suspense fallback={<TablePageSkeleton />}>
      <SystemApi />
    </Suspense>
  ),
})
