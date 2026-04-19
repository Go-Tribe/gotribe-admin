/** 评论（列表项，与接口返回一致） */
export interface Comment {
  id?: number
  commentID?: string
  projectID?: string
  status?: number
  userID?: number
  objectID?: string
  objectType?: number
  comment?: string
  htmlContent?: string
  nickname?: string
  ip?: string
  country?: string
  regionName?: string
  city?: string
  createdAt?: string
  updatedAt?: string
  [key: string]: unknown
}

/** 列表查询参数 */
export interface CommentListParams {
  pageNum?: number
  pageSize?: number
  status?: string | number
  nickname?: string
  projectID?: string
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

/** 列表接口返回 data 结构 */
export interface CommentListResponse {
  comments?: Comment[]
  total?: number
}
