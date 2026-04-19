# 项目审查报告

> 审查日期: 2026-04-19
> 审查范围: 前端项目架构、安全、性能、Bug
> 技术栈: React 19 + TypeScript + Vite + TanStack Router/Query + Zustand + shadcn/ui + Tailwind CSS v4

---

## 修复状态

| # | 问题 | 状态 | 文件 |
|---|------|------|------|
| 1 | CodeSplitDialog 懒加载完全失效 | ✅ 已修复 | `src/components/code-split-dialog.tsx` |

---

## 严重问题 (Critical)

### 2. XSS - dangerouslySetInnerHTML 未净化

- **文件**: `src/components/editor/slate-editor.tsx:712, 731`
- **问题**: MathElement 组件使用 `dangerouslySetInnerHTML` 渲染 KaTeX 生成的 HTML，未对输入的 LaTeX 文本进行净化。如果攻击者控制 LaTeX 内容，可能注入恶意脚本。
- **修复**: 使用 DOMPurify 净化 `katex.renderToString` 的输出。

```tsx
import DOMPurify from 'dompurify'

// 净化后使用
dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(html) }}
```

---

### 3. XSS - URL 协议未验证

- **文件**: `src/components/editor/slate-editor.tsx:846, 850`
- **问题**: 链接元素的 `href` 和 `window.open` 直接使用用户输入的 URL，未验证协议。可能执行 `javascript:` 协议脚本。
- **修复**: 添加 URL 协议白名单验证。

```tsx
const ALLOWED_PROTOCOLS = ['http:', 'https:', 'mailto:', 'tel:']

function isSafeUrl(url: string): boolean {
  try {
    const parsed = new URL(url)
    return ALLOWED_PROTOCOLS.includes(parsed.protocol)
  } catch {
    // 相对路径允许
    return !url.includes(':')
  }
}

// link 元素
href={isSafeUrl(url) ? url : '#'}
onClick={(e) => {
  if ((e.metaKey || e.ctrlKey) && isSafeUrl(url)) {
    window.open(url, '_blank', 'noopener,noreferrer')
  }
}}
```

---

### 4. VirtualDataTable 虚拟滚动 + 分页混合反模式

- **文件**: `src/components/data-table/virtual-data-table.tsx`
- **问题**: 虚拟滚动假设所有数据在内存中，但分页只提供当前页数据。两者同时启用时逻辑冲突，大数据量时虚拟滚动失去意义。
- **修复**: 二选一
  - **方案A**: 纯前端虚拟滚动（不分页，一次性加载所有数据）
  - **方案B**: 服务端分页 + 普通表格（推荐）

---

### 5. Promise.reject 类型断言错误

- **文件**: `src/service/index.ts:105`
- **问题**: `return Promise.reject(new Error(...)) as never` 将 rejected Promise 断言为 `AxiosResponse` 类型，调用方类型与实际不符。
- **修复**: 移除 `as never` 断言，让类型系统正确推断。

```tsx
// 当前（错误）
return Promise.reject(new Error(res.message || '请求失败')) as never

// 修复后
return Promise.reject(new Error(res.message || '请求失败'))
```

---

## 高优先级问题 (High)

### 安全

#### 6. Token 存储在 localStorage

- **文件**: `src/stores/auth-store.ts:56, 71`
- **问题**: accessToken 存储在 localStorage 中，容易受到 XSS 攻击窃取。
- **修复**: 改为 httpOnly + Secure + SameSite=Strict Cookie，由后端设置。

#### 7. 缺少 CSRF 防护

- **文件**: `src/service/index.ts`
- **问题**: axios 请求未携带 CSRF token，如果后端依赖 cookie 认证则存在 CSRF 风险。
- **修复**: 添加 CSRF token 到请求头，或确保后端使用 JWT + 自定义 header 的认证方式。

---

### 架构

#### 8. 跨 Feature 耦合

