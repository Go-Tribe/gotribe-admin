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
  DataTableColumnHeader,
  DataTable,
} from '@/components/data-table'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { ConfigFormDialog } from './components/config-form-dialog'
import type { Config, ConfigCreateParams, ConfigUpdateParams } from './types/config'
import { getConfigList, getConfig, createConfig, updateConfig, deleteConfig } from './service/config'
import { getProjectList } from '@/shared/api'
import { ListPageLayout, DataTableActions } from '@/components'
import { useI18n } from '@/context/i18n-provider'
import { useDataTable } from '@/hooks/use-data-table'
import { useCrudMutations } from '@/hooks/use-crud-mutations'

export function ContentConfig() {
  'use no memo'
  const { t } = useI18n()
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<string | null>(null)
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [editConfigID, setEditConfigID] = useState<string | null>(null)

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

  const queryParams = useMemo(() => ({
    configID: getFilterValue('configID') || undefined,
    title: getFilterValue('title') || undefined,
    pageNum,
    pageSize: pagination.pageSize,
  }), [getFilterValue, pageNum, pagination.pageSize])

  const { data: projectData } = useQuery({
    queryKey: ['projectList', { current: 1, pageNum: 1, pageSize: 1000 }],
    queryFn: () => getProjectList({ current: 1, pageNum: 1, pageSize: 1000 }),
  })
  const projectList = useMemo(
    () => projectData?.projects?.map((p) => ({ projectID: p.projectID, title: p.title ?? p.projectID })) ?? [],
    [projectData?.projects]
  )

  const { data, isPending: isLoading, error } = useQuery({
    queryKey: ['configList', queryParams],
    queryFn: () => getConfigList(queryParams),
    placeholderData: keepPreviousData,
  })

  const { data: configDetail } = useQuery({
    queryKey: ['configDetail', editConfigID],
    queryFn: () => getConfig(editConfigID!),
    enabled: Boolean(editConfigID?.trim()),
  })

  const configData = data?.configs ?? []
  const total = data?.total ?? 0
  const pageCount = Math.ceil(total / pagination.pageSize)

  // 使用统一的 CRUD mutations
  interface ConfigUpdateInput {
    configID: string
    data: ConfigUpdateParams
  }

  const { createMutation, updateMutation, deleteMutation } = useCrudMutations<Config, string, ConfigCreateParams, ConfigUpdateInput>({
    queryKey: ['configList'],
    createFn: (params) => createConfig(params),
    updateFn: (input) => updateConfig(input.configID, input.data),
    deleteFn: (configIds: string) => deleteConfig(configIds),
    messages: {
      createSuccess: t('features.content.config.createSuccess'),
      updateSuccess: t('features.content.config.updateSuccess'),
      deleteSuccess: t('features.content.config.deleteSuccess'),
    },
    onSuccess: () => {
      setCreateDialogOpen(false)
      setEditConfigID(null)
      setDeleteDialogOpen(null)
    },
  })

  const handleDelete = useCallback(() => {
    if (deleteDialogOpen) {
      deleteMutation.mutate(deleteDialogOpen)
    }
  }, [deleteDialogOpen, deleteMutation])

  const handleEdit = useCallback((config: Config) => {
    setEditConfigID(config.configID)
  }, [])

  const columns = useMemo<ColumnDef<Config>[]>(
    () => [
      {
        accessorKey: 'configID',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.config.columns.id')} />
        ),
        cell: ({ row }) => (
          <div className='font-mono text-muted-foreground'>
            {(row.getValue('configID') as string)?.trim() ?? '-'}
          </div>
        ),
      },
      {
        accessorKey: 'alias',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.config.columns.alias')} />
        ),
        cell: ({ row }) => (
          <div className='max-w-[120px] truncate' title={row.getValue('alias') as string}>
            {(row.getValue('alias') as string) ?? '-'}
          </div>
        ),
      },
      {
        accessorKey: 'title',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.config.columns.title')} />
        ),
        cell: ({ row }) => (
          <div className='max-w-[200px] truncate font-medium' title={row.getValue('title') as string}>
            {(row.getValue('title') as string) ?? '-'}
          </div>
        ),
      },
      {
        accessorKey: 'description',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.config.columns.description')} />
        ),
        cell: ({ row }) => {
          const desc = row.getValue('description') as string
          return (
            <div className='max-w-[280px] truncate text-muted-foreground' title={desc}>
              {desc ?? '-'}
            </div>
          )
        },
      },
      {
        id: 'actions',
        header: t('features.content.config.columns.actions'),
        cell: ({ row }) => (
          <DataTableActions
            onEdit={() => handleEdit(row.original)}
            onDelete={() => setDeleteDialogOpen(row.original.configID)}
            deleteConfirmTitle={t('features.content.config.confirmDelete')}
            useDropdown={false}
          />
        ),
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

  const table = useReactTable({
    data: configData,
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

  const configToDelete = configData.find((c) => c.configID === deleteDialogOpen)
  const isDialogLoading = createDialogOpen
    ? createMutation.isPending
    : updateMutation.isPending

  return (
    <ListPageLayout
      title={t('features.content.config.title')}
      description={t('features.content.config.description')}
      actions={
        <Button onClick={() => setCreateDialogOpen(true)}>
          <PlusIcon className='h-4 w-4' />
          {t('features.content.config.createButton')}
        </Button>
      }
      filterContent={
        <div className='flex flex-wrap gap-2'>
          <Input
            type='text'
            placeholder={t('features.content.config.search.id')}
            value={(table.getColumn('configID')?.getFilterValue() as string) ?? ''}
            onChange={(e) => table.getColumn('configID')?.setFilterValue(e.target.value)}
            className='h-8 w-[160px]'
          />
          <Input
            type='text'
            placeholder={t('features.content.config.search.configName')}
            value={(table.getColumn('title')?.getFilterValue() as string) ?? ''}
            onChange={(e) => table.getColumn('title')?.setFilterValue(e.target.value)}
            className='h-8 w-[180px]'
          />
        </div>
      }
      dialogs={
        <>
          <ConfigFormDialog
            open={createDialogOpen || (Boolean(editConfigID?.trim()) && Boolean(configDetail))}
            onOpenChange={(open) => {
              if (!open) {
                setCreateDialogOpen(false)
                setEditConfigID(null)
              }
            }}
            onSubmit={(payload) => createMutation.mutate(payload)}
            onSubmitUpdate={(configID, payload) => updateMutation.mutate({ configID, data: payload } as ConfigUpdateInput)}
            isLoading={isDialogLoading}
            projectList={projectList}
            editConfig={editConfigID && configDetail ? configDetail : null}
          />
          <ConfirmDialog
            open={deleteDialogOpen !== null}
            onOpenChange={(open) => {
              if (!open) setDeleteDialogOpen(null)
            }}
            title={t('features.content.config.confirmDelete')}
            desc={t('features.content.config.confirmDeleteDesc', { title: configToDelete?.title ?? deleteDialogOpen ?? '' })}
            handleConfirm={handleDelete}
            destructive
            confirmText={t('features.content.config.actions.delete')}
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
        loadingText={t('features.content.config.loading')}
        errorText={t('features.content.config.loadError')}
        emptyText={t('features.content.config.noData')}
        bordered={false}
      />
    </ListPageLayout>
  )
}
// MIGRATED
