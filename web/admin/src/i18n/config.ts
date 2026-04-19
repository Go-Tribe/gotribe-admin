import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'

import layoutEnTranslations from './locales/components/layout/en.json'
import layoutZhTranslations from './locales/components/layout/zh.json'
import confirmDialogEnTranslations from './locales/components/confirm-dialog/en.json'
import confirmDialogZhTranslations from './locales/components/confirm-dialog/zh.json'
import languageSwitcherEnTranslations from './locales/components/language-switcher/en.json'
import languageSwitcherZhTranslations from './locales/components/language-switcher/zh.json'
import dataTableEnTranslations from './locales/components/data-table/en.json'
import dataTableZhTranslations from './locales/components/data-table/zh.json'
import profileDropdownEnTranslations from './locales/components/profile-dropdown/en.json'
import profileDropdownZhTranslations from './locales/components/profile-dropdown/zh.json'
import comingSoonEnTranslations from './locales/components/coming-soon/en.json'
import comingSoonZhTranslations from './locales/components/coming-soon/zh.json'
import datePickerEnTranslations from './locales/components/date-picker/en.json'
import datePickerZhTranslations from './locales/components/date-picker/zh.json'
import authEnTranslations from './locales/features/auth/en.json'
import authZhTranslations from './locales/features/auth/zh.json'
import signInEnTranslations from './locales/features/auth/sign-in/en.json'
import signInZhTranslations from './locales/features/auth/sign-in/zh.json'
import signUpEnTranslations from './locales/features/auth/sign-up/en.json'
import signUpZhTranslations from './locales/features/auth/sign-up/zh.json'
import otpEnTranslations from './locales/features/auth/otp/en.json'
import otpZhTranslations from './locales/features/auth/otp/zh.json'
import systemAdminEnTranslations from './locales/features/system/admin/en.json'
import systemAdminZhTranslations from './locales/features/system/admin/zh.json'
import systemRoleEnTranslations from './locales/features/system/role/en.json'
import systemRoleZhTranslations from './locales/features/system/role/zh.json'
import systemMenuEnTranslations from './locales/features/system/menu/en.json'
import systemMenuZhTranslations from './locales/features/system/menu/zh.json'
import systemApiEnTranslations from './locales/features/system/api/en.json'
import systemApiZhTranslations from './locales/features/system/api/zh.json'
import systemConfigEnTranslations from './locales/features/system/config/en.json'
import systemConfigZhTranslations from './locales/features/system/config/zh.json'
import logEnTranslations from './locales/features/log/en.json'
import logZhTranslations from './locales/features/log/zh.json'
import businessProjectEnTranslations from './locales/features/business/project/en.json'
import businessProjectZhTranslations from './locales/features/business/project/zh.json'
import businessUserEnTranslations from './locales/features/business/user/en.json'
import businessUserZhTranslations from './locales/features/business/user/zh.json'
import contentTagEnTranslations from './locales/features/content/tag/en.json'
import contentTagZhTranslations from './locales/features/content/tag/zh.json'
import contentCategoryEnTranslations from './locales/features/content/category/en.json'
import contentCategoryZhTranslations from './locales/features/content/category/zh.json'
import contentArticleEnTranslations from './locales/features/content/article/en.json'
import contentArticleZhTranslations from './locales/features/content/article/zh.json'
import contentResourceEnTranslations from './locales/features/content/resource/en.json'
import contentResourceZhTranslations from './locales/features/content/resource/zh.json'
import contentConfigEnTranslations from './locales/features/content/config/en.json'
import contentConfigZhTranslations from './locales/features/content/config/zh.json'
import contentColumnEnTranslations from './locales/features/content/column/en.json'
import contentColumnZhTranslations from './locales/features/content/column/zh.json'
import editorEnTranslations from './locales/components/editor/en.json'
import editorZhTranslations from './locales/components/editor/zh.json'
import configDrawerEnTranslations from './locales/components/config-drawer/en.json'
import configDrawerZhTranslations from './locales/components/config-drawer/zh.json'
import OperationPointEn from './locales/features/operation/point/en.json'
import OperationPointZh from './locales/features/operation/point/zh.json'
import OperationSceneEn from './locales/features/operation/scene/en.json'
import OperationSceneZh from './locales/features/operation/scene/zh.json'
import OperationAdvertisingEn from './locales/features/operation/advertising/en.json'
import OperationAdvertisingZh from './locales/features/operation/advertising/zh.json'
import OperationCommentEn from './locales/features/operation/comment/en.json'
import OperationCommentZh from './locales/features/operation/comment/zh.json'
import settingsEnTranslations from './locales/features/settings/en.json'
import settingsZhTranslations from './locales/features/settings/zh.json'
import dashboardEnTranslations from './locales/features/dashboard/en.json'
import dashboardZhTranslations from './locales/features/dashboard/zh.json'

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

