import { useState, type SVGProps } from 'react'
import { useTranslation } from 'react-i18next'
import { Root as Radio, Item } from '@radix-ui/react-radio-group'
import { CircleCheck, RotateCcw, Settings } from 'lucide-react'
import { IconDir } from '@/assets/custom/icon-dir'
import { IconLayoutCompact } from '@/assets/custom/icon-layout-compact'
import { IconLayoutDefault } from '@/assets/custom/icon-layout-default'
import { IconLayoutFull } from '@/assets/custom/icon-layout-full'
import { IconSidebarFloating } from '@/assets/custom/icon-sidebar-floating'
import { IconSidebarInset } from '@/assets/custom/icon-sidebar-inset'
import { IconSidebarSidebar } from '@/assets/custom/icon-sidebar-sidebar'
import { IconThemeDark } from '@/assets/custom/icon-theme-dark'
import { IconThemeLight } from '@/assets/custom/icon-theme-light'
import { IconThemeSystem } from '@/assets/custom/icon-theme-system'
import { cn } from '@/lib/utils'
import { useDirection } from '@/context/direction-provider'
import { type Collapsible, useLayout } from '@/context/layout-provider'
import { useTheme } from '@/context/theme-provider'
import { useThemeColor } from '@/context/theme-color-provider'
import { Button } from '@/components/ui/button'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet'
import { useSidebar } from './ui/sidebar'

type ConfigDrawerProps = {
  asMenuItem?: boolean
}

const P = 'components.configDrawer'

export function ConfigDrawer({ asMenuItem = false }: ConfigDrawerProps) {
  const { t } = useTranslation()
  const { setOpen } = useSidebar()
  const { resetDir } = useDirection()
  const { resetTheme } = useTheme()
  const { resetThemeColor } = useThemeColor()
  const { resetLayout } = useLayout()
  const [sheetOpen, setSheetOpen] = useState(false)

  const handleReset = () => {
    setOpen(true)
    resetDir()
    resetTheme()
    resetThemeColor()
    resetLayout()
  }

  if (asMenuItem) {
    return (
      <>
        <Settings className='size-4' />
        <Sheet open={sheetOpen} onOpenChange={setSheetOpen}>
          <SheetTrigger asChild>
            <div className='absolute inset-0' />
          </SheetTrigger>
          <SheetContent className='flex flex-col'>
            <SheetHeader className='pb-0 text-start'>
              <SheetTitle>{t(`${P}.sheetTitle`)}</SheetTitle>
              <SheetDescription id='config-drawer-description'>
                {t(`${P}.sheetDescription`)}
              </SheetDescription>
            </SheetHeader>
            <div className='space-y-6 overflow-y-auto px-4'>
              <ThemeConfig />
              <ThemeColorConfig />
              <SidebarConfig />
              <LayoutConfig />
              <DirConfig />
            </div>
            <SheetFooter className='gap-2'>
              <Button
                variant='destructive'
                onClick={handleReset}
                aria-label={t(`${P}.resetAria`)}
              >
                {t(`${P}.resetButton`)}
              </Button>
            </SheetFooter>
          </SheetContent>
        </Sheet>
      </>
    )
  }

  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button
          size='icon'
          variant='ghost'
          aria-label={t(`${P}.openAria`)}
          aria-describedby='config-drawer-description'
          className='rounded-full'
          onClick={(e) => e.stopPropagation()}
        >
          <Settings aria-hidden='true' />
        </Button>
      </SheetTrigger>
      <SheetContent className='flex flex-col'>
        <SheetHeader className='pb-0 text-start'>
          <SheetTitle>{t(`${P}.sheetTitle`)}</SheetTitle>
          <SheetDescription id='config-drawer-description'>
            {t(`${P}.sheetDescription`)}
          </SheetDescription>
        </SheetHeader>
        <div className='space-y-6 overflow-y-auto px-4'>
          <ThemeConfig />
          <ThemeColorConfig />
          <SidebarConfig />
          <LayoutConfig />
          <DirConfig />
        </div>
        <SheetFooter className='gap-2'>
          <Button
            variant='destructive'
            onClick={handleReset}
            aria-label={t(`${P}.resetAria`)}
          >
            {t(`${P}.resetButton`)}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  )
}

function SectionTitle({
  title,
  showReset = false,
  onReset,
  className,
}: {
  title: string
  showReset?: boolean
  onReset?: () => void
  className?: string
}) {
  return (
    <div
      className={cn(
        'mb-2 flex items-center gap-2 text-sm font-semibold text-muted-foreground',
        className
      )}
    >
      {title}
      {showReset && onReset && (
        <Button
          size='icon'
          variant='secondary'
          className='size-4 rounded-full'
          onClick={onReset}
        >
          <RotateCcw className='size-3' />
        </Button>
      )}
    </div>
  )
}

