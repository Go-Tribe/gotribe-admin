import { useEffect, useMemo } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Loader2, FileIcon, Copy } from 'lucide-react'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useI18n } from '@/context/i18n-provider'
import { getResourceDetail, updateResource } from '../service/resource'
import { FILE_TYPE, type ResourceApiItem } from '../types/resource'
import { toast } from 'sonner'

const createResourceEditFormSchema = (t: (key: string) => string) =>
  z.object({
    title: z
      .string()
      .min(2, t('features.content.resource.form.validation.titleMin'))
      .max(20, t('features.content.resource.form.validation.titleMax')),
    description: z
      .string()
      .min(2, t('features.content.resource.form.validation.descriptionMin'))
      .max(150, t('features.content.resource.form.validation.descriptionMax')),
  })

type ResourceEditFormValues = z.infer<ReturnType<typeof createResourceEditFormSchema>>

const getFileTypeLabels = (t: (key: string) => string): Record<number, string> => ({
  [FILE_TYPE.IMAGE]: t('features.content.resource.filter.image'),
  [FILE_TYPE.VIDEO]: t('features.content.resource.filter.video'),
  [FILE_TYPE.AUDIO]: t('features.content.resource.filter.audio'),
  [FILE_TYPE.ARCHIVE]: t('features.content.resource.filter.archive'),
  [FILE_TYPE.DOCUMENT]: t('features.content.resource.filter.document'),
  [FILE_TYPE.FONT]: t('features.content.resource.filter.font'),
  [FILE_TYPE.APP]: t('features.content.resource.filter.app'),
  [FILE_TYPE.UNKNOWN]: t('features.content.resource.filter.unknown'),
})

function formatFileSize(bytes?: number): string {
  if (bytes == null) return '-'
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
}

function buildFullUrl(url: string, path: string): string {
  const base = (url ?? '').replace(/\/+$/, '')
  const p = (path ?? '').replace(/^\/+/, '')
  return p ? `${base}/${p}` : base || ''
}

type ResourceEditDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  resourceID: number | null
  onSuccess?: () => void
}

