import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

/**
 * 页面级骨架屏组件
 * 用于路由懒加载时的加载状态展示
 */
export function PageSkeleton() {
  return (
    <div className='flex flex-1 flex-col gap-4 p-4 md:p-8'>
      {/* Header 骨架 */}
      <div className='space-y-2'>
        <Skeleton className='h-4 w-32' />
        <Skeleton className='h-8 w-64' />
      </div>

      {/* Stats 骨架 */}
      <div className='grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4'>
        {[...Array(4)].map((_, i) => (
          <Card key={i}>
            <CardContent className='p-6'>
              <div className='flex items-start justify-between'>
                <div className='space-y-2'>
                  <Skeleton className='h-4 w-24' />
                  <Skeleton className='h-8 w-20' />
                  <Skeleton className='h-4 w-28' />
                </div>
                <Skeleton className='h-12 w-12 rounded-xl' />
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Content 骨架 */}
      <div className='grid gap-4 lg:grid-cols-3'>
        <Card className='lg:col-span-2'>
          <CardHeader>
            <Skeleton className='h-5 w-32' />
          </CardHeader>
          <CardContent>
            <Skeleton className='h-64 w-full' />
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <Skeleton className='h-5 w-24' />
          </CardHeader>
          <CardContent>
            <Skeleton className='h-48 w-full' />
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

/**
 * 表格页面骨架屏
 */
export function TablePageSkeleton() {
  return (
    <div className='flex flex-1 flex-col gap-4 p-4 md:p-8'>
      {/* Header 骨架 */}
      <div className='space-y-2'>
        <Skeleton className='h-4 w-24' />
        <div className='flex items-center justify-between'>
          <Skeleton className='h-8 w-48' />
          <Skeleton className='h-10 w-24' />
        </div>
      </div>

      {/* Table 骨架 */}
      <Card>
        <CardContent className='p-6'>
          {/* Filter bar */}
          <div className='mb-4 flex gap-2'>
            <Skeleton className='h-10 w-48' />
            <Skeleton className='h-10 w-32' />
            <Skeleton className='h-10 w-32' />
          </div>

          {/* Table rows */}
          <div className='space-y-3'>
            {/* Header */}
            <div className='flex gap-4 border-b pb-3'>
              <Skeleton className='h-4 w-12' />
              <Skeleton className='h-4 flex-1' />
              <Skeleton className='h-4 w-24' />
              <Skeleton className='h-4 w-24' />
              <Skeleton className='h-4 w-20' />
            </div>
            {/* Rows */}
            {[...Array(8)].map((_, i) => (
              <div key={i} className='flex items-center gap-4 py-3'>
                <Skeleton className='h-4 w-12' />
                <Skeleton className='h-4 flex-1' />
                <Skeleton className='h-4 w-24' />
                <Skeleton className='h-4 w-24' />
                <Skeleton className='h-8 w-20' />
              </div>
            ))}
          </div>

          {/* Pagination */}
          <div className='mt-4 flex items-center justify-between border-t pt-4'>
            <Skeleton className='h-4 w-32' />
            <div className='flex gap-2'>
              <Skeleton className='h-8 w-8' />
              <Skeleton className='h-8 w-8' />
              <Skeleton className='h-8 w-8' />
              <Skeleton className='h-8 w-8' />
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

/**
 * 表单页面骨架屏
 */
export function FormPageSkeleton() {
  return (
    <div className='flex flex-1 flex-col gap-4 p-4 md:p-8'>
      <div className='space-y-2'>
        <Skeleton className='h-4 w-32' />
        <Skeleton className='h-8 w-48' />
      </div>

      <Card>
        <CardContent className='p-6'>
          <div className='space-y-6'>
            {/* Form fields */}
            {[...Array(6)].map((_, i) => (
              <div key={i} className='space-y-2'>
                <Skeleton className='h-4 w-24' />
                <Skeleton className='h-10 w-full' />
              </div>
            ))}
            {/* Textarea */}
            <div className='space-y-2'>
              <Skeleton className='h-4 w-24' />
              <Skeleton className='h-32 w-full' />
            </div>
            {/* Buttons */}
            <div className='flex gap-2 pt-4'>
              <Skeleton className='h-10 w-24' />
              <Skeleton className='h-10 w-24' />
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
