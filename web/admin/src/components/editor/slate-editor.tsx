import { useCallback, useState, useEffect, useRef, useMemo } from 'react'
import isHotkey from 'is-hotkey'
import {
  createEditor,
  type Descendant,
  Editor,
  Element as SlateElement,
  Transforms,
  Range as SlateRange,
  type NodeEntry,
  type Path,
  Text,
} from 'slate'
import { withHistory } from 'slate-history'
import {
  Editable,
  ReactEditor,
  type RenderElementProps,
  type RenderLeafProps,
  Slate,
  useSlate,
  useSelected,
  useFocused,
  withReact,
} from 'slate-react'
import katex from 'katex'
import DOMPurify from 'dompurify'
import {
  Bold,
  Italic,
  Code,
  Quote,
  List,
  ListOrdered,
  Heading1,
  ChevronDown,
  Underline,
  Strikethrough,
  SquareCode,
  Image as ImageIcon,
  Table as TableIcon,
  Minus,
  CheckSquare,
  AlignLeft,
  AlignCenter,
  AlignRight,
  AlignJustify,
  Link as LinkIcon,
  Unlink,
  Undo,
  Redo,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ResourceUpload, FILE_TYPE, type ResourceItem } from '@/components/resource-upload'
import { markdownToSlate, slateToMarkdown } from '@/lib/slate-markdown'
import { useI18n } from '@/context/i18n-provider'
import { cn } from '@/lib/utils'

const I18N_PREFIX = 'components.editor'

/** 允许的 URL 协议白名单 */
const ALLOWED_PROTOCOLS = ['http:', 'https:', 'mailto:', 'tel:']

/** 验证 URL 是否安全（防止 javascript: 等协议注入） */
function isSafeUrl(url: string | undefined): boolean {
  if (!url) return false
  try {
    const parsed = new URL(url)
    return ALLOWED_PROTOCOLS.includes(parsed.protocol)
  } catch {
    // 相对路径允许（不包含协议）
    return !url.includes(':')
  }
}

const MARK_HOTKEYS: Record<string, string> = {
  'mod+b': 'bold',
  'mod+i': 'italic',
  'mod+u': 'underline',
  'mod+shift+x': 'strikethrough',
  'mod+`': 'code',
}

const HEADING_TYPES = ['heading-one', 'heading-two', 'heading-three', 'heading-four', 'heading-five', 'heading-six'] as const
const HEADING_HOTKEYS: Record<string, (typeof HEADING_TYPES)[number]> = {
  'mod+alt+1': 'heading-one',
  'mod+alt+2': 'heading-two',
  'mod+alt+3': 'heading-three',
  'mod+alt+4': 'heading-four',
  'mod+alt+5': 'heading-five',
  'mod+alt+6': 'heading-six',
}

const initialSlateValue: Descendant[] = [{ type: 'paragraph', children: [{ text: '' }] }]

/** 判断 Slate 节点是否为「空块」（空 paragraph / 空 heading / 仅含空白字符） */
function isEmptyBlock(node: Descendant): boolean {
  if (Text.isText(node)) {
    return node.text.trim() === ''
  }
  const el = node as SlateElement
  // 只清洗 paragraph 与 heading 类型的前导空块；保留 list、quote、image 等结构块
  if (el.type === 'paragraph' || (typeof el.type === 'string' && el.type.startsWith('heading-'))) {
    if (!Array.isArray(el.children) || el.children.length === 0) return true
    return el.children.every(isEmptyBlock)
  }
  return false
}

/** 删除节点数组最前面的连续空块；若全部为空则保留一个默认空 paragraph */
function trimLeadingEmptyBlocks(nodes: Descendant[]): Descendant[] {
  let i = 0
  while (i < nodes.length && isEmptyBlock(nodes[i])) {
    i++
  }
  const trimmed = nodes.slice(i)
  return trimmed.length > 0 ? trimmed : initialSlateValue
}

export type SlateEditorProps = {
  value: string
  onChange: (value: string) => void
  placeholder?: string
  className?: string
  minHeight?: string
  outputMode?: 'markdown' | 'json'
  autoHeight?: boolean
}

