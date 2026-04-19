import { useEffect, useCallback, useRef, useState } from 'react'

/**
 * 主题变化回调函数类型
 */
type ThemeChangeCallback = () => void

/**
 * 获取 CSS 变量值
 */
function getCssVariableValue(variableName: string): string {
  if (typeof window === 'undefined') return ''
  return getComputedStyle(document.documentElement).getPropertyValue(variableName).trim()
}

/**
 * 优化后的主题变化监听 Hook
 * 
 * 特性：
 * 1. 使用防抖避免频繁触发
 * 2. 对比 class 变化前后的值，仅在主题相关 class 变化时触发
 * 3. 自动清理和断开 observer
 * 
 * @param callback 主题变化时的回调函数
 * @param options 配置选项
 * 
 * @example
 * useThemeChange(() => {
 *   updateChartColors()
 * }, { debounceMs: 100 })
 */
export function useThemeChange(
  callback: ThemeChangeCallback,
  options: { debounceMs?: number; immediate?: boolean } = {}
): boolean {
  const { debounceMs = 50, immediate = false } = options
  const callbackRef = useRef(callback)
  const observerRef = useRef<MutationObserver | null>(null)
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const lastClassRef = useRef<string>('')
  const [isReady, setIsReady] = useState(false)

  // 保持回调引用最新
  useEffect(() => {
    callbackRef.current = callback
  }, [callback])

  // 防抖执行的回调
  const debouncedCallback = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
    }
    timeoutRef.current = setTimeout(() => {
      callbackRef.current()
    }, debounceMs)
  }, [debounceMs])

  useEffect(() => {
    if (typeof window === 'undefined') return

    const root = document.documentElement
    lastClassRef.current = root.className

    // 使用 requestAnimationFrame 延迟执行，避免同步 setState
    const rafId = requestAnimationFrame(() => {
      // 立即执行一次（可选）
      if (immediate) {
        callbackRef.current()
      }
      setIsReady(true)
    })

    // 创建优化的 MutationObserver
    observerRef.current = new MutationObserver((mutations) => {
      // 检查是否是 class 属性变化
      const classMutation = mutations.find(
        m => m.type === 'attributes' && m.attributeName === 'class'
      )

      if (!classMutation) return

      const currentClass = root.className
      const previousClass = lastClassRef.current

      // 只关注主题相关的 class 变化
      const themeClasses = ['light', 'dark', 'theme-violet', 'theme-sage', 'theme-rose', 'theme-mint']
      const hasThemeChange = themeClasses.some(cls => {
        const hadBefore = previousClass.includes(cls)
        const hasNow = currentClass.includes(cls)
        return hadBefore !== hasNow
      })

      if (hasThemeChange) {
        lastClassRef.current = currentClass
        debouncedCallback()
      }
    })

    // 开始观察
    observerRef.current.observe(root, {
      attributes: true,
      attributeFilter: ['class'],
    })

    return () => {
      cancelAnimationFrame(rafId)
      if (observerRef.current) {
        observerRef.current.disconnect()
        observerRef.current = null
      }
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }
    }
  }, [debouncedCallback, immediate])

  return isReady
}

/**
 * 获取 CSS 变量值的 Hook
 * 在主题变化时自动重新获取
 * 
 * @param variableNames - CSS 变量名数组，必须是稳定的常量
 * @example
 * const colors = useCssVariables(['--chart-1', '--chart-2'])
 */
export function useCssVariables(variableNames: string[]): string[] {
  // 使用 ref 存储变量名，避免依赖变化
  const variableNamesRef = useRef(variableNames)
  
  // 使用变量名长度创建初始空数组，避免在渲染期间访问 ref
  const [values, setValues] = useState<string[]>(
    Array(variableNames.length).fill('')
  )

  const updateValues = useCallback(() => {
    if (typeof window === 'undefined') return
    
    const newValues = variableNamesRef.current.map(name => {
      const value = getCssVariableValue(name)
      // 处理 oklch 格式或其他格式
      return value ? `var(${name})` : ''
    })
    
    setValues(newValues)
  }, [])

  // 初始获取
  useEffect(() => {
    // 使用 requestAnimationFrame 确保 DOM 已准备好
    const rafId = requestAnimationFrame(updateValues)
    return () => cancelAnimationFrame(rafId)
  }, [updateValues])

  // 监听主题变化
  useThemeChange(updateValues, { debounceMs: 50 })

  return values
}
