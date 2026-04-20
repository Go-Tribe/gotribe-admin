import { Link, useLocation } from '@tanstack/react-router'
import { ChevronRight, type LucideIcon } from 'lucide-react'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible'
import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  useSidebar,
} from '@/components/ui/sidebar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '../ui/dropdown-menu'
import { useI18n } from '@/context/i18n-provider'
import { type MenuItem } from './service'
import { getMenuIcon } from './icon-registry'

const MENU_I18N_PREFIX = 'components.layout.menu.'

/**
 * 菜单项展示文案：优先用 name 映射 i18n，无则中文用 title、英文用 name
 */
function getMenuLabel(
  item: MenuItem,
  t: (key: string) => string,
  language: string
): string {
  const key = item.name ? MENU_I18N_PREFIX + item.name : ''
  if (key) {
    const translated = t(key)
    if (translated !== key) return translated
  }
  return language === 'en' ? (item.name || item.title) : (item.title || item.name)
}

/**
 * 拼接路径：处理父子路径的拼接
 */
function joinPath(parentPath: string, childPath: string): string {
  // 如果子路径以 / 开头，则使用子路径
  if (childPath.startsWith('/')) {
    return childPath
  }
  // 否则拼接父路径和子路径
  const parent = parentPath.endsWith('/') ? parentPath.slice(0, -1) : parentPath
  return `${parent}/${childPath}`
}

type NavGroupProps = {
  items: MenuItem[]
}

/** 仅展示启用且未隐藏的菜单 */
function filterVisibleMenus(menus: MenuItem[]): MenuItem[] {
  return (menus || []).filter(
    (m) => m.status === 1 && m.hidden !== 1
  )
}

export function NavGroup({ items }: NavGroupProps) {
  const { t, language } = useI18n()
  const { state, isMobile } = useSidebar()
  const href = useLocation({ select: (location) => location.href })
  const getLabel = (item: MenuItem) => getMenuLabel(item, t, language)

  const visibleItems = filterVisibleMenus(items)
  // 对一级菜单按 sort 字段排序
  const sortedItems = [...visibleItems].sort((a, b) => a.sort - b.sort)

  return (
    <SidebarGroup>
      <SidebarMenu>
        {sortedItems.map((item) => {
          const key = `${item.id}-${item.name}`
          const IconComponent = getMenuIcon(item.icon)
          const filteredChildren = filterVisibleMenus(item.children || [])
          // 对二级菜单按 sort 字段排序
          const sortedChildren = [...filteredChildren].sort((a, b) => a.sort - b.sort)
          const hasChildren = sortedChildren.length > 0
          const visibleChildren = sortedChildren

          if (!hasChildren || visibleChildren.length === 0) {
            return (
              <SidebarMenuLink
                key={key}
                item={item}
                icon={IconComponent}
                href={href}
                getLabel={getLabel}
              />
            )
          }

          if (state === 'collapsed' && !isMobile) {
            return (
              <SidebarMenuCollapsedDropdown
                key={key}
                item={item}
                children={visibleChildren}
                icon={IconComponent}
                href={href}
                getLabel={getLabel}
              />
            )
          }

          return (
            <SidebarMenuCollapsible
              key={key}
              item={item}
              children={visibleChildren}
              icon={IconComponent}
              href={href}
              getLabel={getLabel}
            />
          )
        })}
      </SidebarMenu>
    </SidebarGroup>
  )
}

function SidebarMenuLink({
  item,
  icon: IconComponent,
  href,
  getLabel,
}: {
  item: MenuItem
  icon?: LucideIcon
  href: string
  getLabel: (item: MenuItem) => string
}) {
  const { setOpenMobile } = useSidebar()
  const isActive = checkIsActive(href, item.path)
  const label = getLabel(item)

  return (
    <SidebarMenuItem>
      <SidebarMenuButton
        asChild
        isActive={isActive}
        tooltip={label}
      >
        <Link to={item.path} onClick={() => setOpenMobile(false)}>
          {IconComponent && <IconComponent />}
          <span>{label}</span>
        </Link>
      </SidebarMenuButton>
    </SidebarMenuItem>
  )
}

