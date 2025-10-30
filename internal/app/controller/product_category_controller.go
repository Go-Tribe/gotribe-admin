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

// 获取当前分类信息
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

// 获取分类列表
func (cc ProductCategoryController) GetProductCategorys(c *gin.Context) {
	productCategorys, err := cc.ProductCategoryRepository.GetProductCategorys()
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"productCategorys": productCategorys}, common.Msg(c, common.MsgListSuccess))
}

// 获取分类树
func (cc ProductCategoryController) GetProductCategoryTree(c *gin.Context) {
	productCategoryTree, err := cc.ProductCategoryRepository.GetProductCategoryTree()
	if err != nil {
		response.Fail(c, nil, "获取分类树失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"productCategoryTree": productCategoryTree}, "获取分类树成功")
}

// 创建分类
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

// 更新分类
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

// 批量删除分类
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
