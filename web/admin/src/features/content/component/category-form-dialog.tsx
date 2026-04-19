import { useEffect, useMemo, useState } from 'react'
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
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { ResourceUpload, type ResourceItem } from '@/components/resource-upload'
import { useI18n } from '@/context/i18n-provider'
import type { Category, CategoryParams } from '../types/category'

const createCategoryFormSchema = (t: (key: string) => string) =>
  z.object({
    title: z.string().min(1, t('features.content.category.form.validation.titleRequired')),
    slug: z
      .string()
      .min(1, t('features.content.category.form.validation.slugRequired'))
      .max(30, t('features.content.category.form.validation.slugMax'))
      .regex(/^[a-zA-Z0-9-_]+$/, t('features.content.category.form.validation.slugPattern')),
    description: z.string().optional(),
    icon: z.string().optional(),
    path: z.string().optional(),
    sort: z.number().int().min(1).max(999),
    hidden: z.number().min(1).max(2),
    parentId: z.number().int().min(0),
  })

type CategoryFormValues = z.infer<ReturnType<typeof createCategoryFormSchema>>

/** 递归扁平化分类树为 { id, title } 列表（带层级缩进，参考菜单管理） */
function flattenCategoriesWithLevel(
  nodes: Category[],
  level = 0,
  result: { id: number; title: string }[] = [],
): { id: number; title: string }[] {
  for (const node of nodes) {
    result.push({ id: node.id, title: '  '.repeat(level) + node.title })
    if (node.children?.length) flattenCategoriesWithLevel(node.children, level + 1, result)
  }
  return result
}

type CategoryFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  category: Category | null
  /** 分类树，用于父级选择（新建时） */
  categoryTree?: Category[]
  onSubmit: (data: CategoryParams) => void
  isLoading?: boolean
}

