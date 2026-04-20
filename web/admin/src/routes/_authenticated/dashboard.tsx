import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { PageSkeleton } from '@/components/page-skeleton'

// 懒加载 Dashboard 页面
const DashboardPage = lazy(() => 
  import('@/features/dashboard/dashboard-page').then(m => ({ 
    default: m.DashboardPage 
  }))
)

export const Route = createFileRoute('/_authenticated/dashboard')({
  component: () => (
    <Suspense fallback={<PageSkeleton />}>
      <DashboardPage />
    </Suspense>
  ),
  // 添加错误边界
  errorComponent: ({ error }) => {
    return (
      <div className="p-8">
        <h1 className="text-2xl font-bold text-red-600">加载仪表板失败</h1>
        <p className="mt-4 text-gray-600">{error.message}</p>
        <button 
          className="mt-4 px-4 py-2 bg-blue-500 text-white rounded"
          onClick={() => window.location.reload()}
        >
          刷新页面
        </button>
      </div>
    )
  },
})
