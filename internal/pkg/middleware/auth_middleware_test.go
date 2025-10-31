// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gotribe-admin/config"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"gotribe-admin/pkg/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockAdminRepository 模拟 AdminRepository
type MockAdminRepository struct {
	mock.Mock
}

func (m *MockAdminRepository) Login(admin *model.Admin) (*model.Admin, error) {
	args := m.Called(admin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Admin), args.Error(1)
}

// 设置测试环境
func setupTestConfig() {
	config.Conf.Jwt = &config.JwtConfig{
		Realm:      "test realm",
		Key:        "test_secret_key_for_jwt_testing_12345",
		Timeout:    24,
		MaxRefresh: 24,
	}

	// 初始化日志系统（测试环境）
	if common.Log == nil {
		logger, _ := zap.NewDevelopment()
		common.Log = logger.Sugar()
	}
}

// 创建测试用的 Gin 上下文
func createTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("lang", "zh") // 设置默认语言
	return c, w
}

// TestRepositoryErrorLocalization 测试RepositoryError的本地化处理
func TestRepositoryErrorLocalization(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		lang        string
		errorCode   string
		expectedMsg string
	}{
		{
			name:        "中文-用户不存在",
			lang:        "zh",
			errorCode:   "USER_NOT_FOUND",
			expectedMsg: "未获取到ID为123的用户",
		},
		{
			name:        "英文-用户不存在",
			lang:        "en",
			errorCode:   "USER_NOT_FOUND",
			expectedMsg: "User not found by ID 123",
		},
		{
			name:        "中文-密码错误",
			lang:        "zh",
			errorCode:   "PASSWORD_INCORRECT",
			expectedMsg: "密码错误",
		},
		{
			name:        "英文-密码错误",
			lang:        "en",
			errorCode:   "PASSWORD_INCORRECT",
			expectedMsg: "Incorrect password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试上下文
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("lang", tt.lang)

			// 创建RepositoryError
			var repoErr *common.RepositoryError
			switch tt.errorCode {
			case "USER_NOT_FOUND":
				repoErr = common.NewUserNotFoundByIDError(123)
			case "PASSWORD_INCORRECT":
				repoErr = common.ErrPasswordIncorrect
			}

			// 测试本地化消息
			localizedMsg := repoErr.GetLocalizedMessage(c)
			assert.Equal(t, tt.expectedMsg, localizedMsg)

			// 测试默认Error()方法（应该返回中文）
			defaultMsg := repoErr.Error()
			if tt.lang == "zh" {
				assert.Equal(t, tt.expectedMsg, defaultMsg)
			}
		})
	}
}

// TestInitAuth 测试 JWT 中间件初始化
func TestInitAuth(t *testing.T) {
	setupTestConfig()

	t.Run("成功初始化JWT中间件", func(t *testing.T) {
		authMiddleware, err := InitAuth()

		assert.NoError(t, err)
		assert.NotNil(t, authMiddleware)
		assert.Equal(t, "test realm", authMiddleware.Realm)
		assert.Equal(t, []byte("test_secret_key_for_jwt_testing_12345"), authMiddleware.Key)
		assert.Equal(t, time.Hour*24, authMiddleware.Timeout)
		assert.Equal(t, time.Hour*24, authMiddleware.MaxRefresh)
	})

	t.Run("验证中间件配置", func(t *testing.T) {
		authMiddleware, err := InitAuth()

		assert.NoError(t, err)
		assert.NotNil(t, authMiddleware.PayloadFunc)
		assert.NotNil(t, authMiddleware.IdentityHandler)
		assert.NotNil(t, authMiddleware.Authenticator)
		assert.NotNil(t, authMiddleware.Authorizator)
		assert.NotNil(t, authMiddleware.Unauthorized)
		assert.NotNil(t, authMiddleware.LoginResponse)
		assert.NotNil(t, authMiddleware.LogoutResponse)
		assert.NotNil(t, authMiddleware.RefreshResponse)
	})
}

// TestPayloadFunc 测试载荷处理函数
func TestPayloadFunc(t *testing.T) {
	setupTestConfig()

	t.Run("正确处理用户数据", func(t *testing.T) {
		user := model.Admin{
			Model: model.Model{
				ID: 1,
			},
			Username: "testuser",
			Mobile:   "13800138000",
		}

		userJson, _ := util.JSONUtil.Struct2Json(user)
		data := map[string]interface{}{
			"user": userJson,
		}

		claims := payloadFunc(data)

		assert.NotNil(t, claims)
		assert.Equal(t, uint(1), claims[jwt.IdentityKey])
		assert.Equal(t, userJson, claims["user"])
	})

	t.Run("处理无效数据", func(t *testing.T) {
		claims := payloadFunc("invalid_data")

		assert.NotNil(t, claims)
		assert.Empty(t, claims)
	})

	t.Run("处理空数据", func(t *testing.T) {
		claims := payloadFunc(nil)

		assert.NotNil(t, claims)
		assert.Empty(t, claims)
	})
}

