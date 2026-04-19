/**
 * 更新 HTML meta 标签的工具函数
 * 用于在运行时动态更新页面标题、描述等元数据
 */
export function updateHtmlMeta(config: {
  title?: string
  description?: string
  themeColor?: string
  lang?: string
  icon?: string
}): void {
  if (config.title) {
    document.title = config.title
    const titleMeta = document.querySelector('meta[name="title"]')
    if (titleMeta) {
      titleMeta.setAttribute('content', config.title)
    }
  }

  if (config.description) {
    const descMeta = document.querySelector('meta[name="description"]')
    if (descMeta) {
      descMeta.setAttribute('content', config.description)
    }
  }

  if (config.themeColor) {
    const themeMeta = document.querySelector('meta[name="theme-color"]')
    if (themeMeta) {
      themeMeta.setAttribute('content', config.themeColor)
    }
  }

  if (config.lang) {
    document.documentElement.lang = config.lang
  }

  if (config.icon) {
    // 更新 favicon
    // 查找所有现有的 icon link 标签（可能有多个，如 favicon.svg, favicon.png 等）
    const existingIcons = document.querySelectorAll('link[rel="icon"]')

    // 如果存在现有的 icon，更新第一个的 href
    if (existingIcons.length > 0) {
      const firstIcon = existingIcons[0] as HTMLLinkElement
      firstIcon.href = config.icon
      // 可选：移除其他 icon 标签，只保留一个
      for (let i = 1; i < existingIcons.length; i++) {
        existingIcons[i].remove()
      }
    } else {
      // 如果不存在，创建一个新的 link 标签
      const iconLink = document.createElement('link')
      iconLink.rel = 'icon'
      iconLink.type = 'image/svg+xml'
      iconLink.href = config.icon
      document.head.appendChild(iconLink)
    }
  }
}
