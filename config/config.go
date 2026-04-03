// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

// 系统配置，对应yml
// viper内置了mapstructure, yml文件用"-"区分单词, 转为驼峰方便

// 全局配置变量
var Conf = newConfig()

type config struct {
	System     *SystemConfig    `mapstructure:"system" json:"system"`
	Logs       *LogsConfig      `mapstructure:"logs" json:"logs"`
	Database   *DatabaseConfig  `mapstructure:"database" json:"database"`
	Casbin     *CasbinConfig    `mapstructure:"casbin" json:"casbin"`
	Jwt        *JwtConfig       `mapstructure:"jwt" json:"jwt"`
	CORS       *CORSConfig      `mapstructure:"cors" json:"cors"`
	RateLimit  *RateLimitConfig `mapstructure:"rate-limit" json:"rateLimit"`
	UploadFile *UploadFile      `mapstructure:"upload-file" json:"uploadFile"`
	Jobs       *JobsConfig      `mapstructure:"jobs" json:"jobs"`
}

// 设置读取配置信息
func InitConfig() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("读取应用目录失败:%s \n", err))
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
	bindEnvs()

	// 读取配置信息
	err = viper.ReadInConfig()
	// 热更新配置
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		ensureConfigStructs()
		// 将读取的配置信息保存至全局变量Conf
		if err := viper.Unmarshal(Conf); err != nil {
			// 配置热更新失败时记录错误日志，但不中断服务
			log.Printf("配置热更新失败: %v，保持使用旧配置", err)
			return
		}
		normalizeConfig()
		log.Printf("配置文件已重新加载: %s", e.Name)
	})

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("未找到 config.yml，将使用环境变量和默认值启动")
		} else {
			panic(fmt.Errorf("读取配置文件失败:%s \n", err))
		}
	}

	ensureConfigStructs()
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("初始化配置文件失败:%s \n", err))
	}

	normalizeConfig()
}

func newConfig() *config {
	return &config{
		System:     &SystemConfig{},
		Logs:       &LogsConfig{},
		Database:   &DatabaseConfig{},
		Casbin:     &CasbinConfig{},
		Jwt:        &JwtConfig{},
		CORS:       &CORSConfig{},
		RateLimit:  &RateLimitConfig{},
		UploadFile: &UploadFile{},
		Jobs:       &JobsConfig{},
	}
}

func ensureConfigStructs() {
	if Conf == nil {
		Conf = newConfig()
		return
	}
	if Conf.System == nil {
		Conf.System = &SystemConfig{}
	}
	if Conf.Logs == nil {
		Conf.Logs = &LogsConfig{}
	}
	if Conf.Database == nil {
		Conf.Database = &DatabaseConfig{}
	}
	if Conf.Casbin == nil {
		Conf.Casbin = &CasbinConfig{}
	}
	if Conf.Jwt == nil {
		Conf.Jwt = &JwtConfig{}
	}
	if Conf.CORS == nil {
		Conf.CORS = &CORSConfig{}
	}
	if Conf.RateLimit == nil {
		Conf.RateLimit = &RateLimitConfig{}
	}
	if Conf.UploadFile == nil {
		Conf.UploadFile = &UploadFile{}
	}
	if Conf.Jobs == nil {
		Conf.Jobs = &JobsConfig{}
	}
}

type SystemConfig struct {
	Mode          string `mapstructure:"mode" json:"mode"`
	Host          string `mapstructure:"host" json:"host"`
	UrlPathPrefix string `mapstructure:"url-path-prefix" json:"urlPathPrefix"`
	Port          int    `mapstructure:"port" json:"port"`
	InitData      bool   `mapstructure:"init-data" json:"initData"`
	CDNDomain     string `mapstructure:"cdn-domain" json:"CDNDomain"`
	EnableMigrate bool   `mapstructure:"enable-migrate" json:"enableMigrate"`
	EnableOss     bool   `mapstructure:"enable-oss" json:"enableOss"`
}

type LogsConfig struct {
	Level      zapcore.Level `mapstructure:"level" json:"level"`
	Path       string        `mapstructure:"path" json:"path"`
	MaxSize    int           `mapstructure:"max-size" json:"maxSize"`
	MaxBackups int           `mapstructure:"max-backups" json:"maxBackups"`
	MaxAge     int           `mapstructure:"max-age" json:"maxAge"`
	Compress   bool          `mapstructure:"compress" json:"compress"`
}

