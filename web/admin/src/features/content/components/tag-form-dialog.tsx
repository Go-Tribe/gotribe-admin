import { useMemo } from 'react'
import * as z from 'zod'
import { SchemaFormDialog, type FormFieldConfig } from '@/components/schema-form-dialog'
import { useI18n } from '@/context/i18n-provider'
import type { Tag } from '../types/tag'

const createTagFormSchema = (t: (key: string) => string) =>
  z.object({
    title: z.string().min(1, t('features.content.tag.form.validation.titleRequired')),
    slug: z
      .string()
      .min(2, t('features.content.tag.form.validation.slugRequired'))
      .max(30, t('features.content.tag.form.validation.slugLength')),
    description: z.string().optional(),
    color: z.string().optional(),
  })

type TagFormValues = z.infer<ReturnType<typeof createTagFormSchema>>

type TagFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  tag: Tag | null
  onSubmit: (data: Partial<Tag>) => void
  isLoading?: boolean
}

/**
 * 标签表单对话框 - 垂直紧凑布局
 */
export function TagFormDialog({
  open,
  onOpenChange,
  tag,
  onSubmit,
  isLoading = false,
}: TagFormDialogProps) {
  const { t } = useI18n()
  const isEdit = !!tag
  const tagFormSchema = useMemo(() => createTagFormSchema(t), [t])

  // 表单字段配置 - 垂直紧凑布局
  const fields: FormFieldConfig[] = useMemo(
    () => [
      {
        name: 'title',
        type: 'text',
        label: t('features.content.tag.form.title'),
        required: true,
        placeholder: t('features.content.tag.form.titlePlaceholder'),
      },
      {
        name: 'slug',
        type: 'text',
        label: t('features.content.tag.form.slug'),
        required: true,
        placeholder: t('features.content.tag.form.slugPlaceholder'),
        description: t('features.content.tag.form.slugDesc') || '用于URL的唯一标识，如：default-tag',
      },
      {
        name: 'color',
        type: 'color',
        label: t('features.content.tag.form.color'),
      },
      {
        name: 'description',
        type: 'textarea',
        label: t('features.content.tag.form.description'),
        placeholder: t('features.content.tag.form.descriptionPlaceholder'),
      },
    ],
    [t]
  )

  // 默认值
  const defaultValues: Partial<TagFormValues> = useMemo(
    () =>
      isEdit && tag
        ? {
            title: tag.title,
            slug: tag.slug,
            description: tag.description,
            color: tag.color,
          }
        : {
            title: '',
            slug: '',
            description: '',
            color: '#3b82f6',
          },
    [isEdit, tag]
  )

  const handleSubmit = (values: TagFormValues) => {
    onSubmit({
      ...values,
      id: tag?.id,
    })
  }

  return (
    <SchemaFormDialog
      open={open}
      onOpenChange={onOpenChange}
      title={t('features.content.tag.form.createTitle')}
      editTitle={t('features.content.tag.form.editTitle')}
      description={
        isEdit
          ? t('features.content.tag.form.editDescription')
          : t('features.content.tag.form.createDescription')
      }
      schema={tagFormSchema}
      fields={fields}
      defaultValues={defaultValues}
      isEdit={isEdit}
      onSubmit={handleSubmit}
      isLoading={isLoading}
      maxWidth="md"
      submitText={isEdit ? t('features.content.tag.form.save') : t('features.content.tag.form.create')}
      cancelText={t('features.content.tag.form.cancel')}
    />
  )
}

export default TagFormDialog
