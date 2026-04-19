import { useState, useEffect, useMemo } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { useQuery } from '@tanstack/react-query'
import { Loader2, ImageIcon } from 'lucide-react'
import { SubTitle } from '@/components/sub-title'
import { useI18n } from '@/context/i18n-provider'
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
import { Button } from '@/components/ui/button'
import { ResourceUpload, type ResourceItem } from '@/components/resource-upload'
import { getConfig, updateConfig } from './service/config'
import { useCrudMutations } from '@/hooks/use-crud-mutations'
import type { SystemConfig } from './types/config'

const createConfigFormSchema = (t: (key: string) => string) =>
  z.object({
    title: z.string().min(1, t('features.system.config.form.validation.titleRequired')),
    logo: z.string().min(1, t('features.system.config.form.validation.logoRequired')),
    icon: z.string().min(1, t('features.system.config.form.validation.iconRequired')),
  })

type ConfigFormValues = z.infer<ReturnType<typeof createConfigFormSchema>>

export function SystemConfig() {
  const { t } = useI18n()
  const [logoDialogOpen, setLogoDialogOpen] = useState(false)
  const [iconDialogOpen, setIconDialogOpen] = useState(false)

  const configFormSchema = useMemo(() => createConfigFormSchema(t), [t])

  // 获取系统配置
  const { data, isLoading } = useQuery({
    queryKey: ['systemConfig'],
    queryFn: getConfig,
  })

  const form = useForm<ConfigFormValues>({
    resolver: zodResolver(configFormSchema),
    defaultValues: {
      title: '',
      logo: '',
      icon: '',
    },
  })

  // 数据加载后填充表单
  useEffect(() => {
    if (data?.systemConfig) {
      form.reset({
        title: data.systemConfig.title || '',
        logo: data.systemConfig.logo || '',
        icon: data.systemConfig.icon || '',
      })
    }
  }, [data, form])

  // 使用统一的 CRUD mutation（仅更新操作）
  const { updateMutation } = useCrudMutations<Partial<SystemConfig>, string>({
    queryKey: ['systemConfig'],
    createFn: async () => ({}),
    updateFn: updateConfig,
    deleteFn: async () => ({}),
    messages: {
      updateSuccess: t('features.system.config.form.saveSuccess'),
    },
  })

  // 提交表单
  const handleSubmit = (values: ConfigFormValues) => {
    updateMutation.mutate(values)
  }

  // 选择 Logo
  const handleLogoSelect = (resource: ResourceItem) => {
    form.setValue('logo', resource.url, { shouldValidate: true })
  }

  // 选择图标
  const handleIconSelect = (resource: ResourceItem) => {
    form.setValue('icon', resource.url, { shouldValidate: true })
  }

  if (isLoading) {
    return (
      <div className='flex h-[400px] items-center justify-center'>
        <Loader2 className='h-8 w-8 animate-spin text-muted-foreground' />
      </div>
    )
  }

  return (
    <div className='space-y-4'>
      <div className='flex items-center justify-between px-4 pt-4'>
        <SubTitle
          title={t('features.system.config.title')}
          description={t('features.system.config.description')}
        />
      </div>
      <div className='rounded-md border p-6 mx-4'>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className='space-y-6'>
            <FormField
              control={form.control}
              name='title'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('features.system.config.form.title')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('features.system.config.form.titlePlaceholder')} {...field} />
                  </FormControl>
                  <FormDescription>{t('features.system.config.form.titleDescription')}</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name='logo'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('features.system.config.form.logo')}</FormLabel>
                  <FormControl>
                    <button
                      type='button'
                      onClick={() => setLogoDialogOpen(true)}
                      className='flex h-20 w-40 items-center justify-center overflow-hidden rounded-lg border-2 border-dashed border-muted-foreground/30 bg-muted/50 transition-colors hover:border-primary/50 hover:bg-muted focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2'
                    >
                      {field.value ? (
                        <img
                          src={field.value}
                          alt={t('features.system.config.form.logoAlt')}
                          className='h-full w-full object-contain'
                        />
                      ) : (
                        <ImageIcon className='h-8 w-8 text-muted-foreground' />
                      )}
                    </button>
                  </FormControl>
                  <FormDescription>{t('features.system.config.form.logoDescription')}</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name='icon'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('features.system.config.form.icon')}</FormLabel>
                  <FormControl>
                    <button
                      type='button'
                      onClick={() => setIconDialogOpen(true)}
                      className='flex h-16 w-16 items-center justify-center overflow-hidden rounded-lg border-2 border-dashed border-muted-foreground/30 bg-muted/50 transition-colors hover:border-primary/50 hover:bg-muted focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2'
                    >
                      {field.value ? (
                        <img
                          src={field.value}
                          alt={t('features.system.config.form.iconAlt')}
                          className='h-full w-full object-contain'
                        />
                      ) : (
                        <ImageIcon className='h-6 w-6 text-muted-foreground' />
                      )}
                    </button>
                  </FormControl>
                  <FormDescription>{t('features.system.config.form.iconDescription')}</FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className='pt-4'>
              <Button type='submit' disabled={updateMutation.isPending}>
                {updateMutation.isPending && <Loader2 className='mr-2 h-4 w-4 animate-spin' />}
                {t('features.system.config.form.save')}
              </Button>
            </div>
          </form>
        </Form>
      </div>
      <ResourceUpload
        open={logoDialogOpen}
        onOpenChange={setLogoDialogOpen}
        onSelect={handleLogoSelect}
        type={1}
        title={t('features.system.config.dialog.logoTitle')}
        description={t('features.system.config.dialog.logoDescription')}
      />

      <ResourceUpload
        open={iconDialogOpen}
        onOpenChange={setIconDialogOpen}
        onSelect={handleIconSelect}
        type={1}
        title={t('features.system.config.dialog.iconTitle')}
        description={t('features.system.config.dialog.iconDescription')}
      />
    </div>
  )
}
// MIGRATED
