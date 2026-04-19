import { useEffect, useMemo, memo, useCallback } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import type * as z from 'zod'
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
  FormDescription,
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
import { Checkbox } from '@/components/ui/checkbox'
import { Button } from '@/components/ui/button'
import { PasswordInput } from '@/components/password-input'
import { cn } from '@/lib/utils'
import { ColorPicker } from './color-picker'

/**
 * Section卡片组件
 * 类似ProjectFormDialog中的ProjectSection
 */
function FormSectionCard({
  title,
  description,
  children,
}: {
  title: string
  description?: string
  children: React.ReactNode
}) {
  return (
    <section className='py-2'>
      <div className='mb-4 space-y-1'>
        <h3 className='text-sm font-semibold tracking-tight text-foreground'>{title}</h3>
        {description && <p className='text-sm text-muted-foreground'>{description}</p>}
      </div>
      <div className='space-y-4'>{children}</div>
    </section>
  )
}

/** 表单字段类型 */
export type FieldType = 'text' | 'password' | 'number' | 'textarea' | 'select' | 'checkbox' | 'checkbox-group' | 'color'

/** 选项类型 */
export interface FieldOption {
  label: string
  value: string | number
}

/** 表单字段配置 */
export interface FormFieldConfig {
  /** 字段名 */
  name: string
  /** 字段类型 */
  type: FieldType
  /** 标签 */
  label: string
  /** 是否必填 */
  required?: boolean
  /** 占位符 */
  placeholder?: string
  /** 描述 */
  description?: string
  /** 选项（用于 select/checkbox-group） */
  options?: FieldOption[]
  /** 自定义渲染 */
  render?: (props: { field: any; form: any }) => React.ReactNode
  /** 输入框类型（用于 text） */
  inputType?: 'text' | 'password' | 'email' | 'tel'
  /** 是否禁用 */
  disabled?: boolean
  /** 自定义类名 */
  className?: string
  /** 所属分组 key */
  section?: string
}

/** 表单分组配置 */
export interface FormSectionConfig {
  /** 分组标识 */
  key: string
  /** 分组标题 */
  title: string
  /** 分组描述 */
  description?: string
  /** 默认是否展开（已废弃，现在默认全部展开） */
  defaultOpen?: boolean
}

export interface SchemaFormDialogProps {
  /** 是否打开 */
  open: boolean
  /** 打开状态变化回调 */
  onOpenChange: (open: boolean) => void
  /** 标题 */
  title: string
  /** 编辑模式标题 */
  editTitle?: string
  /** 描述 */
  description?: string
  /** 是否编辑模式 */
  isEdit?: boolean
  /** Zod Schema */
  schema: z.ZodObject<any>
  /** 字段配置 */
  fields: FormFieldConfig[]
  /** 默认值 */
  defaultValues?: Record<string, any>
  /** 提交回调 */
  onSubmit: (values: any) => void
  /** 是否加载中 */
  isLoading?: boolean
  /** 对话框最大宽度 */
  maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
  /** 提交按钮文本 */
  submitText?: string
  /** 取消按钮文本 */
  cancelText?: string
  /** 分组配置 */
  sections?: FormSectionConfig[]
}

const maxWidthMap = {
  sm: 'sm:max-w-sm',
  md: 'sm:max-w-md',
  lg: 'sm:max-w-lg',
  xl: 'sm:max-w-xl',
  full: 'sm:max-w-full',
}

/**
 * 配置化表单对话框组件 - 垂直紧凑布局
 * 
 * 通过配置生成表单，无需编写重复模板代码
 * 
 * @example
 * ```tsx
 * const fields: FormFieldConfig[] = [
 *   { name: 'username', type: 'text', label: '用户名', required: true },
 *   { name: 'status', type: 'select', label: '状态', options: [...] },
 *   { name: 'color', type: 'color', label: '颜色' },
 * ]
 * 
 * <SchemaFormDialog
 *   open={open}
 *   onOpenChange={setOpen}
 *   title="新建用户"
 *   editTitle="编辑用户"
 *   schema={userSchema}
 *   fields={fields}
 *   defaultValues={editingUser}
 *   isEdit={!!editingUser}
 *   onSubmit={handleSubmit}
 * />
 * ```
 */
