/**
 * 认证相关事件系统
 * 用于在 axios 拦截器中触发路由跳转，避免直接使用 window.location
 * 这样可以保持 SPA 的路由状态，避免页面刷新
 */

type AuthEventType = 'unauthorized' | 'forbidden'

interface AuthEvent {
  type: AuthEventType
  redirect?: string
}

type AuthEventListener = (event: AuthEvent) => void

class AuthEventEmitter {
  private listeners: AuthEventListener[] = []

  /**
   * 订阅认证事件
   * @param listener 事件监听器
   * @returns 取消订阅的函数
   */
  subscribe(listener: AuthEventListener): () => void {
    this.listeners.push(listener)
    return () => {
      const index = this.listeners.indexOf(listener)
      if (index > -1) {
        this.listeners.splice(index, 1)
      }
    }
  }

  /**
   * 触发认证事件
   * @param event 事件对象
   */
  emit(event: AuthEvent): void {
    this.listeners.forEach((listener) => {
      try {
        listener(event)
      } catch (error) {
        // eslint-disable-next-line no-console
        console.error('Error in auth event listener:', error)
      }
    })
  }

  /**
   * 触发未授权事件（401）
   * @param redirect 重定向路径（可选）
   */
  emitUnauthorized(redirect?: string): void {
    this.emit({ type: 'unauthorized', redirect })
  }

  /**
   * 触发禁止访问事件（403）
   */
  emitForbidden(): void {
    this.emit({ type: 'forbidden' })
  }
}

/**
 * 全局认证事件发射器
 * 在 main.tsx 中订阅此事件来处理路由跳转
 */
export const authEventEmitter = new AuthEventEmitter()
