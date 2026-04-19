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
import { ListPageLayout } from '@/components/layout'
import { getProjectList } from '@/shared/api'
import { useI18n } from '@/context/i18n-provider'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { SceneFormDialog } from './components/scene-form-dialog'
import { createScene, getSceneList, updateScene, deleteScene } from './service'
import { useDataTable } from '@/hooks/use-data-table'
import { useCrudMutations } from '@/hooks/use-crud-mutations'
import type { Scene, SceneCreateParams, SceneUpdateParams } from './types'

/** 表格行（由 API 数据映射） */
type SceneRow = {
  id: string
  name: string
  description: string
  project: string
  createdAt: string
}

function mapSceneToRow(
  s: Scene,
  projectList: { id: number; title: string }[]
): SceneRow {
  const id = String(s.id ?? '')
  const projectId = s.projectId ?? 0
  const projectTitle = (s.projectTitle ?? '').trim()
  const project =
    projectTitle ||
    (projectList.find((p) => p.id === projectId)?.title ??
      '') ||
    String(projectId) ||
    '-'
  return {
    id,
    name: (s.title ?? '').trim() || '-',
    description: (s.description ?? '').trim() || '-',
    project: project || '-',
    createdAt: (s.createdAt ?? '').trim() || '-',
  }
}

