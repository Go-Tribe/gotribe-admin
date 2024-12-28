// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 获取订单列表结构体
type OrderListRequest struct {
	OrderNumber string `form:"orderID" json:"orderID"`
	PageNum     uint   `json:"pageNum" form:"pageNum"`
	PageSize    uint   `json:"pageSize" form:"pageSize"`
}

// 批量删除订单结构体
type DeleteOrdersRequest struct {
	OrderIds string `json:"orderIds" form:"orderIds"`
}

// 创建订单结构体
type CreateOrderRequest struct {
	AmountPay   float64 `json:"amountPay" form:"amountPay" validate:"required,gte=1"`
	RemarkAdmin string  `json:"remarkAdmin" form:"remarkAdmin"`
	Status      uint    `json:"status" form:"status" validate:"required"`
}
