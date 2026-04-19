import { useMemo, useState, useCallback } from 'react'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import {
  useReactTable,
  getCoreRowModel,
  type ColumnDef,
} from '@tanstack/react-table'
import { PlusIcon } from '@radix-ui/react-icons'
import { Button } from '@/components/ui/button'
import { DataTableFacetedFilter, DataTableColumnHeader, DataTable } from '@/components/data-table'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { AdminFormDialog } from './components/admin-form-dialog'
import type { Admin } from './types/admin'
import { getAdminList, createAdmin, updateAdmin, deleteAdmin } from './service/index'
import { getRoleList } from './service/role'
import { useI18n } from '@/context/i18n-provider'
import { useDataTable } from '@/hooks/use-data-table'
import { useCrudMutations } from '@/hooks/use-crud-mutations'
// 新组件导入
import { ListPageLayout, DebouncedInput, DataTableActions, StatusBadge } from '@/components'

export function SystemAdmin() {
  const { t } = useI18n()
  const [dialogOpen, setDialogOpen] = useState<'create' | 'edit' | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<number | null>(null)
  const [editingAdmin, setEditingAdmin] = useState<Admin | null>(null)

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
    username: getFilterValue('username'),
    nickname: getFilterValue('nickname'),
    mobile: getFilterValue('mobile'),
    status: getFilterValue('status'),
    pageNum,
    pageSize: pagination.pageSize,
    ...getSortingParams(),
  }), [getFilterValue, getSortingParams, pagination.pageSize, pageNum])

  // 获取管理员列表数据
  const { data, isPending: isLoading, error } = useQuery({
    queryKey: ['adminList', queryParams],
    queryFn: () => getAdminList(queryParams),
    placeholderData: keepPreviousData,
  })

  // 获取角色列表数据（缓存）
  const { data: roleData } = useQuery({
    queryKey: ['roleList'],
    queryFn: () => getRoleList()
  })

  // 提取数据
  const adminData = data?.admins || []
  const total = data?.total || 0
  const roleList = roleData?.roles || []
  const pageCount = Math.ceil(total / pagination.pageSize)

  // 使用统一的 CRUD mutations
  const { createMutation, updateMutation, deleteMutation } = useCrudMutations<Admin, number>({
    queryKey: ['adminList'],
    createFn: createAdmin,
    updateFn: updateAdmin,
    deleteFn: deleteAdmin,
    messages: {
      createSuccess: t('features.system.admin.createSuccess'),
      updateSuccess: t('features.system.admin.updateSuccess'),
      deleteSuccess: t('features.system.admin.deleteSuccess'),
    },
    onSuccess: () => {
      setDialogOpen(null)
      setEditingAdmin(null)
      setDeleteDialogOpen(null)
    },
  })

  // 处理新建
  const handleCreate = useCallback(() => {
    setEditingAdmin(null)
    setDialogOpen('create')
  }, [])

  // 处理编辑
  const handleEdit = useCallback((admin: Admin) => {
    setEditingAdmin(admin)
    setDialogOpen('edit')
  }, [])

  // 处理删除
  const handleDelete = useCallback(() => {
    if (deleteDialogOpen) {
      deleteMutation.mutate(deleteDialogOpen)
    }
  }, [deleteDialogOpen, deleteMutation])

  // 列定义
  const columns = useMemo<ColumnDef<Admin>[]>(
    () => [
      {
        accessorKey: 'username',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.admin.username')} />
        ),
        cell: ({ row }) => <div className='font-medium'>{row.getValue('username')}</div>,
      },
      {
        accessorKey: 'nickname',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.admin.nickname')} />
        ),
        cell: ({ row }) => <div>{row.getValue('nickname')}</div>,
      },
      {
        accessorKey: 'status',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.admin.status')} />
        ),
        cell: ({ row }) => (
          <StatusBadge
            value={row.getValue('status') as number}
            type="enabledStatus"
          />
        ),
      },
      {
        accessorKey: 'mobile',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.admin.mobile')} />
        ),
        cell: ({ row }) => <div>{row.getValue('mobile')}</div>,
      },
      {
        accessorKey: 'creator',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.admin.creator')} />
        ),
        cell: ({ row }) => <div>{row.getValue('creator')}</div>,
      },
      {
        accessorKey: 'introduction',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.admin.introduction')} />
        ),
        cell: ({ row }) => {
          const intro = row.getValue('introduction') as string
          return (
            <div className='max-w-[200px] truncate' title={intro}>
              {intro || '-'}
            </div>
          )
        },
      },
      {
        id: 'actions',
        header: t('features.system.admin.actions'),
        cell: ({ row }) => (
          <DataTableActions
            onEdit={() => handleEdit(row.original)}
            onDelete={() => setDeleteDialogOpen(row.original.id)}
            deleteConfirmTitle={t('features.system.admin.confirmDelete')}
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
    data: adminData,
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

  const adminToDelete = adminData.find((a: Admin) => a.id === deleteDialogOpen)
  const isDialogLoading = dialogOpen === 'create' 
    ? createMutation.isPending 
    : updateMutation.isPending

  return (
    <ListPageLayout
      title={t('features.system.admin.title')}
      description={t('features.system.admin.description')}
      cardPadding="md"
      actions={
        <Button onClick={handleCreate}>
          <PlusIcon className='h-4 w-4' />
          {t('features.system.admin.createButton')}
        </Button>
      }
      filterContent={
        <div className='flex flex-wrap gap-2 mb-4'>
          <DebouncedInput
            placeholder={t('features.system.admin.searchUsername')}
            value={(table.getColumn('username')?.getFilterValue() as string) ?? ''}
            onChange={(value) => table.getColumn('username')?.setFilterValue(value)}
            delay={300}
            showSearchIcon={false}
            showClearButton={false}
            className='h-8 w-[150px]'
          />
          <DebouncedInput
            placeholder={t('features.system.admin.searchNickname')}
            value={(table.getColumn('nickname')?.getFilterValue() as string) ?? ''}
            onChange={(value) => table.getColumn('nickname')?.setFilterValue(value)}
            delay={300}
            showSearchIcon={false}
            showClearButton={false}
            className='h-8 w-[150px]'
          />
          <DebouncedInput
            placeholder={t('features.system.admin.searchMobile')}
            value={(table.getColumn('mobile')?.getFilterValue() as string) ?? ''}
            onChange={(value) => table.getColumn('mobile')?.setFilterValue(value)}
            delay={300}
            showSearchIcon={false}
            showClearButton={false}
            className='h-8 w-[150px]'
          />
          <DataTableFacetedFilter
            column={table.getColumn('status')}
            title={t('features.system.admin.status')}
            options={[
              { label: t('features.system.admin.enabled'), value: '1' },
              { label: t('features.system.admin.disabled'), value: '2' },
            ]}
            single
          />
        </div>
      }
      dialogs={
        <>
          <AdminFormDialog
            open={dialogOpen !== null}
            onOpenChange={(open) => {
              if (!open) {
                setDialogOpen(null)
                setEditingAdmin(null)
              }
            }}
            admin={editingAdmin}
            onSubmit={(data) => {
              if (dialogOpen === 'create') {
                createMutation.mutate(data)
              } else {
                updateMutation.mutate(data)
              }
            }}
            isLoading={isDialogLoading}
            roles={roleList}
          />
          <ConfirmDialog
            open={deleteDialogOpen !== null}
            onOpenChange={(open) => {
              if (!open) setDeleteDialogOpen(null)
            }}
            title={t('features.system.admin.confirmDelete')}
            desc={t('features.system.admin.confirmDeleteMessage', { username: adminToDelete?.username })}
            handleConfirm={handleDelete}
            destructive
            confirmText={t('features.system.admin.delete')}
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
        loadingText={t('features.system.admin.loading')}
        errorText={t('features.system.admin.loadError')}
        emptyText={t('features.system.admin.noData')}
        bordered={false}
      />
    </ListPageLayout>
  )
}
