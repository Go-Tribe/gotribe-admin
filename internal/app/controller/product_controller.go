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

type IProductController interface {
	GetProductInfo(c *gin.Context)          // 获取当前产品信息
	GetProducts(c *gin.Context)             // 获取产品列表
	CreateProduct(c *gin.Context)           // 创建产品
	UpdateProductByID(c *gin.Context)       // 更新产品
	BatchDeleteProductByIds(c *gin.Context) // 批量删除产品
}

type ProductController struct {
	ProductRepository     repository.IProductRepository
	ProductSpecRepository repository.IProductSpecRepository
}

// 构造函数
func NewProductController() IProductController {
	productRepository := repository.NewProductRepository()
	productSpecRepository := repository.NewProductSpecRepository()
	productController := ProductController{
		ProductRepository:     productRepository,
		ProductSpecRepository: productSpecRepository,
	}

	return productController
}

// 获取当前产品信息
// 获取当前产品信息
func (tc ProductController) GetProductInfo(c *gin.Context) {
	product, err := tc.ProductRepository.GetProductByProductID(c.Param("productID"))
	if err != nil {
		response.Fail(c, nil, "获取当前产品信息失败: "+err.Error())
		return
	}
	productInfoDto := dto.ToProductInfoDto(product)

	//// 将 SpecIds 字符串转换为字符串切片
	//specIds := strings.Split(product.SpecIds, ",")
	//
	//// 获取关联的规格信息并追加进去
	//productSpecs, err := tc.ProductSpecRepository.GetProductSpecsByProductSpecIDs(specIds)
	//if err != nil {
	//	response.Fail(c, nil, "获取商品规格信息失败: "+err.Error())
	//	return
	//}
	//productInfoDto.Spec = productSpecs

	response.Success(c, gin.H{
		"product": productInfoDto,
	}, "获取当前产品信息成功")
}

// 获取产品列表
func (tc ProductController) GetProducts(c *gin.Context) {
	var req vo.ProductListRequest
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
	product, total, err := tc.ProductRepository.GetProducts(&req)
	if err != nil {
		response.Fail(c, nil, "获取产品列表失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"products": dto.ToProductsDto(product), "total": total}, "获取产品列表成功")
}

// 创建产品
func (tc ProductController) CreateProduct(c *gin.Context) {
	var req vo.CreateProductRequest
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

	product := model.Product{
		Model: model.Model{},
		Title: req.Title,
	}

	productInfo, err := tc.ProductRepository.CreateProduct(&product)
	if err != nil {
		response.Fail(c, nil, "创建产品失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"product": dto.ToProductInfoDto(*productInfo)}, "创建产品成功")
}

// 更新产品
func (tc ProductController) UpdateProductByID(c *gin.Context) {
	var req vo.CreateProductRequest
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

	// 根据path中的ProductID获取产品信息
	oldProduct, err := tc.ProductRepository.GetProductByProductID(c.Param("productID"))
	if err != nil {
		response.Fail(c, nil, "获取需要更新的产品信息失败: "+err.Error())
		return
	}
	oldProduct.Title = req.Title
	// 更新产品
	err = tc.ProductRepository.UpdateProduct(&oldProduct)
	if err != nil {
		response.Fail(c, nil, "更新产品失败: "+err.Error())
		return
	}
	response.Success(c, nil, "更新产品成功")
}

// 批量删除产品
func (tc ProductController) BatchDeleteProductByIds(c *gin.Context) {
	var req vo.DeleteProductsRequest
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

	// 前端传来的产品ID
	reqProductIds := strings.Split(req.ProductIds, ",")
	err := tc.ProductRepository.BatchDeleteProductByIds(reqProductIds)
	if err != nil {
		response.Fail(c, nil, "删除产品失败: "+err.Error())
		return
	}

	response.Success(c, nil, "删除产品成功")

}