function SidebarMenuCollapsible({
  item,
  children: visibleChildren,
  icon: IconComponent,
  href,
  getLabel,
}: {
  item: MenuItem
  children: MenuItem[]
  icon?: LucideIcon
  href: string
  getLabel: (item: MenuItem) => string
}) {
  const { setOpenMobile } = useSidebar()
  const isActive = checkIsActive(href, item.path, true) ||
    visibleChildren.some((child) => {
      const childPath = joinPath(item.path, child.path)
      return checkIsActive(href, childPath)
    })
  const label = getLabel(item)

  return (
    <Collapsible
      asChild
      defaultOpen={isActive}
      className='group/collapsible'
    >
      <SidebarMenuItem>
        <CollapsibleTrigger asChild>
          <SidebarMenuButton tooltip={label}>
            {IconComponent && <IconComponent />}
            <span>{label}</span>
            <ChevronRight className='ms-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90 rtl:rotate-180' />
          </SidebarMenuButton>
        </CollapsibleTrigger>
        <CollapsibleContent className='CollapsibleContent'>
          <SidebarMenuSub>
            {visibleChildren.map((child) => {
              const childPath = joinPath(item.path, child.path)
              const ChildIconComponent = getMenuIcon(child.icon)
              const childIsActive = checkIsActive(href, childPath)
              const childLabel = getLabel(child)

              return (
                <SidebarMenuSubItem key={child.id}>
                  <SidebarMenuSubButton
                    asChild
                    isActive={childIsActive}
                  >
                    <Link to={childPath} onClick={() => setOpenMobile(false)}>
                      {ChildIconComponent && <ChildIconComponent />}
                      <span>{childLabel}</span>
                    </Link>
                  </SidebarMenuSubButton>
                </SidebarMenuSubItem>
              )
            })}
          </SidebarMenuSub>
        </CollapsibleContent>
      </SidebarMenuItem>
    </Collapsible>
  )
}

function SidebarMenuCollapsedDropdown({
  item,
  children: visibleChildren,
  icon: IconComponent,
  href,
  getLabel,
}: {
  item: MenuItem
  children: MenuItem[]
  icon?: LucideIcon
  href: string
  getLabel: (item: MenuItem) => string
}) {
  const isActive = checkIsActive(href, item.path) ||
    visibleChildren.some((child) => {
      const childPath = joinPath(item.path, child.path)
      return checkIsActive(href, childPath)
    })
  const label = getLabel(item)

  return (
    <SidebarMenuItem>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <SidebarMenuButton
            tooltip={label}
            isActive={isActive}
          >
            {IconComponent && <IconComponent />}
            <span>{label}</span>
            <ChevronRight className='ms-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90' />
          </SidebarMenuButton>
        </DropdownMenuTrigger>
        <DropdownMenuContent side='right' align='start' sideOffset={4}>
          <DropdownMenuLabel>{label}</DropdownMenuLabel>
          <DropdownMenuSeparator />
          {visibleChildren.map((child) => {
            const childPath = joinPath(item.path, child.path)
            const ChildIconComponent = getMenuIcon(child.icon)
            const childIsActive = checkIsActive(href, childPath)
            const childLabel = getLabel(child)

            return (
              <DropdownMenuItem key={child.id} asChild>
                <Link
                  to={childPath}
                  className={`${childIsActive ? 'bg-secondary' : ''}`}
                >
                  {ChildIconComponent && <ChildIconComponent />}
                  <span className='max-w-52 text-wrap'>{childLabel}</span>
                </Link>
              </DropdownMenuItem>
            )
          })}
        </DropdownMenuContent>
      </DropdownMenu>
    </SidebarMenuItem>
  )
}

function checkIsActive(href: string, path: string, mainNav = false) {
  const currentPath = href.split('?')[0]

  if (mainNav) {
    // 对于主菜单，检查路径的第一段是否匹配
    const pathSegments = path.split('/').filter(Boolean)
    const currentSegments = currentPath.split('/').filter(Boolean)
    return pathSegments.length > 0 &&
           currentSegments.length > 0 &&
           currentSegments[0] === pathSegments[0]
  }

  return currentPath === path || currentPath.startsWith(`${path}/`)
}
