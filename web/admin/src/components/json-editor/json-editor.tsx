import { useEffect, useLayoutEffect, useRef, useState, useCallback } from 'react'
import { createJSONEditor, Mode, type Content, type MenuItem } from 'vanilla-jsoneditor'
import 'vanilla-jsoneditor/themes/jse-theme-dark.css'
import './json-editor.css'
import { cn } from '@/lib/utils'

const MENU_TEXTS_TO_HIDE = new Set([
  'table', 'sort', 'filter', 'search', 'undo', 'redo',
  '排序', '筛选', '搜索', '撤销', '重做',
])

const MENU_ICON_NAMES_TO_HIDE = new Set([
  'sort', 'filter', 'search', 'undo', 'redo',
  'arrow-up-wide-short', 'arrow-down-wide-short', 'magnifying-glass',
  'arrow-rotate-left', 'arrow-rotate-right', 'rotate', 'rotate-left', 'rotate-right',
])

function getButtonLabel(item: { type: string; text?: string; title?: string; icon?: { iconName?: string } }): string {
  const t = (item.text ?? item.title ?? '').trim().toLowerCase()
  if (t) return t
  const name = item.icon?.iconName?.toLowerCase()
  return name ?? ''
}

function shouldHideMenuItem(item: { type: string; text?: string; title?: string; icon?: { iconName?: string } }): boolean {
  if (item.type !== 'button') return false
  const label = getButtonLabel(item)
  if (MENU_TEXTS_TO_HIDE.has(label)) return true
  if (MENU_ICON_NAMES_TO_HIDE.has(label)) return true
  return false
}

function filterMenuItems<T extends { type: string; text?: string; title?: string; icon?: { iconName?: string } }>(items: T[]): T[] {
  return items.filter((item) => !shouldHideMenuItem(item)) as T[]
}

function valueToContent(value: string): Content {
  const trimmed = (value ?? '').trim()
  if (!trimmed) return { text: '' }
  try {
    const parsed = JSON.parse(trimmed)
    return { json: parsed }
  } catch {
    return { text: value ?? '' }
  }
}

function contentToString(content: Content): string {
  if (content && 'json' in content && content.json !== undefined) {
    try {
      return JSON.stringify(content.json, null, 2)
    } catch {
      return ''
    }
  }
  if (content && 'text' in content && typeof content.text === 'string') {
    return content.text
  }
  return ''
}

export type JsonEditorProps = {
  value: string
  onChange: (value: string) => void
  className?: string
  minHeight?: string
  readOnly?: boolean
  /** 用于编辑回显：传入后首次创建用此内容，避免 value 晚到导致空白 */
  initialContent?: string
}

export function JsonEditor({
  value,
  onChange,
  className,
  minHeight = 'min-h-[300px]',
  readOnly = false,
  initialContent,
}: JsonEditorProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const editorRef = useRef<ReturnType<typeof createJSONEditor> | null>(null)
  const onChangeRef = useRef(onChange)
  const initialContentRef = useRef(initialContent)
  const [containerReady, setContainerReady] = useState(false)

  // 使用 useEffect 更新 ref，避免在 render 期间更新
  useEffect(() => {
    onChangeRef.current = onChange
  }, [onChange])

  useEffect(() => {
    if (initialContent !== undefined) {
      initialContentRef.current = initialContent
    }
  }, [initialContent])

  useLayoutEffect(() => {
    if (containerRef.current && !containerReady) {
      // 使用 requestAnimationFrame 避免同步 setState
      requestAnimationFrame(() => setContainerReady(true))
    }
  }, [containerReady])

  // 将 onRenderMenu 提取为稳定的 useCallback，避免每次 updateProps 触发重建
  const onRenderMenu = useCallback(
    (items: MenuItem[]) => filterMenuItems(items) as MenuItem[],
    []
  )

  useEffect(() => {
    if (!containerReady || !containerRef.current) return
    const contentSource = (initialContentRef.current ?? value).trim()
    const content = contentSource ? valueToContent(contentSource) : { text: '' }
    editorRef.current = createJSONEditor({
      target: containerRef.current,
      props: {
        content,
        mode: Mode.text,
        onChange: (updatedContent: Content) => {
          onChangeRef.current(contentToString(updatedContent))
        },
        readOnly,
        onRenderMenu,
      },
    })
    return () => {
      editorRef.current?.destroy()
      editorRef.current = null
    }
  }, [containerReady, readOnly, onRenderMenu])

  useEffect(() => {
    if (!editorRef.current) return
    editorRef.current.updateProps({
      content: valueToContent(value),
      mode: Mode.text,
      readOnly,
      onRenderMenu,
    })
  }, [value, readOnly, onRenderMenu])

  return (
    <div
      className={cn('json-editor-wrapper rounded-md border overflow-hidden', minHeight, className)}
      data-slot='json-editor'
    >
      <div ref={containerRef} className={cn('h-full', minHeight)} />
    </div>
  )
}
