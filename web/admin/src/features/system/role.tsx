import { useMemo, useState, useCallback } from 'react'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import {
  useReactTable,
  getCoreRowModel,
  type ColumnDef,
} from '@tanstack/react-table'
import { PlusIcon, Pencil1Icon, TrashIcon } from '@radix-ui/react-icons'
import { FolderCog, Settings } from 'lucide-react'
import { ListPageLayout } from '@/components'
import { Button } from '@/components/ui/button'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { Input } from '@/components/ui/input'
import { DataTableFacetedFilter, DataTableColumnHeader, DataTable } from '@/components/data-table'
import { Badge } from '@/components/ui/badge'
import { getStatusVariantByType } from '@/config/status-variants'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { RoleFormDialog } from './components/role-form-dialog'
import { RolePermissionDialog } from './components/role-permission-menu-dialog'
import { RoleApiPermissionDialog } from './components/role-permission-api-dialog'
import type { Role } from './types/admin'
import { getRoleList, createRole, updateRole, deleteRole } from './service/role'
import { useI18n } from '@/context/i18n-provider'
import { useDataTable } from '@/hooks/use-data-table'
import { useCrudMutations } from '@/hooks/use-crud-mutations'

export function SystemRole() {
  const { t } = useI18n()
  const [dialogOpen, setDialogOpen] = useState<'create' | 'edit' | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<number | null>(null)
  const [permissionDialogOpen, setPermissionDialogOpen] = useState<number | null>(null)
  const [apiPermissionDialogOpen, setApiPermissionDialogOpen] = useState<number | null>(null)
  const [editingRole, setEditingRole] = useState<Role | null>(null)

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
    name: getFilterValue('name'),
    keyword: getFilterValue('keyword'),
    status: getFilterValue('status'),
    pageNum,
    pageSize: pagination.pageSize,
    ...getSortingParams(),
  }), [getFilterValue, getSortingParams, pagination.pageSize, pageNum])

  // 获取角色列表数据
  const { data, isPending: isLoading, error } = useQuery({
    queryKey: ['roleList', queryParams],
    queryFn: () => getRoleList(queryParams),
    placeholderData: keepPreviousData,
  })

  // 提取数据
  const roleData = data?.roles || []
  const total = data?.total || 0
  const pageCount = Math.ceil(total / pagination.pageSize)

  // 使用统一的 CRUD mutations
  const { createMutation, updateMutation, deleteMutation } = useCrudMutations<Role, number>({
    queryKey: ['roleList'],
    createFn: createRole,
    updateFn: updateRole,
    deleteFn: deleteRole,
    messages: {
      createSuccess: t('features.system.role.createSuccess'),
      updateSuccess: t('features.system.role.updateSuccess'),
      deleteSuccess: t('features.system.role.deleteSuccess'),
    },
    onSuccess: () => {
      setDialogOpen(null)
      setEditingRole(null)
      setDeleteDialogOpen(null)
    },
  })

  // 处理新建
  const handleCreate = useCallback(() => {
    setEditingRole(null)
    setDialogOpen('create')
  }, [])

  // 处理编辑
  const handleEdit = useCallback((role: Role) => {
    setEditingRole(role)
    setDialogOpen('edit')
  }, [])

  // 处理删除
  const handleDelete = useCallback(() => {
    if (deleteDialogOpen) {
      deleteMutation.mutate(deleteDialogOpen)
    }
  }, [deleteDialogOpen, deleteMutation])

  // 列定义
  const columns = useMemo<ColumnDef<Role>[]>(
    () => [
      {
        accessorKey: 'name',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.role.name')} />
        ),
        cell: ({ row }) => <div className='font-medium'>{row.getValue('name')}</div>,
      },
      {
        accessorKey: 'keyword',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.role.keyword')} />
        ),
        cell: ({ row }) => <div>{row.getValue('keyword')}</div>,
      },
      {
        accessorKey: 'sort',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.role.sort')} />
        ),
        cell: ({ row }) => <div>{row.getValue('sort')}</div>,
      },
      {
        accessorKey: 'status',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.role.status')} />
        ),
        cell: ({ row }) => {
          const status = row.getValue('status') as number
          return (
            <Badge variant={getStatusVariantByType(status, 'enabledStatus')}>
              {status === 1 ? t('features.system.role.enabled') : t('features.system.role.disabled')}
            </Badge>
          )
        },
      },
      {
        accessorKey: 'creator',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.role.creator')} />
        ),
        cell: ({ row }) => <div>{row.getValue('creator')}</div>,
      },
      {
        accessorKey: 'desc',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.role.desc')} />
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
        id: 'actions',
        header: t('features.system.role.actions'),
        cell: ({ row }) => {
          const role = row.original
          return (
            <div className='flex items-center gap-1'>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant='outline'
                    size='icon'
                    className='h-8 w-8 border-border/60'
                    onClick={() => handleEdit(role)}
                  >
                    <Pencil1Icon className='h-4 w-4' />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>
                  {t('features.system.role.edit')}
                </TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant='outline'
                    size='icon'
                    className='h-8 w-8 border-border/60'
                    onClick={() => {
                      setPermissionDialogOpen(role.id)
                    }}
                  >
                    <Settings className='h-4 w-4' />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>
                  {t('features.system.role.permissionSetting')}
                </TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant='outline'
                    size='icon'
                    className='h-8 w-8 border-border/60'
                    onClick={() => {
                      setApiPermissionDialogOpen(role.id)
                    }}
                  >
                    <FolderCog className='h-4 w-4' />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>
                  {t('features.system.role.apiSetting')}
                </TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant='ghost'
                    size='icon'
                    className='h-8 w-8 text-destructive hover:text-destructive'
                    onClick={() => setDeleteDialogOpen(role.id)}
                  >
                    <TrashIcon className='h-4 w-4 text-destructive' />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>
                  {t('features.system.role.delete')}
                </TooltipContent>
              </Tooltip>
            </div>
          )
        },
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
    data: roleData,
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

  const roleToDelete = roleData.find((r: Role) => r.id === deleteDialogOpen)
  const isDialogLoading = dialogOpen === 'create'
    ? createMutation.isPending
    : updateMutation.isPending

  return (
    <ListPageLayout
      title={t('features.system.role.title')}
      description={t('features.system.role.description')}
      actions={
        <Button onClick={handleCreate}>
          <PlusIcon className='h-4 w-4' />
          {t('features.system.role.createButton')}
        </Button>
      }
      filterContent={
        <div className='flex flex-wrap gap-2'>
            <Input
              type='text'
              placeholder={t('features.system.role.searchName')}
              value={(table.getColumn('name')?.getFilterValue() as string) ?? ''}
              onChange={(e) =>
                table.getColumn('name')?.setFilterValue(e.target.value)
              }
              className='h-8 w-[150px]'
            />
            <Input
              type='text'
              placeholder={t('features.system.role.searchKeyword')}
              value={(table.getColumn('keyword')?.getFilterValue() as string) ?? ''}
              onChange={(e) =>
                table.getColumn('keyword')?.setFilterValue(e.target.value)
              }
              className='h-8 w-[150px]'
            />
            <DataTableFacetedFilter
              column={table.getColumn('status')}
              title={t('features.system.role.status')}
              options={[
                { label: t('features.system.role.enabled'), value: '1' },
                { label: t('features.system.role.disabled'), value: '2' },
              ]}
              single
            />
          </div>
      }
      dialogs={
        <>
          <RoleFormDialog
            open={dialogOpen !== null}
            onOpenChange={(open) => {
              if (!open) {
                setDialogOpen(null)
                setEditingRole(null)
              }
            }}
            role={editingRole}
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
            title={t('features.system.role.confirmDelete')}
            desc={t('features.system.role.confirmDeleteMessage', { name: roleToDelete?.name })}
            handleConfirm={handleDelete}
            destructive
            confirmText={t('features.system.role.delete')}
            isLoading={deleteMutation.isPending}
          />
          <RolePermissionDialog
            open={permissionDialogOpen !== null}
            onOpenChange={(open) => {
              if (!open) setPermissionDialogOpen(null)
            }}
            roleId={permissionDialogOpen}
            roleName={roleData.find((r: Role) => r.id === permissionDialogOpen)?.name}
          />
          <RoleApiPermissionDialog
            open={apiPermissionDialogOpen !== null}
            onOpenChange={(open) => {
              if (!open) setApiPermissionDialogOpen(null)
            }}
            roleId={apiPermissionDialogOpen}
            roleName={roleData.find((r: Role) => r.id === apiPermissionDialogOpen)?.name}
          />
        </>
      }
    >
      <DataTable<Role>
        table={table}
        columns={columns}
        isLoading={isLoading}
        error={error}
        loadingText={t('features.system.role.loading')}
        errorText={t('features.system.role.loadError')}
        emptyText={t('features.system.role.noData')}
        bordered={false}
      />
    </ListPageLayout>
  )
}
