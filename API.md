# GoTribe Admin API 文档

## 概述

GoTribe Admin 提供完整的 RESTful API 接口，支持内容管理、用户管理、系统管理等功能。

## 基础信息

- **Base URL**: `http://localhost:8088`
- **API Version**: `v1`
- **Content-Type**: `application/json`
- **认证方式**: JWT Token

## 认证

### 登录获取Token

```http
POST /api/admin/login
Content-Type: application/json

{
  "username": "admin",
  "password": "123456"
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "nickname": "管理员",
      "avatar": "https://example.com/avatar.jpg"
    }
  }
}
```

### 请求头设置

```http
Authorization: Bearer <token>
Content-Type: application/json
```

## 用户管理 API

### 获取用户列表

```http
GET /api/user?page=1&pageSize=10&username=test
```

**查询参数:**
- `page`: 页码 (默认: 1)
- `pageSize`: 每页数量 (默认: 10)
- `username`: 用户名搜索
- `email`: 邮箱搜索
- `phone`: 手机号搜索
- `status`: 状态筛选

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "userID": "gotribe",
        "username": "test",
        "nickname": "测试用户",
        "email": "test@example.com",
        "phone": "13800138000",
        "status": 1,
        "createdAt": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "pageSize": 10
  }
}
```

### 创建用户

```http
POST /api/user
Content-Type: application/json

{
  "username": "newuser",
  "nickname": "新用户",
  "email": "newuser@example.com",
  "phone": "13900139000",
  "password": "123456"
}
```

### 更新用户

```http
PATCH /api/user/:id
Content-Type: application/json

{
  "nickname": "更新后的昵称",
  "email": "updated@example.com"
}
```

### 删除用户

```http
DELETE /api/user
Content-Type: application/json

{
  "ids": [1, 2, 3]
}
```

## 内容管理 API

### 文章管理

#### 获取文章列表

```http
GET /api/post?page=1&pageSize=10&title=测试&categoryID=1&status=2
```

**查询参数:**
- `page`: 页码
- `pageSize`: 每页数量
- `title`: 标题搜索
- `categoryID`: 分类ID
- `status`: 状态 (1: 草稿, 2: 发布)
- `type`: 类型 (1: 文章, 2: 页面, 3: 短文)

#### 创建文章

```http
POST /api/post
Content-Type: application/json

{
  "title": "文章标题",
  "content": "文章内容",
  "htmlContent": "<h1>文章内容</h1>",
  "description": "文章描述",
  "categoryID": "24ejga",
  "tag": "标签1,标签2",
  "status": 2,
  "type": 1,
  "unitPrice": 0
}
```

#### 更新文章

```http
PATCH /api/post/:id
Content-Type: application/json

{
  "title": "更新后的标题",
  "content": "更新后的内容"
}
```

### 分类管理

#### 获取分类列表

```http
GET /api/category?page=1&pageSize=10
```

#### 创建分类

```http
POST /api/category
Content-Type: application/json

{
  "title": "新分类",
  "description": "分类描述",
  "status": 1
}
```

### 标签管理

#### 获取标签列表

```http
GET /api/tag?page=1&pageSize=10
```

#### 创建标签

```http
POST /api/tag
Content-Type: application/json

{
  "title": "新标签",
  "description": "标签描述"
}
```

## 项目管理 API

### 获取项目列表

```http
GET /api/project?page=1&pageSize=10
```

### 创建项目

```http
POST /api/project
Content-Type: application/json

{
  "name": "project-name",
  "title": "项目标题",
  "description": "项目描述",
  "domain": "example.com",
  "logo": "https://example.com/logo.png"
}
```

### 更新项目

```http
PATCH /api/project/:id
Content-Type: application/json

{
  "title": "更新后的标题",
  "description": "更新后的描述"
}
```

## 商品管理 API

### 商品管理

#### 获取商品列表

```http
GET /api/product?page=1&pageSize=10&title=商品&categoryID=1&enable=2
```

**查询参数:**
- `page`: 页码
- `pageSize`: 每页数量
- `title`: 商品标题搜索
- `categoryID`: 分类ID
- `enable`: 启用状态 (1: 下架, 2: 上架)

#### 创建商品

```http
POST /api/product
Content-Type: application/json

