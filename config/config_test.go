// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func resetConfigForTest(t *testing.T) {
	t.Helper()

	viper.Reset()
	Conf = newConfig()
	t.Cleanup(func() {
		viper.Reset()
		Conf = newConfig()
	})
}

func chdirForConfigTest(t *testing.T) {
	t.Helper()

	oldWd, err := os.Getwd()
	require.NoError(t, err)

	tempDir := t.TempDir()
	require.NoError(t, os.Chdir(tempDir))
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(oldWd))
	})
}

func TestInitConfigUsesDefaultsWithoutConfigFile(t *testing.T) {
	resetConfigForTest(t)
	chdirForConfigTest(t)

	t.Setenv("SYSTEM_MODE", "")
	t.Setenv("SYSTEM_HOST", "")
	t.Setenv("SYSTEM_PORT", "")
	t.Setenv("DATABASE_TYPE", "")
	t.Setenv("JWT_TOKEN_LOOKUP", "")
	t.Setenv("RATE_LIMIT_FILL_INTERVAL", "")
	t.Setenv("RATE_LIMIT_CAPACITY", "")

	InitConfig()

	require.Equal(t, "debug", Conf.System.Mode)
	require.Equal(t, "0.0.0.0", Conf.System.Host)
	require.Equal(t, "api", Conf.System.UrlPathPrefix)
	require.Equal(t, 8088, Conf.System.Port)
	require.Equal(t, "mysql", Conf.Database.Type)
	require.Equal(t, "utf8mb4", Conf.Database.Charset)
	require.Equal(t, "utf8mb4_general_ci", Conf.Database.Collation)
	require.Equal(t, "disable", Conf.Database.SSLMode)
	require.Equal(t, "header: Authorization, query: token", Conf.Jwt.TokenLookup)
	require.EqualValues(t, 50, Conf.RateLimit.FillInterval)
	require.EqualValues(t, 200, Conf.RateLimit.Capacity)
	require.Equal(t, 600, Conf.CORS.MaxAge)
}

func TestInitConfigAppliesEnvironmentOverrides(t *testing.T) {
	resetConfigForTest(t)
	chdirForConfigTest(t)

	t.Setenv("SYSTEM_MODE", "release")
	t.Setenv("SYSTEM_HOST", "127.0.0.1")
	t.Setenv("SYSTEM_PORT", "9090")
	t.Setenv("SYSTEM_URL_PATH_PREFIX", "admin-api")
	t.Setenv("DATABASE_TYPE", "postgres")
	t.Setenv("DATABASE_PORT", "5432")
	t.Setenv("DATABASE_SSLMODE", "require")
	t.Setenv("JWT_TOKEN_LOOKUP", "header: Authorization")
	t.Setenv("RATE_LIMIT_FILL_INTERVAL", "25")
	t.Setenv("RATE_LIMIT_CAPACITY", "50")

	InitConfig()

	require.Equal(t, "release", Conf.System.Mode)
	require.Equal(t, "127.0.0.1", Conf.System.Host)
	require.Equal(t, 9090, Conf.System.Port)
	require.Equal(t, "admin-api", Conf.System.UrlPathPrefix)
	require.Equal(t, "postgres", Conf.Database.Type)
	require.Equal(t, 5432, Conf.Database.Port)
	require.Equal(t, "require", Conf.Database.SSLMode)
	require.Equal(t, "header: Authorization", Conf.Jwt.TokenLookup)
	require.EqualValues(t, 25, Conf.RateLimit.FillInterval)
	require.EqualValues(t, 50, Conf.RateLimit.Capacity)
}

func TestUploadProviderCompatibilityFallback(t *testing.T) {
	require.Equal(t, "oss", (&UploadFile{}).GetUploadProvider(true))
	require.Equal(t, "qiniu", (&UploadFile{}).GetUploadProvider(false))
	require.Equal(t, "s3", (&UploadFile{Provider: "s3"}).GetUploadProvider(false))
}
