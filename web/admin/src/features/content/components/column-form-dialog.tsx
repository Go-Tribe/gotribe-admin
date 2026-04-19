import { useEffect, useMemo, useState, lazy, Suspense } from 'react'
import { ImageIcon } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import {
  Dialog,
  DialogContent,
  DialogDescription,
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
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { EditorErrorBoundary } from '@/components/editor-error-boundary'
import { ResourceUpload, type ResourceItem } from '@/components/resource-upload'
import { useI18n } from '@/context/i18n-provider'
import type {
  Column,
  ColumnCreateParams,
  ColumnUpdateParams,
} from '../types/column'

const SlateEditor = lazy(() =>
  import('@/components/editor').then((m) => ({ default: m.SlateEditor }))
)

const createColumnFormSchema = (t: (key: string) => string) =>
  z.object({
    title: z.string().min(1, t('features.content.column.form.validation.titleRequired')),
    info: z.string().min(1, t('features.content.column.form.validation.infoRequired')),
    description: z.string().min(1, t('features.content.column.form.validation.descriptionRequired')),
    projectId: z.number().optional(), // 允许空（编辑时可能无项目；新建时后端会校验）
    icon: z.string().min(1, t('features.content.column.form.validation.iconRequired')),
  })

type ColumnFormValues = z.infer<ReturnType<typeof createColumnFormSchema>>

type ColumnFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSubmit?: (data: ColumnCreateParams) => void
  onSubmitUpdate?: (columnID: string, data: ColumnUpdateParams) => void
  isLoading?: boolean
  projectList: { id: number; title: string }[]
  editColumn?: Column | null
}

