import { useEffect, useMemo } from 'react'
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
import { useI18n } from '@/context/i18n-provider'
import type { SceneCreateParams, SceneUpdateParams, Scene } from '../types'

const createSceneFormSchema = (t: (key: string) => string) =>
  z.object({
    title: z.string().min(1, t('features.operation.scene.form.validation.titleRequired')),
    description: z
      .string()
      .min(1, t('features.operation.scene.form.validation.descriptionRequired')),
    projectId: z.number().min(1, t('features.operation.scene.form.validation.projectRequired')),
  })

type SceneFormValues = z.infer<ReturnType<typeof createSceneFormSchema>>

type SceneFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSubmit?: (data: SceneCreateParams) => void
  onSubmitUpdate?: (adSceneID: string, data: SceneUpdateParams) => void
  isLoading?: boolean
  projectList: { id: number; title: string }[]
  editScene?: Scene | null
}

export function SceneFormDialog({
  open,
  onOpenChange,
  onSubmit,
  onSubmitUpdate,
  isLoading = false,
  projectList,
  editScene,
}: SceneFormDialogProps) {
  const { t } = useI18n()
  const isEdit = Boolean(editScene)
  const sceneFormSchema = useMemo(() => createSceneFormSchema(t), [t])
  const form = useForm<SceneFormValues>({
    resolver: zodResolver(sceneFormSchema),
    defaultValues: {
      title: '',
      description: '',
      projectId: undefined,
    },
  })

  useEffect(() => {
    if (!open) return
    if (editScene) {
      form.reset({
        title: (editScene.title ?? '').trim(),
        description: (editScene.description ?? '').trim(),
        projectId: editScene.projectId ?? undefined,
      })
    } else {
      form.reset({
        title: '',
        description: '',
        projectId: undefined,
      })
    }
  }, [open, editScene, form])

  function handleSubmit(values: SceneFormValues) {
    if (isEdit && editScene && onSubmitUpdate) {
      onSubmitUpdate(String(editScene.id), {
        title: values.title.trim(),
        description: values.description.trim(),
        projectId: values.projectId || undefined,
      })
    } else if (onSubmit) {
      onSubmit({
        title: values.title.trim(),
        description: values.description.trim(),
        projectId: values.projectId,
      })
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-[600px] max-h-[90vh] flex flex-col'>
        <DialogHeader className='shrink-0'>
          <DialogTitle>
            {isEdit
              ? t('features.operation.scene.form.editTitle')
              : t('features.operation.scene.form.dialogTitle')}
          </DialogTitle>
          <DialogDescription>
            {isEdit
              ? t('features.operation.scene.form.editDescription')
              : t('features.operation.scene.form.createDescription')}
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
                  <FormLabel>
                    {t('features.operation.scene.form.title')}
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t('features.operation.scene.form.titlePlaceholder')}
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
                  <FormLabel>
                    {t('features.operation.scene.form.description')}
                  </FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder={t('features.operation.scene.form.descriptionPlaceholder')}
                      className='min-h-[80px] resize-none'
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
                const selectValue = displayValue || '__empty__'
                return (
                  <FormItem className='space-y-2'>
                    <FormLabel>
                      {t('features.operation.scene.form.project')}
                    </FormLabel>
                    <Select
                      key={`project-${editScene?.id ?? 'create'}-${projectList.length}-${selectValue}`}
                      value={selectValue}
                      onValueChange={(v) => field.onChange(v === '__empty__' ? undefined : Number(v))}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue
                            placeholder={t('features.operation.scene.form.projectPlaceholder')}
                          />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value='__empty__'>
                          {t('features.operation.scene.form.projectPlaceholder')}
                        </SelectItem>
                        {projectList
                          .filter((p) => p.id != null)
                          .map((p) => (
                            <SelectItem
                              key={p.id}
                              value={String(p.id)}
                            >
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
            <DialogFooter className='shrink-0'>
              <Button
                type='button'
                variant='outline'
                onClick={() => onOpenChange(false)}
                disabled={isLoading}
              >
                {t('features.operation.scene.form.cancel')}
              </Button>
              <Button type='submit' disabled={isLoading}>
                {isLoading ? '...' : t('features.operation.scene.form.submit')}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
