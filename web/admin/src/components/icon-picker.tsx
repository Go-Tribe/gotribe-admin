import * as React from 'react'
import { ChevronLeft, ChevronRight, Search, X } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Input } from '@/components/ui/input'
import { getMenuIcon, MENU_ICON_OPTIONS } from '@/components/layout/icon-registry'

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

  // 获取当前选中的图标组件
  const SelectedIcon = getMenuIcon(value)

  // 过滤图标
  const filteredIcons = React.useMemo(() => {
    if (!search) return MENU_ICON_OPTIONS

    const searchLower = search.toLowerCase()
    return MENU_ICON_OPTIONS.filter(({ name }) => {
      // lucide-react 图标名称直接使用，不需要移除后缀
      return name.toLowerCase().includes(searchLower)
    })
  }, [search])

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
            <Search className='absolute left-2 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground' />
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
                <X className='h-4 w-4' />
              </button>
            )}
          </div>
        </div>
        <ScrollArea className='h-[280px]'>
          {filteredIcons.length > 0 ? (
            <>
              <div className='grid grid-cols-10 gap-2 p-3'>
                {currentPageIcons.map(({ name, icon: IconComponent }) => {
                  const isSelected = value === name

                  return (
                    <button
                      key={name}
                      type='button'
                      onClick={() => handleSelect(name)}
                      className={cn(
                        'flex h-10 w-10 items-center justify-center rounded-md border transition-colors cursor-pointer',
                        isSelected
                          ? 'bg-accent border-primary text-primary'
                          : 'border-border hover:bg-accent hover:border-primary/50',
                      )}
                      title={name}
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
                      <ChevronLeft className='h-4 w-4' />
                    </Button>
                    <Button
                      variant='outline'
                      size='sm'
                      onClick={handleNextPage}
                      disabled={currentPage === totalPages}
                      className='h-8 w-8 p-0'
                    >
                      <ChevronRight className='h-4 w-4' />
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

export { IconPicker, type IconPickerProps }
