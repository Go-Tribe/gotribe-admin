/**
 * 从系统配置接口获取配置并更新 HTML meta 标签
 * 如果接口返回了配置，则直接更新 HTML meta 标签
 */
export async function updateAppConfigFromApi(): Promise<void> {
  try {
    // 动态导入，避免循环依赖
    const { getConfig } = await import('@/features/system/service/config')
    const { updateHtmlMeta } = await import('@/lib/update-html-meta')

    const response = await getConfig()
    const systemConfig = response?.systemConfig

    if (systemConfig) {
      // 直接从 API 响应中获取配置并更新 HTML meta 标签
      // 如果系统配置接口返回了其他字段，也可以在这里更新
      // 例如：description, themeColor, lang 等
      updateHtmlMeta({
        title: systemConfig.title,
        icon: systemConfig.icon,
        // 如果接口返回了其他字段，可以在这里添加
        // description: systemConfig.description,
        // themeColor: systemConfig.themeColor,
        // lang: systemConfig.lang,
      })
    }
  } catch (error) {
    // 如果获取配置失败，静默处理，不影响应用启动
    // HTML 中已有默认值，所以即使 API 失败也不影响用户体验
    // eslint-disable-next-line no-console
    console.warn('Failed to load system config from API:', error)
  }
}
