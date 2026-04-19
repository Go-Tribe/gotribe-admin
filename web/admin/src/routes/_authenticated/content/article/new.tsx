import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { FormPageSkeleton } from '@/components/page-skeleton'

const ArticleFormPage = lazy(() => 
  import('@/features/content/article-form-page').then(m => ({ 
    default: m.ArticleFormPage 
  }))
)

export const Route = createFileRoute('/_authenticated/content/article/new')({
  component: ArticleNewPage,
})

function ArticleNewPage() {
  return (
    <Suspense fallback={<FormPageSkeleton />}>
      <ArticleFormPage id={null} initialPost={null} />
    </Suspense>
  )
}
