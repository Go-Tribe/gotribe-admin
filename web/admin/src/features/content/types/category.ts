/** 分类类型（树节点含可选 children） */
export interface Category {
  id: number
  categoryID: string
  title: string
  slug?: string
  description: string
  icon: string
  /** 链接，列表展示优先用此字段 */
  path?: string
  route: string
  sort: number
  status?: number
  hidden: number // 1: 显示, 2: 隐藏
  /** 父分类 ID，编辑时用于回显（接口可能返回 parentID 或 parentId） */
  parentId?: number
  parentID?: number
  createdAt: string
  updatedAt: string
  children?: Category[]
}

/** 创建/更新分类请求参数 */
export interface CategoryParams {
  id?: number
  title: string
  slug: string
  description?: string
  icon?: string
  /** 路由地址，创建/更新时接口使用 path 字段 */
  route?: string
  path?: string
  sort?: number
  status?: number
  hidden?: number
  /** 父分类 ID，新建时必传，0 表示根分类；更新时接口字段为 parentID */
  parentId?: number
}

/** 分类列表查询参数 */
export interface CategoryListParams {
  pageNum?: number
  pageSize?: number
  title?: string
  hidden?: string
}

/** 分类列表响应 */
export interface CategoryListResponse {
  categoryTree: Category[]
  total: number
}

/** 删除请求参数（接口要求 body: { categoryIds: string }） */
export interface BatchDeleteParams {
  ids: number[]
}
