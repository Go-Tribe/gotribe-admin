import { useEffect, useMemo, useRef, useState, lazy, Suspense } from 'react'
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
import { EditorErrorBoundary } from '@/components/editor-error-boundary'
import { useI18n } from '@/context/i18n-provider'
import type { Config, ConfigCreateParams, ConfigUpdateParams } from '../types/config'

const SlateEditor = lazy(() =>
  import('@/components/editor').then((m) => ({ default: m.SlateEditor }))
)

const JsonEditor = lazy(() =>
  import('@/components/json-editor').then((m) => ({ default: m.JsonEditor }))
)

const createConfigFormSchema = (t: (key: string) => string) =>
  z.object({
    title: z.string().min(1, t('features.content.config.form.validation.titleRequired')),
    description: z.string().min(1, t('features.content.config.form.validation.descriptionRequired')),
    projectId: z.number().min(1, t('features.content.config.form.validation.projectRequired')),
    alias: z.string().min(1, t('features.content.config.form.validation.aliasRequired')),
    type: z.union([z.literal(1), z.literal(2)]),
    mdContent: z.string().optional(),
  })

type ConfigFormValues = z.infer<ReturnType<typeof createConfigFormSchema>>

type ConfigFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  /** 创建时提交 */
  onSubmit?: (data: ConfigCreateParams) => void
  /** 编辑时提交，configID 由传入的 config 带出 */
  onSubmitUpdate?: (configID: string, data: ConfigUpdateParams) => void
  isLoading?: boolean
  projectList: { id: number; title: string }[]
  /** 编辑时的配置详情（由父组件拉取后传入），有值即为编辑模式 */
  editConfig?: Config | null
}

