// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package common

import (
	"fmt"
	"gotribe-admin/config"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// 全局CasbinEnforcer
var CasbinEnforcer *casbin.Enforcer

// 初始化casbin策略管理器
func InitCasbinEnforcer() {
	e, err := databaseCasbin()
	if err != nil {
		Log.Panicf("初始化Casbin失败：%v", err)
		panic(fmt.Sprintf("初始化Casbin失败：%v", err))
	}

	CasbinEnforcer = e
	Log.Info("初始化Casbin完成!")
}

func databaseCasbin() (*casbin.Enforcer, error) {
	a, err := gormadapter.NewAdapterByDB(DB)
	if err != nil {
		return nil, err
	}

	// 外部路径优先：环境变量 > 配置文件
	modelPath := os.Getenv("RBAC_MODEL_PATH")
	if modelPath == "" {
		modelPath = config.Conf.Casbin.ModelPath
	}

	var e *casbin.Enforcer
	// 如果外部文件存在则优先使用
	if modelPath != "" {
		if _, statErr := os.Stat(modelPath); statErr == nil {
			Log.Infof("加载外部 RBAC 模型: %s", modelPath)
			e, err = casbin.NewEnforcer(modelPath, a)
			if err != nil {
				Log.Warnf("加载外部 RBAC 模型失败(%s): %v，使用内置默认", modelPath, err)
				e = nil
			}
		} else {
			Log.Warnf("外部 RBAC 模型不可用(%s): %v，使用内置默认", modelPath, statErr)
		}
	}

	// 兜底：使用内置默认模型
	if e == nil {
		if embeddedRBACModel == "" {
			return nil, fmt.Errorf("内置 RBAC 模型为空，无法初始化")
		}
		m := model.NewModel()
		if err := m.LoadModelFromText(embeddedRBACModel); err != nil {
			return nil, fmt.Errorf("加载内置 RBAC 模型失败: %w", err)
		}
		Log.Info("加载内置默认 RBAC 模型")
		e, err = casbin.NewEnforcer(m, a)
		if err != nil {
			return nil, err
		}
	}

	if err = e.LoadPolicy(); err != nil {
		return nil, err
	}
	return e, nil
}
