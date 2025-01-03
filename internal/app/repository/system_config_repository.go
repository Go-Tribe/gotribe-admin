// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
)

type ISystemConfigRepository interface {
	GetSystemConfig() (model.SystemConfig, error)              // 获取单个标签
	UpdateSystemConfig(systemConfig *model.SystemConfig) error // 更新标签
}

type SystemConfigRepository struct {
}

// SystemConfigRepository构造函数
func NewSystemConfigRepository() ISystemConfigRepository {
	return SystemConfigRepository{}
}

// 获取单个
func (tr SystemConfigRepository) GetSystemConfig() (model.SystemConfig, error) {
	var systemConfig model.SystemConfig
	err := common.DB.First(&systemConfig).Error
	return systemConfig, err
}

// 更新
func (tr SystemConfigRepository) UpdateSystemConfig(systemConfig *model.SystemConfig) error {
	err := common.DB.Model(systemConfig).Updates(systemConfig).Error
	if err != nil {
		return err
	}

	return err
}
