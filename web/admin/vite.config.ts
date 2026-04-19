import path from 'path'
import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react-swc'
import tailwindcss from '@tailwindcss/vite'
import { tanstackRouter } from '@tanstack/router-plugin/vite'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // 加载对应模式的环境变量
  const env = loadEnv(mode, process.cwd(), '')
  const apiProxyTarget = env.VITE_API_PROXY_TARGET || 'http://127.0.0.1:8089'

  return {
    base: '/',
    plugins: [
      tanstackRouter({
        target: 'react',
        autoCodeSplitting: true,
      }),
      react(),
      tailwindcss(),
    ],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      port: 5173,
      host: true,
      proxy: {
        '/api': {
          target: apiProxyTarget,
          changeOrigin: true,
        },
      },
    },
    preview: {
      port: 5173,
      host: true,
      proxy: {
        '/api': {
          target: apiProxyTarget,
          changeOrigin: true,
        },
      },
    },
    build: {
      // 确保输出到 dist 目录
      outDir: 'dist',
      // 调整 chunk 大小警告限制为 1500KB
      chunkSizeWarningLimit: 1500,
      // 生成 SPA 的 HTML
      rollupOptions: {
        input: {
          main: path.resolve(__dirname, 'index.html'),
        },
        output: {
          // 手动分割 chunk，优化缓存和加载
          manualChunks: {
            // 将 React 相关库打包在一起
            'react-vendor': ['react', 'react-dom', '@tanstack/react-query'],
            // 将路由相关打包在一起
            'router-vendor': ['@tanstack/react-router'],
            // 将 UI 组件库打包在一起
            'ui-vendor': [
              '@radix-ui/react-icons',
              '@radix-ui/react-avatar',
              '@radix-ui/react-dialog',
              '@radix-ui/react-dropdown-menu',
              '@radix-ui/react-popover',
              '@radix-ui/react-select',
              '@radix-ui/react-tooltip',
              '@radix-ui/react-tabs',
              '@radix-ui/react-switch',
              '@radix-ui/react-separator',
            ],
            // 编辑器相关（体积较大）
            'editor': ['slate', 'slate-react', 'slate-history'],
          },
        },
      },
    },
  }
})
