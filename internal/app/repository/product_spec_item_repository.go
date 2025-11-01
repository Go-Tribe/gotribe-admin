// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"errors"
	"fmt"
	"github.com/dengmengmian/ghelper/gconvert"
	"gorm.io/gorm"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type IProductSpecItemRepository interface {
	CreateProductSpecItem(productSpecItem *model.ProductSpecItem) (*model.ProductSpecItem, error)    // 创建商品规格值
	GetProductSpecItemByItemID(productSpecItemID string) (model.ProductSpecItem, error)              // 获取单个商品规格值
	GetProductSpecItems(req *vo.ProductSpecItemListRequest) ([]*model.ProductSpecItem, int64, error) // 获取商品规格值列表
	UpdateProductSpecItem(productSpecItem *model.ProductSpecItem) error                              // 更新商品规格值
	BatchDeleteProductSpecItemByIds(ids []string) error                                              // 批量删除
}

type ProductSpecItemRepository struct {
}

// ProductSpecItemRepository构造函数
func NewProductSpecItemRepository() IProductSpecItemRepository {
	return ProductSpecItemRepository{}
}

// 获取单个商品规格值
func (tr ProductSpecItemRepository) GetProductSpecItemByItemID(productSpecItemID string) (model.ProductSpecItem, error) {
	var productSpecItem model.ProductSpecItem
	err := common.DB.Where("item_id = ?", productSpecItemID).First(&productSpecItem).Error
	return productSpecItem, err
}

// 获取商品规格值列表
func (tr ProductSpecItemRepository) GetProductSpecItems(req *vo.ProductSpecItemListRequest) ([]*model.ProductSpecItem, int64, error) {
	var list []*model.ProductSpecItem
	db := common.DB.Model(&model.ProductSpecItem{}).Order("created_at DESC")

	specID := strings.TrimSpace(req.SpecID)
	if !gconvert.IsEmpty(specID) {
		db = db.Where("spec_id = ?", specID)
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
	return list, total, err
}

// 创建商品规格值
func (tr ProductSpecItemRepository) CreateProductSpecItem(productSpecItem *model.ProductSpecItem) (*model.ProductSpecItem, error) {
	if isProductSpecItemExist(productSpecItem.Title) {
		return nil, fmt.Errorf("%s商品规格值已存在", productSpecItem.Title)
	}
	result := common.DB.Create(productSpecItem)
	if result.Error != nil {
		return nil, result.Error
	}
	return productSpecItem, nil
}

// 更新商品规格值
func (tr ProductSpecItemRepository) UpdateProductSpecItem(productSpecItem *model.ProductSpecItem) error {
	err := common.DB.Model(productSpecItem).Updates(productSpecItem).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (tr ProductSpecItemRepository) BatchDeleteProductSpecItemByIds(ids []string) error {
	var productSpecItems []model.ProductSpecItem
	for _, id := range ids {
		// 根据ID获取商品规格值
		productSpecItem, err := tr.GetProductSpecItemByItemID(id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%s的商品规格值", id)
		}
		productSpecItems = append(productSpecItems, productSpecItem)
	}

	err := common.DB.Unscoped().Delete(&productSpecItems).Error

	return err
}

func isProductSpecItemExist(title string) bool {
	var productSpecItem model.ProductSpecItem
	result := common.DB.Where("title = ?", title).First(&productSpecItem)
	return !errors.Is(result.Error, gorm.ErrRecordNotFound)
}
