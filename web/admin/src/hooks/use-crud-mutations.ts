import { useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { useI18n } from '@/context/i18n-provider'

export interface CrudMutationsConfig<T, ID = number, CreateInput = Partial<T>, UpdateInput = Partial<T>> {
  /** Query key for list data, used for cache invalidation */
  queryKey: string[]
  /** Create function */
  createFn: (data: CreateInput) => Promise<unknown>
  /** Update function */
  updateFn: (data: UpdateInput) => Promise<unknown>
  /** Delete function */
  deleteFn: (id: ID) => Promise<unknown>
  /** Success messages (optional) */
  messages?: {
    createSuccess?: string
    updateSuccess?: string
    deleteSuccess?: string
    createDescription?: string
    updateDescription?: string
    deleteDescription?: string
    createLoading?: string
    updateLoading?: string
    deleteLoading?: string
  }
  /** Callbacks */
  onSuccess?: () => void
  onError?: (error: Error, action: 'create' | 'update' | 'delete') => void
}

export interface UseCrudMutationsReturn<T, ID = number, CreateInput = Partial<T>, UpdateInput = Partial<T>> {
  createMutation: ReturnType<typeof useMutation<unknown, Error, CreateInput>>
  updateMutation: ReturnType<typeof useMutation<unknown, Error, UpdateInput>>
  deleteMutation: ReturnType<typeof useMutation<unknown, Error, ID>>
  isLoading: boolean
}

/**
 * 统一的 CRUD 操作 Hook
 * 
 * 使用示例:
 * ```typescript
 * const { createMutation, updateMutation, deleteMutation, isLoading } = useCrudMutations<Admin>({
 *   queryKey: ['adminList'],
 *   createFn: createAdmin,
 *   updateFn: updateAdmin,
 *   deleteFn: deleteAdmin,
 *   messages: {
 *     createSuccess: t('createSuccess'),
 *     updateSuccess: t('updateSuccess'),
 *     deleteSuccess: t('deleteSuccess'),
 *   },
 *   onSuccess: () => {
 *     setDialogOpen(null)
 *     setEditingEntity(null)
 *   },
 * })
 * ```
 */
export function useCrudMutations<T, ID = number, CreateInput = Partial<T>, UpdateInput = Partial<T>>(
  config: CrudMutationsConfig<T, ID, CreateInput, UpdateInput>
): UseCrudMutationsReturn<T, ID, CreateInput, UpdateInput> {
  const { t } = useI18n()
  const queryClient = useQueryClient()
  
  const {
    queryKey,
    createFn,
    updateFn,
    deleteFn,
    messages = {},
    onSuccess,
    onError,
  } = config

  const {
    createSuccess = t('components.confirmDialog.createSuccess') || 'Created successfully',
    updateSuccess = t('components.confirmDialog.updateSuccess') || 'Updated successfully',
    deleteSuccess = t('components.confirmDialog.deleteSuccess') || 'Deleted successfully',
    createDescription = t('components.confirmDialog.createDescription') || 'The latest changes have been saved and the list will refresh shortly.',
    updateDescription = t('components.confirmDialog.updateDescription') || 'Your changes have been saved and the list will refresh shortly.',
    deleteDescription = t('components.confirmDialog.deleteDescription') || 'The item has been removed and the current list is now up to date.',
    createLoading = t('components.confirmDialog.createLoading') || 'Creating...',
    updateLoading = t('components.confirmDialog.updateLoading') || 'Saving changes...',
    deleteLoading = t('components.confirmDialog.deleteLoading') || 'Deleting...',
  } = messages

  // 默认错误处理
  const handleError = (error: Error, action: 'create' | 'update' | 'delete') => {
    if (onError) {
      onError(error, action)
    }
  }

  // 默认成功处理
  const handleSuccess = async () => {
    // 刷新列表数据
    await queryClient.invalidateQueries({ queryKey })
    if (onSuccess) {
      onSuccess()
    }
  }

  const createMutation = useMutation({
    mutationFn: createFn,
    onMutate: () => {
      return { toastId: toast.loading(createLoading) }
    },
    onSuccess: async (_, __, context) => {
      if (context?.toastId) toast.dismiss(context.toastId)
      toast.success(createSuccess, { description: createDescription })
      await handleSuccess()
    },
    onError: (error, _variables, context) => {
      if (context?.toastId) toast.dismiss(context.toastId)
      handleError(error, 'create')
    },
  })

  const updateMutation = useMutation({
    mutationFn: updateFn,
    onMutate: () => {
      return { toastId: toast.loading(updateLoading) }
    },
    onSuccess: async (_, __, context) => {
      if (context?.toastId) toast.dismiss(context.toastId)
      toast.success(updateSuccess, { description: updateDescription })
      await handleSuccess()
    },
    onError: (error, _variables, context) => {
      if (context?.toastId) toast.dismiss(context.toastId)
      handleError(error, 'update')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: deleteFn,
    onMutate: () => {
      return { toastId: toast.loading(deleteLoading) }
    },
    onSuccess: async (_, __, context) => {
      if (context?.toastId) toast.dismiss(context.toastId)
      toast.success(deleteSuccess, { description: deleteDescription })
      await handleSuccess()
    },
    onError: (error, _variables, context) => {
      if (context?.toastId) toast.dismiss(context.toastId)
      handleError(error, 'delete')
    },
  })

  const isLoading = createMutation.isPending || updateMutation.isPending || deleteMutation.isPending

  return {
    createMutation,
    updateMutation,
    deleteMutation,
    isLoading,
  }
}