export function CategoryFormDialog({
  open,
  onOpenChange,
  category,
  categoryTree = [],
  onSubmit,
  isLoading = false,
}: CategoryFormDialogProps) {
  const { t } = useI18n()
  const isEdit = !!category
  const [sortInputValue, setSortInputValue] = useState<string>('')
  const [iconResourceDialogOpen, setIconResourceDialogOpen] = useState(false)
  const categoryFormSchema = useMemo(() => createCategoryFormSchema(t), [t])
  // 编辑时排除当前分类，避免选自己为父级
  const parentOptions = useMemo(() => {
    const list = flattenCategoriesWithLevel(categoryTree)
    if (isEdit && category) {
      return list.filter((opt) => opt.id !== category.id)
    }
    return list
  }, [categoryTree, isEdit, category])

  const form = useForm<CategoryFormValues>({
    resolver: zodResolver(categoryFormSchema),
    defaultValues: {
      title: '',
      slug: '',
      description: '',
      icon: '',
      path: '',
      sort: 1,
      hidden: 1,
      parentId: 0,
    },
  })

  useEffect(() => {
    if (open) {
      if (isEdit && category) {
        const sortValue = Math.max(1, Math.min(99, category.sort || 1))
        const rawParent = category.parentID ?? category.parentId
        const parentValue = typeof rawParent === 'number' ? rawParent : Number(rawParent) || 0
        form.reset({
          title: category.title,
          slug: category.slug || '',
          description: category.description || '',
          icon: category.icon || '',
          path: category.path ?? category.route ?? '',
          sort: sortValue,
          hidden: category.hidden || 1,
          parentId: parentValue,
        })
        setTimeout(() => setSortInputValue(String(sortValue)), 0)
      } else {
        form.reset({
          title: '',
          slug: '',
          description: '',
          icon: '',
          path: '',
          sort: 1,
          hidden: 1,
          parentId: 0,
        })
        setTimeout(() => setSortInputValue('1'), 0)
      }
    }
  }, [open, isEdit, category, form])

  const handleSubmit = (values: CategoryFormValues) => {
    onSubmit({
      title: values.title,
      slug: values.slug,
      description: values.description,
      icon: values.icon,
      path: values.path ?? '', // 链接，创建/更新接口用 path
      route: values.path ?? '',
      sort: values.sort,
      status: category?.status ?? 1,
      hidden: values.hidden,
      parentId: values.parentId,
      id: category?.id,
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-[600px] max-h-[90vh] flex flex-col'>
        <DialogHeader className='shrink-0'>
          <DialogTitle>{isEdit ? t('features.content.category.form.editTitle') : t('features.content.category.form.createTitle')}</DialogTitle>
          <DialogDescription>
            {isEdit ? t('features.content.category.form.editDescription') : t('features.content.category.form.createDescription')}
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
                  <FormLabel>{t('features.content.category.form.title')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('features.content.category.form.titlePlaceholder')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='slug'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.category.form.slug')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('features.content.category.form.slugPlaceholder')} {...field} />
                  </FormControl>
                  <FormDescription>{t('features.content.category.form.slugDescription')}</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='parentId'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.category.form.parent')}</FormLabel>
                  <Select
                    onValueChange={(value) => field.onChange(Number(value))}
                    value={String(field.value)}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder={t('features.content.category.form.parentPlaceholder')} />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value='0'>{t('features.content.category.form.rootCategory')}</SelectItem>
                      {parentOptions.map((opt) => (
                        <SelectItem key={opt.id} value={String(opt.id)}>
                          {opt.title}
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
              name='description'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.category.form.description')}</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder={t('features.content.category.form.descriptionPlaceholder')}
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
              name='icon'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.category.form.icon')}</FormLabel>
                  <FormControl>
                    <button
                      type='button'
                      onClick={() => setIconResourceDialogOpen(true)}
                      className={cn(
                        'flex flex-col items-center justify-center h-20 w-20 rounded-lg border-2 border-dashed border-border bg-muted/30 overflow-hidden transition-all shrink-0',
                        'cursor-pointer hover:border-primary/50 hover:bg-muted/50'
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
                            {t('features.content.category.form.selectIcon')}
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
              name='path'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.category.form.path')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('features.content.category.form.pathPlaceholder')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='sort'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.category.form.sort')}</FormLabel>
                  <FormControl>
                    <Input
                      type='number'
                      placeholder={t('features.content.category.form.sortPlaceholder')}
                      value={sortInputValue}
                      onChange={(e) => {
                        const value = e.target.value
                        setSortInputValue(value)
                      }}
                      onBlur={() => {
                        const numValue = Number(sortInputValue)
                        let finalValue = 1
                        if (!isNaN(numValue) && sortInputValue !== '') {
                          if (numValue > 999) {
                            finalValue = 999
                          } else if (numValue < 1) {
                            finalValue = 1
                          } else {
                            finalValue = numValue
                          }
                        }
                        setSortInputValue(String(finalValue))
                        field.onChange(finalValue)
                        field.onBlur()
                      }}
                      min={1}
                      max={999}
                    />
                  </FormControl>
                  <FormDescription>{t('features.content.category.form.sortDescription')}</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='hidden'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.content.category.form.hidden')}</FormLabel>
                  <Select
                    onValueChange={(value) => field.onChange(Number(value))}
                    value={String(field.value)}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder={t('features.content.category.form.hiddenPlaceholder')} />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value='1'>{t('features.content.category.hiddenStatus.show')}</SelectItem>
                      <SelectItem value='2'>{t('features.content.category.hiddenStatus.hidden')}</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />
          </form>
        </Form>
        <DialogFooter className='shrink-0 pt-4 border-t mt-4'>
          <Button
            type='button'
            variant='outline'
            onClick={() => onOpenChange(false)}
            disabled={isLoading}
          >
            {t('features.content.category.form.cancel')}
          </Button>
          <Button
            type='button'
            disabled={isLoading}
            onClick={form.handleSubmit(handleSubmit)}
          >
            {isLoading ? t('features.content.category.form.submitting') : isEdit ? t('features.content.category.form.save') : t('features.content.category.form.create')}
          </Button>
        </DialogFooter>
      </DialogContent>
      <ResourceUpload
        open={iconResourceDialogOpen}
        onOpenChange={setIconResourceDialogOpen}
        onSelect={(resource: ResourceItem) => {
          form.setValue('icon', resource.url, { shouldValidate: true })
          setIconResourceDialogOpen(false)
        }}
        type={1}
        title={t('features.content.category.form.selectIconDialogTitle')}
        description={t('features.content.category.form.selectIconDialogDesc')}
      />
    </Dialog>
  )
}
