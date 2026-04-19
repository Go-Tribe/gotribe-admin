import { useState, useRef, useEffect } from 'react'
import { useWatch, type UseFormReturn } from 'react-hook-form'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormDescription,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { Badge } from '@/components/ui/badge'
import { DatePicker } from '@/components/date-picker'
import { Settings, X, ChevronsUpDown, Plus, Check, Trash2 } from 'lucide-react'
import { cn } from '@/lib/utils'
import { toast } from 'sonner'
import { createTag } from '../service/tag'
import { useI18n } from '@/context/i18n-provider'
import type { ArticleFormValues } from '../article-form-page'
import type { User, Project } from '@/shared/types'
import type { Tag } from '../types/tag'

function SettingsSection({
  title,
  description,
  children,
}: {
  title: string
  description: string
  children: React.ReactNode
}) {
  return (
    <section className="rounded-2xl border border-border/60 bg-muted/20 p-4 shadow-sm">
      <div className="mb-4 space-y-1">
        <h3 className="text-sm font-semibold tracking-tight">{title}</h3>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
      <div className="space-y-4">{children}</div>
    </section>
  )
}

// 简单的辅助函数
const getRandomColor = () => {
  const letters = '0123456789ABCDEF'
  let color = '#'
  for (let i = 0; i < 6; i++) {
    color += letters[Math.floor(Math.random() * 16)]
  }
  return color
}

const buildInlineTagSlug = (title: string) => {
  const normalized = title.trim().replace(/\s+/g, '-').slice(0, 30)
  if (normalized.length >= 2) {
    return normalized
  }
  return `tag-${Date.now().toString().slice(-8)}`
}

interface ArticleSettingsSheetProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  form: UseFormReturn<ArticleFormValues>
  userList: User[]
  tagList: Tag[]
  projectList: Project[]
  flattenCategories: Array<{ id: string; title: string; level: number }>
  extFields: Array<{ key: string; value: string }>
  setExtFields: React.Dispatch<React.SetStateAction<Array<{ key: string; value: string }>>>
}

