import * as React from 'react'
import * as PopoverPrimitive from '@radix-ui/react-popover'
import {
  CheckIcon,
  ChevronDownIcon,
  ChevronRightIcon,
  XIcon,
} from 'lucide-react'
import { cn } from '@/lib/utils'

// 树节点数据结构
export interface TreeNode {
  value: string
  label: string
  disabled?: boolean
  children?: TreeNode[]
}

// TreeSelect 组件 Props
export interface TreeSelectProps {
  /** 树形数据 */
  data: TreeNode[]
  /** 当前选中的值 */
  value?: string | string[]
  /** 值变化时的回调 */
  onChange?: (value: string | string[]) => void
  /** 是否多选 */
  multiple?: boolean
  /** 占位符文本 */
  placeholder?: string
  /** 是否禁用 */
  disabled?: boolean
  /** 自定义类名 */
  className?: string
  /** 触发器尺寸 */
  size?: 'sm' | 'default'
  /** 是否默认展开所有节点 */
  defaultExpandAll?: boolean
  /** 是否可清除 */
  clearable?: boolean
  /** 是否可搜索 */
  searchable?: boolean
  /** 搜索占位符 */
  searchPlaceholder?: string
  /** 空数据提示 */
  emptyText?: string
  /** 父节点是否可选 */
  parentSelectable?: boolean
}

// 获取所有节点的值
function getAllNodeValues(nodes: TreeNode[]): string[] {
  const values: string[] = []
  const traverse = (items: TreeNode[]) => {
    items.forEach((item) => {
      values.push(item.value)
      if (item.children?.length) {
        traverse(item.children)
      }
    })
  }
  traverse(nodes)
  return values
}

// 根据值查找节点
function findNodeByValue(
  nodes: TreeNode[],
  value: string
): TreeNode | undefined {
  for (const node of nodes) {
    if (node.value === value) return node
    if (node.children?.length) {
      const found = findNodeByValue(node.children, value)
      if (found) return found
    }
  }
  return undefined
}

/**
 * 过滤树节点并同时收集需要展开的节点值
 * 合并 filterTree 和 getExpandedValuesFromFilter 逻辑，避免重复遍历
 */
function filterTreeWithExpanded(
  nodes: TreeNode[],
  keyword: string
): { filtered: TreeNode[]; expandedValues: Set<string> } {
  const expanded = new Set<string>()
  const lowerKeyword = keyword.toLowerCase()

  const traverse = (items: TreeNode[], parentPath: string[]): TreeNode[] => {
    const levelResult: TreeNode[] = []
    for (const node of items) {
      const matchLabel = node.label.toLowerCase().includes(lowerKeyword)
      const currentPath = [...parentPath, node.value]
      const filteredChildren = node.children?.length
        ? traverse(node.children, currentPath)
        : []

      if (matchLabel || filteredChildren.length > 0) {
        levelResult.push({
          ...node,
          children: filteredChildren.length > 0 ? filteredChildren : node.children,
        })
      }

      // 如果子节点中有匹配，当前节点需要展开
      if (filteredChildren.length > 0) {
        currentPath.forEach((v) => expanded.add(v))
      }
    }
    return levelResult
  }

  const filtered = traverse(nodes, [])
  return { filtered, expandedValues: expanded }
}