export const SchemaFormDialog = memo(function SchemaFormDialog({
  open,
  onOpenChange,
  title,
  editTitle,
  description,
  isEdit = false,
  schema,
  fields,
  defaultValues = {},
  onSubmit,
  isLoading = false,
  maxWidth = 'md',
  submitText,
  cancelText,
  sections,
}: SchemaFormDialogProps) {
  // 提取默认值
  const initialValues = useMemo(() => {
    const values: Record<string, any> = {}
    fields.forEach((field) => {
      if (field.type === 'checkbox-group') {
        values[field.name] = defaultValues[field.name] || []
      } else if (field.type === 'checkbox') {
        values[field.name] = defaultValues[field.name] || false
      } else if (field.type === 'number') {
        values[field.name] = defaultValues[field.name] ?? ''
      } else {
        values[field.name] = defaultValues[field.name] || ''
      }
    })
    return values
  }, [fields, defaultValues])

  const form = useForm({
    resolver: zodResolver(schema),
    defaultValues: initialValues,
  })

  // 打开/编辑时重置表单
  useEffect(() => {
    if (open) {
      form.reset(initialValues)
    }
  }, [open, initialValues, form])

  const handleSubmit = useCallback(
    (values: any) => {
      onSubmit(values)
    },
    [onSubmit]
  )

  const displayTitle = isEdit && editTitle ? editTitle : title
  // 按钮文本
  const displaySubmitText = submitText || (isEdit ? '保存' : '创建')
  const displayCancelText = cancelText || '取消'

  // 分组字段
  const groupedFields = useMemo(() => {
    if (!sections || sections.length === 0) {
      return { ungrouped: fields }
    }
    const groups: Record<string, FormFieldConfig[]> = {}
    const ungrouped: FormFieldConfig[] = []
    
    // 初始化分组
    sections.forEach((section) => {
      groups[section.key] = []
    })
    
    // 分配字段
    fields.forEach((field) => {
      if (field.section && groups[field.section]) {
        groups[field.section].push(field)
      } else {
        ungrouped.push(field)
      }
    })
    
    return { groups, ungrouped }
  }, [fields, sections])

  // 渲染单个字段
  const renderField = useCallback(
    (fieldConfig: FormFieldConfig) => {
      const { name, type, label, required, placeholder, description, options, render, inputType, disabled, className } = fieldConfig

      return (
        <FormField
          key={name}
          control={form.control}
          name={name}
          render={({ field }) => (
            <FormItem className={cn('space-y-2', className)}>
              <FormLabel className="text-sm font-medium">
                {label}
                {required && <span className="text-destructive ml-1">*</span>}
              </FormLabel>
              <FormControl>
                {render ? (
                  render({ field, form })
                ) : type === 'password' ? (
                  <PasswordInput
                    placeholder={placeholder}
                    disabled={disabled}
                    {...field}
                  />
                ) : type === 'text' ? (
                  <Input
                    type={inputType || 'text'}
                    placeholder={placeholder}
                    disabled={disabled}
                    {...field}
                  />
                ) : type === 'number' ? (
                  <Input
                    type="number"
                    placeholder={placeholder}
                    disabled={disabled}
                    onChange={(e) => field.onChange(e.target.valueAsNumber || '')}
                    value={field.value}
                  />
                ) : type === 'textarea' ? (
                  <Textarea
                    placeholder={placeholder}
                    disabled={disabled}
                    rows={3}
                    className="resize-none"
                    {...field}
                  />
                ) : type === 'select' && options ? (
                  <Select
                    value={String(field.value)}
                    onValueChange={field.onChange}
                    disabled={disabled}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder={placeholder} />
                    </SelectTrigger>
                    <SelectContent>
                      {options.map((opt) => (
                        <SelectItem key={opt.value} value={String(opt.value)}>
                          {opt.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                ) : type === 'checkbox' ? (
                  <Checkbox
                    checked={field.value}
                    onCheckedChange={field.onChange}
                    disabled={disabled}
                  />
                ) : type === 'checkbox-group' && options ? (
                  <div className="flex flex-wrap gap-x-6 gap-y-2">
                    {options.map((opt) => (
                      <div key={opt.value} className="flex items-center space-x-2">
                        <Checkbox
                          checked={field.value?.includes(opt.value)}
                          onCheckedChange={(checked) => {
                            const current = field.value || []
                            if (checked) {
                              field.onChange([...current, opt.value])
                            } else {
                              field.onChange(current.filter((v: any) => v !== opt.value))
                            }
                          }}
                          disabled={disabled}
                        />
                        <span className="text-sm">{opt.label}</span>
                      </div>
                    ))}
                  </div>
                ) : type === 'color' ? (
                  <ColorPicker
                    value={field.value || '#3b82f6'}
                    onChange={field.onChange}
                    showPresets
                    showInput
                  />
                ) : null}
              </FormControl>
              {description && type !== 'checkbox-group' && <FormDescription className="text-sm text-muted-foreground">{description}</FormDescription>}
              <FormMessage className="text-xs" />
            </FormItem>
          )}
        />
      )
    },
    [form]
  )

  return (
    <Dialog
      open={open}
      onOpenChange={(nextOpen) => {
        if (isLoading) return
        onOpenChange(nextOpen)
      }}
    >
      <DialogContent className={cn(maxWidthMap[maxWidth], 'max-h-[90vh] flex flex-col')}>
        <DialogHeader className='shrink-0'>
          <DialogTitle>{displayTitle}</DialogTitle>
          {description && <DialogDescription>{description}</DialogDescription>}
        </DialogHeader>

        <Form {...form}>
          <form 
            onSubmit={form.handleSubmit(handleSubmit)} 
            className='flex-1 overflow-y-auto pr-2 space-y-5 min-h-0'
          >
            {/* 分组字段 */}
            {sections && sections.length > 0 && (
              <div className="space-y-5">
                {sections.map((section) => {
                  const sectionFields = groupedFields.groups?.[section.key] || []
                  if (sectionFields.length === 0) return null
                  return (
                    <FormSectionCard
                      key={section.key}
                      title={section.title}
                      description={section.description}
                    >
                      <div className="space-y-4">
                        {sectionFields.map(renderField)}
                      </div>
                    </FormSectionCard>
                  )
                })}
              </div>
            )}
            
            {/* 未分组字段 - 也使用Section卡片包裹 */}
            {groupedFields.ungrouped && groupedFields.ungrouped.length > 0 && (
              <FormSectionCard title="基本信息">
                <div className="space-y-4">
                  {groupedFields.ungrouped.map(renderField)}
                </div>
              </FormSectionCard>
            )}
          </form>
        </Form>
        
        <DialogFooter className='shrink-0 gap-2 pt-4 border-t mt-4'>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isLoading}
          >
            {displayCancelText}
          </Button>
          <Button type="submit" disabled={isLoading} onClick={form.handleSubmit(handleSubmit)}>
            {isLoading ? (
              <span className="flex items-center gap-2">
                <span className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                {displaySubmitText}
              </span>
            ) : (
              submitText
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
})

export default SchemaFormDialog
