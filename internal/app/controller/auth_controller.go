// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"fmt"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"gotribe-admin/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

// 注意：此控制器仅用于生成 Swagger 文档
// 实际的认证逻辑由 JWT 中间件处理，路由在 base_routes.go 中定义
type IAuthController interface {
	Login(c *gin.Context)        // 用户登录
	Logout(c *gin.Context)       // 用户登出
	RefreshToken(c *gin.Context) // 刷新token
}

type AuthController struct {
	AdminRepository repository.IAdminRepository
}

// 构造函数
func NewAuthController() IAuthController {
	adminRepository := repository.NewAdminRepository()
	authController := AuthController{AdminRepository: adminRepository}
	return authController
}

// Login 用户登录
// @Summary      用户登录
// @Description  管理员用户登录接口，返回JWT token
// @Tags         认证管理
// @Accept       json
// @Produce      json
// @Param        request body vo.RegisterAndLoginRequest true "登录请求"
// @Success      200 {object} response.Response{data=object{token=string,expires=string}} "登录成功"
// @Failure      400 {object} response.Response "登录失败"
// @Router       /base/login [post]
func (ac AuthController) Login(c *gin.Context) {
	// 注意：此方法仅用于 Swagger 文档生成，实际登录由 JWT 中间件处理
	var req vo.RegisterAndLoginRequest
	// 请求json绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	u := &model.Admin{
		Username: req.Username,
		Password: req.Password,
	}

	// 密码校验
	user, err := ac.AdminRepository.Login(u)
	if err != nil {
		response.PasswordIncorrect(c, "用户名或密码错误")
		return
	}

	// 将用户以json格式写入, payloadFunc/authorizator会使用到
	userJson, err := util.JSONUtil.Struct2Json(user)
	if err != nil {
		response.InternalServerError(c, fmt.Sprintf("用户信息序列化失败: %v", err))
		return
	}

	// 这里应该调用JWT中间件的LoginHandler，但为了Swagger文档，我们提供一个示例响应
	// 实际的JWT token生成会由中间件处理
	msg := common.Msg(c, common.MsgLoginSuccess)
	response.Success(c, gin.H{
		"token":   "jwt_token_will_be_generated_by_middleware",
		"expires": time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
		"user":    userJson,
	}, msg)
}

// Logout 用户登出
// @Summary      用户登出
// @Description  管理员用户登出接口
// @Tags         认证管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response "登出成功"
// @Failure      400 {object} response.Response "登出失败"
// @Router       /base/logout [post]
// @Security     BearerAuth
func (ac AuthController) Logout(c *gin.Context) {
	// 注意：此方法仅用于 Swagger 文档生成，实际登出由 JWT 中间件处理
	msg := common.Msg(c, common.MsgLogoutSuccess)
	response.Success(c, nil, msg)
}

// RefreshToken 刷新token
// @Summary      刷新token
// @Description  刷新JWT token，延长登录状态
// @Tags         认证管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response{data=object{token=string,expires=string}} "刷新成功"
// @Failure      400 {object} response.Response "刷新失败"
// @Router       /base/refreshToken [post]
// @Security     BearerAuth
func (ac AuthController) RefreshToken(c *gin.Context) {
	// 注意：此方法仅用于 Swagger 文档生成，实际token刷新由 JWT 中间件处理
	msg := common.Msg(c, common.MsgRefreshTokenSuccess)
	response.Success(c, gin.H{
		"token":   "new_jwt_token_will_be_generated_by_middleware",
		"expires": time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	}, msg)
}
