export type Project = {
  projectID: string
  title: string // 项目名称
  description: string // 项目描述
  name: string // 项目别名
  metaTitle?: string // Meta标题
  metaDescription?: string // Meta描述
  keywords: string // Meta关键词
  domain: string // 项目域名
  postUrl: string // 内容链接
  icp: string // icp备案号
  author: string // 项目归属者
  baiduAnalytics: string // 第三方js
  favicon: string // 网站图标
  publicSecurity: string // 公安备案号
  navImage: string // Nav图标
  info?: string // 额外信息
  pushToken?: string
  createdAt?: string
  updatedAt?: string
}

export type ProjectListResponse = {
  projects: Project[]
  total: number
}
