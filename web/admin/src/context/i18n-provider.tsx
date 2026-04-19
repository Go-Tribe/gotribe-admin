import { createContext, useContext, type ReactNode } from 'react'
import { useTranslation } from 'react-i18next'

type Language = 'en' | 'zh'

const LANGUAGE_WHITELIST: Language[] = ['en', 'zh']

function normalizeLanguage(lang: string): Language {
  const lower = lang.split('-')[0]?.toLowerCase() ?? ''
  return LANGUAGE_WHITELIST.includes(lower as Language) ? (lower as Language) : 'zh'
}

interface I18nContextType {
  language: Language
  setLanguage: (lang: Language) => void
  t: (key: string, options?: Record<string, unknown>) => string
}

const I18nContext = createContext<I18nContextType | null>(null)

export function I18nProvider({ children }: { children: ReactNode }) {
  const { i18n, t } = useTranslation()

  const setLanguage = (lang: Language) => {
    i18n.changeLanguage(lang)
    if (typeof window !== 'undefined') {
      localStorage.setItem('i18next_lng', lang)
    }
  }

  return (
    <I18nContext.Provider
      value={{
        language: normalizeLanguage(i18n.language ?? 'zh'),
        setLanguage,
        t,
      }}
    >
      {children}
    </I18nContext.Provider>
  )
}

// eslint-disable-next-line react-refresh/only-export-components
export function useI18n() {
  const context = useContext(I18nContext)
  if (!context) {
    throw new Error('useI18n must be used within I18nProvider')
  }
  return context
}
