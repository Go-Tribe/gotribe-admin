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
	GetIndexData(projectID string) (map[string]interface{}, error)                             // 获取首页数据
	GetTimeRangeData(projectID, timeRange string) (map[string][]map[string]interface{}, error) // 获取时间范围数据
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
		Where("created_at >= ? AND status = 2 AND project_id = ?", startOfDay, projectID).
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
		"visitCount": visitCount,
	}

	return data, nil
}

// GetTimeRangeData retrieves the time range data for the index page based on the specified time range and project ID.
// It calculates the start date based on the time range ("week", "month", "year") and fetches order and user statistics from the database.
// The function returns a map containing the order and user data for each relevant date.
//
// Parameters:
// - projectID: The ID of the project for which to fetch the data.
// - timeRange: The time range for which to fetch the data. Valid values are "week", "month", and "year".
//
// Returns:
// - A map with keys "orders" and "users", each containing a slice of maps with date and corresponding statistics.
// - An error if the time range is invalid or if there is a failure in fetching data from the database.
func (r IndexRepository) GetTimeRangeData(projectID, timeRange string) (map[string][]map[string]interface{}, error) {
	var startDate time.Time
	today := time.Now()

	// Determine the start date based on the specified time range
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

	// Fetch order statistics from the database
	var orderResults []struct {
		Date        string `gorm:"column:date"`
		TotalSales  int64  `gorm:"column:total_sales"`
		TotalOrders int64  `gorm:"column:total_orders"`
	}
	var groupByField string
	if timeRange == "year" {
		groupByField = "DATE_FORMAT(created_at, '%Y-%m')"
	} else {
		groupByField = "DATE(created_at)"
	}
	err := common.DB.Table("order").
		Select(fmt.Sprintf("%s as date, SUM(amount_pay) as total_sales, COUNT(*) as total_orders", groupByField)).
		Where("created_at >= ? AND status = 2 AND project_id = ?", startDate, projectID).
		Group(groupByField).
		Scan(&orderResults).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order data: %v", err)
	}

	// Fetch user statistics from the database
	var userResults []struct {
		Date       string `gorm:"column:date"`
		TotalUsers int64  `gorm:"column:total_users"`
	}
	err = common.DB.Table("user").
		Select(fmt.Sprintf("%s as date, COUNT(*) as total_users", groupByField)).
		Where("created_at >= ? AND project_id = ?", startDate, projectID).
		Group(groupByField).
		Scan(&userResults).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %v", err)
	}

	// Construct the result map for order data
	orderData := make([]map[string]interface{}, len(orderResults))
	for i, res := range orderResults {
		dateStr := res.Date
		if timeRange == "year" {
			// For year, the date is already in the format "2006-01"
		} else {
			// For week and month, parse the date and format it to "2006-01-02"
			parsedDate, err := time.Parse("2006-01-02T15:04:05Z07:00", res.Date)
			if err != nil {
				return nil, fmt.Errorf("failed to parse date: %v", err)
			}
			dateStr = parsedDate.Format("2006-01-02")
		}
		orderData[i] = map[string]interface{}{
			"date":        dateStr,
			"totalSales":  util.FenToYuan(int(res.TotalSales)),
			"totalOrders": res.TotalOrders,
		}
	}

	// Construct the result map for user data
	userData := make([]map[string]interface{}, len(userResults))
	for i, res := range userResults {
		dateStr := res.Date
		if timeRange == "year" {
			// For year, the date is already in the format "2006-01"
		} else {
			// For week and month, parse the date and format it to "2006-01-02"
			parsedDate, err := time.Parse("2006-01-02T15:04:05Z07:00", res.Date)
			if err != nil {
				return nil, fmt.Errorf("failed to parse date: %v", err)
			}
			dateStr = parsedDate.Format("2006-01-02")
		}
		userData[i] = map[string]interface{}{
			"date":       dateStr,
			"totalUsers": res.TotalUsers,
		}
	}

	// Return the combined order and user data
	return map[string][]map[string]interface{}{
		"orders": orderData,
		"users":  userData,
	}, nil
}
