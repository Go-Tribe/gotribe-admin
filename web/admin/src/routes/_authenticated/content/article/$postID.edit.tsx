import { createFileRoute } from '@tanstack/react-router'
import { lazy, Suspense } from 'react'
import { FormPageSkeleton } from '@/components/page-skeleton'

const ArticleFormPage = lazy(() => 
  import('@/features/content/article-form-page').then(m => ({ 
    default: m.ArticleFormPage 
  }))
)

export const Route = createFileRoute('/_authenticated/content/article/$postID/edit')({
  component: ArticleEditPage,
})

function ArticleEditPage() {
  const { postID } = Route.useParams()
  return (
    <Suspense fallback={<FormPageSkeleton />}>
      <ArticleFormPage postID={postID} initialPost={null} />
    </Suspense>
  )
}
