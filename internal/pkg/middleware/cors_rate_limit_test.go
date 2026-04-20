// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gotribe-admin/config"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/pkg/api/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupMiddlewareTest(t *testing.T) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	config.Conf.System = &config.SystemConfig{Mode: gin.DebugMode}
	config.Conf.CORS = &config.CORSConfig{MaxAge: 600}
	if common.Log == nil {
		logger, _ := zap.NewDevelopment()
		common.Log = logger.Sugar()
	}
}

func TestCORSMiddlewareAllowsConfiguredOrigin(t *testing.T) {
	setupMiddlewareTest(t)
	config.Conf.CORS.AllowedOrigins = []string{"https://admin.example.com"}
	config.Conf.CORS.AllowCredentials = true

	router := gin.New()
	router.Use(CORSMiddleware())
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "https://admin.example.com")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "https://admin.example.com", w.Header().Get("Access-Control-Allow-Origin"))
	require.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	require.Equal(t, "Origin", w.Header().Get("Vary"))
}

func TestCORSMiddlewareRejectsUnknownPreflightOrigin(t *testing.T) {
	setupMiddlewareTest(t)
	config.Conf.System.Mode = gin.ReleaseMode
	config.Conf.CORS.AllowedOrigins = []string{"https://admin.example.com"}

	router := gin.New()
	router.Use(CORSMiddleware())
	router.OPTIONS("/ping", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodOptions, "/ping", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
	require.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestRateLimitMiddlewareReturnsTooManyRequests(t *testing.T) {
	setupMiddlewareTest(t)

	router := gin.New()
	router.Use(RateLimitMiddleware(time.Hour, 1))
	router.GET("/limited", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	first := httptest.NewRecorder()
	router.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/limited", nil))
	require.Equal(t, http.StatusOK, first.Code)

	second := httptest.NewRecorder()
	router.ServeHTTP(second, httptest.NewRequest(http.MethodGet, "/limited", nil))
	require.Equal(t, http.StatusTooManyRequests, second.Code)
	require.Contains(t, second.Body.String(), "访问限流")
}

func TestShouldSkipOperationLogPaths(t *testing.T) {
	require.True(t, shouldSkipLog(""))
	require.True(t, shouldSkipLog("/assets/app.js"))
	require.True(t, shouldSkipLog("/favicon.ico"))
	require.False(t, shouldSkipLog("/post"))
	require.Equal(t, http.StatusTooManyRequests, response.GetHTTPStatus(response.CodeTooManyRequests))
}