export function ArticleSettingsSheet({
  open,
  onOpenChange,
  form,
  userList,
  tagList,
  projectList,
  flattenCategories,
  extFields,
  setExtFields,
}: ArticleSettingsSheetProps) {
  const { t } = useI18n()
  const queryClient = useQueryClient()
  const [searchTagValue, setSearchTagValue] = useState('')
  const pendingCreateTitleRef = useRef<string | null>(null)

  // 监听 tagList 变化，如果有待创建的标签且在列表中找到了，则自动选中
  useEffect(() => {
    const pendingTitle = pendingCreateTitleRef.current
    if (pendingTitle && tagList.length > 0) {
      const newTag = tagList.find(t => t.title === pendingTitle)
      if (newTag) {
        const currentTags = form.getValues('tag')
        const currentTagIds = currentTags ? currentTags.split(',').filter(Boolean) : []
        const newTagId = String(newTag.id)
        if (!currentTagIds.includes(newTagId)) {
          const newIds = [...currentTagIds, newTagId]
          form.setValue('tag', newIds.join(','))
          toast.success(t('features.content.tag.createSuccess'))
        }
        pendingCreateTitleRef.current = null
      }
    }
  }, [tagList, form, t])

  const createTagMutation = useMutation({
    mutationFn: createTag,
    onSuccess: () => {
      // 创建成功后只刷新列表，选中逻辑交给 useEffect 处理
      queryClient.invalidateQueries({ queryKey: ['tagList'] })
    },
    onError: () => {
      toast.error(t('features.content.tag.createError'))
      pendingCreateTitleRef.current = null
    }
  })

  const isPasswd = useWatch({
    control: form.control,
    name: 'isPasswd',
  })

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetTrigger asChild>
        <Button variant="outline" size="sm" type="button">
          <Settings className="h-4 w-4 mr-2" />
          {t('features.content.article.form.tabOtherInfo')}
        </Button>
      </SheetTrigger>
      <SheetContent className="w-[400px] sm:w-[540px] p-0">
        <SheetHeader className="px-6 py-4 border-b">
          <SheetTitle>{t('features.content.article.form.tabOtherInfo')}</SheetTitle>
          <SheetDescription>
            {t('features.content.article.form.editDescription')}
          </SheetDescription>
        </SheetHeader>
        <ScrollArea className="h-[calc(100vh-80px)] px-6 py-4">
          <div className="space-y-6 pb-8">
            <div className="rounded-2xl border border-border/60 bg-card/80 p-4 shadow-sm">
              <p className="text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground">
                {t('features.content.article.form.settingsSummaryLabel')}
              </p>
              <p className="mt-2 text-sm text-muted-foreground">
                {t('features.content.article.form.settingsSummaryDescription')}
              </p>
            </div>

            <SettingsSection
              title={t('features.content.article.form.settingsSections.publishTitle')}
              description={t('features.content.article.form.settingsSections.publishDescription')}
            >
              <div className="grid grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name='type'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('features.content.article.form.type')}</FormLabel>
                      <Select onValueChange={(val) => field.onChange(Number(val))} value={field.value?.toString()}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder={t('features.content.article.form.typePlaceholder')} />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value="1">{t('features.content.article.form.typeArticle')}</SelectItem>
                          <SelectItem value="2">{t('features.content.article.form.typePage')}</SelectItem>
                          <SelectItem value="3">{t('features.content.article.form.typeShortPost')}</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormDescription>{t('features.content.article.form.hints.type')}</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name='status'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('features.content.article.form.status')}</FormLabel>
                      <Select onValueChange={(val) => field.onChange(Number(val))} value={field.value?.toString()}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder={t('features.content.article.form.statusPlaceholder')} />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value="1">{t('features.content.article.status.draft')}</SelectItem>
                          <SelectItem value="2">{t('features.content.article.status.published')}</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormDescription>{t('features.content.article.form.hints.status')}</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <FormField
                control={form.control}
                name="showTime"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('features.content.article.form.showTime')}</FormLabel>
                    <FormControl>
                      <DatePicker
                        selected={field.value}
                        onSelect={field.onChange}
                        placeholder={t('features.content.article.form.showTimePlaceholder')}
                      />
                    </FormControl>
                    <FormDescription>{t('features.content.article.form.hints.showTime')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <div className="grid grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name="isTop"
                  render={({ field }) => (
                    <FormItem className="flex flex-row items-center justify-between rounded-2xl border border-border/60 bg-background p-4 shadow-sm">
                      <div className="space-y-1">
                        <FormLabel className="text-sm font-medium">{t('features.content.article.form.isTop')}</FormLabel>
                        <FormDescription>{t('features.content.article.form.hints.isTop')}</FormDescription>
                      </div>
                      <FormControl>
                        <Switch
                          checked={field.value === 2}
                          onCheckedChange={(checked) => field.onChange(checked ? 2 : 1)}
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="isPasswd"
                  render={({ field }) => (
                    <FormItem className="flex flex-row items-center justify-between rounded-2xl border border-border/60 bg-background p-4 shadow-sm">
                      <div className="space-y-1">
                        <FormLabel className="text-sm font-medium">{t('features.content.article.form.isPasswd')}</FormLabel>
                        <FormDescription>{t('features.content.article.form.hints.isPasswd')}</FormDescription>
                      </div>
                      <FormControl>
                        <Switch
                          checked={field.value === 2}
                          onCheckedChange={(checked) => field.onChange(checked ? 2 : 1)}
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
              </div>
              {isPasswd === 2 && (
                <FormField
                  control={form.control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('features.content.article.form.password')}</FormLabel>
                      <FormControl>
                        <Input
                          type="password"
                          placeholder={t('features.content.article.form.passwordPlaceholder')}
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>{t('features.content.article.form.hints.password')}</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              )}
            </SettingsSection>

            <SettingsSection
              title={t('features.content.article.form.settingsSections.taxonomyTitle')}
              description={t('features.content.article.form.settingsSections.taxonomyDescription')}
            >
              <FormField
                control={form.control}
                name='categoryID'
                render={({ field }) => {
                  const categoryValue = (field.value != null && field.value !== '') ? String(field.value).trim() : ''
                  const selectedCategory =
                    categoryValue !== ''
                      ? flattenCategories.find((c) => (c.id ?? '').trim() === categoryValue)
                      : undefined
                  return (
                    <FormItem>
                      <FormLabel>{t('features.content.article.form.category')}</FormLabel>
                      <Select
                        onValueChange={(val) => field.onChange(val ? val.trim() : '')}
                        value={categoryValue || undefined}
                      >
                        <FormControl>
                          <SelectTrigger>
                            {selectedCategory ? (
                              <span className="line-clamp-1 truncate">
                                {selectedCategory.title.trim()}
                              </span>
                            ) : (
                              <SelectValue placeholder={t('features.content.article.form.categoryPlaceholder')} />
                            )}
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {flattenCategories.map((category) => (
                            <SelectItem key={category.id} value={category.id}>
                              <span style={{ paddingLeft: `${category.level * 10}px` }}>{category.title.trim()}</span>
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <FormDescription>{t('features.content.article.form.hints.category')}</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )
                }}
              />
              <FormField
                control={form.control}
                name='projectID'
                render={({ field }) => {
                  const projectValue = (field.value != null && field.value !== '') ? String(field.value).trim() : ''
                  const selectedProject =
                    projectValue !== ''
                      ? projectList.find((p) => (p.projectID ?? '').trim() === projectValue)
                      : undefined
                  return (
                    <FormItem>
                      <FormLabel>{t('features.content.article.filter.project')}</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value || undefined}>
                        <FormControl>
                          <SelectTrigger>
                            {selectedProject ? (
                              <span className="line-clamp-1 flex items-center">{selectedProject.title}</span>
                            ) : (
                              <SelectValue placeholder={t('features.content.article.form.projectPlaceholder')} />
                            )}
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {projectList.map((project) => (
                            <SelectItem key={project.projectID} value={String(project.projectID ?? '').trim()}>
                              {project.title}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <FormDescription>{t('features.content.article.form.hints.project')}</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )
                }}
              />
              <FormField
              control={form.control}
              name='tag'
              render={({ field }) => {
                const selectedTagIds = field.value ? field.value.split(',').filter(Boolean) : []
                const selectedTags = tagList.filter(tag => selectedTagIds.includes(String(tag.id)))

                return (
                  <FormItem>
                    <FormLabel>{t('features.content.article.form.tag')}</FormLabel>
                    <Popover>
                      <PopoverTrigger asChild>
                        <FormControl>
                          <Button
                            variant="outline"
                            role="combobox"
                            className={cn(
                              "w-full justify-between h-auto min-h-10",
                              !selectedTagIds.length && "text-muted-foreground"
                            )}
                          >
                            <div className="flex flex-wrap gap-1">
                              {selectedTags.length > 0 ? (
                                selectedTags.map((tag) => (
                                  <Badge key={tag.id} variant="secondary" className="mr-1">
                                    {tag.title}
                                    <div
                                      className="ml-1 ring-offset-background rounded-full outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 cursor-pointer"
                                      onMouseDown={(e) => {
                                        e.preventDefault()
                                        e.stopPropagation()
                                      }}
                                      onClick={(e) => {
                                        e.preventDefault()
                                        e.stopPropagation()
                                        const newIds = selectedTagIds.filter(id => id !== String(tag.id))
                                        field.onChange(newIds.join(','))
                                      }}
                                    >
                                      <X className="h-3 w-3 text-muted-foreground hover:text-foreground" />
                                    </div>
                                  </Badge>
                                ))
                              ) : (
                                <span>{t('features.content.article.form.tagPlaceholder')}</span>
                              )}
                            </div>
                            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                          </Button>
                        </FormControl>
                      </PopoverTrigger>
                      <PopoverContent className="w-[400px] p-0" align="start">
                        <Command>
                          <CommandInput
                            placeholder={t('features.content.article.form.tagPlaceholder')}
                            value={searchTagValue}
                            onValueChange={setSearchTagValue}
                          />
                          <CommandList>
                            <CommandEmpty>
                              {searchTagValue && !tagList.some(tag => tag.title === searchTagValue) ? (
                                <div
                                  className="flex items-center p-2 text-sm cursor-pointer hover:bg-accent hover:text-accent-foreground"
                                  onMouseDown={(e) => {
                                    e.preventDefault()
                                    e.stopPropagation()
                                  }}
                                  onClick={(e) => {
                                    e.preventDefault()
                                    e.stopPropagation()
                                    pendingCreateTitleRef.current = searchTagValue
                                    createTagMutation.mutate({
                                      title: searchTagValue,
                                      slug: buildInlineTagSlug(searchTagValue),
                                      description: searchTagValue,
                                      color: getRandomColor()
                                    })
                                    setSearchTagValue('')
                                  }}
                                >
                                  <Plus className="mr-2 h-4 w-4" />
                                  {createTagMutation.isPending
                                    ? t('features.content.article.form.inlineTagCreating')
                                    : `${t('features.content.tag.createButton')} "${searchTagValue}"`}
                                </div>
                              ) : (
                                "No tag found."
                              )}
                            </CommandEmpty>
                            <CommandGroup>
                              {tagList.map((tag) => (
                                <CommandItem
                                  value={tag.title}
                                  key={tag.id}
                                  onSelect={() => {
                                    const tagId = String(tag.id)
                                    const isSelected = selectedTagIds.includes(tagId)
                                    let newIds
                                    if (isSelected) {
                                      newIds = selectedTagIds.filter(id => id !== tagId)
                                    } else {
                                      newIds = [...selectedTagIds, tagId]
                                    }
                                    field.onChange(newIds.join(','))
                                  }}
                                >
                                  <Check
                                    className={cn(
                                      "mr-2 h-4 w-4",
                                      selectedTagIds.includes(String(tag.id))
                                        ? "opacity-100"
                                        : "opacity-0"
                                    )}
                                  />
                                  {tag.title}
                                </CommandItem>
                              ))}
                            </CommandGroup>
                          </CommandList>
                        </Command>
                      </PopoverContent>
                    </Popover>
                    <FormDescription>{t('features.content.article.form.hints.tag')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )
              }}
            />
            </SettingsSection>

            <SettingsSection
              title={t('features.content.article.form.settingsSections.authorTitle')}
              description={t('features.content.article.form.settingsSections.authorDescription')}
            >
              <FormField
                control={form.control}
                name="userID"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('features.content.article.form.author')}</FormLabel>
                    <FormControl>
                      <Select
                        onValueChange={(val) => {
                          const user = userList.find((u) => u.userID === val)
                          field.onChange(val)
                          if (user) {
                            form.setValue('author', user.nickname || user.username, { shouldValidate: true })
                          }
                        }}
                        value={field.value || undefined}
                      >
                        <SelectTrigger>
                          <SelectValue placeholder={t('features.content.article.form.authorPlaceholder')} />
                        </SelectTrigger>
                        <SelectContent>
                          {userList.map((user) => (
                            <SelectItem key={user.userID} value={user.userID}>
                              {user.nickname || user.username}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </FormControl>
                    <FormDescription>{t('features.content.article.form.hints.author')}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </SettingsSection>

            <SettingsSection
              title={t('features.content.article.form.settingsSections.extTitle')}
              description={t('features.content.article.form.settingsSections.extDescription')}
            >
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium text-muted-foreground">{t('features.content.article.form.extCustom')}</span>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => setExtFields((prev) => [...prev, { key: '', value: '' }])}
                >
                  <Plus className="mr-1 h-4 w-4" />
                  {t('features.content.article.form.extAdd')}
                </Button>
              </div>
              {extFields.length === 0 ? (
                <div className="rounded-2xl border border-dashed border-border/70 bg-background px-4 py-6 text-sm text-muted-foreground">
                  {t('features.content.article.form.extEmpty')}
                </div>
              ) : (
                <div className="space-y-3">
                  {extFields.map((item, index) => (
                    <div key={index} className="rounded-2xl border border-border/60 bg-background p-3 shadow-sm">
                      <div className="flex items-start gap-2">
                        <div className="grid flex-1 gap-2 md:grid-cols-2">
                          <Input
                            placeholder={t('features.content.article.form.extKeyPlaceholder')}
                            value={item.key}
                            onChange={(e) =>
                              setExtFields((prev) => {
                                const next = [...prev]
                                next[index] = { ...next[index], key: e.target.value }
                                return next
                              })
                            }
                            className="min-w-0 font-mono text-sm"
                          />
                          <Input
                            placeholder={t('features.content.article.form.extValuePlaceholder')}
                            value={item.value}
                            onChange={(e) =>
                              setExtFields((prev) => {
                                const next = [...prev]
                                next[index] = { ...next[index], value: e.target.value }
                                return next
                              })
                            }
                            className="min-w-0 font-mono text-sm"
                          />
                        </div>
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          className="shrink-0 text-muted-foreground hover:text-destructive"
                          onClick={() => setExtFields((prev) => prev.filter((_, i) => i !== index))}
                          aria-label={t('features.content.article.form.extDelete')}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </SettingsSection>
          </div>
        </ScrollArea>
      </SheetContent>
    </Sheet>
  )
}
