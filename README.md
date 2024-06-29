<h1 align="center">gotribe-admin</h1>

<div align="center">
Go + Vue开发的管理系统脚手架, 前后端分离, 仅包含项目开发的必需部分, 基于角色的访问控制(RBAC), 分包合理, 精简易于扩展。 后端Go包含了gin、 gorm、 jwt和casbin等的使用, 前端Vue基于vue-element-admin开发: https://github.com/gotribe/gotribe-ui.git
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

![登录](https://github.com/gnimli/go-web-mini-ui/blob/main/src/assets/GithubImages/login.PNG)
![用户管理](https://github.com/gnimli/go-web-mini-ui/blob/main/src/assets/GithubImages/user.PNG)
![角色管理](https://github.com/gnimli/go-web-mini-ui/blob/main/src/assets/GithubImages/role.PNG)
![角色权限](https://github.com/gnimli/go-web-mini-ui/blob/main/src/assets/GithubImages/rolePermission.PNG)
![菜单管理](https://github.com/gnimli/go-web-mini-ui/blob/main/src/assets/GithubImages/menu.PNG)
![API管理](https://github.com/gnimli/go-web-mini-ui/blob/main/src/assets/GithubImages/api.PNG)

## 项目结构概览

```
├─cmd # 启动项目
├─config # viper读取配置
├─internal # 项目目录
├──gotribe-admin #项目的主要业务逻辑
├───controller # controller层，响应路由请求的方法
├───repository # 数据库操作
├───routes # 所有路由
├──pkg # 非公用包
├───common # casbin mysql zap validator 等公共资源
├───middleware # 中间件
├───model # 结构体模型
├─pkg # 公用包
├──api # 接口相关
├────dto # 返回给前端的数据结构
├────response # 常用返回封装，如Success、Fai
├────vo # 接收前端请求的数据结构
├────known # 常量
├──util # 工具包
└─scripts # 脚本
```
## 前端Vue项目
    gotribe-admin-vue
<https://github.com/go-tribe/gotribe-admin-vue>

## TODO

- 增加支付配置
- 增加商品管理

## MIT License

    Copyright (c) 2024 gotribe
