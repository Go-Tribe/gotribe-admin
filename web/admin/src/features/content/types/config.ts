/** 配置项（列表/详情） */
export interface Config {
  configID: string
  alias: string
  title: string
  description: string
  type?: number
  info?: string
  mdContent?: string
  projectId?: number
  createdAt?: string
  updatedAt?: string
}

/** 列表查询参数 */
export interface ConfigListParams {
  configID?: string
  title?: string
  pageNum?: number
  pageSize?: number
  projectId?: number
  type?: number
}

/** 列表接口返回 */
export interface ConfigListResponse {
  configs: Config[]
  total: number
}

/** 新增/创建配置请求体 */
export interface ConfigCreateParams {
  title: string
  description: string
  info?: string
  projectId: number
  alias: string
  type: number
  mdContent?: string
}

/** 更新配置请求体（编辑时：不含 alias、type） */
export interface ConfigUpdateParams {
  title: string
  description: string
  projectId?: number
  info?: string
  mdContent?: string
}
