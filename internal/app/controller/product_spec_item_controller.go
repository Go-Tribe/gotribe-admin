// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/dto"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"strings"

	"github.com/gin-gonic/gin"
)

type IProductSpecItemController interface {
	GetProductSpecItemInfo(c *gin.Context)          // 获取当前商品规格信息
	GetProductSpecItems(c *gin.Context)             // 获取商品规格列表
	CreateProductSpecItem(c *gin.Context)           // 创建商品规格
	UpdateProductSpecItemByID(c *gin.Context)       // 更新商品规格
	BatchDeleteProductSpecItemByIds(c *gin.Context) // 批量删除商品规格
}

type ProductSpecItemController struct {
	ProductSpecItemRepository repository.IProductSpecItemRepository
}

// 构造函数
func NewProductSpecItemController() IProductSpecItemController {
	productSpecItemRepository := repository.NewProductSpecItemRepository()
	productSpecItemController := ProductSpecItemController{ProductSpecItemRepository: productSpecItemRepository}
	return productSpecItemController
}

// GetProductSpecItemInfo 获取当前商品规格信息
// @Summary 获取商品规格信息
// @Description 根据商品规格ID获取商品规格详细信息
// @Tags 商品规格项管理
// @Accept json
// @Produce json
// @Param productSpecItemID path string true "商品规格项ID"
// @Success 200 {object} response.Response{data=object{productSpecItem=dto.ProductSpecItemDto}} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/product-spec-items/{productSpecItemID} [get]
// @Security BearerAuth
func (tc ProductSpecItemController) GetProductSpecItemInfo(c *gin.Context) {
	productSpecItem, err := tc.ProductSpecItemRepository.GetProductSpecItemByItemID(c.Param("productSpecItemID"))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	productSpecItemInfoDto := dto.ToProductSpecItemInfoDto(&productSpecItem)
	response.Success(c, gin.H{
		"productSpecItem": productSpecItemInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// GetProductSpecItems 获取商品规格列表
// @Summary 获取商品规格列表
// @Description 获取商品规格项列表，支持分页和筛选
// @Tags 商品规格项管理
// @Accept json
// @Produce json
// @Param request query vo.ProductSpecItemListRequest false "查询参数"
// @Success 200 {object} response.Response{data=object{productSpecItems=[]dto.ProductSpecItemDto,total=int}} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/product-spec-items [get]
// @Security BearerAuth
func (tc ProductSpecItemController) GetProductSpecItems(c *gin.Context) {
	var req vo.ProductSpecItemListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	// 获取
	productSpecItem, total, err := tc.ProductSpecItemRepository.GetProductSpecItems(&req)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgListFail)
		return
	}
	response.Success(c, gin.H{"productSpecItems": dto.ToProductSpecItemsDto(productSpecItem), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// CreateProductSpecItem 创建商品规格
// @Summary 创建商品规格
// @Description 创建新的商品规格项
// @Tags 商品规格项管理
// @Accept json
// @Produce json
// @Param request body vo.CreateProductSpecItemRequest true "创建商品规格请求参数"
// @Success 200 {object} response.Response{data=object{productSpecItem=dto.ProductSpecItemDto}} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/product-spec-items [post]
// @Security BearerAuth
func (tc ProductSpecItemController) CreateProductSpecItem(c *gin.Context) {
	var req vo.CreateProductSpecItemRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	productSpecItem := model.ProductSpecItem{
		Title:   req.Title,
		Sort:    req.Sort,
		SpecID:  req.SpecID,
		Enabled: req.Enabled,
	}

	productSpecItemInfo, err := tc.ProductSpecItemRepository.CreateProductSpecItem(&productSpecItem)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgCreateFail)
		return
	}
	response.Success(c, gin.H{"productSpecItem": dto.ToProductSpecItemInfoDto(productSpecItemInfo)}, common.Msg(c, common.MsgCreateSuccess))
}

// UpdateProductSpecItemByID 更新商品规格
// @Summary 更新商品规格
// @Description 根据商品规格ID更新商品规格信息
// @Tags 商品规格项管理
// @Accept json
// @Produce json
// @Param productSpecItemID path string true "商品规格项ID"
// @Param request body vo.CreateProductSpecItemRequest true "更新商品规格请求参数"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/product-spec-items/{productSpecItemID} [put]
// @Security BearerAuth
func (tc ProductSpecItemController) UpdateProductSpecItemByID(c *gin.Context) {
	var req vo.CreateProductSpecItemRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	// 根据path中的ProductSpecItemID获取商品规格信息
	oldProductSpecItem, err := tc.ProductSpecItemRepository.GetProductSpecItemByItemID(c.Param("productSpecItemID"))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	oldProductSpecItem.Title = req.Title
	oldProductSpecItem.Enabled = req.Enabled
	oldProductSpecItem.Sort = req.Sort
	// 更新商品规格
	err = tc.ProductSpecItemRepository.UpdateProductSpecItem(&oldProductSpecItem)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgUpdateFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// BatchDeleteProductSpecItemByIds 批量删除商品规格
// @Summary 批量删除商品规格
// @Description 根据商品规格ID列表批量删除商品规格项
// @Tags 商品规格项管理
// @Accept json
// @Produce json
// @Param request body vo.DeleteProductSpecItemsRequest true "删除商品规格请求参数"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/product-spec-items [delete]
// @Security BearerAuth
func (tc ProductSpecItemController) BatchDeleteProductSpecItemByIds(c *gin.Context) {
	var req vo.DeleteProductSpecItemsRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	// 前端传来的商品规格ID
	reqProductSpecItemIds := strings.Split(req.ProductSpecItemIds, ",")
	err := tc.ProductSpecItemRepository.BatchDeleteProductSpecItemByIds(reqProductSpecItemIds)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgDeleteFail)
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))

}
