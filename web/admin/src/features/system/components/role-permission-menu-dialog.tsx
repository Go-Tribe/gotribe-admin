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
import { getMenus } from '../service/menu'
import { getRoleMenus, updateRoleMenus } from '../service/role'
import { toast } from 'sonner'
import type { Menu } from '../types/menu'
import { ChevronRightIcon, ChevronDownIcon } from '@radix-ui/react-icons'
import { cn } from '@/lib/utils'

type RolePermissionDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  roleId: number | null
  roleName?: string
}

// 获取所有子菜单ID（包括子菜单的子菜单）
function getAllChildrenIds(menuItem: Menu): number[] {
  const ids: number[] = []
  if (menuItem.children) {
    menuItem.children.forEach((child) => {
      ids.push(child.id)
      ids.push(...getAllChildrenIds(child))
    })
  }
  return ids
}

// 获取所有父菜单ID（从当前菜单向上到根）
function getAllParentIds(
  menuId: number,
  parentMap: Map<number, number>
): number[] {
  const ids: number[] = []
  let currentId = menuId
  while (parentMap.has(currentId)) {
    const parentId = parentMap.get(currentId)!
    ids.push(parentId)
    currentId = parentId
  }
  return ids
}

// 检查菜单的所有子节点是否全部选中
function areAllChildrenSelected(
  menu: Menu,
  selectedIds: Set<number>
): boolean {
  if (!menu.children || menu.children.length === 0) {
    return true
  }
  return menu.children.every(
    (child) =>
      selectedIds.has(Number(child.id)) &&
      areAllChildrenSelected(child, selectedIds)
  )
}

// 检查菜单是否有任意子节点被选中
function hasAnyChildSelected(menu: Menu, selectedIds: Set<number>): boolean {
  if (!menu.children || menu.children.length === 0) {
    return false
  }
  return menu.children.some(
    (child) =>
      selectedIds.has(Number(child.id)) ||
      hasAnyChildSelected(child, selectedIds)
  )
}

// 树形菜单项组件
function MenuTreeItem({
  menu,
  selectedMenuIds,
  onToggleWithChildren,
  level = 0,
}: {
  menu: Menu
  selectedMenuIds: Set<number>
  onToggleWithChildren: (menuId: number, checked: boolean, childrenIds: number[]) => void
  level?: number
}) {
  const [expanded, setExpanded] = useState(true)
  const hasChildren = menu.children && menu.children.length > 0

  // 确保 menu.id 是数字类型进行比较
  const menuId = Number(menu.id)
  const isChecked = selectedMenuIds.has(menuId)
  const isIndeterminate = useMemo(() => {
    if (!hasChildren) return false
    // 检查是否有子节点被选中，但不是全部选中
    const hasSelected = hasAnyChildSelected(menu, selectedMenuIds)
    const allSelected = areAllChildrenSelected(menu, selectedMenuIds)
    return hasSelected && !allSelected
  }, [menu, selectedMenuIds, hasChildren])

  const handleCheckboxChange = (checked: boolean | 'indeterminate') => {
    // Radix UI 的 onCheckedChange 可能返回 'indeterminate'，我们需要转换为 boolean
    const isCheckedValue = checked === true
    // 获取所有子菜单ID
    const allChildrenIds = hasChildren ? getAllChildrenIds(menu) : []
    // 调用统一的切换处理函数
    onToggleWithChildren(Number(menu.id), isCheckedValue, allChildrenIds)
  }

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
          {menu.title || menu.name}
        </label>
      </div>
      {hasChildren && expanded && (
        <div>
          {menu.children!.map((child) => (
            <MenuTreeItem
              key={child.id}
              menu={child}
              selectedMenuIds={selectedMenuIds}
              onToggleWithChildren={onToggleWithChildren}
              level={level + 1}
            />
          ))}
        </div>
      )}
    </div>
  )
}