- **文件**: 多处
  - `features/promotion/advertising/components/ad-form-dialog.tsx:49` 导入 `@/features/content/service/post`
  - `features/content/article-form-page.tsx:32-33` 导入 `@/features/business/service/project`
  - `features/operation/comment.tsx:29` 导入 `@/features/business/service/project`
- **问题**: Features 之间直接相互导入，破坏了模块边界。
- **修复**:
  - 建立共享层 `src/shared/` 存放跨 feature 的通用类型和 API
  - 或将共享服务提取到 `src/service/` 顶层目录
  - 使用 barrel exports 控制 feature 的公开接口

#### 9. Service Factory 使用率低

- **文件**: `src/lib/service-factory.ts`（仅 `src/features/system/service/admin-service-simple.ts` 使用）
- **问题**: 项目设计了完善的 `createCrudService` 工厂函数，但 20+ 个 service 文件全部手写重复代码。
- **修复**: 将所有标准 CRUD service 迁移到工厂模式，非标准接口使用 `customMethods` 覆盖。

#### 10. 超大组件违反单一职责原则

- **文件**:
  - `src/features/content/article-form-page.tsx` ~573 行（表单逻辑 + 编辑器集成 + 图片上传 + 标签选择）
  - `src/features/system/component/role-permission-menu-dialog.tsx` ~459 行
  - `src/features/system/menu.tsx` ~469 行
- **修复**: 拆分组件，如 article-form-page 拆分为 `ArticleFormContainer` + `ArticleEditor` + `ArticleSettings` + `CoverUploader`

---

### 性能

#### 11. DataTable 行未使用 React.memo

- **文件**: `src/components/data-table/data-table-optimized.tsx:40-61`
- **问题**: `DataTableRow` 是普通函数组件，每次父组件渲染时所有行都会重新渲染，即使数据未变化。
- **修复**: 用 `memo` 包裹 `DataTableRow` 和 `DataTableHeaderRow`。

```tsx
import { memo } from 'react'

const DataTableRow = memo(function DataTableRow<TData>({ row }: { row: Row<TData> }) {
  return (
    <TableRow data-state={row.getIsSelected() && 'selected'}>
      {row.getVisibleCells().map((cell) => (
        <TableCell key={cell.id}>
          {flexRender(cell.column.columnDef.cell, cell.getContext())}
        </TableCell>
      ))}
    </TableRow>
  )
})
```

#### 12. VirtualDataTable 滚动事件未节流

- **文件**: `src/components/data-table/virtual-data-table.tsx:100-103`
- **问题**: 滚动事件以 60fps+ 频率触发 `setScrollTop`，导致频繁重渲染。
- **修复**: 使用 `requestAnimationFrame` 或 `throttle` 节流。

```tsx
const handleScroll = useCallback((e: React.UIEvent<HTMLDivElement>) => {
  const scrollTop = e.currentTarget.scrollTop
  if (rafRef.current) cancelAnimationFrame(rafRef.current)
  rafRef.current = requestAnimationFrame(() => {
    setScrollTop(scrollTop)
  })
}, [])
```

#### 13. article-form-page 多个 useEffect 依赖缺失

- **文件**: `src/features/content/article-form-page.tsx:183-298`
- **问题**: 四个 `useEffect` 都故意排除了 `form` 依赖（有 eslint-disable），这是一种反模式。如果 `form` 引用变化，这些 effect 不会响应。
- **修复**: 使用 `form.reset` 的返回值或 `useForm` 的 `resetOptions` 来避免时序问题。

#### 14. IconPicker 每次渲染全量遍历 lucide-react

- **文件**: `src/components/icon-picker.tsx:17-53`
- **问题**: `getAllIconNames` 遍历整个 lucide-react 模块（1000+ 图标），低端设备上可能耗时 10-50ms。
- **修复**: 将图标列表提取为模块级常量。

```tsx
const ALL_ICON_NAMES = Object.keys(Icons).filter(...).sort()
// 在模块级别只执行一次
```

---

### Bug

