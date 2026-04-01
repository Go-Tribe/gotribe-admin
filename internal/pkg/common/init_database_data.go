// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package common

import (
	"fmt"
	"gotribe-admin/config"
	"gotribe-admin/internal/pkg/common/seeder"
)

// 初始化数据库数据
func InitData() {
	// 是否初始化数据
	if !config.Conf.System.InitData {
		return
	}

	// 注册所有种子
	registerAllSeeders()

	// 运行所有种子
	if err := seeder.RunSeeders(DB); err != nil {
		Log.Errorf("数据库种子执行失败: %v", err)
	}

	// 种子执行后重置 PostgreSQL 序列，避免自增主键冲突
	if config.Conf.Database.Type == "postgres" {
		if err := resetPostgresSequences(); err != nil {
			Log.Errorf("重置 PostgreSQL 序列失败: %v", err)
		} else {
			Log.Info("PostgreSQL 序列已重置")
		}
	}
}

// registerAllSeeders 注册所有种子
func registerAllSeeders() {
	// 基础数据种子
	seeder.RegisterSeeder(seeder.NewRoleSeeder())
	seeder.RegisterSeeder(seeder.NewAdminSeeder())
	seeder.RegisterSeeder(seeder.NewMenuSeeder())
	seeder.RegisterSeeder(seeder.NewApiSeeder())
	seeder.RegisterSeeder(seeder.NewSystemConfigSeeder())

	// 内容管理种子
	seeder.RegisterSeeder(seeder.NewCategorySeeder())
	seeder.RegisterSeeder(seeder.NewTagSeeder())
	seeder.RegisterSeeder(seeder.NewPostSeeder())

	// 项目管理种子
	seeder.RegisterSeeder(seeder.NewProjectSeeder())
	seeder.RegisterSeeder(seeder.NewUserSeeder())

	// 可以继续添加其他种子...
}

// resetPostgresSequences 重置所有自增序列，解决种子数据写死 ID 后导致序列不同步的问题
type sequenceInfo struct {
	TableName    string
	ColumnName   string
	SequenceName string
}

func resetPostgresSequences() error {
	// 查询 public schema 下所有带自增序列的表和列
	var sequences []sequenceInfo
	sql := `
		SELECT 
			c.relname AS table_name,
			a.attname AS column_name,
			pg_get_serial_sequence(c.relname, a.attname) AS sequence_name
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		JOIN pg_attribute a ON a.attrelid = c.oid
		WHERE c.relkind = 'r'
		  AND n.nspname = 'public'
		  AND a.attnum > 0
		  AND NOT a.attisdropped
		  AND pg_get_serial_sequence(c.relname, a.attname) IS NOT NULL
	`
	if err := DB.Raw(sql).Scan(&sequences).Error; err != nil {
		return fmt.Errorf("查询自增序列失败: %w", err)
	}

	for _, seq := range sequences {
		resetSQL := fmt.Sprintf(
			"SELECT setval('%s', COALESCE(MAX(%s), 0) + 1, false) FROM %s",
			seq.SequenceName, seq.ColumnName, seq.TableName,
		)
		if err := DB.Exec(resetSQL).Error; err != nil {
			Log.Warnf("重置序列失败 %s.%s (%s): %v", seq.TableName, seq.ColumnName, seq.SequenceName, err)
		}
	}

	return nil
}
