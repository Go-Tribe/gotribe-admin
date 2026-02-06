// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

// 系统配置，对应yml
// viper内置了mapstructure, yml文件用"-"区分单词, 转为驼峰方便

// 全局配置变量
var Conf = new(config)

type config struct {
	System     *SystemConfig    `mapstructure:"system" json:"system"`
	Logs       *LogsConfig      `mapstructure:"logs" json:"logs"`
	Database   *DatabaseConfig  `mapstructure:"database" json:"database"`
	Casbin     *CasbinConfig    `mapstructure:"casbin" json:"casbin"`
	Jwt        *JwtConfig       `mapstructure:"jwt" json:"jwt"`
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

	// 读取配置信息
	err = viper.ReadInConfig()
	// 热更新配置
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 将读取的配置信息保存至全局变量Conf
		if err := viper.Unmarshal(Conf); err != nil {
			// 配置热更新失败时记录错误日志，但不中断服务
			log.Printf("配置热更新失败: %v，保持使用旧配置", err)
			return
		}
		log.Printf("配置文件已重新加载: %s", e.Name)
	})

	if err != nil {
		panic(fmt.Errorf("读取配置文件失败:%s \n", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("初始化配置文件失败:%s \n", err))
	}

}

type SystemConfig struct {
	Mode          string `mapstructure:"mode" json:"mode"`
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
	Realm      string `mapstructure:"realm" json:"realm"`
	Key        string `mapstructure:"key" json:"key"`
	Timeout    int    `mapstructure:"timeout" json:"timeout"`
	MaxRefresh int    `mapstructure:"max-refresh" json:"maxRefresh"`
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