export function RolePermissionDialog({
  open,
  onOpenChange,
  roleId,
  roleName,
}: RolePermissionDialogProps) {
  const { t } = useI18n()
  const [selectedMenuIds, setSelectedMenuIds] = useState<Set<number>>(new Set())
  const [isSubmitting, setIsSubmitting] = useState(false)

  // 获取菜单列表
  const { data: menuListData, isLoading: isLoadingMenus } = useQuery({
    queryKey: ['menuList'],
    queryFn: () => getMenus(),
    enabled: open,
  })

  // 获取角色已有的菜单权限
  const { data: roleMenusData, isLoading: isLoadingRoleMenus } = useQuery({
    queryKey: ['roleMenus', roleId],
    queryFn: () => (roleId ? getRoleMenus(roleId) : Promise.resolve({ menus: [] })),
    enabled: open && roleId !== null,
  })

  // 构建菜单树和相关映射
  const { menuTree, menuMap, parentMap } = useMemo(() => {
    if (!menuListData?.menus || menuListData.menus.length === 0) {
      return { menuTree: [], menuMap: new Map<number, Menu>(), parentMap: new Map<number, number>() }
    }

    // 构建树形结构
    const menuMap = new Map<number, Menu>()
    const parentMap = new Map<number, number>() // 子ID -> 父ID
    const rootMenus: Menu[] = []

    // 第一遍：创建所有菜单的映射，并初始化 children
    menuListData.menus.forEach((menu) => {
      menuMap.set(menu.id, { ...menu, children: [] })
    })

    // 第二遍：构建父子关系
    menuListData.menus.forEach((menu) => {
      const menuWithChildren = menuMap.get(menu.id)!
      if (menu.parentID === 0 || !menu.parentID) {
        rootMenus.push(menuWithChildren)
      } else {
        const parent = menuMap.get(menu.parentID)
        if (parent) {
          if (!parent.children) {
            parent.children = []
          }
          parent.children.push(menuWithChildren)
          // 记录父子关系
          parentMap.set(menu.id, menu.parentID)
        } else {
          // 如果找不到父菜单，也作为根菜单
          rootMenus.push(menuWithChildren)
        }
      }
    })

    return { menuTree: rootMenus, menuMap, parentMap }
  }, [menuListData?.menus])

  // 根据子节点选中状态更新父节点状态的辅助函数
  const updateParentStates = useCallback(
    (currentIds: Set<number>): Set<number> => {
      if (menuMap.size === 0) return currentIds

      const newSet = new Set(currentIds)

      // 从叶子节点向上检查每个父节点
      // 遍历所有菜单，检查其父节点状态
      menuMap.forEach((menu, menuId) => {
        if (menu.children && menu.children.length > 0) {
          // 这是一个父节点，检查其所有子节点是否全部选中
          const allChildrenSelected = areAllChildrenSelected(menu, newSet)
          if (allChildrenSelected) {
            newSet.add(menuId)
          } else {
            newSet.delete(menuId)
          }
        }
      })

      return newSet
    },
    [menuMap]
  )

  // 初始化选中的菜单ID
  useEffect(() => {
    if (!open) {
      // 对话框关闭时，重置选中状态
      setSelectedMenuIds(new Set())
      return
    }

    if (!roleId) {
      setSelectedMenuIds(new Set())
      return
    }

    // 当数据加载完成时，设置选中的菜单ID
    if (!isLoadingRoleMenus && menuMap.size > 0) {
      if (roleMenusData?.menus && Array.isArray(roleMenusData.menus)) {
        // 从菜单对象数组中提取ID，确保是数字类型，过滤掉无效值
        const menuIds = roleMenusData.menus
          .map((menu) => Number(menu.id))
          .filter((id) => !isNaN(id) && id > 0)
        const initialSet = new Set(menuIds)
        // 根据子节点状态更新父节点状态
        const updatedSet = updateParentStates(initialSet)
        setSelectedMenuIds(updatedSet)
      } else {
        // 如果加载完成但没有数据，清空选中状态
        setSelectedMenuIds(new Set())
      }
    }
  }, [roleMenusData, open, isLoadingRoleMenus, roleId, menuMap.size, updateParentStates])

  // 切换菜单选择状态（包含子菜单的批量操作）
  const handleToggleWithChildren = useCallback(
    (menuId: number, checked: boolean, childrenIds: number[]) => {
      setSelectedMenuIds((prev) => {
        const newSet = new Set(prev)

        // 更新当前菜单
        if (checked) {
          newSet.add(menuId)
        } else {
          newSet.delete(menuId)
        }

        // 更新所有子菜单
        childrenIds.forEach((childId) => {
          if (checked) {
            newSet.add(Number(childId))
          } else {
            newSet.delete(Number(childId))
          }
        })

        // 更新所有父级菜单状态
        const parentIds = getAllParentIds(menuId, parentMap)
        parentIds.forEach((parentId) => {
          const parentMenu = menuMap.get(parentId)
          if (parentMenu) {
            // 检查该父节点的所有子节点是否全部选中
            const allChildrenSelected = areAllChildrenSelected(parentMenu, newSet)
            if (allChildrenSelected) {
              newSet.add(parentId)
            } else {
              newSet.delete(parentId)
            }
          }
        })

        return newSet
      })
    },
    [parentMap, menuMap]
  )

  // 全选/全不选
  const handleSelectAll = useCallback(() => {
    if (!menuListData?.menus) return

    // 确保所有ID都是数字类型
    const allMenuIds = new Set(menuListData.menus.map((m) => Number(m.id)))
    if (selectedMenuIds.size === allMenuIds.size) {
      // 全不选
      setSelectedMenuIds(new Set())
    } else {
      // 全选
      setSelectedMenuIds(allMenuIds)
    }
  }, [menuListData?.menus, selectedMenuIds.size])

  // 提交保存
  const handleSubmit = async () => {
    if (!roleId) return

    setIsSubmitting(true)
    try {
      await updateRoleMenus(roleId, {
        menuIds: Array.from(selectedMenuIds),
      })
      toast.success(t('features.system.role.permission.updateSuccess'))
      onOpenChange(false)
    } catch {
      // 错误已由响应拦截器统一处理，这里不需要再次弹出错误消息
    } finally {
      setIsSubmitting(false)
    }
  }

  const isLoading = isLoadingMenus || isLoadingRoleMenus

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-[600px] max-h-[90vh] flex flex-col'>
        <DialogHeader className='shrink-0'>
          <DialogTitle>
            {t('features.system.role.permission.title', { name: roleName || '' })}
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
              {selectedMenuIds.size === menuListData?.menus?.length
                ? t('features.system.role.permission.unselectAll')
                : t('features.system.role.permission.selectAll')}
            </Button>
            <span className='text-sm text-muted-foreground'>
              {t('features.system.role.permission.selectedCount', {
                count: selectedMenuIds.size,
                total: menuListData?.menus?.length || 0,
              })}
            </span>
          </div>

          {/* 菜单树 */}
          <div className='flex-1 h-full border rounded-md overflow-hidden overflow-y-auto'>
            <ScrollArea className='h-full'>
              {isLoading ? (
                <div className='flex items-center justify-center h-32 text-muted-foreground'>
                  {t('features.system.role.loading')}
                </div>
              ) : menuTree.length === 0 ? (
                <div className='flex items-center justify-center h-32 text-muted-foreground'>
                  {t('features.system.role.permission.noMenus')}
                </div>
              ) : (
                <div className='p-4'>
                  {menuTree.map((menu) => (
                    <MenuTreeItem
                      key={menu.id}
                      menu={menu}
                      selectedMenuIds={selectedMenuIds}
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
