// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gotribe-admin/config"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/middleware"
	"net/http"

	"time"
)

// 初始化
func InitRoutes(fs embed.FS) *gin.Engine {
	//设置模式
	gin.SetMode(config.Conf.System.Mode)

	// 创建带有默认中间件的路由:
	// 日志与恢复中间件
	r := gin.Default()
	// 创建不带中间件的路由:
	// r := gin.New()
	// r.Use(gin.Recovery())

	// 启用限流中间件
	// 默认每50毫秒填充一个令牌，最多填充200个
	fillInterval := time.Duration(config.Conf.RateLimit.FillInterval)
	capacity := config.Conf.RateLimit.Capacity
	r.Use(middleware.RateLimitMiddleware(time.Millisecond*fillInterval, capacity))

	// 启用全局跨域中间件
	r.Use(middleware.CORSMiddleware())

	// 启用操作日志中间件
	r.Use(middleware.OperationLogMiddleware())

	// 初始化JWT认证中间件
	authMiddleware, err := middleware.InitAuth()
	if err != nil {
		common.Log.Panicf("初始化JWT中间件失败：%v", err)
		panic(fmt.Sprintf("初始化JWT中间件失败：%v", err))
	}
	r.Use(static.Serve("/", static.EmbedFolder(fs, "web/admin/dist")))
	r.NoRoute(func(c *gin.Context) {
		common.Log.Infof("A 404 error occurred, but the specific URL path is not logged to prevent log injection.")
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	// end

	// 路由分组
	apiGroup := r.Group("/" + config.Conf.System.UrlPathPrefix)

	// 注册路由
	InitBaseRoutes(apiGroup, authMiddleware)         // 注册基础路由, 不需要jwt认证中间件,不需要casbin中间件
	InitAdminRoutes(apiGroup, authMiddleware)        // 注册用户路由, jwt认证中间件,casbin鉴权中间件
	InitRoleRoutes(apiGroup, authMiddleware)         // 注册角色路由, jwt认证中间件,casbin鉴权中间件
	InitMenuRoutes(apiGroup, authMiddleware)         // 注册菜单路由, jwt认证中间件,casbin鉴权中间件
	InitApiRoutes(apiGroup, authMiddleware)          // 注册接口路由, jwt认证中间件,casbin鉴权中间件
	InitOperationLogRoutes(apiGroup, authMiddleware) // 注册操作日志路由, jwt认证中间件,casbin鉴权中间件
	InitProjectRoutes(apiGroup, authMiddleware)      // 注册项目管理路由, jwt认证中间件,casbin鉴权中间件
	InitConfigRoutes(apiGroup, authMiddleware)       // 注册配置管理路由, jwt认证中间件,casbin鉴权中间件
	InitTagRoutes(apiGroup, authMiddleware)          // 注册标签管理路由, jwt认证中间件,casbin鉴权中间件
	InitCategoryRoutes(apiGroup, authMiddleware)     // 注册分类管理路由, jwt认证中间件,casbin鉴权中间件
	InitPostRoutes(apiGroup, authMiddleware)         // 注册内容管理路由, jwt认证中间件,casbin鉴权中间件
	InitUserRoutes(apiGroup, authMiddleware)         // 注册用户管理路由, jwt认证中间件,casbin鉴权中间件
	InitResourceRoutes(apiGroup, authMiddleware)     // 注册资源管理路由, jwt认证中间件,casbin鉴权中间件
	InitColumnRoutes(apiGroup, authMiddleware)       // 注册专栏管理路由, jwt认证中间件,casbin鉴权中间件
	InitAdSceneRoutes(apiGroup, authMiddleware)      // 注册推广场景管理路由, jwt认证中间件,casbin鉴权中间件
	InitAdRoutes(apiGroup, authMiddleware)           // 注册广告位管理路由, jwt认证中间件,casbin鉴权中间件
	InitCommentRoutes(apiGroup, authMiddleware)      // 注册评论管理路由, jwt认证中间件,casbin鉴权中间件
	common.Log.Info("初始化路由完成！")
	return r
}
