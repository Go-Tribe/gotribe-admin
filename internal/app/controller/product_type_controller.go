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
	"github.com/go-playground/validator/v10"
)

type IProductTypeController interface {
	GetProductTypeInfo(c *gin.Context)          // 获取当前登录商品类型信息
	GetProductTypes(c *gin.Context)             // 获取商品类型列表
	CreateProductType(c *gin.Context)           // 创建商品类型
	UpdateProductTypeByID(c *gin.Context)       // 更新商品类型
	BatchDeleteProductTypeByIds(c *gin.Context) // 批量删除商品类型
}

type ProductTypeController struct {
	ProductTypeRepository repository.IProductTypeRepository
	ProductSpecRepository repository.IProductSpecRepository
}

// 构造函数
func NewProductTypeController() IProductTypeController {
	productTypeRepository := repository.NewProductTypeRepository()
	productSpecRepository := repository.NewProductSpecRepository()
	productTypeController := ProductTypeController{
		ProductTypeRepository: productTypeRepository,
		ProductSpecRepository: productSpecRepository,
	}

	return productTypeController
}

// 获取当前商品类型信息
// 获取当前商品类型信息
func (tc ProductTypeController) GetProductTypeInfo(c *gin.Context) {
	productType, err := tc.ProductTypeRepository.GetProductTypeByProductTypeID(c.Param("productTypeID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	productTypeInfoDto := dto.ToProductTypeInfoDto(productType)

	// 将 SpecIds 字符串转换为字符串切片
	specIds := strings.Split(productType.SpecIds, ",")

	// 获取关联的规格信息并追加进去
	productSpecs, err := tc.ProductSpecRepository.GetProductSpecsByProductSpecIDs(specIds)
	if err != nil {
		response.Fail(c, nil, "获取商品规格信息失败: "+err.Error())
		return
	}
	productTypeInfoDto.Spec = productSpecs

	response.Success(c, gin.H{
		"productType": productTypeInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// 获取商品类型列表
func (tc ProductTypeController) GetProductTypes(c *gin.Context) {
	var req vo.ProductTypeListRequest
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
	productType, total, err := tc.ProductTypeRepository.GetProductTypes(&req)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"productTypes": dto.ToProductTypesDto(productType), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// 创建商品类型
func (tc ProductTypeController) CreateProductType(c *gin.Context) {
	var req vo.CreateProductTypeRequest
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

	productType := model.ProductType{
		Model:             model.Model{},
		Title:             req.Title,
		Remark:            req.Remark,
		ProductCategoryID: req.CategoryID,
		SpecIds:           req.SpecIds,
	}

	productTypeInfo, err := tc.ProductTypeRepository.CreateProductType(&productType)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgCreateFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"productType": dto.ToProductTypeInfoDto(*productTypeInfo)}, common.Msg(c, common.MsgCreateSuccess))
}

// 更新商品类型
func (tc ProductTypeController) UpdateProductTypeByID(c *gin.Context) {
	var req vo.CreateProductTypeRequest
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

	// 根据path中的ProductTypeID获取商品类型信息
	oldProductType, err := tc.ProductTypeRepository.GetProductTypeByProductTypeID(c.Param("productTypeID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	oldProductType.Title = req.Title
	oldProductType.Remark = req.Remark
	oldProductType.ProductCategoryID = req.CategoryID
	oldProductType.SpecIds = req.SpecIds
	// 更新商品类型
	err = tc.ProductTypeRepository.UpdateProductType(&oldProductType)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgUpdateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// 批量删除商品类型
func (tc ProductTypeController) BatchDeleteProductTypeByIds(c *gin.Context) {
	var req vo.DeleteProductTypesRequest
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

	// 前端传来的商品类型ID
	reqProductTypeIds := strings.Split(req.ProductTypeIds, ",")
	err := tc.ProductTypeRepository.BatchDeleteProductTypeByIds(reqProductTypeIds)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgDeleteFail)+": "+err.Error())
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))

}
