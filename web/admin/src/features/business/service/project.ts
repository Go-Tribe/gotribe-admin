import { request } from '@/service'
import type { Project, ProjectListResponse } from '../types/project'

// 获取项目列表（支持分页和筛选）
export const getProjectList = async (params?: {
  current?: number
  title?: string
  projectID?: string
  pageNum?: number
  pageSize?: number
}) => {
  return request.get<ProjectListResponse>('/api/project', { params })
}

// 获取单个项目
export async function getProject(id: string): Promise<Project> {
  return request.get<Project>(`/api/project/${id}`)
}

// 创建项目
export async function createProject(
  params: Partial<Project>,
): Promise<{ success: boolean }> {
  return request.post<{ success: boolean }>('/api/project', params)
}

// 更新项目
export async function updateProject(
  params: Partial<Project>,
): Promise<{ success: boolean }> {
  return request.patch<{ success: boolean }>(
    `/api/project/${params.projectID}`,
    params,
  )
}

// 删除项目
export async function deleteProject(
  projectID: string,
): Promise<{ success: boolean }> {
  return request.delete<{ success: boolean }>('/api/project', {
    data: { projectIds: projectID },
  })
}
