// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package common

import (
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

// 初始化数据库
func InitDatabase() {
	var dsn string
	var showDsn string

	if config.Conf.Database.Type == "postgres" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s %s",
			config.Conf.Database.Host,
			config.Conf.Database.Username,
			config.Conf.Database.Password,
			config.Conf.Database.Database,
			config.Conf.Database.Port,
			config.Conf.Database.SSLMode,
			config.Conf.Database.Query,
		)
		showDsn = fmt.Sprintf("host=%s user=%s password=****** dbname=%s port=%d sslmode=%s %s",
			config.Conf.Database.Host,
			config.Conf.Database.Username,
			config.Conf.Database.Database,
			config.Conf.Database.Port,
			config.Conf.Database.SSLMode,
			config.Conf.Database.Query,
		)
	} else {
		// MySQL 连接字符串
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
			config.Conf.Database.Username,
			config.Conf.Database.Password,
			config.Conf.Database.Host,
			config.Conf.Database.Port,
			config.Conf.Database.Database,
			config.Conf.Database.Charset,
			config.Conf.Database.Collation,
			config.Conf.Database.Query,
		)
		showDsn = fmt.Sprintf(
			"%s:******@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
			config.Conf.Database.Username,
			config.Conf.Database.Host,
			config.Conf.Database.Port,
			config.Conf.Database.Database,
			config.Conf.Database.Charset,
			config.Conf.Database.Collation,
			config.Conf.Database.Query,
		)
	}

	var db *gorm.DB
	var err error

	if config.Conf.Database.Type == "postgres" {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			// 禁用外键(指定外键时不会在数据库创建真实的外键约束)
			DisableForeignKeyConstraintWhenMigrating: true,
			// 使用单数表名
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	} else {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			// 禁用外键(指定外键时不会在数据库创建真实的外键约束)
			DisableForeignKeyConstraintWhenMigrating: true,
			// 使用单数表名
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	}
	if err != nil {
		Log.Panicf("初始化数据库异常: %v", err)
		panic(fmt.Errorf("初始化数据库异常: %v", err))
	}

	// 开启数据库日志
	if config.Conf.Database.LogMode {
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
	// 全局DB赋值
	DB = db
	// 自动迁移表结构
	if config.Conf.System.EnableMigrate {
		migrate.DBAutoMigrate(DB)
		Log.Infof("数据库迁移完成! dsn: %s", showDsn)
	}
}
