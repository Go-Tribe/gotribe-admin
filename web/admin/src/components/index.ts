/**
 * 通用组件统一导出
 * 
 * 使用方式:
 * ```tsx
 * import { ListPageLayout, DebouncedInput, DataTableActions } from '@/components'
 * ```
 */

// 布局组件
export { ListPageLayout } from './layout/list-page-layout'

// 输入组件
export { DebouncedInput, TableFilterInput } from './debounced-input'

// 状态组件
export { StatusBadge, StatusBadgeGroup } from './status-badge'

// 图片组件
export { LazyImage, ImagePreview } from './lazy-image'

// 弹窗组件
export { SchemaFormDialog } from './schema-form-dialog'
export { CodeSplitDialog, LazyDialog } from './code-split-dialog'

// 表单组件
export { ColorPicker, ColorPreview } from './color-picker'

// 表格组件（从 data-table 子目录导出）
export {
  DataTable,
  DataTableOptimized,
  DataTableCard,
  DataTableActions,
  VirtualDataTable,
  DataTablePagination,
  DataTableColumnHeader,
  DataTableFacetedFilter,
  DataTableViewOptions,
  DataTableToolbar,
  DataTableBulkActions,
} from './data-table'

// 筛选组件
export {
  TableFilters,
  TextFilter,
  SelectFilter,
  DateRangeFilter,
  NumberRangeFilter,
  TableToolbar,
  ResetButton,
} from './table-filters'
