/** 积分项（与接口返回一致） */
export interface PointItem {
  id?: number
  point: number
  userID: number
  reason: string
  nickname: string
  createdAt: string
  updatedAt: string
}

/** 列表接口返回 data 结构 */
export interface PointListResponse {
  points?: PointItem[]
  total?: number
}

/** 列表查询参数 */
export interface PointListParams {
  pageNum?: number
  pageSize?: number
  userID?: string
  nickname?: string
  projectID?: string
}

/** 创建积分请求参数 */
export interface PointCreateParams {
  userID: number
  projectID: string
  point: number
}
