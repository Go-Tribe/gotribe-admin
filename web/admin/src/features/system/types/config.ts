export interface SystemConfig {
  systemConfigID: string
  title: string
  content: string
  logo: string
  icon: string
  footer: string
}

export interface ConfigResponse {
  systemConfig: SystemConfig
}
