import { useEffect, useMemo, useState, useCallback, useRef } from 'react'
import * as z from 'zod'
import { format } from 'date-fns'
import { useForm, type Resolver, type FieldErrors } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { getProjectList, getUserList } from '@/shared/api'
import { ArrowLeft } from 'lucide-react'
import { toast } from 'sonner'
import { slateContentToHtml } from '@/lib/slate-markdown'
import { generateSlug } from '@/lib/slug'
import { useI18n } from '@/context/i18n-provider'
import { Button } from '@/components/ui/button'
import { Form } from '@/components/ui/form'
import { ResourceUpload, type ResourceItem } from '@/components/resource-upload'
import { ArticleEditor } from './components/article-editor'
import { ArticleMediaSheet } from './components/article-media-sheet'
import { ArticleSettingsSheet } from './components/article-settings-sheet'
import { getCategoryTree } from './service/category'
import { createPost, updatePost, getPostDetail } from './service/post'
import { getTagList } from './service/tag'
import type { Category } from './types/category'
import type { Post, PostParams } from './types/post'

/** 根据点号路径写入嵌套对象，如 key="meta.title" value="x" => obj.meta.title = "x" */
function setExtNested(
  obj: Record<string, unknown>,
  path: string,
  value: string
) {
  const parts = path
    .trim()
    .split('.')
    .map((p) => p.trim())
    .filter(Boolean)
  if (parts.length === 0) return
  let current: Record<string, unknown> = obj
  for (let i = 0; i < parts.length - 1; i++) {
    const key = parts[i]!
    if (
      !(key in current) ||
      typeof current[key] !== 'object' ||
      current[key] === null ||
      Array.isArray(current[key])
    ) {
      current[key] = {}
    }
    current = current[key] as Record<string, unknown>
  }
  current[parts[parts.length - 1]!] = value
}

/**
 * 将嵌套对象展平为 key-value 列表，如 { meta: { title: "x" } } => [{ key: "meta.title", value: "x" }]。
 * 仅支持嵌套纯对象与字符串值；数组等会转为 String(v)，无法正确回显。
 */
function flattenExtToKeyValue(
  obj: Record<string, unknown>,
  prefix = ''
): Array<{ key: string; value: string }> {
  const result: Array<{ key: string; value: string }> = []
  for (const [k, v] of Object.entries(obj)) {
    const path = prefix ? `${prefix}.${k}` : k
    if (
      v !== null &&
      typeof v === 'object' &&
      !Array.isArray(v) &&
      Object.getPrototypeOf(v) === Object.prototype
    ) {
      result.push(...flattenExtToKeyValue(v as Record<string, unknown>, path))
    } else {
      result.push({ key: path, value: v != null ? String(v) : '' })
    }
  }
  return result
}

const createArticleFormSchema = (t: (key: string) => string) =>
  z.object({
    title: z
      .string()
      .min(1, t('features.content.article.form.validation.titleRequired')),
    slug: z.string().optional(),
    description: z
      .string()
      .min(
        1,
        t('features.content.article.form.validation.descriptionRequired')
      ),
    author: z
      .string()
      .min(1, t('features.content.article.form.validation.authorRequired')),
    userID: z
      .string()
      .min(1, t('features.content.article.form.validation.authorRequired')),
    content: z
      .string()
      .min(1, t('features.content.article.form.validation.contentRequired')),
    icon: z.string().optional(),
    video: z.string().optional(),
    images: z.array(z.string()).optional(),
    type: z.coerce
      .number()
      .min(1, t('features.content.article.form.validation.typeRequired')),
    status: z.coerce.number().default(1),
    categoryID: z
      .string()
      .min(1, t('features.content.article.form.validation.categoryRequired')),
    projectId: z
      .string()
      .min(1, t('features.content.article.form.validation.projectRequired')),
    isTop: z.coerce.number().optional(),
    isPasswd: z.coerce.number().optional(),
    password: z.string().optional(),
    tag: z.string().optional(),
    showTime: z.date().optional(),
  })

export type ArticleFormValues = z.infer<
  ReturnType<typeof createArticleFormSchema>
