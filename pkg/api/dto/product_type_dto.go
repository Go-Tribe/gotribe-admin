// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

// ProductTypeDto 定义了产品类型信息传输的数据结构。
// 包含产品类型ID、标题、备注、类别ID、规格ID和创建时间等基本信息。
type ProductTypeDto struct {
	ProductTypeID string               `json:"productTypeID"`
	Title         string               `json:"title"`
	Remark        string               `json:"remark"`
	CategoryID    string               `json:"categoryID"`
	Spec          []*model.ProductSpec `json:"spec"`
	SpecIds       string               `json:"specIds"`
	CreatedAt     string               `json:"createdAt"`
}

// toProductTypeDto 将产品类型模型转换为产品类型DTO。
// 参数 productType: 产品类型模型的指针。
// 返回值: 返回一个产品类型DTO，如果参数为nil则返回空DTO。
func toProductTypeDto(productType *model.ProductType) ProductTypeDto {
	if productType == nil {
		return ProductTypeDto{}
	}

	createdAt := ""
	if !productType.CreatedAt.IsZero() {
		createdAt = productType.CreatedAt.Format(known.TimeFormat)
	}

	return ProductTypeDto{
		ProductTypeID: productType.ProductTypeID,
		Title:         productType.Title,
		Remark:        productType.Remark,
		CategoryID:    productType.ProductCategoryID,
		SpecIds:       productType.SpecIds,
		CreatedAt:     createdAt,
	}
}

// ToProductTypeInfoDto 将产品类型模型转换为产品类型信息DTO。
// 参数 productType: 产品类型模型。
// 返回值: 返回一个产品类型DTO。
func ToProductTypeInfoDto(productType model.ProductType) ProductTypeDto {
	return toProductTypeDto(&productType)
}

// ToProductTypesDto 将产品类型模型列表转换为产品类型DTO列表。
// 参数 productTypeList: 产品类型模型列表。
// 返回值: 返回一个产品类型DTO列表。
func ToProductTypesDto(productTypeList []*model.ProductType) []ProductTypeDto {
	var productTypes []ProductTypeDto
	for _, productType := range productTypeList {
		productTypeDto := toProductTypeDto(productType)
		productTypes = append(productTypes, productTypeDto)
	}

	return productTypes
}
