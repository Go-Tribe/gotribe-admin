import { useState, useCallback, useMemo, useRef, useEffect } from 'react'
import { useQuery, useMutation, keepPreviousData, useQueryClient } from '@tanstack/react-query'
import {
  useReactTable,
  getCoreRowModel,
  type ColumnDef,

} from '@tanstack/react-table'
import {
  Upload,
  Video,
  Loader2,
  ImageIcon,
  FileIcon,
  Pencil,
  Trash2,
  Copy,
  X,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { DataTablePagination } from '@/components/data-table'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { ListPageLayout } from '@/components/layout'
import { useI18n } from '@/context/i18n-provider'
import { getResourceList, uploadResource, deleteResource } from './service/resource'
import { ResourceEditDialog } from './component/resource-edit-dialog'
import { FILE_TYPE, type ResourceItem } from './types/resource'
import { toast } from 'sonner'
import { useDataTable } from '@/hooks/use-data-table'
import { useCrudMutations } from '@/hooks/use-crud-mutations'

const getResourceTypeOptions = (t: (key: string) => string) => [
  { value: '0', label: t('features.content.resource.filter.all') },
  { value: String(FILE_TYPE.IMAGE), label: t('features.content.resource.filter.image') },
  { value: String(FILE_TYPE.VIDEO), label: t('features.content.resource.filter.video') },
  { value: String(FILE_TYPE.AUDIO), label: t('features.content.resource.filter.audio') },
  { value: String(FILE_TYPE.ARCHIVE), label: t('features.content.resource.filter.archive') },
  { value: String(FILE_TYPE.DOCUMENT), label: t('features.content.resource.filter.document') },
  { value: String(FILE_TYPE.FONT), label: t('features.content.resource.filter.font') },
  { value: String(FILE_TYPE.APP), label: t('features.content.resource.filter.app') },
  { value: String(FILE_TYPE.UNKNOWN), label: t('features.content.resource.filter.unknown') },
]

function formatName(name: string, maxLen = 14): string {
  if (name.length <= maxLen) return name
  return `${name.slice(0, 8)}...`
}

export function ContentResource() {
  const { t } = useI18n()
  const queryClient = useQueryClient()
  const [resourceType, setResourceType] = useState<string>('0')
  const resourceTypeOptions = useMemo(() => getResourceTypeOptions(t), [t])
  const [deleteId, setDeleteId] = useState<number | null>(null)
  const [editId, setEditId] = useState<number | null>(null)
  const [previewResource, setPreviewResource] = useState<ResourceItem | null>(null)
  const [uploadProgress, setUploadProgress] = useState(0)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const videoInputRef = useRef<HTMLInputElement>(null)

  // 使用统一的表格状态管理
  const {
    pagination,
    setPagination,
  } = useDataTable()

  const queryParams = useMemo(
    () => ({
      type: resourceType === '0' ? undefined : Number(resourceType),
      pageNum: pagination.pageIndex + 1,
      pageSize: pagination.pageSize,
    }),
    [resourceType, pagination.pageIndex, pagination.pageSize]
  )

  const { data, isPending: isLoading, error } = useQuery({
    queryKey: ['resourceList', queryParams],
    queryFn: () => getResourceList(queryParams),
    placeholderData: keepPreviousData,
  })

  const list = useMemo(() => data?.list ?? [], [data?.list])
  const total = data?.total ?? 0
  const pageCount = Math.max(1, Math.ceil(total / pagination.pageSize))

  const table = useReactTable({
    data: list,
    columns: useMemo<ColumnDef<ResourceItem>[]>(() => [{ id: '_placeholder', header: '', cell: () => null }], []),
    getCoreRowModel: getCoreRowModel(),
    manualPagination: true,
    pageCount,
    state: { pagination },
    onPaginationChange: (updater) => {
      const next = typeof updater === 'function' ? updater(pagination) : pagination
      setPagination(
        next.pageSize !== pagination.pageSize ? { ...next, pageIndex: 0 } : next
      )
    },
  })

  // 上传 mutation
  const { mutate: uploadMutate, isPending: isUploadPending } = useMutation({
    mutationFn: ({
      file,
      type,
      onProgress,
    }: {
      file: File
      type: number
      onProgress?: (progress: number) => void
    }) => uploadResource(file, type, onProgress),
    onSuccess: () => {
      toast.success(t('features.content.resource.uploadSuccess'))
      void queryClient.invalidateQueries({ queryKey: ['resourceList'] })
      setUploadProgress(0)
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : t('features.content.resource.uploadError'))
      setUploadProgress(0)
    },
  })

  // 使用统一的 CRUD mutations 处理删除
  const { deleteMutation } = useCrudMutations<ResourceItem, number>({
    queryKey: ['resourceList'],
    createFn: async () => {},
    updateFn: async () => {},
    deleteFn: (id: number) => deleteResource(id),
    messages: {
      deleteSuccess: t('features.content.resource.deleteSuccess'),
    },
    onSuccess: () => {
      setDeleteId(null)
    },
  })

  const handleCopyUrl = useCallback((url: string) => {
    void navigator.clipboard.writeText(url).then(() => toast.success(t('features.content.resource.copySuccess')))
  }, [t])

  const handleEdit = useCallback((resource: ResourceItem) => {
    setEditId(resource.id)
  }, [])

  const handlePreview = useCallback((resource: ResourceItem) => {
    setPreviewResource(resource)
  }, [])

  const handleDelete = useCallback((resource: ResourceItem) => {
    setDeleteId(resource.id)
  }, [])

  const handleUploadResource = useCallback(() => {
    fileInputRef.current?.click()
  }, [])

  const handleUploadVideo = useCallback(() => {
    videoInputRef.current?.click()
  }, [])

  const onFileSelect = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>, uploadType: number) => {
      const file = e.target.files?.[0]
      if (!file) return
      setUploadProgress(0)
      uploadMutate({ file, type: uploadType, onProgress: setUploadProgress })
      e.target.value = ''
    },
    [uploadMutate]
  )

  const isUploading = isUploadPending

  return (
    <ListPageLayout
      title={t('features.content.resource.title')}
      description={t('features.content.resource.description')}
      filterContent={
        <div className='flex flex-wrap items-center gap-2'>
          <label className='text-sm text-muted-foreground whitespace-nowrap'>{t('features.content.resource.filter.typeLabel')}</label>
          <Select
            value={resourceType}
            onValueChange={(v) => {
              setResourceType(v)
              setPagination((p) => ({ ...p, pageIndex: 0 }))
            }}
          >
            <SelectTrigger className='w-[140px] h-9'>
              <SelectValue placeholder={t('features.content.resource.filter.placeholder')} />
            </SelectTrigger>
            <SelectContent>
              {resourceTypeOptions.map((opt) => (
                <SelectItem key={opt.value} value={opt.value}>
                  {opt.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <div className='flex-1' />
          {isUploading && (
            <div className='flex items-center gap-2 w-[180px] shrink-0'>
              <div className='h-2 flex-1 min-w-0 rounded-full bg-muted'>
                <div
                  className='h-full rounded-full bg-primary transition-all duration-300'
                  style={{ width: `${uploadProgress}%` }}
                />
              </div>
              <span className='text-xs text-muted-foreground shrink-0 tabular-nums'>{uploadProgress}%</span>
            </div>
          )}
          <input
            ref={fileInputRef}
            type='file'
            className='hidden'
            accept='*/*'
            onChange={(e) => onFileSelect(e, FILE_TYPE.IMAGE)}
          />
          <input
            ref={videoInputRef}
            type='file'
            className='hidden'
            accept='video/*'
            onChange={(e) => onFileSelect(e, FILE_TYPE.VIDEO)}
          />
          <Button
            variant='default'
            size='sm'
            className='h-9'
            onClick={handleUploadResource}
            disabled={isUploading}
          >
            {isUploading ? (
              <Loader2 className='h-4 w-4 animate-spin' />
            ) : (
              <Upload className='h-4 w-4' />
            )}
            {t('features.content.resource.actions.uploadResource')}
          </Button>
          <Button
            variant='default'
            size='sm'
            className='h-9'
            onClick={handleUploadVideo}
            disabled={isUploading}
          >
            <Video className='h-4 w-4' />
            {t('features.content.resource.actions.uploadVideo')}
          </Button>
        </div>
      }
      dialogs={
        <>
          <ResourceEditDialog
        open={editId !== null}
        onOpenChange={(open) => !open && setEditId(null)}
        resourceID={editId}
        onSuccess={() => setEditId(null)}
      />

      {previewResource && (
        <ResourcePreviewOverlay
          resource={previewResource}
          onClose={() => setPreviewResource(null)}
          onCopyUrl={handleCopyUrl}
          t={t}
        />
      )}

      <ConfirmDialog
        open={deleteId !== null}
        onOpenChange={(open) => !open && setDeleteId(null)}
        title={t('features.content.resource.confirmDelete')}
        desc={t('features.content.resource.confirmDeleteDesc')}
        destructive
        confirmText={t('features.content.resource.actions.delete')}
        handleConfirm={() => deleteId && deleteMutation.mutate(deleteId)}
        isLoading={deleteMutation.isPending}
      />
        </>
      }
    >
      <div className='flex min-h-[520px] flex-col'>
        <div className='flex-1'>
          {isLoading ? (
            <div className='flex h-[320px] items-center justify-center'>
              <Loader2 className='h-8 w-8 animate-spin text-muted-foreground' />
            </div>
          ) : error ? (
            <div className='flex h-[320px] items-center justify-center text-destructive'>
              {t('features.content.resource.loadError')}
            </div>
          ) : list.length === 0 ? (
            <div className='flex h-[320px] flex-col items-center justify-center gap-2 text-muted-foreground'>
              <ImageIcon className='h-12 w-12' />
              <p>{t('features.content.resource.noData')}</p>
            </div>
          ) : (
            <div className='grid grid-cols-3 gap-3 sm:grid-cols-4 md:grid-cols-5 lg:grid-cols-6 xl:grid-cols-8'>
              {list.map((resource) => (
                <ResourceCard
                  key={resource.id}
                  resource={resource}
                  onCopyUrl={handleCopyUrl}
                  onEdit={handleEdit}
                  onDelete={handleDelete}
                  onPreview={handlePreview}
                  t={t}
                />
              ))}
            </div>
          )}
        </div>

        <DataTablePagination table={table} className='mt-6 border-t border-border/60 px-0 py-4' />
      </div>
    </ListPageLayout>
  )
}

type ResourceCardProps = {
  resource: ResourceItem
  onCopyUrl: (url: string) => void
  onEdit: (resource: ResourceItem) => void
  onDelete: (resource: ResourceItem) => void
  onPreview: (resource: ResourceItem) => void
  t: (key: string) => string
}

function ResourceCard({ resource, onCopyUrl, onEdit, onDelete, onPreview, t }: ResourceCardProps) {
  const [imgError, setImgError] = useState(false)
  const isImage = resource.type === FILE_TYPE.IMAGE

  return (
    <div className='group flex w-full min-w-0 flex-col rounded-lg border overflow-hidden bg-card'>
      <div
        className='relative w-full aspect-[4/3] min-h-0 overflow-hidden bg-muted cursor-pointer'
        role='button'
        tabIndex={0}
        onClick={() => onPreview(resource)}
        onKeyDown={(e) => e.key === 'Enter' && onPreview(resource)}
        aria-label={t('features.content.resource.actions.preview')}
      >
        {isImage && !imgError ? (
          <img
            src={resource.url}
            alt={resource.name}
            className='h-full w-full object-cover object-center'
            loading='lazy'
            onError={() => setImgError(true)}
          />
        ) : (
          <div className='absolute inset-0 flex items-center justify-center text-muted-foreground text-sm'>
            {isImage && imgError ? t('features.content.resource.loadFailed') : <FileIcon className='h-10 w-10' />}
          </div>
        )}
        <div className='absolute top-1.5 right-1.5 flex items-center gap-1 rounded-md bg-black/60 px-1 py-1 backdrop-blur-sm opacity-0 transition-opacity group-hover:opacity-100'>
          <Button
            type='button'
            variant='ghost'
            size='icon'
            className='h-7 w-7 text-white hover:bg-white/20 hover:text-white'
            onClick={(e) => {
              e.stopPropagation()
              onEdit(resource)
            }}
            title={t('features.content.resource.actions.edit')}
          >
            <Pencil className='h-3.5 w-3.5' />
          </Button>
          <Button
            type='button'
            variant='ghost'
            size='icon'
            className='h-7 w-7 text-white hover:bg-white/20 hover:text-white'
            onClick={(e) => {
              e.stopPropagation()
              onCopyUrl(resource.url)
            }}
            title={t('features.content.resource.actions.copyUrl')}
          >
            <Copy className='h-3.5 w-3.5' />
          </Button>
          <Button
            type='button'
            variant='ghost'
            size='icon'
            className='h-7 w-7 text-white hover:bg-destructive hover:text-destructive-foreground'
            onClick={(e) => {
              e.stopPropagation()
              onDelete(resource)
            }}
            title={t('features.content.resource.actions.delete')}
          >
            <Trash2 className='h-3.5 w-3.5' />
          </Button>
        </div>
      </div>
      {/* 统一文件名区域高度，单行截断 */}
      <div className='h-9 px-2 flex items-center shrink-0' title={resource.name}>
        <span className='truncate text-xs text-muted-foreground w-full'>
          {formatName(resource.name)}
        </span>
      </div>
    </div>
  )
}

type ResourcePreviewOverlayProps = {
  resource: ResourceItem
  onClose: () => void
  onCopyUrl: (url: string) => void
  t: (key: string) => string
}

function ResourcePreviewOverlay({ resource, onClose, onCopyUrl, t }: ResourcePreviewOverlayProps) {
  useEffect(() => {
    const onEscape = (e: KeyboardEvent) => e.key === 'Escape' && onClose()
    window.addEventListener('keydown', onEscape)
    document.body.style.overflow = 'hidden'
    return () => {
      window.removeEventListener('keydown', onEscape)
      document.body.style.overflow = ''
    }
  }, [onClose])

  return (
    <div
      className='fixed inset-0 z-50 flex items-center justify-center bg-black/90'
      onClick={onClose}
      role='dialog'
      aria-modal
      aria-label={t('features.content.resource.actions.preview')}
    >
      <button
        type='button'
        className='absolute top-4 right-4 z-10 rounded-full p-2 text-white hover:bg-white/20 transition-colors'
        onClick={onClose}
        aria-label={t('features.content.resource.actions.close')}
      >
        <X className='h-6 w-6' />
      </button>
      <div
        className='relative max-h-[90vh] max-w-[90vw] flex items-center justify-center'
        onClick={(e) => e.stopPropagation()}
      >
        {resource.type === FILE_TYPE.IMAGE && (
          <img
            src={resource.url}
            alt={resource.name}
            className='max-h-[90vh] max-w-full object-contain'
          />
        )}
        {resource.type === FILE_TYPE.VIDEO && (
          <video
            src={resource.url}
            controls
            autoPlay
            className='max-h-[90vh] max-w-full'
            onClick={(e) => e.stopPropagation()}
          />
        )}
        {resource.type === FILE_TYPE.AUDIO && (
          <div className='rounded-lg bg-card p-6 shadow-xl'>
            <audio src={resource.url} controls autoPlay />
            <p className='mt-2 text-center text-sm text-muted-foreground'>{resource.name}</p>
          </div>
        )}
        {![FILE_TYPE.IMAGE, FILE_TYPE.VIDEO, FILE_TYPE.AUDIO].includes(resource.type as 1 | 2 | 3) && (
          <div className='rounded-lg bg-card p-8 shadow-xl text-center'>
            <FileIcon className='mx-auto h-16 w-16 text-muted-foreground' />
            <p className='mt-4 max-w-xs truncate font-medium' title={resource.name}>
              {resource.name}
            </p>
            <Button
              type='button'
              variant='outline'
              size='sm'
              className='mt-4'
              onClick={(e) => {
                e.stopPropagation()
                onCopyUrl(resource.url)
              }}
            >
              <Copy className='h-4 w-4' />
              {t('features.content.resource.actions.copyLink')}
            </Button>
          </div>
        )}
      </div>
    </div>
  )
}
