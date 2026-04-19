import { memo } from 'react'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { DebouncedInput } from '@/components/debounced-input'
import { cn } from '@/lib/utils'
import { Search, X } from 'lucide-react'
import { Button } from '@/components/ui/button'

/** 选项类型 */
export interface FilterOption {
  label: string
  value: string | number
}

/** 文本筛选 */
export interface TextFilterProps {
  value: string
  onChange: (value: string) => void
  placeholder?: string
  className?: string
  debounce?: boolean
  debounceMs?: number
}

export const TextFilter = memo(function TextFilter({
  value,
  onChange,
  placeholder = '搜索...',
  className,
  debounce = true,
  debounceMs = 300,
}: TextFilterProps) {
  if (debounce) {
    return (
      <DebouncedInput
        value={value}
        onChange={onChange}
        delay={debounceMs}
        placeholder={placeholder}
        showSearchIcon
        showClearButton
        className={cn('h-8 w-[150px] lg:w-[200px]', className)}
      />
    )
  }

  return (
    <div className="relative">
      <Search className="absolute left-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
      <Input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        className={cn('h-8 w-[150px] pl-9 lg:w-[200px]', className)}
      />
    </div>
  )
})

/** 下拉筛选 */
export interface SelectFilterProps {
  value: string
  onChange: (value: string) => void
  options: FilterOption[]
  placeholder?: string
  className?: string
  allowClear?: boolean
}

export const SelectFilter = memo(function SelectFilter({
  value,
  onChange,
  options,
  placeholder = '全部',
  className,
  allowClear = true,
}: SelectFilterProps) {
  return (
    <Select value={value} onValueChange={onChange}>
      <SelectTrigger className={cn('h-8 w-[150px]', className)}>
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        {allowClear && (
          <SelectItem value="">{placeholder}</SelectItem>
        )}
        {options.map((option) => (
          <SelectItem key={option.value} value={String(option.value)}>
            {option.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
})

/** 日期范围筛选 */
export interface DateRangeFilterProps {
  startDate?: string
  endDate?: string
  onStartDateChange?: (date: string) => void
  onEndDateChange?: (date: string) => void
  onChange?: (range: { startDate: string; endDate: string }) => void
  className?: string
}

export const DateRangeFilter = memo(function DateRangeFilter({
  startDate,
  endDate,
  onStartDateChange,
  onEndDateChange,
  onChange,
  className,
}: DateRangeFilterProps) {
  const handleStartChange = (value: string) => {
    onStartDateChange?.(value)
    if (onChange && endDate !== undefined) {
      onChange({ startDate: value, endDate })
    }
  }

  const handleEndChange = (value: string) => {
    onEndDateChange?.(value)
    if (onChange && startDate !== undefined) {
      onChange({ startDate, endDate: value })
    }
  }

  return (
    <div className={cn('flex items-center gap-2', className)}>
      <Input
        type="date"
        value={startDate}
        onChange={(e) => handleStartChange(e.target.value)}
        className="h-8 w-[140px]"
        placeholder="开始日期"
      />
      <span className="text-muted-foreground">至</span>
      <Input
        type="date"
        value={endDate}
        onChange={(e) => handleEndChange(e.target.value)}
        className="h-8 w-[140px]"
        placeholder="结束日期"
      />
    </div>
  )
})

/** 数字范围筛选 */
export interface NumberRangeFilterProps {
  min?: number
  max?: number
  onMinChange?: (value: number | undefined) => void
  onMaxChange?: (value: number | undefined) => void
  placeholderMin?: string
  placeholderMax?: string
  className?: string
}

export const NumberRangeFilter = memo(function NumberRangeFilter({
  min,
  max,
  onMinChange,
  onMaxChange,
  placeholderMin = '最小值',
  placeholderMax = '最大值',
  className,
}: NumberRangeFilterProps) {
  return (
    <div className={cn('flex items-center gap-2', className)}>
      <Input
        type="number"
        value={min ?? ''}
        onChange={(e) => onMinChange?.(e.target.value ? Number(e.target.value) : undefined)}
        className="h-8 w-[100px]"
        placeholder={placeholderMin}
      />
      <span className="text-muted-foreground">-</span>
      <Input
        type="number"
        value={max ?? ''}
        onChange={(e) => onMaxChange?.(e.target.value ? Number(e.target.value) : undefined)}
        className="h-8 w-[100px]"
        placeholder={placeholderMax}
      />
    </div>
  )
})

/** 筛选工具栏 */
export interface TableToolbarProps {
  children: React.ReactNode
  className?: string
}

export const TableToolbar = memo(function TableToolbar({
  children,
  className,
}: TableToolbarProps) {
  return (
    <div className={cn('flex flex-wrap items-center gap-2', className)}>
      {children}
    </div>
  )
})

/** 重置按钮 */
export interface ResetButtonProps {
  onReset: () => void
  className?: string
  disabled?: boolean
}

export const ResetButton = memo(function ResetButton({
  onReset,
  className,
  disabled,
}: ResetButtonProps) {
  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={onReset}
      disabled={disabled}
      className={cn('h-8 px-2 text-muted-foreground', className)}
    >
      <X className="mr-1 h-4 w-4" />
      重置
    </Button>
  )
})

// 统一导出
export const TableFilters = {
  Text: TextFilter,
  Select: SelectFilter,
  DateRange: DateRangeFilter,
  NumberRange: NumberRangeFilter,
  Toolbar: TableToolbar,
  Reset: ResetButton,
}

export default TableFilters
