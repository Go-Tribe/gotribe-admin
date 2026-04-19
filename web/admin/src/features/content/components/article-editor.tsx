import { Suspense, lazy } from 'react'
import { Textarea } from '@/components/ui/textarea'
import {
  FormField,
  FormItem,
  FormControl,
  FormMessage,
} from '@/components/ui/form'
import { EditorErrorBoundary } from '@/components/editor-error-boundary'
import { useI18n } from '@/context/i18n-provider'
import type { UseFormReturn } from 'react-hook-form'
import type { ArticleFormValues } from '../article-form-page'

const SlateEditor = lazy(() =>
  import('@/components/editor').then((m) => ({ default: m.SlateEditor }))
)

interface ArticleEditorProps {
  form: UseFormReturn<ArticleFormValues>
}

export function ArticleEditor({ form }: ArticleEditorProps) {
  const { t } = useI18n()

  return (
    <div className='mx-auto w-full max-w-[42rem] px-5 sm:px-8 pt-8 pb-24'>
      {/* 标题输入 */}
      <FormField
        control={form.control}
        name="title"
        render={({ field }) => (
          <FormItem className='space-y-0'>
            <FormControl>
              <Textarea
                placeholder={t('features.content.article.form.titlePlaceholder')}
                className="text-3xl sm:text-4xl font-bold border-none resize-none shadow-none focus-visible:ring-0 px-0 py-0 min-h-[2.5rem] overflow-hidden leading-snug placeholder:text-muted-foreground/50 tracking-tight"
                rows={1}
                autoFocus
                onInput={(e) => {
                  const target = e.target as HTMLTextAreaElement
                  target.style.height = 'auto'
                  target.style.height = `${target.scrollHeight}px`
                }}
                {...field}
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />

      {/* 标题与正文之间的留白 */}
      <div className='h-4' aria-hidden />

      {/* Slate 编辑器 */}
      <FormField
        control={form.control}
        name='content'
        render={({ field }) => (
          <FormItem className='flex flex-col min-h-0'>
            <FormControl>
              <EditorErrorBoundary
                fallback={
                  <div className='rounded-md border border-input bg-background flex flex-col flex-1 items-center justify-center text-muted-foreground min-h-[320px]'>
                    <p className='text-sm mb-2'>{t('features.content.article.form.editorLoadError')}</p>
                    <Textarea
                      className='max-w-2xl w-full min-h-[120px] font-mono text-xs resize-none'
                      value={field.value ?? ''}
                      onChange={field.onChange}
                      placeholder={t('features.content.article.form.contentPlaceholder')}
                      readOnly={false}
                    />
                  </div>
                }
              >
                <Suspense
                  fallback={
                    <div className='flex flex-col items-center justify-center text-muted-foreground min-h-[320px]'>
                      <p className='text-sm'>{t('features.content.article.form.editorLoading')}</p>
                    </div>
                  }
                >
                  <SlateEditor
                    value={field.value ?? ''}
                    onChange={field.onChange}
                    minHeight='min-h-[60vh]'
                    outputMode='json'
                    autoHeight={true}
                    className="border-none shadow-none px-0 text-[1.0625rem] leading-[1.75]"
                  />
                </Suspense>
              </EditorErrorBoundary>
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
    </div>
  )
}