// 递归渲染树
function TreeRenderer({
  nodes,
  selectedValues,
  expandedValues,
  onSelect,
  onToggle,
  multiple,
  parentSelectable,
}: {
  nodes: TreeNode[]
  selectedValues: string[]
  expandedValues: Set<string>
  onSelect: (value: string) => void
  onToggle: (value: string) => void
  multiple?: boolean
  parentSelectable?: boolean
}) {
  const renderNode = (node: TreeNode, level: number): React.ReactNode => {
    const hasChildren = node.children && node.children.length > 0
    const isSelected = selectedValues.includes(node.value)
    const isExpanded = expandedValues.has(node.value)
    const canSelect = !node.disabled && (parentSelectable || !hasChildren)

    return (
      <div key={node.value} data-slot='tree-node'>
        <div
          className={cn(
            'flex cursor-pointer items-center gap-1 rounded-sm px-2 py-1.5 text-sm outline-hidden select-none hover:bg-accent hover:text-accent-foreground',
            isSelected && 'bg-accent text-accent-foreground',
            node.disabled && 'pointer-events-none opacity-50'
          )}
          style={{ paddingLeft: `${level * 16 + 8}px` }}
          onClick={() => {
            if (canSelect) {
              onSelect(node.value)
            } else if (hasChildren) {
              onToggle(node.value)
            }
          }}
        >
          {/* 展开/折叠图标 */}
          {hasChildren ? (
            <button
              type='button'
              className='flex size-4 shrink-0 items-center justify-center rounded hover:bg-accent'
              onClick={(e) => {
                e.stopPropagation()
                onToggle(node.value)
              }}
            >
              <ChevronRightIcon
                className={cn(
                  'size-3.5 transition-transform duration-200',
                  isExpanded && 'rotate-90'
                )}
              />
            </button>
          ) : (
            <span className='size-4 shrink-0' />
          )}

          {/* 多选框 */}
          {multiple && canSelect && (
            <div
              className={cn(
                'flex size-4 shrink-0 items-center justify-center rounded-sm border border-primary',
                isSelected && 'bg-primary text-primary-foreground'
              )}
            >
              {isSelected && <CheckIcon className='size-3' />}
            </div>
          )}

          {/* 节点标签 */}
          <span className='flex-1 truncate'>{node.label}</span>

          {/* 单选选中图标 */}
          {!multiple && isSelected && (
            <CheckIcon className='size-4 shrink-0 text-primary' />
          )}
        </div>

        {/* 子节点 */}
        {hasChildren && isExpanded && (
          <div data-slot='tree-children'>
            {node.children!.map((child) => renderNode(child, level + 1))}
          </div>
        )}
      </div>
    )
  }

  return <>{nodes.map((node) => renderNode(node, 0))}</>
}

