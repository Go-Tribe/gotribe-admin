import { useState, useEffect, useRef } from 'react'

interface DateTimeFormat {
  date: string
  time: string
}

/**
 * 格式化日期时间
 */
function formatDateTime(): DateTimeFormat {
  const now = new Date()
  const date = now.toLocaleDateString('zh-CN', { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric',
    weekday: 'long'
  })
  const time = now.toLocaleTimeString('zh-CN', { 
    hour: '2-digit', 
    minute: '2-digit',
    second: '2-digit'
  })
  return { date, time }
}

/**
 * 优化的时间更新 Hook
 * 
 * 特性：
 * 1. 使用 setInterval 进行秒级更新（比 RAF 更节省资源）
 * 2. 页面不可见时自动暂停（减少后台资源占用）
 * 3. 页面重新可见时立即更新（避免显示旧时间）
 * 
 * @example
 * const { date, time } = useOptimizedTime()
 */
export function useOptimizedTime(): DateTimeFormat {
  const [dateTime, setDateTime] = useState<DateTimeFormat>(() => formatDateTime())
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null)

  useEffect(() => {
    // 只在页面可见时更新
    if (document.visibilityState !== 'visible') return

    // 使用 requestAnimationFrame 延迟初始更新，避免同步 setState
    let rafId = requestAnimationFrame(() => {
      setDateTime(formatDateTime())
    })

    // 每秒更新
    intervalRef.current = setInterval(() => {
      setDateTime(formatDateTime())
    }, 1000)

    // 监听页面可见性变化
    const handleVisibilityChange = () => {
      if (document.visibilityState === 'visible') {
        // 页面重新可见时立即更新时间
        rafId = requestAnimationFrame(() => {
          setDateTime(formatDateTime())
        })
        // 重启动画帧循环
        if (!intervalRef.current) {
          intervalRef.current = setInterval(() => {
            setDateTime(formatDateTime())
          }, 1000)
        }
      } else {
        // 页面不可见时暂停
        if (intervalRef.current) {
          clearInterval(intervalRef.current)
          intervalRef.current = null
        }
        if (rafId) {
          cancelAnimationFrame(rafId)
        }
      }
    }

    document.addEventListener('visibilitychange', handleVisibilityChange)

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current)
      }
      if (rafId) {
        cancelAnimationFrame(rafId)
      }
      document.removeEventListener('visibilitychange', handleVisibilityChange)
    }
  }, [])

  return dateTime
}

/**
 * 简化的定时器 Hook（仅分钟级更新）
 * 适用于不需要秒级精度的场景
 * 
 * @example
 * const { date, time } = useMinuteTime()
 */
export function useMinuteTime(): DateTimeFormat {
  const [dateTime, setDateTime] = useState<DateTimeFormat>(() => formatDateTime())

  useEffect(() => {
    // 只在页面可见时更新
    if (document.visibilityState !== 'visible') return

    const update = () => {
      setDateTime(formatDateTime())
    }

    // 计算到下一分钟的毫秒数
    const now = new Date()
    const msToNextMinute = (60 - now.getSeconds()) * 1000 - now.getMilliseconds()

    // 对齐到下一分钟开始
    const timeoutId = setTimeout(() => {
      update()
      // 然后每分钟更新一次
      const intervalId = setInterval(update, 60000)
      
      return () => clearInterval(intervalId)
    }, msToNextMinute)

    return () => clearTimeout(timeoutId)
  }, [])

  return dateTime
}
