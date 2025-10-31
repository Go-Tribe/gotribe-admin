// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package common

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *RepositoryError
		expected string
	}{
		{
			name:     "简单错误消息",
			err:      ErrUserNotFound,
			expected: "用户不存在",
		},
		{
			name:     "带参数的错误消息",
			err:      NewUserNotFoundByIDError(123),
			expected: "未获取到ID为123的用户",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestRepositoryError_GetLocalizedMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		err      *RepositoryError
		lang     string
		expected string
	}{
		{
			name:     "中文-用户不存在",
			err:      ErrUserNotFound,
			lang:     "zh",
			expected: "用户不存在",
		},
		{
			name:     "英文-用户不存在",
			err:      ErrUserNotFound,
			lang:     "en",
			expected: "User not found",
		},
		{
			name:     "中文-用户被禁用",
			err:      ErrUserDisabled,
			lang:     "zh",
			expected: "用户被禁用",
		},
		{
			name:     "英文-用户被禁用",
			err:      ErrUserDisabled,
			lang:     "en",
			expected: "User is disabled",
		},
		{
			name:     "中文-密码错误",
			err:      ErrPasswordIncorrect,
			lang:     "zh",
			expected: "密码错误",
		},
		{
			name:     "英文-密码错误",
			err:      ErrPasswordIncorrect,
			lang:     "en",
			expected: "Incorrect password",
		},
		{
			name:     "中文-带参数的错误",
			err:      NewUserNotFoundByIDError(456),
			lang:     "zh",
			expected: "未获取到ID为456的用户",
		},
		{
			name:     "英文-带参数的错误",
			err:      NewUserNotFoundByIDError(456),
			lang:     "en",
			expected: "User not found by ID 456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试上下文
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// 设置语言
			c.Set("lang", tt.lang)

			// 获取本地化消息
			result := tt.err.GetLocalizedMessage(c)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRepositoryError_GetLocalizedMessage_DefaultLanguage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试默认语言（无语言设置时应该返回中文）
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// 不设置语言，应该默认为中文

	result := ErrUserNotFound.GetLocalizedMessage(c)
	assert.Equal(t, "用户不存在", result)
}

func TestRepositoryError_GetLocalizedMessage_UnsupportedLanguage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试不支持的语言（应该回退到中文）
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("lang", "fr") // 法语，不支持

	result := ErrUserNotFound.GetLocalizedMessage(c)
	assert.Equal(t, "用户不存在", result) // 应该回退到中文
}

func TestNewUserNotFoundByIDError(t *testing.T) {
	err := NewUserNotFoundByIDError(789)

	assert.Equal(t, "user_not_found_by_id", err.Code)
	assert.Equal(t, "未获取到ID为789的用户", err.Error())
	assert.Equal(t, []interface{}{uint(789)}, err.Args)
}

// 基准测试
func BenchmarkRepositoryError_GetLocalizedMessage(b *testing.B) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("lang", "zh")

	err := ErrUserNotFound

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.GetLocalizedMessage(c)
	}
}

func BenchmarkNewUserNotFoundByIDError(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewUserNotFoundByIDError(uint(i))
	}
}
