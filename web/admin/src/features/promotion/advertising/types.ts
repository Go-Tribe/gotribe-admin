/** 广告（列表项，与接口返回一致） */
export interface Ad {
  id?: number
  title?: string
  description?: string
  sceneID?: number
  SceneTitle?: string
  status?: number
  image?: string
  video?: string
  sort?: number
  url?: string
  urlType?: number
  ext?: string
  createdAt?: string
  updatedAt?: string
  [key: string]: unknown
}

/** 列表查询参数 */
export interface AdListParams {
  pageNum?: number
  pageSize?: number
  sceneID?: string
  title?: string
  status?: string | number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

/** 列表接口返回 data 结构 */
export interface AdListResponse {
  ads?: Ad[]
  total?: number
}

/** 创建广告请求参数 */
export interface AdCreateParams {
  title: string
  description: string
  sceneID: number
  status?: number
  image?: string
  video?: string
  sort?: number
  url?: string
  urlType?: number
  ext?: string
}

/** 更新广告请求参数 */
export interface AdUpdateParams {
  title?: string
  description?: string
  sceneID?: number
  status?: number
  image?: string
  video?: string
  sort?: number
  url?: string
  urlType?: number
  ext?: string
}
