import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'

const LANGUAGE_STORAGE_KEY = 'i18next_lng'
const DEFAULT_LANGUAGE = 'zh'

// Get language from localStorage or browser, fallback to default
function getInitialLanguage(): string {
  if (typeof window === 'undefined') return DEFAULT_LANGUAGE

  // Try to get from localStorage first
  const storedLang = localStorage.getItem(LANGUAGE_STORAGE_KEY)
  if (storedLang && ['en', 'zh'].includes(storedLang)) {
    return storedLang
  }

  // Try to detect from browser
  const browserLang = navigator.language.split('-')[0]
  if (['en', 'zh'].includes(browserLang)) {
    return browserLang
  }

  return DEFAULT_LANGUAGE
}

/**
 * 将 kebab-case 文件名转换为 camelCase 键名
 * @example 'confirm-dialog' -> 'confirmDialog'
 */
function toCamelCase(str: string): string {
  return str.replace(/-([a-z])/g, (_, char) => char.toUpperCase())
}

/**
 * 自动加载所有翻译文件
 * 使用 import.meta.glob 避免手动导入，新增语言/模块时无需修改此文件
 */
const localeModules = import.meta.glob('./locales/**/*.json', { eager: true }) as Record<
  string,
  { default: Record<string, unknown> }
>

/**
 * 构建 i18n resources 对象
 * 根据文件路径自动映射到对应的命名空间
 * 路径格式: ./locales/{category}/{module}/{lang}.json
 *           ./locales/{category}/{module}/{submodule}/{lang}.json
 */
const resources: Record<string, { translation: Record<string, unknown> }> = {}

for (const [path, module] of Object.entries(localeModules)) {
  // 路径格式: ./locales/components/layout/en.json
  const match = path.match(/\.\/locales\/(.+)\/(.+)\.json$/)
  if (!match) continue

  const pathParts = match[1].split('/')
  const lang = match[2]

  if (!resources[lang]) {
    resources[lang] = { translation: {} }
  }

  const data = module.default

  // 构建嵌套对象路径
  let current: Record<string, unknown> = resources[lang].translation

  for (let i = 0; i < pathParts.length; i++) {
    const part = toCamelCase(pathParts[i])

    if (i === pathParts.length - 1) {
      // 最后一级：直接赋值翻译数据
      // 如果是子模块（如 auth/sign-in），需要将数据合并到父模块
      const parent = current as Record<string, Record<string, unknown>>
      if (parent[part] && typeof parent[part] === 'object') {
        parent[part] = { ...parent[part], ...data }
      } else {
        parent[part] = data
      }
    } else {
      // 中间级：创建嵌套对象
      if (!current[part] || typeof current[part] !== 'object') {
        current[part] = {}
      }
      current = current[part] as Record<string, unknown>
    }
  }
}

i18n
  .use(initReactI18next)
  .init({
    resources,
    lng: getInitialLanguage(),
    fallbackLng: DEFAULT_LANGUAGE,
    interpolation: {
      escapeValue: false, // React already escapes values
    },
  })

// Save language preference to localStorage when language changes
i18n.on('languageChanged', (lng) => {
  if (typeof window !== 'undefined') {
    localStorage.setItem(LANGUAGE_STORAGE_KEY, lng)
  }
})

export default i18n
