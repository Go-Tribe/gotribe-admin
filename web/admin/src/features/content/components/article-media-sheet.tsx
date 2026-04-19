import { type UseFormReturn } from 'react-hook-form'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
  FormDescription,
} from '@/components/ui/form'
import { Textarea } from '@/components/ui/textarea'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Wand2, ImageIcon, Video, X } from 'lucide-react'
import { cn } from '@/lib/utils'
import { slateToPlainText } from '@/lib/slate-markdown'
import { toast } from 'sonner'
import { useI18n } from '@/context/i18n-provider'
import type { ArticleFormValues } from '../article-form-page'

interface ArticleMediaSheetProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  form: UseFormReturn<ArticleFormValues>
  onOpenResourceUpload: (field: 'icon' | 'video' | 'images') => void
}

export function ArticleMediaSheet({
  open,
  onOpenChange,
  form,
  onOpenResourceUpload,
}: ArticleMediaSheetProps) {
  const { t } = useI18n()

  const handleExtractDescription = () => {
    const content = form.getValues('content')
    if (!content) return

    try {
      const nodes = JSON.parse(content)
      const text = slateToPlainText(nodes)
      const summary = text.slice(0, 200).replace(/\s+/g, ' ').trim()
      if (summary) {
        form.setValue('description', summary, { shouldDirty: true })
        toast.success(t('features.content.article.descriptionExtracted'))
      }
    } catch {
      toast.warning(t('features.content.article.descriptionExtractFailed'))
    }
  }

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetTrigger asChild>
        <Button variant="outline" size="sm" type="button">
          <ImageIcon className="h-4 w-4 mr-2" />
          {t('features.content.article.form.tabMedia')}
        </Button>
      </SheetTrigger>
      <SheetContent className="w-[400px] sm:w-[540px] p-0">
        <SheetHeader className="px-6 py-4 border-b">
          <SheetTitle>{t('features.content.article.form.tabMedia')}</SheetTitle>
          <SheetDescription>
            {t('features.content.article.form.mediaDescription')}
          </SheetDescription>
        </SheetHeader>
        <ScrollArea className="h-[calc(100vh-80px)] px-6 py-4">
          <div className="space-y-6 pb-8">
            {/* 摘要 */}
            {/* Slug */}
            <FormField
              control={form.control}
              name="slug"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('features.content.article.form.slug')}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t('features.content.article.form.slugPlaceholder')}
                      {...field}
                    />
                  </FormControl>
                  <FormDescription className="text-xs text-muted-foreground">
                    {t('features.content.article.form.slugHint')}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <div className="flex items-center justify-between">
                    <FormLabel>{t('features.content.article.form.description')}</FormLabel>
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      className="h-6 px-2 text-xs"
                      onClick={handleExtractDescription}
                    >
                      <Wand2 className="w-3 h-3 mr-1" />
                      {t('features.content.article.form.extractDescription')}
                    </Button>
                  </div>
                  <FormControl>
                    <Textarea
                      placeholder={t('features.content.article.form.descriptionPlaceholder')}
                      className="resize-none min-h-[88px]"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* 封面与视频 */}
            <div className="space-y-2">
              <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                {t('features.content.article.form.cover')} / {t('features.content.article.form.video')}
              </label>
              <div className="flex flex-wrap items-center gap-4">
                <FormField
                  control={form.control}
                  name="icon"
                  render={({ field }) => (
                    <FormItem className="space-y-0">
                      <FormControl>
                        <button
                          type="button"
                          onClick={() => onOpenResourceUpload('icon')}
                          className={cn(
                            'flex flex-col items-center justify-center h-20 w-20 rounded-lg border-2 border-dashed border-border bg-muted/30 overflow-hidden transition-all shrink-0',
                            'cursor-pointer hover:border-primary/50 hover:bg-muted/50'
                          )}
                        >
                          {field.value ? (
                            <img src={field.value} alt="" className="h-full w-full object-cover" />
                          ) : (
                            <>
                              <ImageIcon className="h-8 w-8 text-muted-foreground mb-1" />
                              <span className="text-xs text-muted-foreground">{t('features.content.article.form.selectCover')}</span>
                            </>
                          )}
                        </button>
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="video"
                  render={({ field }) => (
                    <FormItem className="space-y-0">
                      <FormControl>
                        <button
                          type="button"
                          onClick={() => onOpenResourceUpload('video')}
                          className={cn(
                            'flex flex-col items-center justify-center h-20 w-20 rounded-lg border-2 border-dashed border-border bg-muted/30 overflow-hidden transition-all shrink-0',
                            'cursor-pointer hover:border-primary/50 hover:bg-muted/50'
                          )}
                        >
                          {field.value ? (
                            <div className="w-full h-full flex flex-col items-center justify-center bg-muted/30 p-1">
                              <Video className="h-8 w-8 text-muted-foreground mb-1 shrink-0" />
                              <span className="text-xs text-muted-foreground truncate w-full text-center">{field.value.split('/').pop()}</span>
                            </div>
                          ) : (
                            <>
                              <Video className="h-8 w-8 text-muted-foreground mb-1" />
                              <span className="text-xs text-muted-foreground">{t('features.content.article.form.selectVideo')}</span>
                            </>
                          )}
                        </button>
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </div>

            {/* 更多图片 */}
            <FormField
              control={form.control}
              name="images"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-muted-foreground text-sm">{t('features.content.article.form.moreImages')}</FormLabel>
                  <FormControl>
                    <div className="flex flex-wrap gap-2">
                      {(field.value ?? []).map((url, index) => (
                        <div key={`${url}-${index}`} className="relative group">
                          <img src={url} alt="" className="h-20 w-20 rounded-lg border object-cover shrink-0" />
                          <button
                            type="button"
                            onClick={() => {
                              const next = (field.value ?? []).filter((_, i) => i !== index)
                              field.onChange(next)
                            }}
                            className="absolute top-0.5 right-0.5 h-5 w-5 rounded-full bg-background/95 text-foreground shadow-sm border border-border flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity hover:bg-muted"
                            aria-label={t('features.content.article.form.removeImage')}
                          >
                            <X className="h-3 w-3" />
                          </button>
                        </div>
                      ))}
                      <button
                        type="button"
                        onClick={() => onOpenResourceUpload('images')}
                        className={cn(
                          'flex flex-col items-center justify-center h-20 w-20 rounded-lg border-2 border-dashed border-border bg-muted/30 shrink-0',
                          'cursor-pointer hover:border-primary/50 hover:bg-muted/50'
                        )}
                      >
                        <ImageIcon className="h-8 w-8 text-muted-foreground mb-1" />
                        <span className="text-xs text-muted-foreground">{t('features.content.article.form.addImage')}</span>
                      </button>
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>
        </ScrollArea>
      </SheetContent>
    </Sheet>
  )
}