i18n
  .use(initReactI18next)
  .init({
    resources: {
      en: {
        translation: {
          components: {
            layout: layoutEnTranslations,
            confirmDialog: confirmDialogEnTranslations,
            languageSwitcher: languageSwitcherEnTranslations,
            dataTable: dataTableEnTranslations,
            profileDropdown: profileDropdownEnTranslations,
            comingSoon: comingSoonEnTranslations,
            datePicker: datePickerEnTranslations,
            editor: editorEnTranslations,
            configDrawer: configDrawerEnTranslations,
          },
          features: {
            auth: {
              ...authEnTranslations,
              signIn: signInEnTranslations,
              signUp: signUpEnTranslations,
              otp: otpEnTranslations,
            },
            system: {
              admin: systemAdminEnTranslations,
              role: systemRoleEnTranslations,
              menu: systemMenuEnTranslations,
              api: systemApiEnTranslations,
              config: systemConfigEnTranslations,
            },
            log: logEnTranslations,
            business: {
              project: businessProjectEnTranslations,
              user: businessUserEnTranslations,
            },
            content: {
              tag: contentTagEnTranslations,
              category: contentCategoryEnTranslations,
              article: contentArticleEnTranslations,
              resource: contentResourceEnTranslations,
              config: contentConfigEnTranslations,
              column: contentColumnEnTranslations,
            },
            operation: {
              point: OperationPointEn,
              scene: OperationSceneEn,
              advertising: OperationAdvertisingEn,
              comment: OperationCommentEn,
            },
            dashboard: dashboardEnTranslations,
            settings: settingsEnTranslations,
          },
        },
      },
      zh: {
        translation: {
          components: {
            layout: layoutZhTranslations,
            confirmDialog: confirmDialogZhTranslations,
            languageSwitcher: languageSwitcherZhTranslations,
            dataTable: dataTableZhTranslations,
            profileDropdown: profileDropdownZhTranslations,
            comingSoon: comingSoonZhTranslations,
            datePicker: datePickerZhTranslations,
            editor: editorZhTranslations,
            configDrawer: configDrawerZhTranslations,
          },
          features: {
            auth: {
              ...authZhTranslations,
              signIn: signInZhTranslations,
              signUp: signUpZhTranslations,
              otp: otpZhTranslations,
            },
            system: {
              admin: systemAdminZhTranslations,
              role: systemRoleZhTranslations,
              menu: systemMenuZhTranslations,
              api: systemApiZhTranslations,
              config: systemConfigZhTranslations,
            },
            log: logZhTranslations,
            business: {
              project: businessProjectZhTranslations,
              user: businessUserZhTranslations,
            },
            content: {
              tag: contentTagZhTranslations,
              category: contentCategoryZhTranslations,
              article: contentArticleZhTranslations,
              resource: contentResourceZhTranslations,
              config: contentConfigZhTranslations,
              column: contentColumnZhTranslations,
            },
            operation: {
              point: OperationPointZh,
              scene: OperationSceneZh,
              advertising: OperationAdvertisingZh,
              comment: OperationCommentZh,
            },
            dashboard: dashboardZhTranslations,
            settings: settingsZhTranslations,
          },
        },
      },
    },
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