// 主组件
function TreeSelect({
  data,
  value,
  onChange,
  multiple = false,
  placeholder = '请选择',
  disabled = false,
  className,
  size = 'default',
  defaultExpandAll = false,
  clearable = false,
  searchable = false,
  searchPlaceholder = '搜索...',
  emptyText = '暂无数据',
  parentSelectable = false,
}: TreeSelectProps) {
  const [open, setOpen] = React.useState(false)
  const [searchKeyword, setSearchKeyword] = React.useState('')
  const [expandedValues, setExpandedValues] = React.useState<Set<string>>(
    () => {
      if (defaultExpandAll) {
        return new Set(getAllNodeValues(data))
      }
      return new Set()
    }
  )

  // 处理选中值
  const selectedValues = React.useMemo(() => {
    if (value === undefined || value === null) return []
    if (Array.isArray(value)) return value
    return [value]
  }, [value])

  // 过滤后的树数据 + 搜索时自动展开匹配节点的父节点
  // 合并计算，避免重复遍历
  const { filteredData, searchExpandedValues } = React.useMemo(() => {
    if (!searchKeyword) return { filteredData: data, searchExpandedValues: new Set<string>() }
    const { filtered, expandedValues } = filterTreeWithExpanded(data, searchKeyword)
    return { filteredData: filtered, searchExpandedValues: expandedValues }
  }, [data, searchKeyword])

  React.useEffect(() => {
    if (searchKeyword && searchExpandedValues.size > 0) {
      setExpandedValues((prev) => new Set([...prev, ...searchExpandedValues]))
    }
  }, [searchKeyword, searchExpandedValues])

  // 获取显示文本
  const displayText = React.useMemo(() => {
    if (selectedValues.length === 0) return ''
    if (multiple) {
      const labels = selectedValues
        .map((v) => findNodeByValue(data, v)?.label)
        .filter(Boolean)
      return labels.join(', ')
    }
    return findNodeByValue(data, selectedValues[0])?.label || ''
  }, [selectedValues, data, multiple])

  // 处理选择
  const handleSelect = (nodeValue: string) => {
    if (multiple) {
      const newValues = selectedValues.includes(nodeValue)
        ? selectedValues.filter((v) => v !== nodeValue)
        : [...selectedValues, nodeValue]
      onChange?.(newValues)
    } else {
      onChange?.(nodeValue)
      setOpen(false)
    }
  }

  // 处理展开/折叠
  const handleToggle = (nodeValue: string) => {
    setExpandedValues((prev) => {
      const next = new Set(prev)
      if (next.has(nodeValue)) {
        next.delete(nodeValue)
      } else {
        next.add(nodeValue)
      }
      return next
    })
  }

  // 处理清除
  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation()
    onChange?.(multiple ? [] : '')
  }

  // 关闭时清除搜索
  React.useEffect(() => {
    if (!open) {
      setSearchKeyword('')
    }
  }, [open])

  return (
    <PopoverPrimitive.Root open={open} onOpenChange={setOpen}>
      <PopoverPrimitive.Trigger asChild disabled={disabled}>
        <button
          type='button'
          data-slot='tree-select-trigger'
          data-size={size}
          className={cn(
            'flex w-full items-center justify-between gap-2 rounded-md border border-input bg-transparent px-3 py-2 text-sm whitespace-nowrap shadow-xs transition-[color,box-shadow] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 data-[size=default]:h-9 data-[size=sm]:h-8 dark:bg-input/30 dark:hover:bg-input/50',
            !displayText && 'text-muted-foreground',
            className
          )}
        >
          <span className='flex-1 truncate text-left'>
            {displayText || placeholder}
          </span>
          <span className='flex items-center gap-1'>
            {clearable && selectedValues.length > 0 && (
              <span
                className='flex size-4 items-center justify-center rounded-full hover:bg-accent'
                onClick={handleClear}
              >
                <XIcon className='size-3 text-muted-foreground' />
              </span>
            )}
            <ChevronDownIcon
              className={cn(
                'size-4 text-muted-foreground transition-transform duration-200',
                open && 'rotate-180'
              )}
            />
          </span>
        </button>
      </PopoverPrimitive.Trigger>

      <PopoverPrimitive.Portal>
        <PopoverPrimitive.Content
          data-slot='tree-select-content'
          align='start'
          sideOffset={4}
          className={cn(
            'z-50 min-w-(--radix-popover-trigger-width) max-h-[300px] overflow-hidden rounded-md border bg-popover text-popover-foreground shadow-md outline-hidden data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95'
          )}
        >
          {/* 搜索框 */}
          {searchable && (
            <div className='border-b p-2'>
              <input
                type='text'
                className='w-full rounded-md border border-input bg-transparent px-3 py-1.5 text-sm outline-none placeholder:text-muted-foreground focus:border-ring focus:ring-1 focus:ring-ring'
                placeholder={searchPlaceholder}
                value={searchKeyword}
                onChange={(e) => setSearchKeyword(e.target.value)}
              />
            </div>
          )}

          {/* 树内容 */}
          <div className='max-h-[250px] overflow-y-auto p-1'>
            {filteredData.length > 0 ? (
              <TreeRenderer
                nodes={filteredData}
                selectedValues={selectedValues}
                expandedValues={expandedValues}
                onSelect={handleSelect}
                onToggle={handleToggle}
                multiple={multiple}
                parentSelectable={parentSelectable}
              />
            ) : (
              <div className='py-6 text-center text-sm text-muted-foreground'>
                {emptyText}
              </div>
            )}
          </div>
        </PopoverPrimitive.Content>
      </PopoverPrimitive.Portal>
    </PopoverPrimitive.Root>
  )
}

export { TreeSelect }
