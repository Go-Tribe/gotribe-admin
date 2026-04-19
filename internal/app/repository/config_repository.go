// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"
	"fmt"
	"github.com/dengmengmian/ghelper/gconvert"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type IConfigRepository interface {
	CreateConfig(ctx context.Context, config *model.Config) error                              // 创建配置
	GetConfigByConfigID(ctx context.Context, configID string) (model.Config, error)            // 获取单个配置
	GetConfigs(ctx context.Context, req *vo.ConfigListRequest) ([]*model.Config, int64, error) // 获取配置列表
	UpdateConfig(ctx context.Context, config *model.Config) error                              // 更新配置
	BatchDeleteConfigByIds(ctx context.Context, ids []string) error                            // 批量删除
}

type ConfigRepository struct {
}

// ConfigRepository构造函数
func NewConfigRepository() IConfigRepository {
	return ConfigRepository{}
}

// 获取单个配置
func (cr ConfigRepository) GetConfigByConfigID(ctx context.Context, configID string) (model.Config, error) {
	var config model.Config
	err := common.WithContext(ctx).DB().Where("config_id = ?", configID).First(&config).Error
	return config, err
}

// 获取配置列表
func (cr ConfigRepository) GetConfigs(ctx context.Context, req *vo.ConfigListRequest) ([]*model.Config, int64, error) {
	var list []*model.Config
	db := common.WithContext(ctx).DB().Model(&model.Config{}).Order("created_at DESC")

	title := strings.TrimSpace(req.Title)
	if !gconvert.IsEmpty(title) {
		db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title))
	}
	configID := strings.TrimSpace(req.ConfigID)
	if !gconvert.IsEmpty(configID) {
		db = db.Where("config_id = ?", configID)
	}
	if req.ProjectId > 0 {
		db = db.Where("project_id = ?", req.ProjectId)
	}
	reqType := req.Type
	if reqType != 0 {
		db = db.Where("type = ?", reqType)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return GetConfigOther(list), total, err
}

// 获取配置其他信息
func GetConfigOther(configs []*model.Config) []*model.Config {
	if len(configs) == 0 {
		return configs
	}

	// 批量收集并去重 project_id，避免 N+1 查询
	projectIDSet := make(map[uint]struct{}, len(configs))
	projectIDs := make([]uint, 0, len(configs))
	for _, cfg := range configs {
		if cfg.ProjectId == 0 {
			continue
		}
		if _, exists := projectIDSet[cfg.ProjectId]; exists {
			continue
		}
		projectIDSet[cfg.ProjectId] = struct{}{}
		projectIDs = append(projectIDs, cfg.ProjectId)
	}

	if len(projectIDs) == 0 {
		return configs
	}

	var projects []model.Project
	if err := common.DB.Where("id IN (?)", projectIDs).Find(&projects).Error; err != nil {
		return configs
	}

	projectMap := make(map[uint]*model.Project, len(projects))
	for i := range projects {
		projectMap[projects[i].ID] = &projects[i]
	}

	for _, m := range configs {
		if project, ok := projectMap[m.ProjectId]; ok {
			m.Project = project
		}
	}
	return configs
}

// 创建配置
func (cr ConfigRepository) CreateConfig(ctx context.Context, config *model.Config) error {
	err := common.WithContext(ctx).DB().Create(config).Error
	return err
}

// 更新配置
func (cr ConfigRepository) UpdateConfig(ctx context.Context, config *model.Config) error {
	err := common.WithContext(ctx).DB().Model(config).Updates(config).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (cr ConfigRepository) BatchDeleteConfigByIds(ctx context.Context, ids []string) error {
	var configs []model.Config
	for _, id := range ids {
		// 根据ID获取标签
		config, err := cr.GetConfigByConfigID(ctx, id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%s的配置", id)
		}
		configs = append(configs, config)
	}

	err := common.WithContext(ctx).DB().Unscoped().Delete(&configs).Error

	return err
}
