import { useMemo } from 'react'
import * as z from 'zod'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { SchemaFormDialog, type FormFieldConfig, type FormSectionConfig } from '@/components/schema-form-dialog'
import { IconPicker } from '@/components/icon-picker'
import type { Menu } from '../types/menu'
import { useI18n } from '@/context/i18n-provider'

const createMenuFormSchema = (t: (key: string) => string, isEdit: boolean) =>
  z.object({
    title: z.string().min(1, t('features.system.menu.form.validation.titleRequired')),
    name: z.string().min(1, t('features.system.menu.form.validation.nameRequired')),
    icon: z.string().optional(),
    path: z.string().min(1, t('features.system.menu.form.validation.pathRequired')),
    component: isEdit
      ? z.string().optional()
      : z.string().min(1, t('features.system.menu.form.validation.componentRequired')),
    redirect: z.string().optional(),
    sort: z.number().int().min(1, t('features.system.menu.form.validation.sortRequired')),
    status: z.number().min(1).max(2),
    hidden: z.number().min(1).max(2),
    noCache: z.number().min(1).max(2),
    activeMenu: z.string().optional(),
    parentID: z.number().optional(),
  })

type MenuFormValues = z.infer<ReturnType<typeof createMenuFormSchema>>

type MenuFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  menu: Menu | null
  onSubmit: (data: Partial<Menu>) => void
  isLoading?: boolean
  parentMenus?: Array<{ id: number; title: string; parentID: number | null }>
}

/**
 * 菜单表单对话框 - 使用分组折叠布局
 */
