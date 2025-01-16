// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"fmt"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/pkg/util"
	"time"
)

type IIndexRepository interface {
}

type IndexRepository struct {
}

// IndexRepository构造函数
func NewIndexRepository() IIndexRepository {
	return IndexRepository{}
}

// 获取当日 销售额，订单量，新增用户，访问量
func (r IndexRepository) GetIndexData(projectID string) (map[string]interface{}, error) {
	// 动态生成当天的时间范围
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	// 获取订单销售额和订单数量
	var result struct {
		TotalSales  int64
		TotalOrders int64
	}
	err := common.DB.Table("orders").
		Select("SUM(amount_pay) as total_sales, COUNT(*) as total_orders").
		Where("created_at >= ? AND pay_status = 2 AND project_id = ?", startOfDay, projectID).
		Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order data: %v", err)
	}

	// 获取新增用户数
	var totalUsers int64
	err = common.DB.Table("users").
		Select("COUNT(*)").
		Where("created_at >= ? AND project_id = ?", startOfDay, projectID).
		Count(&totalUsers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %v", err)
	}

	// 返回结果
	data := map[string]interface{}{
		"sales":      util.FenToYuan(int(result.TotalSales)),
		"orders":     result.TotalOrders,
		"newUsers":   totalUsers,
		"visitCount": 0,
	}

	return data, nil
}
