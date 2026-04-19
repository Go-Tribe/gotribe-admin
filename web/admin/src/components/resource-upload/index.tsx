import * as React from 'react'
import { useState, useCallback, useMemo } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { FileIcon, UploadIcon, Loader2, CheckIcon, ChevronLeftIcon, ChevronRightIcon } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog'
import { ScrollArea } from '@/components/ui/scroll-area'
import { getResourceList, uploadResource } from '@/features/content/service/resource'
import type { ResourceItem } from '@/features/content/types/resource'
// eslint-disable-next-line no-duplicate-imports
import { FILE_TYPE } from '@/features/content/types/resource'
import { toast } from 'sonner'

// 对外暴露类型，与 features/content 一致
export type { ResourceItem }
export { FILE_TYPE }

/** 文件类型（与后端一致，用于 type 入参） */
export type ResourceType = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8

// 组件 Props
export interface ResourceUploadProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSelect?: (resource: ResourceItem) => void
  /** 资源类型：1 图片 2 视频 3 音频 4 压缩包 5 文档 6 字体 7 应用 8 未知。封面图传 1，视频传 2 */
  type?: ResourceType
  accept?: string
  title?: string
  description?: string
  /** 每页条数，默认 20 */
  pageSize?: number
}

// 格式化文件大小
function formatFileSize(bytes?: number): string {
  if (!bytes) return '-'
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

const EMPTY_TEXT_BY_TYPE: Record<number, string> = {
  [FILE_TYPE.IMAGE]: '暂无图片',
  [FILE_TYPE.VIDEO]: '暂无视频',
  [FILE_TYPE.AUDIO]: '暂无音频',
  [FILE_TYPE.ARCHIVE]: '暂无压缩包',
  [FILE_TYPE.DOCUMENT]: '暂无文档',
  [FILE_TYPE.FONT]: '暂无字体',
  [FILE_TYPE.APP]: '暂无应用',
  [FILE_TYPE.UNKNOWN]: '暂无资源',
}

const UPLOAD_BUTTON_TEXT_BY_TYPE: Record<number, string> = {
  [FILE_TYPE.IMAGE]: '上传图片',
  [FILE_TYPE.VIDEO]: '上传视频',
  [FILE_TYPE.AUDIO]: '上传音频',
  [FILE_TYPE.ARCHIVE]: '上传',
  [FILE_TYPE.DOCUMENT]: '上传',
  [FILE_TYPE.FONT]: '上传',
  [FILE_TYPE.APP]: '上传',
  [FILE_TYPE.UNKNOWN]: '上传',
}

const DEFAULT_ACCEPT_BY_TYPE: Record<number, string> = {
  [FILE_TYPE.IMAGE]: 'image/*',
  [FILE_TYPE.VIDEO]: 'video/*',
  [FILE_TYPE.AUDIO]: 'audio/*',
  [FILE_TYPE.ARCHIVE]: '*/*',
  [FILE_TYPE.DOCUMENT]: '*/*',
  [FILE_TYPE.FONT]: '*/*',
  [FILE_TYPE.APP]: '*/*',
  [FILE_TYPE.UNKNOWN]: '*/*',
}

const DEFAULT_PAGE_SIZE = 20

function ResourceUpload({
  open,
  onOpenChange,
  onSelect,
  type = FILE_TYPE.IMAGE,
  accept,
  title = '选择资源',
  description = '上传或选择已有资源',
  pageSize = DEFAULT_PAGE_SIZE,
}: ResourceUploadProps) {
  const queryClient = useQueryClient()
  const [uploadProgress, setUploadProgress] = useState<number>(0)
  const [isUploading, setIsUploading] = useState(false)
  const [selectedId, setSelectedId] = useState<number | null>(null)
  const [pageNum, setPageNum] = useState(1)
  const [failedImageIds, setFailedImageIds] = useState<number[]>([])
  const fileInputRef = React.useRef<HTMLInputElement>(null)

  const resourceType = type

  const { data, isLoading, error } = useQuery({
    queryKey: ['resourceList', resourceType, pageNum, pageSize],
    queryFn: () => getResourceList({ type: resourceType, pageNum: pageNum, pageSize: pageSize }),
    enabled: open,
  })

  const resourceList = useMemo(() => data?.list ?? [], [data?.list])
  const total = data?.total ?? 0
  const pageCount = Math.max(1, Math.ceil(total / pageSize))
  const hasPrev = pageNum > 1
  const hasNext = pageNum < pageCount

  const { mutate: uploadMutate } = useMutation({
    mutationFn: ({ file, type: uploadType }: { file: File; type: number }) =>
      uploadResource(file, uploadType, setUploadProgress),
    onSuccess: async () => {
      toast.success('上传成功')
      await queryClient.invalidateQueries({ queryKey: ['resourceList'] })
      setUploadProgress(0)
      setIsUploading(false)
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : '上传失败')
      setUploadProgress(0)
      setIsUploading(false)
    },
  })

  const handleFileSelect = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0]
      if (!file) return
      setIsUploading(true)
      setUploadProgress(0)
      uploadMutate({ file, type: resourceType })
      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
    },
    [resourceType, uploadMutate]
  )

  const handleUploadClick = useCallback(() => {
    fileInputRef.current?.click()
  }, [])

  const handleResourceSelect = useCallback((resource: ResourceItem) => {
    setSelectedId(resource.id)
  }, [])

  const handleConfirmSelect = useCallback(() => {
    const selectedResource = resourceList.find((r) => r.id === selectedId)
    if (selectedResource) {
      onSelect?.(selectedResource)
      onOpenChange(false)
      setSelectedId(null)
    }
  }, [resourceList, selectedId, onSelect, onOpenChange])

  const handleDoubleClick = useCallback(
    (resource: ResourceItem) => {
      onSelect?.(resource)
      onOpenChange(false)
      setSelectedId(null)
    },
    [onSelect, onOpenChange]
  )

  const handleOpenChange = useCallback(
    (newOpen: boolean) => {
      if (!newOpen) {
        setSelectedId(null)
        setUploadProgress(0)
        setIsUploading(false)
        setPageNum(1)
        setFailedImageIds([])
      }
      onOpenChange(newOpen)
    },
    [onOpenChange]
  )

  const inputAccept = accept ?? DEFAULT_ACCEPT_BY_TYPE[resourceType] ?? '*/*'
  const emptyText = EMPTY_TEXT_BY_TYPE[resourceType] ?? '暂无资源'
  const uploadButtonText = UPLOAD_BUTTON_TEXT_BY_TYPE[resourceType] ?? '上传'
  const isImageType = resourceType === FILE_TYPE.IMAGE

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className='sm:max-w-[700px]'>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description}</DialogDescription>
        </DialogHeader>

        <div className='flex items-center justify-between gap-2'>
          <span className='text-sm text-muted-foreground'>选择已有资源</span>
          <input
            ref={fileInputRef}
            type='file'
            accept={inputAccept}
            onChange={handleFileSelect}
            className='hidden'
          />
          <Button
            variant='outline'
            size='sm'
            onClick={handleUploadClick}
            disabled={isUploading}
          >
            {isUploading ? (
              <>
                <Loader2 className='h-4 w-4 animate-spin' />
                {uploadProgress}%
              </>
            ) : (
              <>
                <UploadIcon className='h-4 w-4' />
                {uploadButtonText}
              </>
            )}
          </Button>
        </div>

        {isUploading && (
          <div className='mt-2'>
            <div className='h-2 w-full rounded-full bg-muted'>
              <div
                className='h-full rounded-full bg-primary transition-all duration-300'
                style={{ width: `${uploadProgress}%` }}
              />
            </div>
          </div>
        )}

        <ScrollArea
          className={cn(
            'mt-4 rounded-md border p-2',
            isImageType ? 'h-[400px]' : 'h-[400px]'
          )}
        >
          {isLoading ? (
            <div className='flex h-full min-h-[200px] items-center justify-center'>
              <Loader2 className='h-8 w-8 animate-spin text-muted-foreground' />
            </div>
          ) : error ? (
            <div className='flex h-full min-h-[200px] items-center justify-center text-destructive'>
              加载失败，请重试
            </div>
          ) : resourceList.length === 0 ? (
            <div className='flex min-h-[120px] items-center justify-center text-sm text-muted-foreground'>
              {emptyText}
            </div>
          ) : isImageType ? (
            <div className='grid w-full min-w-0 grid-cols-4 gap-3'>
              {resourceList.map((resource) => (
                <div
                  key={resource.id}
                  className={cn(
                    'group relative min-w-0 cursor-pointer overflow-hidden rounded-lg border-2 bg-muted transition-all',
                    'aspect-square',
                    selectedId === resource.id
                      ? 'border-primary ring-2 ring-primary/30'
                      : 'border-transparent hover:border-primary/50'
                  )}
                  onClick={() => handleResourceSelect(resource)}
                  onDoubleClick={() => handleDoubleClick(resource)}
                >
                  {failedImageIds.includes(resource.id) ? (
                    <div className='absolute inset-0 flex flex-col items-center justify-center gap-2 text-muted-foreground'>
                      <FileIcon className='h-8 w-8' />
                      <span className='px-3 text-center text-xs'>预览不可用</span>
                    </div>
                  ) : (
                    <div className='absolute inset-0 flex items-center justify-center p-3'>
                      <img
                        src={resource.url}
                        alt={resource.name}
                        className='max-h-full max-w-full rounded object-contain'
                        loading='lazy'
                        onError={() => {
                          setFailedImageIds((prev) => (
                            prev.includes(resource.id) ? prev : [...prev, resource.id]
                          ))
                        }}
                      />
                    </div>
                  )}
                  {selectedId === resource.id && (
                    <div className='absolute top-1 right-1 flex h-5 w-5 items-center justify-center rounded-full bg-primary text-primary-foreground'>
                      <CheckIcon className='h-3 w-3' />
                    </div>
                  )}
                  <div className='absolute inset-x-0 bottom-0 bg-linear-to-t from-black/60 to-transparent p-2 opacity-0 transition-opacity group-hover:opacity-100'>
                    <p className='truncate text-xs text-white'>{resource.name}</p>
                    <p className='text-xs text-white/70'>{formatFileSize(resource.size)}</p>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className='divide-y'>
              {resourceList.map((resource) => (
                <div
                  key={resource.id}
                  className={cn(
                    'flex cursor-pointer items-center gap-3 p-3 transition-colors',
                    selectedId === resource.id
                      ? 'bg-primary/10'
                      : 'hover:bg-muted/50'
                  )}
                  onClick={() => handleResourceSelect(resource)}
                  onDoubleClick={() => handleDoubleClick(resource)}
                >
                  <div className='flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-muted'>
                    <FileIcon className='h-5 w-5 text-muted-foreground' />
                  </div>
                  <div className='flex-1 overflow-hidden'>
                    <p className='truncate font-medium'>{resource.name}</p>
                    <p className='text-sm text-muted-foreground'>
                      {formatFileSize(resource.size)}
                      {resource.createdAt && ` · ${resource.createdAt}`}
                    </p>
                  </div>
                  {selectedId === resource.id && (
                    <div className='flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground'>
                      <CheckIcon className='h-4 w-4' />
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </ScrollArea>

        {!isLoading && !error && resourceList.length > 0 && total > pageSize && (
          <div className='flex items-center justify-center gap-2 border-t pt-3'>
            <Button
              type='button'
              variant='outline'
              size='sm'
              disabled={!hasPrev}
              onClick={() => setPageNum((p) => Math.max(1, p - 1))}
            >
              <ChevronLeftIcon className='h-4 w-4' />
              上一页
            </Button>
            <span className='text-sm text-muted-foreground'>
              第 {pageNum} / {pageCount} 页，共 {total} 条
            </span>
            <Button
              type='button'
              variant='outline'
              size='sm'
              disabled={!hasNext}
              onClick={() => setPageNum((p) => Math.min(pageCount, p + 1))}
            >
              下一页
              <ChevronRightIcon className='h-4 w-4' />
            </Button>
          </div>
        )}

        <div className='flex items-center justify-between border-t pt-4'>
          <div className='text-sm text-muted-foreground'>
            {selectedId ? '已选择 1 个资源，双击可直接确认' : '请选择一个资源'}
          </div>
          <div className='flex gap-2'>
            <Button variant='outline' onClick={() => handleOpenChange(false)}>
              取消
            </Button>
            <Button onClick={handleConfirmSelect} disabled={!selectedId}>
              确认选择
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}

export { ResourceUpload }