export function MenuFormDialog({
  open,
  onOpenChange,
  menu,
  onSubmit,
  isLoading = false,
  parentMenus = [],
}: MenuFormDialogProps) {
  const { t } = useI18n()
  const isEdit = !!menu

  const menuFormSchema = useMemo(() => createMenuFormSchema(t, isEdit), [t, isEdit])

  // 分组配置
  const sections: FormSectionConfig[] = useMemo(
    () => [
      { key: 'basic', title: t('features.system.menu.sections.basic') || '基本信息', defaultOpen: true },
      { key: 'display', title: t('features.system.menu.sections.display') || '显示设置', defaultOpen: true },
      { key: 'advanced', title: t('features.system.menu.sections.advanced') || '高级设置', defaultOpen: false },
    ],
    [t]
  )

  // 字段配置
  const fields: FormFieldConfig[] = useMemo(
    () => [
      // 基本信息
      {
        name: 'title',
        type: 'text',
        label: t('features.system.menu.menuTitle'),
        required: true,
        placeholder: t('features.system.menu.form.titlePlaceholder'),
        section: 'basic',
      },
      {
        name: 'name',
        type: 'text',
        label: t('features.system.menu.name'),
        required: true,
        placeholder: t('features.system.menu.form.namePlaceholder'),
        section: 'basic',
      },
      {
        name: 'path',
        type: 'text',
        label: t('features.system.menu.path'),
        required: true,
        placeholder: t('features.system.menu.form.pathPlaceholder'),
        section: 'basic',
      },
      {
        name: 'component',
        type: 'text',
        label: t('features.system.menu.component'),
        placeholder: t('features.system.menu.form.componentPlaceholder'),
        description: t('features.system.menu.form.componentDescription'),
        section: 'basic',
      },
      // 显示设置
      {
        name: 'icon',
        type: 'text',
        label: t('features.system.menu.icon'),
        render: ({ field }) => (
          <IconPicker
            value={field.value}
            onValueChange={field.onChange}
            placeholder={t('features.system.menu.form.iconPlaceholder')}
          />
        ),
        description: t('features.system.menu.form.iconDescription'),
        section: 'display',
      },
      {
        name: 'sort',
        type: 'number',
        label: t('features.system.menu.sort'),
        required: true,
        placeholder: t('features.system.menu.form.sortPlaceholder'),
        section: 'display',
      },
      {
        name: 'status',
        type: 'select',
        label: t('features.system.menu.status'),
        required: true,
        options: [
          { label: t('features.system.menu.enabled'), value: 1 },
          { label: t('features.system.menu.disabled'), value: 2 },
        ],
        section: 'display',
      },
      {
        name: 'hidden',
        type: 'select',
        label: t('features.system.menu.hidden'),
        required: true,
        options: [
          { label: t('features.system.menu.visible'), value: 2 },
          { label: t('features.system.menu.hidden'), value: 1 },
        ],
        section: 'display',
      },
      {
        name: 'noCache',
        type: 'select',
        label: t('features.system.menu.cache'),
        required: true,
        options: [
          { label: t('features.system.menu.cached'), value: 2 },
          { label: t('features.system.menu.noCache'), value: 1 },
        ],
        section: 'display',
      },
      // 高级设置
      {
        name: 'parentID',
        type: 'select',
        label: t('features.system.menu.parentMenu'),
        description: t('features.system.menu.form.parentMenuDescription'),
        section: 'advanced',
        render: ({ field }) => (
          <Select
            value={field.value ? String(field.value) : 'none'}
            onValueChange={(value) => field.onChange(value === 'none' ? undefined : Number(value))}
          >
            <SelectTrigger>
              <SelectValue placeholder={t('features.system.menu.form.parentMenuPlaceholder')} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="none">{t('features.system.menu.form.noParentMenu')}</SelectItem>
              {parentMenus.map((p) => (
                <SelectItem key={p.id} value={String(p.id)}>
                  {p.title}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        ),
      },
      {
        name: 'redirect',
        type: 'text',
        label: t('features.system.menu.redirect'),
        placeholder: t('features.system.menu.form.redirectPlaceholder'),
        description: t('features.system.menu.form.redirectDescription'),
        section: 'advanced',
      },
      {
        name: 'activeMenu',
        type: 'text',
        label: t('features.system.menu.activeMenu'),
        placeholder: t('features.system.menu.form.activeMenuPlaceholder'),
        description: t('features.system.menu.form.activeMenuDescription'),
        section: 'advanced',
      },
    ],
    [t, parentMenus]
  )

  // 默认值
  const defaultValues = useMemo(
    () =>
      isEdit && menu
        ? {
            title: menu.title,
            name: menu.name,
            icon: menu.icon || '',
            path: menu.path,
            component: menu.component || '',
            redirect: menu.redirect || '',
            sort: menu.sort,
            status: menu.status,
            hidden: menu.hidden,
            noCache: menu.noCache,
            activeMenu: menu.activeMenu || '',
            parentID: menu.parentID,
          }
        : {
            title: '',
            name: '',
            icon: '',
            path: '',
            component: '',
            redirect: '',
            sort: 1,
            status: 1,
            hidden: 2,
            noCache: 2,
            activeMenu: '',
            parentID: undefined,
          },
    [isEdit, menu]
  )

  const handleSubmit = (values: MenuFormValues) => {
    onSubmit({
      ...values,
      parentID: values.parentID || 0,
    })
  }

  return (
    <SchemaFormDialog
      open={open}
      onOpenChange={onOpenChange}
      title={t('features.system.menu.form.createTitle')}
      editTitle={t('features.system.menu.form.editTitle')}
      description={
        isEdit
          ? t('features.system.menu.form.editDescription')
          : t('features.system.menu.form.createDescription')
      }
      schema={menuFormSchema}
      fields={fields}
      sections={sections}
      defaultValues={defaultValues}
      isEdit={isEdit}
      onSubmit={handleSubmit}
      isLoading={isLoading}
      maxWidth="md"
      submitText={isEdit ? t('features.system.menu.form.save') : t('features.system.menu.form.create')}
      cancelText={t('features.system.menu.form.cancel')}
    />
  )
}

export default MenuFormDialog
