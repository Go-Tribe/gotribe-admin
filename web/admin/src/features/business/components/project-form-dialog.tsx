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
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { ResourceUpload, type ResourceItem } from '@/components/resource-upload'
import type { Project } from '../types/project'
import { useI18n } from '@/context/i18n-provider'

const createProjectFormSchema = (t: (key: string) => string) =>
  z.object({
    title: z.string().min(1, t('features.business.project.form.validation.titleRequired')),
    description: z.string().optional(),
    name: z.string().min(1, t('features.business.project.form.validation.nameRequired')),
    metaDescription: z.string().optional(),
    keywords: z.string().optional(),
    domain: z.string().optional(),
    postUrl: z.string().optional(),
    icp: z.string().optional(),
    author: z.string().optional(),
    baiduAnalytics: z.string().optional(),
    favicon: z.string().optional(),
    publicSecurity: z.string().optional(),
    navImage: z.string().optional(),
  })

type ProjectFormValues = z.infer<ReturnType<typeof createProjectFormSchema>>

type ProjectFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  project: Project | null
  onSubmit: (data: Partial<Project>) => void
  isLoading?: boolean
}

function ProjectSection({
  title,
  description,
  children,
}: {
  title: string
  description: string
  children: React.ReactNode
}) {
  return (
    <section className='rounded-2xl border border-border/60 bg-muted/20 p-4 shadow-sm'>
      <div className='mb-4 space-y-1'>
        <h3 className='text-sm font-semibold tracking-tight text-foreground'>{title}</h3>
        <p className='text-sm text-muted-foreground'>{description}</p>
      </div>
      <div className='space-y-4'>{children}</div>
    </section>
  )
}

