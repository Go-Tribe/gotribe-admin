import z from 'zod'
import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { TablePageSkeleton } from '@/components/page-skeleton'

const SystemMenu = lazy(() => 
  import('@/features/system/menu').then(m => ({ 
    default: m.SystemMenu 
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

export const Route = createFileRoute('/_authenticated/system/menu')({
  validateSearch: appsSearchSchema,
  component: () => (
    <Suspense fallback={<TablePageSkeleton />}>
      <SystemMenu />
    </Suspense>
  ),
})
