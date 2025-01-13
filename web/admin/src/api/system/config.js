import request from '@/utils/request'

export function getConfig() {
  return request({
    url: '/api/base/config',
    method: 'get'
  })
}

export function updateConfig(data) {
  return request({
    url: '/api/system',
    method: 'patch',
    data
  })
}
