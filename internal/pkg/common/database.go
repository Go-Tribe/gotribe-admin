// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package common

import (
	"database/sql"
	"fmt"
	"gotribe-admin/config"
	"gotribe-admin/internal/pkg/model/migrate"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 全局数据库变量
var DB *gorm.DB

// getGormConfig 获取 GORM 通用配置
func getGormConfig() *gorm.Config {
	return &gorm.Config{
		// 禁用外键(指定外键时不会在数据库创建真实的外键约束)
		DisableForeignKeyConstraintWhenMigrating: true,
		// 使用单数表名
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}

// buildPostgresDSN 构建 PostgreSQL 连接字符串
func buildPostgresDSN(maskPassword bool) string {
	db := config.Conf.Database

	// 处理 SSLMode 默认值
	sslMode := db.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	// 处理 Query 参数，如果为空则不添加
	queryPart := ""
	if db.Query != "" {
		queryPart = " " + db.Query
	}

	password := db.Password
	if maskPassword {
		password = "******"
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s%s",
		db.Host, db.Username, password, db.Database, db.Port, sslMode, queryPart)
}

// buildMySQLDSN 构建 MySQL 连接字符串
func buildMySQLDSN(maskPassword bool) string {
	db := config.Conf.Database

	// 处理 Query 参数，如果为空则不添加 & 前缀
	queryPart := ""
	if db.Query != "" {
		queryPart = "&" + db.Query
	}

	password := db.Password
	if maskPassword {
		password = "******"
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s%s",
		db.Username, password, db.Host, db.Port, db.Database,
		db.Charset, db.Collation, queryPart)
}

// setupDatabaseLogger 配置数据库日志
func setupDatabaseLogger(db *gorm.DB) {
	if !config.Conf.Database.LogMode {
		return
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db.Debug()
	db.Logger = newLogger
}

// openDatabase 打开数据库连接
func openDatabase(dsn string) (*gorm.DB, error) {
	dbType := config.Conf.Database.Type
	gormConfig := getGormConfig()

	if dbType == "postgres" {
		return gorm.Open(postgres.Open(dsn), gormConfig)
	}
	// 默认使用 MySQL
	return gorm.Open(mysql.Open(dsn), gormConfig)
}

// 初始化数据库
func InitDatabase() {
	var dsn, showDsn string

	// 根据数据库类型构建连接字符串
	if config.Conf.Database.Type == "postgres" {
		dsn = buildPostgresDSN(false)
		showDsn = buildPostgresDSN(true)
	} else {
		dsn = buildMySQLDSN(false)
		showDsn = buildMySQLDSN(true)
	}

	// 打开数据库连接
	db, err := openDatabase(dsn)
	if err != nil {
		Log.Panicf("初始化数据库异常: %v", err)
		panic(fmt.Errorf("初始化数据库异常: %v", err))
	}

	// 配置数据库日志
	setupDatabaseLogger(db)

	// 配置连接池参数，避免连接频繁关闭和重连导致的日志输出
	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil {
		Log.Panicf("获取数据库连接失败: %v", err)
		panic(fmt.Errorf("获取数据库连接失败: %v", err))
	}

	// 设置连接池参数
	// SetMaxOpenConns: 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)
	// SetMaxIdleConns: 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetConnMaxLifetime: 设置连接可复用的最大时间（应该小于 MySQL 的 wait_timeout）
	// MySQL 默认 wait_timeout 是 8 小时，这里设置为 7 小时，避免连接被 MySQL 主动关闭
	sqlDB.SetConnMaxLifetime(7 * time.Hour)
	// SetConnMaxIdleTime: 设置连接的最大空闲时间，超过此时间的空闲连接会被关闭
	sqlDB.SetConnMaxIdleTime(1 * time.Hour)

	// 全局DB赋值
	DB = db

	// 自动迁移表结构
	if config.Conf.System.EnableMigrate {
		migrate.DBAutoMigrate(DB)
		Log.Infof("数据库迁移完成! dsn: %s", showDsn)
	}
}
