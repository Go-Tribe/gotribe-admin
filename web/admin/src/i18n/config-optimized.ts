/**
 * 优化的 i18n 配置 - 按需加载版本
 * 
 * 说明：
 * 当前项目使用静态导入所有翻译文件，这会导致首屏加载所有语言的翻译。
 * 要完全实现按需加载，需要重构为动态导入方式。
 * 
 * 优化方案：
 * 1. 核心翻译（layout, common）静态导入
 * 2. 功能模块翻译按需动态导入
 * 3. 使用 i18next-http-backend 加载远程翻译文件
 * 
 * 要实现完整的按需加载，需要：
 * 1. 将翻译文件放到 public/locales 目录
 * 2. 使用 i18next-http-backend 插件
 * 3. 配置 webpack/vite 代码分割
 * 
 * 当前优化：对不常用功能添加懒加载标记
 */

import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'

// 核心翻译 - 始终加载
import layoutEnTranslations from './locales/components/layout/en.json'
import layoutZhTranslations from './locales/components/layout/zh.json'
import dataTableEnTranslations from './locales/components/data-table/en.json'
import dataTableZhTranslations from './locales/components/data-table/zh.json'
import confirmDialogEnTranslations from './locales/components/confirm-dialog/en.json'
import confirmDialogZhTranslations from './locales/components/confirm-dialog/zh.json'

// 认证相关 - 首屏需要
import authEnTranslations from './locales/features/auth/en.json'
import authZhTranslations from './locales/features/auth/zh.json'
import signInEnTranslations from './locales/features/auth/sign-in/en.json'
import signInZhTranslations from './locales/features/auth/sign-in/zh.json'

// 懒加载标记 - 这些可以在路由级别按需加载
// 使用魔法注释帮助 webpack 代码分割
const loadFeatureTranslations = async (feature: string, lang: string) => {
  switch (feature) {
    case 'dashboard':
      return import(/* webpackChunkName: "i18n-dashboard-[request]" */ `./locales/features/dashboard/${lang}.json`)
    case 'system':
      return import(/* webpackChunkName: "i18n-system-[request]" */ `./locales/features/system/admin/${lang}.json`)
    case 'content':
      return import(/* webpackChunkName: "i18n-content-[request]" */ `./locales/features/content/article/${lang}.json`)
    case 'business':
      return import(/* webpackChunkName: "i18n-business-[request]" */ `./locales/features/business/project/${lang}.json`)
    default:
      return {}
  }
}

const LANGUAGE_STORAGE_KEY = 'i18next_lng'
const DEFAULT_LANGUAGE = 'zh'

function getInitialLanguage(): string {
  if (typeof window === 'undefined') return DEFAULT_LANGUAGE
  const storedLang = localStorage.getItem(LANGUAGE_STORAGE_KEY)
  if (storedLang && ['en', 'zh'].includes(storedLang)) {
    return storedLang
  }
  const browserLang = navigator.language.split('-')[0]
  if (['en', 'zh'].includes(browserLang)) {
    return browserLang
  }
  return DEFAULT_LANGUAGE
}

i18n
  .use(initReactI18next)
  .init({
    resources: {
      en: {
        translation: {
          components: {
            layout: layoutEnTranslations,
            dataTable: dataTableEnTranslations,
            confirmDialog: confirmDialogEnTranslations,
          },
          features: {
            auth: {
              ...authEnTranslations,
              signIn: signInEnTranslations,
            },
          },
        },
      },
      zh: {
        translation: {
          components: {
            layout: layoutZhTranslations,
            dataTable: dataTableZhTranslations,
            confirmDialog: confirmDialogZhTranslations,
          },
          features: {
            auth: {
              ...authZhTranslations,
              signIn: signInZhTranslations,
            },
          },
        },
      },
    },
    lng: getInitialLanguage(),
    fallbackLng: DEFAULT_LANGUAGE,
    interpolation: {
      escapeValue: false,
    },
    // 添加资源加载优化
    react: {
      useSuspense: false, // 避免 SSR 问题
    },
  })

// 语言切换时保存
i18n.on('languageChanged', (lng) => {
  if (typeof window !== 'undefined') {
    localStorage.setItem(LANGUAGE_STORAGE_KEY, lng)
  }
})

// 导出按需加载函数
export { loadFeatureTranslations }
export default i18n