export function ColumnFormDialog({
  open,
  onOpenChange,
  onSubmit,
  onSubmitUpdate,
  isLoading = false,
  projectList,
  editColumn,
}: ColumnFormDialogProps) {
  const { t } = useI18n()
  const isEdit = Boolean(editColumn)
  const [iconResourceDialogOpen, setIconResourceDialogOpen] = useState(false)
  const columnFormSchema = useMemo(() => createColumnFormSchema(t), [t])
  const form = useForm<ColumnFormValues>({
    resolver: zodResolver(columnFormSchema),
    defaultValues: {
      title: '',
      info: '',
      description: '',
      projectId: undefined,
      icon: '',
    },
  })

  useEffect(() => {
    if (!open) return
    if (editColumn) {
      const projectId = editColumn.projectId ?? undefined
      form.reset({
        title: editColumn.title ?? '',
        info: editColumn.info ?? '',
        description: editColumn.description ?? '',
        projectId,
        icon: editColumn.icon ?? '',
      })
    } else {
      form.reset({
        title: '',
        info: '',
        description: '',
        projectId: undefined,
        icon: '',
      })
    }
  }, [open, editColumn, form])

  function handleSubmit(values: ColumnFormValues) {
    const projectId = values.projectId ?? 0
    if (isEdit && editColumn) {
      onSubmitUpdate?.(String(editColumn.id), {
        title: values.title,
        info: values.info ?? '',
        description: values.description ?? '',
        icon: values.icon ?? '',
        projectId: projectId || undefined,
      })
    } else {
      onSubmit?.({
        title: values.title,
        description: values.description ?? '',
        info: values.info ?? '',
        projectId,
        icon: values.icon ?? '',
      })
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-[600px] max-h-[90vh] flex flex-col'>
        <DialogHeader className='shrink-0'>
          <DialogTitle>
            {isEdit
              ? t('features.content.column.form.editTitle')
              : t('features.content.column.form.createTitle')}
          </DialogTitle>
          <DialogDescription>
            {isEdit
              ? t('features.content.column.form.editDescription')
              : t('features.content.column.form.createDescription')}
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className='flex-1 overflow-y-auto pr-2 space-y-4 min-h-0'
          >
            <FormField
              control={form.control}
              name='title'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.column.form.title')}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t('features.content.column.form.titlePlaceholder')}
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
                  <FormLabel>{t('features.content.column.form.description')}</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder={t('features.content.column.form.descriptionPlaceholder')}
                      className='resize-none'
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='projectId'
              render={({ field }) => {
                const displayValue = field.value != null ? String(field.value) : ''
                const selectValue = displayValue ? displayValue : '__empty__'
                return (
                  <FormItem className='space-y-2'>
                    <FormLabel>{t('features.content.column.form.project')}</FormLabel>
                    <Select
                      key={`project-${editColumn?.id ?? 'create'}-${projectList.length}-${selectValue}`}
                      onValueChange={(v) => field.onChange(v === '__empty__' ? undefined : Number(v))}
                      value={selectValue}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue
                            placeholder={t('features.content.column.form.projectPlaceholder')}
                          />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value='__empty__'>
                          {t('features.content.column.form.projectPlaceholder')}
                        </SelectItem>
                        {projectList
                          .filter((p) => p.id != null)
                          .map((p) => (
                            <SelectItem key={p.id} value={String(p.id)}>
                              {p.title}
                            </SelectItem>
                          ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )
              }}
            />
            <FormField
              control={form.control}
              name='icon'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.column.form.icon')}</FormLabel>
                  <FormControl>
                    <button
                      type='button'
                      disabled={isEdit}
                      onClick={() => !isEdit && setIconResourceDialogOpen(true)}
                      className={cn(
                        'flex flex-col items-center justify-center h-20 w-20 rounded-lg border-2 border-dashed border-border bg-muted/30 overflow-hidden transition-all shrink-0',
                        !isEdit && 'cursor-pointer hover:border-primary/50 hover:bg-muted/50',
                        isEdit && 'cursor-default border-solid'
                      )}
                    >
                      {field.value ? (
                        <img
                          src={field.value}
                          alt=''
                          className='h-full w-full object-cover'
                        />
                      ) : (
                        <>
                          <ImageIcon className='h-8 w-8 text-muted-foreground mb-1' />
                          <span className='text-xs text-muted-foreground'>
                            {t('features.content.column.form.selectIcon')}
                          </span>
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
              name='info'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.column.form.info')}</FormLabel>
                  <FormControl>
                    <div className='min-h-[280px] rounded-md border border-border overflow-hidden'>
                      <EditorErrorBoundary
                        fallback={
                          <Textarea
                            className='min-h-[280px] font-mono text-sm resize-none border-0 rounded-md'
                            value={field.value ?? ''}
                            onChange={(e) => field.onChange(e.target.value)}
                            placeholder={t('features.content.column.form.infoPlaceholder')}
                          />
                        }
                      >
                        <Suspense
                          fallback={
                            <div className='min-h-[280px] flex items-center justify-center text-muted-foreground text-sm bg-muted/30'>
                              {t('features.content.column.form.infoLoading')}
                            </div>
                          }
                        >
                          <SlateEditor
                            value={field.value ?? ''}
                            onChange={field.onChange}
                            minHeight='min-h-[280px]'
                            outputMode='json'
                          />
                        </Suspense>
                      </EditorErrorBoundary>
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter className='shrink-0 pt-4 border-t mt-4'>
              <Button
                type='button'
                variant='outline'
                onClick={() => onOpenChange(false)}
              >
                {t('features.content.column.form.cancel')}
              </Button>
              <Button type='submit' disabled={isLoading}>
                {isLoading
                  ? t('features.content.column.form.submitting')
                  : t('features.content.column.form.submit')}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
      <ResourceUpload
        open={iconResourceDialogOpen}
        onOpenChange={setIconResourceDialogOpen}
        onSelect={(resource: ResourceItem) => {
          form.setValue('icon', resource.url, { shouldValidate: true })
          setIconResourceDialogOpen(false)
        }}
        type={1}
        title={t('features.content.column.form.selectIconDialogTitle')}
        description={t('features.content.column.form.selectIconDialogDesc')}
      />
    </Dialog>
  )
}
