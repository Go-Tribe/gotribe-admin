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
      chunkSizeWarningLimit: 900,
      // 生成 SPA 的 HTML
      rollupOptions: {
        input: {
          main: path.resolve(__dirname, 'index.html'),
        },
        output: {
          manualChunks(id) {
            const nodeModulesMarker = '/node_modules/'
            const nodeModulesIndex = id.lastIndexOf(nodeModulesMarker)
            if (nodeModulesIndex === -1) return undefined

            const nodeModule = id.slice(nodeModulesIndex + nodeModulesMarker.length)

            if (
              nodeModule.startsWith('react/') ||
              nodeModule.startsWith('react-dom/') ||
              nodeModule.startsWith('scheduler/') ||
              nodeModule.startsWith('@tanstack/react-query/')
            ) {
              return 'react-vendor'
            }

            if (nodeModule.startsWith('@tanstack/react-router/')) {
              return 'router-vendor'
            }

            if (
              nodeModule.startsWith('@radix-ui/') ||
              nodeModule.startsWith('cmdk/') ||
              nodeModule.startsWith('sonner/')
            ) {
              return 'ui-vendor'
            }

            if (
              nodeModule.startsWith('vanilla-jsoneditor/')
            ) {
              return 'json-editor-vendor'
            }

            if (
              nodeModule.startsWith('@codemirror/') ||
              nodeModule.startsWith('@lezer/')
            ) {
              return 'codemirror-vendor'
            }

            if (
              nodeModule.startsWith('ajv/') ||
              nodeModule.startsWith('json-source-map/')
            ) {
              return 'json-editor-vendor'
            }

            if (
              nodeModule.startsWith('@tanstack/react-table/') ||
              nodeModule.startsWith('@tanstack/react-virtual/')
            ) {
              return 'table-vendor'
            }

            return undefined
          },
        },
      },
    },
  }
})
