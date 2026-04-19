import { useMemo, useState } from 'react'
import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query'
import {
  useReactTable,
  getCoreRowModel,
  type ColumnDef,
} from '@tanstack/react-table'
import { CheckCircledIcon, CrossCircledIcon } from '@radix-ui/react-icons'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { getStatusVariantByType } from '@/config/status-variants'
import { DataTableColumnHeader, DataTable, DataTableFacetedFilter } from '@/components/data-table'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { ListPageLayout } from '@/components/layout'
import { useI18n } from '@/context/i18n-provider'
import { toast } from 'sonner'
import { getProjectList } from '@/shared/api'
import { getCommentList, approveComment, rejectComment } from './service/comment'
import { useDataTable } from '@/hooks/use-data-table'
import type { Comment } from './types/comment'

/** 表格行（由 API 数据映射） */
type CommentRow = {
  id: string
  comment: string
  nickname: string
  userID: string
  objectType: number
  status: number
  createdAt: string
}

function mapCommentToRow(c: Comment): CommentRow {
  const id = (c.commentID ?? '').trim() || String(c.id ?? '')
  return {
    id,
    comment: (c.comment ?? '').trim() || '-',
    nickname: (c.nickname ?? '').trim() || '-',
    userID: c.userID != null ? String(c.userID) : '-',
    objectType: c.objectType ?? 0,
    status: c.status ?? 1,
    createdAt: (c.createdAt as string)?.trim() || '-',
  }
}

