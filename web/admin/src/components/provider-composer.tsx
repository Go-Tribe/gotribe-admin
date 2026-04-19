import { memo, type ReactNode } from 'react'

interface ProviderComposerProps {
  providers: Array<React.FC<{ children: ReactNode }> | [React.FC<{ children: ReactNode }>, Record<string, unknown>]>
  children: ReactNode
}

/**
 * Provider 组合器
 * 将多个 Provider 扁平化嵌套，避免代码深度缩进
 * 
 * @example
 * <ProviderComposer
 *   providers={[
 *     QueryClientProvider,
 *     I18nProvider,
 *     [ThemeProvider, { defaultTheme: 'dark' }],
 *   ]}
 * >
 *   <App />
 * </ProviderComposer>
 */
export const ProviderComposer = memo(function ProviderComposer({
  providers,
  children,
}: ProviderComposerProps) {
  return providers.reduceRight<ReactNode>((acc, provider) => {
    if (Array.isArray(provider)) {
      const [Component, props] = provider
      return <Component {...props}>{acc}</Component>
    }
    const Component = provider
    return <Component>{acc}</Component>
  }, children)
})