{
  "title": "商品标题",
  "productNumber": "SKU001",
  "description": "商品描述",
  "image": "https://example.com/product.jpg",
  "categoryID": "24ejga",
  "productSpec": "规格信息",
  "content": "商品详情",
  "tag": "标签1,标签2",
  "enable": 2,
  "skus": [
    {
      "title": "规格1",
      "costPrice": 1000,
      "unitPrice": 1500,
      "marketPrice": 2000,
      "unitPoint": 100,
      "quantity": 100,
      "enableDefault": 1
    }
  ]
}
```

### 订单管理

#### 获取订单列表

```http
GET /api/order?page=1&pageSize=10&status=1&userID=gotribe
```

**查询参数:**
- `page`: 页码
- `pageSize`: 每页数量
- `status`: 订单状态
- `userID`: 用户ID
- `orderID`: 订单号

#### 更新订单状态

```http
PATCH /api/order/:id
Content-Type: application/json

{
  "status": 2,
  "amountPay": 1500,
  "remarkAdmin": "管理员备注"
}
```

## 系统管理 API

### 菜单管理

#### 获取菜单列表

```http
GET /api/menu?page=1&pageSize=10
```

#### 创建菜单

```http
POST /api/menu
Content-Type: application/json

{
  "name": "MenuName",
  "title": "菜单标题",
  "path": "/menu-path",
  "component": "/menu/component",
  "icon": "icon-name",
  "parentID": 0,
  "sort": 1
}
```

### 角色管理

#### 获取角色列表

```http
GET /api/role?page=1&pageSize=10
```

#### 创建角色

```http
POST /api/role
Content-Type: application/json

{
  "name": "角色名称",
  "keyword": "role_keyword",
  "desc": "角色描述",
  "sort": 1
}
```

### 管理员管理

#### 获取管理员列表

```http
GET /api/admin?page=1&pageSize=10
```

#### 创建管理员

```http
POST /api/admin
Content-Type: application/json

{
  "username": "newadmin",
  "password": "123456",
  "mobile": "13800138000",
  "nickname": "新管理员",
  "avatar": "https://example.com/avatar.jpg",
  "roleIDs": [1, 2]
}
```

## 资源管理 API

### 文件上传

```http
POST /api/resource/upload
Content-Type: multipart/form-data

file: [文件]
```

**响应示例:**
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "id": 1,
    "name": "example.jpg",
    "url": "https://example.com/uploads/example.jpg",
    "size": 1024000,
    "type": 1,
    "domain": "https://cdn.example.com"
  }
}
```

### 获取资源列表

```http
GET /api/resource?page=1&pageSize=10&type=1
```

**查询参数:**
- `page`: 页码
- `pageSize`: 每页数量
- `type`: 文件类型 (1: 图片, 2: 音频, 3: 视频, 4: 文档, 5: 应用, 6: 压缩包, 7: 字体, 8: 未知)

## 积分管理 API

### 获取积分记录

```http
GET /api/point?page=1&pageSize=10&userID=gotribe&type=add
```

**查询参数:**
- `page`: 页码
- `pageSize`: 每页数量
- `userID`: 用户ID
- `type`: 积分类型 (add: 增加, deduct: 扣除)

### 创建积分记录

```http
POST /api/point
Content-Type: application/json

{
  "userID": "gotribe",
  "points": 100,
  "reason": "签到奖励",
  "type": "add",
  "eventID": "event_001"
}
```

## 系统配置 API

### 获取系统配置

