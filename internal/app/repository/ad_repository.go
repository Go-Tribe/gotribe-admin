// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"errors"
	"fmt"
	"github.com/dengmengmian/ghelper/gconvert"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type IAdRepository interface {
	CreateAd(ad *model.Ad) error                              // 创建推广场景
	GetAdByAdID(adID string) (model.Ad, error)                // 获取单个推广场景
	GetAds(req *vo.AdListRequest) ([]*model.Ad, int64, error) // 获取推广场景列表
	UpdateAd(ad *model.Ad) error                              // 更新推广场景
	BatchDeleteAdByIds(ids []string) error                    // 批量删除
}

type AdRepository struct {
}

// AdRepository构造函数
func NewAdRepository() IAdRepository {
	return AdRepository{}
}

// 获取单个推广场景
func (cr AdRepository) GetAdByAdID(adID string) (model.Ad, error) {
	var ad model.Ad
	err := common.DB.Where("ad_id = ?", adID).First(&ad).Error
	return ad, err
}

// 获取推广场景列表
func (cr AdRepository) GetAds(req *vo.AdListRequest) ([]*model.Ad, int64, error) {
	var list []*model.Ad
	db := common.DB.Model(&model.Ad{}).Order("created_at DESC")

	adSceneID := strings.TrimSpace(req.SceneID)
	if !gconvert.IsEmpty(adSceneID) {
		db = db.Where("scene_id = ?", fmt.Sprintf("%s", adSceneID))
	}
	if !gconvert.IsEmpty(req.Title) {
		db = db.Where("title like ?", fmt.Sprintf("%%%s%%", req.Title))
	}
	if !gconvert.IsEmpty(req.Status) {
		db = db.Where("staus = ?", req.Status)
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
	return GetAdOther(list), total, err
}

// 获取推广场景其他信息
func GetAdOther(ads []*model.Ad) []*model.Ad {
	for _, m := range ads {
		var adScene *model.AdScene
		_ = common.DB.Where("ad_scene_id = ?", m.SceneID).First(&adScene).Error
		m.Scene = adScene
	}
	return ads
}

// 创建推广场景
func (cr AdRepository) CreateAd(ad *model.Ad) error {
	err := common.DB.Create(ad).Error
	return err
}

// 更新推广场景
func (cr AdRepository) UpdateAd(ad *model.Ad) error {
	err := common.DB.Model(ad).Updates(ad).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (cr AdRepository) BatchDeleteAdByIds(ids []string) error {
	var ads []model.Ad
	for _, id := range ids {
		// 根据ID获取标签
		ad, err := cr.GetAdByAdID(id)
		if err != nil {
			return errors.New(fmt.Sprintf("未获取到ID为%s的推广场景", id))
		}
		ads = append(ads, ad)
	}

	err := common.DB.Unscoped().Delete(&ads).Error

	return err
}