function RadioGroupItem({
  item,
  isTheme = false,
}: {
  item: {
    value: string
    label: string
    icon: (props: SVGProps<SVGSVGElement>) => React.ReactElement
  }
  isTheme?: boolean
}) {
  return (
    <Item
      value={item.value}
      className={cn('group outline-none', 'transition duration-200 ease-in')}
      aria-label={`Select ${item.label.toLowerCase()}`}
      aria-describedby={`${item.value}-description`}
    >
      <div
        className={cn(
          'relative rounded-[6px] ring-[1px] ring-border',
          'group-data-[state=checked]:shadow-2xl group-data-[state=checked]:ring-primary',
          'group-focus-visible:ring-2'
        )}
        role='img'
        aria-hidden='false'
        aria-label={`${item.label} option preview`}
      >
        <CircleCheck
          className={cn(
            'size-6 fill-primary stroke-white',
            'group-data-[state=unchecked]:hidden',
            'absolute top-0 right-0 translate-x-1/2 -translate-y-1/2'
          )}
          aria-hidden='true'
        />
        <item.icon
          className={cn(
            !isTheme &&
              'fill-primary stroke-primary group-data-[state=unchecked]:fill-muted-foreground group-data-[state=unchecked]:stroke-muted-foreground'
          )}
          aria-hidden='true'
        />
      </div>
      <div
        className='mt-1 text-xs'
        id={`${item.value}-description`}
        aria-live='polite'
      >
        {item.label}
      </div>
    </Item>
  )
}

function ThemeConfig() {
  const { t } = useTranslation()
  const { defaultTheme, theme, setTheme } = useTheme()
  const items = [
    { value: 'system', labelKey: `${P}.theme.system`, icon: IconThemeSystem },
    { value: 'light', labelKey: `${P}.theme.light`, icon: IconThemeLight },
    { value: 'dark', labelKey: `${P}.theme.dark`, icon: IconThemeDark },
  ]
  return (
    <div>
      <SectionTitle
        title={t(`${P}.theme.title`)}
        showReset={theme !== defaultTheme}
        onReset={() => setTheme(defaultTheme)}
      />
      <Radio
        value={theme}
        onValueChange={setTheme}
        className='grid w-full max-w-md grid-cols-3 gap-4'
        aria-label={t(`${P}.theme.ariaLabel`)}
        aria-describedby='theme-description'
      >
        {items.map((item) => (
          <RadioGroupItem
            key={item.value}
            item={{ ...item, label: t(item.labelKey) }}
            isTheme
          />
        ))}
      </Radio>
      <div id='theme-description' className='sr-only'>
        {t(`${P}.theme.description`)}
      </div>
    </div>
  )
}

const THEME_COLOR_OPTIONS: { value: 'default' | 'violet' | 'sage' | 'rose' | 'mint'; swatch: string }[] = [
  { value: 'default', swatch: 'oklch(0.208 0.042 265.755)' },
  { value: 'violet', swatch: 'hsl(280 60% 50%)' },
  { value: 'sage', swatch: 'oklch(0.635 0.062 7.954)' },
  { value: 'rose', swatch: 'oklch(0.637 0.218 5.211)' },
  { value: 'mint', swatch: 'oklch(0.531 0.09 125.252)' },
]

function ThemeColorConfig() {
  const { t } = useTranslation()
  const { themeColor, setThemeColor, resetThemeColor } = useThemeColor()
  const defaultThemeColor: 'default' | 'violet' | 'sage' | 'rose' | 'mint' = 'default'
  return (
    <div>
      <SectionTitle
        title={t(`${P}.themeColor.title`)}
        showReset={themeColor !== defaultThemeColor}
        onReset={resetThemeColor}
      />
      <Radio
        value={themeColor}
        onValueChange={(v) => setThemeColor(v as 'default' | 'violet' | 'sage' | 'rose' | 'mint')}
        className='grid w-full max-w-md grid-cols-5 gap-4'
        aria-label={t(`${P}.themeColor.ariaLabel`)}
        aria-describedby='theme-color-description'
      >
        {THEME_COLOR_OPTIONS.map((item) => (
          <Item
            key={item.value}
            value={item.value}
            className='group outline-none transition duration-200 ease-in'
            aria-label={`${t(`${P}.themeColor.ariaLabel`)}: ${t(`${P}.themeColor.${item.value}`)}`}
          >
            <div
              className={cn(
                'relative rounded-[6px] ring-[1px] ring-border p-3',
                'group-data-[state=checked]:ring-primary group-data-[state=checked]:ring-2',
                'group-focus-visible:ring-2'
              )}
            >
              <div
                className='mx-auto size-8 rounded-full shadow-inner'
                style={{ backgroundColor: item.swatch }}
              />
              <CircleCheck
                className={cn(
                  'size-5 fill-primary stroke-white',
                  'group-data-[state=unchecked]:hidden',
                  'absolute top-2 right-2'
                )}
              />
              <div className='mt-1 text-xs text-center'>{t(`${P}.themeColor.${item.value}`)}</div>
            </div>
          </Item>
        ))}
      </Radio>
      <div id='theme-color-description' className='sr-only'>
        {t(`${P}.themeColor.description`)}
      </div>
    </div>
  )
}

