// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"fmt"
	"github.com/thoas/go-funk"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
	"gotribe-admin/pkg/api/vo"
	"strings"
	"time"
)

type IOrderRepository interface {
	GetOrderByOrderID(orderID string) (*model.Order, error)            // 获取单个订单
	GetOrders(req *vo.OrderListRequest) ([]*model.Order, int64, error) // 获取订单列表
	UpdateOrder(order *model.Order) error                              // 更新订单
	BatchDeleteOrderByIds(ids []string) error                          // 批量删除
}

type OrderRepository struct {
}

// OrderRepository构造函数
func NewOrderRepository() IOrderRepository {
	return OrderRepository{}
}

// 获取单个订单
func (tr OrderRepository) GetOrderByOrderID(orderID string) (*model.Order, error) {
	var order model.Order
	err := common.DB.Where("order_id = ?", orderID).First(&order).Error
	return getOrdertUser(&order), err
}

// 获取订单详情里的用户信息
func getOrdertUser(order *model.Order) *model.Order {
	// 通过 order.userID 获取用户信息
	var user model.User
	err := common.DB.Where("user_id = ?", order.UserID).First(&user).Error
	if err != nil {
		return order
	}
	order.User = &user
	return order
}

// 获取订单列表
func (tr OrderRepository) GetOrders(req *vo.OrderListRequest) ([]*model.Order, int64, error) {
	var list []*model.Order
	db := common.DB.Model(&model.Order{}).Order("created_at DESC")

	orderID := strings.TrimSpace(req.OrderNumber)
	if req.OrderNumber != "" {
		db = db.Where("order_number = ?", orderID)
	}
	if req.Title != "" {
		db = db.Where("product_name LIKE ?", fmt.Sprintf("%%%s%%", req.Title))
	}
	if len(req.StartTime) > 0 {
		t, _ := time.Parse(known.TIME_FORMAT_SHORT, req.StartTime)
		db = db.Where("date(created_at) >= ?", t)
	}
	if len(req.EndTime) > 0 {
		t, _ := time.Parse(known.TIME_FORMAT_SHORT, req.StartTime)
		db = db.Where("date(created_at) >= ?", t)
	}
	if req.UserID != "" {
		db = db.Where("user_id = ?", req.UserID)
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return getOrdertOther(list), total, err
}

func getOrdertOther(orders []*model.Order) []*model.Order {
	// 拿出所有用户ID，去重后去 user表查出用户信息
	var userIDs []string
	for _, order := range orders {
		userIDs = append(userIDs, order.UserID)
	}
	userIDs = funk.UniqString(userIDs)
	var users []model.User
	err := common.DB.Where("user_id in (?)", userIDs).Find(&users).Error
	if err != nil {
		return orders
	}

	// 使用映射存储用户信息
	userMap := make(map[string]model.User)
	for _, user := range users {
		userMap[user.UserID] = user
	}

	// 分配用户信息
	for _, order := range orders {
		if user, exists := userMap[order.UserID]; exists {
			order.User = &user
		}
	}
	return orders
}

// 更新订单
func (tr OrderRepository) UpdateOrder(order *model.Order) error {
	err := common.DB.Model(order).Updates(order).Error
	if err != nil {
		return err
	}
	return err
}

// 批量删除
func (tr OrderRepository) BatchDeleteOrderByIds(ids []string) error {
	var orders []*model.Order
	for _, id := range ids {
		// 根据ID获取订单
		order, err := tr.GetOrderByOrderID(id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%s的订单", id)
		}
		orders = append(orders, order)
	}

	err := common.DB.Unscoped().Delete(&orders).Error

	return err
}
