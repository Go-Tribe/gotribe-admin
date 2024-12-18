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
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/dto"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"strings"
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

// 获取当前商品规格信息
func (tc ProductSpecItemController) GetProductSpecItemInfo(c *gin.Context) {
	productSpecItem, err := tc.ProductSpecItemRepository.GetProductSpecItemByProductSpecItemID(c.Param("productSpecItemID"))
	if err != nil {
		response.Fail(c, nil, "获取当前商品规格信息失败: "+err.Error())
		return
	}
	productSpecItemInfoDto := dto.ToProductSpecItemInfoDto(productSpecItem)
	response.Success(c, gin.H{
		"productSpecItem": productSpecItemInfoDto,
	}, "获取当前商品规格信息成功")
}

// 获取商品规格列表
func (tc ProductSpecItemController) GetProductSpecItems(c *gin.Context) {
	var req vo.ProductSpecItemListRequest
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
	productSpecItem, total, err := tc.ProductSpecItemRepository.GetProductSpecItems(&req)
	if err != nil {
		response.Fail(c, nil, "获取商品规格列表失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"productSpecItems": dto.ToProductSpecItemsDto(productSpecItem), "total": total}, "获取商品规格列表成功")
}

// 创建商品规格
func (tc ProductSpecItemController) CreateProductSpecItem(c *gin.Context) {
	var req vo.CreateProductSpecItemRequest
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

	productSpecItem := model.ProductSpecItem{
		Model:   model.Model{},
		Title:   req.Title,
		Sort:    req.Sort,
		SpecID:  req.SpecID,
		Enabled: req.Enabled,
	}

	productSpecItemInfo, err := tc.ProductSpecItemRepository.CreateProductSpecItem(&productSpecItem)
	if err != nil {
		response.Fail(c, nil, "创建商品规格失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"productSpecItem": dto.ToProductSpecItemInfoDto(*productSpecItemInfo)}, "创建商品规格成功")
}

// 更新商品规格
func (tc ProductSpecItemController) UpdateProductSpecItemByID(c *gin.Context) {
	var req vo.CreateProductSpecItemRequest
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

	// 根据path中的ProductSpecItemID获取商品规格信息
	oldProductSpecItem, err := tc.ProductSpecItemRepository.GetProductSpecItemByProductSpecItemID(c.Param("productSpecItemID"))
	if err != nil {
		response.Fail(c, nil, "获取需要更新的商品规格信息失败: "+err.Error())
		return
	}
	oldProductSpecItem.Title = req.Title
	oldProductSpecItem.Enabled = req.Enabled
	oldProductSpecItem.Sort = req.Sort
	// 更新商品规格
	err = tc.ProductSpecItemRepository.UpdateProductSpecItem(&oldProductSpecItem)
	if err != nil {
		response.Fail(c, nil, "更新商品规格失败: "+err.Error())
		return
	}
	response.Success(c, nil, "更新商品规格成功")
}

// 批量删除商品规格
func (tc ProductSpecItemController) BatchDeleteProductSpecItemByIds(c *gin.Context) {
	var req vo.DeleteProductSpecItemsRequest
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

	// 前端传来的商品规格ID
	reqProductSpecItemIds := strings.Split(req.ProductSpecItemIds, ",")
	err := tc.ProductSpecItemRepository.BatchDeleteProductSpecItemByIds(reqProductSpecItemIds)
	if err != nil {
		response.Fail(c, nil, "删除商品规格失败: "+err.Error())
		return
	}

	response.Success(c, nil, "删除商品规格成功")

}
