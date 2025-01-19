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
	GetIndexData(projectID string) (map[string]interface{}, error)                                    // 获取首页数据
	GetTimeRangeData(timeRange string, projectID string) (map[string][]map[string]interface{}, error) // 获取时间范围数据
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
	err := common.DB.Table("order").
		Select("SUM(amount_pay) as total_sales, COUNT(*) as total_orders").
		Where("created_at >= ? AND pay_status = 2 AND project_id = ?", startOfDay, projectID).
		Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order data: %v", err)
	}

	// 获取新增用户数
	var totalUsers int64
	err = common.DB.Table("user").
		Select("COUNT(*)").
		Where("created_at >= ? AND project_id = ?", startOfDay, projectID).
		Count(&totalUsers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %v", err)
	}
	// 获取浏览数据
	var visitCount int64
	err = common.DB.Table("user_event").
		Select("COUNT(*)").
		Where("event_type = 1 AND created_at >= ? AND project_id = ?", startOfDay, projectID).
		Count(&visitCount).Error
	// 返回结果
	data := map[string]interface{}{
		"sales":      util.FenToYuan(int(result.TotalSales)),
		"orders":     result.TotalOrders,
		"newUsers":   totalUsers,
		"visitCount": 0,
	}

	return data, nil
}

func (r IndexRepository) GetTimeRangeData(timeRange string, projectID string) (map[string][]map[string]interface{}, error) {
	var startDate time.Time
	today := time.Now()

	switch timeRange {
	case "week":
		startDate = today.AddDate(0, 0, -7)
	case "month":
		startDate = today.AddDate(0, -1, 0)
	case "year":
		startDate = today.AddDate(-1, 0, 0)
	default:
		return nil, fmt.Errorf("invalid time range: %s", timeRange)
	}

	// 获取订单统计数据
	var orderResults []struct {
		Date        time.Time `gorm:"column:date"`
		TotalSales  int64     `gorm:"column:total_sales"`
		TotalOrders int64     `gorm:"column:total_orders"`
	}
	err := common.DB.Table("orders").
		Select("DATE(created_at) as date, SUM(amount_pay) as total_sales, COUNT(*) as total_orders").
		Where("created_at >= ? AND pay_status = 2 AND project_id = ?", startDate, projectID).
		Group("DATE(created_at)").
		Scan(&orderResults).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order data: %v", err)
	}

	// 获取用户统计数据
	var userResults []struct {
		Date       time.Time `gorm:"column:date"`
		TotalUsers int64     `gorm:"column:total_users"`
	}
	err = common.DB.Table("users").
		Select("DATE(created_at) as date, COUNT(*) as total_users").
		Where("created_at >= ? AND project_id = ?", startDate, projectID).
		Group("DATE(created_at)").
		Scan(&userResults).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %v", err)
	}

	// 构建结果
	orderData := make([]map[string]interface{}, len(orderResults))
	for i, res := range orderResults {
		orderData[i] = map[string]interface{}{
			"date":        res.Date.Format("2006-01-02"),
			"totalSales":  util.FenToYuan(int(res.TotalSales)),
			"totalOrders": res.TotalOrders,
		}
	}

	userData := make([]map[string]interface{}, len(userResults))
	for i, res := range userResults {
		userData[i] = map[string]interface{}{
			"date":       res.Date.Format("2006-01-02"),
			"totalUsers": res.TotalUsers,
		}
	}

	return map[string][]map[string]interface{}{
		"orders": orderData,
		"users":  userData,
	}, nil
}
