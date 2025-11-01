// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"fmt"
	"github.com/dengmengmian/ghelper/gconvert"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type IAdSceneRepository interface {
	CreateAdScene(adScene *model.AdScene) error                              // 创建推广场景
	GetAdSceneByAdSceneID(adSceneID string) (model.AdScene, error)           // 获取单个推广场景
	GetAdScenes(req *vo.AdSceneListRequest) ([]*model.AdScene, int64, error) // 获取推广场景列表
	UpdateAdScene(adScene *model.AdScene) error                              // 更新推广场景
	BatchDeleteAdSceneByIds(ids []string) error                              // 批量删除
}

type AdSceneRepository struct {
}

// AdSceneRepository构造函数
func NewAdSceneRepository() IAdSceneRepository {
	return AdSceneRepository{}
}

// 获取单个推广场景
func (cr AdSceneRepository) GetAdSceneByAdSceneID(adSceneID string) (model.AdScene, error) {
	var adScene model.AdScene
	err := common.DB.Where("ad_scene_id = ?", adSceneID).First(&adScene).Error
	return adScene, err
}

// 获取推广场景列表
func (cr AdSceneRepository) GetAdScenes(req *vo.AdSceneListRequest) ([]*model.AdScene, int64, error) {
	var list []*model.AdScene
	db := common.DB.Model(&model.AdScene{}).Order("created_at DESC")

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
	return GetAdSceneOther(list), total, err
}

// 获取推广场景其他信息
func GetAdSceneOther(adScenes []*model.AdScene) []*model.AdScene {
	for _, m := range adScenes {
		var project *model.Project
		_ = common.DB.Where("project_id = ?", m.ProjectID).First(&project).Error
		m.Project = project
	}
	return adScenes
}

// 创建推广场景
func (cr AdSceneRepository) CreateAdScene(adScene *model.AdScene) error {
	err := common.DB.Create(adScene).Error
	return err
}

// 更新推广场景
func (cr AdSceneRepository) UpdateAdScene(adScene *model.AdScene) error {
	err := common.DB.Model(adScene).Updates(adScene).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (cr AdSceneRepository) BatchDeleteAdSceneByIds(ids []string) error {
	var adScenes []model.AdScene
	for _, id := range ids {
		// 根据ID获取标签
		adScene, err := cr.GetAdSceneByAdSceneID(id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%s的推广场景", id)
		}
		adScenes = append(adScenes, adScene)
	}

	err := common.DB.Unscoped().Delete(&adScenes).Error

	return err
}
