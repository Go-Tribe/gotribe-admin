import { useEffect, useMemo, useState } from 'react'
import { ImageIcon, Video, Check, ChevronsUpDown } from 'lucide-react'
import { useForm, useWatch } from 'react-hook-form'
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
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { ResourceUpload, type ResourceItem, FILE_TYPE } from '@/components/resource-upload'
import { useI18n } from '@/context/i18n-provider'
import { useQuery } from '@tanstack/react-query'
import { getPostList } from '@/features/content/service/post'
import type { AdCreateParams, AdUpdateParams, Ad } from '../types'

const createAdFormSchema = (t: (key: string) => string) =>
  z
    .object({
      title: z.string().min(1, t('features.operation.advertising.form.validation.titleRequired')),
      description: z
        .string()
        .min(1, t('features.operation.advertising.form.validation.descriptionRequired')),
      sceneID: z.string().min(1, t('features.operation.advertising.form.validation.sceneRequired')),
      status: z.number().optional(),
      urlType: z.number().min(1).max(2),
      image: z.string().optional(),
      video: z.string().optional(),
      ext: z.string().optional(),
      sort: z.number().min(1, t('features.operation.advertising.form.validation.sortRequired')),
      url: z.string().optional(),
      postID: z.string().optional(),
    })
    .refine(
      (data) => {
        // 如果类型是链接(1)，url 必填
        if (data.urlType === 1) {
          return data.url && data.url.trim().length > 0
        }
        // 如果类型是文章(2)，postID 必填
        if (data.urlType === 2) {
          return data.postID && data.postID.trim().length > 0
        }
        return true
      },
      {
        message: t('features.operation.advertising.form.validation.urlOrPostRequired'),
        path: ['url'],
      }
    )

type AdFormValues = z.infer<ReturnType<typeof createAdFormSchema>>

type AdFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSubmit?: (data: AdCreateParams) => void
  onSubmitUpdate?: (adID: string, data: AdUpdateParams) => void
  isLoading?: boolean
  sceneList: { adSceneID: string; title: string }[]
  editAd?: Ad | null
}