```http
GET /api/system
```

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "siteName": "GoTribe",
    "siteDescription": "GoTribe CMS",
    "siteKeywords": "cms,go,vue",
    "siteLogo": "https://example.com/logo.png",
    "siteFavicon": "https://example.com/favicon.ico",
    "siteCopyright": "© 2023 GoTribe",
    "siteICP": "京ICP备12345678号",
    "siteBeian": "京公网安备12345678901234号",
    "siteAnalytics": "UA-123456789-1",
    "siteBaiduAnalytics": "1234567890",
    "siteBaiduSeoPush": "https://example.com/push",
    "siteCDNDomain": "https://cdn.example.com",
    "siteUploadPath": "/uploads",
    "siteUploadMaxSize": 10485760,
    "siteUploadAllowedExts": "jpg,jpeg,png,gif,mp4,mp3,pdf,doc,docx",
    "siteUploadAllowedTypes": "1,2,3,4,5,6,7,8",
    "siteUploadAllowedMimes": "image/jpeg,image/png,image/gif,video/mp4,audio/mp3,application/pdf",
    "siteUploadAllowedSize": 10485760,
    "siteUploadAllowedWidth": 1920,
    "siteUploadAllowedHeight": 1080,
    "siteUploadAllowedRatio": "16:9",
    "siteUploadAllowedQuality": 80,
    "siteUploadAllowedWatermark": true,
    "siteUploadAllowedWatermarkText": "GoTribe",
    "siteUploadAllowedWatermarkImage": "https://example.com/watermark.png",
    "siteUploadAllowedWatermarkPosition": "bottom-right",
    "siteUploadAllowedWatermarkOpacity": 0.5,
    "siteUploadAllowedWatermarkSize": 0.1,
    "siteUploadAllowedWatermarkMargin": 10,
    "siteUploadAllowedWatermarkColor": "#ffffff",
    "siteUploadAllowedWatermarkFontSize": 14,
    "siteUploadAllowedWatermarkFontFamily": "Arial",
    "siteUploadAllowedWatermarkFontWeight": "bold",
    "siteUploadAllowedWatermarkFontStyle": "normal",
    "siteUploadAllowedWatermarkFontDecoration": "none",
    "siteUploadAllowedWatermarkFontStretch": "normal",
    "siteUploadAllowedWatermarkFontKerning": "auto",
    "siteUploadAllowedWatermarkFontVariant": "normal",
    "siteUploadAllowedWatermarkFontFeatureSettings": "normal",
    "siteUploadAllowedWatermarkFontVariationSettings": "normal",
    "siteUploadAllowedWatermarkFontDisplay": "auto",
    "siteUploadAllowedWatermarkFontSynthesis": "auto",
    "siteUploadAllowedWatermarkFontVariantLigatures": "normal",
    "siteUploadAllowedWatermarkFontVariantCaps": "normal",
    "siteUploadAllowedWatermarkFontVariantNumeric": "normal",
    "siteUploadAllowedWatermarkFontVariantAlternates": "normal",
    "siteUploadAllowedWatermarkFontVariantEastAsian": "normal",
    "siteUploadAllowedWatermarkFontVariantLigatures": "normal",
    "siteUploadAllowedWatermarkFontVariantCaps": "normal",
    "siteUploadAllowedWatermarkFontVariantNumeric": "normal",
    "siteUploadAllowedWatermarkFontVariantAlternates": "normal",
    "siteUploadAllowedWatermarkFontVariantEastAsian": "normal"
  }
}
```

### 更新系统配置

```http
PATCH /api/system
Content-Type: application/json

{
  "siteName": "新的站点名称",
  "siteDescription": "新的站点描述"
}
```

## 统计信息 API

### 获取首页统计

```http
GET /api/index
```

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "sales": 12345.67,
    "orders": 100,
    "newUsers": 50,
    "visitCount": 1000
  }
}
```

### 获取时间范围数据

```http
GET /api/index/data?startDate=2023-01-01&endDate=2023-01-31&type=sales
```

**查询参数:**
- `startDate`: 开始日期 (YYYY-MM-DD)
- `endDate`: 结束日期 (YYYY-MM-DD)
- `type`: 数据类型 (sales: 销售, orders: 订单, users: 用户, visits: 访问)

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 响应格式

### 成功响应

```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": "2023-01-01T00:00:00Z"
}
```

### 错误响应

```json
{
  "code": 400,
  "message": "参数错误",
  "data": null,
  "timestamp": "2023-01-01T00:00:00Z"
}
```

### 分页响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "pageSize": 10
  },
  "timestamp": "2023-01-01T00:00:00Z"
}
```

## 限流说明

- **登录接口**: 5次/分钟
- **其他接口**: 100次/分钟
- **上传接口**: 10次/分钟

## 注意事项

1. 所有时间格式均为 ISO 8601 格式
2. 分页参数 `page` 从 1 开始
3. 文件上传需要设置正确的 Content-Type
4. 所有接口都需要在请求头中携带认证信息
5. 金额字段单位为分（整数），显示时需要转换为元
6. 积分字段单位为分（整数），显示时需要转换为元