export function SlateEditor({
  value,
  onChange,
  placeholder = '',
  className,
  minHeight = 'min-h-[280px]',
  outputMode = 'markdown',
  autoHeight = false,
}: SlateEditorProps) {
  const { t } = useI18n()
  const editor = useMemo(() => {
    const baseEditor = withHistory(withReact(createEditor()))
    const { normalizeNode, isInline, isVoid } = baseEditor
    baseEditor.isInline = (element) => {
      return ['link', 'inline-math'].includes(element.type) || isInline(element)
    }
    baseEditor.isVoid = (element) => {
      return ['image', 'divider', 'thematic-break'].includes(element.type) || isVoid(element)
    }
    baseEditor.normalizeNode = (entry: NodeEntry, options?: Parameters<typeof normalizeNode>[1]) => {
      const [node] = entry
      if (SlateElement.isElement(node) && node.children.length === 0) {
        Transforms.insertNodes(baseEditor, { text: '' }, { at: [...entry[1], 0] })
        return
      }
      normalizeNode(entry, options)
    }
    return baseEditor
  }, [])

  const parseValue = useCallback(
    (val: string) => {
      if (!val?.trim()) return initialSlateValue
      if (outputMode === 'markdown') {
        return trimLeadingEmptyBlocks(markdownToSlate(val))
      }
      try {
        const parsed = JSON.parse(val)
        const nodes = Array.isArray(parsed) && parsed.length > 0 ? parsed : initialSlateValue
        return trimLeadingEmptyBlocks(nodes)
      } catch {
        return initialSlateValue
      }
    },
    [outputMode]
  )

  const [slateValue, setSlateValue] = useState<Descendant[]>(() => {
    const v = parseValue(value)
    return Array.isArray(v) && v.length > 0 ? v : initialSlateValue
  })
  const isInternalChange = useRef(false)

  // Update slateValue when value prop changes externally
  useEffect(() => {
    if (isInternalChange.current) {
      isInternalChange.current = false
      return
    }
    const next = parseValue(value)
    const safe = Array.isArray(next) && next.length > 0 ? next : initialSlateValue

    if (editor) {
      // 使用 Slate API 而非直接修改 editor.children，避免破坏内部状态
      editor.withoutNormalizing(() => {
        // 先选中所有内容
        Transforms.select(editor, {
          anchor: Editor.start(editor, []),
          focus: Editor.end(editor, []),
        })
        // 删除所有内容
        Transforms.delete(editor)
        // 插入新内容
        Transforms.insertNodes(editor, safe)
      })
      // Reset selection to avoid "Cannot get the leaf node" error if path becomes invalid
      editor.selection = null
    }
    setSlateValue(safe)
  }, [value, editor, parseValue])

  const handleSlateChange = useCallback(
    (next: Descendant[]) => {
      setSlateValue(next)
      isInternalChange.current = true

      try {
        let output = ''
        if (outputMode === 'markdown') {
          output = slateToMarkdown(next)
        } else {
          output = JSON.stringify(next)
        }
        onChange(output ?? '')
      } catch {
        onChange('')
      }
    },
    [onChange, outputMode]
  )

  const renderElement = useCallback((props: RenderElementProps) => <Element {...props} />, [])
  const renderLeaf = useCallback((props: RenderLeafProps) => <Leaf {...props} />, [])

  const handleKeyDown = useCallback(
    (event: React.KeyboardEvent) => {
      const { selection } = editor
      if (!selection) return

      // 1. Markdown 行首转换：Space 时 # / ## ... -> 标题；- * -> 列表；[] [ ] [x] -> 任务列表
      if (event.key === ' ') {
        const blockEntry = Editor.above(editor, { match: (n) => SlateElement.isElement(n) })
        if (blockEntry) {
          const [, path] = blockEntry
          const blockText = Editor.string(editor, path)
          const headingOrList = /^(#{1,6})$|^(-|\*)$|^(\d+\.)$/.exec(blockText)
          const checklistMatch = /^\[\s?([xX])?\s?\]$/.exec(blockText)
          if (checklistMatch) {
            event.preventDefault()
            const range = { anchor: { path, offset: 0 }, focus: { path, offset: blockText.length } }
            Transforms.delete(editor, { at: range })
            const [node] = Editor.node(editor, path)
            if (SlateElement.isElement(node) && node.children.length === 0) {
              Transforms.insertNodes(editor, { text: '' }, { at: [...path, 0] })
            }
            const checked = !!checklistMatch[1]
            setBlockType(editor, 'check-list', path)
            const listItemPath = [...path, 0]
            Transforms.setNodes(editor, { checked }, { at: listItemPath })
            Transforms.insertText(editor, ' ', { at: Editor.start(editor, path) })
            return
          }
          if (headingOrList) {
            event.preventDefault()
            // 不要删除整行，而是直接替换内容
            // 这样可以避免节点变空导致的 "Cannot get the leaf node" 错误
            // 先将当前 block 的所有子节点选中并替换为空文本（保留了结构）
            Transforms.select(editor, path)
            Transforms.delete(editor) // 删除内容，但不删除节点本身（如果它是 block）

            // 此时节点内容为空，但 Slate 会保留至少一个空文本节点如果它是 block
            // 但为了安全起见，我们不依赖删除后的状态，而是直接设置属性

            if (headingOrList[1]) {
              const depth = headingOrList[1].length
              Transforms.setNodes(editor, { type: HEADING_TYPES[depth - 1] as string }, { at: path })
              // 不需要额外插入空格，因为删除后光标在开头，用户刚按了空格
              // 但通常 Markdown 转换后用户希望光标在后面
              // 不过这里我们的逻辑是把 "# " 变成了标题样式，内容被清空了？
              // 不对，用户输入 "#" 然后按空格，意思是把当前行变成 H1，并清空 "#" 字符
              // 所以上面的逻辑是对的：清空文本，改变样式
            } else if (headingOrList[2]) {
              setBlockType(editor, 'bulleted-list', path)
            } else if (headingOrList[3]) {
              setBlockType(editor, 'numbered-list', path)
            }
            return
          }
        }
      }

      // 2. 标题快捷键 mod+alt+1..6（避免与 mod+b 等冲突）
      for (const hotkey in HEADING_HOTKEYS) {
        if (isHotkey(hotkey, event.nativeEvent)) {
          event.preventDefault()
          setBlockType(editor, HEADING_HOTKEYS[hotkey as keyof typeof HEADING_HOTKEYS])
          return
        }
      }

      // 3. 格式快捷键 mod+b / mod+i / mod+u / mod+shift+x / mod+`
      for (const hotkey in MARK_HOTKEYS) {
        if (isHotkey(hotkey, event.nativeEvent)) {
          event.preventDefault()
          const mark = MARK_HOTKEYS[hotkey] as 'bold' | 'italic' | 'underline' | 'strikethrough' | 'code'
          toggleMark(editor, mark)
          return
        }
      }

      // 4. Soft break / Exit block
      // 回车键：如果当前 block 为空且不是 paragraph，则转换为 paragraph 并跳出
      if (event.key === 'Enter') {
        const { selection } = editor
        if (selection && SlateRange.isCollapsed(selection)) {
          const blockEntry = Editor.above(editor, { match: (n) => SlateElement.isElement(n) })
          if (blockEntry) {
            const [block, path] = blockEntry
            const isBlockEmpty = Editor.string(editor, path) === ''
            const type = (block as SlateElement).type

            // 如果是在列表项中，且内容为空
            if (['list-item', 'check-list-item'].includes(type) && isBlockEmpty) {
              event.preventDefault()
              unwrapLists(editor)
              return
            }

            // 如果是在标题中，且内容为空
            if (type.startsWith('heading-') && isBlockEmpty) {
              event.preventDefault()
              setBlockType(editor, 'paragraph')
              return
            }

            // 如果是引用或代码块，且内容为空（通常引用块回车应该换行，但如果是空的连续回车可能想跳出）
            // 这里我们简单处理：如果是 block-quote 且空，跳出
            if (type === 'block-quote' && isBlockEmpty) {
              event.preventDefault()
              setBlockType(editor, 'paragraph')
              return
            }
          }
        }
      }

      // Backspace 键：如果当前 block 为空且不是 paragraph，则转换为 paragraph
      if (event.key === 'Backspace') {
        const { selection } = editor
        if (selection && SlateRange.isCollapsed(selection)) {
          const blockEntry = Editor.above(editor, { match: (n) => SlateElement.isElement(n) })
          if (blockEntry) {
            const [block, path] = blockEntry
            const isBlockEmpty = Editor.string(editor, path) === ''
            const type = (block as SlateElement).type

            if (type !== 'paragraph' && isBlockEmpty) {
              // 比如是 heading 或 list-item，按删除键应该变回 paragraph
              event.preventDefault()
              if (['list-item', 'check-list-item'].includes(type)) {
                 unwrapLists(editor)
              } else {
                 setBlockType(editor, 'paragraph')
              }
              return
            }
          }
        }
      }
      // 5. Select All (Cmd+A / Ctrl+A)
      if (isHotkey('mod+a', event.nativeEvent)) {
        event.preventDefault()
        Transforms.select(editor, {
          anchor: Editor.start(editor, []),
          focus: Editor.end(editor, []),
        })
        return
      }
    },
    [editor]
  )

  const handlePaste = useCallback(
    (event: React.ClipboardEvent) => {
      // 1. 如果当前在代码块中，走默认粘贴（纯文本）
      const [match] = Editor.nodes(editor, {
        match: (n) => SlateElement.isElement(n) && n.type === 'code-block',
      })
      if (match) {
        return
      }

      // 2. 获取剪贴板文本
      const text = event.clipboardData.getData('text/plain')
      if (!text) return

      // 3. 简单的启发式检查：如果包含 Markdown 特殊字符或换行，尝试解析
      // 如果只是普通无格式单行文本，交给 Slate 默认处理（通常更好）
      const isMarkdownLike =
        text.includes('\n') || /[*_`[\]#>-]/.test(text) || /^(\d+\.)/.test(text)

      if (!isMarkdownLike) return

      // 4. 尝试解析 Markdown
      const fragment = markdownToSlate(text)

      // 如果解析结果有效且不仅仅是空文本
      if (Array.isArray(fragment) && fragment.length > 0) {
        event.preventDefault()

        // 优化：如果只有一个 paragraph，提取 children 插入以保持行内样式（如果可能）
        if (fragment.length === 1 && fragment[0].type === 'paragraph') {
          // 强制类型断言，因为 fragment[0] 是 Descendant，可能是 Text，但 type==paragraph 意味着是 Element
          Transforms.insertFragment(editor, (fragment[0] as SlateElement).children)
        } else {
          Transforms.insertFragment(editor, fragment)
        }
      }
    },
    [editor]
  )

  return (
    <div
      className={cn(
        'rounded-md border border-input bg-background flex flex-col',
        autoHeight ? 'h-auto max-h-none overflow-visible' : 'min-h-[320px] h-[55vh] max-h-[55vh] overflow-hidden',
        minHeight,
        className
      )}
    >
      <Slate editor={editor} initialValue={slateValue} onChange={handleSlateChange}>
        <div className='shrink-0 border-b border-input bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60 sticky top-0 z-10'>
          <div className='flex flex-wrap items-center gap-1 pl-0 pr-3 py-2 min-h-[40px]'>
            <ToolbarGroup>
              <ToolbarItem label={t(`${I18N_PREFIX}.undo`)}>
                <UndoButton />
              </ToolbarItem>
              <ToolbarItem label={t(`${I18N_PREFIX}.redo`)}>
                <RedoButton />
              </ToolbarItem>
            </ToolbarGroup>
            <ToolbarSeparator />
            <ToolbarGroup>
              <ToolbarItem label={t(`${I18N_PREFIX}.bold`)}>
                <MarkButton format='bold' icon={<Bold className='h-4 w-4' />} />
              </ToolbarItem>
              <ToolbarItem label={t(`${I18N_PREFIX}.italic`)}>
                <MarkButton format='italic' icon={<Italic className='h-4 w-4' />} />
              </ToolbarItem>
              <ToolbarItem label={t(`${I18N_PREFIX}.underline`)}>
                <MarkButton format='underline' icon={<Underline className='h-4 w-4' />} />
              </ToolbarItem>
              <ToolbarItem label={t(`${I18N_PREFIX}.strikethrough`)}>
                <MarkButton format='strikethrough' icon={<Strikethrough className='h-4 w-4' />} />
              </ToolbarItem>
            </ToolbarGroup>
            <ToolbarSeparator />
            <ToolbarGroup>
              <ToolbarItem label={t(`${I18N_PREFIX}.code`)}>
                <MarkButton format='code' icon={<Code className='h-4 w-4' />} />
              </ToolbarItem>
              <HeadingDropdown />
            </ToolbarGroup>
            <ToolbarSeparator />
            <ToolbarGroup>
              <ToolbarItem label={t(`${I18N_PREFIX}.quote`)}>
                <BlockButton format='block-quote' icon={<Quote className='h-4 w-4' />} />
              </ToolbarItem>
              <ToolbarItem label={t(`${I18N_PREFIX}.codeBlock`)}>
                <BlockButton format='code-block' icon={<SquareCode className='h-4 w-4' />} />
              </ToolbarItem>
              <ToolbarItem label={t(`${I18N_PREFIX}.bulletList`)}>
                <BlockButton format='bulleted-list' icon={<List className='h-4 w-4' />} />
              </ToolbarItem>
              <ToolbarItem label={t(`${I18N_PREFIX}.orderedList`)}>
                <BlockButton format='numbered-list' icon={<ListOrdered className='h-4 w-4' />} />
              </ToolbarItem>
              <ToolbarItem label={t(`${I18N_PREFIX}.checklist`)}>
                <BlockButton format='check-list' icon={<CheckSquare className='h-4 w-4' />} />
              </ToolbarItem>
            </ToolbarGroup>
            <ToolbarSeparator />
            <ToolbarGroup>
              <ToolbarItem label={t(`${I18N_PREFIX}.insertImage`)}>
                <InsertImageButton />
              </ToolbarItem>
              <ToolbarItem label={t(`${I18N_PREFIX}.link`)}>
                <InsertLinkButton />
              </ToolbarItem>
              <ToolbarItem label={t(`${I18N_PREFIX}.table`)}>
                <InsertTableButton />
              </ToolbarItem>
              <ToolbarItem label={t(`${I18N_PREFIX}.divider`)}>
                <InsertDividerButton />
              </ToolbarItem>
            </ToolbarGroup>
            <ToolbarSeparator />
            <ToolbarGroup>
              <AlignmentDropdown />
            </ToolbarGroup>
          </div>
        </div>
        <div className={cn('flex-1 min-h-0', autoHeight ? 'overflow-visible' : 'overflow-auto')}>
          <Editable
            renderElement={renderElement}
            renderLeaf={renderLeaf}
            onKeyDown={handleKeyDown}
            onPaste={handlePaste}
            placeholder={placeholder}
            className={cn(
              'prose prose-sm dark:prose-invert max-w-none w-full pl-0 pr-3 py-3 outline-none empty:before:content-[attr(placeholder)] empty:before:text-muted-foreground block first:[&>*]:mt-0',
              minHeight
            )}
            spellCheck
          />
        </div>
      </Slate>
    </div>
  )
}

function toggleMark(editor: Editor, format: string) {
  const isActive = isMarkActive(editor, format)
  if (isActive) {
    Editor.removeMark(editor, format)
  } else {
    Editor.addMark(editor, format, true)
  }
}

function isMarkActive(editor: Editor, format: string) {
  try {
    const marks = Editor.marks(editor)
    return marks ? (marks as Record<string, unknown>)[format] === true : false
  } catch {
    return false
  }
}

function unwrapLists(editor: Editor) {
  Transforms.unwrapNodes(editor, {
    match: (n) =>
      SlateElement.isElement(n) &&
      ['numbered-list', 'bulleted-list', 'check-list'].includes((n as SlateElement).type),
    split: true,
  })
}

function setBlockType(editor: Editor, format: string, at?: Path) {
  unwrapLists(editor)
  const isList = ['numbered-list', 'bulleted-list', 'check-list'].includes(format)
  const targetPath = at ?? (Editor.above(editor, { match: (n) => SlateElement.isElement(n) })?.[1])
  if (targetPath === undefined) return
  const listItemType = format === 'check-list' ? 'check-list-item' : 'list-item'
  Transforms.setNodes<SlateElement>(
    editor,
    { type: isList ? listItemType : format, ...(format === 'check-list' ? { checked: false } : {}) },
    { at: targetPath }
  )
  if (isList) {
    const block = { type: format, children: [] }
    Transforms.wrapNodes(editor, block, { at: targetPath })
  }
}

function toggleBlock(editor: Editor, format: string) {
  const isActive = isBlockActive(editor, format)
  const isList = ['numbered-list', 'bulleted-list', 'check-list'].includes(format)
  const listItemType = format === 'check-list' ? 'check-list-item' : 'list-item'
  unwrapLists(editor)
  const newProperties: Partial<SlateElement> = {
    type: isActive ? 'paragraph' : isList ? listItemType : format,
    ...(format === 'check-list' && !isActive ? { checked: false } : {}),
  }
  Transforms.setNodes<SlateElement>(editor, newProperties)
  if (!isActive && isList) {
    const block = { type: format, children: [] }
    Transforms.wrapNodes(editor, block)
  }
}

function isBlockActive(editor: Editor, format: string) {
  try {
    const { selection } = editor
    if (!selection) return false
    const [match] = Array.from(
      Editor.nodes(editor, {
        at: Editor.unhangRange(editor, selection),
        match: (n) => SlateElement.isElement(n) && (n as SlateElement).type === format,
      })
    )
    return !!match
  } catch {
    return false
  }
}

function getBlockAlign(editor: Editor): 'left' | 'center' | 'right' | 'justify' | undefined {
  try {
    const { selection } = editor
    if (!selection) return undefined
    const [node] = Editor.nodes(editor, {
      at: Editor.unhangRange(editor, selection),
      match: (n) => SlateElement.isElement(n) && !Editor.isEditor(n),
    })
    if (!node) return undefined
    const el = node[0] as SlateElement & { align?: string }
    return el.align as 'left' | 'center' | 'right' | 'justify' | undefined
  } catch {
    return undefined
  }
}

function setBlockAlign(editor: Editor, align: 'left' | 'center' | 'right' | 'justify') {
  const { selection } = editor
  if (!selection) return
  const [match] = Array.from(
    Editor.nodes(editor, {
      at: selection,
      match: (n) =>
        SlateElement.isElement(n) &&
        ['paragraph', 'heading-one', 'heading-two', 'heading-three', 'heading-four', 'heading-five', 'heading-six'].includes((n as SlateElement).type),
    })
  )
  if (match) Transforms.setNodes(editor, { align }, { at: match[1] })
}

function insertImage(editor: Editor, url: string) {
  if (!url?.trim()) return
  const image = {
    type: 'image',
    url: url.trim(),
    width: 40,
    align: 'left',
    children: [{ text: '' }],
  }
  const trailingParagraph = {
    type: 'paragraph',
    children: [{ text: '' }],
  }
  Transforms.insertNodes(editor, [image, trailingParagraph])
  window.setTimeout(() => {
    ReactEditor.focus(editor)
    Transforms.move(editor)
  }, 0)
}

function isLinkActive(editor: Editor) {
  const [link] = Editor.nodes(editor, {
    match: (n) =>
      !Editor.isEditor(n) && SlateElement.isElement(n) && (n as SlateElement).type === 'link',
  })
  return !!link
}

function unwrapLink(editor: Editor) {
  Transforms.unwrapNodes(editor, {
    match: (n) =>
      !Editor.isEditor(n) && SlateElement.isElement(n) && (n as SlateElement).type === 'link',
  })
}

function wrapLink(editor: Editor, url: string) {
  if (isLinkActive(editor)) {
    unwrapLink(editor)
  }

  const { selection } = editor
  const isCollapsed = selection && SlateRange.isCollapsed(selection)
  const link: SlateElement = {
    type: 'link',
    url,
    children: isCollapsed ? [{ text: url }] : [],
  }

  if (isCollapsed) {
    Transforms.insertNodes(editor, link)
  } else {
    Transforms.wrapNodes(editor, link, { split: true })
    Transforms.collapse(editor, { edge: 'end' })
  }
}

function insertTable(editor: Editor, rows: number, cols: number) {
  const r = Math.max(1, Math.min(10, rows))
  const c = Math.max(1, Math.min(10, cols))
  const tableRows = Array.from({ length: r }, () => ({
    type: 'table-row',
    children: Array.from({ length: c }, () => ({
      type: 'table-cell',
      children: [{ type: 'paragraph', children: [{ text: '' }] }],
    })),
  })) as unknown as Descendant[]
  const table = { type: 'table', children: tableRows } as unknown as Descendant
  Transforms.insertNodes(editor, table)
}

function insertDivider(editor: Editor) {
  Transforms.insertNodes(editor, { type: 'divider', children: [{ text: '' }] })
}

const HEADING_CLASSES: Record<string, string> = {
  h1: 'text-2xl font-bold mt-4 mb-2',
  h2: 'text-xl font-semibold mt-3 mb-2',
  h3: 'text-lg font-semibold mt-3 mb-1',
  h4: 'text-base font-semibold mt-2 mb-1',
  h5: 'text-sm font-semibold mt-2 mb-1',
  h6: 'text-sm font-medium mt-2 mb-1 text-muted-foreground',
}

function blockStyle(el: SlateElement & { align?: string }) {
  const align = el.align
  if (align === 'center' || align === 'right' || align === 'justify')
    return { textAlign: align }
  return undefined
}

function MathElement({ attributes, children, element }: RenderElementProps) {
  const selected = useSelected()
  const focused = useFocused()
  const isInline = element.type === 'inline-math'
  const text = useMemo(() => {
    const node = element.children[0]
    return Text.isText(node) ? node.text : ''
  }, [element])

  const html = useMemo(() => {
    try {
      return katex.renderToString(text, {
        throwOnError: false,
        displayMode: !isInline,
      })
    } catch {
      return text
    }
  }, [text, isInline])

  const showSource = selected && focused

  if (isInline) {
    return (
      <span {...attributes} className={cn(showSource && 'bg-muted rounded px-1')}>
        <span
          contentEditable={false}
          className={cn(showSource ? 'hidden' : '')}
          dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(html) }}
        />
        <span
          className={cn(
            'font-mono text-sm',
            showSource ? '' : 'absolute opacity-0 pointer-events-none h-0 w-0 overflow-hidden'
          )}
        >
          {children}
        </span>
      </span>
    )
  }

  return (
    <div {...attributes} className='my-2 relative'>
      <div
        contentEditable={false}
        className={cn('select-none', showSource ? 'hidden' : 'block')}
        dangerouslySetInnerHTML={{ __html: html }}
      />
      <div
        className={cn(
          'font-mono text-sm bg-muted p-2 rounded',
          showSource ? 'block' : 'absolute opacity-0 pointer-events-none h-0 w-0 overflow-hidden'
        )}
      >
        {children}
      </div>
    </div>
  )
}

function Element({ attributes, children, element }: RenderElementProps) {
  const el = element as SlateElement & { depth?: number; align?: string; width?: number }
  if (el.type === 'heading' && typeof el.depth === 'number') {
    const d = Math.min(6, Math.max(1, el.depth))
    const Tag = `h${d}` as 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6'
    return (
      <Tag {...attributes} className={HEADING_CLASSES[Tag]} style={blockStyle(el)}>
        {children}
      </Tag>
    )
  }
  switch (element.type) {
    case 'block-quote':
    case 'blockquote':
      return (
        <blockquote {...attributes} className='border-l-4 border-primary pl-4 my-2'>
          {children}
        </blockquote>
      )
    case 'heading-one':
      return (
        <h1 {...attributes} className={HEADING_CLASSES.h1} style={blockStyle(el)}>
          {children}
        </h1>
      )
    case 'heading-two':
      return (
        <h2 {...attributes} className={HEADING_CLASSES.h2} style={blockStyle(el)}>
          {children}
        </h2>
      )
    case 'heading-three':
      return (
        <h3 {...attributes} className={HEADING_CLASSES.h3} style={blockStyle(el)}>
          {children}
        </h3>
      )
    case 'heading-four':
      return (
        <h4 {...attributes} className={HEADING_CLASSES.h4} style={blockStyle(el)}>
          {children}
        </h4>
      )
    case 'heading-five':
      return (
        <h5 {...attributes} className={HEADING_CLASSES.h5} style={blockStyle(el)}>
          {children}
        </h5>
      )
    case 'heading-six':
      return (
        <h6 {...attributes} className={HEADING_CLASSES.h6} style={blockStyle(el)}>
          {children}
        </h6>
      )
    case 'list-item':
      return (
        <li {...attributes} className='ml-1'>
          {children}
        </li>
      )
    case 'numbered-list':
      return (
        <ol {...attributes} className='my-2 pl-5'>
          {children}
        </ol>
      )
    case 'bulleted-list':
      return (
        <ul {...attributes} className='my-2 pl-5'>
          {children}
        </ul>
      )
    case 'check-list':
      return (
        <ul {...attributes} className='list-none my-2 space-y-1 pl-0'>
          {children}
        </ul>
      )
    case 'check-list-item': {
      const checked = !!(el as SlateElement & { checked?: boolean }).checked
      return (
        <li
          {...attributes}
          className='flex items-start gap-2 list-none my-0.5'
          data-checked={checked}
        >
          <CheckboxElement checked={checked} element={el} />
          <span className={cn('flex-1', checked && 'line-through text-muted-foreground')}>
            {children}
          </span>
        </li>
      )
    }
    case 'image': {
      return <ImageElement attributes={attributes} element={el}>{children}</ImageElement>
    }
    case 'link': {
      const linkUrl = (el as SlateElement & { url?: string }).url
      const safeUrl = isSafeUrl(linkUrl) ? linkUrl : '#'
      return (
        <a
          {...attributes}
          href={safeUrl}
          className='text-primary underline underline-offset-4 cursor-pointer'
          onClick={(e) => {
            if ((e.metaKey || e.ctrlKey) && isSafeUrl(linkUrl)) {
              window.open(linkUrl, '_blank', 'noopener,noreferrer')
            }
          }}
        >
          {children}
        </a>
      )
    }
    case 'table':
      return (
        <div {...attributes} className='my-3 overflow-x-auto'>
          <table className='w-full border-collapse border border-border'>
            <tbody>{children}</tbody>
          </table>
        </div>
      )
    case 'table-row':
      return <tr {...attributes}>{children}</tr>
    case 'table-cell':
      return (
        <td {...attributes} className='border border-border p-2 align-top'>
          {children}
        </td>
      )
    case 'divider':
    case 'thematic-break':
      return (
        <div {...attributes} className='my-3'>
          <hr className='border-border' />
          {children}
        </div>
      )
    case 'code-block':
    case 'code':
      return (
        <pre {...attributes} className='my-2 rounded-md bg-secondary p-3 overflow-x-auto border border-border text-secondary-foreground'>
          <code className='text-sm font-mono'>{children}</code>
        </pre>
      )
    case 'math':
    case 'inline-math':
      return <MathElement attributes={attributes} children={children} element={element} />
    default:
      return (
        <p {...attributes} className='my-1 min-h-6 cursor-text' style={blockStyle(el)}>
          {children}
        </p>
      )
  }
}

function CheckboxElement({
  checked,
  element,
}: {
  checked: boolean
  element: SlateElement
}) {
  const editor = useSlate()
  return (
    <span
      contentEditable={false}
      className='mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded border border-input'
      onMouseDown={(e) => e.preventDefault()}
      onClick={() => {
        const path = ReactEditor.findPath(editor, element)
        Transforms.setNodes(editor, { checked: !checked }, { at: path })
      }}
    >
      {checked ? (
        <span className='text-primary text-sm leading-none'>✓</span>
      ) : null}
    </span>
  )
}

function ImageElement({
  attributes,
  children,
  element,
}: {
  attributes: RenderElementProps['attributes']
  children: RenderElementProps['children']
  element: SlateElement & { url?: string; width?: number; align?: 'left' | 'center' | 'right' | 'justify' }
}) {
  const editor = useSlate()
  const selected = useSelected()
  const { t } = useI18n()
  const [controlsOpen, setControlsOpen] = useState(false)
  const url = element.url
  const width = Math.min(100, Math.max(20, element.width ?? 100))
  const align = element.align ?? 'left'
  const showControls = selected && controlsOpen

  useEffect(() => {
    if (!selected) {
      setControlsOpen(false)
    }
  }, [selected])

  const updateWidth = (nextWidth: number) => {
    const path = ReactEditor.findPath(editor, element)
    Transforms.setNodes(editor, { width: Math.min(100, Math.max(20, nextWidth)) }, { at: path })
  }

  const updateAlign = (nextAlign: 'left' | 'center' | 'right') => {
    const path = ReactEditor.findPath(editor, element)
    Transforms.setNodes(editor, { align: nextAlign }, { at: path })
  }

  return (
    <div {...attributes} className='my-3'>
      <div
        contentEditable={false}
        className={cn(
          'relative flex w-full max-w-full',
          align === 'center' && 'justify-center',
          align === 'right' && 'justify-end',
          align === 'left' && 'justify-start'
        )}
      >
        <div
          className={cn(
            'overflow-hidden rounded-xl bg-transparent p-2 transition-shadow',
            selected && 'ring-2 ring-primary/20 shadow-sm border bg-muted/30'
          )}
          style={{ width: `${width}%`, maxWidth: '100%' }}
          onMouseDown={(e) => {
            e.preventDefault()
            setControlsOpen(true)
            const path = ReactEditor.findPath(editor, element)
            Transforms.select(editor, path)
            ReactEditor.focus(editor)
          }}
        >
          {url ? (
            <img src={url} alt='' className='block h-auto max-w-full rounded-md object-contain' />
          ) : (
            <span className='text-muted-foreground text-sm'>[图片]</span>
          )}
        </div>

        {showControls ? (
          <div className='absolute right-2 top-2 z-10 flex max-w-[calc(100vw-4rem)] items-center gap-1.5 rounded-full border bg-background/95 px-2 py-1.5 shadow-lg backdrop-blur'>
            <span className='text-xs text-muted-foreground'>{width}%</span>
            <input
              type='range'
              min={20}
              max={100}
              step={5}
              value={width}
              className='h-1.5 w-24 cursor-pointer accent-primary'
              onMouseDown={(e) => e.stopPropagation()}
              onChange={(e) => updateWidth(Number(e.target.value))}
            />
            <span className='mx-1 h-5 w-px bg-border' />
            <div className='flex items-center gap-1'>
              {[
                { value: 'left' as const, icon: <AlignLeft className='h-4 w-4' /> },
                { value: 'center' as const, icon: <AlignCenter className='h-4 w-4' /> },
                { value: 'right' as const, icon: <AlignRight className='h-4 w-4' /> },
              ].map((option) => (
                <Button
                  key={option.value}
                  type='button'
                  size='icon'
                  variant={align === option.value ? 'default' : 'outline'}
                  className='h-7 w-7 shadow-none'
                  onMouseDown={(e) => e.preventDefault()}
                  onClick={() => updateAlign(option.value)}
                  title={t(`${I18N_PREFIX}.imageAlign${option.value.charAt(0).toUpperCase()}${option.value.slice(1)}`)}
                >
                  {option.icon}
                </Button>
              ))}
            </div>
          </div>
        ) : null}
      </div>
      {children}
    </div>
  )
}

function Leaf({ attributes, children, leaf }: RenderLeafProps) {
  if (leaf.bold) children = <strong>{children}</strong>
  if (leaf.italic) children = <em>{children}</em>
  if (leaf.underline) children = <u>{children}</u>
  if (leaf.strikethrough) children = <s>{children}</s>
  if (leaf.code) children = <code className='rounded bg-secondary px-1.5 py-0.5 font-mono text-sm border border-border text-secondary-foreground'>{children}</code>
  return <span {...attributes}>{children}</span>
}

function ToolbarGroup({ children }: { children: React.ReactNode }) {
  return <span className='inline-flex items-center gap-0.5'>{children}</span>
}

function ToolbarSeparator() {
  return <span className='mx-1.5 h-5 w-px bg-border' aria-hidden />
}

function ToolbarItem({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <span className='inline-flex items-center gap-1 rounded px-1 py-0.5' title={label}>
      {children}
    </span>
  )
}

function MarkButton({ format, icon }: { format: string; icon: React.ReactNode }) {
  const editor = useSlate()
  const isActive = isMarkActive(editor, format)
  return (
    <Button
      type='button'
      variant='ghost'
      size='icon'
      className={cn('h-8 w-8', isActive && 'bg-muted')}
      onMouseDown={(e) => e.preventDefault()}
      onClick={() => toggleMark(editor, format)}
    >
      {icon}
    </Button>
  )
}

function BlockButton({ format, icon }: { format: string; icon: React.ReactNode }) {
  const editor = useSlate()
  const isActive = isBlockActive(editor, format)
  return (
    <Button
      type='button'
      variant='ghost'
      size='icon'
      className={cn('h-8 w-8', isActive && 'bg-muted')}
      onMouseDown={(e) => e.preventDefault()}
      onClick={() => toggleBlock(editor, format)}
    >
      {icon}
    </Button>
  )
}

const EDITOR_H_KEYS = ['h1', 'h2', 'h3', 'h4', 'h5', 'h6'] as const

function InsertImageButton() {
  const editor = useSlate()
  const { t } = useI18n()
  const [open, setOpen] = useState(false)
  const handleInsert = (resource: ResourceItem) => {
    insertImage(editor, resource.url)
    setOpen(false)
  }
  return (
    <>
      <Button
        type='button'
        variant='ghost'
        size='icon'
        className='h-8 w-8'
        onMouseDown={(e) => e.preventDefault()}
        onClick={() => setOpen(true)}
      >
        <ImageIcon className='h-4 w-4' />
      </Button>
      <ResourceUpload
        open={open}
        onOpenChange={setOpen}
        onSelect={handleInsert}
        type={FILE_TYPE.IMAGE}
        title={t(`${I18N_PREFIX}.insertImage`)}
        description={t(`${I18N_PREFIX}.imageResourceDescription`)}
      />
    </>
  )
}

function InsertTableButton() {
  const editor = useSlate()
  const { t } = useI18n()
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button type='button' variant='ghost' size='icon' className='h-8 w-8' onMouseDown={(e) => e.preventDefault()}>
          <TableIcon className='h-4 w-4' />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align='start' className='min-w-32'>
        {[
          [2, 2],
          [3, 3],
          [4, 4],
          [5, 5],
        ].map(([r, c]) => (
          <DropdownMenuItem
            key={`${r}x${c}`}
            onMouseDown={(e) => {
              e.preventDefault()
              insertTable(editor, r, c)
            }}
          >
            {r} × {c} {t(`${I18N_PREFIX}.table`)}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

function InsertDividerButton() {
  const editor = useSlate()
  return (
    <Button
      type='button'
      variant='ghost'
      size='icon'
      className='h-8 w-8'
      onMouseDown={(e) => e.preventDefault()}
      onClick={() => insertDivider(editor)}
    >
      <Minus className='h-4 w-4' />
    </Button>
  )
}

const ALIGN_OPTIONS: { value: 'left' | 'center' | 'right' | 'justify'; icon: React.ReactNode }[] = [
  { value: 'left', icon: <AlignLeft className='h-4 w-4' /> },
  { value: 'center', icon: <AlignCenter className='h-4 w-4' /> },
  { value: 'right', icon: <AlignRight className='h-4 w-4' /> },
  { value: 'justify', icon: <AlignJustify className='h-4 w-4' /> },
]

function AlignmentDropdown() {
  const editor = useSlate()
  const { t } = useI18n()
  const currentAlign = getBlockAlign(editor)
  const labelKey = currentAlign ? `align${currentAlign.charAt(0).toUpperCase()}${currentAlign.slice(1)}` : 'align'
  const label = t(`${I18N_PREFIX}.${labelKey}`)
  const Icon = ALIGN_OPTIONS.find((o) => o.value === currentAlign)?.icon ?? <AlignLeft className='h-4 w-4' />
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          type='button'
          variant='ghost'
          size='sm'
          className='h-8 gap-1 px-2'
          onMouseDown={(e) => e.preventDefault()}
        >
          {Icon}
          <span className='text-xs hidden sm:inline'>{label}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align='start' className='min-w-32'>
        {ALIGN_OPTIONS.map((opt) => (
          <DropdownMenuItem
            key={opt.value}
            onMouseDown={(e) => {
              e.preventDefault()
              setBlockAlign(editor, opt.value)
            }}
          >
            {opt.icon}
            <span className='ml-2'>{t(`${I18N_PREFIX}.align${opt.value.charAt(0).toUpperCase()}${opt.value.slice(1)}`)}</span>
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

function HeadingDropdown() {
  const editor = useSlate()
  const { t } = useI18n()
  const headingIndex = HEADING_TYPES.findIndex((f) => isBlockActive(editor, f))
  const label =
    headingIndex >= 0
      ? t(`${I18N_PREFIX}.${EDITOR_H_KEYS[headingIndex]}`)
      : t(`${I18N_PREFIX}.heading`)
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          type='button'
          variant='ghost'
          size='sm'
          className={cn('h-8 gap-1 px-2', headingIndex >= 0 && 'bg-muted')}
          onMouseDown={(e) => e.preventDefault()}
        >
          <Heading1 className='h-4 w-4' />
          <span className='text-xs max-w-16 truncate'>{label}</span>
          <ChevronDown className='h-3 w-3 opacity-50' />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align='start' className='min-w-40'>
        {HEADING_TYPES.map((format, i) => (
          <DropdownMenuItem
            key={format}
            onMouseDown={(e) => {
              e.preventDefault()
              setBlockType(editor, format)
            }}
          >
            <span className={cn('font-medium', HEADING_CLASSES[`h${i + 1}`])}>
              {t(`${I18N_PREFIX}.${EDITOR_H_KEYS[i]}`)}
            </span>
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

function InsertLinkButton() {
  const editor = useSlate()
  const { t } = useI18n()
  const [open, setOpen] = useState(false)
  const [url, setUrl] = useState('')

  const isActive = isLinkActive(editor)

  const handleInsert = () => {
    if (isActive) {
      unwrapLink(editor)
      return
    }
    if (!url) return
    wrapLink(editor, url)
    setUrl('')
    setOpen(false)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          type='button'
          variant='ghost'
          size='icon'
          className={cn('h-8 w-8', isActive && 'bg-muted')}
          onMouseDown={(e) => e.preventDefault()}
          onClick={(e) => {
            if (isActive) {
              e.preventDefault()
              unwrapLink(editor)
            }
          }}
        >
          {isActive ? <Unlink className='h-4 w-4' /> : <LinkIcon className='h-4 w-4' />}
        </Button>
      </DialogTrigger>
      <DialogContent className='sm:max-w-md'>
        <DialogHeader>
          <DialogTitle>{t(`${I18N_PREFIX}.insertLink`)}</DialogTitle>
        </DialogHeader>
        <div className='grid gap-2 py-2'>
          <Label htmlFor='link-url'>{t(`${I18N_PREFIX}.linkUrl`)}</Label>
          <Input
            id='link-url'
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder='https://...'
            onKeyDown={(e) => {
               if (e.key === 'Enter') {
                 e.preventDefault()
                 handleInsert()
               }
            }}
          />
        </div>
        <DialogFooter>
          <Button type='button' variant='outline' onMouseDown={(e) => e.preventDefault()} onClick={() => setOpen(false)}>
            {t(`${I18N_PREFIX}.cancel`)}
          </Button>
          <Button type='button' onMouseDown={(e) => e.preventDefault()} onClick={handleInsert}>
            {t(`${I18N_PREFIX}.insert`)}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

function UndoButton() {
  const editor = useSlate()
  return (
    <Button
      type='button'
      variant='ghost'
      size='icon'
      className='h-8 w-8'
      onMouseDown={(e) => e.preventDefault()}
      onClick={() => editor.undo()}
    >
      <Undo className='h-4 w-4' />
    </Button>
  )
}

function RedoButton() {
  const editor = useSlate()
  return (
    <Button
      type='button'
      variant='ghost'
      size='icon'
      className='h-8 w-8'
      onMouseDown={(e) => e.preventDefault()}
      onClick={() => editor.redo()}
    >
      <Redo className='h-4 w-4' />
    </Button>
  )
}
