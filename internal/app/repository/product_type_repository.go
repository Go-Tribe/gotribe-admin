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

type IProductTypeRepository interface {
	CreateProductType(productType *model.ProductType) (*model.ProductType, error)        // 创建商品类型
	GetProductTypeByProductTypeID(productTypeID string) (model.ProductType, error)       // 获取单个商品类型
	GetProductTypes(req *vo.ProductTypeListRequest) ([]*model.ProductType, int64, error) // 获取商品类型列表
	UpdateProductType(productType *model.ProductType) error                              // 更新商品类型
	BatchDeleteProductTypeByIds(ids []string) error                                      // 批量删除
}

type ProductTypeRepository struct {
}

// ProductTypeRepository构造函数
func NewProductTypeRepository() IProductTypeRepository {
	return ProductTypeRepository{}
}

// 获取单个商品类型
func (tr ProductTypeRepository) GetProductTypeByProductTypeID(productTypeID string) (model.ProductType, error) {
	var productType model.ProductType
	err := common.DB.Where("product_type_id = ?", productTypeID).First(&productType).Error
	return productType, err
}

// 获取商品类型列表
func (tr ProductTypeRepository) GetProductTypes(req *vo.ProductTypeListRequest) ([]*model.ProductType, int64, error) {
	var list []*model.ProductType
	db := common.DB.Model(&model.ProductType{}).Order("created_at DESC")

	title := strings.TrimSpace(req.Title)
	if title != "" {
		db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title))
	}
	productTypeID := strings.TrimSpace(req.ProductTypeID)
	if req.ProductTypeID != "" {
		db = db.Where("productType_id = ?", fmt.Sprintf("%s", productTypeID))
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

// 创建商品类型
func (tr ProductTypeRepository) CreateProductType(productType *model.ProductType) (*model.ProductType, error) {
	if isProductTypeExist(productType.Title) {
		return nil, fmt.Errorf("%s商品类型已存在", productType.Title)
	}
	result := common.DB.Create(productType)
	if result.Error != nil {
		return nil, result.Error
	}
	return productType, nil
}

// 更新商品类型
func (tr ProductTypeRepository) UpdateProductType(productType *model.ProductType) error {
	err := common.DB.Model(productType).Updates(productType).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (tr ProductTypeRepository) BatchDeleteProductTypeByIds(ids []string) error {
	var productTypes []model.ProductType
	for _, id := range ids {
		// 根据ID获取商品类型
		productType, err := tr.GetProductTypeByProductTypeID(id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%s的商品类型", id)
		}
		productTypes = append(productTypes, productType)
	}

	err := common.DB.Unscoped().Delete(&productTypes).Error

	return err
}

func isProductTypeExist(title string) bool {
	var productType model.ProductType
	result := common.DB.Where("title = ?", title).First(&productType)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