export function ResourceEditDialog({
  open,
  onOpenChange,
  resourceID,
  onSuccess,
}: ResourceEditDialogProps) {
  const { t } = useI18n()
  const queryClient = useQueryClient()
  const resourceEditSchema = useMemo(() => createResourceEditFormSchema(t), [t])
  const fileTypeLabels = useMemo(() => getFileTypeLabels(t), [t])

  const form = useForm<ResourceEditFormValues>({
    resolver: zodResolver(resourceEditSchema),
    defaultValues: { title: '', description: '' },
  })

  const { data: detail, isLoading: loadingDetail } = useQuery({
    queryKey: ['resourceDetail', resourceID],
    queryFn: () => getResourceDetail(resourceID!),
    enabled: open && !!resourceID,
  })

  useEffect(() => {
    if (detail) {
      form.reset({
        title: (detail.title ?? '').trim(),
        description: (detail.description ?? '').trim(),
      })
    }
  }, [detail, form])

  const { mutate: submitMutate, isPending: submitting } = useMutation({
    mutationFn: (values: ResourceEditFormValues) =>
      updateResource(resourceID!, {
        title: values.title.trim(),
        description: values.description.trim(),
      }),
    onSuccess: () => {
      toast.success(t('features.content.resource.form.saveSuccess'))
      onSuccess?.()
      queryClient.invalidateQueries({ queryKey: ['resourceList'] })
      onOpenChange(false)
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : t('features.content.resource.form.saveFailed'))
    },
  })

  const handleCopy = (text: string) => {
    void navigator.clipboard.writeText(text).then(() => toast.success(t('features.content.resource.copySuccessShort')))
  }

  const handleSubmit = (values: ResourceEditFormValues) => {
    submitMutate(values)
  }

  if (!open) return null

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-[500px]'>
        <DialogHeader>
          <DialogTitle>{t('features.content.resource.form.editTitle')}</DialogTitle>
        </DialogHeader>

        {loadingDetail ? (
          <div className='flex min-h-[200px] items-center justify-center py-8'>
            <Loader2 className='h-8 w-8 animate-spin text-muted-foreground' />
          </div>
        ) : detail ? (
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(handleSubmit)}
              className='space-y-4'
            >
              <div>
                <Label className='text-muted-foreground'>{t('features.content.resource.form.preview')}</Label>
                <div className='mt-1.5 flex h-20 w-20 items-center justify-center overflow-hidden rounded border bg-muted'>
                  {detail.fileType === FILE_TYPE.IMAGE ? (
                    <img
                      src={buildFullUrl(detail.url, detail.path)}
                      alt={detail.title}
                      className='h-full w-full object-cover'
                    />
                  ) : (
                    <FileIcon className='h-10 w-10 text-muted-foreground' />
                  )}
                </div>
              </div>

              {/* 资源名称、描述 */}
              <FormField
                control={form.control}
                name='title'
                render={({ field }) => (
                  <FormItem className='space-y-2'>
                    <FormLabel htmlFor='resource-title'>
                      {t('features.content.resource.form.title')}
                    </FormLabel>
                    <FormControl>
                      <Input
                        id='resource-title'
                        placeholder={t('features.content.resource.form.titlePlaceholder')}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='description'
                render={({ field }) => (
                  <FormItem className='space-y-2'>
                    <FormLabel htmlFor='resource-desc'>
                      {t('features.content.resource.form.description')}
                    </FormLabel>
                    <FormControl>
                      <Input
                        id='resource-desc'
                        placeholder={t('features.content.resource.form.descriptionPlaceholder')}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className='grid grid-cols-3 gap-4'>
                <div className='space-y-1'>
                  <Label className='text-muted-foreground'>{t('features.content.resource.form.resourceType')}</Label>
                  <p className='text-sm'>
                    {fileTypeLabels[detail.fileType] ?? t('features.content.resource.form.typeFallback', { type: detail.fileType })}
                  </p>
                </div>
                <div className='space-y-1'>
                  <Label className='text-muted-foreground'>{t('features.content.resource.form.uploadTime')}</Label>
                  <p className='text-sm'>{detail.createdAt ?? '-'}</p>
                </div>
                <div className='space-y-1'>
                  <Label className='text-muted-foreground'>{t('features.content.resource.form.resourceSize')}</Label>
                  <p className='text-sm'>{formatFileSize(detail.size)}</p>
                </div>
              </div>

              <LinkSection detail={detail} onCopy={handleCopy} t={t} />

              <DialogFooter>
                <Button type='button' variant='outline' onClick={() => onOpenChange(false)}>
                  {t('features.content.resource.form.cancel')}
                </Button>
                <Button type='submit' disabled={submitting}>
                  {submitting ? (
                    <span className='inline-flex items-center gap-2'>
                      <Loader2 className='h-4 w-4 animate-spin' />
                      {t('features.content.resource.form.submitting')}
                    </span>
                  ) : (
                    t('features.content.resource.form.submit')
                  )}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        ) : (
          <div className='py-6 text-center text-muted-foreground'>{t('features.content.resource.loadFailed')}</div>
        )}
      </DialogContent>
    </Dialog>
  )
}

function LinkSection({
  detail,
  onCopy,
  t,
}: {
  detail: ResourceApiItem
  onCopy: (text: string) => void
  t: (key: string) => string
}) {
  const fullUrl = buildFullUrl(detail.url, detail.path)
  const path = (detail.path ?? '').trim()
  const title = (detail.title ?? '').trim()
  const html = `<img src="${fullUrl}" alt="${title.replace(/"/g, '&quot;')}" />`
  const markdown = `![${title}](${fullUrl})`
  const copyTitle = t('features.content.resource.actions.copy')

  return (
    <div className='space-y-3'>
      <Label className='text-muted-foreground'>{t('features.content.resource.form.links')}</Label>
      <div className='space-y-2'>
        <div className='space-y-1'>
          <span className='text-xs text-muted-foreground'>URL</span>
          <div className='flex gap-2'>
            <Input readOnly value={path || fullUrl} className='font-mono text-xs' />
            <Button type='button' variant='outline' size='icon' onClick={() => onCopy(path || fullUrl)} title={copyTitle}>
              <Copy className='h-4 w-4' />
            </Button>
          </div>
        </div>
        <div className='space-y-1'>
          <span className='text-xs text-muted-foreground'>HTML</span>
          <div className='flex gap-2'>
            <Input readOnly value={html} className='font-mono text-xs' />
            <Button type='button' variant='outline' size='icon' onClick={() => onCopy(html)} title={copyTitle}>
              <Copy className='h-4 w-4' />
            </Button>
          </div>
        </div>
        <div className='space-y-1'>
          <span className='text-xs text-muted-foreground'>Markdown</span>
          <div className='flex gap-2'>
            <Input readOnly value={markdown} className='font-mono text-xs' />
            <Button type='button' variant='outline' size='icon' onClick={() => onCopy(markdown)} title={copyTitle}>
              <Copy className='h-4 w-4' />
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
