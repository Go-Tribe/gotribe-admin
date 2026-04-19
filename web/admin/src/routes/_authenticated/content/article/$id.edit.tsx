import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { FormPageSkeleton } from '@/components/page-skeleton'

const ArticleFormPage = lazy(() =>
  import('@/features/content/article-form-page').then(m => ({
    default: m.ArticleFormPage
  }))
)

export const Route = createFileRoute('/_authenticated/content/article/$id/edit')({
  component: ArticleEditPage,
})

function ArticleEditPage() {
  const { id } = Route.useParams()
  const numericId = Number(id)
  return (
    <Suspense fallback={<FormPageSkeleton />}>
      <ArticleFormPage id={isNaN(numericId) ? null : numericId} initialPost={null} />
    </Suspense>
  )
}