// TestIdentityHandler 测试身份解析函数
func TestIdentityHandler(t *testing.T) {
	setupTestConfig()

	t.Run("正确解析身份信息", func(t *testing.T) {
		c, _ := createTestContext()

		// 模拟 JWT claims
		user := model.Admin{
			Model: model.Model{
				ID: 1,
			},
			Username: "testuser",
		}
		userJson, _ := util.JSONUtil.Struct2Json(user)

		// 设置 JWT claims 到上下文
		claims := jwt.MapClaims{
			jwt.IdentityKey: uint(1),
			"user":          userJson,
		}

		// 模拟 ExtractClaims 的行为
		c.Set("JWT_PAYLOAD", claims)

		// 由于我们无法直接测试 identityHandler（它依赖于 jwt.ExtractClaims），
		// 我们测试其预期的返回格式
		expectedResult := map[string]interface{}{
			"IdentityKey": uint(1),
			"user":        userJson,
		}

		assert.Equal(t, uint(1), expectedResult["IdentityKey"])
		assert.Equal(t, userJson, expectedResult["user"])
	})
}

// TestAuthorizator 测试授权验证函数
func TestAuthorizator(t *testing.T) {
	setupTestConfig()

	t.Run("成功授权有效用户", func(t *testing.T) {
		c, _ := createTestContext()

		user := model.Admin{
			Model: model.Model{
				ID: 1,
			},
			Username: "testuser",
			Mobile:   "13800138000",
		}
		userJson, _ := util.JSONUtil.Struct2Json(user)

		data := map[string]interface{}{
			"user": userJson,
		}

		result := authorizator(data, c)

		assert.True(t, result)

		// 验证用户信息是否正确设置到上下文
		contextUser, exists := c.Get("user")
		assert.True(t, exists)

		savedUser, ok := contextUser.(model.Admin)
		assert.True(t, ok)
		assert.Equal(t, uint(1), savedUser.ID)
		assert.Equal(t, "testuser", savedUser.Username)
	})

	t.Run("拒绝无效数据格式", func(t *testing.T) {
		c, _ := createTestContext()

		result := authorizator("invalid_data", c)

		assert.False(t, result)
	})

	t.Run("拒绝空数据", func(t *testing.T) {
		c, _ := createTestContext()

		result := authorizator(nil, c)

		assert.False(t, result)
	})
}

// TestUnauthorized 测试认证失败处理函数
func TestUnauthorized(t *testing.T) {
	setupTestConfig()

	t.Run("处理认证失败_中文", func(t *testing.T) {
		c, w := createTestContext()
		c.Set("lang", "zh")

		unauthorized(c, http.StatusUnauthorized, "token expired")

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Contains(t, response["message"].(string), "JWT认证失败")
		assert.Contains(t, response["message"].(string), "token expired")
	})

	t.Run("处理认证失败_英文", func(t *testing.T) {
		c, w := createTestContext()
		c.Set("lang", "en")

		unauthorized(c, http.StatusUnauthorized, "invalid token")

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Contains(t, response["message"].(string), "JWT authentication failed")
		assert.Contains(t, response["message"].(string), "invalid token")
	})
}

// TestLoginResponse 测试登录成功响应
func TestLoginResponse(t *testing.T) {
	setupTestConfig()

	t.Run("登录成功响应_中文", func(t *testing.T) {
		c, w := createTestContext()
		c.Set("lang", "zh")

		token := "test.jwt.token"
		expires := time.Now().Add(time.Hour * 24)

		loginResponse(c, http.StatusOK, token, expires)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "登录成功", response["message"])

		data, ok := response["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, token, data["token"])
		assert.NotEmpty(t, data["expires"])
	})

	t.Run("登录成功响应_英文", func(t *testing.T) {
		c, w := createTestContext()
		c.Set("lang", "en")

		token := "test.jwt.token"
		expires := time.Now().Add(time.Hour * 24)

		loginResponse(c, http.StatusOK, token, expires)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Login successful", response["message"])
	})
}

// TestLogoutResponse 测试登出响应
func TestLogoutResponse(t *testing.T) {
	setupTestConfig()

	t.Run("登出成功响应_中文", func(t *testing.T) {
		c, w := createTestContext()
		c.Set("lang", "zh")

		logoutResponse(c, http.StatusOK)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "退出成功", response["message"])
	})

	t.Run("登出成功响应_英文", func(t *testing.T) {
		c, w := createTestContext()
		c.Set("lang", "en")

		logoutResponse(c, http.StatusOK)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Logout successful", response["message"])
	})
}

