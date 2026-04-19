import { useQuery } from '@tanstack/react-query'
import { Logo } from '@/assets/logo'
import { useI18n } from '@/context/i18n-provider'
import { getConfig } from '@/shared/api'

/** 登录/注册等页卡片顶部品牌区：Logo + 系统名称（从后台配置读取） */
export function AuthBrand() {
  const { t } = useI18n()
  const { data: configData } = useQuery({
    queryKey: ['systemConfig'],
    queryFn: () => getConfig(),
    staleTime: 5 * 60 * 1000,
  })
  const systemConfig = configData?.systemConfig
  const title = systemConfig?.title?.trim() || t('features.auth.layout.appName')
  const logoUrl = systemConfig?.logo?.trim()

  return (
    <div className='flex flex-col items-center gap-3 pb-6 text-center'>
      <div className='flex h-16 w-16 shrink-0 items-center justify-center overflow-hidden rounded-2xl border border-border/60 bg-background shadow-sm'>
        {logoUrl ? (
          <img src={logoUrl} alt='' className='h-10 w-10 object-contain' />
        ) : (
          <Logo className='h-9 w-9 opacity-80' />
        )}
      </div>
      <div className='space-y-1'>
        <p className='text-xs font-semibold uppercase tracking-[0.18em] text-muted-foreground'>Admin workspace</p>
        <h1 className='text-xl font-semibold tracking-tight text-foreground'>{title}</h1>
      </div>
    </div>
  )
}
