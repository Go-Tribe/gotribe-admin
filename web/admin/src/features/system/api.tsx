import { useMemo, useState, useCallback } from 'react'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import {
  useReactTable,
  getCoreRowModel,
  type ColumnDef,
} from '@tanstack/react-table'
import { PlusIcon } from '@radix-ui/react-icons'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { DataTableFacetedFilter, DataTableColumnHeader, DataTable } from '@/components/data-table'
import { Badge } from '@/components/ui/badge'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { ApiFormDialog } from './component/api-form-dialog'
import type { Api } from './types/api'
import { getApiList, createApi, updateApi, deleteApi } from './service/api'
import { useI18n } from '@/context/i18n-provider'
import { useDataTable } from '@/hooks/use-data-table'
import { useCrudMutations } from '@/hooks/use-crud-mutations'
import { ListPageLayout, DataTableActions } from '@/components'

const HTTP_METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS']

export function SystemApi() {
  const { t } = useI18n()
  const [dialogOpen, setDialogOpen] = useState<'create' | 'edit' | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<number | null>(null)
  const [editingApi, setEditingApi] = useState<Api | null>(null)

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
    path: getFilterValue('path'),
    category: getFilterValue('category'),
    method: getFilterValue('method'),
    creator: getFilterValue('creator'),
    pageNum,
    pageSize: pagination.pageSize,
    ...getSortingParams(),
  }), [getFilterValue, getSortingParams, pagination.pageSize, pageNum])

  // 获取API列表数据
  const { data, isPending: isLoading, error } = useQuery({
    queryKey: ['apiList', queryParams],
    queryFn: () => getApiList(queryParams),
    placeholderData: keepPreviousData,
  })

  // 提取数据
  const apiData = data?.apis || []
  const total = data?.total || 0
  const pageCount = Math.ceil(total / pagination.pageSize)

  // 使用统一的 CRUD mutations
  const { createMutation, updateMutation, deleteMutation } = useCrudMutations<Api, number>({
    queryKey: ['apiList'],
    createFn: createApi,
    updateFn: updateApi,
    deleteFn: deleteApi,
    messages: {
      createSuccess: t('features.system.api.createSuccess'),
      updateSuccess: t('features.system.api.updateSuccess'),
      deleteSuccess: t('features.system.api.deleteSuccess'),
    },
    onSuccess: () => {
      setDialogOpen(null)
      setEditingApi(null)
      setDeleteDialogOpen(null)
    },
  })

  // 处理新建
  const handleCreate = useCallback(() => {
    setEditingApi(null)
    setDialogOpen('create')
  }, [])

  // 处理编辑
  const handleEdit = useCallback((api: Api) => {
    setEditingApi(api)
    setDialogOpen('edit')
  }, [])

  // 处理删除
  const handleDelete = useCallback(() => {
    if (deleteDialogOpen) {
      deleteMutation.mutate(deleteDialogOpen)
    }
  }, [deleteDialogOpen, deleteMutation])

  // 列定义
  const columns = useMemo<ColumnDef<Api>[]>(
    () => [
      {
        accessorKey: 'method',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.api.method')} />
        ),
        cell: ({ row }) => {
          const method = row.getValue('method') as string
          const methodColors: Record<string, string> = {
            GET: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300',
            POST: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
            PUT: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300',
            PATCH: 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300',
            DELETE: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300',
          }
          const colorClass = methodColors[method] || 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300'
          return (
            <Badge className={colorClass}>
              {method}
            </Badge>
          )
        },
      },
      {
        accessorKey: 'path',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.api.path')} />
        ),
        cell: ({ row }) => <div className='font-mono text-sm'>{row.getValue('path')}</div>,
      },
      {
        accessorKey: 'category',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.api.category')} />
        ),
        cell: ({ row }) => <div>{row.getValue('category')}</div>,
      },
      {
        accessorKey: 'desc',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.api.desc')} />
        ),
        cell: ({ row }) => {
          const desc = row.getValue('desc') as string
          return (
            <div className='max-w-[200px] truncate' title={desc}>
              {desc || '-'}
            </div>
          )
        },
      },
      {
        accessorKey: 'creator',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.api.creator')} />
        ),
        cell: ({ row }) => <div>{row.getValue('creator')}</div>,
      },
      {
        id: 'actions',
        header: t('features.system.api.actions'),
        cell: ({ row }) => (
          <DataTableActions
            onEdit={() => handleEdit(row.original)}
            onDelete={() => setDeleteDialogOpen(row.original.id)}
            deleteConfirmTitle={t('features.system.api.confirmDelete')}
            useDropdown={false}
          />
        ),
        enableHiding: false,
        meta: {
          className: 'sticky right-0 bg-background z-10 shadow-[inset_-1px_0_0_0_hsl(var(--border))]',
          thClassName: 'sticky right-0 bg-background z-10 shadow-[inset_-1px_0_0_0_hsl(var(--border))]',
          tdClassName: 'sticky right-0 bg-background z-10 shadow-[inset_-1px_0_0_0_hsl(var(--border))]',
        }
      },
    ],
    [t, handleEdit]
  )

  // 表格实例
  const table = useReactTable({
    data: apiData,
    columns,
    getCoreRowModel: getCoreRowModel(),
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

  const apiToDelete = apiData.find((a: Api) => a.id === deleteDialogOpen)
  const isDialogLoading = dialogOpen === 'create' 
    ? createMutation.isPending 
    : updateMutation.isPending

  return (
    <ListPageLayout
        title={t('features.system.api.title')}
        description={t('features.system.api.description')}
        actions={
          <Button onClick={handleCreate}>
            <PlusIcon className='h-4 w-4' />
            {t('features.system.api.createButton')}
          </Button>
        }
        filterContent={
          <div className='flex flex-wrap gap-2'>
            <Input
              type='text'
              placeholder={t('features.system.api.searchPath')}
              value={(table.getColumn('path')?.getFilterValue() as string) ?? ''}
              onChange={(e) =>
                table.getColumn('path')?.setFilterValue(e.target.value)
              }
              className='h-8 w-[150px]'
            />
            <Input
              type='text'
              placeholder={t('features.system.api.searchCategory')}
              value={(table.getColumn('category')?.getFilterValue() as string) ?? ''}
              onChange={(e) =>
                table.getColumn('category')?.setFilterValue(e.target.value)
              }
              className='h-8 w-[150px]'
            />
            <Input
              type='text'
              placeholder={t('features.system.api.searchCreator')}
              value={(table.getColumn('creator')?.getFilterValue() as string) ?? ''}
              onChange={(e) =>
                table.getColumn('creator')?.setFilterValue(e.target.value)
              }
              className='h-8 w-[150px]'
            />
            <DataTableFacetedFilter
              column={table.getColumn('method')}
              title={t('features.system.api.method')}
              options={HTTP_METHODS.map((method) => ({ label: method, value: method }))}
              single
            />
          </div>
        }
        dialogs={
          <>
            <ApiFormDialog
              open={dialogOpen !== null}
              onOpenChange={(open) => {
                if (!open) {
                  setDialogOpen(null)
                  setEditingApi(null)
                }
              }}
              api={editingApi}
              onSubmit={(data) => {
                if (dialogOpen === 'create') {
                  createMutation.mutate(data)
                } else {
                  updateMutation.mutate(data)
                }
              }}
              isLoading={isDialogLoading}
            />
            <ConfirmDialog
              open={deleteDialogOpen !== null}
              onOpenChange={(open) => {
                if (!open) setDeleteDialogOpen(null)
              }}
              title={t('features.system.api.confirmDelete')}
              desc={t('features.system.api.confirmDeleteMessage', { path: apiToDelete?.path })}
              handleConfirm={handleDelete}
              destructive
              confirmText={t('features.system.api.delete')}
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
          loadingText={t('features.system.api.loading')}
          errorText={t('features.system.api.loadError')}
          emptyText={t('features.system.api.noData')}
          bordered={false}
        />
      </ListPageLayout>
  )
}