export function OperationComment() {
  'use no memo'
  const { t } = useI18n()
  const queryClient = useQueryClient()
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
    getSortingParams,
  } = useDataTable({
    defaultPageSize: 10,
    defaultPageIndex: 0,
  })

  const { data: projectData } = useQuery({
    queryKey: ['projectList', { pageNum: 1, pageSize: 1000 }],
    queryFn: () => getProjectList({ pageNum: 1, pageSize: 1000 }),
  })
  const projectList = useMemo(
    () =>
      projectData?.projects?.map((p) => ({
        id: p.id,
        title: p.title ?? String(p.id),
      })) ?? [],
    [projectData?.projects]
  )

  const queryParams = useMemo(() => ({
    pageNum,
    pageSize: pagination.pageSize,
    status: getFilterValue('status') || undefined,
    nickname: getFilterValue('nickname') || undefined,
    projectId:
      projectFilter && projectFilter !== '__all__'
        ? Number(projectFilter) || undefined
        : undefined,
    ...getSortingParams(),
  }), [pageNum, pagination.pageSize, getFilterValue, getSortingParams, projectFilter])

  const { data: commentListData, isPending: isLoading, error } = useQuery({
    queryKey: ['commentList', queryParams],
    queryFn: () => getCommentList(queryParams),
    placeholderData: keepPreviousData,
  })

  const rawList = commentListData?.comments ?? []
  const commentData = useMemo(
    () => rawList.map((c) => mapCommentToRow(c)),
    [rawList]
  )
  const total = commentListData?.total ?? 0
  const pageCount = Math.max(1, Math.ceil(total / pagination.pageSize))
  // 使用独立的 mutations (approve/reject 是特殊操作，不适合 useCrudMutations)
  const { mutate: approveMutate, isPending: isApprovePending } = useMutation({
    mutationFn: (commentID: string) => approveComment(commentID),
    onSuccess: () => {
      toast.success(t('features.operation.comment.approveSuccess'))
      queryClient.invalidateQueries({ queryKey: ['commentList'] })
    },
    onError: (err) => {
      toast.error(
        err instanceof Error ? err.message : t('features.operation.comment.approveError')
      )
    },
  })

  const { mutate: rejectMutate, isPending: isRejectPending } = useMutation({
    mutationFn: (commentID: string) => rejectComment(commentID),
    onSuccess: () => {
      toast.success(t('features.operation.comment.rejectSuccess'))
      queryClient.invalidateQueries({ queryKey: ['commentList'] })
    },
    onError: (err) => {
      toast.error(
        err instanceof Error ? err.message : t('features.operation.comment.rejectError')
      )
    },
  })

  const columns = useMemo<ColumnDef<CommentRow>[]>(
    () => [
      {
        accessorKey: 'id',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.comment.comment.columns.id')}
          />
        ),
        cell: ({ row }) => (
          <div className='font-mono text-muted-foreground max-w-[100px] truncate'>
            {(row.getValue('id') as string) ?? '-'}
          </div>
        ),
      },
      {
        accessorKey: 'comment',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.comment.comment.columns.comment')}
          />
        ),
        cell: ({ row }) => {
          const comment = row.getValue('comment') as string
          return (
            <div
              className='max-w-[200px] truncate font-medium'
              title={comment}
            >
              {comment ?? '-'}
            </div>
          )
        },
      },
      {
        accessorKey: 'nickname',
        enableSorting: false,
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.comment.comment.columns.nickname')}
          />
        ),
        cell: ({ row }) => (
          <div className='max-w-[120px] truncate text-muted-foreground'>
            {(row.getValue('nickname') as string) ?? '-'}
          </div>
        ),
      },
      {
        accessorKey: 'userID',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.comment.comment.columns.userID')}
          />
        ),
        cell: ({ row }) => (
          <div className='font-mono text-muted-foreground max-w-[80px] truncate'>
            {(row.getValue('userID') as string) ?? '-'}
          </div>
        ),
      },
      {
        accessorKey: 'status',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.comment.comment.columns.status')}
          />
        ),
        cell: ({ row }) => {
          const status = row.getValue('status') as number
          return (
            <Badge variant={getStatusVariantByType(status, 'approvalStatus')}>
                {status === 2
                  ? t('features.operation.comment.comment.status.approved')
                  : t('features.operation.comment.comment.status.pending')}
            </Badge>
          )
        },
      },
      {
        accessorKey: 'createdAt',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t(
              'features.operation.comment.comment.columns.createdAt'
            )}
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
        header: t('features.operation.comment.comment.columns.actions'),
        cell: ({ row }) => {
          const commentRow = row.original
          const isPending = commentRow.status === 1
          const isApproved = commentRow.status === 2
          return (
            <div className='flex items-center gap-2'>
              {isPending && (
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      variant='outline'
                      size='sm'
                      className='h-8 border-border/60 text-green-700 hover:text-green-700'
                      onClick={() => approveMutate(commentRow.id)}
                      disabled={isApprovePending}
                    >
                      <CheckCircledIcon className='h-4 w-4 text-green-600' />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>
                    {t('features.operation.comment.comment.actions.approve')}
                  </TooltipContent>
                </Tooltip>
              )}
              {isApproved && (
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      variant='outline'
                      size='sm'
                      className='h-8 border-border/60 text-red-700 hover:text-red-700'
                      onClick={() => rejectMutate(commentRow.id)}
                      disabled={isRejectPending}
                    >
                      <CrossCircledIcon className='h-4 w-4 text-red-600' />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>
                    {t('features.operation.comment.comment.actions.reject')}
                  </TooltipContent>
                </Tooltip>
              )}
            </div>
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
    [t, approveMutate, isApprovePending, rejectMutate, isRejectPending]
  )

  const table = useReactTable({
    data: commentData,
    columns,
    getCoreRowModel: getCoreRowModel(),
    manualSorting: true,
    manualPagination: true,
    pageCount,
    getRowId: (row) => row.id,
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
      title={t('features.operation.comment.title')}
      description={t('features.operation.comment.description')}
      filterContent={
        <div className='flex flex-wrap gap-2'>
          <Select value={projectFilter} onValueChange={setProjectFilter}>
            <SelectTrigger className='h-8 w-[160px]'>
              <SelectValue
                placeholder={t(
                  'features.operation.comment.comment.search.project'
                )}
              />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value='__all__'>
                {t(
                  'features.operation.comment.comment.search.projectAll'
                )}
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
          <Input
            type='text'
            placeholder={t(
              'features.operation.comment.comment.search.nicknamePlaceholder'
            )}
            value={(table.getColumn('nickname')?.getFilterValue() as string) ?? ''}
            onChange={(e) => table.getColumn('nickname')?.setFilterValue(e.target.value)}
            className='h-8 w-[160px]'
          />
          <DataTableFacetedFilter
            column={table.getColumn('status')}
            title={t('features.operation.comment.comment.search.status')}
            options={[
              { label: t('features.operation.comment.comment.status.pending'), value: '1' },
              { label: t('features.operation.comment.comment.status.approved'), value: '2' },
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
        loadingText={t('features.operation.comment.comment.loading')}
        errorText={t('features.operation.comment.comment.loadError')}
        emptyText={t('features.operation.comment.comment.noData')}
      />
    </ListPageLayout>
  )
}
