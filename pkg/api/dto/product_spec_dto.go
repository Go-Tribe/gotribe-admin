// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

// ProductSpecDto 定义了产品类型信息传输的数据结构。
// 包含产品类型ID、标题、备注、类别ID、规格ID和创建时间等基本信息。
type ProductSpecDto struct {
	ProductSpecID string               `json:"productSpecID"`
	Title         string               `json:"title"`
	Remark        string               `json:"remark"`
	Format        uint                 `json:"format"`
	Image         string               `json:"image"`
	Sort          uint                 `json:"sort"`
	CreatedAt     string               `json:"createdAt"`
	Items         []ProductSpecItemDto `json:"item"`
}

// toProductSpecDto 将产品类型模型转换为产品类型DTO。
// 参数 productSpec: 产品类型模型的指针。
// 返回值: 返回一个产品类型DTO，如果参数为nil则返回空DTO。
func toProductSpecDto(productSpec *model.ProductSpec) ProductSpecDto {
	if productSpec == nil {
		return ProductSpecDto{}
	}

	createdAt := ""
	if !productSpec.CreatedAt.IsZero() {
		createdAt = productSpec.CreatedAt.Format(known.TimeFormat)
	}
	items := make([]ProductSpecItemDto, 0, len(productSpec.Items))
	for _, item := range productSpec.Items {
		items = append(items, ToProductSpecItemInfoDto(item))
	}

	return ProductSpecDto{
		ProductSpecID: productSpec.ProductSpecID,
		Title:         productSpec.Title,
		Remark:        productSpec.Remark,
		Format:        productSpec.Format,
		Image:         productSpec.Image,
		Sort:          productSpec.Sort,
		CreatedAt:     createdAt,
		Items:         items,
	}
}

// ToProductSpecInfoDto 将产品类型模型转换为产品类型信息DTO。
// 参数 productSpec: 产品类型模型。
// 返回值: 返回一个产品类型DTO。
func ToProductSpecInfoDto(productSpec model.ProductSpec) ProductSpecDto {
	return toProductSpecDto(&productSpec)
}

// ToProductSpecsDto 将产品类型模型列表转换为产品类型DTO列表。
// 参数 productSpecList: 产品类型模型列表。
// 返回值: 返回一个产品类型DTO列表。
func ToProductSpecsDto(productSpecList []*model.ProductSpec) []ProductSpecDto {
	var productSpecs []ProductSpecDto
	for _, productSpec := range productSpecList {
		productSpecDto := toProductSpecDto(productSpec)
		productSpecs = append(productSpecs, productSpecDto)
	}

	return productSpecs
}