type DatabaseConfig struct {
	Type      string `mapstructure:"type" json:"type"`
	Username  string `mapstructure:"username" json:"username"`
	Password  string `mapstructure:"password" json:"password"`
	Database  string `mapstructure:"database" json:"database"`
	Host      string `mapstructure:"host" json:"host"`
	Port      int    `mapstructure:"port" json:"port"`
	Query     string `mapstructure:"query" json:"query"`
	LogMode   bool   `mapstructure:"log-mode" json:"logMode"`
	Charset   string `mapstructure:"charset" json:"charset"`
	Collation string `mapstructure:"collation" json:"collation"`
	SSLMode   string `mapstructure:"sslmode" json:"sslmode"`
}

type CasbinConfig struct {
	ModelPath string `mapstructure:"model-path" json:"modelPath"`
}

type JwtConfig struct {
	Realm       string `mapstructure:"realm" json:"realm"`
	Key         string `mapstructure:"key" json:"key"`
	Timeout     int    `mapstructure:"timeout" json:"timeout"`
	MaxRefresh  int    `mapstructure:"max-refresh" json:"maxRefresh"`
	TokenLookup string `mapstructure:"token-lookup" json:"tokenLookup"`
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed-origins" json:"allowedOrigins"`
	AllowCredentials bool     `mapstructure:"allow-credentials" json:"allowCredentials"`
	MaxAge           int      `mapstructure:"max-age" json:"maxAge"`
}

type RateLimitConfig struct {
	FillInterval int64 `mapstructure:"fill-interval" json:"fillInterval"`
	Capacity     int64 `mapstructure:"capacity" json:"capacity"`
}

type UploadFile struct {
	// Provider 上传服务商: qiniu(七牛), oss(阿里云OSS), s3(亚马逊S3，预留)
	Provider  string `mapstructure:"provider" json:"provider"`
	Accesskey string `mapstructure:"access-key" json:"accesskey"`
	Secretkey string `mapstructure:"secret-key" json:"secretkey"`
	Bucket    string `mapstructure:"bucket" json:"bucket"`
	Endpoint  string `mapstructure:"endpoint" json:"endpoint"`
}

// GetUploadProvider 返回当前生效的上传服务商。若未配置 provider 则按 system.enable-oss 兼容：true=oss, false=qiniu
func (u *UploadFile) GetUploadProvider(enableOss bool) string {
	if u != nil && u.Provider != "" {
		return u.Provider
	}
	if enableOss {
		return "oss"
	}
	return "qiniu"
}

type JobsConfig struct {
	Enabled bool                 `mapstructure:"enabled" json:"enabled"`
	List    map[string]JobConfig `mapstructure:"list" json:"list"`
}

type JobConfig struct {
	Name        string        `mapstructure:"name" json:"name"`
	Description string        `mapstructure:"description" json:"description"`
	Schedule    string        `mapstructure:"schedule" json:"schedule"`
	Enabled     bool          `mapstructure:"enabled" json:"enabled"`
	Timeout     time.Duration `mapstructure:"timeout" json:"timeout"`
	RetryCount  int           `mapstructure:"retry-count" json:"retryCount"`
}

