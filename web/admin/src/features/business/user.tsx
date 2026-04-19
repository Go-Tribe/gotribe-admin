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
import { DataTableColumnHeader, DataTable } from '@/components/data-table'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { UserFormDialog } from './components/user-form-dialog'
import { UserDetailDialog } from './components/user-detail-dialog'
import type { User } from './types/user'
import { getUserList, createUser, updateUser } from './service/user'
import { getProjectList } from './service/project'
import { useI18n } from '@/context/i18n-provider'
import { useDataTable } from '@/hooks/use-data-table'
import { useCrudMutations } from '@/hooks/use-crud-mutations'
// 新组件导入
import { ListPageLayout, DebouncedInput, DataTableActions, LazyImage } from '@/components'

export function BusinessUser() {
  const { t } = useI18n()
  const [dialogOpen, setDialogOpen] = useState<'create' | 'edit' | null>(null)
  const [editingUser, setEditingUser] = useState<User | null>(null)
  const [detailUserID, setDetailUserID] = useState<number | null>(null)

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

  // 获取项目列表用于筛选
  const { data: projectData } = useQuery({
    queryKey: ['projectList', { current: 1, pageNum: 1, pageSize: 1000 }],
    queryFn: () => getProjectList({ current: 1, pageNum: 1, pageSize: 1000 }),
  })

  const projectList = projectData?.projects || []

  // 构建查询参数
  const queryParams = useMemo(() => ({
    current: 1,
    userID: getFilterValue('id') || undefined,
    projectID: getFilterValue('projectID') || undefined,
    pageNum,
    pageSize: pagination.pageSize,
  }), [getFilterValue, pagination.pageSize, pageNum])

  // 获取用户列表数据
  const { data, isLoading, error } = useQuery({
    queryKey: ['userList', queryParams],
    queryFn: () => getUserList(queryParams),
    placeholderData: keepPreviousData,
  })

  // 提取数据
  const userData = data?.users || []
  const total = data?.total || 0
  const pageCount = Math.ceil(total / pagination.pageSize)

  // 使用统一的 CRUD mutations
  const { createMutation, updateMutation } = useCrudMutations<User, number>({
    queryKey: ['userList'],
    createFn: createUser,
    updateFn: updateUser,
    deleteFn: async () => {}, // 用户管理没有删除功能
    messages: {
      createSuccess: t('features.business.user.createSuccess'),
      updateSuccess: t('features.business.user.updateSuccess'),
    },
    onSuccess: () => {
      setDialogOpen(null)
      setEditingUser(null)
    },
  })

  // 格式化日期时间
  const formatDateTime = useCallback((dateTime: string) => {
    if (!dateTime) return '-'
    try {
      return dateTime.replace(' ', ' ')
    } catch {
      return dateTime
    }
  }, [])

  // 格式化生日
  const formatBirthday = useCallback((birthday: string) => {
    if (!birthday) return '-'
    return birthday
  }, [])

  // 格式化性别
  const formatSex = useCallback((sex: string) => {
    if (sex === 'M') return t('features.business.user.sex.male')
    if (sex === 'F') return t('features.business.user.sex.female')
    return t('features.business.user.sex.unknown')
  }, [t])

  // 处理新建
  const handleCreate = useCallback(() => {
    setEditingUser(null)
    setDialogOpen('create')
  }, [])

  // 处理编辑
  const handleEdit = useCallback((user: User) => {
    setEditingUser(user)
    setDialogOpen('edit')
  }, [])

  // 列定义
  const columns = useMemo<ColumnDef<User>[]>(
    () => [
      {
        accessorKey: 'projectID',
        header: () => null,
        cell: () => null,
        enableHiding: false,
        enableSorting: false,
        enableColumnFilter: true,
      },
      {
        accessorKey: 'avatarURL',
        header: t('features.business.user.columns.avatar'),
        cell: ({ row }) => {
          const avatarURL = row.getValue('avatarURL') as string
          const username = row.original.username
          // 使用 LazyImage 优化头像加载
          return (
            <LazyImage
              src={avatarURL}
              width={40}
              height={40}
              rounded="full"
              objectFit="cover"
              placeholder={
                <Avatar className="h-10 w-10">
                  <AvatarFallback>
                    {username?.charAt(0)?.toUpperCase() || 'U'}
                  </AvatarFallback>
                </Avatar>
              }
            />
          )
        },
      },
      {
        accessorKey: 'id',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.business.user.columns.id')} />
        ),
        cell: ({ row }) => (
          <div className='font-medium'>{row.getValue('id')}</div>
        ),
      },
      {
        accessorKey: 'createdAt',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.business.user.columns.createdAt')} />
        ),
        cell: ({ row }) => {
          const createdAt = row.getValue('createdAt') as string
          return <div>{formatDateTime(createdAt)}</div>
        },
      },
      {
        accessorKey: 'username',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.business.user.columns.username')} />
        ),
        cell: ({ row }) => (
          <div className='font-medium'>{row.getValue('username')}</div>
        ),
      },
      {
        accessorKey: 'nickname',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.business.user.columns.nickname')} />
        ),
        cell: ({ row }) => <div>{row.getValue('nickname')}</div>,
      },
      {
        accessorKey: 'point',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.business.user.columns.point')} />
        ),
        cell: ({ row }) => {
          const point = row.getValue('point') as number
          return <div>{point || 0}</div>
        },
      },
      {
        accessorKey: 'birthday',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.business.user.columns.birthday')} />
        ),
        cell: ({ row }) => {
          const birthday = row.getValue('birthday') as string
          return <div>{formatBirthday(birthday)}</div>
        },
      },
      {
        accessorKey: 'sex',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.business.user.columns.sex')} />
        ),
        cell: ({ row }) => {
          const sex = row.getValue('sex') as string
          return <div>{formatSex(sex)}</div>
        },
      },
      {
        id: 'actions',
        header: t('features.business.user.columns.actions'),
        cell: ({ row }) => (
          <DataTableActions
            onView={() => setDetailUserID(row.original.id)}
            onEdit={() => handleEdit(row.original)}
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
    [t, formatDateTime, formatBirthday, formatSex, handleEdit],
  )

  // 表格实例
  const table = useReactTable({
    data: userData,
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
      columnVisibility: {
        projectID: false,
        ...columnVisibility,
      },
      pagination,
    },
  })

  const isDialogLoading = dialogOpen === 'create' 
    ? createMutation.isPending 
    : updateMutation.isPending

  return (
    <ListPageLayout
      title={t('features.business.user.title')}
      description={t('features.business.user.description')}
      actions={
        <Button onClick={handleCreate}>
          <PlusIcon className='h-4 w-4' />
          {t('features.business.user.createButton')}
        </Button>
      }
      filterContent={
        <div className='flex flex-wrap gap-2'>
          <DebouncedInput
            placeholder={t('features.business.user.search.userID')}
            value={(table.getColumn('id')?.getFilterValue() as string) ?? ''}
            onChange={(value) => table.getColumn('id')?.setFilterValue(value)}
            delay={300}
            className='h-8 w-[150px]'
          />
          <Select
            value={
              (table.getColumn('projectID')?.getFilterValue() as string) ||
              'all'
            }
            onValueChange={(value) =>
              table
                .getColumn('projectID')
                ?.setFilterValue(value === 'all' ? undefined : value)
            }
          >
            <SelectTrigger className='h-8 w-[180px]'>
              <SelectValue placeholder={t('features.business.user.search.project')} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value='all'>{t('features.business.user.search.allProjects')}</SelectItem>
              {projectList.map((project) => (
                <SelectItem key={project.projectID} value={project.projectID}>
                  {project.title}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      }
      dialogs={
        <>
          <UserFormDialog
            open={dialogOpen !== null}
            onOpenChange={(open) => {
              if (!open) {
                setDialogOpen(null)
                setEditingUser(null)
              }
            }}
            user={editingUser}
            onSubmit={(data) => {
              if (dialogOpen === 'create') {
                createMutation.mutate(data)
              } else {
                updateMutation.mutate(data)
              }
            }}
            isLoading={isDialogLoading}
          />
          <UserDetailDialog
            open={!!detailUserID}
            onOpenChange={(open) => !open && setDetailUserID(null)}
            userID={detailUserID}
          />
        </>
      }
    >
      <DataTable
        table={table}
        columns={columns}
        isLoading={isLoading}
        error={error}
        loadingText={t('features.business.user.loading')}
        errorText={t('features.business.user.loadError')}
        emptyText={t('features.business.user.noData')}
      />
    </ListPageLayout>
  )
}
