import { unified } from 'unified'
import remarkParse from 'remark-parse'
import remarkStringify from 'remark-stringify'
import remarkGfm from 'remark-gfm'
import remarkMath from 'remark-math'
import { remarkToSlate, slateToRemark } from 'remark-slate-transformer'
import { type Descendant, type Element, Node, Text } from 'slate'

/** remark-slate-transformer 的 overrides 使用具体节点形状，避免 any；库的 Plugin 类型未导出，传参处用断言兼容 */
const remarkToSlateOptions = {
  overrides: {
    math: (node: { value: string }) => ({ type: 'math', children: [{ text: node.value }] }),
    inlineMath: (node: { value: string }) => ({ type: 'inline-math', children: [{ text: node.value }] }),
    heading: (node: { depth?: number; children?: unknown[] }, next: (nodes: unknown) => unknown[]) => {
      const depth = Math.min(Math.max(1, node.depth ?? 1), 6)
      const types = ['heading-one', 'heading-two', 'heading-three', 'heading-four', 'heading-five', 'heading-six']
      return { type: types[depth - 1], children: next(node.children ?? []) }
    },
    list: (node: { children?: unknown[]; ordered?: boolean }, next: (nodes: unknown) => unknown[]) => {
      const children = node.children ?? []
      const hasChecked = children.some(
        (child: unknown) => typeof child === 'object' && child !== null && 'checked' in (child as object) && (child as { checked?: unknown }).checked != null
      )
      if (hasChecked) return { type: 'check-list', children: next(children) }
      return { type: node.ordered ? 'numbered-list' : 'bulleted-list', children: next(children) }
    },
    listItem: (node: { checked?: boolean | null; children?: unknown[] }, next: (nodes: unknown) => unknown[]) => {
      if (node.checked !== null && node.checked !== undefined) {
        return { type: 'check-list-item', checked: node.checked, children: next(node.children ?? []) }
      }
      return { type: 'list-item', children: next(node.children ?? []) }
    },
    blockquote: (node: { children?: unknown[] }, next: (nodes: unknown) => unknown[]) => ({ type: 'block-quote', children: next(node.children ?? []) }),
    code: (node: { value?: string }) => ({ type: 'code-block', children: [{ text: node.value ?? '' }] }),
    thematicBreak: () => ({ type: 'divider', children: [{ text: '' }] }),
    image: (node: { url?: string; width?: number }) => ({
      type: 'image',
      url: node.url ?? '',
      width: typeof node.width === 'number' ? node.width : 100,
      children: [{ text: '' }],
    }),
    link: (node: { url?: string; children?: unknown[] }, next: (nodes: unknown) => unknown[]) => ({ type: 'link', url: node.url ?? '', children: next(node.children ?? []) }),
    table: (node: { children?: unknown[] }, next: (nodes: unknown) => unknown[]) => ({ type: 'table', children: next(node.children ?? []) }),
    tableRow: (node: { children?: unknown[] }, next: (nodes: unknown) => unknown[]) => ({ type: 'table-row', children: next(node.children ?? []) }),
    tableCell: (node: { children?: unknown[] }, next: (nodes: unknown) => unknown[]) => ({
      type: 'table-cell',
      children: [{ type: 'paragraph', children: next(node.children ?? []) }],
    }),
    delete: (node: { children?: unknown[] }, next: (nodes: unknown) => unknown[]) => {
      const children = next(node.children ?? []) as { text?: string; strikethrough?: boolean }[]
      return children.map((child) => (child.text !== undefined ? { ...child, strikethrough: true } : child))
    },
  },
}

const markdownToSlateProcessor = unified()
  .use(remarkParse)
  .use(remarkGfm)
  .use(remarkMath)
  // 库类型为 [boolean] | [Options?]，overrides 与 Options 兼容但类型未导出，此处断言以通过类型检查
  .use(remarkToSlate, remarkToSlateOptions as never)

const markdownStringifyProcessor = unified()
  .use(remarkStringify)
  .use(remarkGfm)
  .use(remarkMath)

const EMPTY_SLATE: Descendant[] = [{ type: 'paragraph', children: [{ text: '' }] }]

/** 将 Markdown 字符串转为 Slate 文档（用于编辑器初始值） */
export function markdownToSlate(markdown: string): Descendant[] {
  if (typeof markdown !== 'string' || !markdown.trim()) {
    return EMPTY_SLATE
  }
  try {
    const result = markdownToSlateProcessor.processSync(markdown).result
    if (!Array.isArray(result) || result.length === 0) return EMPTY_SLATE

    // Validate and normalize top-level nodes
    const validated = result.map(validateNode)

    // Ensure all top-level nodes are blocks (Elements)
    const normalized = validated.map(node => {
      if (Text.isText(node)) {
        return { type: 'paragraph', children: [node] }
      }
      return node
    }) as Descendant[]

    return normalized
  } catch {
    return [{ type: 'paragraph', children: [{ text: markdown.slice(0, 500) }] }]
  }
}

/** 提取 Slate 文档的纯文本内容 */
export function slateToPlainText(nodes: Descendant[]): string {
  return nodes.map(node => Node.string(node)).join('\n')
}

