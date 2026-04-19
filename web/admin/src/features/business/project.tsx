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
import { DataTableColumnHeader, DataTable } from '@/components/data-table'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { ProjectFormDialog } from './component/project-form-dialog'
import type { Project } from './types/project'
import {
  getProjectList,
  createProject,
  updateProject,
  deleteProject,
} from './service/project'
import { ListPageLayout, DataTableActions } from '@/components'
import { useI18n } from '@/context/i18n-provider'
import { useDataTable } from '@/hooks/use-data-table'
import { useCrudMutations } from '@/hooks/use-crud-mutations'

export function BusinessProject() {
  const { t } = useI18n()
  const [dialogOpen, setDialogOpen] = useState<'create' | 'edit' | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<string | null>(null)
  const [editingProject, setEditingProject] = useState<Project | null>(null)

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
    current: 1,
    title: getFilterValue('title'),
    projectID: getFilterValue('projectID'),
    pageNum,
    pageSize: pagination.pageSize,
  }), [getFilterValue, pagination.pageSize, pageNum])

  // 获取项目列表数据
  const { data, isLoading, error } = useQuery({
    queryKey: ['projectList', queryParams],
    queryFn: () => getProjectList(queryParams),
    placeholderData: keepPreviousData,
  })

  // 提取数据
  const projectData = data?.projects || []
  const total = data?.total || 0
  const pageCount = Math.ceil(total / pagination.pageSize)

  // 使用统一的 CRUD mutations
  const { createMutation, updateMutation, deleteMutation } = useCrudMutations<Project, string>({
    queryKey: ['projectList'],
    createFn: createProject,
    updateFn: updateProject,
    deleteFn: deleteProject,
    messages: {
      createSuccess: t('features.business.project.createSuccess'),
      updateSuccess: t('features.business.project.updateSuccess'),
      deleteSuccess: t('features.business.project.deleteSuccess'),
    },
    onSuccess: () => {
      setDialogOpen(null)
      setEditingProject(null)
      setDeleteDialogOpen(null)
    },
  })

  // 处理新建
  const handleCreate = useCallback(() => {
    setEditingProject(null)
    setDialogOpen('create')
  }, [])

  // 处理编辑
  const handleEdit = useCallback((project: Project) => {
    setEditingProject(project)
    setDialogOpen('edit')
  }, [])

  // 处理删除
  const handleDelete = useCallback(() => {
    if (deleteDialogOpen) {
      deleteMutation.mutate(deleteDialogOpen)
    }
  }, [deleteDialogOpen, deleteMutation])

  // 列定义
  const columns = useMemo<ColumnDef<Project>[]>(
    () => [
      {
        accessorKey: 'projectID',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.business.project.columns.id')} />
        ),
        cell: ({ row }) => (
          <div className='font-medium'>{row.getValue('projectID')}</div>
        ),
      },
      {
        accessorKey: 'title',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.business.project.columns.projectName')} />
        ),
        cell: ({ row }) => <div className='font-medium'>{row.getValue('title')}</div>,
      },
      {
        accessorKey: 'description',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.business.project.columns.projectDescription')} />
        ),
        cell: ({ row }) => {
          const desc = row.getValue('description') as string
          return (
            <div className='max-w-[300px] truncate' title={desc}>
              {desc || '-'}
            </div>
          )
        },
      },
      {
        id: 'actions',
        header: t('features.business.project.columns.actions'),
        cell: ({ row }) => (
          <DataTableActions
            onEdit={() => handleEdit(row.original)}
            onDelete={() => setDeleteDialogOpen(row.original.projectID)}
            deleteConfirmTitle={t('features.business.project.confirmDelete')}
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
    [t, handleEdit],
  )

  // 表格实例
  const table = useReactTable({
    data: projectData,
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

  const projectToDelete = projectData.find(
    (p: Project) => p.projectID === deleteDialogOpen,
  )
  const isDialogLoading = dialogOpen === 'create' 
    ? createMutation.isPending 
    : updateMutation.isPending

  return (
    <ListPageLayout
      title={t('features.business.project.title')}
      description={t('features.business.project.description')}
      actions={
        <Button onClick={handleCreate}>
          <PlusIcon className='h-4 w-4' />
          {t('features.business.project.createButton')}
        </Button>
      }
      filterContent={
        <div className='flex flex-wrap gap-2'>
          <Input
            type='text'
            placeholder={t('features.business.project.search.projectName')}
            value={(table.getColumn('title')?.getFilterValue() as string) ?? ''}
            onChange={(e) =>
              table.getColumn('title')?.setFilterValue(e.target.value)
            }
            className='h-8 w-[150px]'
          />
          <Input
            type='text'
            placeholder={t('features.business.project.search.projectID')}
            value={
              (table.getColumn('projectID')?.getFilterValue() as string) ?? ''
            }
            onChange={(e) =>
              table.getColumn('projectID')?.setFilterValue(e.target.value)
            }
            className='h-8 w-[150px]'
          />
        </div>
      }
      dialogs={
        <>
          <ProjectFormDialog
            open={dialogOpen !== null}
            onOpenChange={(open) => {
              if (!open) {
                setDialogOpen(null)
                setEditingProject(null)
              }
            }}
            project={editingProject}
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
            title={t('features.business.project.confirmDelete')}
            desc={t('features.business.project.confirmDeleteMessage', { title: projectToDelete?.title })}
            handleConfirm={handleDelete}
            destructive
            confirmText={t('features.business.project.actions.delete')}
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
        loadingText={t('features.business.project.loading')}
        errorText={t('features.business.project.loadError')}
        emptyText={t('features.business.project.noData')}
        bordered={false}
      />
    </ListPageLayout>
  )
}
// MIGRATED