export function ProjectFormDialog({
  open,
  onOpenChange,
  project,
  onSubmit,
  isLoading = false,
}: ProjectFormDialogProps) {
  const { t } = useI18n()
  const isEdit = !!project
  const [iconResourceField, setIconResourceField] = useState<'favicon' | 'navImage' | null>(null)

  const projectFormSchema = useMemo(() => createProjectFormSchema(t), [t])

  const form = useForm<ProjectFormValues>({
    resolver: zodResolver(projectFormSchema),
    defaultValues: {
      title: '',
      description: '',
      name: '',
      metaDescription: '',
      keywords: '',
      domain: '',
      postUrl: '',
      icp: '',
      author: '',
      baiduAnalytics: '',
      favicon: '',
      publicSecurity: '',
      navImage: '',
    },
  })

  useEffect(() => {
    if (open) {
      if (isEdit && project) {
        form.reset({
          title: project.title || '',
          description: project.description || '',
          name: project.name || '',
          metaDescription: project.info || '',
          keywords: project.keywords || '',
          domain: project.domain || '',
          postUrl: project.postUrl || '',
          icp: project.icp || '',
          author: project.author || '',
          baiduAnalytics: project.baiduAnalytics || '',
          favicon: project.favicon || '',
          publicSecurity: project.publicSecurity || '',
          navImage: project.navImage || '',
        })
      } else {
        form.reset({
          title: '',
          description: '',
          name: '',
          metaDescription: '',
          keywords: '',
          domain: '',
          postUrl: '',
          icp: '',
          author: '',
          baiduAnalytics: '',
          favicon: '',
          publicSecurity: '',
          navImage: '',
        })
      }
    }
  }, [open, isEdit, project, form])

  const handleSubmit = (values: ProjectFormValues) => {
    onSubmit({
      ...values,
      id: project?.id,
      // 将 metaDescription 映射到 info 字段（如果 API 期望的是 info）
      info: values.metaDescription,
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-[700px] max-h-[90vh] flex flex-col'>
        <DialogHeader className='shrink-0'>
          <DialogTitle>
            {isEdit ? t('features.business.project.form.editTitle') : t('features.business.project.form.createTitle')}
          </DialogTitle>
          <DialogDescription>
            {isEdit
              ? t('features.business.project.form.editDescription')
              : t('features.business.project.form.createDescription')}
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className='flex-1 overflow-y-auto pr-2 space-y-5 min-h-0'
          >
            <div className='rounded-2xl border border-border/60 bg-card/80 p-4 shadow-sm'>
              <p className='text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground'>
                {t('features.business.project.form.summaryLabel')}
              </p>
              <div className='mt-2 flex flex-wrap gap-2 text-sm text-muted-foreground'>
                <span className='rounded-full bg-muted px-3 py-1'>
                  {t('features.business.project.form.summaryBasic')}
                </span>
                <span className='rounded-full bg-muted px-3 py-1'>
                  {t('features.business.project.form.summarySeo')}
                </span>
                <span className='rounded-full bg-muted px-3 py-1'>
                  {t('features.business.project.form.summaryBrand')}
                </span>
                <span className='rounded-full bg-muted px-3 py-1'>
                  {t('features.business.project.form.summaryCompliance')}
                </span>
              </div>
            </div>

            <ProjectSection
              title={t('features.business.project.form.sections.basicTitle')}
              description={t('features.business.project.form.sections.basicDescription')}
            >
              <div className='grid gap-4 md:grid-cols-2'>
                <FormField
                  control={form.control}
                  name='title'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('features.business.project.form.fields.title')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('features.business.project.form.fields.titlePlaceholder')} {...field} />
                      </FormControl>
                      <FormDescription>{t('features.business.project.form.hints.title')}</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name='name'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('features.business.project.form.fields.name')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('features.business.project.form.fields.namePlaceholder')} {...field} />
                      </FormControl>
                      <FormDescription>{t('features.business.project.form.hints.name')}</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <FormField
                control={form.control}
                name='description'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('features.business.project.form.fields.description')}</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder={t('features.business.project.form.fields.descriptionPlaceholder')}
                        className='min-h-28 resize-none'
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>{t('features.business.project.form.hints.description')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='author'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('features.business.project.form.fields.author')}</FormLabel>
                    <FormControl>
                      <Input placeholder={t('features.business.project.form.fields.authorPlaceholder')} {...field} />
                    </FormControl>
                    <FormDescription>{t('features.business.project.form.hints.author')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </ProjectSection>

            <ProjectSection
              title={t('features.business.project.form.sections.seoTitle')}
              description={t('features.business.project.form.sections.seoDescription')}
            >
              <div className='grid gap-4 md:grid-cols-2'>
                <FormField
                  control={form.control}
                  name='domain'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('features.business.project.form.fields.domain')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('features.business.project.form.fields.domainPlaceholder')} {...field} />
                      </FormControl>
                      <FormDescription>{t('features.business.project.form.hints.domain')}</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name='postUrl'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('features.business.project.form.fields.postUrl')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('features.business.project.form.fields.postUrlPlaceholder')} {...field} />
                      </FormControl>
                      <FormDescription>{t('features.business.project.form.hints.postUrl')}</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <FormField
                control={form.control}
                name='metaDescription'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('features.business.project.form.fields.metaDescription')}</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder={t('features.business.project.form.fields.metaDescriptionPlaceholder')}
                        className='min-h-24 resize-none'
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>{t('features.business.project.form.hints.metaDescription')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name='keywords'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('features.business.project.form.fields.keywords')}</FormLabel>
                    <FormControl>
                      <Input placeholder={t('features.business.project.form.fields.keywordsPlaceholder')} {...field} />
                    </FormControl>
                    <FormDescription>{t('features.business.project.form.hints.keywords')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </ProjectSection>

            <ProjectSection
              title={t('features.business.project.form.sections.brandTitle')}
              description={t('features.business.project.form.sections.brandDescription')}
            >
              <div className='grid gap-4 md:grid-cols-2'>
                <FormField
                  control={form.control}
                  name='favicon'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('features.business.project.form.fields.favicon')}</FormLabel>
                      <FormControl>
                        <button
                          type='button'
                          onClick={() => setIconResourceField('favicon')}
                          className={cn(
                            'flex min-h-36 w-full items-center gap-4 rounded-2xl border border-dashed border-border bg-background px-4 py-4 text-left transition-all',
                            'cursor-pointer hover:border-primary/50 hover:bg-muted/40'
                          )}
                        >
                          <div className='flex h-20 w-20 shrink-0 items-center justify-center overflow-hidden rounded-2xl border bg-muted/50'>
                            {field.value ? (
                              <img src={field.value} alt='' className='h-full w-full object-cover' />
                            ) : (
                              <ImageIcon className='h-8 w-8 text-muted-foreground' />
                            )}
                          </div>
                          <div className='space-y-1'>
                            <div className='font-medium text-foreground'>
                              {field.value
                                ? t('features.business.project.form.fields.replaceImage')
                                : t('features.business.project.form.fields.selectImage')}
                            </div>
                            <p className='text-sm text-muted-foreground'>
                              {t('features.business.project.form.hints.favicon')}
                            </p>
                          </div>
                        </button>
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name='navImage'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('features.business.project.form.fields.navImage')}</FormLabel>
                      <FormControl>
                        <button
                          type='button'
                          onClick={() => setIconResourceField('navImage')}
                          className={cn(
                            'flex min-h-36 w-full items-center gap-4 rounded-2xl border border-dashed border-border bg-background px-4 py-4 text-left transition-all',
                            'cursor-pointer hover:border-primary/50 hover:bg-muted/40'
                          )}
                        >
                          <div className='flex h-20 w-20 shrink-0 items-center justify-center overflow-hidden rounded-2xl border bg-muted/50'>
                            {field.value ? (
                              <img src={field.value} alt='' className='h-full w-full object-cover' />
                            ) : (
                              <ImageIcon className='h-8 w-8 text-muted-foreground' />
                            )}
                          </div>
                          <div className='space-y-1'>
                            <div className='font-medium text-foreground'>
                              {field.value
                                ? t('features.business.project.form.fields.replaceImage')
                                : t('features.business.project.form.fields.selectImage')}
                            </div>
                            <p className='text-sm text-muted-foreground'>
                              {t('features.business.project.form.hints.navImage')}
                            </p>
                          </div>
                        </button>
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </ProjectSection>

            <ProjectSection
              title={t('features.business.project.form.sections.complianceTitle')}
              description={t('features.business.project.form.sections.complianceDescription')}
            >
              <div className='grid gap-4 md:grid-cols-2'>
                <FormField
                  control={form.control}
                  name='icp'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('features.business.project.form.fields.icp')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('features.business.project.form.fields.icpPlaceholder')} {...field} />
                      </FormControl>
                      <FormDescription>{t('features.business.project.form.hints.icp')}</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name='publicSecurity'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('features.business.project.form.fields.publicSecurity')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('features.business.project.form.fields.publicSecurityPlaceholder')} {...field} />
                      </FormControl>
                      <FormDescription>{t('features.business.project.form.hints.publicSecurity')}</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <FormField
                control={form.control}
                name='baiduAnalytics'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('features.business.project.form.fields.baiduAnalytics')}</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder={t('features.business.project.form.fields.baiduAnalyticsPlaceholder')}
                        className='min-h-32 resize-none font-mono text-sm'
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>{t('features.business.project.form.hints.baiduAnalytics')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </ProjectSection>
          </form>
        </Form>
        <DialogFooter className='shrink-0 pt-4 border-t mt-4'>
          <Button
            type='button'
            variant='outline'
            onClick={() => onOpenChange(false)}
            disabled={isLoading}
          >
            {t('features.business.project.form.cancel')}
          </Button>
          <Button
            type='button'
            disabled={isLoading}
            onClick={form.handleSubmit(handleSubmit)}
          >
            {isLoading ? t('features.business.project.form.submitting') : isEdit ? t('features.business.project.form.save') : t('features.business.project.form.create')}
          </Button>
        </DialogFooter>
      </DialogContent>
      <ResourceUpload
        open={iconResourceField !== null}
        onOpenChange={(open) => !open && setIconResourceField(null)}
        onSelect={(resource: ResourceItem) => {
          if (iconResourceField) {
            form.setValue(iconResourceField, resource.url, { shouldValidate: true })
            setIconResourceField(null)
          }
        }}
        type={1}
        title={t('features.business.project.form.fields.selectImage')}
        description={t('features.business.project.form.fields.favicon')}
      />
    </Dialog>
  )
}
