import { useQuery } from '@tanstack/react-query'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { getUserDetail } from '../service/user'
import { useI18n } from '@/context/i18n-provider'

type UserDetailDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  userID: number | null
}

function DetailRow({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className='flex py-2 border-b border-border/50 last:border-0'>
      <span className='w-24 shrink-0 text-muted-foreground text-sm'>{label}</span>
      <span className='text-sm break-all'>{value ?? '-'}</span>
    </div>
  )
}

export function UserDetailDialog({
  open,
  onOpenChange,
  userID,
}: UserDetailDialogProps) {
  const { t } = useI18n()
  const { data: user, isLoading, error } = useQuery({
    queryKey: ['userDetail', userID],
    queryFn: () => getUserDetail(userID!),
    enabled: open && !!userID,
  })

  function formatSex(sex: string) {
    if (sex === 'M') return t('features.business.user.sex.male')
    if (sex === 'F') return t('features.business.user.sex.female')
    return t('features.business.user.sex.unknown')
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-[425px] max-h-[80vh] overflow-y-auto'>
        <DialogHeader>
          <DialogTitle>{t('features.business.user.detail.title')}</DialogTitle>
        </DialogHeader>
        {isLoading && (
          <div className='py-8 text-center text-muted-foreground text-sm'>
            {t('features.business.user.loading')}
          </div>
        )}
        {error && (
          <div className='py-8 text-center text-destructive text-sm'>
            {t('features.business.user.loadError')}
          </div>
        )}
        {user && !isLoading && (
          <div className='space-y-6'>
            <div className='flex flex-col items-center gap-2'>
              <Avatar className='h-20 w-20'>
                <AvatarImage src={user.avatarURL} alt={user.nickname} />
                <AvatarFallback className='text-lg'>
                  {user.username?.charAt(0)?.toUpperCase() || 'U'}
                </AvatarFallback>
              </Avatar>
              <div className='text-center'>
                <div className='font-semibold text-lg'>{user.nickname}</div>
                <div className='text-sm text-muted-foreground'>@{user.username}</div>
              </div>
            </div>

            <div className='rounded-lg border p-4 space-y-0'>
              <DetailRow label={t('features.business.user.detail.userID')} value={user.id} />
              <DetailRow label={t('features.business.user.detail.email')} value={user.email} />
              <DetailRow label={t('features.business.user.detail.phone')} value={user.phone} />
              <DetailRow label={t('features.business.user.detail.sex')} value={formatSex(user.sex || '')} />
              <DetailRow label={t('features.business.user.detail.projectID')} value={user.projectId} />
              <DetailRow label={t('features.business.user.detail.status')} value={user.status} />
              <DetailRow label={t('features.business.user.detail.birthday')} value={user.birthday} />
              <DetailRow label={t('features.business.user.detail.point')} value={user.point} />
              <DetailRow label={t('features.business.user.detail.createdAt')} value={user.createdAt} />
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  )
}
