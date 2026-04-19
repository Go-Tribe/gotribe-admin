import { useState, useEffect, useCallback, useRef } from 'react'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'
import { Search, X } from 'lucide-react'
import { Button } from '@/components/ui/button'

export interface DebouncedInputProps
  extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'onChange'> {
  /** 当前值 */
  value: string
  /** 值变化回调（防抖后触发） */
  onChange: (value: string) => void
  /** 防抖延迟（毫秒） */
  delay?: number
  /** 是否显示搜索图标 */
  showSearchIcon?: boolean
  /** 是否显示清除按钮 */
  showClearButton?: boolean
  /** 自定义类名 */
  wrapperClassName?: string
  /** 占位符 */
  placeholder?: string
}

/**
 * 防抖输入组件
 * 
 * 性能优化：减少频繁输入时的回调触发，降低服务器压力
 * 
 * @example
 * ```tsx
 * // 基础用法
 * <DebouncedInput
 *   value={searchValue}
 *   onChange={setSearchValue}
 *   placeholder="搜索用户..."
 * />
 * 
 * // 表格筛选
 * <DebouncedInput
 *   value={getFilterValue('username')}
 *   onChange={(value) => setColumnFilters([{ id: 'username', value }])}
 *   delay={300}
 *   showSearchIcon
 *   showClearButton
 * />
 * ```
 */
export function DebouncedInput({
  value,
  onChange,
  delay = 300,
  showSearchIcon = true,
  showClearButton = true,
  className,
  wrapperClassName,
  placeholder,
  ...props
}: DebouncedInputProps) {
  // 内部状态用于即时显示输入
  const [inputValue, setInputValue] = useState(value)
  const timeoutRef = useRef<NodeJS.Timeout | null>(null)

  // 外部值变化时同步内部值
  useEffect(() => {
    setInputValue(value)
  }, [value])

  // 清理定时器
  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }
    }
  }, [])

  const handleChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const newValue = e.target.value
      setInputValue(newValue)

      // 清除之前的定时器
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }

      // 设置新的定时器
      timeoutRef.current = setTimeout(() => {
        onChange(newValue)
      }, delay)
    },
    [onChange, delay]
  )

  const handleClear = useCallback(() => {
    setInputValue('')
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
    }
    onChange('')
  }, [onChange])

  return (
    <div className={cn('relative', wrapperClassName)}>
      {showSearchIcon && (
        <Search className="absolute left-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
      )}
      <Input
        type="text"
        value={inputValue}
        onChange={handleChange}
        placeholder={placeholder}
        className={cn(
          showSearchIcon && 'pl-9',
          showClearButton && inputValue && 'pr-9',
          className
        )}
        {...props}
      />
      {showClearButton && inputValue && (
        <Button
          variant="ghost"
          size="sm"
          className="absolute right-0 top-1/2 h-full -translate-y-1/2 px-2 py-0 hover:bg-transparent"
          onClick={handleClear}
        >
          <X className="h-4 w-4 text-muted-foreground hover:text-foreground" />
        </Button>
      )}
    </div>
  )
}

/**
 * 表格筛选专用防抖输入
 * 
 * 简化的 API，专为 DataTable 筛选设计
 */
export interface TableFilterInputProps {
  columnId: string
  placeholder?: string
  value: string
  onChange: (value: string) => void
  delay?: number
  className?: string
}

export function TableFilterInput({
  columnId: _columnId,
  placeholder,
  value,
  onChange,
  delay = 300,
  className,
}: TableFilterInputProps) {
  return (
    <DebouncedInput
      value={value}
      onChange={onChange}
      delay={delay}
      placeholder={placeholder}
      showSearchIcon
      showClearButton
      className={cn('h-8 w-[150px] lg:w-[250px]', className)}
    />
  )
}
