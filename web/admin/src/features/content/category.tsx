import { useMemo, useState, useCallback } from 'react'
import { useQuery, useQueryClient, useMutation } from '@tanstack/react-query'
import {
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  getExpandedRowModel,
  type ColumnDef,
  type ExpandedState,
} from '@tanstack/react-table'
import { PlusIcon, Pencil1Icon, TrashIcon } from '@radix-ui/react-icons'
import { ChevronRight, ChevronDown } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { Input } from '@/components/ui/input'
import { DataTableFacetedFilter, DataTableColumnHeader, TreeDataTable } from '@/components/data-table'
import { SubTitle } from '@/components/sub-title'
import { Badge } from '@/components/ui/badge'
import { getStatusVariantByType } from '@/config/status-variants'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { CategoryFormDialog } from './component/category-form-dialog'
import type { Category, CategoryParams } from './types/category'
import { getCategoryTree, createCategory, updateCategory, batchDeleteCategory } from './service/category'
import { toast } from 'sonner'
import { useI18n } from '@/context/i18n-provider'
import { useDataTable } from '@/hooks/use-data-table'

export function ContentCategory() {
  'use no memo'
  const { t } = useI18n()
  const [dialogOpen, setDialogOpen] = useState<'create' | 'edit' | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<number | null>(null)
  const [editingCategory, setEditingCategory] = useState<Category | null>(null)

  // 使用统一的表格状态管理
  const {
    columnFilters,
    setColumnFilters,
    sorting,
    setSorting,
    columnVisibility,
    setColumnVisibility,
  } = useDataTable()

  // 树形表格的展开状态（保持独立，因为树形表格特殊）
  const [expanded, setExpanded] = useState<ExpandedState>({})

  // 获取分类树数据
  const { data, isLoading, error } = useQuery({
    queryKey: ['categoryList'],
    queryFn: () => getCategoryTree(),
  })
  const categoryData = data?.categoryTree || []

  // 扁平化分类树，用于查找删除项等
  const flattenCategoryData = useCallback((nodes: Category[]): Category[] => {
    const result: Category[] = []
    nodes.forEach((node) => {
      result.push(node)
      if (node.children?.length) {
        result.push(...flattenCategoryData(node.children))
      }
    })
    return result
  }, [])
  const allCategories = useMemo(() => flattenCategoryData(categoryData), [categoryData, flattenCategoryData])

  // 在树中按 id 查找节点（含 children），用于删除时收集自身及所有子分类 ID
  const findCategoryInTree = useCallback((nodes: Category[], categoryID: number): Category | null => {
    for (const node of nodes) {
      if (node.id === categoryID) return node
      if (node.children?.length) {
        const found = findCategoryInTree(node.children, categoryID)
        if (found) return found
      }
    }
    return null
  }, [])

  // 收集分类及其所有子分类的 id
  const collectCategoryIds = useCallback((category: Category): number[] => {
    const ids = [category.id]
    if (category.children?.length) {
      category.children.forEach((child) => ids.push(...collectCategoryIds(child)))
    }
    return ids
  }, [])

  const queryClient = useQueryClient()

  // 使用统一的 CRUD mutations（category 接口特殊，需要自定义）
  const createMutation = useMutation({
    mutationFn: (data: CategoryParams) => createCategory(data),
    onSuccess: () => {
      toast.success(t('features.content.category.createSuccess'))
      queryClient.invalidateQueries({ queryKey: ['categoryList'] })
      setDialogOpen(null)
      setEditingCategory(null)
    },
    onError: () => {},
  })

  const updateMutation = useMutation({
    mutationFn: (data: CategoryParams & { id: number }) => updateCategory(data.id, data),
    onSuccess: () => {
      toast.success(t('features.content.category.updateSuccess'))
      queryClient.invalidateQueries({ queryKey: ['categoryList'] })
      setDialogOpen(null)
      setEditingCategory(null)
    },
    onError: () => {},
  })

  const { mutate: deleteCategoryMutate, isPending: isDeletePending } = useMutation({
    mutationFn: (data: { ids: number[] }) => batchDeleteCategory(data),
    onSuccess: () => {
      toast.success(t('features.content.category.deleteSuccess'))
      queryClient.invalidateQueries({ queryKey: ['categoryList'] })
      setDeleteDialogOpen(null)
    },
    onError: () => {},
  })

  // 处理新建
  const handleCreate = useCallback(() => {
    setEditingCategory(null)
    setDialogOpen('create')
  }, [])

  // 处理编辑
  const handleEdit = useCallback((category: Category) => {
    setEditingCategory(category)
    setDialogOpen('edit')
  }, [])

  // 处理删除：有子分类时传父+所有子分类的 id
  const handleDelete = useCallback(() => {
    if (deleteDialogOpen == null) return
    const nodeInTree = findCategoryInTree(categoryData, deleteDialogOpen)
    const idsToDelete = nodeInTree
      ? collectCategoryIds(nodeInTree)
      : [deleteDialogOpen]
    deleteCategoryMutate({ ids: idsToDelete })
  }, [deleteDialogOpen, categoryData, findCategoryInTree, collectCategoryIds, deleteCategoryMutate])

  // 列定义
  const columns = useMemo<ColumnDef<Category>[]>(
    () => [
      {
        accessorKey: 'id',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.category.columns.id')} />
        ),
        cell: ({ row }) => (
          <div className='font-mono text-muted-foreground'>{String(row.getValue('id') ?? '')}</div>
        ),
      },
      {
        accessorKey: 'title',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.category.columns.title')} />
        ),
        cell: ({ row }) => {
          const canExpand = row.getCanExpand()
          const isExpanded = row.getIsExpanded()
          const level = row.depth
          return (
            <div 
              className='flex items-center gap-2 font-medium cursor-pointer' 
              style={{ paddingLeft: `${level * 20}px` }}
              onClick={() => canExpand && row.toggleExpanded()}
            >
              {canExpand && (
                <button
                  type='button'
                  onClick={(e) => {
                    e.stopPropagation()
                    row.toggleExpanded()
                  }}
                  className='flex items-center justify-center w-4 h-4 hover:bg-accent rounded'
                >
                  {isExpanded ? (
                    <ChevronDown className='h-4 w-4' />
                  ) : (
                    <ChevronRight className='h-4 w-4' />
                  )}
                </button>
              )}
              {!canExpand && <div className='w-4' />}
              <span className={canExpand ? 'hover:text-primary' : ''}>{row.getValue('title')}</span>
            </div>
          )
        },
      },
      {
        accessorKey: 'description',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.category.columns.description')} />
        ),
        cell: ({ row }) => {
          const desc = row.getValue('description') as string
          return (
            <div className='max-w-[200px] truncate' title={desc}>
              {desc || '-'}
            </div>
          )
        },
      },
      {
        accessorKey: 'icon',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.category.columns.icon')} />
        ),
        cell: ({ row }) => {
          const icon = row.getValue('icon') as string
          if (!icon) return <span className='text-muted-foreground'>-</span>
          const isUrl = /^https?:\/\//i.test(icon) || icon.startsWith('/')
          if (isUrl) {
            return (
              <img
                src={icon}
                alt=''
                className='h-8 w-8 object-cover rounded border border-border'
                title={icon}
              />
            )
          }
          return <div className='max-w-[100px] truncate' title={icon}>{icon}</div>
        },
      },
      {
        id: 'path',
        accessorFn: (row) => row.path ?? row.route,
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.category.columns.path')} />
        ),
        cell: ({ row }) => {
          const link = row.original.path ?? row.original.route
          if (!link) return <span className='text-muted-foreground'>-</span>
          const href = /^https?:\/\//i.test(link) ? link : `${window.location.origin}${link.startsWith('/') ? '' : '/'}${link}`
          return (
            <a
              href={href}
              target='_blank'
              rel='noopener noreferrer'
              className='text-primary hover:underline max-w-[150px] truncate inline-block'
              title={href}
            >
              {link}
            </a>
          )
        },
      },
      {
        accessorKey: 'sort',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.category.columns.sort')} />
        ),
        cell: ({ row }) => <div>{row.getValue('sort')}</div>,
      },
      {
        accessorKey: 'hidden',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.category.columns.hidden')} />
        ),
        cell: ({ row }) => {
          const hidden = row.getValue('hidden') as number
          return (
            <Badge variant={getStatusVariantByType(hidden, 'visibility')}>
              {hidden === 1 ? t('features.content.category.hiddenStatus.show') : t('features.content.category.hiddenStatus.hidden')}
            </Badge>
          )
        },
      },
      {
        id: 'actions',
        header: t('features.content.category.columns.actions'),
        cell: ({ row }) => {
          const category = row.original
          return (
            <div className='flex items-center gap-1'>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant='outline'
                    size='icon'
                    className='h-8 w-8 border-border/60'
                    onClick={() => handleEdit(category)}
                  >
                    <Pencil1Icon className='h-4 w-4' />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>{t('features.content.category.actions.edit')}</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant='ghost'
                    size='icon'
                    className='h-8 w-8 text-destructive hover:text-destructive'
                    onClick={() => setDeleteDialogOpen(category.id)}
                  >
                    <TrashIcon className='h-4 w-4 text-destructive' />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>{t('features.content.category.actions.delete')}</TooltipContent>
              </Tooltip>
            </div>
          )
        },
        enableHiding: false,
        meta: {
          className: 'sticky right-0 bg-background z-10 shadow-[inset_-1px_0_0_0_hsl(var(--border))]',
          thClassName: 'sticky right-0 bg-background z-10 shadow-[inset_-1px_0_0_0_hsl(var(--border))]',
          tdClassName: 'sticky right-0 bg-background z-10 shadow-[inset_-1px_0_0_0_hsl(var(--border))]',
        },
      },
    ],
    [t, handleEdit]
  )

  // 表格实例（树形表格）
  // eslint-disable-next-line react-hooks/incompatible-library
  const table = useReactTable({
    data: categoryData,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getExpandedRowModel: getExpandedRowModel(),
    getSubRows: (row) => row.children,
    onExpandedChange: setExpanded,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onColumnVisibilityChange: setColumnVisibility,
    state: {
      expanded,
      sorting,
      columnFilters,
      columnVisibility,
    },
  })

  // 处理表单提交
  const handleFormSubmit = (data: CategoryParams) => {
    if (dialogOpen === 'create') {
      createMutation.mutate(data)
    } else if (dialogOpen === 'edit' && data.id) {
      updateMutation.mutate({ id: data.id, ...data })
    }
  }

  const categoryToDelete = allCategories.find((c: Category) => c.id === deleteDialogOpen)
  const isDialogLoading = dialogOpen === 'create'
    ? createMutation.isPending
    : updateMutation.isPending

  return (
    <div className='space-y-4'>
      <div className='flex items-center justify-between px-4 pt-4'>
        <SubTitle
          title={t('features.content.category.title')}
          description={t('features.content.category.description')}
          children={
            <Button onClick={handleCreate}>
              <PlusIcon className='h-4 w-4' />
              {t('features.content.category.createButton')}
            </Button>
          }
        />
      </div>

      <div className='rounded-md border p-6 mx-4'>
        <TreeDataTable
          table={table}
          columns={columns}
          isLoading={isLoading}
          error={error}
          loadingText={t('features.content.category.loading')}
          errorText={t('features.content.category.loadError')}
          emptyText={t('features.content.category.noData')}
          bordered={false}
        >
          {/* 自定义搜索输入框 */}
          <div className='flex flex-wrap gap-2'>
            <Input
              type='text'
              placeholder={t('features.content.category.search.title')}
              value={(table.getColumn('title')?.getFilterValue() as string) ?? ''}
              onChange={(e) =>
                table.getColumn('title')?.setFilterValue(e.target.value)
              }
              className='h-8 w-[150px]'
            />
            <DataTableFacetedFilter
              column={table.getColumn('hidden')}
              title={t('features.content.category.filter.hidden')}
              options={[
                { label: t('features.content.category.hiddenStatus.show'), value: '1' },
                { label: t('features.content.category.hiddenStatus.hidden'), value: '2' },
              ]}
              single
            />
          </div>
        </TreeDataTable>
      </div>

      {/* 新建/编辑对话框 */}
      <CategoryFormDialog
        open={dialogOpen !== null}
        onOpenChange={(open) => {
          if (!open) {
            setDialogOpen(null)
            setEditingCategory(null)
          }
        }}
        category={editingCategory}
        categoryTree={categoryData}
        onSubmit={handleFormSubmit}
        isLoading={isDialogLoading}
      />

      {/* 删除确认对话框 */}
      <ConfirmDialog
        open={deleteDialogOpen !== null}
        onOpenChange={(open) => {
          if (!open) setDeleteDialogOpen(null)
        }}
        title={t('features.content.category.confirmDelete')}
        desc={t('features.content.category.confirmDeleteMessage', { title: categoryToDelete?.title ?? '' })}
        handleConfirm={handleDelete}
        destructive
        confirmText={t('features.content.category.actions.delete')}
        isLoading={isDeletePending}
      />
    </div>
  )
}
