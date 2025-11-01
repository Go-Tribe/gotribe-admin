// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type IProductSpecRepository interface {
	CreateProductSpec(productSpec *model.ProductSpec) (*model.ProductSpec, error)          // 创建商品规格
	GetProductSpecByProductSpecID(productSpecID string) (model.ProductSpec, error)         // 获取单个商品规格
	GetProductSpecs(req *vo.ProductSpecListRequest) ([]*model.ProductSpec, int64, error)   // 获取商品规格列表
	UpdateProductSpec(productSpec *model.ProductSpec) error                                // 更新商品规格
	BatchDeleteProductSpecByIds(ids []string) error                                        // 批量删除
	GetProductSpecsByProductSpecIDs(productSpecIDs []string) ([]*model.ProductSpec, error) // 获取多个商品规格
	GetProductSpecAndItem(categoryID string) ([]*model.ProductSpec, error)                 // 通过商品分类获取商品规格和规格项
}

type ProductSpecRepository struct {
}

// ProductSpecRepository构造函数
func NewProductSpecRepository() IProductSpecRepository {
	return ProductSpecRepository{}
}

// 获取单个商品规格
func (tr ProductSpecRepository) GetProductSpecByProductSpecID(productSpecID string) (model.ProductSpec, error) {
	var productSpec model.ProductSpec
	err := common.DB.Where("product_spec_id = ?", productSpecID).First(&productSpec).Error
	return productSpec, err
}

// 获取商品规格列表
func (tr ProductSpecRepository) GetProductSpecs(req *vo.ProductSpecListRequest) ([]*model.ProductSpec, int64, error) {
	var list []*model.ProductSpec
	db := common.DB.Model(&model.ProductSpec{}).Order("created_at DESC")

	title := strings.TrimSpace(req.Title)
	if title != "" {
		db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title))
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

// 创建商品规格
func (tr ProductSpecRepository) CreateProductSpec(productSpec *model.ProductSpec) (*model.ProductSpec, error) {
	if isProductSpecExist(productSpec.Title) {
		return nil, fmt.Errorf("%s商品规格已存在", productSpec.Title)
	}
	result := common.DB.Create(productSpec)
	if result.Error != nil {
		return nil, result.Error
	}
	return productSpec, nil
}

// 更新商品规格
func (tr ProductSpecRepository) UpdateProductSpec(productSpec *model.ProductSpec) error {
	err := common.DB.Model(productSpec).Updates(productSpec).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (tr ProductSpecRepository) BatchDeleteProductSpecByIds(ids []string) error {
	var productSpecs []model.ProductSpec
	for _, id := range ids {
		// 根据ID获取商品规格
		productSpec, err := tr.GetProductSpecByProductSpecID(id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%s的商品规格", id)
		}
		productSpecs = append(productSpecs, productSpec)
	}

	err := common.DB.Unscoped().Delete(&productSpecs).Error

	return err
}

func isProductSpecExist(title string) bool {
	var productSpec model.ProductSpec
	result := common.DB.Where("title = ?", title).First(&productSpec)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

// 获取多个商品规格
func (tr ProductSpecRepository) GetProductSpecsByProductSpecIDs(productSpecIDs []string) ([]*model.ProductSpec, error) {
	var productSpecs []*model.ProductSpec
	err := common.DB.Where("product_spec_id IN (?)", productSpecIDs).Find(&productSpecs).Error
	return productSpecs, err
}

func (tr ProductSpecRepository) GetProductSpecAndItem(categoryID string) ([]*model.ProductSpec, error) {
	// 通过分类获取商品类型里的spec_ids
	var productType model.ProductType
	err := common.DB.Where("product_category_id = ?", categoryID).First(&productType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 或者根据业务需求返回适当的错误信息
		}
		return nil, err
	}

	// 拿到的spec_ids去查所有规格信息
	specIDs := strings.Split(productType.SpecIds, ",")
	if len(specIDs) == 0 {
		return nil, nil // 或者根据业务需求返回适当的错误信息
	}

	var productSpecList []*model.ProductSpec
	err = common.DB.Where("product_spec_id in (?)", specIDs).Find(&productSpecList).Error
	if err != nil {
		return nil, err
	}

	// 获取所有规格项信息
	specIDsForItems := make([]string, 0, len(productSpecList))
	for _, productSpec := range productSpecList {
		specIDsForItems = append(specIDsForItems, productSpec.ProductSpecID)
	}

	var productSpecItemList []model.ProductSpecItem
	if len(specIDsForItems) > 0 {
		err = common.DB.Where("spec_id IN (?)", specIDsForItems).Find(&productSpecItemList).Error
		if err != nil {
			return nil, err
		}
	}

	// 组合返回规格和规格项信息
	for i := range productSpecList {
		productSpecList[i].Items = filterItemsBySpecID(productSpecItemList, productSpecList[i].ProductSpecID)
	}

	return productSpecList, nil
}

func filterItemsBySpecID(items []model.ProductSpecItem, specID string) []*model.ProductSpecItem {
	result := make([]*model.ProductSpecItem, 0)
	for i := range items {
		if items[i].SpecID == specID {
			result = append(result, &items[i])
		}
	}
	return result
}
