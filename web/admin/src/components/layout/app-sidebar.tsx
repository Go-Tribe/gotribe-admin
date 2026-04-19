import { useLayout } from '@/context/layout-provider'
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarRail,
  SidebarFooter
} from '@/components/ui/sidebar'
import { NavGroup } from './nav-group'
import { NavUser } from './nav-user'
import { SidebarBrand } from './sidebar-brand'
import { useAuthUser } from '@/stores/auth-store'
import { useCachedFetch } from '@/hooks/use-cached-fetch'
import { getMenuAccessTree, type MenuItem } from './service'
import { Skeleton } from '@/components/ui/skeleton'

/**
 * 菜单骨架屏
 */
function MenuSkeleton() {
  return (
    <div className="space-y-2 p-4">
      {[...Array(5)].map((_, i) => (
        <Skeleton key={i} className="h-8 w-full" />
      ))}
    </div>
  )
}

/**
 * 优化的应用侧边栏
 * 
 * 优化点：
 * 1. 使用 useAuthUser 只订阅需要的 user 字段
 * 2. 使用 useCachedFetch 缓存菜单数据
 * 3. 页面可见性变化时有冷却时间，避免频繁刷新
 */
export function AppSidebar() {
  const { collapsible, variant } = useLayout()
  const user = useAuthUser()
  const userId = user?.id

  // 使用缓存的数据获取 hook
  const { data: menuItems, isLoading } = useCachedFetch<MenuItem[]>(
    userId ? `menu-${userId}` : null,
    async () => {
      const res = await getMenuAccessTree(userId as number)
      return res?.menuTree || []
    },
    {
      cacheTime: 10 * 60 * 1000, // 缓存 10 分钟
      debounceMs: 200,
      refetchOnVisible: true,
      refetchCooldown: 60 * 1000, // 冷却 60 秒
      initialData: [],
    }
  )

  return (
    <Sidebar collapsible={collapsible} variant={variant}>
      <SidebarHeader>
        <SidebarBrand />
      </SidebarHeader>
      <SidebarContent>
        {isLoading ? (
          <MenuSkeleton />
        ) : (
          menuItems && menuItems.length > 0 && <NavGroup items={menuItems} />
        )}
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
