// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"errors"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"

	"gorm.io/gorm"
)

type IProductCategoryRepository interface {
	GetConfigByProductCategoryID(productCategoryID string) (model.ProductCategory, error)
	GetProductCategorys() ([]*model.ProductCategory, error)                                           // 获取分类列表
	GetProductCategoryTree() ([]*model.ProductCategory, error)                                        // 获取分类树
	CreateProductCategory(productCategory *model.ProductCategory) error                               // 创建分类
	UpdateProductCategoryByID(productCategoryID string, productCategory *model.ProductCategory) error // 更新分类
	BatchDeleteProductCategoryByIds(productCategoryIds []string) error                                // 批量删除分类
}

type ProductCategoryRepository struct {
}

func NewProductCategoryRepository() IProductCategoryRepository {
	return ProductCategoryRepository{}
}

// 获取单个分类详情
func (cr ProductCategoryRepository) GetConfigByProductCategoryID(productCategoryID string) (model.ProductCategory, error) {
	var productCategory model.ProductCategory
	err := common.DB.Where("product_category_id = ?", productCategoryID).First(&productCategory).Error
	return productCategory, err
}

// 获取分类列表
func (cr ProductCategoryRepository) GetProductCategorys() ([]*model.ProductCategory, error) {
	var productCategorys []*model.ProductCategory
	err := common.DB.Order("sort").Find(&productCategorys).Error
	return productCategorys, err
}

// 获取分类树
func (cr ProductCategoryRepository) GetProductCategoryTree() ([]*model.ProductCategory, error) {
	var productCategorys []*model.ProductCategory
	err := common.DB.Order("sort").Find(&productCategorys).Error
	return GenProductCategoryTree(0, productCategorys), err
}

func GenProductCategoryTree(parentID uint, productCategorys []*model.ProductCategory) []*model.ProductCategory {
	tree := make([]*model.ProductCategory, 0)

	for _, m := range productCategorys {
		if *m.ParentID == parentID {
			children := GenProductCategoryTree(m.ID, productCategorys)
			m.Children = children
			tree = append(tree, m)
		}
	}
	return tree
}

// 创建分类
func (cr ProductCategoryRepository) CreateProductCategory(productCategory *model.ProductCategory) error {
	err := common.DB.Create(productCategory).Error
	return err
}

// 更新分类
func (cr ProductCategoryRepository) UpdateProductCategoryByID(productCategoryID string, productCategory *model.ProductCategory) error {
	err := common.DB.Model(productCategory).Where("productCategory_id = ?", productCategoryID).Updates(productCategory).Error
	return err
}

// 批量删除分类
func (cr ProductCategoryRepository) BatchDeleteProductCategoryByIds(productCategoryIds []string) error {
	var productCategorys []*model.ProductCategory

	err := common.DB.Where("product_category_id IN (?)", productCategoryIds).Find(&productCategorys).Error
	if err != nil {
		return err
	}
	j := 0
	for _, productCategory := range productCategorys {
		if productCategory.ID != known.DefulatID && !isProductCategoryPID(int64((productCategory.ID))) {
			productCategorys[j] = productCategory
			j++
		}
	}
	// Slice productCategorys to new size.
	productCategorys = productCategorys[:j]
	err = common.DB.Unscoped().Delete(&productCategorys).Error
	return err
}

// isPID 判断是否为别人的父类 ID
// 存在 true 不存在 false
func isProductCategoryPID(ID int64) bool {
	var productCategory model.ProductCategory
	if err := common.DB.Where("parent_id = ?", ID).First(&productCategory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		} else {
			common.Log.Error(err.Error())
			return false
		}
	}
	return true
}
