// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/pkg/api/dto"
	"gotribe-admin/pkg/api/known"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"gotribe-admin/pkg/util"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IOrderController interface {
	GetOrderInfo(c *gin.Context)          // 获取当前订单信息
	GetOrders(c *gin.Context)             // 获取订单列表
	UpdateOrderByID(c *gin.Context)       // 更新订单
	BatchDeleteOrderByIds(c *gin.Context) // 批量删除订单
	GetOrderLogs(c *gin.Context)          //获取订单记录
	UpdateLogistics(c *gin.Context)
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

// GetOrderInfo 获取当前订单信息
// @Summary 获取订单详情
// @Description 根据订单ID获取订单详细信息
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param orderID path string true "订单ID"
// @Success 200 {object} response.Response{data=map[string]dto.OrderDto} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /order/{orderID} [get]
// @Security BearerAuth
func (tc OrderController) GetOrderInfo(c *gin.Context) {
	order, err := tc.OrderRepository.GetOrderByOrderID(c.Param("orderID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	orderInfoDto := dto.ToOrderInfoDto(order)
	response.Success(c, gin.H{
		"order": orderInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// GetOrders 获取订单列表
// @Summary 获取订单列表
// @Description 根据查询条件获取订单列表，支持分页
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param orderID query string false "订单ID"
// @Param userID query string false "用户ID"
// @Param status query int false "订单状态"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /order [get]
// @Security BearerAuth
func (tc OrderController) GetOrders(c *gin.Context) {
	var req vo.OrderListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.Fail(c, nil, errStr)
		return
	}

	// 获取
	order, total, err := tc.OrderRepository.GetOrders(&req)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"orders": dto.ToOrdersDto(order), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// UpdateOrderByID 更新订单
// @Summary 更新订单
// @Description 根据订单ID更新订单信息
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param orderID path string true "订单ID"
// @Param order body vo.CreateOrderRequest true "订单信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /order/{orderID} [patch]
// @Security BearerAuth
func (tc OrderController) UpdateOrderByID(c *gin.Context) {
	var req vo.CreateOrderRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.Fail(c, nil, errStr)
		return
	}

	// 根据path中的OrderID获取订单信息
	oldOrder, err := tc.OrderRepository.GetOrderByOrderID(c.Param("orderID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	oldOrder.AmountPay = util.MoneyUtil.YuanToCents(req.AmountPay)
	oldOrder.Status = req.Status
	oldOrder.RemarkAdmin = req.RemarkAdmin
	// 更新订单
	err = tc.OrderRepository.UpdateOrder(oldOrder)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgUpdateFail)+": "+err.Error())
		return
	}
	// 增加修改记录
	err = tc.OrderLogRepository.CreateOrderLog(c.Param("orderID"), "后台编辑")
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// BatchDeleteOrderByIds 批量删除订单
// @Summary 批量删除订单
// @Description 根据订单ID列表批量删除订单
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param orders body vo.DeleteOrdersRequest true "订单ID列表"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /order [delete]
// @Security BearerAuth
func (tc OrderController) BatchDeleteOrderByIds(c *gin.Context) {
	var req vo.DeleteOrdersRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.Fail(c, nil, errStr)
		return
	}

	// 前端传来的订单ID
	reqOrderIds := strings.Split(req.OrderIds, ",")
	err := tc.OrderRepository.BatchDeleteOrderByIds(reqOrderIds)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgDeleteFail)+": "+err.Error())
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))
}

// GetOrderLogs 获取订单记录
// @Summary 获取订单记录
// @Description 根据订单ID获取订单操作记录
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param orderID path string true "订单ID"
// @Success 200 {object} response.Response{data=map[string]interface{}} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /order/{orderID}/logs [get]
// @Security BearerAuth
func (tc OrderController) GetOrderLogs(c *gin.Context) {
	orderLogs, total, err := tc.OrderLogRepository.GetOrderLogs(c.Param("orderID"))
	if err != nil {
		response.Fail(c, nil, "获取订单记录失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{
		"orderLogs": dto.ToOrderLogsDto(orderLogs),
		"total":     total,
	}, common.Msg(c, common.MsgListSuccess))
}

// UpdateLogistics 更新物流信息
// @Summary 更新物流信息
// @Description 根据订单ID更新物流信息
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param orderID path string true "订单ID"
// @Param logistics body vo.CreateOrderLogisticsRequest true "物流信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /order/{orderID}/logistics [patch]
// @Security BearerAuth
func (tc OrderController) UpdateLogistics(c *gin.Context) {
	var req vo.CreateOrderLogisticsRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.Fail(c, nil, errStr)
		return
	}

	// 根据path中的OrderID获取订单信息
	oldOrder, err := tc.OrderRepository.GetOrderByOrderID(c.Param("orderID"))
	if err != nil {
		response.Fail(c, nil, "获取需要更新的订单信息失败: "+err.Error())
		return
	}
	oldOrder.LogisticsNumber = req.Number
	oldOrder.LogisticsCompany = req.Company
	oldOrder.Status = known.OrderStatusShipped
	// 更新物流信息
	err = tc.OrderRepository.UpdateOrder(oldOrder)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgUpdateFail)+": "+err.Error())
		return
	}
	// 增加修改记录
	err = tc.OrderLogRepository.CreateOrderLog(c.Param("orderID"), "更新物流")
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}
