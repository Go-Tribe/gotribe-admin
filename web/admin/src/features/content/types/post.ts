/** 文章（接口 /api/post 返回项） */
export interface Post {
  id: number
  slug: string
  title: string
  description: string
  categoryID: number
  projectId: number
  userID: number
  author: string
  content: string
  htmlContent: string
  ext: string
  icon: string
  tag: string
  type: number
  isTop: number
  isPasswd: number
  category: unknown | null
  tags: unknown | null
  project: unknown | null
  createdAt: string
  status: number
  location: string
  people: string
  time: string
  images: string[] | null
  unitPrice: number
  video: string
  password?: string
  showTime?: string
}

/** 文章列表查询参数 */
export interface PostListParams {
  id?: number
  title?: string
  status?: string
  projectId?: number
  pageNum?: number
  pageSize?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

/** 创建/更新文章请求参数 */
export interface PostParams {
  title: string
  slug?: string
  description?: string
  author?: string
  userID?: number
  content?: string
  status?: number
  icon?: string
  categoryID?: number
  isTop?: number
  tag?: string
  // 新增字段
  type?: number
  projectId?: number
  isPasswd?: number
  video?: string
  images?: string[]
  password?: string
  htmlContent?: string
  showTime?: string
  /** 自定义字段 JSON 字符串，存 key-value 对象 */
  ext?: string
}

/** 文章列表响应 */
export interface PostListResponse {
  posts: Post[]
  total: number
}
