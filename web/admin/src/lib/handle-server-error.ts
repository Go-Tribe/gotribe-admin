import { AxiosError } from 'axios'
import { toast } from 'sonner'

/**
 * 处理服务器错误
 * 提供多级错误消息提取，确保始终显示有意义的错误信息
 * @param error 错误对象
 */
export function handleServerError(error: unknown): void {
  // 使用 console.error 记录错误，便于调试
  // eslint-disable-next-line no-console
  console.error('Server error:', error)

  let errMsg = 'Something went wrong!'

  // 处理特殊状态码
  if (
    error &&
    typeof error === 'object' &&
    'status' in error &&
    Number(error.status) === 204
  ) {
    errMsg = 'Content not found.'
  }

  // 处理 AxiosError
  if (error instanceof AxiosError) {
    const data: unknown = error.response?.data

    if (data !== null && data !== undefined) {
      if (typeof data === 'string') {
        errMsg = data
      } else if (typeof data === 'object') {
        const obj = data as Record<string, unknown>
        const msg =
          (typeof obj.message === 'string' && obj.message) ||
          (typeof obj.msg === 'string' && obj.msg) ||
          (typeof obj.error === 'string' && obj.error) ||
          (typeof obj.errorMessage === 'string' && obj.errorMessage) ||
          (typeof obj.code !== 'undefined' && typeof obj.message === 'string' ? obj.message : null)
        if (msg) errMsg = msg
      }
    }

    // 3. 如果仍未提取到有效信息（即仍为默认值），或者是网络层面的默认消息，则回退到 error.message
    if (errMsg === 'Something went wrong!' && error.message) {
      errMsg = error.message
    }
  }

  // 处理普通 Error 对象（AxiosError 也是 Error 的子类，所以要排除 AxiosError 以免覆盖上面的逻辑）
  if (error instanceof Error && !(error instanceof AxiosError)) {
    errMsg = error.message || errMsg
  }

  // 显示错误提示
  toast.error(errMsg)
}
