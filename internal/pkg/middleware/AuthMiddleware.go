// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package middleware

import (
	"fmt"
	"gotribe-admin/config"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"gotribe-admin/pkg/util"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitAuth 初始化JWT中间件，配置认证、授权、登录、登出等处理函数
// InitAuth initializes JWT middleware and configures authentication, authorization, login, logout handlers
func InitAuth() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           config.Conf.Jwt.Realm,                                 // jwt标识
		Key:             []byte(config.Conf.Jwt.Key),                           // 服务端密钥
		Timeout:         time.Hour * time.Duration(config.Conf.Jwt.Timeout),    // token过期时间
		MaxRefresh:      time.Hour * time.Duration(config.Conf.Jwt.MaxRefresh), // token最大刷新时间(RefreshToken过期时间=Timeout+MaxRefresh)
		PayloadFunc:     payloadFunc,                                           // 有效载荷处理
		IdentityHandler: identityHandler,                                       // 解析Claims
		Authenticator:   login,                                                 // 校验token的正确性, 处理登录逻辑
		Authorizator:    authorizator,                                          // 用户登录校验成功处理
		Unauthorized:    unauthorized,                                          // 用户登录校验失败处理
		LoginResponse:   loginResponse,                                         // 登录成功后的响应
		LogoutResponse:  logoutResponse,                                        // 登出后的响应
		RefreshResponse: refreshResponse,                                       // 刷新token后的响应
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",    // 自动在这几个地方寻找请求中的token
		TokenHeadName:   "Bearer",                                              // header名称
		TimeFunc:        time.Now,
	})
	return authMiddleware, err
}

// payloadFunc 处理JWT有效载荷，将用户数据转换为Claims格式
// payloadFunc handles JWT payload and converts user data to Claims format
func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		var user model.Admin
		// 将用户json转为结构体
		util.JSONUtil.JsonI2Struct(v["user"], &user)
		return jwt.MapClaims{
			jwt.IdentityKey: user.ID,
			"user":          v["user"],
		}
	}
	return jwt.MapClaims{}
}

// identityHandler 解析JWT Claims，提取用户身份信息
// identityHandler parses JWT Claims and extracts user identity information
func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	// 此处返回值类型map[string]interface{}与payloadFunc和authorizator的data类型必须一致, 否则会导致授权失败还不容易找到原因
	return map[string]interface{}{
		"IdentityKey": claims[jwt.IdentityKey],
		"user":        claims["user"],
	}
}

// login 校验用户登录凭据，验证用户名和密码的正确性
// login validates user login credentials and verifies username and password
func login(c *gin.Context) (interface{}, error) {
	var req vo.RegisterAndLoginRequest
	// 请求json绑定
	if err := c.ShouldBind(&req); err != nil {
		return "", err
	}

	u := &model.Admin{
		Username: req.Username,
		Password: req.Password,
	}

	// 密码校验
	userRepository := repository.NewAdminRepository()
	user, err := userRepository.Login(u)
	if err != nil {
		// 检查是否为RepositoryError，如果是则返回本地化错误消息
		if repoErr, ok := err.(*common.RepositoryError); ok {
			localizedMsg := repoErr.GetLocalizedMessage(c)
			return nil, fmt.Errorf("%s", localizedMsg)
		}
		return nil, err
	}
	// 将用户以json格式写入, payloadFunc/authorizator会使用到
	userJson, err := util.JSONUtil.Struct2Json(user)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", common.Msg(c, common.MsgUserSerializeFail), err)
	}
	return map[string]interface{}{
		"user": userJson,
	}, nil
}

// authorizator 处理用户登录校验成功后的授权逻辑，将用户信息保存到上下文
// authorizator handles authorization logic after successful login validation and saves user info to context
func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(map[string]interface{}); ok {
		userStr := v["user"].(string)
		var user model.Admin
		// 将用户json转为结构体
		util.JSONUtil.Json2Struct(userStr, &user)
		// 将用户保存到context, api调用时取数据方便
		c.Set("user", user)
		return true
	}
	return false
}

// unauthorized 处理用户登录校验失败的情况，返回本地化的错误信息
// unauthorized handles user login validation failure and returns localized error messages
func unauthorized(c *gin.Context, code int, message string) {
	common.Log.Debugf("JWT认证失败, 错误码: %d, 错误信息: %s", code, message)
	// 使用本地化的JWT认证失败消息，并包含具体的错误信息
	localizedMsg := fmt.Sprintf("%s: %s", common.Msg(c, common.MsgJWTAuthFail), message)
	response.Unauthorized(c, localizedMsg)
}

// loginResponse 处理用户登录成功后的响应，返回JWT令牌和过期时间
// loginResponse handles the response after successful user login and returns JWT token and expiration time
func loginResponse(c *gin.Context, _code int, token string, expires time.Time) {
	msg := common.Msg(c, common.MsgLoginSuccess)
	response.Success(c,
		gin.H{
			"token":   token,
			"expires": expires.Format(known.TIME_FORMAT),
		},
		msg)
}

// logoutResponse 处理用户登出后的响应，返回本地化的成功消息
// logoutResponse handles the response after user logout and returns a localized success message
func logoutResponse(c *gin.Context, _code int) {
	msg := common.Msg(c, common.MsgLogoutSuccess)
	response.Success(c, nil, msg)
}

// refreshResponse 处理刷新JWT令牌后的响应，返回新的令牌和过期时间
// refreshResponse handles the response after JWT token refresh and returns new token and expiration time
func refreshResponse(c *gin.Context, _code int, token string, expires time.Time) {
	msg := common.Msg(c, common.MsgRefreshTokenSuccess)
	response.Success(c,
		gin.H{
			"token":   token,
			"expires": expires,
		},
		msg)
}
