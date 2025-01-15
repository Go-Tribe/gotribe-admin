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

type IOrderLogRepository interface {
	GetOrderLogs(orderID string) ([]*model.OrderLog, int64, error) // 获取订单记录列表
	CreateOrderLog(orderID, remark string) error
}

type OrderLogRepository struct {
}

// OrderLogRepository构造函数
func NewOrderLogRepository() IOrderLogRepository {
	return OrderLogRepository{}
}

// 获取单个订单记录
func (tr OrderLogRepository) GetOrderLogByOrderLogID(orderLogID string) (model.OrderLog, error) {
	var orderLog model.OrderLog
	err := common.DB.Where("orderLog_id = ?", orderLogID).First(&orderLog).Error
	return orderLog, err
}

// 获取订单记录列表
func (tr OrderLogRepository) GetOrderLogs(orderID string) ([]*model.OrderLog, int64, error) {
	var list []*model.OrderLog
	db := common.DB.Model(&model.OrderLog{}).Order("created_at DESC")

	if orderID != "" {
		db = db.Where("order_id = ?", fmt.Sprintf("%s", orderID))
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	err = db.Find(&list).Error
	return list, total, err
}

func (tr OrderLogRepository) CreateOrderLog(orderID, remark string) error {
	return common.DB.Create(&model.OrderLog{
		OrderID: orderID,
		Remark:  remark,
	}).Error
}
