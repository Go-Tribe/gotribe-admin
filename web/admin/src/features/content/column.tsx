import { useMemo, useState, useCallback } from 'react'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import {
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  type ColumnDef,
} from '@tanstack/react-table'
import { PlusIcon } from '@radix-ui/react-icons'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  DataTableColumnHeader,
  DataTable,
  DataTableActions,
} from '@/components/data-table'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { ListPageLayout } from '@/components/layout'
import { ColumnFormDialog } from './component/column-form-dialog'
import type { Column as ColumnType, ColumnCreateParams } from './types/column'
import { getColumnList, createColumn, updateColumn, deleteColumn } from './service/column'
import { getProjectList } from '@/features/business/service/project'
import { useI18n } from '@/context/i18n-provider'
import { useDataTable } from '@/hooks/use-data-table'
import { useCrudMutations } from '@/hooks/use-crud-mutations'

// 适配器函数：将 Partial<ColumnType> 转换为服务函数期望的类型
const adaptedCreateColumn = async (data: Partial<ColumnType>): Promise<unknown> => {
  return createColumn(data as ColumnCreateParams)
}

const adaptedUpdateColumn = async (data: Partial<ColumnType>): Promise<unknown> => {
  if (data.id == null) {
    throw new Error('Column ID is required')
  }
  const { id, ...updateData } = data
  return updateColumn(String(id), updateData)
}

