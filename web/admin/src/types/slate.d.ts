import type { BaseEditor } from 'slate'
import type { ReactEditor } from 'slate-react'
import type { HistoryEditor } from 'slate-history'

/** 与 remark-slate-transformer 输出兼容的块级元素 */
export type CustomElement = {
  type: string
  children: (CustomElement | CustomText)[]
  /** 文本对齐：left | center | right | justify */
  align?: 'left' | 'center' | 'right' | 'justify'
  /** 图片 URL（type === 'image'） */
  url?: string
  /** 图片宽度百分比（type === 'image'） */
  width?: number
  /** 任务项是否勾选（type === 'check-list-item'） */
  checked?: boolean
  [key: string]: unknown
}

export type CustomText = {
  text: string
  bold?: boolean
  italic?: boolean
  underline?: boolean
  strikethrough?: boolean
  code?: boolean
  [key: string]: unknown
}

declare module 'slate' {
  interface CustomTypes {
    Editor: BaseEditor & ReactEditor & HistoryEditor
    Element: CustomElement
    Text: CustomText
  }
}
