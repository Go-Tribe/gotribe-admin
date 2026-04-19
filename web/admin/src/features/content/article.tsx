import { useMemo, useState } from 'react'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import {
  useReactTable,
  getCoreRowModel,
  type ColumnDef,
} from '@tanstack/react-table'
import { PlusIcon, Pencil1Icon, TrashIcon } from '@radix-ui/react-icons'
import { Button } from '@/components/ui/button'
import { SubTitle } from '@/components/sub-title'
import { Input } from '@/components/ui/input'
import { DataTableColumnHeader, DataTable } from '@/components/data-table'
import { DataTableFacetedFilter } from '@/components/data-table/faceted-filter'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { getStatusVariantByType } from '@/config/status-variants'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { useNavigate } from '@tanstack/react-router'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { useI18n } from '@/context/i18n-provider'
import type { Post } from './types/post'
import { getPostList, deletePost } from './service/post'
import { getProjectList } from '@/features/business/service/project'
import { useDataTable } from '@/hooks/use-data-table'
import { useCrudMutations } from '@/hooks/use-crud-mutations'

export function ContentArticle() {
  'use no memo'
  const { t } = useI18n()
  const navigate = useNavigate()
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<string | null>(null)

  const {
    columnFilters,
    setColumnFilters,
    pagination,
    setPagination,
    sorting,
    setSorting,
    columnVisibility,
    setColumnVisibility,
    getFilterValue,
    pageNum,
  } = useDataTable()

  const queryParams = useMemo(() => {
    const sortBy = sorting[0]?.id
    const sortOrder = sorting[0]?.desc
      ? ('desc' as const)
      : sorting[0]
        ? ('asc' as const)
        : undefined

    return {
      postID: getFilterValue('postID'),
      title: getFilterValue('title'),
      status: getFilterValue('status') || undefined,
      projectID: getFilterValue('projectID') || undefined,
      pageNum,
      pageSize: pagination.pageSize,
      sortBy,
      sortOrder,
    }
  }, [getFilterValue, pagination, sorting, pageNum])

  const { data: projectData } = useQuery({
    queryKey: ['projectList', { current: 1, pageNum: 1, pageSize: 1000 }],
    queryFn: () => getProjectList({ current: 1, pageNum: 1, pageSize: 1000 }),
  })
  const projectList = useMemo(() => projectData?.projects ?? [], [projectData?.projects])

  const { data, isLoading, error } = useQuery({
    queryKey: ['postList', queryParams],
    queryFn: () => getPostList(queryParams),
    placeholderData: keepPreviousData,
  })

  const postData = data?.posts ?? []
  const total = data?.total ?? 0
  const pageCount = Math.ceil(total / pagination.pageSize)

  const columns = useMemo<ColumnDef<Post>[]>(
    () => [
      {
        accessorKey: 'postID',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.article.columns.postID')} />
        ),
        cell: ({ row }) => (
          <div className='font-mono text-muted-foreground'>{row.getValue('postID') as string}</div>
        ),
      },
      {
        accessorKey: 'title',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.article.columns.title')} />
        ),
        cell: ({ row }) => (
          <div className='flex items-center gap-2'>
            {row.original.isTop === 2 && (
              <Badge variant='destructive' className='px-1 py-0 text-[10px] h-5 whitespace-nowrap'>
                TOP
              </Badge>
            )}
            <div className='max-w-[200px] truncate font-medium' title={row.getValue('title') as string}>
              {row.getValue('title') as string}
            </div>
          </div>
        ),
      },
      {
        accessorKey: 'author',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.article.columns.author')} />
        ),
        cell: ({ row }) => <div>{row.getValue('author') as string}</div>,
      },
      {
        accessorKey: 'description',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.article.columns.description')} />
        ),
        cell: ({ row }) => {
          const desc = row.getValue('description') as string
          return (
            <div className='max-w-[280px] truncate text-muted-foreground' title={desc}>
              {desc || '-'}
            </div>
          )
        },
      },
      {
        id: 'projectID',
        accessorKey: 'projectID',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.article.filter.project')} />
        ),
        cell: ({ row }) => {
          const id = row.getValue('projectID') as string
          const project = projectList.find((p) => p.projectID === id)
          return <div className='text-muted-foreground'>{project?.title ?? id ?? '-'}</div>
        },
      },
      {
        accessorKey: 'status',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.article.columns.status')} />
        ),
        cell: ({ row }) => {
          const status = row.getValue('status') as number
          return (
            <Badge variant={getStatusVariantByType(status, 'articlePublish')}>
              {status === 2 ? t('features.content.article.status.published') : t('features.content.article.status.draft')}
            </Badge>
          )
        },
      },
      {
        accessorKey: 'createdAt',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.content.article.columns.createdAt')} />
        ),
        cell: ({ row }) => {
          const time = row.original.showTime || row.original.createdAt
          return (
            <div className='text-muted-foreground'>
              {time ? String(time).slice(0, 19) : '-'}
            </div>
          )
        },
      },
      {
        id: 'actions',
        header: t('features.content.article.columns.actions'),
        cell: ({ row }) => {
          const post = row.original
          return (
            <div className='flex items-center gap-2'>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant='outline'
                    size='sm'
                    className='h-8 border-border/60'
                    onClick={() => navigate({ to: '/content/article/$postID/edit', params: { postID: post.postID } })}
                  >
                    <Pencil1Icon className='h-4 w-4' />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>{t('features.content.article.actions.edit')}</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant='ghost'
                    size='sm'
                    className='h-8 text-destructive hover:text-destructive'
                    onClick={() => setDeleteDialogOpen(post.postID)}
                  >
                    <TrashIcon className='h-4 w-4 text-destructive' />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>{t('features.content.article.actions.delete')}</TooltipContent>
              </Tooltip>
            </div>
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
    [t, navigate, projectList]
  )

  // eslint-disable-next-line react-hooks/incompatible-library
  const table = useReactTable({
    data: postData,
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

  const { deleteMutation, isLoading: isMutating } = useCrudMutations<Post, string>({
    queryKey: ['postList'],
    createFn: async () => { throw new Error('Not implemented') },
    updateFn: async () => { throw new Error('Not implemented') },
    deleteFn: deletePost,
    messages: {
      deleteSuccess: t('features.content.article.deleteSuccess'),
    },
    onSuccess: () => {
      setDeleteDialogOpen(null)
    },
  })

  const handleDelete = () => {
    if (deleteDialogOpen) {
      deleteMutation.mutate(deleteDialogOpen)
    }
  }

  const postToDelete = postData.find((p) => p.postID === deleteDialogOpen)

  return (
    <div className='space-y-4'>
      <div className='flex items-center justify-between px-4 pt-4'>
        <SubTitle
          title={t('features.content.article.title')}
          description={t('features.content.article.description')}
          children={
            <Button onClick={() => navigate({ to: '/content/article/new' })}>
              <PlusIcon className='h-4 w-4' />
              {t('features.content.article.createButton')}
            </Button>
          }
        />
      </div>

      <DataTable
        table={table}
        columns={columns}
        isLoading={isLoading}
        error={error}
        loadingText={t('features.content.article.loading')}
        errorText={t('features.content.article.loadError')}
        emptyText={t('features.content.article.noData')}
      >
        <div className='flex flex-wrap items-center gap-2 pb-4'>
          <Input
            type='text'
            placeholder={t('features.content.article.search.postID')}
            value={(table.getColumn('postID')?.getFilterValue() as string) ?? ''}
            onChange={(e) => table.getColumn('postID')?.setFilterValue(e.target.value)}
            className='h-8 w-[160px]'
          />
          <Input
            type='text'
            placeholder={t('features.content.article.search.title')}
            value={(table.getColumn('title')?.getFilterValue() as string) ?? ''}
            onChange={(e) => table.getColumn('title')?.setFilterValue(e.target.value)}
            className='h-8 w-[180px]'
          />
          <Select
            value={
              (table.getColumn('projectID')?.getFilterValue() as string) || 'all'
            }
            onValueChange={(value) =>
              table
                .getColumn('projectID')
                ?.setFilterValue(value === 'all' ? undefined : value)
            }
          >
            <SelectTrigger className='h-8 w-[180px]'>
              <SelectValue placeholder={t('features.content.article.filter.project')} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value='all'>{t('features.content.article.search.allProjects')}</SelectItem>
              {projectList.map((project) => (
                <SelectItem key={project.projectID} value={project.projectID}>
                  {project.title}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <DataTableFacetedFilter
            column={table.getColumn('status')}
            title={t('features.content.article.filter.status')}
            options={[
              { label: t('features.content.article.status.draft'), value: '1' },
              { label: t('features.content.article.status.published'), value: '2' },
            ]}
            single
          />
        </div>
      </DataTable>

      <ConfirmDialog
        open={deleteDialogOpen !== null}
        onOpenChange={(open) => {
          if (!open) setDeleteDialogOpen(null)
        }}
        title={t('features.content.article.confirmDelete')}
        desc={t('features.content.article.confirmDeleteDesc', { title: postToDelete?.title ?? '' })}
        handleConfirm={handleDelete}
        destructive
        confirmText={t('features.content.article.actions.delete')}
        isLoading={isMutating}
      />
    </div>
  )
}
