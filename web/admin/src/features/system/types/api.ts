export interface Api {
  id: number
  createdAt: string
  updatedAt: string
  deletedAt: string | null
  method: string
  path: string
  category: string
  desc: string
  creator: string
}