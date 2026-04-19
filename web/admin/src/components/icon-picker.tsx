import * as React from 'react'
import * as Icons from 'lucide-react'
import type { LucideIcon } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Input } from '@/components/ui/input'

// 动态获取所有 lucide-react 图标名称
// lucide-react 导出两种格式：IconName 和 IconNameIcon（别名）
// 我们只保留不带 Icon 后缀的版本，避免重复
const getAllIconNames = (): string[] => {
  const iconNames: string[] = []
  const excludeNames = new Set([
    'createLucideIcon',
    'IconNode',
    'Icon',
    'IconProps',
    'LucideProps',
    'LucideIcon',
    'default',
  ])
  
  for (const name in Icons) {
    const icon = Icons[name as keyof typeof Icons]
    
    // 检查是否是有效的图标组件
    // lucide-react 的图标是 React 组件（forwardRef），具有 $$typeof 属性
    const isValidIcon = 
      icon !== undefined &&
      icon !== null &&
      typeof icon === 'object' &&
      ('$$typeof' in icon || typeof icon === 'function')
    
    // 排除：1. 排除列表中的名称 2. 以 Icon 结尾的（别名） 3. 不以大写字母开头的 4. 以下划线开头的
    if (
      !excludeNames.has(name) &&
      !name.endsWith('Icon') &&
      name[0] === name[0].toUpperCase() &&
      !name.startsWith('_') &&
      isValidIcon
    ) {
      iconNames.push(name)
    }
  }
  
  return iconNames.sort()
}

type IconPickerProps = {
  value?: string
  onValueChange?: (value: string) => void
  placeholder?: string
  className?: string
}

// 每页显示的图标数量
const ICONS_PER_PAGE = 100

function IconPicker({
  value,
  onValueChange,
  placeholder = '选择图标',
  className,
}: IconPickerProps) {
  const [open, setOpen] = React.useState(false)
  const [search, setSearch] = React.useState('')
  const [currentPage, setCurrentPage] = React.useState(1)

  // 在组件内部动态获取所有图标名称（使用 useMemo 确保只计算一次）
  const iconNames = React.useMemo(() => getAllIconNames(), [])

  // 获取当前选中的图标组件
  const SelectedIcon = value && value in Icons 
    ? (Icons as unknown as Record<string, LucideIcon>)[value] 
    : null

  // 过滤图标
  const filteredIcons = React.useMemo(() => {
    if (!search) return iconNames

    const searchLower = search.toLowerCase()
    return iconNames.filter((name) => {
      // lucide-react 图标名称直接使用，不需要移除后缀
      return name.toLowerCase().includes(searchLower)
    })
  }, [search, iconNames])

  // 当搜索内容改变时，重置到第一页
  React.useEffect(() => {
    setCurrentPage(1)
  }, [search])

  // 计算分页数据
  const totalPages = Math.ceil(filteredIcons.length / ICONS_PER_PAGE)
  const startIndex = (currentPage - 1) * ICONS_PER_PAGE
  const endIndex = startIndex + ICONS_PER_PAGE
  const currentPageIcons = filteredIcons.slice(startIndex, endIndex)

  const handleSelect = (iconName: string) => {
    onValueChange?.(iconName)
    setOpen(false)
    setSearch('')
    setCurrentPage(1)
  }

  const handlePreviousPage = () => {
    setCurrentPage((prev) => Math.max(1, prev - 1))
  }

  const handleNextPage = () => {
    setCurrentPage((prev) => Math.min(totalPages, prev + 1))
  }

  return (
      <Popover 
      open={open} 
      onOpenChange={(newOpen) => {
        setOpen(newOpen)
        if (!newOpen) {
          // 关闭时重置搜索和分页
          setSearch('')
          setCurrentPage(1)
        }
      }}
    >
      <PopoverTrigger asChild>
        <Button
          variant='outline'
          role='combobox'
          aria-expanded={open}
          className={cn('w-full justify-between', className)}
        >
          <div className='flex items-center gap-2'>
            {SelectedIcon ? (
              <>
                {React.createElement(SelectedIcon, { className: 'h-4 w-4' })}
                <span className='truncate'>{value}</span>
              </>
            ) : (
              <span className='text-muted-foreground'>{placeholder}</span>
            )}
          </div>
        </Button>
      </PopoverTrigger>
      <PopoverContent 
        className='w-[500px] p-0' 
        align='start'
        onWheel={(e) => {
          // 允许 PopoverContent 内部的滚动事件正常传播
          // 防止 Radix UI Popover 阻止内部 ScrollArea 的滚动
          e.stopPropagation()
        }}
      >
        <div className='p-3 border-b'>
          <div className='relative'>
            {React.createElement(
              (Icons as unknown as Record<string, LucideIcon>).Search,
              { className: 'absolute left-2 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground' }
            )}
            <Input
              type='text'
              placeholder='搜索图标...'
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className='pl-8'
              autoFocus
            />
            {search && (
              <button
                type='button'
                onClick={() => setSearch('')}
                className='absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground'
                aria-label='清除搜索'
              >
                {React.createElement(
                  (Icons as unknown as Record<string, LucideIcon>).X,
                  { className: 'h-4 w-4' }
                )}
              </button>
            )}
          </div>
        </div>
        <ScrollArea className='h-[280px]'>
          {filteredIcons.length > 0 ? (
            <>
              <div className='grid grid-cols-10 gap-2 p-3'>
                {currentPageIcons.map((iconName) => {
                  const IconComponent = (Icons as unknown as Record<string, LucideIcon>)[iconName]
                  if (!IconComponent) return null

                  const isSelected = value === iconName

                  return (
                    <button
                      key={iconName}
                      type='button'
                      onClick={() => handleSelect(iconName)}
                      className={cn(
                        'flex h-10 w-10 items-center justify-center rounded-md border transition-colors cursor-pointer',
                        isSelected
                          ? 'bg-accent border-primary text-primary'
                          : 'border-border hover:bg-accent hover:border-primary/50',
                      )}
                      title={iconName}
                    >
                      {React.createElement(IconComponent, { className: 'h-4 w-4' })}
                    </button>
                  )
                })}
              </div>
              {totalPages > 1 && (
                <div className='flex items-center justify-between border-t px-4 py-3'>
                  <div className='text-sm text-muted-foreground'>
                    第 {currentPage} / {totalPages} 页，共 {filteredIcons.length} 个图标
                  </div>
                  <div className='flex items-center gap-2'>
                    <Button
                      variant='outline'
                      size='sm'
                      onClick={handlePreviousPage}
                      disabled={currentPage === 1}
                      className='h-8 w-8 p-0'
                    >
                      {React.createElement(
                        (Icons as unknown as Record<string, LucideIcon>).ChevronLeft,
                        { className: 'h-4 w-4' }
                      )}
                    </Button>
                    <Button
                      variant='outline'
                      size='sm'
                      onClick={handleNextPage}
                      disabled={currentPage === totalPages}
                      className='h-8 w-8 p-0'
                    >
                      {React.createElement(
                        (Icons as unknown as Record<string, LucideIcon>).ChevronRight,
                        { className: 'h-4 w-4' }
                      )}
                    </Button>
                  </div>
                </div>
              )}
            </>
          ) : (
            <div className='flex items-center justify-center h-[280px] text-muted-foreground'>
              <p>未找到匹配的图标</p>
            </div>
          )}
        </ScrollArea>
      </PopoverContent>
    </Popover>
  )
}
export  {  IconPicker, type IconPickerProps }