import { useMemo } from 'react'
import * as z from 'zod'
import { SchemaFormDialog, type FormFieldConfig } from '@/components/schema-form-dialog'
import { useI18n } from '@/context/i18n-provider'
import type { Api } from '../types/api'

const HTTP_METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS']

const createApiFormSchema = (t: (key: string) => string) =>
  z.object({
    method: z.string().min(1, t('features.system.api.form.validation.methodRequired')),
    path: z.string().min(1, t('features.system.api.form.validation.pathRequired')),
    category: z.string().min(1, t('features.system.api.form.validation.categoryRequired')),
    desc: z.string().optional(),
  })

type ApiFormValues = z.infer<ReturnType<typeof createApiFormSchema>>

type ApiFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  api: Api | null
  onSubmit: (data: Partial<Api>) => void
  isLoading?: boolean
}

/**
 * API表单对话框 - 使用 SchemaFormDialog 通用组件
 */
export function ApiFormDialog({
  open,
  onOpenChange,
  api,
  onSubmit,
  isLoading = false,
}: ApiFormDialogProps) {
  const { t } = useI18n()
  const isEdit = !!api

  const apiFormSchema = useMemo(() => createApiFormSchema(t), [t])

  // 字段配置
  const fields: FormFieldConfig[] = useMemo(
    () => [
      {
        name: 'method',
        type: 'select',
        label: t('features.system.api.method'),
        required: true,
        placeholder: t('features.system.api.form.methodPlaceholder'),
        options: HTTP_METHODS.map((method) => ({ label: method, value: method })),
      },
      {
        name: 'path',
        type: 'text',
        label: t('features.system.api.path'),
        required: true,
        placeholder: t('features.system.api.form.pathPlaceholder'),
        description: t('features.system.api.form.pathDescription'),
      },
      {
        name: 'category',
        type: 'text',
        label: t('features.system.api.category'),
        required: true,
        placeholder: t('features.system.api.form.categoryPlaceholder'),
        description: t('features.system.api.form.categoryDescription'),
      },
      {
        name: 'desc',
        type: 'textarea',
        label: t('features.system.api.desc'),
        placeholder: t('features.system.api.form.descPlaceholder'),
      },
    ],
    [t]
  )

  // 默认值
  const defaultValues = useMemo(
    () =>
      isEdit && api
        ? {
            method: api.method,
            path: api.path,
            category: api.category,
            desc: api.desc || '',
          }
        : {
            method: 'GET',
            path: '',
            category: '',
            desc: '',
          },
    [isEdit, api]
  )

  const handleSubmit = (values: ApiFormValues) => {
    onSubmit({
      ...values,
      id: api?.id,
    })
  }

  return (
    <SchemaFormDialog
      open={open}
      onOpenChange={onOpenChange}
      title={t('features.system.api.form.createTitle')}
      editTitle={t('features.system.api.form.editTitle')}
      description={
        isEdit
          ? t('features.system.api.form.editDescription')
          : t('features.system.api.form.createDescription')
      }
      schema={apiFormSchema}
      fields={fields}
      defaultValues={defaultValues}
      isEdit={isEdit}
      onSubmit={handleSubmit}
      isLoading={isLoading}
      maxWidth="md"
      submitText={isEdit ? t('features.system.api.form.save') : t('features.system.api.form.create')}
      cancelText={t('features.system.api.form.cancel')}
    />
  )
}

export default ApiFormDialog
