import { useState } from 'react'

/**
 * 用于确认类弹窗的状态：传入某值时若当前已是该值则关闭，否则打开为该值（切换语义）。
 * @param initialState 初始值，通常为 null 表示关闭
 * @returns [当前值, setOpen]；setOpen(value) 会切换：若当前 === value 则置为 null，否则置为 value
 * @example const [open, setOpen] = useDialogState<"approve" | "reject">()
 */
export default function useDialogState<T extends string | boolean>(
  initialState: T | null = null
) {
  const [open, _setOpen] = useState<T | null>(initialState)

  const setOpen = (str: T | null) =>
    _setOpen((prev) => (prev === str ? null : str))

  return [open, setOpen] as const
}
