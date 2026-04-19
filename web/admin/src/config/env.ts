/**
 * 环境变量配置和验证
 * 在应用启动时验证必需的环境变量，避免运行时错误
 */

interface EnvConfig {
  VITE_API_BASE_URL: string
}

/**
 * 验证环境变量
 * @throws {Error} 如果必需的环境变量缺失或无效
 */
function validateEnv(): EnvConfig {
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL

  if (!apiBaseUrl) {
    throw new Error(
      'Missing required environment variable: VITE_API_BASE_URL. ' +
        'Please check your .env file.'
    )
  }

  // 验证 URL 格式：相对路径（如 /api）或绝对 URL；相对路径不做 new URL() 校验
  if (!apiBaseUrl.startsWith('/') && !apiBaseUrl.startsWith('http')) {
    throw new Error(
      `Invalid VITE_API_BASE_URL format: ${apiBaseUrl}. ` +
        'Please provide a valid URL (e.g., http://localhost:8088 or https://api.example.com) or a relative path (e.g. /api)'
    )
  }

  try {
    if (apiBaseUrl.startsWith('http')) {
      new URL(apiBaseUrl)
    }
  } catch {
    throw new Error(
      `Invalid VITE_API_BASE_URL format: ${apiBaseUrl}. ` +
        'Please provide a valid URL (e.g., http://localhost:8088 or https://api.example.com)'
    )
  }

  return {
    VITE_API_BASE_URL: apiBaseUrl,
  }
}

/**
 * 已验证的环境变量配置
 * 在应用启动时验证，确保所有必需的环境变量都存在且有效
 */
export const env = validateEnv()
