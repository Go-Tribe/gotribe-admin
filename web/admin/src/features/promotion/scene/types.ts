/** 广告场景（列表项，与接口返回一致） */
export interface Scene {
  id: number
  title: string
  description?: string
  projectId: number
  projectTitle?: string
  createdAt?: string
  updatedAt?: string
}

/** 列表查询参数 */
export interface SceneListParams {
  pageNum?: number
  pageSize?: number
  projectId?: number
}

/** 列表接口返回 data 结构 */
export interface SceneListResponse {
  adScenes: Scene[]
  total: number
}

/** 新建广告场景请求参数 */
export interface SceneCreateParams {
  title: string
  description: string
  projectId: number
}

/** 更新广告场景请求参数（PATCH body） */
export interface SceneUpdateParams {
  title: string
  description: string
  projectId?: number
}