#### 15. console.error 循环引用崩溃

- **文件**: `src/lib/handle-server-error.ts:10-12`
- **问题**: `console.error('Server error:', error)` 如果 error 对象含循环引用，会导致异常。
- **修复**: 安全打印错误信息。

```ts
import { AxiosError } from 'axios'

function safeStringifyError(error: unknown): string {
  if (error instanceof AxiosError) {
    return error.message
  }
  if (error instanceof Error) {
    return error.message
  }
  try {
    return JSON.stringify(error)
  } catch {
    return String(error)
  }
}

export function handleServerError(error: unknown): void {
  // eslint-disable-next-line no-console
  console.error('Server error:', safeStringifyError(error))
  // ...
}
```

#### 16. meta 更新异步竞态

- **文件**: `src/routes/__root.tsx:65-69`
- **问题**: `beforeLoad` 中发起异步 `updateAppConfigFromApi`，不阻塞路由但可能覆盖后续配置。
- **修复**: 在组件 mount 后执行，或使用 AbortController 取消旧请求。

#### 17. QueryCache 重复 toast

- **文件**: `src/lib/query-client.ts:64-91`
- **问题**: 401 错误时 axios 拦截器和 QueryCache 都可能触发 toast，导致重复提示。
- **修复**: 统一错误处理，避免重复 toast。

#### 18. LazyImage 重复加载

- **文件**: `src/components/lazy-image.tsx:132-147`
- **问题**: `useEffect` 依赖 `onLoad/onError` 回调函数，若父组件传递内联函数会导致重复加载图片。
- **修复**: 将 `onLoad`/`onError` 用 `useRef` 存储，不放入 effect 依赖数组。

```tsx
const onLoadRef = useRef(onLoad)
onLoadRef.current = onLoad

useEffect(() => {
  if (!isInView || !src) return
  const img = new Image()
  img.src = src
  img.onload = () => {
    setState('loaded')
    onLoadRef.current?.()
  }
  img.onerror = () => {
    setState('error')
    onErrorRef.current?.(new Error(`Failed to load image: ${src}`))
  }
}, [isInView, src]) // 移除 onLoad, onError
```

---

## 中优先级问题 (Medium)

### 架构

| # | 问题 | 文件 | 修复建议 |
|---|------|------|----------|
| 19 | 目录命名不一致 | 全局 | `component/` 统一改为 `components/` |
| 20 | 类型文件扩展名不统一 | 全局 | `.d.ts` 统一改为 `.ts`（仅声明文件保留 `.d.ts`） |
| 21 | i18n 手动导入翻译文件 | `src/i18n/config.ts` | 改用 `import.meta.glob` 自动加载 |
| 22 | routeTree.gen.ts 不应提交到 Git | `src/routeTree.gen.ts` | 添加到 `.gitignore`，构建时生成 |
| 23 | AuthUser 类型过于宽松 | `src/stores/auth-store.ts:11` | 移除 `extends Record<string, unknown>`，明确定义字段 |

### 性能

| # | 问题 | 文件 | 修复建议 |
|---|------|------|----------|
| 24 | useDebouncedQuery 频繁 JSON.stringify | `src/hooks/use-debounced-query.ts:77` | 使用 `useRef` 存储上一次的值进行比较 |
| 25 | useCachedFetch cache 无上限 | `src/hooks/use-cached-fetch.ts` | 为 cache Map 添加 LRU 淘汰机制 |
| 26 | TreeSelect 重复 filterTree | `src/components/ui/tree-select.tsx:265-270` | 合并逻辑，filteredData 计算时同时收集展开的节点 |
| 27 | JsonEditor 频繁 updateProps | `src/components/json-editor/json-editor.tsx:126-134` | 将 `onRenderMenu` 提取到组件外部或使用 `useCallback` |
| 28 | SlateEditor 直接修改 editor.children | `src/components/editor/slate-editor.tsx` | 使用 Slate API 而非直接修改数组，避免破坏内部状态 |

