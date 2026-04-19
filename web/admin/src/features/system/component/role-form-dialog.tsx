import { useMemo } from 'react'
import * as z from 'zod'
import { Input } from '@/components/ui/input'
import { SchemaFormDialog, type FormFieldConfig } from '@/components/schema-form-dialog'
import type { Role } from '../types/admin'
import { useI18n } from '@/context/i18n-provider'

const createRoleFormSchema = (t: (key: string) => string) =>
  z.object({
    name: z.string().min(1, t('features.system.role.form.validation.nameRequired')),
    keyword: z
      .string()
      .min(1, t('features.system.role.form.validation.keywordRequired'))
      .regex(/^[a-zA-Z0-9]+$/, t('features.system.role.form.validation.keywordInvalid')),
    sort: z.number().int().min(1).max(99),
    status: z.number().min(1).max(2),
    desc: z.string().optional(),
  })

type RoleFormValues = z.infer<ReturnType<typeof createRoleFormSchema>>

type RoleFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  role: Role | null
  onSubmit: (data: Partial<Role>) => void
  isLoading?: boolean
}

/**
 * 角色表单对话框 - 使用 SchemaFormDialog 通用组件
 */
export function RoleFormDialog({
  open,
  onOpenChange,
  role,
  onSubmit,
  isLoading = false,
}: RoleFormDialogProps) {
  const { t } = useI18n()
  const isEdit = !!role

  const roleFormSchema = useMemo(() => createRoleFormSchema(t), [t])

  // 字段配置
  const fields: FormFieldConfig[] = useMemo(
    () => [
      {
        name: 'name',
        type: 'text',
        label: t('features.system.role.name'),
        required: true,
        placeholder: t('features.system.role.form.namePlaceholder'),
      },
      {
        name: 'keyword',
        type: 'text',
        label: t('features.system.role.keyword'),
        required: true,
        placeholder: t('features.system.role.form.keywordPlaceholder'),
        description: t('features.system.role.form.keywordDescription') || '仅支持字母和数字',
        render: ({ field }) => (
          <Input
            placeholder={t('features.system.role.form.keywordPlaceholder')}
            {...field}
            onChange={(e) => {
              const value = e.target.value.replace(/[^a-zA-Z0-9]/g, '')
              field.onChange(value)
            }}
          />
        ),
      },
      {
        name: 'sort',
        type: 'number',
        label: t('features.system.role.sort'),
        required: true,
        placeholder: t('features.system.role.form.sortPlaceholder'),
        description: t('features.system.role.form.sortDescription'),
      },
      {
        name: 'status',
        type: 'select',
        label: t('features.system.role.status'),
        required: true,
        options: [
          { label: t('features.system.role.enabled'), value: 1 },
          { label: t('features.system.role.disabled'), value: 2 },
        ],
      },
      {
        name: 'desc',
        type: 'textarea',
        label: t('features.system.role.desc'),
        placeholder: t('features.system.role.form.descPlaceholder'),
      },
    ],
    [t]
  )

  // 默认值
  const defaultValues = useMemo(
    () =>
      isEdit && role
        ? {
            name: role.name,
            keyword: role.keyword,
            sort: Math.max(1, Math.min(99, role.sort || 1)),
            status: role.status,
            desc: role.desc,
          }
        : {
            name: '',
            keyword: '',
            sort: 1,
            status: 1,
            desc: '',
          },
    [isEdit, role]
  )

  const handleSubmit = (values: RoleFormValues) => {
    onSubmit({
      ...values,
      id: role?.id,
    })
  }

  return (
    <SchemaFormDialog
      open={open}
      onOpenChange={onOpenChange}
      title={t('features.system.role.form.createTitle')}
      editTitle={t('features.system.role.form.editTitle')}
      description={
        isEdit
          ? t('features.system.role.form.editDescription')
          : t('features.system.role.form.createDescription')
      }
      schema={roleFormSchema}
      fields={fields}
      defaultValues={defaultValues}
      isEdit={isEdit}
      onSubmit={handleSubmit}
      isLoading={isLoading}
      maxWidth="md"
      submitText={isEdit ? t('features.system.role.form.save') : t('features.system.role.form.create')}
      cancelText={t('features.system.role.form.cancel')}
    />
  )
}

export default RoleFormDialog
