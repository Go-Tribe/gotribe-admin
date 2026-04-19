import { useMemo } from 'react'
import * as z from 'zod'
import { Checkbox } from '@/components/ui/checkbox'
import { FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { SchemaFormDialog, type FormFieldConfig } from '@/components/schema-form-dialog'
import type { Admin, Role } from '../types/admin'
import { useI18n } from '@/context/i18n-provider'

const createAdminFormSchema = (t: (key: string) => string) =>
  z.object({
    username: z.string().min(1, t('features.system.admin.form.validation.usernameRequired')),
    nickname: z.string().min(1, t('features.system.admin.form.validation.nicknameRequired')),
    mobile: z.string().regex(/^1[3-9]\d{9}$/, t('features.system.admin.form.validation.mobileInvalid')),
    introduction: z.string().optional(),
    status: z.number().min(1).max(2),
    roleIds: z.array(z.number()).min(1, t('features.system.admin.form.validation.rolesRequired')),
  })

type AdminFormValues = z.infer<ReturnType<typeof createAdminFormSchema>>

type AdminFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  admin: Admin | null
  onSubmit: (data: Partial<Admin>) => void
  isLoading?: boolean
  roles?: Role[]
}

/**
 * 管理员表单对话框 - 使用 SchemaFormDialog 通用组件
 */
export function AdminFormDialog({
  open,
  onOpenChange,
  admin,
  onSubmit,
  isLoading = false,
  roles = [],
}: AdminFormDialogProps) {
  const { t } = useI18n()
  const isEdit = !!admin

  const adminFormSchema = useMemo(() => createAdminFormSchema(t), [t])

  // 字段配置
  const fields: FormFieldConfig[] = useMemo(
    () => [
      {
        name: 'username',
        type: 'text',
        label: t('features.system.admin.username'),
        required: true,
        placeholder: t('features.system.admin.form.usernamePlaceholder'),
      },
      {
        name: 'nickname',
        type: 'text',
        label: t('features.system.admin.nickname'),
        required: true,
        placeholder: t('features.system.admin.form.nicknamePlaceholder'),
      },
      {
        name: 'mobile',
        type: 'text',
        label: t('features.system.admin.mobile'),
        required: true,
        placeholder: t('features.system.admin.form.mobilePlaceholder'),
        description: t('features.system.admin.form.mobileDescription'),
      },
      {
        name: 'status',
        type: 'select',
        label: t('features.system.admin.status'),
        required: true,
        options: [
          { label: t('features.system.admin.enabled'), value: 1 },
          { label: t('features.system.admin.disabled'), value: 2 },
        ],
      },
      {
        name: 'roleIds',
        type: 'checkbox-group',
        label: t('features.system.admin.form.roles'),
        required: true,
        description: t('features.system.admin.form.rolesDescription'),
        options: roles.map((role) => ({
          label: role.desc ? `${role.name} (${role.desc})` : role.name,
          value: role.id,
        })),
        render: ({ form }) => (
          <div className="space-y-1">
            {roles.length === 0 ? (
              <div className="text-sm text-muted-foreground pt-2">
                {t('features.system.admin.form.noRoles')}
              </div>
            ) : (
              <div className="space-y-3 pt-2">
                {roles.map((role) => (
                  <FormField
                    key={role.id}
                    control={form.control}
                    name="roleIds"
                    render={({ field }) => (
                      <FormItem className="flex flex-row items-start space-x-3 space-y-0">
                        <FormControl>
                          <Checkbox
                            checked={field.value?.includes(role.id)}
                            onCheckedChange={(checked) => {
                              const currentValue = field.value || []
                              return checked
                                ? field.onChange([...currentValue, role.id])
                                : field.onChange(currentValue.filter((id: number) => id !== role.id))
                            }}
                            disabled={isLoading}
                          />
                        </FormControl>
                        <FormLabel className="font-normal cursor-pointer">
                          <div className="flex items-center gap-2">
                            <span>{role.name}</span>
                            {role.desc && (
                              <span className="text-xs text-muted-foreground">({role.desc})</span>
                            )}
                          </div>
                        </FormLabel>
                      </FormItem>
                    )}
                  />
                ))}
              </div>
            )}
            <FormDescription className="text-xs mt-2">
              {t('features.system.admin.form.rolesDescription')}
            </FormDescription>
            <FormMessage />
          </div>
        ),
      },
      {
        name: 'introduction',
        type: 'textarea',
        label: t('features.system.admin.introduction'),
        placeholder: t('features.system.admin.form.introductionPlaceholder'),
      },
    ],
    [t, roles, isLoading]
  )

  // 默认值
  const defaultValues = useMemo(
    () =>
      isEdit && admin
        ? {
            username: admin.username,
            nickname: admin.nickname,
            mobile: admin.mobile,
            introduction: admin.introduction,
            status: admin.status,
            roleIds: admin.roleIds || [],
          }
        : {
            username: '',
            nickname: '',
            mobile: '',
            introduction: '',
            status: 1,
            roleIds: [],
          },
    [isEdit, admin]
  )

  const handleSubmit = (values: AdminFormValues) => {
    onSubmit({
      ...values,
      id: admin?.id,
      avatar: admin?.avatar || '',
      roleIds: values.roleIds || [],
    })
  }

  return (
    <SchemaFormDialog
      open={open}
      onOpenChange={onOpenChange}
      title={t('features.system.admin.form.createTitle')}
      editTitle={t('features.system.admin.form.editTitle')}
      description={
        isEdit
          ? t('features.system.admin.form.editDescription')
          : t('features.system.admin.form.createDescription')
      }
      schema={adminFormSchema}
      fields={fields}
      defaultValues={defaultValues}
      isEdit={isEdit}
      onSubmit={handleSubmit}
      isLoading={isLoading}
      maxWidth="md"
      submitText={isEdit ? t('features.system.admin.form.save') : t('features.system.admin.form.create')}
      cancelText={t('features.system.admin.form.cancel')}
    />
  )
}

export default AdminFormDialog
