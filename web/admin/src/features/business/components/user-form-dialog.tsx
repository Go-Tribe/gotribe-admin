import { useEffect, useMemo, useState } from 'react'
import { useForm, type Resolver } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { useQuery } from '@tanstack/react-query'
import { ImageIcon } from 'lucide-react'
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
import { PasswordInput } from '@/components/password-input'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { ResourceUpload } from '@/components/resource-upload'
import type { User } from '../types/user'
import { getProjectList } from '../service/project'
import { useI18n } from '@/context/i18n-provider'

const alphanumRegex = /^[a-zA-Z0-9]+$/

const createUserFormSchema = (t: (key: string) => string, isEdit: boolean) =>
  z.object({
    avatarURL: z.string().optional(),
    username: isEdit
      ? z.string().optional()
      : z
          .string()
          .min(2, t('features.business.user.form.validation.usernameLength'))
          .max(20, t('features.business.user.form.validation.usernameLength'))
          .regex(alphanumRegex, t('features.business.user.form.validation.usernamePattern')),
    nickname: z.string().min(2, t('features.business.user.form.validation.nicknameLength')).max(20, t('features.business.user.form.validation.nicknameLength')),
    password: isEdit
      ? z.string().optional()
      : z.string().min(6, t('features.business.user.form.validation.passwordLength')).max(20, t('features.business.user.form.validation.passwordLength')),
    phone: z.string().min(1, t('features.business.user.form.validation.phoneRequired')),
    email: z.string().min(1, t('features.business.user.form.validation.emailRequired')),
    projectId: isEdit
      ? z.number().optional()
      : z.number().min(1, t('features.business.user.form.validation.projectIDRequired')),
  })

type UserFormValues = {
  avatarURL?: string
  username?: string
  nickname: string
  password?: string
  phone?: string
  email?: string
  projectId?: number
}

type UserFormDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  user: User | null
  onSubmit: (data: Partial<User>) => void
  isLoading?: boolean
}

