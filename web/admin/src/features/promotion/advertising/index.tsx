import { useMemo, useState, useCallback } from 'react'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import {
  useReactTable,
  getCoreRowModel,
  type ColumnDef,
} from '@tanstack/react-table'
import { PlusIcon } from '@radix-ui/react-icons'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { getStatusVariantByType } from '@/config/status-variants'
import { DataTableColumnHeader, DataTable, DataTableFacetedFilter, DataTableActions } from '@/components/data-table'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { ListPageLayout } from '@/components/layout'
import { getSceneList } from '@/shared/api'
import type { Scene } from '@/shared/types'
import { useI18n } from '@/context/i18n-provider'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { AdFormDialog } from './components/ad-form-dialog'
import { createAd, getAdList, updateAd, deleteAd } from './service'
import type { AdCreateParams, AdUpdateParams, Ad } from './types'
import { useDataTable } from '@/hooks/use-data-table'
import { useCrudMutations } from '@/hooks/use-crud-mutations'

/** 表格行（由 API 数据映射） */
type AdRow = {
  id: string
  title: string
  description: string
  status: number | string
  createdAt: string
}

function mapAdToRow(ad: Ad): AdRow {
  const id = String(ad.id ?? '')
  return {
    id,
    title: (ad.title ?? '').trim() || '-',
    description: (ad.description ?? '').trim() || '-',
    status: ad.status ?? 1,
    createdAt: (ad.createdAt as string)?.trim() || '-',
  }
}

