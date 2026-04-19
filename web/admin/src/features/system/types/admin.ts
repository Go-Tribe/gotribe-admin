export type Admin = {
  id: number
  username: string
  mobile: string
  avatar: string
  nickname: string
  introduction: string
  status: number // 1: 启用, 0: 禁用
  creator: string
  roleIds: number[]
}

export type AdminListResponse = {
  admins: Admin[]
  total: number
}

export type Role = {
  id: number
  name: string
  keyword: string
  desc: string
  status: number
  sort: number
  creator: string
  createdAt: string
  updatedAt: string
  deletedAt: string | null
  admins: unknown[] | null
  menus: unknown[] | null
}
