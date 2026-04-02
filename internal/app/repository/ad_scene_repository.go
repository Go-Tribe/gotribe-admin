// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/dengmengmian/ghelper/gconvert"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
)

type IAdSceneRepository interface {
	CreateAdScene(ctx context.Context, adScene *model.AdScene) error                              // 创建推广场景
	GetAdSceneByID(ctx context.Context, id uint) (model.AdScene, error)                           // 获取单个推广场景
	GetAdScenes(ctx context.Context, req *vo.AdSceneListRequest) ([]*model.AdScene, int64, error) // 获取推广场景列表
	UpdateAdScene(ctx context.Context, adScene *model.AdScene) error                              // 更新推广场景
	BatchDeleteAdSceneByIds(ctx context.Context, ids []uint) error                                // 批量删除
}

type AdSceneRepository struct {
}

// AdSceneRepository构造函数
func NewAdSceneRepository() IAdSceneRepository {
	return AdSceneRepository{}
}

// 获取单个推广场景
func (cr AdSceneRepository) GetAdSceneByID(ctx context.Context, id uint) (model.AdScene, error) {
	var adScene model.AdScene
	err := common.WithContext(ctx).DB().Where("id = ?", id).First(&adScene).Error
	return adScene, err
}

// 获取推广场景列表
func (cr AdSceneRepository) GetAdScenes(ctx context.Context, req *vo.AdSceneListRequest) ([]*model.AdScene, int64, error) {
	var list []*model.AdScene
	db := common.WithContext(ctx).DB().Model(&model.AdScene{}).Order("created_at DESC")

	projectID := strings.TrimSpace(req.ProjectID)
	if !gconvert.IsEmpty(projectID) {
		db = db.Where("project_id = ?", projectID)
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
	list, err = GetAdSceneOther(ctx, list)
	return list, total, err
}

// 获取推广场景其他信息
func GetAdSceneOther(ctx context.Context, adScenes []*model.AdScene) ([]*model.AdScene, error) {
	if len(adScenes) == 0 {
		return adScenes, nil
	}

	// 收集所有 ProjectID
	projectIDs := make([]string, 0, len(adScenes))
	for _, m := range adScenes {
		if m.ProjectID != "" {
			projectIDs = append(projectIDs, m.ProjectID)
		}
	}

	if len(projectIDs) == 0 {
		return adScenes, nil
	}

	// 批量查询
	var projects []*model.Project
	if err := common.WithContext(ctx).DB().Where("project_id IN ?", projectIDs).Find(&projects).Error; err != nil {
		return adScenes, err
	}

	// 建立映射
	projectMap := make(map[string]*model.Project)
	for _, project := range projects {
		projectMap[project.ProjectID] = project
	}

	// 赋值
	for _, m := range adScenes {
		if project, ok := projectMap[m.ProjectID]; ok {
			m.Project = project
		}
	}
	return adScenes, nil
}

// 创建推广场景
func (cr AdSceneRepository) CreateAdScene(ctx context.Context, adScene *model.AdScene) error {
	err := common.WithContext(ctx).DB().Create(adScene).Error
	return err
}

// 更新推广场景
func (cr AdSceneRepository) UpdateAdScene(ctx context.Context, adScene *model.AdScene) error {
	err := common.WithContext(ctx).DB().Model(adScene).Updates(adScene).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (cr AdSceneRepository) BatchDeleteAdSceneByIds(ctx context.Context, ids []uint) error {
	var adScenes []model.AdScene
	for _, id := range ids {
		// 根据ID获取标签
		adScene, err := cr.GetAdSceneByID(ctx, id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%d的推广场景", id)
		}
		adScenes = append(adScenes, adScene)
	}

	err := common.WithContext(ctx).DB().Unscoped().Delete(&adScenes).Error

	return err
}
