import { useEffect, useMemo, useState } from 'react'
import { Check, ChevronsUpDown } from 'lucide-react'
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
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
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
import { useI18n } from '@/context/i18n-provider'
import { useQuery } from '@tanstack/react-query'
import { getUserList } from '@/shared/api'
import type { PointCreateParams } from '../types/point'

const createPointFormSchema = (t: (key: string) => string) =>
  z.object({
    userID: z.string().min(1, t('features.operation.point.form.validation.userIDRequired')),
    projectId: z.number().min(1, t('features.operation.point.form.validation.projectIDRequired')),
    point: z.number().min(1, t('features.operation.point.form.validation.pointRequired')),
  })

type PointFormValues = z.infer<ReturnType<typeof createPointFormSchema>>

type PointFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSubmit: (data: PointCreateParams) => void
  isLoading?: boolean
  projectList: { id: number; title: string }[]
}

export function PointFormDialog({
  open,
  onOpenChange,
  onSubmit,
  isLoading = false,
  projectList,
}: PointFormDialogProps) {
  const { t } = useI18n()
  const [userSearchOpen, setUserSearchOpen] = useState(false)
  const [userSearchKeyword, setUserSearchKeyword] = useState('')
  const [selectedProjectId, setSelectedProjectId] = useState<number | undefined>(undefined)

  const pointFormSchema = useMemo(() => createPointFormSchema(t), [t])
  const form = useForm<PointFormValues>({
    resolver: zodResolver(pointFormSchema),
    defaultValues: {
      userID: '',
      projectId: undefined,
      point: 0,
    },
  })

  // 获取用户列表（支持搜索）
  const { data: userData, isLoading: isLoadingUsers } = useQuery({
    queryKey: ['userList', { current: 1, pageNum: 1, pageSize: 1000, projectId: selectedProjectId || undefined }],
    queryFn: () => getUserList({ current: 1, pageNum: 1, pageSize: 1000, projectId: selectedProjectId || undefined }),
    enabled: open && selectedProjectId != null,
  })

  const userList = useMemo(() => {
    const users = userData?.users || []
    if (!userSearchKeyword) return users
    const keyword = userSearchKeyword.toLowerCase()
    return users.filter(
      (user) =>
        user.userID.toLowerCase().includes(keyword) ||
        user.nickname?.toLowerCase().includes(keyword) ||
        user.username?.toLowerCase().includes(keyword)
    )
  }, [userData?.users, userSearchKeyword])

  useEffect(() => {
    if (!open) {
      form.reset({
        userID: '',
        projectId: undefined,
        point: 0,
      })
      setUserSearchKeyword('')
      setSelectedProjectId(undefined)
    }
  }, [open, form])

  // 当项目改变时，清空用户选择
  useEffect(() => {
    const currentProjectId = form.getValues('projectId')
    if (selectedProjectId !== currentProjectId) {
      setSelectedProjectId(currentProjectId)
      form.setValue('userID', '')
      setUserSearchKeyword('')
    }
  }, [form.watch('projectId'), form])

  function handleSubmit(values: PointFormValues) {
    onSubmit({
      userID: Number(values.userID),
      projectId: values.projectId,
      point: values.point,
    })
  }

  const selectedUser = userList.find((user) => user.userID === form.watch('userID'))

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-[600px] max-h-[90vh] flex flex-col'>
        <DialogHeader className='shrink-0'>
          <DialogTitle>
            {t('features.operation.point.form.dialogTitle')}
          </DialogTitle>
          <DialogDescription>
            {t('features.operation.point.form.createDescription')}
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className='flex-1 overflow-y-auto pr-2 space-y-4 min-h-0'
          >
            <FormField
              control={form.control}
              name='projectId'
              render={({ field }) => {
                const displayValue = field.value != null ? String(field.value) : ''
                const selectValue = displayValue || '__empty__'
                return (
                  <FormItem className='space-y-2'>
                    <FormLabel>
                      {t('features.operation.point.form.project')}
                    </FormLabel>
                    <Select
                      key={`project-${selectValue}`}
                      value={selectValue}
                      onValueChange={(v) => {
                        const projectId = v === '__empty__' ? undefined : Number(v)
                        field.onChange(projectId)
                        setSelectedProjectId(projectId)
                      }}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue
                            placeholder={t('features.operation.point.form.projectPlaceholder')}
                          />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value='__empty__'>
                          {t('features.operation.point.form.projectPlaceholder')}
                        </SelectItem>
                        {projectList
                          .filter((p) => p.id != null)
                          .map((p) => (
                            <SelectItem
                              key={p.id}
                              value={String(p.id)}
                            >
                              {p.title}
                            </SelectItem>
                          ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )
              }}
            />
            <FormField
              control={form.control}
              name='userID'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>
                    {t('features.operation.point.form.user')}
                  </FormLabel>
                  <Popover open={userSearchOpen} onOpenChange={setUserSearchOpen}>
                    <PopoverTrigger asChild>
                      <FormControl>
                        <Button
                          variant='outline'
                          role='combobox'
                          className={cn(
                            'w-full justify-between',
                            !field.value && 'text-muted-foreground'
                          )}
                          disabled={selectedProjectId == null}
                        >
                          {field.value && selectedUser
                            ? `${selectedUser.nickname || selectedUser.username} (${selectedUser.userID})`
                            : t('features.operation.point.form.selectUser')}
                          <ChevronsUpDown className='ml-2 h-4 w-4 shrink-0 opacity-50' />
                        </Button>
                      </FormControl>
                    </PopoverTrigger>
                    <PopoverContent className='w-[400px] p-0' align='start'>
                      <Command shouldFilter={false}>
                        <CommandInput
                          placeholder={t('features.operation.point.form.searchUser')}
                          value={userSearchKeyword}
                          onValueChange={(value) => {
                            setUserSearchKeyword(value)
                          }}
                        />
                        <CommandList>
                          {isLoadingUsers ? (
                            <div className='py-6 text-center text-sm text-muted-foreground'>
                              {t('features.operation.point.form.loadingUsers')}
                            </div>
                          ) : userList.length === 0 ? (
                            <CommandEmpty>
                              {t('features.operation.point.form.noUsers')}
                            </CommandEmpty>
                          ) : (
                            <CommandGroup>
                              {userList.map((user) => (
                                <CommandItem
                                  key={user.userID}
                                  value={`${user.userID}-${user.nickname || user.username}`}
                                  onSelect={() => {
                                    form.setValue('userID', user.userID, {
                                      shouldValidate: true,
                                    })
                                    setUserSearchOpen(false)
                                    setUserSearchKeyword('')
                                  }}
                                >
                                  <Check
                                    className={cn(
                                      'mr-2 h-4 w-4',
                                      field.value === user.userID
                                        ? 'opacity-100'
                                        : 'opacity-0'
                                    )}
                                  />
                                  {user.nickname || user.username} ({user.userID})
                                </CommandItem>
                              ))}
                            </CommandGroup>
                          )}
                        </CommandList>
                      </Command>
                    </PopoverContent>
                  </Popover>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='point'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>
                    {t('features.operation.point.form.point')}
                  </FormLabel>
                  <FormControl>
                    <Input
                      type='number'
                      placeholder={t('features.operation.point.form.pointPlaceholder')}
                      {...field}
                      onChange={(e) => {
                        const value = e.target.value === '' ? 0 : Number(e.target.value)
                        field.onChange(value)
                      }}
                      value={field.value || ''}
                    />
                  </FormControl>
                  <FormMessage />
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
                {t('features.operation.point.form.cancel')}
              </Button>
              <Button type='submit' disabled={isLoading}>
                {isLoading ? '...' : t('features.operation.point.form.submit')}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
