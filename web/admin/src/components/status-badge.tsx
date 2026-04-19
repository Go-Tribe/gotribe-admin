import { memo, useMemo } from 'react'
import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'
import {
  getStatusVariantByType,
  getHttpStatusVariant,
  type StatusType,
  type BadgeVariant,
} from '@/config/status-variants'

/** 状态选项配置 */
export interface StatusOption {
  label: string
  value: number | string
  variant?: BadgeVariant
  color?: string
}

/** 状态映射配置 */
export type StatusMap = Record<number | string, { label: string; variant: BadgeVariant }>

export interface StatusBadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  /** 状态值 */
  value: number | string
  /** 状态类型（使用预定义配置） */
  type?: StatusType
  /** 状态选项数组 */
  options?: StatusOption[]
  /** 状态映射对象 */
  map?: StatusMap
  /** HTTP 状态码模式 */
  httpStatus?: boolean
  /** 自定义标签 */
  label?: string
  /** 自定义颜色 */
  variant?: BadgeVariant
  /** 是否使用点状样式 */
  dot?: boolean
  /** 自定义类名 */
  className?: string
}

/**
 * 状态徽章组件
 * 
 * 统一状态展示样式，支持多种配置方式
 * 
 * @example
 * ```tsx
 * // 使用预定义类型
 * <StatusBadge value={1} type="enabledStatus" />
 * 
 * // 使用选项数组
 * <StatusBadge
 *   value={status}
 *   options={[
 *     { label: '启用', value: 1, variant: 'success' },
 *     { label: '禁用', value: 2, variant: 'secondary' },
 *   ]}
 * />
 * 
 * // 使用映射对象
 * <StatusBadge
 *   value={status}
 *   map={{
 *     1: { label: '正常', variant: 'success' },
 *     2: { label: '异常', variant: 'destructive' },
 *   }}
 * />
 * 
 * // HTTP 状态码
 * <StatusBadge value={200} httpStatus />
 * 
 * // 点状样式
 * <StatusBadge value={1} type="enabledStatus" dot />
 * ```
 */
export const StatusBadge = memo(function StatusBadge({
  value,
  type,
  options,
  map,
  httpStatus = false,
  label: customLabel,
  variant: customVariant,
  dot = false,
  className,
  ...props
}: StatusBadgeProps) {
  // 计算状态配置
  const config = useMemo(() => {
    // 优先级 1: 自定义标签和样式
    if (customLabel && customVariant) {
      return { label: customLabel, variant: customVariant }
    }

    // 优先级 2: HTTP 状态码
    if (httpStatus && typeof value === 'number') {
      return {
        label: String(value),
        variant: getHttpStatusVariant(value),
      }
    }

    // 优先级 3: 预定义类型
    if (type && typeof value === 'number') {
      return {
        label: getLabelByType(value, type),
        variant: getStatusVariantByType(value, type),
      }
    }

    // 优先级 4: 选项数组
    if (options) {
      const option = options.find((opt) => opt.value === value)
      if (option) {
        return {
          label: option.label,
          variant: option.variant || 'secondary',
        }
      }
    }

    // 优先级 5: 映射对象
    if (map && map[value] !== undefined) {
      return map[value]
    }

    // 默认回退
    return {
      label: String(value),
      variant: 'secondary' as BadgeVariant,
    }
  }, [value, type, options, map, httpStatus, customLabel, customVariant])

  const { label, variant } = config

  return (
    <Badge
      variant={variant}
      className={cn(
        'font-normal',
        dot && 'gap-1.5 pl-2',
        className
      )}
      {...props}
    >
      {dot && (
        <span
          className={cn(
            'h-1.5 w-1.5 rounded-full',
            variant === 'success' && 'bg-green-500',
            variant === 'destructive' && 'bg-red-500',
            variant === 'warning' && 'bg-yellow-500',
            variant === 'secondary' && 'bg-gray-400',
            variant === 'default' && 'bg-primary',
            variant === 'outline' && 'bg-foreground',
          )}
        />
      )}
      {label}
    </Badge>
  )
})

/** 根据类型获取标签 */
function getLabelByType(value: number, type: StatusType): string {
  const labels: Record<StatusType, Record<number, string>> = {
    enabledStatus: { 0: '禁用', 1: '启用', 2: '禁用' },
    publishStatus: { 1: '未发布', 2: '已发布' },
    approvalStatus: { 1: '待审核', 2: '已通过' },
    articlePublish: { 1: '草稿', 2: '已发布' },
    visibility: { 1: '显示', 2: '隐藏' },
    visibilityInverse: { 1: '隐藏', 2: '显示' },
  }

  return labels[type]?.[value] || String(value)
}

/** 批量状态展示 */
export interface StatusBadgeGroupProps {
  statuses: Array<{ value: number | string; label: string }>
  type?: StatusType
  className?: string
}

export const StatusBadgeGroup = memo(function StatusBadgeGroup({
  statuses,
  type,
  className,
}: StatusBadgeGroupProps) {
  return (
    <div className={cn('flex flex-wrap gap-1', className)}>
      {statuses.map((status, index) => (
        <StatusBadge
          key={index}
          value={status.value}
          type={type}
          label={status.label}
          variant="outline"
        />
      ))}
    </div>
  )
})

export default StatusBadge
