// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"errors"
	"fmt"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type IProductRepository interface {
	CreateProduct(product *model.Product) (*model.Product, error)            // 创建产品
	GetProductByProductID(productID string) (model.Product, error)           // 获取单个产品
	GetProducts(req *vo.ProductListRequest) ([]*model.Product, int64, error) // 获取产品列表
	UpdateProduct(product *model.Product) error                              // 更新产品
	BatchDeleteProductByIds(ids []string) error                              // 批量删除
}

type ProductRepository struct {
}

// ProductRepository构造函数
func NewProductRepository() IProductRepository {
	return ProductRepository{}
}

// 获取单个产品
func (tr ProductRepository) GetProductByProductID(productID string) (model.Product, error) {
	var product model.Product
	err := common.DB.Where("product_type_id = ?", productID).First(&product).Error
	return product, err
}

// 获取产品列表
func (tr ProductRepository) GetProducts(req *vo.ProductListRequest) ([]*model.Product, int64, error) {
	var list []*model.Product
	db := common.DB.Model(&model.Product{}).Order("created_at DESC")

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

// 创建产品
func (tr ProductRepository) CreateProduct(product *model.Product) (*model.Product, error) {
	result := common.DB.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

// 更新产品
func (tr ProductRepository) UpdateProduct(product *model.Product) error {
	err := common.DB.Model(product).Updates(product).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (tr ProductRepository) BatchDeleteProductByIds(ids []string) error {
	var products []model.Product
	for _, id := range ids {
		// 根据ID获取产品
		product, err := tr.GetProductByProductID(id)
		if err != nil {
			return errors.New(fmt.Sprintf("未获取到ID为%s的产品", id))
		}
		products = append(products, product)
	}

	err := common.DB.Unscoped().Delete(&products).Error

	return err
}
