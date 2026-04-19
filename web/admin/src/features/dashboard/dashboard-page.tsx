import { useState, useEffect, lazy, Suspense } from 'react'
import { Link } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { 
  Search,
  Plus,
  FileText,
  FileEdit,
  MessageCircle,
  Eye,
  TrendingUp,
  TrendingDown,

  Settings,
  Image,
  Tags,
  FolderOpen,
  AlertCircle,
  CheckCircle2,
  Clock,
  ChevronRight,

  Zap,
  Shield,
  RefreshCw,
  AlertTriangle
} from 'lucide-react'
import { useI18n } from '@/context/i18n-provider'
import { useAuthUser } from '@/stores/auth-store'
import { getProjectList } from '@/features/business/service/project'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { cn } from '@/lib/utils'

import { useCssVariables } from '@/hooks/use-theme-change'

// 懒加载图表组件
const VisitTrendChart = lazy(() => import('./components/visit-trend-chart'))

/** 统计卡片数据 */
const STATS = [
  { 
    label: '总文章数', 
    value: 156, 
    trend: 12, 
    icon: FileText, 
    color: 'blue',
    path: '/content/article'
  },
  { 
    label: '草稿数', 
    value: 18, 
    trend: -3, 
    icon: FileEdit, 
    color: 'amber',
    path: '/content/article'
  },
  { 
    label: '待审核评论', 
    value: 23, 
    trend: 8, 
    icon: MessageCircle, 
    color: 'rose',
    path: '/operation/comment'
  },
  { 
    label: '近7日访问量', 
    value: 2847, 
    trend: 15, 
    icon: Eye, 
    color: 'emerald',
    path: '/dashboard'
  },
] as const

/** 快捷操作 */
const QUICK_ACTIONS = [
  { label: '新建文章', icon: Plus, path: '/content/article/new', color: 'bg-blue-500 hover:bg-blue-600' },
  { label: '上传图片', icon: Image, path: '/content/resource', color: 'bg-violet-500 hover:bg-violet-600' },
  { label: '管理标签', icon: Tags, path: '/content/tag', color: 'bg-amber-500 hover:bg-amber-600' },
  { label: '查看分类', icon: FolderOpen, path: '/content/category', color: 'bg-emerald-500 hover:bg-emerald-600' },
] as const

/** 待处理提醒 */
const PENDING_ITEMS = [
  { label: '待审核文章', count: 5, path: '/content/article', priority: 'high' },
  { label: '待审核评论', count: 23, path: '/operation/comment', priority: 'high' },
  { label: '即将过期资源', count: 8, path: '/content/resource', priority: 'medium' },
  { label: '系统通知', count: 3, path: '/system/config', priority: 'low' },
] as const

/** 最近文章 */
const RECENT_ARTICLES = [
  { id: 1, title: '2024年前端技术趋势总结', status: 'published', date: '10分钟前', views: 128 },
  { id: 2, title: 'React 19 新特性详解', status: 'published', date: '1小时前', views: 256 },
  { id: 3, title: 'TypeScript 最佳实践指南', status: 'draft', date: '3小时前', views: 0 },
  { id: 4, title: 'Node.js 性能优化技巧', status: 'published', date: '昨天', views: 512 },
  { id: 5, title: 'CSS Grid 布局完全指南', status: 'reviewing', date: '昨天', views: 89 },
] as const

/** 最近评论 */
const RECENT_COMMENTS = [
  { id: 1, author: '张三', content: '写得很详细，学到了很多！', article: 'React 19 新特性详解', time: '5分钟前' },
  { id: 2, author: '李四', content: '期待更多这样的文章', article: '2024年前端技术趋势总结', time: '15分钟前' },
  { id: 3, author: '王五', content: '请问有源码吗？', article: 'TypeScript 最佳实践指南', time: '1小时前' },
] as const

/** 热门文章 */
const POPULAR_ARTICLES = [
  { id: 1, title: 'Vue 3.0 组合式 API 入门', views: 3421 },
  { id: 2, title: 'Next.js 14 新功能介绍', views: 2893 },
  { id: 3, title: 'Tailwind CSS 实战技巧', views: 2156 },
] as const

/** 系统状态 */
const SYSTEM_STATUS = [
  { label: '系统运行时间', value: '15天 3小时', icon: Clock, status: 'normal' },
  { label: '数据库连接', value: '正常', icon: CheckCircle2, status: 'normal' },
  { label: '缓存状态', value: '已启用', icon: Zap, status: 'normal' },
  { label: '备份状态', value: '3小时前', icon: Shield, status: 'warning' },
] as const

