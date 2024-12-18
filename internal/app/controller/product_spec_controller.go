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

type IProductSpecController interface {
	GetProductSpecInfo(c *gin.Context)          // 获取当前商品规格信息
	GetProductSpecs(c *gin.Context)             // 获取商品规格列表
	CreateProductSpec(c *gin.Context)           // 创建商品规格
	UpdateProductSpecByID(c *gin.Context)       // 更新商品规格
	BatchDeleteProductSpecByIds(c *gin.Context) // 批量删除商品规格
}

type ProductSpecController struct {
	ProductSpecRepository repository.IProductSpecRepository
}

// 构造函数
func NewProductSpecController() IProductSpecController {
	productSpecRepository := repository.NewProductSpecRepository()
	productSpecController := ProductSpecController{ProductSpecRepository: productSpecRepository}
	return productSpecController
}

// 获取当前商品规格信息
func (tc ProductSpecController) GetProductSpecInfo(c *gin.Context) {
	productSpec, err := tc.ProductSpecRepository.GetProductSpecByProductSpecID(c.Param("productSpecID"))
	if err != nil {
		response.Fail(c, nil, "获取当前商品规格信息失败: "+err.Error())
		return
	}
	productSpecInfoDto := dto.ToProductSpecInfoDto(productSpec)
	response.Success(c, gin.H{
		"productSpec": productSpecInfoDto,
	}, "获取当前商品规格信息成功")
}

// 获取商品规格列表
func (tc ProductSpecController) GetProductSpecs(c *gin.Context) {
	var req vo.ProductSpecListRequest
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
	productSpec, total, err := tc.ProductSpecRepository.GetProductSpecs(&req)
	if err != nil {
		response.Fail(c, nil, "获取商品规格列表失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"productSpecs": dto.ToProductSpecsDto(productSpec), "total": total}, "获取商品规格列表成功")
}

// 创建商品规格
func (tc ProductSpecController) CreateProductSpec(c *gin.Context) {
	var req vo.CreateProductSpecRequest
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

	productSpec := model.ProductSpec{
		Model:  model.Model{},
		Title:  req.Title,
		Remark: req.Remark,
		Format: req.Format,
		Image:  req.Image,
		Sort:   req.Sort,
	}

	productSpecInfo, err := tc.ProductSpecRepository.CreateProductSpec(&productSpec)
	if err != nil {
		response.Fail(c, nil, "创建商品规格失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"productSpec": dto.ToProductSpecInfoDto(*productSpecInfo)}, "创建商品规格成功")
}

// 更新商品规格
func (tc ProductSpecController) UpdateProductSpecByID(c *gin.Context) {
	var req vo.CreateProductSpecRequest
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

	// 根据path中的ProductSpecID获取商品规格信息
	oldProductSpec, err := tc.ProductSpecRepository.GetProductSpecByProductSpecID(c.Param("productSpecID"))
	if err != nil {
		response.Fail(c, nil, "获取需要更新的商品规格信息失败: "+err.Error())
		return
	}
	oldProductSpec.Title = req.Title
	oldProductSpec.Remark = req.Remark
	oldProductSpec.Format = req.Format
	oldProductSpec.Image = req.Image
	oldProductSpec.Sort = req.Sort
	// 更新商品规格
	err = tc.ProductSpecRepository.UpdateProductSpec(&oldProductSpec)
	if err != nil {
		response.Fail(c, nil, "更新商品规格失败: "+err.Error())
		return
	}
	response.Success(c, nil, "更新商品规格成功")
}

// 批量删除商品规格
func (tc ProductSpecController) BatchDeleteProductSpecByIds(c *gin.Context) {
	var req vo.DeleteProductSpecRequest
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
	reqProductSpecIds := strings.Split(req.ProductSpecIds, ",")
	err := tc.ProductSpecRepository.BatchDeleteProductSpecByIds(reqProductSpecIds)
	if err != nil {
		response.Fail(c, nil, "删除商品规格失败: "+err.Error())
		return
	}

	response.Success(c, nil, "删除商品规格成功")

}
