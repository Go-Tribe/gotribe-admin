/**
 * 全局状态与 Badge 颜色映射配置（方案 B）
 * 语义 → Badge variant，颜色由主题（CSS 变量）控制，定制主题时改 theme 即可。
 */

export type BadgeVariant = 'default' | 'secondary' | 'destructive' | 'success' | 'warning' | 'outline'

/** 语义键 → Badge variant；positive=绿(启用等)，negative=中性灰(草稿/待处理/停用等)，danger=红 */
const SEMANTIC_TO_VARIANT: Record<'positive' | 'negative' | 'danger', BadgeVariant> = {
  positive: 'success',
  negative: 'secondary',
  danger: 'destructive',
}

/** 启用/禁用类：1=启用(正面)，0/2=禁用(负面) - admin(0/1), role, menu.status(1/2) */
const ENABLED_STATUS_MAP: Record<number, keyof typeof SEMANTIC_TO_VARIANT> = {
  0: 'negative',
  1: 'positive',
  2: 'negative',
}

/** 发布状态(广告)：1=未发布(负面)，2=已发布(正面) */
const PUBLISH_STATUS_MAP: Record<number, keyof typeof SEMANTIC_TO_VARIANT> = {
  1: 'negative',
  2: 'positive',
}

/** 审核状态(评论)：1=待审核(负面)，2=通过(正面) */
const APPROVAL_STATUS_MAP: Record<number, keyof typeof SEMANTIC_TO_VARIANT> = {
  1: 'negative',
  2: 'positive',
}

/** 文章发布状态：1=草稿(负面)，2=已发布(正面) */
const ARTICLE_PUBLISH_MAP: Record<number, keyof typeof SEMANTIC_TO_VARIANT> = {
  1: 'negative',
  2: 'positive',
}

/** 显示/隐藏(分类)：1=显示(正面)，2=隐藏(负面) */
const VISIBILITY_MAP: Record<number, keyof typeof SEMANTIC_TO_VARIANT> = {
  1: 'positive',
  2: 'negative',
}

/** 显示/隐藏取反(菜单)：1=隐藏(负面)，2=显示(正面)；或 1=不缓存(负面)，2=缓存(正面) */
const VISIBILITY_INVERSE_MAP: Record<number, keyof typeof SEMANTIC_TO_VARIANT> = {
  1: 'negative',
  2: 'positive',
}

export type StatusType =
  | 'enabledStatus'
  | 'publishStatus'
  | 'approvalStatus'
  | 'articlePublish'
  | 'visibility'
  | 'visibilityInverse'

const VALUE_MAPS: Record<StatusType, Record<number, keyof typeof SEMANTIC_TO_VARIANT>> = {
  enabledStatus: ENABLED_STATUS_MAP,
  publishStatus: PUBLISH_STATUS_MAP,
  approvalStatus: APPROVAL_STATUS_MAP,
  articlePublish: ARTICLE_PUBLISH_MAP,
  visibility: VISIBILITY_MAP,
  visibilityInverse: VISIBILITY_INVERSE_MAP,
}

/**
 * 根据状态类型和数值得到 Badge variant，用于列表等展示。
 * 未匹配时回退为 secondary。
 */
export function getStatusVariantByType(value: number, type: StatusType): BadgeVariant {
  const map = VALUE_MAPS[type]
  const semantic = map[value] ?? 'negative'
  return SEMANTIC_TO_VARIANT[semantic]
}

/**
 * HTTP 状态码 → Badge variant（操作日志等）。
 * 2xx 正面，3xx 中性/负面，4xx/5xx 危险，1xx 及其它视为负面。
 */
export function getHttpStatusVariant(code: number): BadgeVariant {
  if (code >= 200 && code < 300) return SEMANTIC_TO_VARIANT.positive
  if (code >= 300 && code < 400) return SEMANTIC_TO_VARIANT.negative
  if (code >= 400) return SEMANTIC_TO_VARIANT.danger
  return SEMANTIC_TO_VARIANT.negative
}
