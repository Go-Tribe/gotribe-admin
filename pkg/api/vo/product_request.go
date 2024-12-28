// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建产品结构体
type CreateProductRequest struct {
	Title         string   `form:"title" json:"title" validate:"required,min=2,max=60"`
	CategoryID    string   `form:"categoryID" json:"categoryID" validate:"required"`
	ProductNumber string   `form:"productNumber" json:"productNumber"`
	ProjectID     string   `form:"projectID" json:"projectID"  validate:"required"`
	Description   string   `form:"description" json:"description"`
	Image         []string `form:"images" json:"images"`
	Video         string   `form:"video" json:"video"`
	BuyLimit      uint     `form:"buyLimit" json:"buyLimit" validate:"required,min=1,max=100"`
	ProductSpec   string   `form:"productSpec" json:"productSpec"`
	Content       string   `form:"content" json:"content"`
	HtmlContent   string   `form:"Htmlcontent" json:"Htmlcontent"`
	Enable        uint     `form:"enable" json:"enable" validate:"oneof=1 2"`
	SKU           []Sku    `form:"sku" json:"sku" validate:"required"`
}

// 获取产品列表结构体
type ProductListRequest struct {
	CategoryID string `form:"categoryID" json:"categoryID"`
	ProjectID  string `form:"projectID" json:"projectID"`
	Title      string `form:"title" json:"title"`
	PageNum    uint   `json:"pageNum" form:"pageNum"`
	PageSize   uint   `json:"pageSize" form:"pageSize"`
}

// 批量删除产品结构体
type DeleteProductsRequest struct {
	ProductIds string `json:"productIds" form:"productIds"`
}

type Sku struct {
	SKUID         string  `form:"skuID" json:"skuID"`
	Title         string  `form:"skuTitle" json:"skuTitle" validate:"required"`
	Image         string  `form:"image" json:"image" validate:"required"`
	CostPrice     float64 `json:"costPrice" validate:"required"`
	MarketPrice   float64 `json:"marketPrice" validate:"required"`
	UnitPrice     float64 `json:"unitPrice" validate:"required"`
	UnitPoint     float64 `json:"unitPoint" validate:"required"`
	Quantity      uint    `json:"quantity" validate:"required"`
	EnableDefault uint    `json:"enableDefault"`
}