export function ContentColumn() {
  'use no memo'
  const { t } = useI18n()
  const [dialogOpen, setDialogOpen] = useState<'create' | 'edit' | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<string | null>(null)
  const [editingColumn, setEditingColumn] = useState<ColumnType | null>(null)
  const [projectFilter, setProjectFilter] = useState<string>('__all__')

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
  } = useDataTable()

  // 构建查询参数
  const queryParams = useMemo(() => ({
    projectID: projectFilter && projectFilter !== '__all__' ? projectFilter.trim() || undefined : undefined,
    title: getFilterValue('title'),
    pageNum,
    pageSize: pagination.pageSize,
  }), [getFilterValue, pagination.pageSize, pageNum, projectFilter])

  // 获取项目列表数据
  const { data: projectData } = useQuery({
    queryKey: ['projectList', { current: 1, pageNum: 1, pageSize: 1000 }],
    queryFn: () => getProjectList({ current: 1, pageNum: 1, pageSize: 1000 }),
  })

  const projectList = useMemo(
    () =>
      projectData?.projects?.map((p) => ({
        projectID: p.projectID,
        title: p.title ?? p.projectID,
      })) ?? [],
    [projectData?.projects]
  )

  // 获取专栏列表数据
  const { data, isLoading, error } = useQuery({
    queryKey: ['columnList', queryParams],
    queryFn: () => getColumnList(queryParams),
    placeholderData: keepPreviousData,
  })

  // 提取数据
  const columnData = data?.columns ?? []
  const total = data?.total ?? 0
  const pageCount = Math.ceil(total / pagination.pageSize)

  // 使用统一的 CRUD mutations
  const { createMutation, updateMutation, deleteMutation } = useCrudMutations<ColumnType, string>({
    queryKey: ['columnList'],
    createFn: adaptedCreateColumn,
    updateFn: adaptedUpdateColumn,
    deleteFn: deleteColumn,
    messages: {
      createSuccess: t('features.content.column.createSuccess'),
      updateSuccess: t('features.content.column.updateSuccess'),
      deleteSuccess: t('features.content.column.deleteSuccess'),
    },
    onSuccess: () => {
      setDialogOpen(null)
      setEditingColumn(null)
      setDeleteDialogOpen(null)
    },
  })

  // 处理新建
  const handleCreate = useCallback(() => {
    setEditingColumn(null)
    setDialogOpen('create')
  }, [])

  // 处理编辑
  const handleEdit = useCallback((column: ColumnType) => {
    setEditingColumn(column)
    setDialogOpen('edit')
  }, [])

  // 处理删除
  const handleDelete = useCallback(() => {
    if (deleteDialogOpen) {
      deleteMutation.mutate(deleteDialogOpen)
    }
  }, [deleteDialogOpen, deleteMutation])

  // 列定义
  const columns = useMemo<ColumnDef<ColumnType>[]>(
    () => [
      {
        accessorKey: 'columnID',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.content.column.columns.id')}
          />
        ),
        cell: ({ row }) => (
          <div className='font-mono text-muted-foreground'>
            {String(row.original.id ?? '-')}
          </div>
        ),
      },
      {
        accessorKey: 'title',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.content.column.columns.title')}
          />
        ),
        cell: ({ row }) => (
          <div
            className='max-w-[160px] truncate font-medium'
            title={row.getValue('title') as string}
          >
            {(row.getValue('title') as string) ?? '-'}
          </div>
        ),
      },
      {
        accessorKey: 'description',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.content.column.columns.description')}
          />
        ),
        cell: ({ row }) => {
          const desc = row.getValue('description') as string
          return (
            <div
              className='max-w-[240px] truncate text-muted-foreground'
              title={desc}
            >
              {desc ?? '-'}
            </div>
          )
        },
      },
      {
        accessorKey: 'createdAt',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.content.column.columns.createdAt')}
          />
        ),
        cell: ({ row }) => (
          <div className='text-muted-foreground whitespace-nowrap'>
            {(row.getValue('createdAt') as string) ?? '-'}
          </div>
        ),
      },
      {
        id: 'actions',
        header: t('features.content.column.columns.actions'),
        cell: ({ row }) => {
          const columnItem = row.original
          return (
            <DataTableActions
              onEdit={() => handleEdit(columnItem)}
              onDelete={() => setDeleteDialogOpen(String(columnItem.id))}
              deleteConfirmTitle={t('features.content.column.confirmDelete')}
              useDropdown={false}
            />
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
    data: columnData,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
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

  const columnToDelete = columnData.find((c) => String(c.id) === deleteDialogOpen)
  const isDialogLoading = dialogOpen === 'create'
    ? createMutation.isPending
    : updateMutation.isPending

  return (
    <ListPageLayout
      title={t('features.content.column.title')}
      description={t('features.content.column.description')}
      actions={
        <Button onClick={handleCreate}>
          <PlusIcon className='h-4 w-4' />
          {t('features.content.column.createButton')}
        </Button>
      }
      filterContent={
        <div className='flex flex-wrap gap-2'>
          <Input
            type='text'
            placeholder={t('features.content.column.search.title')}
            value={(table.getColumn('title')?.getFilterValue() as string) ?? ''}
            onChange={(e) =>
              table.getColumn('title')?.setFilterValue(e.target.value)
            }
            className='h-8 w-[180px]'
          />
          <Select
            value={projectFilter}
            onValueChange={setProjectFilter}
          >
            <SelectTrigger className='h-8 w-[160px]'>
              <SelectValue
                placeholder={t('features.content.column.search.project')}
              />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value='__all__'>
                {t('features.content.column.search.projectAll')}
              </SelectItem>
              {projectList.map((p) => (
                <SelectItem key={p.projectID} value={(p.projectID ?? '').trim()}>
                  {p.title}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      }
      dialogs={
        <>
          {/* 新建/编辑对话框 */}
          <ColumnFormDialog
        open={dialogOpen !== null}
        onOpenChange={(open) => {
          if (!open) {
            setDialogOpen(null)
            setEditingColumn(null)
          }
        }}
        onSubmit={(payload) => createMutation.mutate(payload)}
        onSubmitUpdate={(columnID, payload) =>
          updateMutation.mutate({ columnID, ...payload } as Partial<ColumnType>)
        }
        isLoading={isDialogLoading}
        projectList={projectList}
        editColumn={editingColumn}
      />

      {/* 删除确认对话框 */}
      <ConfirmDialog
        open={deleteDialogOpen !== null}
        onOpenChange={(open) => {
          if (!open) setDeleteDialogOpen(null)
        }}
        title={t('features.content.column.confirmDelete')}
        desc={t('features.content.column.confirmDeleteDesc', {
          title: columnToDelete?.title ?? deleteDialogOpen ?? '',
        })}
        handleConfirm={handleDelete}
        destructive
        confirmText={t('features.content.column.actions.delete')}
        isLoading={deleteMutation.isPending}
      />
        </>
      }
    >
      <DataTable
        table={table}
        columns={columns}
        isLoading={isLoading}
        error={error}
        loadingText={t('features.content.column.loading')}
        errorText={t('features.content.column.loadError')}
        emptyText={t('features.content.column.noData')}
      />
    </ListPageLayout>
  )
}
