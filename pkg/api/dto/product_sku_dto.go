// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

// ProductSkuDto 定义了产品类型信息传输的数据结构。
// 包含产品类型ID、标题、备注、类别ID、规格ID和创建时间等基本信息。
type ProductSkuDto struct {
	SKUID         string  `json:"skuID"`
	Title         string  `form:"skuTitle" json:"skuTitle" validate:"required"`
	Image         string  `form:"image" json:"image" validate:"required"`
	CostPrice     float64 `json:"cost_price" validate:"required"`
	MarketPrice   float64 `json:"market_price" validate:"required"`
	UnitPrice     float64 `json:"unit_price" validate:"required"`
	UnitPoint     float64 `json:"unit_point" validate:"required"`
	Quantity      uint    `json:"quantity" validate:"required"`
	EnableDefault uint    `json:"enable_default"`
	CreatedAt     string  `json:"createdAt"`
}

// toProductSkuDto 将产品类型模型转换为产品类型DTO。
// 参数 productSku: 产品类型模型的指针。
// 返回值: 返回一个产品类型DTO，如果参数为nil则返回空DTO。
func toProductSkuDto(productSku *model.ProductSku) ProductSkuDto {
	if productSku == nil {
		return ProductSkuDto{}
	}

	createdAt := ""
	if !productSku.CreatedAt.IsZero() {
		createdAt = productSku.CreatedAt.Format(known.TimeFormat)
	}

	return ProductSkuDto{
		SKUID:         productSku.SKUID,
		Title:         productSku.Title,
		Image:         productSku.Image,
		CostPrice:     float64(productSku.CostPrice / 100),
		MarketPrice:   float64(productSku.MarketPrice / 100),
		UnitPrice:     float64(productSku.UnitPrice / 100),
		UnitPoint:     float64(productSku.UnitPoint / 100),
		Quantity:      productSku.Quantity,
		EnableDefault: productSku.EnableDefault,
		CreatedAt:     createdAt,
	}
}

// ToProductSkuInfoDto 将产品类型模型转换为产品类型信息DTO。
// 参数 productSku: 产品类型模型。
// 返回值: 返回一个产品类型DTO。
func ToProductSkuInfoDto(productSku model.ProductSku) ProductSkuDto {
	return toProductSkuDto(&productSku)
}

// ToProductSkusDto 将产品类型模型列表转换为产品类型DTO列表。
// 参数 productSkuList: 产品类型模型列表。
// 返回值: 返回一个产品类型DTO列表。
func ToProductSkusDto(productSkuList []*model.ProductSku) []ProductSkuDto {
	var productSkus []ProductSkuDto
	for _, productSku := range productSkuList {
		productSkuDto := toProductSkuDto(productSku)
		productSkus = append(productSkus, productSkuDto)
	}

	return productSkus
}
