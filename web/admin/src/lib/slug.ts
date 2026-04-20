/**
 * 从标题生成 URL 友好的 slug
 * - 英文/ASCII 标题：转小写，空格替换为连字符
 * - 中文/非 ASCII 标题：返回空字符串，建议用户手动输入
 */
export function generateSlug(title: string): string {
  const trimmed = title.trim()
  if (!trimmed) return ''

  // 纯 ASCII（英文、数字、空格、连字符、下划线）
  if (Array.from(trimmed).every((char) => char.charCodeAt(0) <= 0x7f)) {
    return trimmed
      .toLowerCase()
      .replace(/[^a-z0-9\s-]/g, '')
      .replace(/\s+/g, '-')
      .replace(/-+/g, '-')
      .replace(/^-|-$/g, '')
  }

  // 非 ASCII（中文等）：返回空字符串，由用户手动输入
  return ''
}

/** 校验 slug 格式是否合法 */
export function isValidSlug(slug: string): boolean {
  if (!slug) return true // 空 slug 视为合法（可选字段）
  return /^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(slug)
}
