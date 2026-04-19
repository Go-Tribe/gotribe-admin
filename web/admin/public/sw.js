/**
 * Service Worker for Go-Tribe Admin
 * 
 * 功能：
 * 1. 缓存静态资源（JS、CSS、图片）
 * 2. 网络优先策略获取 API 数据
 * 3. 离线时返回缓存内容
 */

const CACHE_NAME = 'go-tribe-admin-v1'
const STATIC_ASSETS = [
  '/',
  '/index.html',
  '/images/favicon.svg',
]

// 安装时缓存静态资源
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.addAll(STATIC_ASSETS)
    })
  )
  self.skipWaiting()
})

// 激活时清理旧缓存
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames
          .filter((name) => name !== CACHE_NAME)
          .map((name) => caches.delete(name))
      )
    })
  )
  self.clients.claim()
})

// 拦截请求
self.addEventListener('fetch', (event) => {
  const { request } = event
  const url = new URL(request.url)

  // API 请求使用网络优先策略
  if (url.pathname.startsWith('/api/')) {
    event.respondWith(networkFirst(request))
    return
  }

  // 静态资源使用缓存优先策略
  if (isStaticAsset(request)) {
    event.respondWith(cacheFirst(request))
    return
  }

  // 其他请求使用网络优先
  event.respondWith(networkFirst(request))
})

// 缓存优先策略
async function cacheFirst(request) {
  const cache = await caches.open(CACHE_NAME)
  const cached = await cache.match(request)

  if (cached) {
    return cached
  }

  try {
    const response = await fetch(request)
    if (shouldCache(request, response)) {
      await cache.put(request, response.clone())
    }
    return response
  } catch (error) {
    return new Response('Offline', { status: 503 })
  }
}

// 网络优先策略
async function networkFirst(request) {
  const cache = await caches.open(CACHE_NAME)

  try {
    const networkResponse = await fetch(request)
    if (shouldCache(request, networkResponse)) {
      await cache.put(request, networkResponse.clone())
    }
    return networkResponse
  } catch (error) {
    const cached = await cache.match(request)
    if (cached) {
      return cached
    }
    throw error
  }
}

// 判断是否是静态资源
function isStaticAsset(request) {
  return request.destination === 'script' ||
    request.destination === 'style' ||
    request.destination === 'image' ||
    request.destination === 'font'
}

// 判断是否应该缓存
function shouldCache(request, response) {
  if (request.method !== 'GET') return false
  if (response.status !== 200) return false
  if (!response || response.type === 'error') return false
  return true
}
