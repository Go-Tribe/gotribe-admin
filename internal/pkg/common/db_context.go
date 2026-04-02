// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package common

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// DBContext 封装带 Context 的数据库操作
type DBContext struct {
	ctx context.Context
	db  *gorm.DB
}

// WithContext 创建带 Context 的数据库操作实例
func WithContext(ctx context.Context) *DBContext {
	if ctx == nil {
		ctx = context.Background()
	}
	return &DBContext{
		ctx: ctx,
		db:  DB.WithContext(ctx),
	}
}

// WithTimeout 创建带超时的数据库操作上下文
func WithTimeout(parent context.Context, timeout time.Duration) (*DBContext, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithTimeout(parent, timeout)
	return &DBContext{
		ctx: ctx,
		db:  DB.WithContext(ctx),
	}, cancel
}

// DB 返回 GORM DB 实例
func (d *DBContext) DB() *gorm.DB {
	return d.db
}

// Context 返回当前上下文
func (d *DBContext) Context() context.Context {
	return d.ctx
}

// CheckHealth 检查数据库连接健康状态
func CheckHealth(ctx context.Context) error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
