import {
  LayoutDashboard,
  FileText,
  Users,
  Coins,
  Megaphone,
  Settings,
  UserCircle,
  Bug,
  LogIn,
} from 'lucide-react'
import { type SidebarData } from '../types'
import { Logo } from '@/assets/logo'

/**
 * 仅用于 Cmd+K 命令菜单的导航数据，仅包含实际存在的路由。
 * 侧边栏菜单来自接口 getMenuAccessTree，不依赖本文件。
 */
export const sidebarData: SidebarData = {
  team: {
    name: 'Go-Tribe',
    logo: Logo,
    plan: 'Go-Tribe',
  },
  navGroups: [
    {
      name: 'General',
      items: [
        { title: 'Dashboard', name: 'dashboard', url: '/dashboard', icon: LayoutDashboard },
        {
          title: 'Content',
          name: 'content',
          icon: FileText,
          items: [
            { title: 'Article', name: 'article', url: '/content/article' },
            { title: 'Category', name: 'category', url: '/content/category' },
            { title: 'Tag', name: 'tag', url: '/content/tag' },
            { title: 'Column', name: 'column', url: '/content/column' },
            { title: 'Data', name: 'config', url: '/content/config' },
            { title: 'Resource', name: 'resource', url: '/content/resource' },
          ],
        },
        {
          title: 'Business',
          name: 'business',
          icon: Users,
          items: [
            { title: 'User', name: 'user', url: '/business/user' },
            { title: 'Project', name: 'project', url: '/business/project' },
          ],
        },
        {
          title: 'Operation',
          icon: Coins,
          items: [
            { title: 'Point', url: '/operation/point' },
            { title: 'Comment', url: '/operation/comment' },
          ],
        },
        {
          title: 'Promotion',
          name: 'promotion',
          icon: Megaphone,
          items: [
            { title: 'Scene', name: 'scene', url: '/promotion/scene' },
            { title: 'Ad Content', name: 'advertising', url: '/promotion/advertising' },
          ],
        },
        {
          title: 'System',
          name: 'system',
          icon: Settings,
          items: [
            { title: 'Admin', name: 'admin', url: '/system/admin' },
            { title: 'Role', name: 'role', url: '/system/role' },
            { title: 'Menu', name: 'menu', url: '/system/menu' },
            { title: 'API', name: 'api', url: '/system/api' },
            { title: 'Operation Log', url: '/system/operation-log' },
            { title: 'Site Config', name: 'config_system', url: '/system/config' },
          ],
        },
      ],
    },
    {
      name: 'Account',
      items: [
        {
          title: 'Personal Center',
          icon: UserCircle,
          items: [{ title: 'Profile', url: '/personal-center' }],
        },
      ],
    },
    {
      name: 'Pages',
      items: [
        {
          title: 'Auth',
          icon: LogIn,
          items: [
            { title: 'Sign In', url: '/sign-in' },
            { title: 'Sign Up', url: '/sign-up' },
            { title: 'OTP', url: '/otp' },
          ],
        },
        {
          title: 'Errors',
          icon: Bug,
          items: [
            { title: '401 Unauthorized', url: '/401' },
            { title: '403 Forbidden', url: '/403' },
            { title: '404 Not Found', url: '/404' },
            { title: '500 Server Error', url: '/500' },
            { title: '503 Maintenance', url: '/503' },
          ],
        },
      ],
    },
  ],
}
