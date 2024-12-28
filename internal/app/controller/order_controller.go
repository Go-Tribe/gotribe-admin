// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/pkg/api/dto"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type IOrderController interface {
	GetOrderInfo(c *gin.Context)          // 获取当前订单信息
	GetOrders(c *gin.Context)             // 获取订单列表
	UpdateOrderByID(c *gin.Context)       // 更新订单
	BatchDeleteOrderByIds(c *gin.Context) // 批量删除订单
	GetOrderLogs(c *gin.Context)          //获取订单记录
}

type OrderController struct {
	OrderRepository    repository.IOrderRepository
	OrderLogRepository repository.IOrderLogRepository
}

// 构造函数
func NewOrderController() IOrderController {
	orderRepository := repository.NewOrderRepository()
	orderLogRepository := repository.NewOrderLogRepository()
	orderController := OrderController{
		OrderRepository:    orderRepository,
		OrderLogRepository: orderLogRepository,
	}
	return orderController
}

// 获取当前订单信息
func (tc OrderController) GetOrderInfo(c *gin.Context) {
	order, err := tc.OrderRepository.GetOrderByOrderID(c.Param("orderID"))
	if err != nil {
		response.Fail(c, nil, "获取当前订单信息失败: "+err.Error())
		return
	}
	orderInfoDto := dto.ToOrderInfoDto(&order)
	response.Success(c, gin.H{
		"order": orderInfoDto,
	}, "获取当前订单信息成功")
}

// 获取订单列表
func (tc OrderController) GetOrders(c *gin.Context) {
	var req vo.OrderListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}

	// 获取
	order, total, err := tc.OrderRepository.GetOrders(&req)
	if err != nil {
		response.Fail(c, nil, "获取订单列表失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"orders": dto.ToOrdersDto(order), "total": total}, "获取订单列表成功")
}

// 更新订单
func (tc OrderController) UpdateOrderByID(c *gin.Context) {
	var req vo.CreateOrderRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}

	// 根据path中的OrderID获取订单信息
	oldOrder, err := tc.OrderRepository.GetOrderByOrderID(c.Param("orderID"))
	if err != nil {
		response.Fail(c, nil, "获取需要更新的订单信息失败: "+err.Error())
		return
	}
	oldOrder.AmountPay = uint(req.AmountPay / 100)
	oldOrder.Status = req.Status
	oldOrder.Remark = req.RemarkAdmin
	// 更新订单
	err = tc.OrderRepository.UpdateOrder(&oldOrder)
	if err != nil {
		response.Fail(c, nil, "更新订单失败: "+err.Error())
		return
	}
	response.Success(c, nil, "更新订单成功")
}

// 批量删除订单
func (tc OrderController) BatchDeleteOrderByIds(c *gin.Context) {
	var req vo.DeleteOrdersRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}

	// 前端传来的订单ID
	reqOrderIds := strings.Split(req.OrderIds, ",")
	err := tc.OrderRepository.BatchDeleteOrderByIds(reqOrderIds)
	if err != nil {
		response.Fail(c, nil, "删除订单失败: "+err.Error())
		return
	}

	response.Success(c, nil, "删除订单成功")
}

// 获取订单记录
func (tc OrderController) GetOrderLogs(c *gin.Context) {
	var req vo.OrderLogListRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	orderLogs, total, err := tc.OrderLogRepository.GetOrderLogs(&req)
	if err != nil {
		response.Fail(c, nil, "获取订单记录失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{
		"orderLogs": dto.ToOrderLogsDto(orderLogs),
		"total":     total,
	}, "获取订单记录成功")
}
