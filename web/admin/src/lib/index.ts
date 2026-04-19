/**
 * 工具函数统一导出
 */

export { cn, sleep, getPageNumbers } from './utils'

// 新增 Service Factory
export {
  createCrudService,
  createCachedCrudService,
  type CrudService,
  type ListResponse,
  type ListParams,
  type ServiceFactoryOptions,
} from './service-factory'
