import { request } from '@/service'

export interface Role {
  id: number
  createdAt: string
  updatedAt: string
  deletedAt: string | null
  name: string
  keyword: string
  desc: string
  status: number
  sort: number
  creator: string
  admins: unknown[]
  menus: unknown[]
}

export interface RoleListData {
  roles: Role[]
  total: number
}

export type MenuItem = {
  id: number;
  createdAt: string;
  updatedAt: string;
  deletedAt: string | null;
  name: string;
  title: string;
  icon: string;
  path: string;
  redirect: string | null;
  component: string | null;
  sort: number;
  status: number;
  hidden: number;
  noCache: number;
  alwaysShow: number;
  breadcrumb: number;
  activeMenu: string | null;
  parentID: number;
  creator: string;
  children: MenuItem[];
  roles: unknown[] | null;
};

export type MenuList = {
  menuTree?: MenuItem[];
};

export async function getRoleList(): Promise<RoleListData> {
  return request.get<RoleListData>('/api/role/list')
}

export async function getMenuAccessTree(parentId: number): Promise<MenuList> {
  return request.get<MenuList>(`/api/menu/access/tree/${parentId}`)
}