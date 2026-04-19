import { useMemo, useState } from 'react'
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
} from '@/components/data-table'
import { ListPageLayout } from '@/components'
import { useI18n } from '@/context/i18n-provider'
import { getProjectList } from '@/shared/api'
import { getPointList, createPoint } from './service/point'
import { type PointItem, type PointCreateParams } from './types/point'
import { PointFormDialog } from './components/point-form-dialog'
import { useDataTable } from '@/hooks/use-data-table'
import { useCrudMutations } from '@/hooks/use-crud-mutations'

export function OperationPoint() {
  const { t } = useI18n()
  const [createDialogOpen, setCreateDialogOpen] = useState(false)

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

  const [projectFilter, setProjectFilter] = useState<string>('__all__')

  // 获取项目列表用于筛选
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
  const queryParams = useMemo(() => ({
    pageNum,
    pageSize: pagination.pageSize,
    userID: getFilterValue('userID') || undefined,
    nickname: getFilterValue('nickname') || undefined,
    projectId:
      projectFilter && projectFilter !== '__all__'
        ? Number(projectFilter) || undefined
        : undefined,
  }), [pageNum, pagination.pageSize, getFilterValue, projectFilter])

  // 获取积分列表数据
  const { data, isPending: isLoading, error } = useQuery({
    queryKey: ['pointList', queryParams],
    queryFn: () => getPointList(queryParams),
    placeholderData: keepPreviousData,
  })

  // 提取数据
  const pointData = useMemo(() => data?.points || [], [data?.points])
  const total = data?.total || 0
  const pageCount = Math.max(1, Math.ceil(total / pagination.pageSize))

  // 使用统一的 CRUD mutations
  const { createMutation } = useCrudMutations<PointItem, string, PointCreateParams>({
    queryKey: ['pointList'],
    createFn: (params) => createPoint(params),
    updateFn: async () => {},
    deleteFn: async () => {},
    messages: {
      createSuccess: t('features.operation.point.createSuccess'),
    },
    onSuccess: () => {
      setCreateDialogOpen(false)
    },
  })

  const columns = useMemo<ColumnDef<PointItem>[]>(
    () => [
      {
        accessorKey: 'userID',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.point.point.columns.userID')}
          />
        ),
        cell: ({ row }) => (
          <div className='font-mono'>
            {String(row.getValue('userID') ?? '-')}
          </div>
        ),
      },
      {
        accessorKey: 'nickname',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.point.point.columns.nickname')}
          />
        ),
        cell: ({ row }) => (
          <div className='font-medium'>
            {(row.getValue('nickname') as string) ?? '-'}
          </div>
        ),
      },
      {
        accessorKey: 'point',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.point.point.columns.point')}
          />
        ),
        cell: ({ row }) => (
          <div className='font-medium'>
            {row.getValue('point') ?? '-'}
          </div>
        ),
      },
      {
        accessorKey: 'reason',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.point.point.columns.reason')}
          />
        ),
        cell: ({ row }) => {
          const reason = row.getValue('reason') as string
          return (
            <div
              className='max-w-[200px] truncate text-muted-foreground'
              title={reason}
            >
              {reason ?? '-'}
            </div>
          )
        },
      },
      {
        accessorKey: 'createdAt',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.point.point.columns.createdAt')}
          />
        ),
        cell: ({ row }) => (
          <div className='text-muted-foreground whitespace-nowrap'>
            {(row.getValue('createdAt') as string) ?? '-'}
          </div>
        ),
      },
    ],
    [t]
  )

  // 表格实例（服务器端分页，不使用客户端过滤和分页）
  const table = useReactTable({
    data: pointData,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    manualPagination: true,
    pageCount,
    getRowId: (row) => String(row.id ?? ''),
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

  return (
    <ListPageLayout
      title={t('features.operation.point.title')}
      description={t('features.operation.point.description')}
      actions={
        <Button onClick={() => setCreateDialogOpen(true)}>
          <PlusIcon className='h-4 w-4' />
          {t('features.operation.point.point.createButton')}
        </Button>
      }
      filterContent={
        <div className='flex flex-wrap gap-2'>
          <Input
            type='text'
            placeholder={t(
              'features.operation.point.point.search.userIDPlaceholder'
            )}
            value={(table.getColumn('userID')?.getFilterValue() as string) ?? ''}
            onChange={(e) =>
              table.getColumn('userID')?.setFilterValue(e.target.value)
            }
            className='h-8 w-[200px]'
          />
          <Input
            type='text'
            placeholder={t(
              'features.operation.point.point.search.nicknamePlaceholder'
            )}
            value={(table.getColumn('nickname')?.getFilterValue() as string) ?? ''}
            onChange={(e) =>
              table.getColumn('nickname')?.setFilterValue(e.target.value)
            }
            className='h-8 w-[200px]'
          />
          <Select value={projectFilter} onValueChange={setProjectFilter}>
            <SelectTrigger className='h-8 w-[160px]'>
              <SelectValue
                placeholder={t(
                  'features.operation.point.point.search.project'
                )}
              />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value='__all__'>
                {t('features.operation.point.point.search.projectAll')}
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
        <PointFormDialog
          open={createDialogOpen}
          onOpenChange={(open) => {
            if (!open) setCreateDialogOpen(false)
          }}
          onSubmit={(payload) => createMutation.mutate(payload)}
          isLoading={createMutation.isPending}
          projectList={projectList}
        />
      }
    >
      <DataTable
        table={table}
        columns={columns}
        isLoading={isLoading}
        error={error}
        loadingText={t('features.operation.point.point.loading')}
        errorText={t('features.operation.point.point.loadError')}
        emptyText={t('features.operation.point.point.noData')}
        bordered={false}
      />
    </ListPageLayout>
  )
}
// MIGRATED
