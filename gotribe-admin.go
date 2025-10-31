// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

// @title           Gotribe Admin API
// @version         1.0
// @description     这是一个使用 Swagger 进行接口文档管理的 Go API 服务器。
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8088
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 请输入 Bearer {token}，注意 Bearer 和 token 之间有一个空格

package main

import (
	"context"
	"embed"
	"fmt"
	"gotribe-admin/config"
	_ "gotribe-admin/docs" // swagger docs
	"gotribe-admin/internal/app/jobs"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/app/routes"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
)

//go:embed web/admin/dist/*
var content embed.FS

func main() {

	// 加载配置文件到全局配置结构体
	config.InitConfig()

	// 初始化日志
	common.InitLogger()

	// 初始化数据库
	common.InitDatabase()

	// 初始化casbin策略管理器
	common.InitCasbinEnforcer()

	// 初始化Validator数据校验
	common.InitValidate()

	// 初始化数据库数据
	common.InitData()

	// 初始化定时任务
	if err := jobs.InitJobs(); err != nil {
		common.Log.Errorf("Failed to initialize jobs: %v", err)
		return
	}

	if err := jobs.StartAllJobs(); err != nil {
		common.Log.Errorf("Failed to start jobs: %v", err)
		return
	}
	defer jobs.StopAllJobs()

	// 操作日志中间件处理日志时没有将日志发送到rabbitmq或者kafka中, 而是发送到了channel中
	// 这里开启3个goroutine处理channel将日志记录到数据库
	logRepository := repository.NewOperationLogRepository()
	for i := 0; i < 3; i++ {
		go logRepository.SaveOperationLogChannel(middleware.OperationLogChan)
	}

	// 注册所有路由
	r := routes.InitRoutes(content)

	host := "localhost"
	port := config.Conf.System.Port

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			common.Log.Fatalf("listen: %s\n", err)
		}
	}()

	colorFg := color.New(color.FgCyan, color.Bold)
	colorFg.Println(`
	░██████╗░░█████╗░████████╗██████╗░██╗██████╗░███████╗
	██╔════╝░██╔══██╗╚══██╔══╝██╔══██╗██║██╔══██╗██╔════╝
	██║░░██╗░██║░░██║░░░██║░░░██████╔╝██║██████╦╝█████╗░░
	██║░░╚██╗██║░░██║░░░██║░░░██╔══██╗██║██╔══██╗██╔══╝░░
	╚██████╔╝╚█████╔╝░░░██║░░░██║░░██║██║██████╦╝███████╗
	░╚═════╝░░╚════╝░░░░╚═╝░░░╚═╝░░╚═╝╚═╝╚═════╝░╚══════╝
`)
	fmt.Println("	App running at:")
	fmt.Println(fmt.Sprintf("	- Local: %s%s:%d", "http://", host, port))
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	common.Log.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		common.Log.Fatal("Server forced to shutdown:", err)
	}

	common.Log.Info("Server exiting!")
}