function SidebarConfig() {
  const { t } = useTranslation()
  const { defaultVariant, variant, setVariant } = useLayout()
  const items = [
    { value: 'inset', labelKey: `${P}.sidebar.inset`, icon: IconSidebarInset },
    { value: 'floating', labelKey: `${P}.sidebar.floating`, icon: IconSidebarFloating },
    { value: 'sidebar', labelKey: `${P}.sidebar.sidebar`, icon: IconSidebarSidebar },
  ]
  return (
    <div className='max-md:hidden'>
      <SectionTitle
        title={t(`${P}.sidebar.title`)}
        showReset={defaultVariant !== variant}
        onReset={() => setVariant(defaultVariant)}
      />
      <Radio
        value={variant}
        onValueChange={setVariant}
        className='grid w-full max-w-md grid-cols-3 gap-4'
        aria-label={t(`${P}.sidebar.ariaLabel`)}
        aria-describedby='sidebar-description'
      >
        {items.map((item) => (
          <RadioGroupItem
            key={item.value}
            item={{ ...item, label: t(item.labelKey) }}
          />
        ))}
      </Radio>
      <div id='sidebar-description' className='sr-only'>
        {t(`${P}.sidebar.description`)}
      </div>
    </div>
  )
}

function LayoutConfig() {
  const { t } = useTranslation()
  const { open, setOpen } = useSidebar()
  const { defaultCollapsible, collapsible, setCollapsible } = useLayout()

  const radioState = open ? 'default' : collapsible
  const items = [
    { value: 'default', labelKey: `${P}.layout.default`, icon: IconLayoutDefault },
    { value: 'icon', labelKey: `${P}.layout.compact`, icon: IconLayoutCompact },
    { value: 'offcanvas', labelKey: `${P}.layout.fullLayout`, icon: IconLayoutFull },
  ]

  return (
    <div className='max-md:hidden'>
      <SectionTitle
        title={t(`${P}.layout.title`)}
        showReset={radioState !== 'default'}
        onReset={() => {
          setOpen(true)
          setCollapsible(defaultCollapsible)
        }}
      />
      <Radio
        value={radioState}
        onValueChange={(v) => {
          if (v === 'default') {
            setOpen(true)
            return
          }
          setOpen(false)
          setCollapsible(v as Collapsible)
        }}
        className='grid w-full max-w-md grid-cols-3 gap-4'
        aria-label={t(`${P}.layout.ariaLabel`)}
        aria-describedby='layout-description'
      >
        {items.map((item) => (
          <RadioGroupItem
            key={item.value}
            item={{ ...item, label: t(item.labelKey) }}
          />
        ))}
      </Radio>
      <div id='layout-description' className='sr-only'>
        {t(`${P}.layout.description`)}
      </div>
    </div>
  )
}

function DirConfig() {
  const { t } = useTranslation()
  const { defaultDir, dir, setDir } = useDirection()
  const items = [
    {
      value: 'ltr',
      labelKey: `${P}.direction.ltr`,
      icon: (props: SVGProps<SVGSVGElement>) => <IconDir dir='ltr' {...props} />,
    },
    {
      value: 'rtl',
      labelKey: `${P}.direction.rtl`,
      icon: (props: SVGProps<SVGSVGElement>) => <IconDir dir='rtl' {...props} />,
    },
  ]
  return (
    <div>
      <SectionTitle
        title={t(`${P}.direction.title`)}
        showReset={defaultDir !== dir}
        onReset={() => setDir(defaultDir)}
      />
      <Radio
        value={dir}
        onValueChange={setDir}
        className='grid w-full max-w-md grid-cols-3 gap-4'
        aria-label={t(`${P}.direction.ariaLabel`)}
        aria-describedby='direction-description'
      >
        {items.map((item) => (
          <RadioGroupItem
            key={item.value}
            item={{ ...item, label: t(item.labelKey) }}
          />
        ))}
      </Radio>
      <div id='direction-description' className='sr-only'>
        {t(`${P}.direction.description`)}
      </div>
    </div>
  )
}
