import { useMemo } from 'react'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'

/** 预设颜色 */
const PRESET_COLORS = [
  '#ef4444', // red
  '#f97316', // orange
  '#f59e0b', // amber
  '#84cc16', // lime
  '#22c55e', // green
  '#10b981', // emerald
  '#14b8a6', // teal
  '#06b6d4', // cyan
  '#0ea5e9', // sky
  '#3b82f6', // blue
  '#6366f1', // indigo
  '#8b5cf6', // violet
  '#a855f7', // purple
  '#d946ef', // fuchsia
  '#ec4899', // pink
  '#f43f5e', // rose
  '#6b7280', // gray
  '#171717', // neutral
]

export interface ColorPickerProps {
  /** 当前颜色值 */
  value?: string
  /** 颜色变化回调 */
  onChange: (value: string) => void
  /** 是否显示预设颜色 */
  showPresets?: boolean
  /** 自定义预设颜色 */
  presetColors?: string[]
  /** 是否显示输入框 */
  showInput?: boolean
  /** 自定义类名 */
  className?: string
}

/**
 * 颜色选择器组件
 * 
 * 支持预设颜色面板 + 自定义颜色输入
 * 
 * @example
 * ```tsx
 * <ColorPicker
 *   value={color}
 *   onChange={setColor}
 *   showPresets
 *   showInput
 * />
 * ```
 */
export function ColorPicker({
  value = '#3b82f6',
  onChange,
  showPresets = true,
  presetColors = PRESET_COLORS,
  showInput = true,
  className,
}: ColorPickerProps) {
  // 处理颜色输入
  const handleTextChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const inputValue = e.target.value
    if (/^#[0-9A-Fa-f]{0,6}$/.test(inputValue) || inputValue === '') {
      onChange(inputValue || '#3b82f6')
    }
  }

  // 处理颜色选择器
  const handleColorChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange(e.target.value)
  }

  return (
    <div className={cn('space-y-3', className)}>
      {/* 预设颜色面板 */}
      {showPresets && (
        <div className="flex flex-wrap gap-2">
          {presetColors.map((color) => (
            <button
              key={color}
              type="button"
              onClick={() => onChange(color)}
              className={cn(
                'h-6 w-6 rounded-md transition-all',
                'hover:scale-110 hover:shadow-md',
                'focus:outline-none focus:ring-2 focus:ring-offset-1',
                value === color && 'ring-2 ring-offset-1 ring-primary scale-110'
              )}
              style={{ backgroundColor: color }}
              title={color}
            />
          ))}
        </div>
      )}

      {/* 颜色输入区域 */}
      <div className="flex items-center gap-3">
        <input
          type="color"
          value={value}
          onChange={handleColorChange}
          className="h-9 w-16 cursor-pointer rounded border border-input bg-background p-1"
        />
        {showInput && (
          <Input
            type="text"
            placeholder="#3b82f6"
            value={value}
            onChange={handleTextChange}
            className="flex-1 font-mono text-sm uppercase"
            maxLength={7}
          />
        )}
      </div>
    </div>
  )
}

/** 颜色预览组件 */
export interface ColorPreviewProps {
  /** 颜色值 */
  color?: string
  /** 标签文字 */
  label?: string
  /** 尺寸 */
  size?: 'sm' | 'md' | 'lg'
  /** 自定义类名 */
  className?: string
}

const sizeMap = {
  sm: 'h-5 w-5 text-xs',
  md: 'h-8 w-8 text-sm',
  lg: 'h-12 w-12 text-base',
}

/**
 * 颜色预览组件
 * 
 * 显示带颜色的标签效果
 */
export function ColorPreview({
  color = '#3b82f6',
  label = '预览',
  size = 'md',
  className,
}: ColorPreviewProps) {
  // 计算对比色（黑/白）
  const contrastColor = useMemo(() => {
    const hex = color.replace('#', '')
    const r = parseInt(hex.substr(0, 2), 16)
    const g = parseInt(hex.substr(2, 2), 16)
    const b = parseInt(hex.substr(4, 2), 16)
    const brightness = (r * 299 + g * 587 + b * 114) / 1000
    return brightness > 128 ? '#000000' : '#ffffff'
  }, [color])

  return (
    <div
      className={cn(
        'inline-flex items-center justify-center rounded-md font-medium transition-colors',
        sizeMap[size],
        className
      )}
      style={{
        backgroundColor: color,
        color: contrastColor,
      }}
    >
      {label}
    </div>
  )
}

export default ColorPicker
