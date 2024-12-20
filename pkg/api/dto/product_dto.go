// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

// ProductDto 定义了产品类型信息传输的数据结构。
// 包含产品类型ID、标题、备注、类别ID、规格ID和创建时间等基本信息。
type ProductDto struct {
	ProductID  string               `json:"productID"`
	Title      string               `json:"title"`
	Remark     string               `json:"remark"`
	CategoryID string               `json:"categoryID"`
	Spec       []*model.ProductSpec `json:"spec"`
	SpecIds    string               `json:"specIds"`
	CreatedAt  string               `json:"createdAt"`
}

// toProductDto 将产品类型模型转换为产品类型DTO。
// 参数 product: 产品类型模型的指针。
// 返回值: 返回一个产品类型DTO，如果参数为nil则返回空DTO。
func toProductDto(product *model.Product) ProductDto {
	if product == nil {
		return ProductDto{}
	}

	createdAt := ""
	if !product.CreatedAt.IsZero() {
		createdAt = product.CreatedAt.Format(known.TimeFormat)
	}

	return ProductDto{
		ProductID: product.ProductID,
		Title:     product.Title,
		CreatedAt: createdAt,
	}
}

// ToProductInfoDto 将产品类型模型转换为产品类型信息DTO。
// 参数 product: 产品类型模型。
// 返回值: 返回一个产品类型DTO。
func ToProductInfoDto(product model.Product) ProductDto {
	return toProductDto(&product)
}

// ToProductsDto 将产品类型模型列表转换为产品类型DTO列表。
// 参数 productList: 产品类型模型列表。
// 返回值: 返回一个产品类型DTO列表。
func ToProductsDto(productList []*model.Product) []ProductDto {
	var products []ProductDto
	for _, product := range productList {
		productDto := toProductDto(product)
		products = append(products, productDto)
	}

	return products
}
