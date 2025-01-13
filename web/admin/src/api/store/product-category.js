import request from '@/utils/request'

// 获取菜单树
export function getCategoryTree() {
  return request({
    url: '/api/product/category/tree',
    method: 'get'
  })
}

// 创建菜单
export function createCategory(data) {
  return request({
    url: '/api/product/category',
    method: 'post',
    data
  })
}

// 更新菜单
export function updateCategory(Id, data) {
  return request({
    url: '/api/product/category/' + Id,
    method: 'patch',
    data
  })
}

// 批量删除菜单
export function batchDeleteCategory(data) {
  return request({
    url: '/api/product/category',
    method: 'delete',
    data
  })
}