export function UserFormDialog({
  open,
  onOpenChange,
  user,
  onSubmit,
  isLoading = false,
}: UserFormDialogProps) {
  const { t } = useI18n()
  const isEdit = !!user
  const [avatarResourceDialogOpen, setAvatarResourceDialogOpen] = useState(false)

  const userFormSchema = useMemo(
    () => createUserFormSchema(t, isEdit),
    [t, isEdit]
  )

  // 获取项目列表
  const { data: projectData } = useQuery({
    queryKey: ['projectList', { current: 1, pageNum: 1, pageSize: 1000 }],
    queryFn: () => getProjectList({ current: 1, pageNum: 1, pageSize: 1000 }),
  })

  const projectList = projectData?.projects || []

  const form = useForm<UserFormValues>({
    resolver: zodResolver(userFormSchema) as Resolver<UserFormValues>,
    defaultValues: {
      avatarURL: '',
      username: '',
      nickname: '',
      password: '',
      phone: '',
      email: '',
      projectId: undefined,
    },
  })

  useEffect(() => {
    if (open) {
      if (isEdit && user) {
        const u = user as User & { Phone?: string; Email?: string }
        form.reset({
          avatarURL: u.avatarURL || '',
          username: u.username || '',
          nickname: u.nickname || '',
          password: '',
          phone: u.phone,
          email: u.email,
          projectId: u.projectId,
        })
      } else {
        form.reset({
          avatarURL: '',
          username: '',
          nickname: '',
          password: '',
          phone: '',
          email: '',
          projectId: undefined,
        })
      }
    }
  }, [open, isEdit, user, form])

  const handleSubmit = (values: UserFormValues) => {
    if (isEdit && user) {
      // UpdateUserRequest: Nickname, Email, Phone, Password
      const payload: Partial<User> = {
        id: user.id,
        nickname: values.nickname,
        email: values.email?.trim() || undefined,
        phone: values.phone?.trim() || undefined,
      }
      if (values.avatarURL?.trim()) payload.avatarURL = values.avatarURL.trim()
      if (values.password?.trim()) payload.password = values.password
      onSubmit(payload)
    } else {
      // CreateUserRequest: Username, Nickname, Email, Phone, ProjectID, Password
      onSubmit({
        username: values.username!.trim(),
        nickname: values.nickname.trim(),
        email: values.email?.trim() || undefined,
        phone: values.phone?.trim() || undefined,
        projectId: values.projectId!,
        password: values.password!.trim(),
        avatarURL: values.avatarURL?.trim(),
      })
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-[600px] max-h-[90vh] flex flex-col'>
        <DialogHeader className='shrink-0'>
          <DialogTitle>{isEdit ? t('features.business.user.form.editTitle') : t('features.business.user.form.createTitle')}</DialogTitle>
          <DialogDescription>
            {isEdit ? t('features.business.user.form.editDescription') : t('features.business.user.form.createDescription')}
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className='flex-1 overflow-y-auto pr-2 space-y-4 min-h-0'
          >
            {!isEdit && (
              <FormField
                control={form.control}
                name='avatarURL'
                render={({ field }) => (
                  <FormItem className='space-y-2'>
                    <FormLabel>{t('features.business.user.form.fields.avatarURL')}
                    </FormLabel>
                    <FormControl>
                      <button
                        type='button'
                        onClick={() => setAvatarResourceDialogOpen(true)}
                        className='flex items-center justify-center h-14 w-14 rounded-full border-2 border-dashed border-muted-foreground/30 bg-muted/50 hover:border-primary/50 hover:bg-muted transition-colors cursor-pointer overflow-hidden focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2'
                      >
                        {field.value ? (
                          <img
                            src={field.value}
                            alt=''
                            className='h-full w-full object-cover'
                          />
                        ) : (
                          <ImageIcon className='h-7 w-7 text-muted-foreground' />
                        )}
                      </button>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}
            <FormField
              control={form.control}
              name='username'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.business.user.form.fields.username')}
                  </FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t('features.business.user.form.fields.usernamePlaceholder')}
                      {...field}
                      disabled={isEdit}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='nickname'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.business.user.form.fields.nickname')}
                  </FormLabel>
                  <FormControl>
                    <Input placeholder={t('features.business.user.form.fields.nicknamePlaceholder')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='password'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.business.user.form.fields.password')}
                  </FormLabel>
                  <FormControl>
                    <PasswordInput
                      placeholder={isEdit ? t('features.business.user.form.fields.passwordEditPlaceholder') : t('features.business.user.form.fields.passwordPlaceholder')}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='phone'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.business.user.form.fields.phone')}
                  </FormLabel>
                  <FormControl>
                    <Input placeholder={t('features.business.user.form.fields.phonePlaceholder')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='email'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.business.user.form.fields.email')}
                  </FormLabel>
                  <FormControl>
                    <Input type='email' placeholder={t('features.business.user.form.fields.emailPlaceholder')} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='projectId'
              render={({ field }) => (
                <FormItem className='space-y-2'>
                  <FormLabel>{t('features.business.user.form.fields.project')}
                  </FormLabel>
                  <Select
                    onValueChange={(v) => field.onChange(Number(v))}
                    value={field.value != null ? String(field.value) : ''}
                    disabled={isEdit}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder={t('features.business.user.form.fields.projectPlaceholder')} />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {projectList.map((project) => (
                        <SelectItem key={project.id} value={String(project.id)}>
                          {project.title}
                        </SelectItem>
                      ))}
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
            {t('features.business.user.form.cancel')}
          </Button>
          <Button
            type='button'
            disabled={isLoading}
            onClick={form.handleSubmit(handleSubmit)}
          >
            {isLoading ? t('features.business.user.form.submitting') : isEdit ? t('features.business.user.form.save') : t('features.business.user.form.create')}
          </Button>
        </DialogFooter>
      </DialogContent>
      <ResourceUpload
        open={avatarResourceDialogOpen}
        onOpenChange={setAvatarResourceDialogOpen}
        onSelect={(resource) => {
          form.setValue('avatarURL', resource.url, { shouldValidate: true })
          setAvatarResourceDialogOpen(false)
        }}
        type={1}
        title={t('features.business.user.form.fields.selectAvatarDialogTitle')}
        description={t('features.business.user.form.fields.selectAvatarDialogDesc')}
      />
    </Dialog>
  )
}