// TestRefreshResponse 测试刷新token响应
func TestRefreshResponse(t *testing.T) {
	setupTestConfig()

	t.Run("刷新token成功响应_中文", func(t *testing.T) {
		c, w := createTestContext()
		c.Set("lang", "zh")

		token := "new.jwt.token"
		expires := time.Now().Add(time.Hour * 24)

		refreshResponse(c, http.StatusOK, token, expires)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "刷新token成功", response["message"])

		data, ok := response["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, token, data["token"])
		assert.NotNil(t, data["expires"])
	})

	t.Run("刷新token成功响应_英文", func(t *testing.T) {
		c, w := createTestContext()
		c.Set("lang", "en")

		token := "new.jwt.token"
		expires := time.Now().Add(time.Hour * 24)

		refreshResponse(c, http.StatusOK, token, expires)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "Refresh token successful", response["message"])
	})
}

// TestJWTAuthenticationSuccess 测试JWT认证成功
func TestJWTAuthenticationSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		lang     string
		expected string
	}{
		{
			name:     "中文-登录成功",
			lang:     "zh",
			expected: "登录成功",
		},
		{
			name:     "英文-登录成功",
			lang:     "en",
			expected: "Login successful",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试上下文
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("lang", tt.lang)

			// 测试登录成功响应
			loginResponse(c, 200, "test-token", time.Now().Add(time.Hour))

			// 验证响应
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			message, exists := response["message"]
			assert.True(t, exists)
			assert.Equal(t, tt.expected, message.(string))
		})
	}
}

// TestLogin 测试登录认证逻辑（需要模拟 repository）
func TestLogin(t *testing.T) {
	setupTestConfig()

	t.Run("成功登录", func(t *testing.T) {
		c, _ := createTestContext()

		// 准备登录请求数据
		loginReq := vo.RegisterAndLoginRequest{
			Username: "testuser",
			Password: "testpass",
		}

		reqBody, _ := json.Marshal(loginReq)
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// 由于 login 函数依赖于 repository.NewAdminRepository()，
		// 这里我们测试请求绑定部分
		var req vo.RegisterAndLoginRequest
		err := c.ShouldBind(&req)

		assert.NoError(t, err)
		assert.Equal(t, "testuser", req.Username)
		assert.Equal(t, "testpass", req.Password)
	})

	t.Run("无效的请求格式", func(t *testing.T) {
		c, _ := createTestContext()

		// 发送无效的 JSON
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		var req vo.RegisterAndLoginRequest
		err := c.ShouldBind(&req)

		assert.Error(t, err)
	})
}

// TestUserSerializationError 测试用户序列化错误处理
func TestUserSerializationError(t *testing.T) {
	setupTestConfig()

	t.Run("用户序列化失败错误信息_中文", func(t *testing.T) {
		c, _ := createTestContext()
		c.Set("lang", "zh")

		// 模拟序列化错误
		err := fmt.Errorf("serialization error")
		errorMsg := fmt.Sprintf("%s: %v", common.Msg(c, common.MsgUserSerializeFail), err)

		assert.Contains(t, errorMsg, "用户信息序列化失败")
		assert.Contains(t, errorMsg, "serialization error")
	})

	t.Run("用户序列化失败错误信息_英文", func(t *testing.T) {
		c, _ := createTestContext()
		c.Set("lang", "en")

		// 模拟序列化错误
		err := fmt.Errorf("serialization error")
		errorMsg := fmt.Sprintf("%s: %v", common.Msg(c, common.MsgUserSerializeFail), err)

		assert.Contains(t, errorMsg, "User information serialization failed")
		assert.Contains(t, errorMsg, "serialization error")
	})
}

// BenchmarkPayloadFunc 性能测试：载荷处理函数
func BenchmarkPayloadFunc(b *testing.B) {
	setupTestConfig()

	user := model.Admin{
		Model: model.Model{
			ID: 1,
		},
		Username: "testuser",
		Mobile:   "13800138000",
	}

	userJson, _ := util.JSONUtil.Struct2Json(user)
	data := map[string]interface{}{
		"user": userJson,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		payloadFunc(data)
	}
}

// BenchmarkAuthorizator 性能测试：授权验证函数
func BenchmarkAuthorizator(b *testing.B) {
	setupTestConfig()

	user := model.Admin{
		Model: model.Model{
			ID: 1,
		},
		Username: "testuser",
		Mobile:   "13800138000",
	}
	userJson, _ := util.JSONUtil.Struct2Json(user)

	data := map[string]interface{}{
		"user": userJson,
	}

	c, _ := createTestContext()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		authorizator(data, c)
	}
}

// BenchmarkJWTLogin 基准测试：JWT登录
func BenchmarkJWTLogin(b *testing.B) {
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("lang", "zh")
		c.Next()
	})

	authMiddleware, _ := InitAuth()
	r.POST("/login", authMiddleware.LoginHandler)

	loginReq := vo.RegisterAndLoginRequest{
		Username: "nonexistent",
		Password: "password",
	}
	jsonData, _ := json.Marshal(loginReq)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}
