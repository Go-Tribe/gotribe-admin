import { lazy, Suspense, type ReactNode } from 'react'
import { Skeleton } from '@/components/ui/skeleton'

// 图表骨架屏
function ChartSkeleton({ height = 256 }: { height?: number }) {
  return (
    <div style={{ height }} className='w-full'>
      <Skeleton className='h-full w-full' />
    </div>
  )
}

// 预加载 recharts
export function preloadRecharts() {
  return import('recharts')
}

// 导出常用图表组件的懒加载版本（供直接导入使用）
export const LazyAreaChart = lazy(() => 
  import('recharts').then(m => ({ default: m.AreaChart }))
)

export const LazyBarChart = lazy(() => 
  import('recharts').then(m => ({ default: m.BarChart }))
)

export const LazyPieChart = lazy(() => 
  import('recharts').then(m => ({ default: m.PieChart }))
)

export const LazyLineChart = lazy(() => 
  import('recharts').then(m => ({ default: m.LineChart }))
)

// 辅助组件懒加载
export const LazyArea = lazy(() => 
  import('recharts').then(m => ({ default: m.Area }))
)

export const LazyBar = lazy(() => 
  import('recharts').then(m => ({ default: m.Bar }))
)

export const LazyPie = lazy(() => 
  import('recharts').then(m => ({ default: m.Pie }))
)

export const LazyXAxis = lazy(() => 
  import('recharts').then(m => ({ default: m.XAxis }))
)

export const LazyYAxis = lazy(() => 
  import('recharts').then(m => ({ default: m.YAxis }))
)

export const LazyTooltip = lazy(() => 
  import('recharts').then(m => ({ default: m.Tooltip }))
)

export const LazyLegend = lazy(() => 
  import('recharts').then(m => ({ default: m.Legend }))
)

export const LazyResponsiveContainer = lazy(() => 
  import('recharts').then(m => ({ default: m.ResponsiveContainer }))
)

export const LazyCell = lazy(() => 
  import('recharts').then(m => ({ default: m.Cell }))
)

interface LazyChartContainerProps {
  children: ReactNode
  fallback?: ReactNode
  height?: number
}

/**
 * 懒加载图表容器
 * 包裹图表组件，提供加载骨架屏
 * 
 * @example
 * <LazyChartContainer height={300}>
 *   <LazyAreaChart data={data}>
 *     ...
 *   </LazyAreaChart>
 * </LazyChartContainer>
 */
export function LazyChartContainer({ 
  children, 
  fallback, 
  height = 256 
}: LazyChartContainerProps) {
  return (
    <Suspense fallback={fallback ?? <ChartSkeleton height={height} />}>
      {children}
    </Suspense>
  )
}
