import { useMemo, useState, useCallback } from 'react'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import {
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  type ColumnDef,
} from '@tanstack/react-table'
import { PlusIcon, Pencil1Icon, TrashIcon } from '@radix-ui/react-icons'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  DataTableColumnHeader,
  DataTablePagination,
} from '@/components/data-table'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { useI18n } from '@/context/i18n-provider'
import { TagFormDialog } from './component/tag-form-dialog'
import type { Tag } from './types/tag'
import { getTagList, createTag, updateTag, deleteTag } from './service/tag'
import { useDataTable } from '@/hooks/use-data-table'
import { useCrudMutations } from '@/hooks/use-crud-mutations'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

export function ContentTag() {
  const { t } = useI18n()
  const [dialogOpen, setDialogOpen] = useState<'create' | 'edit' | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<number | null>(null)
  const [editingTag, setEditingTag] = useState<Tag | null>(null)

  // 使用统一的表格状态管理
  const {
    columnFilters,
    setColumnFilters,
    pagination,
    setPagination,
    sorting,
    setSorting,
    columnVisibility,
    setColumnVisibility,
    pageNum,
    getFilterValue,
    getSortingParams,
  } = useDataTable()

  // 构建查询参数
  const queryParams = useMemo(() => ({
    current: 1,
    title: getFilterValue('title'),
    description: getFilterValue('description'),
    pageNum,
    pageSize: pagination.pageSize,
    ...getSortingParams(),
  }), [getFilterValue, getSortingParams, pagination.pageSize, pageNum])

  // 获取标签列表数据
  const { data, isLoading, error } = useQuery({
    queryKey: ['tagList', queryParams],
    queryFn: () => getTagList(queryParams),
    placeholderData: keepPreviousData,
  })

  // 提取数据
  const tagData = data?.tags || []
  const total = data?.total || 0
  const pageCount = Math.ceil(total / pagination.pageSize)

  // 使用统一的 CRUD mutations
  const { createMutation, updateMutation, deleteMutation } = useCrudMutations<Tag, number>({
    queryKey: ['tagList'],
    createFn: createTag,
    updateFn: (data) => {
      if (!data.id) {
        throw new Error(t('features.content.tag.form.validation.titleRequired'))
      }
      return updateTag(data.id, data)
    },
    deleteFn: deleteTag,
    messages: {
      createSuccess: t('features.content.tag.createSuccess'),
      updateSuccess: t('features.content.tag.updateSuccess'),
      deleteSuccess: t('features.content.tag.deleteSuccess'),
    },
    onSuccess: () => {
      setDialogOpen(null)
      setEditingTag(null)
      setDeleteDialogOpen(null)
    },
  })

  // 处理新建
  const handleCreate = useCallback(() => {
    setEditingTag(null)
    setDialogOpen('create')
  }, [])

  // 处理编辑
  const handleEdit = useCallback((tag: Tag) => {
    setEditingTag(tag)
    setDialogOpen('edit')
  }, [])

  // 处理删除
  const { mutate: deleteTagMutate } = deleteMutation
  const handleDelete = useCallback(() => {
    if (deleteDialogOpen) {
      deleteTagMutate(deleteDialogOpen)
    }
  }, [deleteDialogOpen, deleteTagMutate])

  // 列定义
  const columns = useMemo<ColumnDef<Tag>[]>(
    () => [
      {
        accessorKey: 'id',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.tag.columns.id')} />
        ),
        cell: ({ row }) => {
          const id = row.getValue('id') as number
          return <div className='font-medium'>{id}</div>
        },
      },
      {
        accessorKey: 'title',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.tag.columns.title')} />
        ),
        cell: ({ row }) => {
          const tag = row.original
          return (
            <div className='flex items-center gap-2'>
              {tag.color && (
                <div
                  className='h-4 w-4 rounded-full border border-border'
                  style={{ backgroundColor: tag.color }}
                />
              )}
              <div className='font-medium'>{row.getValue('title')}</div>
            </div>
          )
        },
      },
      {
        accessorKey: 'description',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.tag.columns.description')} />
        ),
        cell: ({ row }) => {
          const description = row.getValue('description') as string
          return (
            <div className='max-w-[300px] truncate' title={description}>
              {description || '-'}
            </div>
          )
        },
      },
      {
        id: 'actions',
        header: t('features.content.tag.columns.actions'),
        cell: ({ row }) => {
          const tag = row.original
          return (
            <div className='flex items-center gap-1'>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant='ghost'
                    size='icon'
                    className='h-8 w-8'
                    onClick={() => handleEdit(tag)}
                  >
                    <Pencil1Icon className='h-4 w-4' />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>{t('features.content.tag.actions.edit')}</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant='ghost'
                    size='icon'
                    className='h-8 w-8'
                    onClick={() => setDeleteDialogOpen(tag.id)}
                  >
                    <TrashIcon className='h-4 w-4 text-destructive' />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>{t('features.content.tag.actions.delete')}</TooltipContent>
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

  // 表格实例
  const table = useReactTable({
    data: tagData,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    manualSorting: true,
    manualPagination: true,
    pageCount,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onColumnVisibilityChange: setColumnVisibility,
    onPaginationChange: setPagination,
    state: {
      sorting,
      columnFilters,
      columnVisibility,
      pagination,
    },
  })

  const tagToDelete = tagData.find((t: Tag) => t.id === deleteDialogOpen)
  const isDialogLoading = dialogOpen === 'create' 
    ? createMutation.isPending 
    : updateMutation.isPending

  return (
    <div className='space-y-4'>
      <div className='flex items-center justify-between px-4 pt-4'>
        <div className='flex w-full items-center justify-between'>
          <div className='flex flex-col gap-2'>
            <h2 className='text-2xl font-bold'>{t('features.content.tag.title')}</h2>
            <p className='text-sm text-muted-foreground'>{t('features.content.tag.description')}</p>
          </div>
          <div>
            <Button onClick={handleCreate}>
              <PlusIcon className='h-4 w-4' />
              {t('features.content.tag.createButton')}
            </Button>
          </div>
        </div>
      </div>

      <div className='rounded-md border p-6 mx-4'>
        <div className='flex flex-wrap gap-2 pb-4'>
          <Input
            type='text'
            placeholder={t('features.content.tag.search.title')}
            value={(table.getColumn('title')?.getFilterValue() as string) ?? ''}
            onChange={(e) =>
              table.getColumn('title')?.setFilterValue(e.target.value)
            }
            className='h-8 w-[150px]'
          />
          <Input
            type='text'
            placeholder={t('features.content.tag.search.description')}
            value={(table.getColumn('description')?.getFilterValue() as string) ?? ''}
            onChange={(e) =>
              table.getColumn('description')?.setFilterValue(e.target.value)
            }
            className='h-8 w-[150px]'
          />
        </div>

        <div className='overflow-x-auto'>
          <Table>
            <TableHeader>
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id}>
                  {headerGroup.headers.map((header) => {
                    const headerContent = header.isPlaceholder
                      ? null
                      : typeof header.column.columnDef.header === 'function'
                        ? header.column.columnDef.header(header.getContext())
                        : header.column.columnDef.header
                    const meta = header.column.columnDef.meta as { thClassName?: string } | undefined
                    return (
                      <TableHead key={header.id} className={meta?.thClassName}>
                        {headerContent}
                      </TableHead>
                    )
                  })}
                </TableRow>
              ))}
            </TableHeader>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={columns.length} className='h-24 text-center'>
                    {t('features.content.tag.loading')}
                  </TableCell>
                </TableRow>
              ) : error ? (
                <TableRow>
                  <TableCell colSpan={columns.length} className='h-24 text-center text-destructive'>
                    {t('features.content.tag.loadError')}
                  </TableCell>
                </TableRow>
              ) : table.getRowModel().rows?.length ? (
                table.getRowModel().rows.map((row) => (
                  <TableRow key={row.id} data-state={row.getIsSelected() && 'selected'}>
                    {row.getVisibleCells().map((cell) => {
                      const cellContent =
                        typeof cell.column.columnDef.cell === 'function'
                          ? cell.column.columnDef.cell(cell.getContext())
                          : cell.column.columnDef.cell
                      const meta = cell.column.columnDef.meta as { tdClassName?: string } | undefined
                      return (
                        <TableCell key={cell.id} className={meta?.tdClassName}>
                          {cellContent}
                        </TableCell>
                      )
                    })}
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={columns.length} className='h-24 text-center'>
                    {t('features.content.tag.noData')}
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>

        <DataTablePagination table={table} className='border-t border-border/60 px-4 py-4' />
      </div>

      {/* 新建/编辑对话框 */}
      <TagFormDialog
        open={dialogOpen !== null}
        onOpenChange={(open) => {
          if (!open) {
            setDialogOpen(null)
            setEditingTag(null)
          }
        }}
        tag={editingTag}
        onSubmit={(data) => {
          if (dialogOpen === 'create') {
            createMutation.mutate(data)
          } else {
            updateMutation.mutate(data)
          }
        }}
        isLoading={isDialogLoading}
      />

      {/* 删除确认对话框 */}
      <ConfirmDialog
        open={deleteDialogOpen !== null}
        onOpenChange={(open) => {
          if (!open) setDeleteDialogOpen(null)
        }}
        title={t('features.content.tag.confirmDelete')}
        desc={t('features.content.tag.confirmDeleteDesc', { title: tagToDelete?.title ?? '' })}
        handleConfirm={handleDelete}
        destructive
        confirmText={t('features.content.tag.actions.delete')}
        isLoading={deleteMutation.isPending}
      />
    </div>
  )
}
// MIGRATED
