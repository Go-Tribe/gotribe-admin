// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"fmt"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
)

type IProductSkuRepository interface {
	CreateProductSku(productSku *model.ProductSku) (*model.ProductSku, error)  // 创建sku
	GetProductSkuByProductSkuID(productSkuID string) (model.ProductSku, error) // 获取单个sku
	UpdateProductSku(productSku *model.ProductSku) error                       // 更新sku
	BatchDeleteProductSkuByIds(ids []string) error                             // 批量删除
	GetProductSkuByProductID(productID string) ([]*model.ProductSku, error)
}

type ProductSkuRepository struct {
}

// ProductSkuRepository构造函数
func NewProductSkuRepository() IProductSkuRepository {
	return ProductSkuRepository{}
}

// 获取单个sku
func (tr ProductSkuRepository) GetProductSkuByProductSkuID(productSkuID string) (model.ProductSku, error) {
	var productSku model.ProductSku
	err := common.DB.Where("sku_id = ?", productSkuID).First(&productSku).Error
	return productSku, err
}

// 通过商品ID获取sku
func (tr ProductSkuRepository) GetProductSkuByProductID(productID string) ([]*model.ProductSku, error) {
	var productSkus []*model.ProductSku
	err := common.DB.Where("product_id = ?", productID).Find(&productSkus).Error
	return productSkus, err
}

// 创建sku
func (tr ProductSkuRepository) CreateProductSku(productSku *model.ProductSku) (*model.ProductSku, error) {
	result := common.DB.Create(productSku)
	if result.Error != nil {
		return nil, result.Error
	}
	return productSku, nil
}

// 更新sku
func (tr ProductSkuRepository) UpdateProductSku(productSku *model.ProductSku) error {
	err := common.DB.Model(productSku).Updates(productSku).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (tr ProductSkuRepository) BatchDeleteProductSkuByIds(ids []string) error {
	var productSkus []model.ProductSku
	for _, id := range ids {
		// 根据ID获取sku
		productSku, err := tr.GetProductSkuByProductSkuID(id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%s的sku", id)
		}
		productSkus = append(productSkus, productSku)
	}

	err := common.DB.Unscoped().Delete(&productSkus).Error

	return err
}
