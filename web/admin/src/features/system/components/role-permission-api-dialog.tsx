import { useEffect, useMemo, useState, useCallback } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { ScrollArea } from '@/components/ui/scroll-area'
import { useI18n } from '@/context/i18n-provider'
import { getRoleApis, updateRoleApis, getApis } from '../service/role'
import { toast } from 'sonner'
import type { Api } from '../types/role'
import { ChevronRightIcon, ChevronDownIcon } from '@radix-ui/react-icons'
import { cn } from '@/lib/utils'

type RolePermissionDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  roleId: number | null
  roleName?: string
}

// 获取所有子API的ID（包括子API的子API），只获取有效ID（大于0）
function getAllChildrenIds(apiItem: Api): number[] {
  const ids: number[] = []
  if (apiItem.children && apiItem.children.length > 0) {
    apiItem.children.forEach((child) => {
      if (child.id > 0) {
        ids.push(child.id)
      }
      ids.push(...getAllChildrenIds(child))
    })
  }
  return ids
}

// 获取节点的选中状态
type CheckedState = 'checked' | 'unchecked' | 'indeterminate'

function getNodeCheckedState(api: Api, selectedIds: Set<number>): CheckedState {
  const hasChildren = api.children && api.children.length > 0

  if (!hasChildren) {
    // 叶子节点：直接检查是否在选中集合中
    return selectedIds.has(api.id) ? 'checked' : 'unchecked'
  }

  // 父节点：检查所有子节点的状态
  const childStates = api.children.map((child) => getNodeCheckedState(child, selectedIds))
  const allChecked = childStates.every((s) => s === 'checked')
  const allUnchecked = childStates.every((s) => s === 'unchecked')

  if (allChecked) return 'checked'
  if (allUnchecked) return 'unchecked'
  return 'indeterminate'
}

// 树形API项组件
function ApiTreeItem({
  api,
  selectedApiIds,
  onToggle,
  onToggleWithChildren,
  level = 0,
}: {
  api: Api
  selectedApiIds: Set<number>
  onToggle: (apiId: number, checked: boolean) => void
  onToggleWithChildren: (api: Api, checked: boolean) => void
  level?: number
}) {
  const [expanded, setExpanded] = useState(true)
  const hasChildren = api.children && api.children.length > 0

  // 使用统一的状态计算函数
  const checkedState = useMemo(() => {
    return getNodeCheckedState(api, selectedApiIds)
  }, [api, selectedApiIds])

  const isChecked = checkedState === 'checked'
  const isIndeterminate = checkedState === 'indeterminate'

  const handleCheckboxChange = (checked: boolean | 'indeterminate') => {
    // Radix UI 的 onCheckedChange 可能返回 'indeterminate'，我们需要转换为 boolean
    const isCheckedValue = checked === true
    // 使用新的处理函数，会自动处理子节点和父节点状态
    onToggleWithChildren(api, isCheckedValue)
  }

  // 显示名称：优先显示描述，否则显示分类
  const displayName = api.desc || api.category || `${api.method} ${api.path}`

  return (
    <div className='select-none'>
      <div
        className={cn(
          'flex items-center gap-2 py-1.5 px-2 rounded-md hover:bg-accent transition-colors',
          level > 0 && 'ml-4'
        )}
        style={{ paddingLeft: `${level * 16 + 8}px` }}
      >
        {hasChildren ? (
          <button
            type='button'
            onClick={() => setExpanded(!expanded)}
            className='flex items-center justify-center w-4 h-4 hover:bg-accent rounded'
          >
            {expanded ? (
              <ChevronDownIcon className='h-3 w-3' />
            ) : (
              <ChevronRightIcon className='h-3 w-3' />
            )}
          </button>
        ) : (
          <div className='w-4' />
        )}
        <div className='relative flex items-center justify-center'>
          <Checkbox
            checked={isChecked}
            onCheckedChange={handleCheckboxChange}
            className={cn(
              isIndeterminate && !isChecked && 'border-primary bg-primary/50'
            )}
          />
          {isIndeterminate && !isChecked && (
            <div className='absolute inset-0 flex items-center justify-center pointer-events-none'>
              <div className='w-2 h-0.5 bg-primary-foreground rounded' />
            </div>
          )}
        </div>
        <label
          className='flex-1 cursor-pointer text-sm'
          onClick={() => handleCheckboxChange(!isChecked)}
        >
          {displayName}
        </label>
      </div>
      {hasChildren && expanded && (
        <div>
          {api.children.map((child) => (
            <ApiTreeItem
              key={child.id}
              api={child}
              selectedApiIds={selectedApiIds}
              onToggle={onToggle}
              onToggleWithChildren={onToggleWithChildren}
              level={level + 1}
            />
          ))}
        </div>
      )}
    </div>
  )
}

