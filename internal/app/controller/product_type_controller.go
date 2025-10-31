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

// GetProductTypeInfo 获取当前商品类型信息
// @Summary 获取商品类型详情
// @Description 根据商品类型ID获取商品类型详细信息，包含关联的规格信息
// @Tags 商品类型管理
// @Accept json
// @Produce json
// @Param productTypeID path string true "商品类型ID"
// @Success 200 {object} response.Response{data=map[string]dto.ProductTypeDto} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /product-type/{productTypeID} [get]
// @Security BearerAuth
func (tc ProductTypeController) GetProductTypeInfo(c *gin.Context) {
	productType, err := tc.ProductTypeRepository.GetProductTypeByProductTypeID(c.Param("productTypeID"))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	productTypeInfoDto := dto.ToProductTypeInfoDto(productType)

	// 将 SpecIds 字符串转换为字符串切片
	specIds := strings.Split(productType.SpecIds, ",")

	// 获取关联的规格信息并追加进去
	productSpecs, err := tc.ProductSpecRepository.GetProductSpecsByProductSpecIDs(specIds)
	if err != nil {
		response.InternalServerError(c, "获取商品规格信息失败: "+err.Error())
		return
	}
	productTypeInfoDto.Spec = dto.ToProductSpecsDto(productSpecs)

	response.Success(c, gin.H{
		"productType": productTypeInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// GetProductTypes 获取商品类型列表
// @Summary 获取商品类型列表
// @Description 根据查询条件获取商品类型列表，支持分页
// @Tags 商品类型管理
// @Accept json
// @Produce json
// @Param productTypeID query string false "商品类型ID"
// @Param title query string false "商品类型标题"
// @Param categoryID query string false "分类ID"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /product-type [get]
// @Security BearerAuth
func (tc ProductTypeController) GetProductTypes(c *gin.Context) {
	var req vo.ProductTypeListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.ValidationFail(c, errStr)
		return
	}

	// 获取
	productType, total, err := tc.ProductTypeRepository.GetProductTypes(&req)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgListFail)
		return
	}
	response.Success(c, gin.H{"productTypes": dto.ToProductTypesDto(productType), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// CreateProductType 创建商品类型
// @Summary 创建商品类型
// @Description 创建新的商品类型
// @Tags 商品类型管理
// @Accept json
// @Produce json
// @Param productType body vo.CreateProductTypeRequest true "商品类型信息"
// @Success 200 {object} response.Response{data=map[string]dto.ProductTypeDto} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /product-type [post]
// @Security BearerAuth
func (tc ProductTypeController) CreateProductType(c *gin.Context) {
	var req vo.CreateProductTypeRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.ValidationFail(c, errStr)
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
		response.HandleDatabaseError(c, err, common.MsgCreateFail)
		return
	}
	response.Success(c, gin.H{"productType": dto.ToProductTypeInfoDto(*productTypeInfo)}, common.Msg(c, common.MsgCreateSuccess))
}

// UpdateProductTypeByID 更新商品类型
// @Summary 更新商品类型
// @Description 根据商品类型ID更新商品类型信息
// @Tags 商品类型管理
// @Accept json
// @Produce json
// @Param productTypeID path string true "商品类型ID"
// @Param productType body vo.CreateProductTypeRequest true "商品类型信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /product-type/{productTypeID} [patch]
// @Security BearerAuth
func (tc ProductTypeController) UpdateProductTypeByID(c *gin.Context) {
	var req vo.CreateProductTypeRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.ValidationFail(c, errStr)
		return
	}

	// 根据path中的ProductTypeID获取商品类型信息
	oldProductType, err := tc.ProductTypeRepository.GetProductTypeByProductTypeID(c.Param("productTypeID"))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	oldProductType.Title = req.Title
	oldProductType.Remark = req.Remark
	oldProductType.ProductCategoryID = req.CategoryID
	oldProductType.SpecIds = req.SpecIds
	// 更新商品类型
	err = tc.ProductTypeRepository.UpdateProductType(&oldProductType)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgUpdateFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// BatchDeleteProductTypeByIds 批量删除商品类型
// @Summary 批量删除商品类型
// @Description 根据商品类型ID列表批量删除商品类型
// @Tags 商品类型管理
// @Accept json
// @Produce json
// @Param productTypes body vo.DeleteProductTypesRequest true "商品类型ID列表"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /product-type [delete]
// @Security BearerAuth
func (tc ProductTypeController) BatchDeleteProductTypeByIds(c *gin.Context) {
	var req vo.DeleteProductTypesRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.ValidationFail(c, errStr)
		return
	}

	// 前端传来的商品类型ID
	reqProductTypeIds := strings.Split(req.ProductTypeIds, ",")
	err := tc.ProductTypeRepository.BatchDeleteProductTypeByIds(reqProductTypeIds)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgDeleteFail)
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))

}
