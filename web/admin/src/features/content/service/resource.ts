import { request } from '@/service'
import type {
  ResourceApiItem,
  ResourceItem,
  ResourceListParams,
  ResourceListResponse,
  UploadResourceApiItem,
} from '../types/resource'

/** 将接口返回项映射为前端 ResourceItem */
function mapApiItemToResource(r: ResourceApiItem): ResourceItem {
  const baseUrl = (r.url ?? '').replace(/\/+$/, '')
  const path = (r.path ?? '').replace(/^\/+/, '')
  const url = path ? `${baseUrl}/${path}` : baseUrl || ''
  return {
    id: r.id,
    url,
    name: (r.title ?? '').trim(),
    type: r.fileType ?? 0,
    size: r.size,
    createdAt: r.createdAt,
  }
}

/** 获取资源列表（分页、按类型筛选）；接口返回 data.resources、data.total */
export async function getResourceList(
  params?: ResourceListParams
): Promise<ResourceListResponse> {
  const query: Record<string, string | number | undefined> = {}
  if (params?.type != null && params.type !== 0) query.type = params.type
  if (params?.pageNum != null) query.pageNum = params.pageNum
  if (params?.pageSize != null) query.pageSize = params.pageSize
  const data = await request.get<{
    resources?: ResourceApiItem[]
    total?: number
  }>('/api/resource', { params: query })
  const resources = data?.resources ?? []
  const total = data?.total ?? 0
  return {
    list: resources.map(mapApiItemToResource),
    total,
  }
}

/** 上传资源；接口若返回与列表项相同结构则映射为 ResourceItem */
export async function uploadResource(
  file: File,
  type: number,
  onProgress?: (progress: number) => void
): Promise<UploadResourceApiItem> {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('type', type.toString())
  const data = await request.upload<{ upload?: UploadResourceApiItem } | UploadResourceApiItem>(
    '/api/resource/upload',
    formData,
    (e) => {
      if (e.total && onProgress) {
        onProgress(Math.round((e.loaded * 100) / e.total))
      }
    }
  )
  const upload = (data as { upload?: UploadResourceApiItem }).upload ?? (data as UploadResourceApiItem)
  return upload
}

/** 获取资源详情；GET /api/resource/:resourceID，返回 data.resource */
export async function getResourceDetail(resourceID: number): Promise<ResourceApiItem> {
  const data = await request.get<{ resource?: ResourceApiItem } | ResourceApiItem>(
    `/api/resource/${resourceID}`
  )
  const resource =
    (data as { resource?: ResourceApiItem }).resource ?? (data as ResourceApiItem)
  if (!resource?.id) throw new Error('资源不存在')
  return resource
}

/** 更新资源；PATCH /api/resource/:resourceID，请求体 { title, description } */
export async function updateResource(
  resourceID: number,
  body: { title: string; description?: string }
): Promise<void> {
  await request.patch(`/api/resource/${resourceID}`, body)
}

/** 删除资源；请求体为 { resourceID: string } */
export async function deleteResource(id: number): Promise<void> {
  await request.delete('/api/resource', { data: { ids: [id] } })
}
