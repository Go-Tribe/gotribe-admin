import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { FormPageSkeleton } from '@/components/page-skeleton'

const SettingsProfile = lazy(() => 
  import('@/features/settings/profile').then(m => ({ 
    default: m.SettingsProfile 
  }))
)

export const Route = createFileRoute('/_authenticated/personal-center/')({
  component: () => (
    <Suspense fallback={<FormPageSkeleton />}>
      <SettingsProfile />
    </Suspense>
  ),
})