export function ConfigFormDialog({
  open,
  onOpenChange,
  onSubmit,
  onSubmitUpdate,
  isLoading = false,
  projectList,
  editConfig,
}: ConfigFormDialogProps) {
  const { t } = useI18n()
  const isEdit = Boolean(editConfig)
  const [editValuesApplied, setEditValuesApplied] = useState(false)
  const [jsonEditorMounted, setJsonEditorMounted] = useState(false)
  const configFormSchema = useMemo(() => createConfigFormSchema(t), [t])
  const form = useForm<ConfigFormValues>({
    resolver: zodResolver(configFormSchema),
    defaultValues: {
      title: '',
      description: '',
      projectId: undefined,
      alias: '',
      type: 1,
      mdContent: '',
    },
  })

  const typeValue = form.watch('type')

  useEffect(() => {
    if (!open) {
      setEditValuesApplied(false)
      setJsonEditorMounted(false)
      return
    }
    if (editConfig) {
      form.reset({
        title: editConfig.title ?? '',
        description: editConfig.description ?? '',
        projectId: editConfig.projectId ?? undefined,
        alias: (editConfig.alias ?? '').trim(),
        type: (editConfig.type === 2 ? 2 : 1) as 1 | 2,
        mdContent: editConfig.info ?? editConfig.mdContent ?? '',
      })
      setEditValuesApplied(true)
      if (editConfig.type === 2) {
        const id = requestAnimationFrame(() => setJsonEditorMounted(true))
        return () => cancelAnimationFrame(id)
      }
      setJsonEditorMounted(true)
    } else {
      form.reset({
        title: '',
        description: '',
        projectId: undefined,
        alias: '',
        type: 1,
        mdContent: '',
      })
      setEditValuesApplied(true)
      setJsonEditorMounted(true)
    }
  }, [open, editConfig, form])

  const prevTypeRef = useRef<1 | 2>(1)
  useEffect(() => {
    if (!open) return
    if (prevTypeRef.current === 1 && typeValue === 2) {
      form.setValue('mdContent', '')
    }
    prevTypeRef.current = typeValue
  }, [open, typeValue, form])

  const handleSubmit = (values: ConfigFormValues) => {
    const mdContent = values.mdContent ?? ''
    if (isEdit && editConfig && onSubmitUpdate) {
      onSubmitUpdate(editConfig.configID.trim(), {
        title: values.title,
        description: values.description,
        projectId: values.projectId || undefined,
        info: mdContent,
        mdContent,
      })
    } else if (onSubmit) {
      onSubmit({
        title: values.title,
        description: values.description,
        projectId: values.projectId,
        alias: values.alias,
        type: values.type,
        mdContent,
        info: mdContent,
      })
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-[600px] max-h-[90vh] flex flex-col'>
        <DialogHeader className='shrink-0'>
          <DialogTitle>{isEdit ? t('features.content.config.form.editTitle') : t('features.content.config.form.createTitle')}</DialogTitle>
          <DialogDescription>
            {isEdit ? t('features.content.config.form.editDescription') : t('features.content.config.form.createDescription')}
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
                  <FormLabel>{t('features.content.config.form.title')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('features.content.config.form.titlePlaceholder')} {...field} />
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
                  <FormLabel>{t('features.content.config.form.description')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('features.content.config.form.descriptionPlaceholder')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='projectId'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.config.form.project')}</FormLabel>
                  <Select onValueChange={(v) => field.onChange(Number(v))} value={field.value != null ? String(field.value) : ''}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder={t('features.content.config.form.projectPlaceholder')} />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {projectList.map((p) => (
                        <SelectItem key={p.id} value={String(p.id)}>
                          {p.title}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='alias'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.config.form.alias')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('features.content.config.form.aliasPlaceholder')} {...field} disabled={isEdit} className={isEdit ? 'opacity-60' : ''} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='type'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.config.form.type')}</FormLabel>
                  <Select
                    onValueChange={(v) => field.onChange(Number(v) as 1 | 2)}
                    value={String(field.value)}
                    disabled={isEdit}
                  >
                    <FormControl>
                      <SelectTrigger className={isEdit ? 'opacity-60' : ''}>
                        <SelectValue placeholder={t('features.content.config.form.typePlaceholder')} />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value='1'>{t('features.content.config.form.typeRichText')}</SelectItem>
                      <SelectItem value='2'>{t('features.content.config.form.typeJson')}</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='mdContent'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.config.form.content')}</FormLabel>
                  <FormControl>
                    <div className='min-h-[280px] rounded-md border border-border overflow-hidden'>
                      {isEdit && !editValuesApplied ? (
                        <div className='min-h-[280px] flex items-center justify-center text-muted-foreground text-sm bg-muted/30'>
                          {t('features.content.config.form.contentLoading')}
                        </div>
                      ) : typeValue === 2 && isEdit && !jsonEditorMounted ? (
                        <div className='min-h-[280px] flex items-center justify-center text-muted-foreground text-sm bg-muted/30'>
                          {t('features.content.config.form.contentLoading')}
                        </div>
                      ) : typeValue === 1 ? (
                        <EditorErrorBoundary
                          fallback={
                            <Textarea
                              className='min-h-[280px] font-mono text-sm resize-none border-0 rounded-md'
                              value={field.value ?? ''}
                              onChange={(e) => field.onChange(e.target.value)}
                              placeholder={t('features.content.config.form.contentPlaceholder')}
                            />
                          }
                        >
                          <Suspense
                            fallback={
                              <div className='min-h-[280px] flex items-center justify-center text-muted-foreground text-sm bg-muted/30'>
                                {t('features.content.config.form.contentLoading')}
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
                      ) : (
                        <Suspense
                          fallback={
                            <div className='min-h-[280px] flex items-center justify-center text-muted-foreground text-sm bg-muted/30'>
                              {t('features.content.config.form.contentLoading')}
                            </div>
                          }
                        >
                          <JsonEditor
                            key={isEdit && editConfig ? `edit-${editConfig.configID.trim()}` : 'create'}
                            value={field.value ?? ''}
                            onChange={field.onChange}
                            minHeight='min-h-[280px]'
                            initialContent={isEdit && editConfig ? (editConfig.info ?? editConfig.mdContent ?? '') : undefined}
                          />
                        </Suspense>
                      )}
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
                disabled={isLoading}
              >
                {t('features.content.config.form.cancel')}
              </Button>
              <Button type='submit' disabled={isLoading}>
                {isLoading ? t('features.content.config.form.submitting') : t('features.content.config.form.submit')}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
