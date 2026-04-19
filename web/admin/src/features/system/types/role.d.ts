// {
//   "id": 0,
//   "desc": "ad",
//   "category": "ad",
//   "children": [
//       {
//           "id": 84,
//           "createdAt": "2025-01-03T22:01:27.346+08:00",
//           "updatedAt": "2025-01-03T22:01:27.346+08:00",
//           "deletedAt": null,
//           "method": "DELETE",
//           "path": "/ad",
//           "category": "ad",
//           "desc": "删除广告",
//           "creator": "系统"
//       },
//       {
//           "id": 83,
//           "createdAt": "2025-01-03T22:01:27.346+08:00",
//           "updatedAt": "2025-01-03T22:01:27.346+08:00",
//           "deletedAt": null,
//           "method": "PATCH",
//           "path": "/ad/:adID",
//           "category": "ad",
//           "desc": "更新广告",
//           "creator": "系统"
//       },
//       {
//           "id": 82,
//           "createdAt": "2025-01-03T22:01:27.346+08:00",
//           "updatedAt": "2025-01-03T22:01:27.346+08:00",
//           "deletedAt": null,
//           "method": "POST",
//           "path": "/ad",
//           "category": "ad",
//           "desc": "创建广告",
//           "creator": "系统"
//       },
//       {
//           "id": 81,
//           "createdAt": "2025-01-03T22:01:27.346+08:00",
//           "updatedAt": "2025-01-03T22:01:27.346+08:00",
//           "deletedAt": null,
//           "method": "GET",
//           "path": "/ad",
//           "category": "ad",
//           "desc": "获取广告列表",
//           "creator": "系统"
//       },
//       {
//           "id": 80,
//           "createdAt": "2025-01-03T22:01:27.346+08:00",
//           "updatedAt": "2025-01-03T22:01:27.346+08:00",
//           "deletedAt": null,
//           "method": "GET",
//           "path": "/ad/:adID",
//           "category": "ad",
//           "desc": "获取单个广告位",
//           "creator": "系统"
//       }
//   ]
// }

export type Api = {
  id: number
  createdAt: string
  updatedAt: string
  deletedAt: string | null
  method: string
  path: string
  category: string
  desc: string
  creator: string,
  children: Api[]
}

export type ApiListResponse = {
  apis: Api[]
  total: number
}