<h1 align="center">gotribe-admin</h1>

<div align="center">
Go + Vue开发的小型Cms解决方案, 前后端分离。 由管理端 API，用户端 API，管理后台。
<p align="center">
<img src="https://img.shields.io/github/go-mod/go-version/go-tribe/gotribe-admin" alt="Go version"/>
<img src="https://img.shields.io/badge/Gin-1.9.1-brightgreen" alt="Gin version"/>
<img src="https://img.shields.io/badge/Gorm-1.25.8-brightgreen" alt="Gorm version"/>
<img src="https://img.shields.io/github/license/go-tribe/gotribe-admin" alt="License"/>
</p>
</div>

## 特性

- `Gin` 一个类似于martini但拥有更好性能的API框架, 由于使用了httprouter, 速度提高了近40倍
- `MySQL` 采用的是MySql数据库
- `Jwt` 使用JWT轻量级认证, 并提供活跃用户Token刷新功能
- `Casbin` Casbin是一个强大的、高效的开源访问控制框架，其权限管理机制支持多种访问控制模型
- `Gorm` 采用Gorm 2.0版本开发, 包含一对多、多对多、事务等操作
- `Validator` 使用validator v10做参数校验, 严密校验前端传入参数
- `Lumberjack` 设置日志文件大小、保存数量、保存时间和压缩等 ```
- `Viper` Go应用程序的完整配置解决方案, 支持配置热更新

## 中间件

- `AuthMiddleware` 权限认证中间件 -- 处理登录、登出、无状态token校验
- `RateLimitMiddleware` 基于令牌桶的限流中间件 -- 限制用户的请求次数
- `OperationLogMiddleware` 操作日志中间件 -- 记录所有用户操作
- `CORSMiddleware` -- 跨域中间件 -- 解决跨域问题
- `CasbinMiddleware` 访问控制中间件 -- 基于Casbin RBAC, 精细控制接口访问

## 项目截图

![登录](https://github.com/go-tribe/gotribe-admin/docs/images/login.PNG)
![后台首页](https://github.com/go-tribe/gotribe-admin/docs/images/index.PNG)
![系统管理](https://github.com/go-tribe/gotribe-admin/docs/images/system.PNG)
![日志管理](https://github.com/go-tribe/gotribe-admin/docs/images/log.PNG)
![业务管理](https://github.com/go-tribe/gotribe-admin/docs/images/project.PNG)
![内容管理](https://github.com/go-tribe/gotribe-admin/docs/images/content.PNG)

## 项目合集
| 项目 | 描述       |地址|
| --- |----------| --- |
| gotribe-admin | 后台管理 api | https://github.com/go-tribe/gotribe-admin.git |
| gotribe | 业务端 api  | https://github.com/go-tribe/gotribe.git |
| gotribe-ui | 前端管理后台   | https://github.com/go-tribe/gotribe-admin-vue.git |

## TODO

- 增加支付配置
- 增加商品管理

## MIT License

    Copyright (c) 2024 gotribe