export function RoleApiPermissionDialog({
  open,
  onOpenChange,
  roleId,
  roleName,
}: RolePermissionDialogProps) {
  const { t } = useI18n()
  const [selectedApiIds, setSelectedApiIds] = useState<Set<number>>(new Set())
  const [isSubmitting, setIsSubmitting] = useState(false)

  // 获取API列表
  const { data: apiListData, isLoading: isLoadingApis } = useQuery({
    queryKey: ['apiList'],
    queryFn: () => getApis(),
    enabled: open,
  })

  // 获取角色已有的API权限
  const { data: roleApisData, isLoading: isLoadingRoleApis, isFetched } = useQuery({
    queryKey: ['roleApis', roleId],
    queryFn: () => (roleId ? getRoleApis(roleId) : Promise.resolve({ apis: [] })),
    enabled: open && roleId !== null,
    staleTime: 0, // 确保每次打开都获取最新数据
  })

  // API树（接口返回的已经是树形结构）
  const apiTree = useMemo(() => {
    if (!apiListData?.apis || apiListData.apis.length === 0) {
      return []
    }
    return apiListData.apis
  }, [apiListData?.apis])

  // 获取所有API的ID（用于全选功能），过滤掉小于等于0的ID
  const getAllApiIds = useCallback((apis: Api[]): number[] => {
    const ids: number[] = []
    apis.forEach((api) => {
      if (api.id > 0) {
        ids.push(api.id)
      }
      if (api.children && api.children.length > 0) {
        ids.push(...getAllApiIds(api.children))
      }
    })
    return ids
  }, [])

  const allApiIds = useMemo(() => {
    if (!apiListData?.apis) return []
    return getAllApiIds(apiListData.apis)
  }, [apiListData?.apis, getAllApiIds])

  // 初始化选中的API ID
  useEffect(() => {
    if (!open) {
      // 对话框关闭时，重置选中状态
      setSelectedApiIds(new Set())
      return
    }

    if (!roleId) {
      setSelectedApiIds(new Set())
      return
    }

    // 当数据获取完成时，设置选中的API ID
    if (isFetched && !isLoadingRoleApis) {
      if (roleApisData?.apis && Array.isArray(roleApisData.apis)) {
        // API ID 已经是数字数组，过滤掉小于等于0的ID
        const apiIds = roleApisData.apis
          .map((item) => Number(item.id))
          .filter((id) => !isNaN(id) && id > 0)
        setSelectedApiIds(new Set(apiIds))
      } else {
        // 如果加载完成但没有数据，清空选中状态
        setSelectedApiIds(new Set())
      }
    }
  }, [roleApisData, open, isLoadingRoleApis, isFetched, roleId])

  // 切换单个API选择状态
  const handleToggleApi = useCallback((apiId: number, checked: boolean) => {
    if (apiId <= 0) return // 忽略无效ID
    setSelectedApiIds((prev) => {
      const newSet = new Set(prev)
      if (checked) {
        newSet.add(apiId)
      } else {
        newSet.delete(apiId)
      }
      return newSet
    })
  }, [])

  // 切换API及其所有子节点的选择状态
  const handleToggleWithChildren = useCallback((api: Api, checked: boolean) => {
    setSelectedApiIds((prev) => {
      const newSet = new Set(prev)

      // 处理当前节点（只有有效ID才加入）
      if (api.id > 0) {
        if (checked) {
          newSet.add(api.id)
        } else {
          newSet.delete(api.id)
        }
      }

      // 处理所有子节点
      const allChildrenIds = getAllChildrenIds(api)
      allChildrenIds.forEach((childId) => {
        if (checked) {
          newSet.add(childId)
        } else {
          newSet.delete(childId)
        }
      })

      return newSet
    })
  }, [])

  // 全选/全不选
  const handleSelectAll = useCallback(() => {
    if (allApiIds.length === 0) return

    if (selectedApiIds.size === allApiIds.length) {
      // 全不选
      setSelectedApiIds(new Set())
    } else {
      // 全选
      setSelectedApiIds(new Set(allApiIds))
    }
  }, [allApiIds, selectedApiIds.size])

  // 提交保存
  const handleSubmit = async () => {
    if (!roleId) return

    setIsSubmitting(true)
    try {
      await updateRoleApis(roleId, {
        apiIds: Array.from(selectedApiIds),
      })
      toast.success(t('features.system.role.permission.updateSuccess'))
      onOpenChange(false)
    } catch {
      // 错误已由响应拦截器统一处理，这里不需要再次弹出错误消息
    } finally {
      setIsSubmitting(false)
    }
  }

  const isLoading = isLoadingApis || isLoadingRoleApis

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-[600px] max-h-[90vh] flex flex-col'>
        <DialogHeader className='shrink-0'>
          <DialogTitle>
            {t('features.system.role.permission.apiTitle', { name: roleName || '' })}
          </DialogTitle>
        </DialogHeader>

        <div className='flex-1 min-h-0 flex flex-col space-y-4'>
          {/* 操作栏 */}
          <div className='flex items-center justify-between shrink-0'>
            <Button
              type='button'
              variant='outline'
              size='sm'
              onClick={handleSelectAll}
              disabled={isLoading}
            >
              {selectedApiIds.size === allApiIds.length
                ? t('features.system.role.permission.unselectAll')
                : t('features.system.role.permission.selectAll')}
            </Button>
            <span className='text-sm text-muted-foreground'>
              {t('features.system.role.permission.selectedCount', {
                count: selectedApiIds.size,
                total: allApiIds.length,
              })}
            </span>
          </div>

          {/* API树 */}
          <div className='flex-1 h-full border rounded-md overflow-hidden overflow-y-auto'>
            <ScrollArea className='h-full'>
              {isLoading ? (
                <div className='flex items-center justify-center h-32 text-muted-foreground'>
                  {t('features.system.role.loading')}
                </div>
              ) : apiTree.length === 0 ? (
                <div className='flex items-center justify-center h-32 text-muted-foreground'>
                  {t('features.system.role.permission.noMenus')}
                </div>
              ) : (
                <div className='p-4'>
                  {apiTree.map((api) => (
                    <ApiTreeItem
                      key={api.id || api.category}
                      api={api}
                      selectedApiIds={selectedApiIds}
                      onToggle={handleToggleApi}
                      onToggleWithChildren={handleToggleWithChildren}
                    />
                  ))}
                </div>
              )}
            </ScrollArea>
          </div>
        </div>

        <DialogFooter className='shrink-0 pt-4 border-t mt-4'>
          <Button
            type='button'
            variant='outline'
            onClick={() => onOpenChange(false)}
            disabled={isSubmitting}
          >
            {t('features.system.role.form.cancel')}
          </Button>
          <Button
            type='button'
            disabled={isSubmitting || !roleId}
            onClick={handleSubmit}
          >
            {isSubmitting
              ? t('features.system.role.form.submitting')
              : t('features.system.role.permission.save')}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