/** SEO 提醒 */
const SEO_ALERTS = [
  { type: 'warning', message: '3篇文章缺少 meta description' },
  { type: 'info', message: '2张图片未添加 alt 属性' },
  { type: 'success', message: '站点地图已更新' },
] as const

/** CSS 变量名 - 图表颜色 */
const CHART_CSS_VARS = [
  '--chart-1',
  '--chart-2', 
  '--chart-3',
  '--chart-4',
  '--chart-5',
]

function getColorClasses(color: string) {
  const colors: Record<string, { bg: string; text: string; border: string }> = {
    blue: {
      bg: 'bg-blue-500/10',
      text: 'text-blue-600',
      border: 'border-blue-200/60',
    },
    emerald: {
      bg: 'bg-emerald-500/10',
      text: 'text-emerald-600',
      border: 'border-emerald-200/60',
    },
    amber: {
      bg: 'bg-amber-500/10',
      text: 'text-amber-600',
      border: 'border-amber-200/60',
    },
    rose: {
      bg: 'bg-rose-500/10',
      text: 'text-rose-600',
      border: 'border-rose-200/60',
    },
  }
  return colors[color] || colors.blue
}

function getPriorityColor(priority: string) {
  switch (priority) {
    case 'high':
      return 'bg-rose-500 text-white'
    case 'medium':
      return 'bg-amber-500 text-white'
    default:
      return 'bg-slate-400 text-white'
  }
}

function getStatusBadge(status: string) {
  switch (status) {
    case 'published':
      return <Badge variant="default" className="bg-emerald-500 hover:bg-emerald-600">已发布</Badge>
    case 'draft':
      return <Badge variant="secondary">草稿</Badge>
    case 'reviewing':
      return <Badge variant="outline" className="border-amber-500 text-amber-600">审核中</Badge>
    default:
      return <Badge variant="outline">未知</Badge>
  }
}

/** 当前时间组件 */
function CurrentTime() {
  const [time, setTime] = useState(new Date())
  
  useEffect(() => {
    const timer = setInterval(() => setTime(new Date()), 1000)
    return () => clearInterval(timer)
  }, [])
  
  const timeStr = time.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
  })
  
  const dateStr = time.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
    weekday: 'short'
  })
  
  return (
    <span className="text-muted-foreground/80 ml-2 text-xs tabular-nums">
      · {dateStr} {timeStr}
    </span>
  )
}