>

type ArticleFormPageProps = {
  id?: number | null
  initialPost?: Post | null
}

export function ArticleFormPage({ id, initialPost }: ArticleFormPageProps) {
  const { t } = useI18n()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const isEdit = !!id
  const articleFormSchema = useMemo(() => createArticleFormSchema(t), [t])

  const [resourceOpen, setResourceOpen] = useState(false)
  const [resourceType, setResourceType] = useState<1 | 2>(1)
  const [currentUploadField, setCurrentUploadField] = useState<
    'icon' | 'video' | 'images' | null
  >(null)
  /** 自定义字段 key-value 列表，最终序列化为 JSON 存到 ext */
  const [extFields, setExtFields] = useState<
    Array<{ key: string; value: string }>
  >([])
  const [sheetOpen, setSheetOpen] = useState(false)
  const [mediaSheetOpen, setMediaSheetOpen] = useState(false)

  const { data: categoryData } = useQuery({
    queryKey: ['categoryTree'],
    queryFn: getCategoryTree,
  })

  const { data: projectData } = useQuery({
    queryKey: ['projectList', { current: 1, pageNum: 1, pageSize: 1000 }],
    queryFn: () => getProjectList({ current: 1, pageNum: 1, pageSize: 1000 }),
  })
  const projectList = useMemo(
    () => projectData?.projects ?? [],
    [projectData?.projects]
  )

  const { data: userData } = useQuery({
    queryKey: ['userList', { current: 1, pageNum: 1, pageSize: 1000 }],
    queryFn: () => getUserList({ current: 1, pageNum: 1, pageSize: 1000 }),
  })
  const userList = useMemo(() => userData?.users ?? [], [userData?.users])

  const { data: tagData } = useQuery({
    queryKey: ['tagList', { pageNum: 1, pageSize: 1000 }],
    queryFn: () => getTagList({ pageNum: 1, pageSize: 1000 }),
  })
  const tagList = useMemo(() => tagData?.tags ?? [], [tagData?.tags])

  const flattenCategories = useMemo(() => {
    const categories: { id: string; title: string; level: number }[] = []
    const traverse = (nodes: Category[], level = 0) => {
      nodes.forEach((node) => {
        categories.push({
          id: String(node.id ?? '').trim(),
          title: node.title,
          level,
        })
        if (node.children?.length) {
          traverse(node.children, level + 1)
        }
      })
    }
    if (categoryData?.categoryTree) {
      traverse(categoryData.categoryTree)
    }
    return categories
  }, [categoryData])

  const form = useForm<ArticleFormValues>({
    resolver: zodResolver(articleFormSchema) as Resolver<ArticleFormValues>,
    defaultValues: {
      title: '',
      slug: '',
      description: '',
      author: '',
      userID: '',
      content: '',
      icon: '',
      video: '',
      images: [],
      type: 1,
      status: 1,
      categoryID: '',
      projectId: '',
      isTop: 1,
      isPasswd: 1,
      password: '',
      tag: '',
      showTime: new Date(),
    },
  })

  // 使用 ref 稳定 form 方法引用，避免 useEffect 依赖变化导致无限重渲染
  const formRef = useRef(form)
  formRef.current = form

  const { data: postDataRes, isLoading: postLoading } = useQuery({
    queryKey: ['post', id],
    queryFn: () => getPostDetail(id!),
    enabled: !!isEdit && !!id && !initialPost,
  })
  const post = id ? (initialPost ?? postDataRes ?? null) : null
  const isLoadingPost = isEdit && !!id && !initialPost && postLoading
  const loadFinishedNoPost =
    isEdit && !!id && !initialPost && !postLoading && !post

  // 同步 post 数据到表单
  useEffect(() => {
    const currentForm = formRef.current
    if (post) {
      const rawCategoryID =
        (post.categoryID != null ? String(post.categoryID) : '') ||
        (post.category &&
        typeof post.category === 'object' &&
        'id' in post.category
          ? String((post.category as { id: string | number }).id)
          : '')
      const categoryID = rawCategoryID.trim()
      const rawProjectId =
        (post.projectId != null ? String(post.projectId) : '') ||
        (post.project &&
        typeof post.project === 'object' &&
        'id' in post.project
          ? String((post.project as { id: number }).id)
          : '')
      const projectId = rawProjectId.trim()
      const values = {
        title: post.title || '',
        slug: post.slug || '',
        description: post.description || '',
        author: post.author || '',
        userID: post.userID ? String(post.userID) : '',
        content: post.content || '',
        icon: post.icon || '',
        video: post.video || '',
        images: Array.isArray(post.images) ? post.images : [],
        type: post.type || 1,
        status: [1, 2].includes(Number(post.status)) ? Number(post.status) : 1,
        categoryID: categoryID || '',
        projectId: projectId || '',
        isTop: post.isTop ?? 1,
        isPasswd: post.isPasswd ?? 1,
        password: post.password || '',
        tag: post.tag || '',
        showTime: post.showTime
          ? new Date(post.showTime)
          : post.createdAt
            ? new Date(post.createdAt)
            : new Date(),
      }
      currentForm.reset(values)
      // 解析 ext 为自定义字段列表（支持嵌套对象，展平为 key 如 meta.title）
      let newExtFields: Array<{ key: string; value: string }> = []
      try {
        const extStr = (post.ext ?? '').trim()
        if (extStr) {
          const obj = JSON.parse(extStr) as Record<string, unknown>
          if (obj && typeof obj === 'object' && !Array.isArray(obj)) {
            newExtFields = flattenExtToKeyValue(obj)
          }
        }
      } catch {
        // ignore
      }
      // 下一帧再 setValue 一次，避免 reset 与首帧渲染时序导致 field.value 未更新、分类/项目不回显
      const tid = setTimeout(() => {
        setExtFields(newExtFields)
        if (categoryID) currentForm.setValue('categoryID', categoryID)
        if (projectId) currentForm.setValue('projectId', projectId)
        currentForm.setValue('status', values.status)
      }, 0)
      return () => clearTimeout(tid)
    } else if (!isEdit) {
      currentForm.reset({
        title: '',
        slug: '',
        description: '',
        author: '', // Will be set by useEffect
        userID: '',
        content: '',
        icon: '',
        video: '',
        images: [],
        type: 1, // Default Article
        status: 1,
        categoryID: '',
        projectId: '', // 由 useEffect 在 projectList 加载后设为第一项
        isTop: 1,
        isPasswd: 1,
        password: '',
        tag: '',
        showTime: new Date(),
      })
      const tid = setTimeout(() => setExtFields([]), 0)
      return () => clearTimeout(tid)
    }
  }, [post, isEdit])

  // 标题变化时自动生成 slug（仅当 slug 为空时）
  useEffect(() => {
    const subscription = form.watch((value, { name }) => {
      if (name === 'title' && value.title) {
        const currentSlug = form.getValues('slug')
        if (!currentSlug) {
          const generated = generateSlug(value.title)
          if (generated) {
            form.setValue('slug', generated, { shouldValidate: false })
          }
        }
      }
    })
    return () => subscription.unsubscribe()
  }, [form])

  // Set default author for new posts
  useEffect(() => {
    if (!isEdit && userList.length > 0) {
      const currentAuthor = formRef.current.getValues('author')
      if (!currentAuthor) {
        formRef.current.setValue(
          'author',
          userList[0].nickname || userList[0].username
        )
        formRef.current.setValue('userID', userList[0].userID)
      }
    }
  }, [isEdit, userList])

  // Set default category for new posts（跳过 id 为空的节点，避免校验失败）
  useEffect(() => {
    if (!isEdit && flattenCategories.length > 0) {
      const currentCategory = formRef.current.getValues('categoryID')
      if (!currentCategory?.trim()) {
        const first = flattenCategories.find((c) => (c.id ?? '').trim() !== '')
        if (first) formRef.current.setValue('categoryID', first.id)
      }
    }
  }, [isEdit, flattenCategories])

  // Set default project for new posts
  useEffect(() => {
    if (!isEdit && projectList.length > 0) {
      const currentProject = formRef.current.getValues('projectId')
      if (!currentProject) {
        formRef.current.setValue('projectId', String(projectList[0].id))
      }
    }
  }, [isEdit, projectList])

  const createMutation = useMutation({
    mutationFn: (data: PostParams) => createPost(data),
    onSuccess: () => {
      toast.success(t('features.content.article.createSuccess'), {
        description: t('features.content.article.form.returnToListAfterCreate'),
      })
      queryClient.invalidateQueries({ queryKey: ['postList'] })
      navigate({ to: '/content/article' })
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id: postId, ...params }: PostParams & { id: number }) =>
      updatePost(postId, params),
    onSuccess: () => {
      toast.success(t('features.content.article.updateSuccess'), {
        description: t('features.content.article.form.returnToListAfterUpdate'),
      })
      queryClient.invalidateQueries({ queryKey: ['postList'] })
      navigate({ to: '/content/article' })
    },
  })

  const createPostMutate = createMutation.mutate
  const updatePostMutate = updateMutation.mutate

  const onInvalid = useCallback(
    (errors: FieldErrors<ArticleFormValues>) => {
      const mediaFields = ['description', 'icon', 'video', 'images']
      const settingsFields = [
        'type',
        'status',
        'categoryID',
        'projectId',
        'tag',
        'showTime',
        'isTop',
        'isPasswd',
        'password',
        'author',
        'userID',
      ]

      const hasMediaError = Object.keys(errors).some((key) =>
        mediaFields.includes(key)
      )
      if (hasMediaError) {
        setMediaSheetOpen(true)
      }

      const hasSettingsError = Object.keys(errors).some((key) =>
        settingsFields.includes(key)
      )
      if (hasSettingsError) {
        setSheetOpen(true)
      }

      if (hasMediaError || hasSettingsError) {
        toast.error(t('features.content.article.form.validationErrorInSheet'))
      }
    },
    [t]
  )

  const handleSubmit = useCallback(
    async (values: ArticleFormValues) => {
      if (createMutation.isPending || updateMutation.isPending) return
      // content 存 Slate 原数据（JSON），用于回显；htmlContent 存转换后的 HTML
      const htmlContent = slateContentToHtml(values.content ?? '')
      const extObj: Record<string, unknown> = {}
      extFields
        .filter((e) => e.key.trim() !== '')
        .forEach((e) => setExtNested(extObj, e.key.trim(), e.value.trim()))
      const ext = Object.keys(extObj).length > 0 ? JSON.stringify(extObj) : ''
      // 显式构造 payload，确保 status 等字段正确传递
      const payload = {
        title: values.title,
        slug: values.slug,
        description: values.description,
        author: values.author,
        userID: values.userID ? Number(values.userID) : undefined,
        content: values.content,
        htmlContent,
        ext,
        icon: values.icon,
        video: values.video,
        images: values.images,
        type: values.type,
        status: values.status ?? 1,
        categoryID: values.categoryID ? Number(values.categoryID) : undefined,
        projectId: values.projectId ? Number(values.projectId) : undefined,
        isTop: values.isTop,
        isPasswd: values.isPasswd,
        password: values.password,
        tag: values.tag,
        showTime: values.showTime
          ? format(values.showTime, 'yyyy-MM-dd HH:mm:ss')
          : undefined,
      }

      if (isEdit && id) {
        updatePostMutate({ ...payload, id })
      } else {
        createPostMutate(payload)
      }
    },
    [
      extFields,
      isEdit,
      id,
      createPostMutate,
      updatePostMutate,
      createMutation.isPending,
      updateMutation.isPending,
    ]
  )

  // Keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 's') {
        e.preventDefault()
        form.handleSubmit(handleSubmit, onInvalid)()
      }
    }
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [form, handleSubmit, onInvalid])

  const handleResourceSelect = (resource: ResourceItem) => {
    if (currentUploadField === 'images') {
      const current = form.getValues('images') ?? []
      form.setValue('images', [...current, resource.url])
    } else if (currentUploadField) {
      form.setValue(currentUploadField, resource.url)
    }
    setResourceOpen(false)
  }

  const openResourceUpload = (field: 'icon' | 'video' | 'images') => {
    setCurrentUploadField(field)
    setResourceType(field === 'video' ? 2 : 1) // 1: Image, 2: Video
    setResourceOpen(true)
  }

  const isLoading = createMutation.isPending || updateMutation.isPending

  if (isLoadingPost) {
    return (
      <div className='flex min-h-[200px] items-center justify-center px-4'>
        <p className='text-muted-foreground'>
          {t('features.content.article.loading')}
        </p>
      </div>
    )
  }

  if (loadFinishedNoPost) {
    return (
      <div className='flex min-h-[200px] flex-col items-center justify-center gap-4 px-4'>
        <p className='text-muted-foreground'>
          {t('features.content.article.notFound')}
        </p>
        <Button
          variant='outline'
          onClick={() => navigate({ to: '/content/article' })}
        >
          <ArrowLeft className='mr-2 h-4 w-4' />
          {t('features.content.article.form.backToList')}
        </Button>
      </div>
    )
  }

  return (
    <div className='flex h-[calc(100vh-4rem)] flex-col bg-background'>
      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(handleSubmit, onInvalid)}
          className='flex h-full flex-col'
        >
          {/* 顶栏：sticky 方便长文时随时保存，同时保持编辑器页更紧凑的层级 */}
          <div className='sticky top-0 z-10 shrink-0 border-b border-border/60 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/80'>
            <div className='mx-auto flex w-full max-w-[88rem] items-center justify-between gap-4 px-4 py-4 sm:px-6 lg:px-8'>
              <div className='flex min-w-0 items-center gap-3'>
                <Button
                  type='button'
                  variant='ghost'
                  size='icon'
                  className='shrink-0'
                  onClick={() => navigate({ to: '/content/article' })}
                  aria-label={t('features.content.article.form.cancel')}
                >
                  <ArrowLeft className='h-4 w-4' />
                </Button>
                <div className='min-w-0'>
                  <div className='text-[11px] font-semibold tracking-[0.18em] text-muted-foreground uppercase'>
                    Editor
                  </div>
                  <h1 className='truncate text-xl font-semibold tracking-tight text-foreground sm:text-2xl'>
                    {isEdit
                      ? t('features.content.article.form.editTitle')
                      : t('features.content.article.form.createTitle')}
                  </h1>
                  <p className='line-clamp-1 text-sm text-muted-foreground'>
                    {isEdit
                      ? t('features.content.article.form.editDescription')
                      : t('features.content.article.form.createDescription')}
                  </p>
                </div>
              </div>
              <div className='flex shrink-0 items-center gap-2'>
                <ArticleMediaSheet
                  open={mediaSheetOpen}
                  onOpenChange={setMediaSheetOpen}
                  form={form}
                  onOpenResourceUpload={openResourceUpload}
                />
                <ArticleSettingsSheet
                  open={sheetOpen}
                  onOpenChange={setSheetOpen}
                  form={form}
                  userList={userList}
                  tagList={tagList}
                  projectList={projectList}
                  flattenCategories={flattenCategories}
                  extFields={extFields}
                  setExtFields={setExtFields}
                />
                <Button
                  type='submit'
                  disabled={isLoading}
                  className='min-w-24 rounded-full px-5'
                >
                  {isLoading
                    ? t('features.content.article.form.submitting')
                    : isEdit
                      ? t('features.content.article.form.save')
                      : t('features.content.article.form.create')}
                </Button>
              </div>
            </div>
          </div>

          {/* 书写区：固定宽度 + 留白，标题与正文分区清晰 */}
          <div className='min-h-0 flex-1 overflow-y-auto'>
            <ArticleEditor form={form} isEdit={isEdit} />
          </div>
        </form>
      </Form>

      <ResourceUpload
        open={resourceOpen}
        onOpenChange={setResourceOpen}
        type={resourceType}
        onSelect={handleResourceSelect}
      />
    </div>
  )
}
