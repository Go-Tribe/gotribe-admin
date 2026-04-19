export interface OperationLog {
  id: number
  createdAt: string
  updatedAt: string
  deletedAt: string | null
  username: string
  ip: string
  ipLocation: string
  method: string
  path: string
  desc: string
  status: number
  startTime: string
  timeCost: number
  userAgent: string
}
