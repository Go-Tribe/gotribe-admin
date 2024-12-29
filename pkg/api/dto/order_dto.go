// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/util"
)

type OrderDto struct {
	OrderID           string  `json:"orderID"`
	OrderNumber       string  `json:"orderNumber"`
	OrderType         uint    `json:"orderType"`
	PayMethod         uint    `json:"payMethod"`
	PayStatus         uint    `json:"payStatus"`
	PayTime           string  `json:"payTime"`
	ProductID         string  `json:"productID"`
	ProductName       string  `json:"productName"`
	ProductSku        string  `json:"productSku"`
	ProjectID         string  `json:"projectID"`
	Quantity          uint    `json:"quantity"`
	RefundStatus      uint    `json:"refundStatus"`
	RefundTime        string  `json:"refundTime"`
	Remark            string  `json:"remark"`
	RemarkAdmin       string  `json:"remarkAdmin"`
	Status            uint    `json:"status"`
	UnitPrice         float64 `json:"unitPrice"`
	UserID            string  `json:"userID"`
	Username          string  `json:"username"`
	Amount            float64 `json:"amount"`
	AmountPay         float64 `json:"amountPay"`
	ConsigneeName     string  `json:"consigneeName"`
	ConsigneePhone    string  `json:"consigneePhone"`
	ConsigneeAddress  string  `json:"consigneeAddress"`
	ConsigneeProvince string  `json:"consigneeProvince"`
	ProductImage      string  `json:"productImage"`
	ConsigneeStreet   string  `json:"consigneeStreet"`
	ConsigneeDistrict string  `json:"consigneeDistrict"`
	ConsigneeCity     string  `json:"consigneeCity"`
	CreatedAt         string  `json:"createdAt"`
}

func ToOrderInfoDto(order *model.Order) OrderDto {
	if order == nil {
		return OrderDto{}
	}

	return OrderDto{
		OrderID:           order.OrderID,
		OrderNumber:       order.OrderNumber,
		OrderType:         order.OrderType,
		PayMethod:         order.PayMethod,
		PayStatus:         order.PayStatus,
		PayTime:           util.FormatTime(order.PayTime),
		ProductID:         order.ProductID,
		ProductName:       order.ProductName,
		ProductSku:        order.ProductSku,
		ProjectID:         order.ProjectID,
		Quantity:          order.Quantity,
		RefundStatus:      order.RefundStatus,
		RefundTime:        util.FormatTime(order.RefundTime),
		Remark:            order.Remark,
		RemarkAdmin:       order.RemarkAdmin,
		Status:            order.Status,
		UnitPrice:         util.FenToYuan(order.UnitPrice),
		UserID:            order.UserID,
		Username:          order.Username,
		Amount:            util.FenToYuan(order.Amount),
		AmountPay:         util.FenToYuan(order.AmountPay),
		ConsigneeName:     order.ConsigneeName,
		ConsigneePhone:    order.ConsigneePhone,
		ConsigneeAddress:  order.ConsigneeAddress,
		ConsigneeProvince: order.ConsigneeProvince,
		ProductImage:      order.ProductImage,
		ConsigneeStreet:   order.ConsigneeStreet,
		ConsigneeDistrict: order.ConsigneeDistrict,
		ConsigneeCity:     order.ConsigneeCity,
		CreatedAt:         util.FormatTime(order.CreatedAt),
	}
}

func ToOrdersDto(orderList []*model.Order) []OrderDto {
	if orderList == nil {
		return nil
	}

	orders := make([]OrderDto, 0, len(orderList))
	for _, order := range orderList {
		if order == nil {
			continue
		}
		orderDto := ToOrderInfoDto(order)
		orders = append(orders, orderDto)
	}

	return orders
}
