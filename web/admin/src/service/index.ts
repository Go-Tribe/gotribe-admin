import axios, {
  type AxiosInstance,
  type AxiosProgressEvent,
  type AxiosRequestConfig,
  type AxiosResponse,
  type InternalAxiosRequestConfig,
} from 'axios'
import { useAuthStore } from '@/stores/auth-store'
import { handleServerError } from '@/lib/handle-server-error'
import { authEventEmitter } from '@/lib/auth-events'
import { env } from '@/config/env'

/**
 * API 响应数据接口
 * 所有 API 响应都应遵循此格式
 */
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

/**
 * 创建 axios 实例
 * 使用已验证的环境变量配置
 */
const service: AxiosInstance = axios.create({
  baseURL: env.VITE_API_BASE_URL,
  timeout: 30000, // 30秒超时
  headers: {
    'Content-Type': 'application/json;charset=UTF-8',
  },
})

/**
 * 请求拦截器
 * 自动添加认证 token 到请求头
 */
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 从 store 获取 token
    const token = useAuthStore.getState().auth.accessToken

    // 如果存在 token，添加到请求头
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }

    return config
  },
  (error) => {
    // 请求错误处理
    // eslint-disable-next-line no-console
    if (import.meta.env.DEV) console.error('Request error:', error)
    return Promise.reject(error)
  }
)

/**
 * 认证相关的 API 路径
 * 这些 API 的 401 错误不应该触发清除 token 或未授权事件
 * 因为这是正常的业务错误（如登录失败）
 */
const AUTH_API_PATHS = ['/api/base/login', '/api/base/signup', '/api/base/sign-up', '/api/base/register']

/**
 * 检查请求 URL 是否是认证相关的 API
 * @param url 请求 URL
 * @returns 是否是认证相关的 API
 */
function isAuthApi(url: string): boolean {
  return AUTH_API_PATHS.some((path) => url.includes(path))
}

/**
 * 响应拦截器
 * 统一处理 API 响应和错误
 * 注意：401/403 错误只在这里清除 token 和触发事件，路由跳转由 main.tsx 中的事件监听器处理
 * 但是，认证相关的 API（登录、注册等）的 401 错误不应该触发这些操作
 */
service.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>): AxiosResponse => {
    const res = response.data
    const requestUrl = response.config.url || ''

    // 如果响应状态码不是 200，视为错误
    if (res.code && res.code !== 200) {
      // 处理业务错误码
      if (res.code === 401) {
        // token 过期或未授权，清除 token 并触发事件
        // 但是，如果是认证相关的 API（登录、注册等），不应该清除 token
        // 因为这是正常的业务错误（如密码错误），不是 token 过期
        if (!isAuthApi(requestUrl)) {
          useAuthStore.getState().auth.reset()
          // 触发未授权事件，由 main.tsx 中的监听器处理路由跳转
          authEventEmitter.emitUnauthorized()
        }
      }

      // 显示错误消息（401 错误的消息由 QueryCache 处理，这里不重复显示）
      if (res.code !== 401 && res.message) {
        handleServerError(new Error(res.message))
      }

      return Promise.reject(new Error(res.message || '请求失败')) as never
    }

    // 直接返回后端 data 字段，Promise<T> 解析为 T；非 JSON 响应（如 blob）保持原 response
    const isJsonApi =
      res &&
      typeof res === 'object' &&
      'code' in res &&
      'data' in res
    if (isJsonApi) {
      return (res as ApiResponse).data as AxiosResponse as never
    }
    return response
  },
  (error) => {
    // HTTP 状态码错误处理
    if (error.response) {
      const { status } = error.response
      const requestUrl = error.config?.url || ''

      switch (status) {
        case 401:
          // 未授权，清除 token 并触发事件
          // 但是，如果是认证相关的 API（登录、注册等），不应该清除 token
          // 因为这是正常的业务错误（如密码错误），不是 token 过期
          if (!isAuthApi(requestUrl)) {
            useAuthStore.getState().auth.reset()
            authEventEmitter.emitUnauthorized()
            // 不在这里显示错误消息，由 QueryCache 统一处理
          } else {
            // 认证 API 的 401 错误，只显示错误消息，不清除 token
            handleServerError(error)
          }
          break
        case 403:
          // 禁止访问
          authEventEmitter.emitForbidden()
          handleServerError(error)
          break
        case 404:
          // 资源不存在
          handleServerError(error)
          break
        case 500:
          // 服务器错误
          handleServerError(error)
          break
        default:
          handleServerError(error)
      }
    } else if (error.request) {
      // 请求已发出但没有收到响应
      // eslint-disable-next-line no-console
      if (import.meta.env.DEV) console.error('No response received:', error.request)
      handleServerError(new Error('网络错误，请检查网络连接'))
    } else {
      // 其他错误
      // eslint-disable-next-line no-console
      if (import.meta.env.DEV) console.error('Request setup error:', error.message)
      handleServerError(error)
    }

    return Promise.reject(error)
  }
)