### 安全

| # | 问题 | 文件 | 修复建议 |
|---|------|------|----------|
| 29 | Cookie 安全属性缺失 | `src/lib/cookies.ts:33` | 添加 `Secure; SameSite=Strict` |
| 30 | updateHtmlMeta XSS 风险 | `src/lib/update-html-meta.ts` | API 返回值做 HTML 实体编码 |
| 31 | 路由权限仅依赖前端 | `src/routes/_authenticated/route.tsx` | 后端必须校验每个 API 请求的权限 |
| 32 | 缺少 CSP / X-Frame-Options | `index.html` | 添加 CSP meta 标签或响应头 |

### Bug

| # | 问题 | 文件 | 修复建议 |
|---|------|------|----------|
| 33 | JSON.parse 静默失败 | `src/features/content/components/article-media-sheet.tsx:47` | catch 中添加错误处理或 toast 提示 |
| 34 | VirtualDataTable key 使用索引 | `src/components/data-table/virtual-data-table.tsx:160` | 使用数据中稳定的唯一标识（如 ID）作为 key |
| 35 | rootElement.innerHTML 检查不可靠 | `src/main.tsx:75` | 使用 `getElementById` 后判断是否存在，或使用 hydrate 模式 |

---

## 低优先级问题 (Low)

| # | 问题 | 文件 | 修复建议 |
|---|------|------|----------|
| 36 | DashboardPage CurrentTime 每秒重渲染 | `src/features/dashboard/dashboard-page.tsx` | 替换为 `useOptimizedTime` hook，页面不可见时暂停更新 |
| 37 | lucide-react 未代码分割 | `vite.config.ts` | 将 `lucide-react` 单独分块或按需导入 |
| 38 | 开发环境代理配置硬编码 | `vite.config.ts:11` | 完全通过环境变量配置 |
| 39 | 错误日志可能泄露敏感信息 | `src/lib/handle-server-error.ts:10-12` | 生产环境禁用详细错误日志 |
| 40 | 构建配置缺少 sourcemap/minify | `vite.config.ts` | 添加 `build.sourcemap` 和 `esbuild.drop` 配置 |

---

## 正面发现

| 方面 | 说明 |
|------|------|
| 状态管理 | Zustand 选择器优化（`useAuthUser`, `useAccessToken`），避免不必要重渲染 |
| API 封装 | Axios 拦截器处理 token、错误统一，区分认证/非认证 API 的 401 处理 |
| 懒加载 | `LazyImage` 使用 Intersection Observer，`VirtualDataTable` 提供虚拟滚动能力 |
| 防抖 | `useDebouncedQuery`、`useDebouncedSearch`、`useDebouncedMutation` 减少服务器请求 |
| Hooks | `useDataTable`、`useCrudMutations` 消除大量重复代码 |
| 安全 | 外部链接使用 `rel="noopener noreferrer"`，`escapeHtml` 处理 Slate 内容 |
| 构建 | 手动 chunk 分割策略合理（react-vendor, router-vendor, ui-vendor, editor） |
| 性能 | Service Worker 注册，路由级 `lazy` + `Suspense` 代码分割 |

---

## 建议修复优先级

### 第一阶段（安全 + 核心功能）
1. ✅ CodeSplitDialog 懒加载失效（已修复）
2. XSS - dangerouslySetInnerHTML 未净化
3. XSS - URL 协议未验证
4. Token 存储在 localStorage
5. console.error 循环引用崩溃

### 第二阶段（架构 + 性能）
6. VirtualDataTable 虚拟滚动+分页混合
7. 跨 Feature 耦合
8. DataTable 行未 memo
9. VirtualDataTable 滚动未节流
10. article-form-page useEffect 依赖问题
11. Service Factory 推广

### 第三阶段（代码质量）
12. 超大组件拆分
13. 目录命名统一
14. 类型文件扩展名统一
15. i18n 自动加载
16. routeTree.gen.ts 排除
