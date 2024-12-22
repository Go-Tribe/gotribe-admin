// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

// ProductSpecItemDto 定义了产品类型信息传输的数据结构。
// 包含产品类型ID、标题、备注、类别ID、规格ID和创建时间等基本信息。
type ProductSpecItemDto struct {
	ItemID    string `json:"productSpecItemID"`
	Title     string `json:"title"`
	SpecID    string `json:"specID"`
	Sort      uint   `json:"sort"`
	Enabled   uint   `json:"enabled"`
	CreatedAt string `json:"createdAt"`
}

// toProductSpecItemDto 将产品类型模型转换为产品类型DTO。
// 参数 productSpecItem: 产品类型模型的指针。
// 返回值: 返回一个产品类型DTO，如果参数为nil则返回空DTO。
func toProductSpecItemDto(productSpecItem *model.ProductSpecItem) ProductSpecItemDto {
	if productSpecItem == nil {
		return ProductSpecItemDto{}
	}

	createdAt := ""
	if !productSpecItem.CreatedAt.IsZero() {
		createdAt = productSpecItem.CreatedAt.Format(known.TimeFormat)
	}

	return ProductSpecItemDto{
		ItemID:    productSpecItem.ItemID,
		Title:     productSpecItem.Title,
		Sort:      productSpecItem.Sort,
		SpecID:    productSpecItem.SpecID,
		Enabled:   productSpecItem.Enabled,
		CreatedAt: createdAt,
	}
}

// ToProductSpecItemInfoDto 将产品类型模型转换为产品类型信息DTO。
// 参数 productSpecItem: 产品类型模型。
// 返回值: 返回一个产品类型DTO。
func ToProductSpecItemInfoDto(productSpecItem model.ProductSpecItem) ProductSpecItemDto {
	return toProductSpecItemDto(&productSpecItem)
}

// ToProductSpecItemsDto 将产品类型模型列表转换为产品类型DTO列表。
// 参数 productSpecItemList: 产品类型模型列表。
// 返回值: 返回一个产品类型DTO列表。
func ToProductSpecItemsDto(productSpecItemList []*model.ProductSpecItem) []ProductSpecItemDto {
	var productSpecItems []ProductSpecItemDto
	for _, productSpecItem := range productSpecItemList {
		productSpecItemDto := toProductSpecItemDto(productSpecItem)
		productSpecItems = append(productSpecItems, productSpecItemDto)
	}

	return productSpecItems
}
