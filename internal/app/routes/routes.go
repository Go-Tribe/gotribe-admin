// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"embed"
	"fmt"
	"gotribe-admin/config"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/middleware"
	"net/http"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RouteInitializer 路由初始化器类型
type RouteInitializer func(*gin.RouterGroup, *jwt.GinJWTMiddleware) gin.IRoutes

// InitRoutes 初始化所有路由
func InitRoutes(fs embed.FS) *gin.Engine {
	// 设置运行模式
	gin.SetMode(config.Conf.System.Mode)

	// 创建 Gin 引擎
	r := gin.Default()

	// 设置中间件
	setupMiddlewares(r)

	// 初始化 JWT 认证中间件
	authMiddleware := initAuthMiddleware()

	// 设置静态文件服务
	setupStaticFiles(r, fs)

	// 设置 Swagger 文档路由
	setupSwaggerRoutes(r)

	// 注册 API 路由
	registerAPIRoutes(r, authMiddleware)

	// 注册任务管理路由（特殊处理，不在 API 分组中）
	JobRoutes(r, authMiddleware)

	common.Log.Info("初始化路由完成！")
	return r
}

// setupMiddlewares 设置全局中间件
func setupMiddlewares(r *gin.Engine) {
	// 限流中间件
	fillInterval := time.Duration(config.Conf.RateLimit.FillInterval)
	capacity := config.Conf.RateLimit.Capacity
	r.Use(middleware.RateLimitMiddleware(time.Millisecond*fillInterval, capacity))

	// 全局跨域中间件
	r.Use(middleware.CORSMiddleware())

	// 语言协商中间件（zh/en），默认中文
	r.Use(middleware.LangMiddleware())

	// 操作日志中间件
	r.Use(middleware.OperationLogMiddleware())
}

// initAuthMiddleware 初始化 JWT 认证中间件
func initAuthMiddleware() *jwt.GinJWTMiddleware {
	authMiddleware, err := middleware.InitAuth()
	if err != nil {
		errMsg := fmt.Sprintf("初始化JWT中间件失败：%v", err)
		common.Log.Panicf(errMsg)
		panic(errMsg)
	}
	return authMiddleware
}

// setupStaticFiles 设置静态文件服务
func setupStaticFiles(r *gin.Engine, fs embed.FS) {
	embedFS, err := static.EmbedFolder(fs, "web/admin/dist")
	if err != nil {
		common.Log.Errorf("设置静态文件服务失败：%v", err)
		return
	}

	// 1. 尝试从 embed.FS 中读取 index.html 内容
	// 注意：这里需要使用完整的嵌入路径
	indexData, err := fs.ReadFile("web/admin/dist/index.html")
	if err != nil {
		common.Log.Errorf("读取 index.html 失败：%v", err)
	}

	// 2. 静态资源服务 (处理 /assets, /favicon.ico 等)
	r.Use(static.Serve("/", embedFS))

	// 3. SPA 兜底处理
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// 如果是静态资源请求（以/assets/开头）或API请求，返回404
		if strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/swagger/") {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// 对于前端路由，直接返回内存中的 index.html 内容
		if len(indexData) > 0 {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusOK, string(indexData))
		} else {
			c.String(http.StatusNotFound, "Index file not found")
		}
	})
}

// setupSwaggerRoutes 设置 Swagger 文档路由
func setupSwaggerRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// registerAPIRoutes 注册所有 API 路由
func registerAPIRoutes(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	apiGroup := r.Group("/" + config.Conf.System.UrlPathPrefix)

	// 路由初始化器列表
	routeInitializers := []RouteInitializer{
		InitBaseRoutes,            // 基础路由（无需认证）
		InitAdminRoutes,           // 管理员路由
		InitRoleRoutes,            // 角色管理
		InitMenuRoutes,            // 菜单管理
		InitApiRoutes,             // 接口管理
		InitOperationLogRoutes,    // 操作日志
		InitProjectRoutes,         // 项目管理
		InitConfigRoutes,          // 配置管理
		InitTagRoutes,             // 标签管理
		InitCategoryRoutes,        // 分类管理
		InitPostRoutes,            // 内容管理
		InitUserRoutes,            // 用户管理
		InitResourceRoutes,        // 资源管理
		InitColumnRoutes,          // 专栏管理
		InitAdSceneRoutes,         // 推广场景管理
		InitAdRoutes,              // 广告位管理
		InitCommentRoutes,         // 评论管理
		InitPointRoutes,           // 积分管理
		InitProductCategoryRoutes, // 商品分类管理
		InitProductTypeRoutes,     // 商品类型管理
		InitProductSpecRoutes,     // 商品规格管理
		InitProductSpecItemRoutes, // 商品规格项管理
		InitProductRoutes,         // 商品管理
		InitOrderRoutes,           // 订单管理
		InitSystemConfigRoutes,    // 系统配置管理
		InitFeedbackRoutes,        // 反馈管理
		InitIndexRoutes,           // 首页数据
	}

	// 批量注册路由
	for _, initFunc := range routeInitializers {
		initFunc(apiGroup, authMiddleware)
	}
}