/**
 * 封装的请求方法
 * 提供类型安全的 API 调用接口
 */
export const request = {
  /**
   * GET 请求
   * @param url 请求 URL
   * @param config 请求配置
   * @returns Promise<T> 返回类型化的数据
   */
  get<T = unknown>(
    url: string,
    config?: AxiosRequestConfig
  ): Promise<T> {
    return service.get<ApiResponse<T>, T>(url, config)
  },

  /**
   * POST 请求
   * @param url 请求 URL
   * @param data 请求数据
   * @param config 请求配置
   * @returns Promise<T> 返回类型化的数据
   */
  post<T = unknown>(
    url: string,
    data?: unknown,
    config?: AxiosRequestConfig
  ): Promise<T> {
    return service.post<ApiResponse<T>, T>(url, data, config)
  },

  /**
   * PUT 请求
   * @param url 请求 URL
   * @param data 请求数据
   * @param config 请求配置
   * @returns Promise<T> 返回类型化的数据
   */
  put<T = unknown>(
    url: string,
    data?: unknown,
    config?: AxiosRequestConfig
  ): Promise<T> {
    return service.put<ApiResponse<T>, T>(url, data, config)
  },

  /**
   * PATCH 请求
   * @param url 请求 URL
   * @param data 请求数据
   * @param config 请求配置
   * @returns Promise<T> 返回类型化的数据
   */
  patch<T = unknown>(
    url: string,
    data?: unknown,
    config?: AxiosRequestConfig
  ): Promise<T> {
    return service.patch<ApiResponse<T>, T>(url, data, config)
  },

  /**
   * DELETE 请求
   * @param url 请求 URL
   * @param config 请求配置
   * @returns Promise<T> 返回类型化的数据
   */
  delete<T = unknown>(
    url: string,
    config?: AxiosRequestConfig
  ): Promise<T> {
    return service.delete<ApiResponse<T>, T>(url, config)
  },

  /**
   * 上传文件
   * @param url 上传 URL
   * @param file 文件或 FormData
   * @param onUploadProgress 上传进度回调
   * @returns Promise<T> 返回类型化的数据
   */
  upload<T = unknown>(
    url: string,
    file: File | FormData,
    onUploadProgress?: (progressEvent: AxiosProgressEvent) => void
  ): Promise<T> {
    const formData = file instanceof FormData ? file : new FormData()
    if (file instanceof File) {
      formData.append('file', file)
    }

    return service.post<ApiResponse<T>, T>(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress,
    })
  },

  /**
   * 下载文件
   * @param url 下载 URL
   * @param filename 文件名（可选）
   * @param config 请求配置
   * @returns Promise<void>
   */
  download(
    url: string,
    filename?: string,
    config?: AxiosRequestConfig
  ): Promise<void> {
    return service
      .get(url, {
        ...config,
        responseType: 'blob',
      })
      .then((response: AxiosResponse<Blob>) => {
        const blob = new Blob([response.data])
        const downloadUrl = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = downloadUrl
        link.download = filename || 'download'
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(downloadUrl)
      })
  },
}

/**
 * 导出 axios 实例
 * 仅在需要直接使用 axios 功能时使用，通常应使用封装的 request 方法
 */
export default service
