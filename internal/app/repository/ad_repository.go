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

type IAdRepository interface {
	CreateAd(ctx context.Context, ad *model.Ad) error                              // 创建推广场景
	GetAdByID(ctx context.Context, id uint) (model.Ad, error)                      // 获取单个推广场景
	GetAds(ctx context.Context, req *vo.AdListRequest) ([]*model.Ad, int64, error) // 获取推广场景列表
	UpdateAd(ctx context.Context, ad *model.Ad) error                              // 更新推广场景
	BatchDeleteAdByIds(ctx context.Context, ids []uint) error                      // 批量删除
}

type AdRepository struct {
}

// AdRepository构造函数
func NewAdRepository() IAdRepository {
	return AdRepository{}
}

func buildAdOrder(req *vo.AdListRequest) string {
	sortByMap := map[string]string{
		"id":          "id",
		"title":       "title",
		"description": "description",
		"status":      "status",
		"createdAt":   "created_at",
		"created_at":  "created_at",
	}

	column, ok := sortByMap[strings.TrimSpace(req.SortBy)]
	if !ok {
		return "created_at DESC"
	}

	direction := "ASC"
	if strings.EqualFold(strings.TrimSpace(req.SortOrder), "desc") {
		direction = "DESC"
	}

	return fmt.Sprintf("%s %s", column, direction)
}

// 获取单个推广场景
func (cr AdRepository) GetAdByID(ctx context.Context, id uint) (model.Ad, error) {
	var ad model.Ad
	err := common.WithContext(ctx).DB().Where("id = ?", id).First(&ad).Error
	return ad, err
}

// 获取推广场景列表
func (cr AdRepository) GetAds(ctx context.Context, req *vo.AdListRequest) ([]*model.Ad, int64, error) {
	var list []*model.Ad
	db := common.WithContext(ctx).DB().Model(&model.Ad{}).Order(buildAdOrder(req))

	if req.SceneID > 0 {
		db = db.Where("scene_id = ?", req.SceneID)
	}
	if !gconvert.IsEmpty(req.Title) {
		db = db.Where("title like ?", fmt.Sprintf("%%%s%%", req.Title))
	}
	if !gconvert.IsEmpty(req.Status) {
		db = db.Where("status = ?", fmt.Sprintf("%d", req.Status))
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
	list, err = GetAdOther(ctx, list)
	return list, total, err
}

// 获取推广场景其他信息
func GetAdOther(ctx context.Context, ads []*model.Ad) ([]*model.Ad, error) {
	if len(ads) == 0 {
		return ads, nil
	}

	// 收集所有 SceneID
	sceneIDs := make([]uint, 0, len(ads))
	for _, m := range ads {
		if m.SceneID > 0 {
			sceneIDs = append(sceneIDs, m.SceneID)
		}
	}

	if len(sceneIDs) == 0 {
		return ads, nil
	}

	// 批量查询
	var adScenes []*model.AdScene
	if err := common.WithContext(ctx).DB().Where("id IN ?", sceneIDs).Find(&adScenes).Error; err != nil {
		return ads, err
	}

	// 建立映射
	adSceneMap := make(map[uint]*model.AdScene)
	for _, scene := range adScenes {
		adSceneMap[scene.ID] = scene
	}

	// 赋值
	for _, m := range ads {
		if scene, ok := adSceneMap[m.SceneID]; ok {
			m.Scene = scene
		}
	}
	return ads, nil
}

// 创建推广场景
func (cr AdRepository) CreateAd(ctx context.Context, ad *model.Ad) error {
	err := common.WithContext(ctx).DB().Create(ad).Error
	return err
}

// 更新推广场景
func (cr AdRepository) UpdateAd(ctx context.Context, ad *model.Ad) error {
	err := common.WithContext(ctx).DB().Model(ad).Updates(ad).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (cr AdRepository) BatchDeleteAdByIds(ctx context.Context, ids []uint) error {
	var ads []model.Ad
	for _, id := range ids {
		// 根据ID获取标签
		ad, err := cr.GetAdByID(ctx, id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%d的推广场景", id)
		}
		ads = append(ads, ad)
	}

	err := common.WithContext(ctx).DB().Unscoped().Delete(&ads).Error

	return err
}