function validateNode(node: unknown): Descendant {
  if (!node || typeof node !== 'object') {
    return { text: '' }
  }

  if (Text.isText(node)) {
    return node
  }

  const el = node as Element
  // Ensure element has children
  if (!el.children || !Array.isArray(el.children) || el.children.length === 0) {
    return { ...el, children: [{ text: '' }] }
  }

  // Recursively validate children
  const validatedChildren = el.children.map(validateNode)

  // Double check that we didn't end up with empty children after validation
  if (validatedChildren.length === 0) {
    return { ...el, children: [{ text: '' }] }
  }

  return {
    ...el,
    children: validatedChildren,
  }
}

/** 将工具栏产生的 Slate 类型转为 remark-slate-transformer / mdast 期望的格式后再序列化 */
function normalizeForRemark(nodes: Descendant[]): Descendant[] {
  return nodes.map((node) => {
    if (!node || typeof node !== 'object') return { type: 'paragraph', children: [{ text: '' }] }
    if (Text.isText(node)) return node
    const el = node as Element & { depth?: number; url?: string; width?: number; align?: string; checked?: boolean }
    const children = Array.isArray(el.children) ? normalizeForRemark(el.children as Descendant[]) : [{ text: '' }]
    switch (el.type) {
      case 'heading-one':
        return { type: 'heading', depth: 1, children }
      case 'heading-two':
        return { type: 'heading', depth: 2, children }
      case 'heading-three':
        return { type: 'heading', depth: 3, children }
      case 'heading-four':
        return { type: 'heading', depth: 4, children }
      case 'heading-five':
        return { type: 'heading', depth: 5, children }
      case 'heading-six':
        return { type: 'heading', depth: 6, children }
      case 'block-quote':
        return { type: 'blockquote', children }
      case 'code-block': {
        const value = (el.children as Descendant[])
          .map((c) => (Text.isText(c) ? c.text : ''))
          .join('')
        return { type: 'code', value, lang: null, meta: null }
      }
      case 'numbered-list':
        return { type: 'list', ordered: true, spread: false, children }
      case 'bulleted-list':
        return { type: 'list', ordered: false, spread: false, children }
      case 'list-item':
        return { type: 'listItem', spread: false, checked: null, children }
      case 'check-list':
        return { type: 'list', ordered: false, spread: false, children }
      case 'check-list-item':
        return { type: 'listItem', spread: false, checked: el.checked ?? null, children }
      case 'image':
        return { type: 'image', url: el.url ?? '', width: el.width ?? 100, align: el.align ?? 'left', alt: '', children: [] }
      case 'link':
        return { type: 'link', url: el.url ?? '', children }
      case 'table': {
        const rows = (el.children as Descendant[]).map((row) => {
          const rowEl = row as Element
          const cells = (rowEl.children as Descendant[]).map((cell) => ({
            type: 'tableCell',
            children: normalizeForRemark((cell as Element).children as Descendant[]),
          }))
          return { type: 'tableRow', children: cells }
        })
        return { type: 'table', align: null, children: rows }
      }
      case 'table-row': {
        const cells = (el.children as Descendant[]).map((cell) => ({
          type: 'tableCell',
          children: normalizeForRemark((cell as Element).children as Descendant[]),
        }))
        return { type: 'tableRow', children: cells }
      }
      case 'table-cell':
        return { type: 'tableCell', children }
      case 'divider':
      case 'thematic-break':
        return { type: 'thematicBreak', children: [] }
      case 'math':
        return { type: 'math', children }
      case 'inline-math':
        return { type: 'inline-math', children }
      default:
        return { ...el, children }
    }
  }) as Descendant[]
}

/** 将 Slate 文档转为 Markdown 字符串（用于提交/存储） */
export function slateToMarkdown(value: Descendant[]): string {
  if (!Array.isArray(value) || value.length === 0) return ''
  try {
    const normalized = normalizeForRemark(value)
    const ast = slateToRemark(normalized, {
      overrides: {
        math: (node: unknown) => {
          const n = node as { children: { text: string }[] }
          return { type: 'math', value: n.children[0]?.text || '' }
        },
        'inline-math': (node: unknown) => {
          const n = node as { children: { text: string }[] }
          return { type: 'inlineMath', value: n.children[0]?.text || '' }
        },
      },
    })
    const out = markdownStringifyProcessor.stringify(ast)
    return typeof out === 'string' ? out : ''
  } catch {
    return ''
  }
}

/** 转义 HTML 特殊字符 */
function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

/** 文本叶子节点：支持 bold/strong、italic/emphasis、underline、strikethrough、code/inlineCode */
function leafToHtml(leaf: { text: string; bold?: boolean; strong?: boolean; italic?: boolean; emphasis?: boolean; underline?: boolean; strikethrough?: boolean; code?: boolean; inlineCode?: boolean }): string {
  let t = escapeHtml(leaf.text)
  if (!t) return ''
  if (leaf.bold || leaf.strong) t = `<strong>${t}</strong>`
  if (leaf.italic || leaf.emphasis) t = `<em>${t}</em>`
  if (leaf.underline) t = `<u>${t}</u>`
  if (leaf.strikethrough) t = `<s>${t}</s>`
  if (leaf.code || leaf.inlineCode) t = `<code>${t}</code>`
  return t
}

