import { request } from '@/service'
import { createCrudService } from '@/lib'
import type { Admin } from '../types/admin'

/** 获取当前管理员信息；GET /api/admin/info，用于 Profile 页展示 */
export interface AdminInfoParams {
  project_id?: string
  title?: string
  page_num?: number
  page_size?: number
}

export async function getAdminInfo(
  params?: AdminInfoParams
): Promise<Admin | null> {
  const data = await request.get<{ admin: Admin }>(
    '/api/admin/info',
    { params: params ?? { page_num: 1, page_size: 3 } }
  )
  return data?.admin ?? null
}

/** 修改当前管理员密码；PUT /api/admin/changePwd */
export async function changePassword(params: {
  OldPassword: string
  NewPassword: string
}): Promise<void> {
  await request.put<void>('/api/admin/changePwd', params)
}

/**
 * 使用 CrudServiceFactory 创建的 Admin Service
 * 
 * 对比原版本：
 * - 原版本：需要手动编写 getList, create, update, delete 等方法
 * - 简化版：一行代码生成所有 CRUD 方法
 */
export const adminService = createCrudService<Admin>({
  baseUrl: '/api/admin',
  // 自定义 list 路径
  listPath: '/list',
  // 非 RESTful 风格（update 使用 POST）
  restful: false,
  // 自定义方法覆盖
  customMethods: {
    // 更新使用 patch
    update: async (id, params) => {
      return request.patch<{ success: boolean }>(`/api/admin/update/${id}`, params)
    },
    // 删除使用特殊路径
    delete: async (id) => {
      await request.delete<{ success: boolean }>('/api/admin/delete/batch', {
        data: { userIds: [id] },
      })
    },
  },
})

// 导出兼容原版本的函数
export const getAdminList = adminService.getList
export const createAdmin = adminService.create
export const updateAdmin = (params: Partial<Admin>) => adminService.update(params.id!, params) as Promise<Admin>
export const deleteAdmin = adminService.delete

export default adminService
