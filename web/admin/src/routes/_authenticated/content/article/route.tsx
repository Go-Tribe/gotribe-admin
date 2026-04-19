import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/content/article')({
  component: ArticleLayout,
})

function ArticleLayout() {
  return <Outlet />
}