func bindEnvs() {
	mustBindEnv("system.mode", "SYSTEM_MODE", "GIN_MODE")
	mustBindEnv("system.host", "SYSTEM_HOST", "HOST")
	mustBindEnv("system.port", "SYSTEM_PORT", "PORT")
	mustBindEnv("system.url-path-prefix", "SYSTEM_URL_PATH_PREFIX", "URL_PATH_PREFIX")
	mustBindEnv("system.init-data", "SYSTEM_INIT_DATA")
	mustBindEnv("system.cdn-domain", "SYSTEM_CDN_DOMAIN", "CDN_DOMAIN")
	mustBindEnv("system.enable-migrate", "SYSTEM_ENABLE_MIGRATE", "ENABLE_MIGRATE")
	mustBindEnv("system.enable-oss", "SYSTEM_ENABLE_OSS", "ENABLE_OSS")

	mustBindEnv("logs.level", "LOGS_LEVEL")
	mustBindEnv("logs.path", "LOGS_PATH")
	mustBindEnv("logs.max-size", "LOGS_MAX_SIZE")
	mustBindEnv("logs.max-backups", "LOGS_MAX_BACKUPS")
	mustBindEnv("logs.max-age", "LOGS_MAX_AGE")
	mustBindEnv("logs.compress", "LOGS_COMPRESS")

	mustBindEnv("database.type", "DATABASE_TYPE", "DB_TYPE")
	mustBindEnv("database.username", "DATABASE_USERNAME", "DB_USERNAME")
	mustBindEnv("database.password", "DATABASE_PASSWORD", "DB_PASSWORD")
	mustBindEnv("database.database", "DATABASE_DATABASE", "DB_DATABASE")
	mustBindEnv("database.host", "DATABASE_HOST", "DB_HOST")
	mustBindEnv("database.port", "DATABASE_PORT", "DB_PORT")
	mustBindEnv("database.query", "DATABASE_QUERY", "DB_QUERY")
	mustBindEnv("database.log-mode", "DATABASE_LOG_MODE", "DB_LOG_MODE")
	mustBindEnv("database.charset", "DATABASE_CHARSET", "DB_CHARSET")
	mustBindEnv("database.collation", "DATABASE_COLLATION", "DB_COLLATION")
	mustBindEnv("database.sslmode", "DATABASE_SSLMODE", "DB_SSLMODE")

	mustBindEnv("casbin.model-path", "RBAC_MODEL_PATH", "CASBIN_MODEL_PATH")

	mustBindEnv("jwt.realm", "JWT_REALM")
	mustBindEnv("jwt.key", "JWT_KEY")
	mustBindEnv("jwt.timeout", "JWT_TIMEOUT")
	mustBindEnv("jwt.max-refresh", "JWT_MAX_REFRESH")
	mustBindEnv("jwt.token-lookup", "JWT_TOKEN_LOOKUP")

	mustBindEnv("rate-limit.fill-interval", "RATE_LIMIT_FILL_INTERVAL")
	mustBindEnv("rate-limit.capacity", "RATE_LIMIT_CAPACITY")

	mustBindEnv("upload-file.provider", "UPLOAD_FILE_PROVIDER")
	mustBindEnv("upload-file.access-key", "UPLOAD_FILE_ACCESS_KEY")
	mustBindEnv("upload-file.secret-key", "UPLOAD_FILE_SECRET_KEY")
	mustBindEnv("upload-file.bucket", "UPLOAD_FILE_BUCKET")
	mustBindEnv("upload-file.endpoint", "UPLOAD_FILE_ENDPOINT")

	mustBindEnv("cors.max-age", "CORS_MAX_AGE")
	mustBindEnv("cors.allow-credentials", "CORS_ALLOW_CREDENTIALS")
}

func mustBindEnv(key string, envs ...string) {
	args := append([]string{key}, envs...)
	if err := viper.BindEnv(args...); err != nil {
		panic(fmt.Errorf("绑定环境变量失败(%s): %w", key, err))
	}
}

func normalizeConfig() {
	if Conf.System != nil {
		if Conf.System.Mode == "" {
			Conf.System.Mode = "debug"
		}
		if Conf.System.UrlPathPrefix == "" {
			Conf.System.UrlPathPrefix = "api"
		}
		if Conf.System.Port == 0 {
			Conf.System.Port = 8088
		}
	}
	if Conf.System != nil && Conf.System.Host == "" {
		Conf.System.Host = "0.0.0.0"
	}

	if Conf.Logs != nil && Conf.Logs.Path == "" {
		Conf.Logs.Path = "logs"
	}

	if Conf.Database != nil {
		if Conf.Database.Type == "" {
			Conf.Database.Type = "mysql"
		}
		if Conf.Database.Charset == "" {
			Conf.Database.Charset = "utf8mb4"
		}
		if Conf.Database.Collation == "" {
			Conf.Database.Collation = "utf8mb4_general_ci"
		}
		if Conf.Database.SSLMode == "" {
			Conf.Database.SSLMode = "disable"
		}
	}

	if Conf.Jwt != nil && Conf.Jwt.TokenLookup == "" {
		Conf.Jwt.TokenLookup = "header: Authorization, query: token"
	}

	if Conf.CORS == nil {
		Conf.CORS = &CORSConfig{}
	}
	if Conf.CORS.MaxAge <= 0 {
		Conf.CORS.MaxAge = 600
	}

	if Conf.RateLimit != nil {
		if Conf.RateLimit.FillInterval <= 0 {
			Conf.RateLimit.FillInterval = 50
		}
		if Conf.RateLimit.Capacity <= 0 {
			Conf.RateLimit.Capacity = 200
		}
	}
}