export function DashboardPage() {
  const { t: _t } = useI18n()
  const user = useAuthUser()
  const chartColors = useCssVariables(CHART_CSS_VARS)
  const [mounted, setMounted] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
  
  useEffect(() => {
    const timer = setTimeout(() => setMounted(true), 100)
    return () => clearTimeout(timer)
  }, [])

  const { data: _projectData } = useQuery({
    queryKey: ['projectList', { pageNum: 1, pageSize: 500 }],
    queryFn: () => getProjectList({ pageNum: 1, pageSize: 500 }),
  })

  if (!mounted) {
    return <DashboardSkeleton />
  }

  return (
    <div className="min-h-full p-6 space-y-6">
      {/* 顶部 Header */}
      <header className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">控制台</h1>
          <p className="text-sm text-muted-foreground">
            欢迎回来，{user?.nickname || user?.username || '管理员'}
            <CurrentTime />
          </p>
        </div>
        
        <div className="flex items-center gap-3">
          {/* 搜索框 */}
          <div className="relative w-full md:w-80">
            <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              type="search"
              placeholder="搜索文章、标签..."
              className="pl-9"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>
          
          {/* 新建文章按钮 */}
          <Link to="/content/article/new">
            <Button className="gap-2">
              <Plus className="size-4" />
              <span className="hidden sm:inline">新建文章</span>
            </Button>
          </Link>
          
          {/* 用户菜单 */}
          <Avatar className="size-9 cursor-pointer">
            <AvatarImage src={user?.avatar as string} />
            <AvatarFallback>{(user?.nickname || user?.username || 'A').charAt(0)}</AvatarFallback>
          </Avatar>
        </div>
      </header>

      {/* 第一排：统计卡片 */}
      <section className="grid grid-cols-2 gap-4 lg:grid-cols-4">
        {STATS.map(({ label, value, trend, icon: Icon, color, path }) => {
          const colors = getColorClasses(color)
          return (
            <Link key={label} to={path}>
              <Card className="transition-all hover:shadow-md hover:-translate-y-0.5">
                <CardContent className="p-4">
                  <div className="flex items-start justify-between">
                    <div className="space-y-2">
                      <p className="text-xs text-muted-foreground">{label}</p>
                      <p className="text-2xl font-bold">{value.toLocaleString()}</p>
                      <div className="flex items-center gap-1 text-xs">
                        {trend > 0 ? (
                          <>
                            <TrendingUp className="size-3 text-emerald-500" />
                            <span className="text-emerald-500">+{trend}%</span>
                          </>
                        ) : (
                          <>
                            <TrendingDown className="size-3 text-rose-500" />
                            <span className="text-rose-500">{trend}%</span>
                          </>
                        )}
                      </div>
                    </div>
                    <div className={cn('rounded-lg p-2', colors.bg, colors.text)}>
                      <Icon className="size-5" />
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>
          )
        })}
      </section>

      {/* 第二排：访问趋势 + 快捷操作/待处理 */}
      <section className="grid gap-4 lg:grid-cols-3">
        {/* 左：访问趋势图 */}
        <Card className="lg:col-span-2">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <div>
              <CardTitle className="text-base font-medium">访问趋势</CardTitle>
              <CardDescription>近7日网站访问统计</CardDescription>
            </div>
            <Badge variant="secondary">本周</Badge>
          </CardHeader>
          <CardContent>
            <Suspense fallback={<Skeleton className="h-[280px] w-full" />}>
              <VisitTrendChart chartColors={chartColors} />
            </Suspense>
          </CardContent>
        </Card>

        {/* 右：快捷操作 + 待处理提醒 */}
        <div className="space-y-4">
          {/* 快捷操作 */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base font-medium">快捷操作</CardTitle>
            </CardHeader>
            <CardContent className="grid grid-cols-2 gap-2">
              {QUICK_ACTIONS.map(({ label, icon: Icon, path, color }) => (
                <Link key={label} to={path}>
                  <Button 
                    variant="secondary" 
                    className={cn("w-full justify-start gap-2 text-white", color)}
                  >
                    <Icon className="size-4" />
                    {label}
                  </Button>
                </Link>
              ))}
            </CardContent>
          </Card>

          {/* 待处理提醒 */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base font-medium">待处理提醒</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              {PENDING_ITEMS.map(({ label, count, path, priority }) => (
                <Link 
                  key={label} 
                  to={path}
                  className="flex items-center justify-between rounded-lg border p-3 transition-colors hover:bg-muted/50"
                >
                  <div className="flex items-center gap-3">
                    <span className={cn("flex size-5 items-center justify-center rounded-full text-[10px] font-bold", getPriorityColor(priority))}>
                      {count}
                    </span>
                    <span className="text-sm">{label}</span>
                  </div>
                  <ChevronRight className="size-4 text-muted-foreground" />
                </Link>
              ))}
            </CardContent>
          </Card>
        </div>
      </section>

      {/* 第三排：最近文章 + 最近评论/热门文章 */}
      <section className="grid gap-4 lg:grid-cols-3">
        {/* 左：最近文章 */}
        <Card className="lg:col-span-2">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <div>
              <CardTitle className="text-base font-medium">最近文章</CardTitle>
            </div>
            <Link to="/content/article">
              <Button variant="ghost" size="sm" className="gap-1">
                查看全部
                <ChevronRight className="size-4" />
              </Button>
            </Link>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {RECENT_ARTICLES.map((article) => (
                <div 
                  key={article.id} 
                  className="flex items-center justify-between rounded-lg border p-3 transition-colors hover:bg-muted/30"
                >
                  <div className="flex items-center gap-3">
                    {getStatusBadge(article.status)}
                    <span className="text-sm font-medium">{article.title}</span>
                  </div>
                  <div className="flex items-center gap-4 text-xs text-muted-foreground">
                    <span className="flex items-center gap-1">
                      <Eye className="size-3" />
                      {article.views}
                    </span>
                    <span>{article.date}</span>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* 右：最近评论 + 热门文章 */}
        <div className="space-y-4">
          {/* 最近评论 */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base font-medium">最近评论</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              {RECENT_COMMENTS.map((comment) => (
                <div key={comment.id} className="space-y-1 rounded-lg border p-3">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium">{comment.author}</span>
                    <span className="text-xs text-muted-foreground">{comment.time}</span>
                  </div>
                  <p className="text-xs text-muted-foreground line-clamp-2">{comment.content}</p>
                  <p className="text-xs text-blue-600 truncate">《{comment.article}》</p>
                </div>
              ))}
            </CardContent>
          </Card>

          {/* 热门文章 */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base font-medium">热门文章</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              {POPULAR_ARTICLES.map((article, index) => (
                <div 
                  key={article.id} 
                  className="flex items-center justify-between rounded-lg p-2 transition-colors hover:bg-muted/30"
                >
                  <div className="flex items-center gap-2">
                    <span className={cn(
                      "flex size-5 items-center justify-center rounded text-[10px] font-bold",
                      index === 0 ? "bg-amber-500 text-white" : 
                      index === 1 ? "bg-slate-400 text-white" : 
                      index === 2 ? "bg-orange-400 text-white" : "bg-slate-200 text-slate-600"
                    )}>
                      {index + 1}
                    </span>
                    <span className="text-sm truncate max-w-[140px]">{article.title}</span>
                  </div>
                  <span className="text-xs text-muted-foreground">{article.views.toLocaleString()}</span>
                </div>
              ))}
            </CardContent>
          </Card>
        </div>
      </section>

      {/* 第四排：系统状态 / 备份 / 缓存 / SEO 提醒 */}
      <section className="grid gap-4 lg:grid-cols-3">
        {/* 系统状态 */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-medium">系统状态</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {SYSTEM_STATUS.map((item) => (
              <div key={item.label} className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <item.icon className={cn(
                    "size-4",
                    item.status === 'normal' ? "text-emerald-500" : "text-amber-500"
                  )} />
                  <span className="text-sm">{item.label}</span>
                </div>
                <span className={cn(
                  "text-sm font-medium",
                  item.status === 'normal' ? "text-emerald-600" : "text-amber-600"
                )}>
                  {item.value}
                </span>
              </div>
            ))}
          </CardContent>
          <CardFooter className="border-t pt-3">
            <Button variant="outline" size="sm" className="w-full gap-2">
              <RefreshCw className="size-4" />
              立即备份
            </Button>
          </CardFooter>
        </Card>

        {/* 缓存管理 */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-medium">缓存管理</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <div className="flex items-center justify-between text-sm">
                <span className="text-muted-foreground">已使用缓存</span>
                <span className="font-medium">128 MB</span>
              </div>
              <div className="h-2 w-full rounded-full bg-muted overflow-hidden">
                <div className="h-full w-[45%] rounded-full bg-blue-500" />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-2">
              <Button variant="outline" size="sm" className="gap-2">
                <Zap className="size-4" />
                清理缓存
              </Button>
              <Button variant="outline" size="sm" className="gap-2">
                <Settings className="size-4" />
                设置
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* SEO 提醒 */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-medium">SEO 提醒</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            {SEO_ALERTS.map((alert, index) => (
              <div 
                key={index} 
                className={cn(
                  "flex items-start gap-2 rounded-lg p-2",
                  alert.type === 'warning' ? "bg-amber-500/10" :
                  alert.type === 'success' ? "bg-emerald-500/10" :
                  "bg-blue-500/10"
                )}
              >
                {alert.type === 'warning' ? <AlertTriangle className="size-4 text-amber-500 mt-0.5" /> :
                 alert.type === 'success' ? <CheckCircle2 className="size-4 text-emerald-500 mt-0.5" /> :
                 <AlertCircle className="size-4 text-blue-500 mt-0.5" />}
                <span className="text-xs">{alert.message}</span>
              </div>
            ))}
          </CardContent>
          <CardFooter className="border-t pt-3">
            <Link to="/system/config" className="w-full">
              <Button variant="outline" size="sm" className="w-full">
                SEO 设置
              </Button>
            </Link>
          </CardFooter>
        </Card>
      </section>
    </div>
  )
}

