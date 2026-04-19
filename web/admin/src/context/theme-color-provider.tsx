import { createContext, useContext, useEffect, useState } from 'react'
import { getCookie, setCookie, removeCookie } from '@/lib/cookies'

export type ThemeColorId = 'default' | 'violet' | 'sage' | 'rose' | 'mint'

const THEME_COLOR_COOKIE_NAME = 'vite-ui-theme-color'
const THEME_COLOR_COOKIE_MAX_AGE = 60 * 60 * 24 * 365 // 1 year

const DEFAULT_THEME_COLOR: ThemeColorId = 'default'

const VALID_THEME_COLOR_IDS: ThemeColorId[] = ['default', 'violet', 'sage', 'rose', 'mint']

/** CSS class applied to documentElement when using a preset (e.g. theme-violet) */
const THEME_COLOR_CLASS: Record<ThemeColorId, string | null> = {
  default: null,
  violet: 'theme-violet',
  sage: 'theme-sage',
  rose: 'theme-rose',
  mint: 'theme-mint',
}

type ThemeColorProviderState = {
  themeColor: ThemeColorId
  setThemeColor: (id: ThemeColorId) => void
  resetThemeColor: () => void
}

const initialState: ThemeColorProviderState = {
  themeColor: DEFAULT_THEME_COLOR,
  setThemeColor: () => null,
  resetThemeColor: () => null,
}

const ThemeColorContext = createContext<ThemeColorProviderState>(initialState)

type ThemeColorProviderProps = {
  children: React.ReactNode
  storageKey?: string
}

export function ThemeColorProvider({
  children,
  storageKey = THEME_COLOR_COOKIE_NAME,
}: ThemeColorProviderProps) {
  const [themeColor, _setThemeColor] = useState<ThemeColorId>(() => {
    const saved = getCookie(storageKey)
    if (VALID_THEME_COLOR_IDS.includes(saved as ThemeColorId)) return saved as ThemeColorId
    return DEFAULT_THEME_COLOR
  })

  useEffect(() => {
    const root = window.document.documentElement
    const toRemove = (Object.keys(THEME_COLOR_CLASS) as ThemeColorId[])
      .map((id) => THEME_COLOR_CLASS[id])
      .filter(Boolean) as string[]
    const toAdd = THEME_COLOR_CLASS[themeColor]
    toRemove.forEach((c) => root.classList.remove(c))
    if (toAdd) root.classList.add(toAdd)
  }, [themeColor])

  const setThemeColor = (id: ThemeColorId) => {
    setCookie(storageKey, id, THEME_COLOR_COOKIE_MAX_AGE)
    _setThemeColor(id)
  }

  const resetThemeColor = () => {
    removeCookie(storageKey)
    _setThemeColor(DEFAULT_THEME_COLOR)
  }

  return (
    <ThemeColorContext.Provider
      value={{
        themeColor,
        setThemeColor,
        resetThemeColor,
      }}
    >
      {children}
    </ThemeColorContext.Provider>
  )
}

export function useThemeColor() {
  const context = useContext(ThemeColorContext)
  if (!context) {
    throw new Error('useThemeColor must be used within ThemeColorProvider')
  }
  return context
}