export function PromotionAdvertising() {
  'use no memo'
  const { t } = useI18n()
  const [dialogOpen, setDialogOpen] = useState<'create' | 'edit' | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<string | null>(null)
  const [editAdId, setEditAdId] = useState<string | null>(null)
  const [sceneFilter, setSceneFilter] = useState<string>('__all__')

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

  // 获取场景列表用于筛选
  const { data: sceneListData } = useQuery({
    queryKey: ['adSceneList', { pageNum: 1, pageSize: 1000 }],
    queryFn: () => getSceneList({ pageNum: 1, pageSize: 1000 }),
  })
  const sceneList = useMemo(
    () =>
      sceneListData?.adScenes?.map((s: Scene) => ({
        adSceneID: String(s.id),
        title: s.title ?? String(s.id),
      })) ?? [],
    [sceneListData?.adScenes]
  )

  // 构建查询参数
  const queryParams = useMemo(() => ({
    pageNum,
    pageSize: pagination.pageSize,
    sceneID:
      sceneFilter && sceneFilter !== '__all__'
        ? sceneFilter.trim() || undefined
        : undefined,
    title: getFilterValue('title') || undefined,
    status: getFilterValue('status') || undefined,
    ...getSortingParams(),
  }), [pageNum, pagination.pageSize, sceneFilter, getFilterValue, getSortingParams])

  // 获取广告列表数据
  const { data: adListData, isPending: isLoading, error } = useQuery({
    queryKey: ['adList', queryParams],
    queryFn: () => getAdList(queryParams),
    placeholderData: keepPreviousData,
  })

  const rawList = adListData?.ads ?? []
  const adData = useMemo(
    () => rawList.map((ad) => mapAdToRow(ad)),
    [rawList]
  )
  const total = adListData?.total ?? 0
  const pageCount = Math.max(1, Math.ceil(total / pagination.pageSize))
  const editAd = useMemo(() => {
    if (editAdId == null) return null
    return rawList.find((ad) => String(ad.id ?? '') === editAdId) ?? null
  }, [editAdId, rawList])

  // 使用统一的 CRUD mutations
  const { createMutation, updateMutation, deleteMutation } = useCrudMutations<Ad, string>({
    queryKey: ['adList'],
    createFn: async (data: Partial<Ad>) => {
      const payload = data as AdCreateParams
      await createAd({
        ...payload,
        sceneID: Number(payload.sceneID),
      })
    },
    updateFn: async (data: Partial<Ad>) => {
      if (data.id != null) {
        const payload = data as AdUpdateParams
        await updateAd(String(data.id), {
          ...payload,
          sceneID: payload.sceneID != null ? Number(payload.sceneID) : undefined,
        })
      }
    },
    deleteFn: async (id: string) => {
      await deleteAd(id)
    },
    messages: {
      createSuccess: t('features.operation.advertising.createSuccess'),
      updateSuccess: t('features.operation.advertising.updateSuccess'),
      deleteSuccess: t('features.operation.advertising.deleteSuccess'),
    },
    onSuccess: () => {
      setDialogOpen(null)
      setEditAdId(null)
      setDeleteDialogOpen(null)
    },
  })

  // 处理新建
  const handleCreate = useCallback(() => {
    setEditAdId(null)
    setDialogOpen('create')
  }, [])

  // 处理编辑
  const handleEdit = useCallback((ad: Ad) => {
    setEditAdId(ad.id != null ? String(ad.id) : null)
    setDialogOpen('edit')
  }, [])

  // 处理删除
  const handleDelete = useCallback(() => {
    if (deleteDialogOpen) {
      deleteMutation.mutate(deleteDialogOpen)
    }
  }, [deleteDialogOpen, deleteMutation])

  const adToDelete = adData.find((r) => r.id === deleteDialogOpen)

  const columns = useMemo<ColumnDef<AdRow>[]>(
    () => [
      {
        accessorKey: 'id',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.advertising.advertising.columns.id')}
          />
        ),
        cell: ({ row }) => (
          <div className='font-mono text-muted-foreground'>
            {(row.getValue('id') as string)?.trim() ?? '-'}
          </div>
        ),
      },
      {
        accessorKey: 'title',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.advertising.advertising.columns.title')}
          />
        ),
        cell: ({ row }) => (
          <div
            className='max-w-[160px] truncate font-medium'
            title={row.getValue('title') as string}
          >
            {(row.getValue('title') as string) ?? '-'}
          </div>
        ),
      },
      {
        accessorKey: 'description',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.advertising.advertising.columns.description')}
          />
        ),
        cell: ({ row }) => {
          const desc = row.getValue('description') as string
          return (
            <div
              className='max-w-[280px] truncate text-muted-foreground'
              title={desc}
            >
              {desc ?? '-'}
            </div>
          )
        },
      },
      {
        accessorKey: 'status',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('features.operation.advertising.advertising.columns.status')}
          />
        ),
        cell: ({ row }) => {
          const status = row.getValue('status') as number
          return (
            <Badge variant={getStatusVariantByType(status, 'publishStatus')}>
                {status === 2
                  ? t('features.operation.advertising.advertising.status.published')
                  : t('features.operation.advertising.advertising.status.unpublished')}
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
              'features.operation.advertising.advertising.columns.createdAt'
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
        header: t('features.operation.advertising.advertising.columns.actions'),
        cell: ({ row }) => {
          const adRow = row.original
          const ad = rawList.find((a) => String(a.id ?? '') === adRow.id)
          return (
            <DataTableActions
              onEdit={() => ad && handleEdit(ad)}
              onDelete={() => setDeleteDialogOpen(adRow.id)}
              deleteConfirmTitle={t('features.operation.advertising.confirmDelete')}
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
    data: adData,
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

  const isDialogLoading = dialogOpen === 'create'
    ? createMutation.isPending
    : updateMutation.isPending

  return (
    <ListPageLayout
      title={t('features.operation.advertising.title')}
      description={t('features.operation.advertising.description')}
      actions={
        <Button onClick={handleCreate}>
          <PlusIcon className='h-4 w-4' />
          {t('features.operation.advertising.advertising.createButton')}
        </Button>
      }
      filterContent={
        <div className='flex flex-wrap gap-2'>
          <Select value={sceneFilter} onValueChange={setSceneFilter}>
            <SelectTrigger className='h-8 w-[160px]'>
              <SelectValue
                placeholder={t(
                  'features.operation.advertising.advertising.search.scene'
                )}
              />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value='__all__'>
                {t(
                  'features.operation.advertising.advertising.search.sceneAll'
                )}
              </SelectItem>
              {sceneList.map((s: { adSceneID: string; title: string }) => (
                <SelectItem
                  key={s.adSceneID}
                  value={(s.adSceneID ?? '').trim()}
                >
                  {s.title}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Input
            type='text'
            placeholder={t(
              'features.operation.advertising.advertising.search.titlePlaceholder'
            )}
            value={(table.getColumn('title')?.getFilterValue() as string) ?? ''}
            onChange={(e) => table.getColumn('title')?.setFilterValue(e.target.value)}
            className='h-8 w-[200px]'
          />
          <DataTableFacetedFilter
            column={table.getColumn('status')}
            title={t('features.operation.advertising.advertising.search.status')}
            options={[
              { label: t('features.operation.advertising.advertising.status.unpublished'), value: '1' },
              { label: t('features.operation.advertising.advertising.status.published'), value: '2' },
            ]}
            single
          />
        </div>
      }
      dialogs={
        <>
          {/* 新建/编辑对话框 */}
          <AdFormDialog
        open={dialogOpen !== null}
        onOpenChange={(open) => {
          if (!open) {
            setDialogOpen(null)
            setEditAdId(null)
          }
        }}
        onSubmit={(payload) => createMutation.mutate(payload as Partial<Ad>)}
        onSubmitUpdate={(adID, payload) =>
          updateMutation.mutate({ id: Number(adID), ...payload } as Partial<Ad>)
        }
        isLoading={isDialogLoading}
        sceneList={sceneList}
        editAd={editAd}
      />

      {/* 删除确认对话框 */}
      <ConfirmDialog
        open={deleteDialogOpen !== null}
        onOpenChange={(open) => {
          if (!open) setDeleteDialogOpen(null)
        }}
        title={t('features.operation.advertising.confirmDelete')}
        desc={t('features.operation.advertising.confirmDeleteDesc', {
          name: adToDelete?.title ?? deleteDialogOpen ?? '',
        })}
        handleConfirm={handleDelete}
        destructive
        confirmText={t('features.operation.advertising.advertising.actions.delete')}
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
        loadingText={t('features.operation.advertising.advertising.loading')}
        errorText={t('features.operation.advertising.advertising.loadError')}
        emptyText={t('features.operation.advertising.advertising.noData')}
      />
    </ListPageLayout>
  )
}