/** 将 Slate 节点数组直接序列化为 HTML（不经过 Markdown，避免 remark 链失败） */
function slateNodesToHtml(nodes: Descendant[], inlineContext = false): string {
  if (!Array.isArray(nodes) || nodes.length === 0) return ''
  const parts: string[] = []
  for (const node of nodes) {
    if (Text.isText(node)) {
      parts.push(leafToHtml(node))
      continue
    }
    const el = node as Element & { type: string; url?: string; width?: number; align?: string; checked?: boolean }
    const children = Array.isArray(el.children)
      ? slateNodesToHtml(el.children as Descendant[], el.type === 'paragraph' || el.type === 'link')
      : ''
    switch (el.type) {
      case 'paragraph':
        parts.push(children ? `<p>${children}</p>` : '<p><br></p>')
        break
      case 'heading-one':
      case 'heading-two':
      case 'heading-three':
      case 'heading-four':
      case 'heading-five':
      case 'heading-six': {
        const headingLevel: Record<string, number> = {
          'heading-one': 1,
          'heading-two': 2,
          'heading-three': 3,
          'heading-four': 4,
          'heading-five': 5,
          'heading-six': 6,
        }
        const n = headingLevel[el.type] ?? 1
        const tag = `h${n}`
        parts.push(children ? `<${tag}>${children}</${tag}>` : `<${tag}><br></${tag}>`)
        break
      }
      case 'block-quote':
        parts.push(`<blockquote>${children}</blockquote>`)
        break
      case 'code-block': {
        const code = (el.children as Descendant[])
          .map((c) => (Text.isText(c) ? escapeHtml(c.text) : ''))
          .join('')
        parts.push(`<pre><code>${code || '\n'}</code></pre>`)
        break
      }
      case 'bulleted-list':
        parts.push(`<ul>${children}</ul>`)
        break
      case 'numbered-list':
        parts.push(`<ol>${children}</ol>`)
        break
      case 'list-item':
        parts.push(`<li>${children || '<br>'}</li>`)
        break
      case 'check-list':
        parts.push(`<ul data-type="taskList">${children}</ul>`)
        break
      case 'check-list-item': {
        const checked = el.checked ? ' data-checked="true"' : ''
        parts.push(`<li${checked}>${children || '<br>'}</li>`)
        break
      }
      case 'image': {
        const url = el.url ? escapeHtml(el.url) : ''
        const width = typeof el.width === 'number' ? Math.min(100, Math.max(20, el.width)) : 100
        const align = el.align === 'center' || el.align === 'right' ? el.align : 'left'
        const marginStyle =
          align === 'center'
            ? 'display:block;margin-left:auto;margin-right:auto;'
            : align === 'right'
              ? 'display:block;margin-left:auto;margin-right:0;'
              : 'display:block;margin-left:0;margin-right:auto;'
        const style = ` style="width:${width}%;max-width:100%;height:auto;${marginStyle}"`
        if (inlineContext) {
          parts.push(url ? `<img src="${url}" alt=""${style} />` : '[图片]')
        } else {
          parts.push(url ? `<p><img src="${url}" alt=""${style} /></p>` : '<p>[图片]</p>')
        }
        break
      }
      case 'link': {
        const href = el.url ? escapeHtml(el.url) : '#'
        parts.push(children ? `<a href="${href}">${children}</a>` : href)
        break
      }
      case 'table':
        parts.push(`<table><tbody>${children}</tbody></table>`)
        break
      case 'table-row':
        parts.push(`<tr>${children}</tr>`)
        break
      case 'table-cell':
        parts.push(`<td>${children || '<br>'}</td>`)
        break
      case 'divider':
      case 'thematic-break':
        parts.push('<hr />')
        break
      case 'math':
      case 'inline-math': {
        const latex = (el.children as Descendant[]).map((c) => (Text.isText(c) ? escapeHtml(c.text) : '')).join('')
        parts.push(latex ? `<span class="math" data-latex="${latex}">${latex}</span>` : '')
        break
      }
      default:
        parts.push(children ? `<p>${children}</p>` : '')
    }
  }
  return parts.join('')
}

/**
 * 将 Slate 编辑器的 content（JSON 字符串）转为 HTML 字符串，用于存 htmlContent。
 * content 保持原样用于回显，htmlContent 用于展示/SEO 等。
 * 直接走 Slate→HTML 序列化，避免 Slate→Markdown→HTML 链在复杂内容下失败。
 * 若 content 已是 HTML（如旧数据），则原样返回。
 */
export function slateContentToHtml(content: string): string {
  if (typeof content !== 'string' || !content.trim()) return ''
  const trimmed = content.trim()
  if (trimmed.startsWith('<')) return trimmed
  try {
    const parsed = JSON.parse(content) as unknown
    const nodes = Array.isArray(parsed) && parsed.length > 0 ? (parsed as Descendant[]) : []
    return slateNodesToHtml(nodes)
  } catch {
    return ''
  }
}
