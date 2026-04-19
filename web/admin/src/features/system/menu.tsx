import { useMemo, useState, useCallback } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  getExpandedRowModel,
  type ColumnDef,
  type ExpandedState,
} from '@tanstack/react-table'
import { PlusIcon } from '@radix-ui/react-icons'
import { ChevronRight, ChevronDown } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DataTableColumnHeader, TreeDataTable } from '@/components/data-table'
import { ListPageLayout, DataTableActions } from '@/components'
import { Badge } from '@/components/ui/badge'
import { getStatusVariantByType } from '@/config/status-variants'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { MenuFormDialog } from './components/menu-form-dialog'
import type { Menu } from './types/menu'
import { createMenu, updateMenuById, batchDeleteMenuByIds, getMenuTree, getMenus } from './service/menu'
import { useI18n } from '@/context/i18n-provider'
import { useDataTable } from '@/hooks/use-data-table'
import { useCrudMutations } from '@/hooks/use-crud-mutations'

export function SystemMenu() {
  const { t } = useI18n()
  const [dialogOpen, setDialogOpen] = useState<'create' | 'edit' | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<number | null>(null)
  const [editingMenu, setEditingMenu] = useState<Menu | null>(null)

  // 使用统一的表格状态管理
  const {
    columnFilters,
    setColumnFilters,
    sorting,
    setSorting,
    columnVisibility,
    setColumnVisibility,
  } = useDataTable()

  // 树形表格的展开状态（保持独立，因为树形表格特殊）
  const [expanded, setExpanded] = useState<ExpandedState>({})

  // 获取菜单列表数据（用于构建树形结构）
  const { data: menuListData, isLoading, error } = useQuery({
    queryKey: ['menuList'],
    queryFn: () => getMenus(),
  })

  // 获取菜单树数据（用于选择父菜单）
  const { data: menuTreeData } = useQuery({
    queryKey: ['menuTree'],
    queryFn: () => getMenuTree(),
  })

  // 从列表数据构建树形结构
  const menuTree = useMemo(() => {
    if (!menuListData?.menus || menuListData.menus.length === 0) {
      return []
    }
    
    // 构建树形结构
    const menuMap = new Map<number, Menu>()
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
        } else {
          // 如果找不到父菜单，也作为根菜单
          rootMenus.push(menuWithChildren)
        }
      }
    })
    
    return rootMenus
  }, [menuListData?.menus])

  // 用于选择父菜单的树形数据（优先使用 API 返回的树形数据，否则使用构建的树）
  const menuTreeForParent = useMemo(() => {
    return menuTreeData?.menuTree && menuTreeData.menuTree.length > 0
      ? menuTreeData.menuTree
      : menuTree
  }, [menuTreeData?.menuTree, menuTree])

  // 扁平化菜单数据用于查找（编辑时使用）
  const flattenMenuData = useCallback((menus: Menu[]): Menu[] => {
    const result: Menu[] = []
    menus.forEach((menu) => {
      result.push(menu)
      if (menu.children && menu.children.length > 0) {
        result.push(...flattenMenuData(menu.children))
      }
    })
    return result
  }, [])
  
  const allMenus = useMemo(() => flattenMenuData(menuTree), [menuTree, flattenMenuData])

  // 使用统一的 CRUD mutations（menu 接口需要适配）
  const { createMutation, updateMutation, deleteMutation } = useCrudMutations<Menu, number>({
    queryKey: ['menuList'],
    createFn: (data) => createMenu(data),
    // updateMenuById 需要 id 和 data 分开，这里进行适配
    updateFn: (data) => {
      const { id, ...rest } = data as Partial<Menu> & { id: number }
      return updateMenuById(id, rest)
    },
    // batchDeleteMenuByIds 需要包装成接受单个 id 的函数
    deleteFn: (id) => batchDeleteMenuByIds({ menuIds: [id] }),
    messages: {
      createSuccess: t('features.system.menu.createSuccess'),
      updateSuccess: t('features.system.menu.updateSuccess'),
      deleteSuccess: t('features.system.menu.deleteSuccess'),
    },
    onSuccess: () => {
      setDialogOpen(null)
      setEditingMenu(null)
      setDeleteDialogOpen(null)
    },
  })

  // 列定义
  const columns = useMemo<ColumnDef<Menu>[]>(
    () => [
      {
        accessorKey: 'title',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.menu.menuTitle')} />
        ),
        cell: ({ row }) => {
          const canExpand = row.getCanExpand()
          const isExpanded = row.getIsExpanded()
          const level = row.depth
          return (
            <div 
              className='flex items-center gap-2 font-medium cursor-pointer' 
              style={{ paddingLeft: `${level * 20}px` }}
              onClick={() => canExpand && row.toggleExpanded()}
            >
              {canExpand && (
                <button
                  type='button'
                  onClick={(e) => {
                    e.stopPropagation()
                    row.toggleExpanded()
                  }}
                  className='flex items-center justify-center w-4 h-4 hover:bg-accent rounded'
                >
                  {isExpanded ? (
                    <ChevronDown className='h-4 w-4' />
                  ) : (
                    <ChevronRight className='h-4 w-4' />
                  )}
                </button>
              )}
              {!canExpand && <div className='w-4' />}
              <span className={canExpand ? 'hover:text-primary' : ''}>{row.getValue('title')}</span>
            </div>
          )
        },
      },
      {
        accessorKey: 'name',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.menu.name')} />
        ),
        cell: ({ row }) => <div>{row.getValue('name')}</div>,
      },
      {
        accessorKey: 'icon',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.menu.icon')} />
        ),
        cell: ({ row }) => {
          const icon = row.getValue('icon') as string
          return <div className='max-w-[100px] truncate' title={icon}>{icon || '-'}</div>
        },
      },
      {
        accessorKey: 'path',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.menu.path')} />
        ),
        cell: ({ row }) => {
          const path = row.getValue('path') as string
          return <div className='max-w-[150px] truncate' title={path}>{path || '-'}</div>
        },
      },
      {
        accessorKey: 'component',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.menu.component')} />
        ),
        cell: ({ row }) => {
          const component = row.getValue('component') as string
          return <div className='max-w-[150px] truncate' title={component}>{component || '-'}</div>
        },
      },
      {
        accessorKey: 'redirect',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.menu.redirect')} />
        ),
        cell: ({ row }) => {
          const redirect = row.getValue('redirect') as string
          return <div className='max-w-[150px] truncate' title={redirect || ''}>{redirect || '-'}</div>
        },
      },
      {
        accessorKey: 'sort',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.menu.sort')} />
        ),
        cell: ({ row }) => <div>{row.getValue('sort')}</div>,
      },
      {
        accessorKey: 'status',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.menu.status')} />
        ),
        cell: ({ row }) => {
          const status = row.getValue('status') as number
          return (
            <Badge variant={getStatusVariantByType(status, 'enabledStatus')}>
              {status === 1 ? t('features.system.menu.enabled') : t('features.system.menu.disabled')}
            </Badge>
          )
        },
      },
      {
        accessorKey: 'hidden',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.menu.hidden')} />
        ),
        cell: ({ row }) => {
          const hidden = row.getValue('hidden') as number
          return (
            <Badge variant={getStatusVariantByType(hidden, 'visibilityInverse')}>
              {hidden === 1 ? t('features.system.menu.hidden') : t('features.system.menu.visible')}
            </Badge>
          )
        },
      },
      {
        accessorKey: 'noCache',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.menu.cache')} />
        ),
        cell: ({ row }) => {
          const noCache = row.getValue('noCache') as number
          return (
            <Badge variant={getStatusVariantByType(noCache, 'visibilityInverse')}>
              {noCache === 1 ? t('features.system.menu.noCache') : t('features.system.menu.cached')}
            </Badge>
          )
        },
      },
      {
        accessorKey: 'activeMenu',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('features.system.menu.activeMenu')} />
        ),
        cell: ({ row }) => {
          const activeMenu = row.getValue('activeMenu') as string
          return <div className='max-w-[150px] truncate' title={activeMenu || ''}>{activeMenu || '-'}</div>
        },
      },
      {
        id: 'actions',
        header: t('features.system.menu.actions'),
        cell: ({ row }) => (
          <DataTableActions
            onEdit={() => {
              setEditingMenu(row.original)
              setDialogOpen('edit')
            }}
            onDelete={() => setDeleteDialogOpen(row.original.id)}
            deleteConfirmTitle={t('features.system.menu.confirmDelete')}
            useDropdown={false}
          />
        ),
        enableHiding: false,
        meta: {
          className: 'sticky right-0 bg-background z-10 shadow-[inset_-1px_0_0_0_hsl(var(--border))]',
          thClassName: 'sticky right-0 bg-background z-10 shadow-[inset_-1px_0_0_0_hsl(var(--border))]',
          tdClassName: 'sticky right-0 bg-background z-10 shadow-[inset_-1px_0_0_0_hsl(var(--border))]',
        },
      },
    ],
    [t]
  )

  // 表格实例（树形表格）
  // eslint-disable-next-line react-hooks/incompatible-library
  const table = useReactTable({
    data: menuTree,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getExpandedRowModel: getExpandedRowModel(),
    getSubRows: (row) => row.children,
    onExpandedChange: setExpanded,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    onColumnVisibilityChange: setColumnVisibility,
    state: {
      expanded,
      sorting,
      columnFilters,
      columnVisibility,
    },
  })

  // 处理新建
  const handleCreate = useCallback(() => {
    setEditingMenu(null)
    setDialogOpen('create')
  }, [])

  // 处理表单提交
  const handleFormSubmit = (data: Partial<Menu>) => {
    if (dialogOpen === 'create') {
      createMutation.mutate(data)
    } else if (dialogOpen === 'edit' && editingMenu) {
      updateMutation.mutate({ ...data, id: editingMenu.id })
    }
  }

  // 处理删除
  const handleDelete = useCallback(() => {
    if (deleteDialogOpen) {
      deleteMutation.mutate(deleteDialogOpen)
    }
  }, [deleteDialogOpen, deleteMutation])

  const menuToDelete = allMenus.find((m: Menu) => m.id === deleteDialogOpen)

  // 收集菜单及其所有子菜单的 ID
  const collectMenuIds = useCallback((menu: Menu): number[] => {
    const ids = [menu.id]
    if (menu.children && menu.children.length > 0) {
      menu.children.forEach((child) => {
        ids.push(...collectMenuIds(child))
      })
    }
    return ids
  }, [])

  // 扁平化菜单树，用于选择父菜单（只包含有 children 的菜单）
  const flattenMenuTree = useCallback((menus: Menu[], parentId: number | null = null, level: number = 0, excludeIds: number[] = []): Array<{ id: number; title: string; parentID: number | null; level: number }> => {
    const result: Array<{ id: number; title: string; parentID: number | null; level: number }> = []
    menus.forEach((menu) => {
      if (!excludeIds.includes(menu.id)) {
        // 只添加有 children 的菜单项
        if (menu.children && menu.children.length > 0) {
          result.push({
            id: menu.id,
            title: menu.title,
            parentID: parentId,
            level,
          })
          result.push(...flattenMenuTree(menu.children, menu.id, level + 1, excludeIds))
        }
      }
    })
    return result
  }, [])

  // 获取需要排除的菜单 ID 列表（编辑菜单及其所有子菜单）
  const excludeIds = useMemo(() => {
    if (!editingMenu) return []
    // 从菜单树中找到编辑菜单的完整信息（包含 children）
    const findMenuInTree = (menus: Menu[]): Menu | null => {
      for (const menu of menus) {
        if (menu.id === editingMenu.id) {
          return menu
        }
        if (menu.children && menu.children.length > 0) {
          const found = findMenuInTree(menu.children)
          if (found) return found
        }
      }
      return null
    }
    const fullMenu = findMenuInTree(menuTreeForParent)
    return fullMenu ? collectMenuIds(fullMenu) : [editingMenu.id]
  }, [editingMenu, menuTreeForParent, collectMenuIds])

  const parentMenuOptions = useMemo(() => {
    return flattenMenuTree(menuTreeForParent, null, 0, excludeIds).map(menu => ({
      id: menu.id,
      title: '  '.repeat(menu.level) + menu.title,
      parentID: menu.parentID,
    }))
  }, [menuTreeForParent, excludeIds, flattenMenuTree])

  const isDialogLoading = dialogOpen === 'create'
    ? createMutation.isPending
    : updateMutation.isPending

  return (
    <ListPageLayout
      title={t('features.system.menu.title')}
      description={t('features.system.menu.description')}
      actions={
        <Button onClick={handleCreate}>
          <PlusIcon className='h-4 w-4' />
          {t('features.system.menu.createButton')}
        </Button>
      }
      dialogs={
        <>
          <MenuFormDialog
            open={dialogOpen !== null}
            onOpenChange={(open) => {
              if (!open) {
                setDialogOpen(null)
                setEditingMenu(null)
              }
            }}
            menu={editingMenu}
            onSubmit={handleFormSubmit}
            isLoading={isDialogLoading}
            parentMenus={parentMenuOptions}
          />
          <ConfirmDialog
            open={deleteDialogOpen !== null}
            onOpenChange={(open) => {
              if (!open) setDeleteDialogOpen(null)
            }}
            title={t('features.system.menu.confirmDelete')}
            desc={t('features.system.menu.confirmDeleteMessage', { title: menuToDelete?.title })}
            handleConfirm={handleDelete}
            destructive
            confirmText={t('features.system.menu.delete')}
            isLoading={deleteMutation.isPending}
          />
        </>
      }
    >
      <TreeDataTable
        table={table}
        columns={columns}
        isLoading={isLoading}
        error={error}
        loadingText={t('features.system.menu.loading')}
        errorText={t('features.system.menu.loadError')}
        emptyText={t('features.system.menu.noData')}
        bordered={false}
      />
    </ListPageLayout>
  )
}
