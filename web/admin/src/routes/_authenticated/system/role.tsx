import z from 'zod'
import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const SystemRole = lazy(() => 
  import('@/features/system/role').then(m => ({ 
    default: m.SystemRole 
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

export const Route = createFileRoute('/_authenticated/system/role')({
  validateSearch: appsSearchSchema,
  component: () => (
    <Suspense fallback={<TablePageSkeleton />}>
      <SystemRole />
    </Suspense>
  ),
})
