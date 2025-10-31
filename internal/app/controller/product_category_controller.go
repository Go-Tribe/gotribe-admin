// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IProductCategoryController interface {
	GetProductCategorys(c *gin.Context) // 获取分类列表
	GetProductCategoryInfo(c *gin.Context)
	GetProductCategoryTree(c *gin.Context)          // 获取分类树
	CreateProductCategory(c *gin.Context)           // 创建分类
	UpdateProductCategoryByID(c *gin.Context)       // 更新分类
	BatchDeleteProductCategoryByIds(c *gin.Context) // 批量删除分
}

type ProductCategoryController struct {
	ProductCategoryRepository repository.IProductCategoryRepository
}

func NewProductCategoryController() IProductCategoryController {
	productCategoryRepository := repository.NewProductCategoryRepository()
	productCategoryController := ProductCategoryController{ProductCategoryRepository: productCategoryRepository}
	return productCategoryController
}

// GetProductCategoryInfo 获取当前分类信息
// @Summary 获取分类信息
// @Description 根据分类ID获取分类详细信息
// @Tags 商品分类管理
// @Accept json
// @Produce json
// @Param productCategoryID path string true "分类ID"
// @Success 200 {object} response.Response{data=object{product_category=object}} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/product-categories/{productCategoryID} [get]
// @Security BearerAuth
func (cc ProductCategoryController) GetProductCategoryInfo(c *gin.Context) {
	productCategory, err := cc.ProductCategoryRepository.GetConfigByProductCategoryID(c.Param("productCategoryID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{
		"product_category": productCategory,
	}, common.Msg(c, common.MsgGetSuccess))
}

// GetProductCategorys 获取分类列表
// @Summary 获取分类列表
// @Description 获取所有商品分类列表
// @Tags 商品分类管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=object{productCategorys=[]object}} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/product-categories [get]
// @Security BearerAuth
func (cc ProductCategoryController) GetProductCategorys(c *gin.Context) {
	productCategorys, err := cc.ProductCategoryRepository.GetProductCategorys()
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"productCategorys": productCategorys}, common.Msg(c, common.MsgListSuccess))
}

// GetProductCategoryTree 获取分类树
// @Summary 获取分类树
// @Description 获取商品分类的树形结构
// @Tags 商品分类管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=object{productCategoryTree=[]object}} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/product-categories/tree [get]
// @Security BearerAuth
func (cc ProductCategoryController) GetProductCategoryTree(c *gin.Context) {
	productCategoryTree, err := cc.ProductCategoryRepository.GetProductCategoryTree()
	if err != nil {
		response.Fail(c, nil, "获取分类树失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"productCategoryTree": productCategoryTree}, "获取分类树成功")
}

// CreateProductCategory 创建分类
// @Summary 创建分类
// @Description 创建新的商品分类
// @Tags 商品分类管理
// @Accept json
// @Produce json
// @Param request body vo.CreateProductCategoryRequest true "创建分类请求参数"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/product-categories [post]
// @Security BearerAuth
func (cc ProductCategoryController) CreateProductCategory(c *gin.Context) {
	var req vo.CreateProductCategoryRequest
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

	productCategory := model.ProductCategory{
		Title:       req.Title,
		Icon:        req.Icon,
		Path:        req.Path,
		ProjectID:   req.ProjectID,
		Sort:        req.Sort,
		Status:      req.Status,
		Hidden:      req.Hidden,
		ParentID:    req.ParentID,
		Description: req.Description,
	}

	err := cc.ProductCategoryRepository.CreateProductCategory(&productCategory)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgCreateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgCreateSuccess))
}

// UpdateProductCategoryByID 更新分类
// @Summary 更新分类
// @Description 根据分类ID更新分类信息
// @Tags 商品分类管理
// @Accept json
// @Produce json
// @Param productCategoryID path string true "分类ID"
// @Param request body vo.UpdateProductCategoryRequest true "更新分类请求参数"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/product-categories/{productCategoryID} [put]
// @Security BearerAuth
func (cc ProductCategoryController) UpdateProductCategoryByID(c *gin.Context) {
	var req vo.UpdateProductCategoryRequest
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
	productCategoryID := c.Param("productCategoryID")
	// 校验父级分类ID
	productCategory, err := cc.ProductCategoryRepository.GetConfigByProductCategoryID(productCategoryID)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	if req.ParentID == &productCategory.ID {
		response.Fail(c, nil, "不能把自己设为父分类")
		return
	}
	productCategory.Title = req.Title
	productCategory.Icon = req.Icon
	productCategory.Path = req.Path
	productCategory.Sort = req.Sort
	productCategory.Status = req.Status
	productCategory.Hidden = req.Hidden
	productCategory.ParentID = req.ParentID
	productCategory.Description = req.Description
	productCategory.ProjectID = req.ProjectID
	err = cc.ProductCategoryRepository.UpdateProductCategoryByID(productCategoryID, &productCategory)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgUpdateFail)+": "+err.Error())
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))

}

// BatchDeleteProductCategoryByIds 批量删除分类
// @Summary 批量删除分类
// @Description 根据分类ID列表批量删除分类
// @Tags 商品分类管理
// @Accept json
// @Produce json
// @Param request body vo.DeleteProductCategoryRequest true "删除分类请求参数"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/product-categories [delete]
// @Security BearerAuth
func (cc ProductCategoryController) BatchDeleteProductCategoryByIds(c *gin.Context) {
	var req vo.DeleteProductCategoryRequest
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
	reqProductCategoryIds := strings.Split(req.ProductCategoryIds, ",")
	err := cc.ProductCategoryRepository.BatchDeleteProductCategoryByIds(reqProductCategoryIds)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgDeleteFail)+": "+err.Error())
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))
}
