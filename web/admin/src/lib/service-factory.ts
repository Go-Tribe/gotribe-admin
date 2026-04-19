import { request } from '@/service'

/** 通用响应类型 */
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

/** 列表响应类型 */
export interface ListResponse<T = unknown> {
  list: T[]
  total: number
}

/** 列表查询参数 */
export interface ListParams {
  pageNum?: number
  pageSize?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  [key: string]: unknown
}

/** CRUD Service 接口 */
export interface CrudService<T = unknown, CreateDTO = Partial<T>, UpdateDTO = Partial<T>> {
  /** 获取列表 */
  getList: (params?: ListParams) => Promise<ListResponse<T>>
  /** 获取详情 */
  getById: (id: number | string) => Promise<T>
  /** 创建 */
  create: (data: CreateDTO) => Promise<T>
  /** 更新 */
  update: (id: number | string, data: UpdateDTO) => Promise<T>
  /** 删除 */
  delete: (id: number | string) => Promise<void>
  /** 批量删除 */
  batchDelete?: (ids: (number | string)[]) => Promise<void>
}

/** Service 配置选项 */
export interface ServiceFactoryOptions {
  /** 基础 URL */
  baseUrl: string
  /** 自定义列表路径 */
  listPath?: string
  /** 自定义详情路径 */
  detailPath?: string
  /** 是否使用 RESTful 路径 */
  restful?: boolean
  /** 自定义请求方法 */
  customMethods?: Partial<CrudService>
}

/**
 * CRUD Service 工厂函数
 * 
 * 自动生成标准的增删改查 API 方法，消除样板代码
 * 
 * @example
 * ```typescript
 * // 基础用法
 * const userService = createCrudService<User>({
 *   baseUrl: '/api/user'
 * })
 * 
 * // 使用
 * const { list, total } = await userService.getList({ pageNum: 1, pageSize: 10 })
 * const user = await userService.create({ name: '张三', email: 'zhangsan@example.com' })
 * await userService.update(1, { name: '李四' })
 * await userService.delete(1)
 * 
 * // 自定义配置
 * const customService = createCrudService({
 *   baseUrl: '/api/custom',
 *   listPath: '/search',
 *   restful: false,
 *   customMethods: {
 *     getList: (params) => request.post('/api/custom/search', params)
 *   }
 * })
 * ```
 */
export function createCrudService<T = unknown, CreateDTO = Partial<T>, UpdateDTO = Partial<T>>(
  options: ServiceFactoryOptions
): CrudService<T, CreateDTO, UpdateDTO> {
  const {
    baseUrl,
    listPath = '/list',
    detailPath = '',
    restful = true,
    customMethods = {},
  } = options

  // 标准化 URL
  const normalizedBaseUrl = baseUrl.endsWith('/') ? baseUrl.slice(0, -1) : baseUrl

  // 构建详情 URL
  const buildDetailUrl = (id: number | string) => {
    if (detailPath) {
      return `${normalizedBaseUrl}${detailPath}/${id}`
    }
    return `${normalizedBaseUrl}/${id}`
  }

  const service: CrudService<T, CreateDTO, UpdateDTO> = {
    /** 获取列表 */
    getList: async (params?: ListParams) => {
      const url = `${normalizedBaseUrl}${listPath}`
      const res = await request.get<ApiResponse<ListResponse<T>>>(url, { params })
      return res.data
    },

    /** 获取详情 */
    getById: async (id: number | string) => {
      const url = buildDetailUrl(id)
      const res = await request.get<ApiResponse<T>>(url)
      return res.data
    },

    /** 创建 */
    create: async (data: CreateDTO) => {
      const res = await request.post<ApiResponse<T>>(normalizedBaseUrl, data)
      return res.data
    },

    /** 更新 */
    update: async (id: number | string, data: UpdateDTO) => {
      if (restful) {
        const url = buildDetailUrl(id)
        const res = await request.put<ApiResponse<T>>(url, data)
        return res.data
      } else {
        const res = await request.post<ApiResponse<T>>(`${normalizedBaseUrl}/update`, {
          id,
          ...data,
        })
        return res.data
      }
    },

    /** 删除 */
    delete: async (id: number | string) => {
      if (restful) {
        const url = buildDetailUrl(id)
        await request.delete(url)
      } else {
        await request.post(`${normalizedBaseUrl}/delete`, { id })
      }
    },

    /** 批量删除 */
    batchDelete: async (ids: (number | string)[]) => {
      await request.post(`${normalizedBaseUrl}/batch-delete`, { ids })
    },
  }

  // 合并自定义方法
  return { ...service, ...customMethods } as CrudService<T, CreateDTO, UpdateDTO>
}

/**
 * 创建带缓存的 Service
 * 
 * 自动缓存 GET 请求结果，提升性能
 */
export function createCachedCrudService<T = unknown, CreateDTO = Partial<T>, UpdateDTO = Partial<T>>(
  options: ServiceFactoryOptions & { cacheTime?: number }
): CrudService<T, CreateDTO, UpdateDTO> {
  const cache = new Map<string, { data: unknown; timestamp: number }>()
  const cacheTime = options.cacheTime || 60000 // 默认 1 分钟

  const service = createCrudService<T, CreateDTO, UpdateDTO>(options)

  const getCacheKey = (method: string, params?: unknown) => {
    return `${method}:${JSON.stringify(params)}`
  }

  const getCachedData = <R>(key: string): R | null => {
    const cached = cache.get(key)
    if (cached && Date.now() - cached.timestamp < cacheTime) {
      return cached.data as R
    }
    return null
  }

  const setCachedData = (key: string, data: unknown) => {
    cache.set(key, { data, timestamp: Date.now() })
  }

  const clearCache = () => {
    cache.clear()
  }

  return {
    ...service,
    getList: async (params?: ListParams) => {
      const key = getCacheKey('getList', params)
      const cached = getCachedData<ListResponse<T>>(key)
      if (cached) return cached

      const data = await service.getList(params)
      setCachedData(key, data)
      return data
    },
    getById: async (id: number | string) => {
      const key = getCacheKey('getById', { id })
      const cached = getCachedData<T>(key)
      if (cached) return cached

      const data = await service.getById(id)
      setCachedData(key, data)
      return data
    },
    create: async (data: CreateDTO) => {
      clearCache()
      return service.create(data)
    },
    update: async (id: number | string, data: UpdateDTO) => {
      clearCache()
      return service.update(id, data)
    },
    delete: async (id: number | string) => {
      clearCache()
      return service.delete(id)
    },
  }
}

export default createCrudService