export function OperationScene() {
  'use no memo'
  const { t } = useI18n()
  const [dialogOpen, setDialogOpen] = useState<'create' | 'edit' | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<string | null>(null)
  const [editingScene, setEditingScene] = useState<Scene | null>(null)
  const [projectFilter, setProjectFilter] = useState<string>('__all__')

  // 使用统一的表格状态管理
  const {
    pagination,
    setPagination,
    sorting,
    setSorting,
    columnVisibility,
    setColumnVisibility,
    pageNum,
  } = useDataTable()

  const { data: projectData } = useQuery({
    queryKey: ['projectList', { current: 1, pageNum: 1, pageSize: 1000 }],
    queryFn: () => getProjectList({ current: 1, pageNum: 1, pageSize: 1000 }),
  })
  const projectList = useMemo(
    () =>
      projectData?.projects?.map((p) => ({
        id: p.id,
        title: p.title ?? String(p.id),
      })) ?? [],
    [projectData?.projects]
  )

  // 构建查询参数
  const queryParams = useMemo(
    () => ({
      pageNum,
      pageSize: pagination.pageSize,
      projectId:
        projectFilter && projectFilter !== '__all__'
          ? Number(projectFilter) || undefined
          : undefined,
    }),
    [pageNum, pagination.pageSize, projectFilter]
  )

  const { data: sceneListData, isPending: isLoading, error } = useQuery({
    queryKey: ['adSceneList', queryParams],
    queryFn: () => getSceneList(queryParams),
    placeholderData: keepPreviousData,
  })

  const rawList = sceneListData?.adScenes ?? []
  const scenarioData = useMemo(
    () => rawList.map((s: Scene) => mapSceneToRow(s, projectList)),
    [rawList, projectList]
  )
  const total = sceneListData?.total ?? 0
  const pageCount = Math.max(1, Math.ceil(total / pagination.pageSize))

  // 包装创建函数以适配 useCrudMutations 签名
  const createSceneWrapper = useCallback(async (data: unknown) => {
    const params = data as SceneCreateParams
    return createScene(params)
  }, [])

  // 包装更新函数以适配 useCrudMutations 签名
  const updateSceneWrapper = useCallback(async (data: unknown) => {
    const params = data as SceneUpdateParams & { id: number }
    if (params.id == null) throw new Error('Scene ID is required')
    const { id, ...updateData } = params
    return updateScene(String(id), updateData)
  }, [])

  // 使用统一的 CRUD mutations
  const { createMutation, updateMutation, deleteMutation } = useCrudMutations<unknown, string>({
    queryKey: ['adSceneList'],
    createFn: createSceneWrapper,
    updateFn: updateSceneWrapper,
    deleteFn: deleteScene,
    messages: {
      createSuccess: t('features.operation.scene.createSuccess'),
      updateSuccess: t('features.operation.scene.updateSuccess'),
      deleteSuccess: t('features.operation.scene.deleteSuccess'),
    },
    onSuccess: () => {
      setDialogOpen(null)
      setEditingScene(null)
      setDeleteDialogOpen(null)
    },
  })

  const handleCreate = useCallback(() => {
    setEditingScene(null)
    setDialogOpen('create')
  }, [])

  const handleEdit = useCallback((scene: Scene) => {
    setEditingScene(scene)
    setDialogOpen('edit')
  }, [])

  const handleDelete = useCallback(() => {
    if (deleteDialogOpen) {
      deleteMutation.mutate(deleteDialogOpen)
    }
  }, [deleteDialogOpen, deleteMutation])

  const sceneToDelete = scenarioData.find((r: SceneRow) => r.id === deleteDialogOpen)

  const columns = useMemo<ColumnDef<SceneRow>[]>(
    () => [
      {
        accessorKey: 'id',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.scene.scenario.columns.id')}
          />
        ),
        cell: ({ row }) => (
          <div className='font-mono text-muted-foreground'>
            {(row.getValue('id') as string)?.trim() ?? '-'}
          </div>
        ),
      },
      {
        accessorKey: 'name',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.scene.scenario.columns.name')}
          />
        ),
        cell: ({ row }) => (
          <div
            className='max-w-[160px] truncate font-medium'
            title={row.getValue('name') as string}
          >
            {(row.getValue('name') as string) ?? '-'}
          </div>
        ),
      },
      {
        accessorKey: 'description',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.scene.scenario.columns.description')}
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
        accessorKey: 'project',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.scene.scenario.columns.project')}
          />
        ),
        cell: ({ row }) => (
          <div className='text-muted-foreground whitespace-nowrap'>
            {(row.getValue('project') as string) ?? '-'}
          </div>
        ),
      },
      {
        accessorKey: 'createdAt',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.scene.scenario.columns.createdAt')}
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
        header: t('features.operation.scene.scenario.columns.actions'),
        cell: ({ row }) => {
          const sceneRow = row.original
          const scene = rawList.find((s: Scene) => String(s.id ?? '') === sceneRow.id)
          return (
            <DataTableActions
              onEdit={() => scene && handleEdit(scene)}
              onDelete={() => setDeleteDialogOpen(sceneRow.id)}
              deleteConfirmTitle={t('features.operation.scene.confirmDelete')}
              useDropdown={false}
            />
          )
        },
        enableHiding: false,
        meta: {
          thClassName:
            'sticky right-0 bg-background z-10 shadow-[inset_-1px_0_0_0_hsl(var(--border))]',
          tdClassName:
            'sticky right-0 bg-background z-10 shadow-[inset_-1px_0_0_0_hsl(var(--border))]',
        },
      },
    ],
    [t, rawList, handleEdit]
  )

  const table = useReactTable({
    data: scenarioData,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    manualPagination: true,
    pageCount,
    getRowId: (row) => row.id,
    onSortingChange: setSorting,
    onColumnVisibilityChange: setColumnVisibility,
    onPaginationChange: setPagination,
    state: {
      sorting,
      columnVisibility,
      pagination,
    },
  })

  const isDialogLoading = dialogOpen === 'create' 
    ? createMutation.isPending 
    : updateMutation.isPending

  return (
    <ListPageLayout
      title={t('features.operation.scene.title')}
      description={t('features.operation.scene.description')}
      actions={
        <Button onClick={handleCreate}>
          <PlusIcon className='h-4 w-4' />
          {t('features.operation.scene.scenario.createButton')}
        </Button>
      }
      filterContent={
        <div className='flex flex-wrap gap-2'>
          <Select value={projectFilter} onValueChange={setProjectFilter}>
            <SelectTrigger className='h-8 w-[160px]'>
              <SelectValue
                placeholder={t(
                  'features.operation.scene.scenario.search.project'
                )}
              />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value='__all__'>
                {t('features.operation.scene.scenario.search.projectAll')}
              </SelectItem>
              {projectList.map((p) => (
                <SelectItem
                  key={p.id}
                  value={String(p.id)}
                >
                  {p.title}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      }
      dialogs={
        <>
          <SceneFormDialog
        open={dialogOpen !== null}
        onOpenChange={(open) => {
          if (!open) {
            setDialogOpen(null)
            setEditingScene(null)
          }
        }}
        onSubmit={(payload) => createMutation.mutate(payload)}
        onSubmitUpdate={(adSceneID, payload) =>
          updateMutation.mutate({
            id: Number(adSceneID),
            ...payload,
          })
        }
        isLoading={isDialogLoading}
        projectList={projectList}
        editScene={editingScene}
      />

      <ConfirmDialog
        open={deleteDialogOpen !== null}
        onOpenChange={(open) => {
          if (!open) setDeleteDialogOpen(null)
        }}
        title={t('features.operation.scene.confirmDelete')}
        desc={t('features.operation.scene.confirmDeleteDesc', {
          name: sceneToDelete?.name ?? deleteDialogOpen ?? '',
        })}
        handleConfirm={handleDelete}
        destructive
        confirmText={t('features.operation.scene.scenario.actions.delete')}
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
        loadingText={t('features.operation.scene.scenario.loading')}
        errorText={t('features.operation.scene.scenario.loadError')}
        emptyText={t('features.operation.scene.scenario.noData')}
      />
    </ListPageLayout>
  )
}