export function AdFormDialog({
  open,
  onOpenChange,
  onSubmit,
  onSubmitUpdate,
  isLoading = false,
  sceneList,
  editAd,
}: AdFormDialogProps) {
  const { t } = useI18n()
  const isEdit = Boolean(editAd)
  const [imageResourceDialogOpen, setImageResourceDialogOpen] = useState(false)
  const [videoResourceDialogOpen, setVideoResourceDialogOpen] = useState(false)
  const [postSearchOpen, setPostSearchOpen] = useState(false)
  const [postSearchKeyword, setPostSearchKeyword] = useState('')
  const adFormSchema = useMemo(() => createAdFormSchema(t), [t])
  const form = useForm<AdFormValues>({
    resolver: zodResolver(adFormSchema),
    defaultValues: {
      title: '',
      description: '',
      sceneID: '',
      status: 1,
      urlType: 1,
      image: '',
      video: '',
      ext: '',
      sort: 1,
      url: '',
      postID: '',
    },
  })

  const urlType = useWatch({ control: form.control, name: 'urlType' })

  // 获取文章列表（支持标题搜索）
  const { data: postListData, isLoading: isLoadingPosts } = useQuery({
    queryKey: ['postList', { title: postSearchKeyword, page_num: 1, page_size: 100 }],
    queryFn: () =>
      getPostList({
        title: postSearchKeyword || undefined,
        pageNum: 1,
        pageSize: 100,
      }),
    enabled: postSearchOpen && urlType === 2,
  })

  const postList = useMemo(() => postListData?.posts ?? [], [postListData?.posts])

  useEffect(() => {
    if (!open) return
    if (editAd) {
      // 根据 urlType 判断是链接还是文章
      const urlTypeValue = editAd.urlType ?? 1
      form.reset({
        title: (editAd.title ?? '').trim(),
        description: (editAd.description ?? '').trim(),
        sceneID: editAd.sceneID != null ? String(editAd.sceneID) : '',
        status: editAd.status ?? 1,
        urlType: urlTypeValue,
        image: (editAd.image ?? '').trim(),
        video: (editAd.video ?? '').trim(),
        ext: (editAd.ext ?? '').trim(),
        sort: editAd.sort ?? 1,
        url: urlTypeValue === 1 ? (editAd.url ?? '').trim() : '',
        postID: urlTypeValue === 2 ? (editAd.url ?? '').trim() : '',
      })
    } else {
      form.reset({
        title: '',
        description: '',
        sceneID: '',
        status: 1,
        urlType: 1,
        image: '',
        video: '',
        ext: '',
        sort: 1,
        url: '',
        postID: '',
      })
    }
  }, [open, editAd, form])

  function handleSubmit(values: AdFormValues) {
    // 根据 urlType 决定使用 url 还是 postID
    const urlValue = values.urlType === 1 ? values.url?.trim() : undefined
    const postIDValue = values.urlType === 2 ? values.postID?.trim() : undefined

    if (isEdit && editAd && onSubmitUpdate) {
      onSubmitUpdate(editAd.id != null ? String(editAd.id) : '', {
        title: values.title.trim(),
        description: values.description.trim(),
        sceneID: values.sceneID ? Number(values.sceneID) : undefined,
        status: values.status,
        urlType: values.urlType,
        image: (values.image ?? '').trim() || undefined,
        video: (values.video ?? '').trim() || undefined,
        ext: (values.ext ?? '').trim() || undefined,
        sort: values.sort,
        url: urlValue || postIDValue || undefined,
      })
    } else if (onSubmit) {
      const sceneID = (values.sceneID ?? '').trim()
      onSubmit({
        title: values.title.trim(),
        description: values.description.trim(),
        sceneID: Number(sceneID),
        status: values.status ?? 1,
        urlType: values.urlType,
        image: (values.image ?? '').trim() || undefined,
        video: (values.video ?? '').trim() || undefined,
        ext: (values.ext ?? '').trim() || undefined,
        sort: values.sort,
        url: urlValue || postIDValue || undefined,
      })
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-[600px] max-h-[90vh] flex flex-col'>
        <DialogHeader className='shrink-0'>
          <DialogTitle>
            {isEdit
              ? t('features.operation.advertising.form.editTitle')
              : t('features.operation.advertising.form.dialogTitle')}
          </DialogTitle>
          <DialogDescription>
            {isEdit
              ? t('features.operation.advertising.form.editDescription')
              : t('features.operation.advertising.form.createDescription')}
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
                  <FormLabel >
                    {t('features.operation.advertising.form.title')}
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t('features.operation.advertising.form.titlePlaceholder')}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage  />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='description'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel >
                    {t('features.operation.advertising.form.description')}
                  </FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder={t('features.operation.advertising.form.descriptionPlaceholder')}
                      className='min-h-[80px] resize-none'
                      {...field}
                    />
                  </FormControl>
                  <FormMessage  />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='sceneID'
              render={({ field }) => {
                const displayValue = (field.value ?? '').trim()
                const selectValue = displayValue || '__empty__'
                return (
                  <FormItem className='space-y-2'>
                    <FormLabel >
                      {t('features.operation.advertising.form.scene')}
                    </FormLabel>
                    <Select
                      key={`scene-${editAd?.id ?? 'create'}-${sceneList.length}-${selectValue}`}
                      value={selectValue}
                      onValueChange={(v) => field.onChange(v === '__empty__' ? '' : v)}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue
                            placeholder={t('features.operation.advertising.form.scenePlaceholder')}
                          />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value='__empty__'>
                          {t('features.operation.advertising.form.scenePlaceholder')}
                        </SelectItem>
                        {sceneList
                          .filter((s) => (s.adSceneID ?? '').trim() !== '')
                          .map((s) => (
                            <SelectItem
                              key={s.adSceneID}
                              value={(s.adSceneID ?? '').trim()}
                            >
                              {s.title}
                            </SelectItem>
                          ))}
                      </SelectContent>
                    </Select>
                    <FormMessage  />
                  </FormItem>
                )
              }}
            />
            <FormField
              control={form.control}
              name='urlType'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel >
                    {t('features.operation.advertising.form.urlType.label')}
                  </FormLabel>
                  <FormControl>
                    <Select
                      onValueChange={(val) => {
                        field.onChange(Number(val))
                        // 切换类型时清空 url 和 postID
                        form.setValue('url', '')
                        form.setValue('postID', '')
                      }}
                      value={String(field.value ?? 1)}
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value='1'>
                          {t('features.operation.advertising.form.urlType.link')}
                        </SelectItem>
                        <SelectItem value='2'>
                          {t('features.operation.advertising.form.urlType.article')}
                        </SelectItem>
                      </SelectContent>
                    </Select>
                  </FormControl>
                  <FormMessage  />
                </FormItem>
              )}
            />
            {urlType === 1 ? (
              <FormField
                control={form.control}
                name='url'
                render={({ field }) => (
                  <FormItem className='space-y-2'>
                    <FormLabel >
                      {t('features.operation.advertising.form.url')}
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t('features.operation.advertising.form.urlPlaceholder')}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage  />
                  </FormItem>
                )}
              />
            ) : (
              <FormField
                control={form.control}
                name='postID'
                render={({ field }) => (
                  <FormItem className='space-y-2'>
                    <FormLabel >
                      {t('features.operation.advertising.form.article')}
                    </FormLabel>
                    <Popover
                      open={postSearchOpen}
                      onOpenChange={(open) => {
                        setPostSearchOpen(open)
                        if (!open) {
                          setPostSearchKeyword('')
                        }
                      }}
                    >
                      <PopoverTrigger asChild>
                        <FormControl>
                          <Button
                            variant='outline'
                            role='combobox'
                            className={cn(
                              'w-full justify-between',
                              !field.value && 'text-muted-foreground'
                            )}
                          >
                            {field.value
                              ? postList.find((post) => post.postID === field.value)?.title ||
                                t('features.operation.advertising.form.selectArticle')
                              : t('features.operation.advertising.form.selectArticle')}
                            <ChevronsUpDown className='ml-2 h-4 w-4 shrink-0 opacity-50' />
                          </Button>
                        </FormControl>
                      </PopoverTrigger>
                      <PopoverContent className='w-[400px] p-0' align='start'>
                        <Command shouldFilter={false}>
                          <CommandInput
                            placeholder={t('features.operation.advertising.form.searchArticle')}
                            value={postSearchKeyword}
                            onValueChange={(value) => {
                              setPostSearchKeyword(value)
                            }}
                          />
                          <CommandList>
                            {isLoadingPosts ? (
                              <div className='py-6 text-center text-sm text-muted-foreground'>
                                {t('features.operation.advertising.form.loadingArticles')}
                              </div>
                            ) : postList.length === 0 ? (
                              <CommandEmpty>
                                {t('features.operation.advertising.form.noArticles')}
                              </CommandEmpty>
                            ) : (
                              <CommandGroup>
                                {postList.map((post) => (
                                  <CommandItem
                                    key={post.postID}
                                    value={`${post.postID}-${post.title}`}
                                    onSelect={() => {
                                      form.setValue('postID', post.postID, {
                                        shouldValidate: true,
                                      })
                                      setPostSearchOpen(false)
                                      setPostSearchKeyword('')
                                    }}
                                  >
                                    <Check
                                      className={cn(
                                        'mr-2 h-4 w-4',
                                        field.value === post.postID
                                          ? 'opacity-100'
                                          : 'opacity-0'
                                      )}
                                    />
                                    {post.title}
                                  </CommandItem>
                                ))}
                              </CommandGroup>
                            )}
                          </CommandList>
                        </Command>
                      </PopoverContent>
                    </Popover>
                    <FormMessage  />
                  </FormItem>
                )}
              />
            )}
            <FormField
              control={form.control}
              name='image'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel >
                    {t('features.operation.advertising.form.image')}
                  </FormLabel>
                  <FormControl>
                    <button
                      type='button'
                      onClick={() => setImageResourceDialogOpen(true)}
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
                            {t('features.operation.advertising.form.selectImage')}
                          </span>
                        </>
                      )}
                    </button>
                  </FormControl>
                  <FormMessage  />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='video'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel >
                    {t('features.operation.advertising.form.video')}
                  </FormLabel>
                  <FormControl>
                    <button
                      type='button'
                      onClick={() => setVideoResourceDialogOpen(true)}
                      className={cn(
                        'flex flex-col items-center justify-center h-20 w-20 rounded-lg border-2 border-dashed border-border bg-muted/30 overflow-hidden transition-all shrink-0',
                        'cursor-pointer hover:border-primary/50 hover:bg-muted/50'
                      )}
                    >
                      {field.value ? (
                        <div className='w-full h-full flex flex-col items-center justify-center bg-muted/30 p-1'>
                          <Video className='h-8 w-8 text-muted-foreground mb-1 shrink-0' />
                          <span className='text-xs text-muted-foreground truncate w-full text-center'>
                            {field.value.split('/').pop()}
                          </span>
                        </div>
                      ) : (
                        <>
                          <Video className='h-8 w-8 text-muted-foreground mb-1' />
                          <span className='text-xs text-muted-foreground'>
                            {t('features.operation.advertising.form.selectVideo')}
                          </span>
                        </>
                      )}
                    </button>
                  </FormControl>
                  <FormMessage  />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='ext'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel >
                    {t('features.operation.advertising.form.ext')}
                  </FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder={t('features.operation.advertising.form.extPlaceholder')}
                      className='min-h-[80px] resize-none'
                      {...field}
                    />
                  </FormControl>
                  <FormMessage  />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='sort'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel >
                    {t('features.operation.advertising.form.sort')}
                  </FormLabel>
                  <FormControl>
                    <Input
                      type='number'
                      min={1}
                      placeholder={t('features.operation.advertising.form.sortPlaceholder')}
                      value={field.value ?? ''}
                      onChange={(e) => field.onChange(Number(e.target.value))}
                    />
                  </FormControl>
                  <FormMessage  />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='status'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel >
                    {t('features.operation.advertising.form.status')}
                  </FormLabel>
                  <Select
                    onValueChange={(val) => field.onChange(Number(val))}
                    value={field.value ? String(field.value) : '1'}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value='1'>
                        {t('features.operation.advertising.advertising.status.unpublished')}
                      </SelectItem>
                      <SelectItem value='2'>
                        {t('features.operation.advertising.advertising.status.published')}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage  />
                </FormItem>
              )}
            />
            <DialogFooter className='shrink-0'>
              <Button
                type='button'
                variant='outline'
                onClick={() => onOpenChange(false)}
                disabled={isLoading}
              >
                {t('features.operation.advertising.form.cancel')}
              </Button>
              <Button type='submit' disabled={isLoading}>
                {isLoading ? '...' : t('features.operation.advertising.form.submit')}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
      <ResourceUpload
        open={imageResourceDialogOpen}
        onOpenChange={setImageResourceDialogOpen}
        onSelect={(resource: ResourceItem) => {
          form.setValue('image', resource.url, { shouldValidate: true })
          setImageResourceDialogOpen(false)
        }}
        type={FILE_TYPE.IMAGE}
        title={t('features.operation.advertising.form.selectImageDialogTitle')}
        description={t('features.operation.advertising.form.selectImageDialogDesc')}
      />
      <ResourceUpload
        open={videoResourceDialogOpen}
        onOpenChange={setVideoResourceDialogOpen}
        onSelect={(resource: ResourceItem) => {
          form.setValue('video', resource.url, { shouldValidate: true })
          setVideoResourceDialogOpen(false)
        }}
        type={FILE_TYPE.VIDEO}
        title={t('features.operation.advertising.form.selectVideoDialogTitle')}
        description={t('features.operation.advertising.form.selectVideoDialogDesc')}
      />
    </Dialog>
  )
}
