import type { ReactNode } from 'react'
import { QueryClientProvider } from '@tanstack/react-query'
import { DirectionProvider } from '@/context/direction-provider'
import { FontProvider } from '@/context/font-provider'
import { ThemeProvider } from '@/context/theme-provider'
import { ThemeColorProvider } from '@/context/theme-color-provider'
import { I18nProvider } from '@/context/i18n-provider'
import { ProjectProvider } from '@/context/project-provider'
import { ProviderComposer } from '@/components/provider-composer'
import { queryClient } from '@/lib/query-client'

/**
 * 创建应用 Provider 栈
 * 使用 ProviderComposer 避免深层嵌套
 */
export function AppProviders({ children }: { children: ReactNode }) {
  return (
    <ProviderComposer
      providers={[
        // QueryClient 必须最外层，因为它被其他 hooks 使用
        [(props) => <QueryClientProvider client={queryClient} {...props} />, {}],
        I18nProvider,
        ProjectProvider,
        ThemeProvider,
        ThemeColorProvider,
        FontProvider,
        DirectionProvider,
      ]}
    >
      {children}
    </ProviderComposer>
  )
}