/** 骨架屏加载状态 */
export function DashboardSkeleton() {
  return (
    <div className="min-h-full p-6 space-y-6">
      {/* Header Skeleton */}
      <div className="flex items-center justify-between">
        <div className="space-y-2">
          <Skeleton className="h-8 w-32" />
          <Skeleton className="h-4 w-48" />
        </div>
        <div className="flex items-center gap-3">
          <Skeleton className="h-10 w-80" />
          <Skeleton className="h-10 w-28" />
          <Skeleton className="h-9 w-9 rounded-full" />
        </div>
      </div>

      {/* Stats Skeleton */}
      <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
        {[...Array(4)].map((_, i) => (
          <Skeleton key={i} className="h-24 w-full rounded-lg" />
        ))}
      </div>

      {/* Charts + Quick Actions Skeleton */}
      <div className="grid gap-4 lg:grid-cols-3">
        <Skeleton className="h-[360px] w-full rounded-lg lg:col-span-2" />
        <div className="space-y-4">
          <Skeleton className="h-48 w-full rounded-lg" />
          <Skeleton className="h-48 w-full rounded-lg" />
        </div>
      </div>

      {/* Recent Articles + Comments Skeleton */}
      <div className="grid gap-4 lg:grid-cols-3">
        <Skeleton className="h-[400px] w-full rounded-lg lg:col-span-2" />
        <div className="space-y-4">
          <Skeleton className="h-48 w-full rounded-lg" />
          <Skeleton className="h-48 w-full rounded-lg" />
        </div>
      </div>
    </div>
  )
}
