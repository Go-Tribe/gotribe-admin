import { useMemo } from 'react'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import {
  useReactTable,
  getCoreRowModel,
  type ColumnDef,
} from '@tanstack/react-table'
import { Input } from '@/components/ui/input'
import { DataTableFacetedFilter, DataTableColumnHeader, DataTable } from '@/components/data-table'
import { Badge } from '@/components/ui/badge'
import { getHttpStatusVariant } from '@/config/status-variants'
import type { OperationLog } from './types/log'
import { getLogList } from './service/logs'
import { ListPageLayout } from '@/components'
import { useI18n } from '@/context/i18n-provider'
import { useDataTable } from '@/hooks/use-data-table'

export function OperationLog() {
  const { t } = useI18n()

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
    ip: getFilterValue('ip'),
    path: getFilterValue('path'),
    status: getFilterValue('status'),
    pageNum,
    pageSize: pagination.pageSize,
    ...getSortingParams(),
  }), [getFilterValue, getSortingParams, pageNum, pagination.pageSize])

  // 获取日志列表数据
  const { data, isPending: isLoading, error } = useQuery({
    queryKey: ['logList', queryParams],
    queryFn: () => getLogList(queryParams),
    placeholderData: keepPreviousData,
  })

  // 提取数据
  const logData = data?.logs || []
  const total = data?.total || 0
  const pageCount = Math.ceil(total / pagination.pageSize)

  // 格式化时间
  const formatTime = (timeStr: string) => {
    try {
      const date = new Date(timeStr)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
      })
    } catch {
      return timeStr
    }
  }

  // 列定义
  const columns = useMemo<ColumnDef<OperationLog>[]>(
    () => [
      {
        accessorKey: 'username',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.log.operationLog.username')} />
        ),
        cell: ({ row }) => <div className='font-medium'>{row.getValue('username')}</div>,
      },
      {
        accessorKey: 'ip',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.log.operationLog.ip')} />
        ),
        cell: ({ row }) => <div>{row.getValue('ip')}</div>,
      },
      {
        accessorKey: 'path',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.log.operationLog.path')} />
        ),
        cell: ({ row }) => {
          const path = row.getValue('path') as string
          return (
            <div className='max-w-[200px] truncate' title={path}>
              {path}
            </div>
          )
        },
      },
      {
        accessorKey: 'status',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.log.operationLog.status')} />
        ),
        cell: ({ row }) => {
          const status = row.getValue('status') as number
          return (
            <Badge variant={getHttpStatusVariant(status)}>
              {status}
            </Badge>
          )
        },
      },
      {
        accessorKey: 'startTime',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.log.operationLog.startTime')} />
        ),
        cell: ({ row }) => {
          const startTime = row.getValue('startTime') as string
          return <div>{formatTime(startTime)}</div>
        },
      },
      {
        accessorKey: 'timeCost',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.log.operationLog.timeCost')} />
        ),
        cell: ({ row }) => {
          const timeCost = row.getValue('timeCost') as number
          return <div>{timeCost}ms</div>
        },
      },
      {
        accessorKey: 'desc',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.log.operationLog.desc')} />
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
    ],
    [t]
  )

  // 表格实例（服务器端分页，不使用客户端过滤和分页）
  const table = useReactTable({
    data: logData,
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

  return (
    <ListPageLayout
      title={t('features.log.operationLog.title')}
      description={t('features.log.operationLog.description')}
      filterContent={
        <div className='flex flex-wrap gap-2'>
          <Input
            type='text'
            placeholder={t('features.log.operationLog.searchUsername')}
            value={(table.getColumn('username')?.getFilterValue() as string) ?? ''}
            onChange={(e) =>
              table.getColumn('username')?.setFilterValue(e.target.value)
            }
            className='h-8 w-[150px]'
          />
          <Input
            type='text'
            placeholder={t('features.log.operationLog.searchIp')}
            value={(table.getColumn('ip')?.getFilterValue() as string) ?? ''}
            onChange={(e) =>
              table.getColumn('ip')?.setFilterValue(e.target.value)
            }
            className='h-8 w-[150px]'
          />
          <Input
            type='text'
            placeholder={t('features.log.operationLog.searchPath')}
            value={(table.getColumn('path')?.getFilterValue() as string) ?? ''}
            onChange={(e) =>
              table.getColumn('path')?.setFilterValue(e.target.value)
            }
            className='h-8 w-[150px]'
          />
          <DataTableFacetedFilter
            column={table.getColumn('status')}
            title={t('features.log.operationLog.status')}
            options={[
              { label: t('features.log.operationLog.statusOptions.200'), value: '200' },
              { label: t('features.log.operationLog.statusOptions.300'), value: '300' },
              { label: t('features.log.operationLog.statusOptions.400'), value: '400' },
              { label: t('features.log.operationLog.statusOptions.500'), value: '500' },
            ]}
            single
          />
        </div>
      }
    >
      <DataTable
        table={table}
        columns={columns}
        isLoading={isLoading}
        error={error}
        loadingText={t('features.log.operationLog.loading')}
        errorText={t('features.log.operationLog.loadError')}
        emptyText={t('features.log.operationLog.noData')}
        bordered={false}
      />
    </ListPageLayout>
  )
}
// MIGRATED
