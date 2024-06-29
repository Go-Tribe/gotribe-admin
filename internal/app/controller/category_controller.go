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
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type ICategoryController interface {
	GetCategorys(c *gin.Context) // 获取分类列表
	GetCategoryInfo(c *gin.Context)
	GetCategoryTree(c *gin.Context)          // 获取分类树
	CreateCategory(c *gin.Context)           // 创建分类
	UpdateCategoryByID(c *gin.Context)       // 更新分类
	BatchDeleteCategoryByIds(c *gin.Context) // 批量删除分
}

type CategoryController struct {
	CategoryRepository repository.ICategoryRepository
}

func NewCategoryController() ICategoryController {
	categoryRepository := repository.NewCategoryRepository()
	categoryController := CategoryController{CategoryRepository: categoryRepository}
	return categoryController
}

// 获取当前分类信息
func (cc CategoryController) GetCategoryInfo(c *gin.Context) {
	category, err := cc.CategoryRepository.GetConfigByCategoryID(c.Param("categoryID"))
	if err != nil {
		response.Fail(c, nil, "获取当前分类信息失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{
		"category": category,
	}, "获取当前分类信息成功")
}

// 获取分类列表
func (cc CategoryController) GetCategorys(c *gin.Context) {
	categorys, err := cc.CategoryRepository.GetCategorys()
	if err != nil {
		response.Fail(c, nil, "获取分类列表失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"categorys": categorys}, "获取分类列表成功")
}

// 获取分类树
func (cc CategoryController) GetCategoryTree(c *gin.Context) {
	categoryTree, err := cc.CategoryRepository.GetCategoryTree()
	if err != nil {
		response.Fail(c, nil, "获取分类树失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"categoryTree": categoryTree}, "获取分类树成功")
}

// 创建分类
func (cc CategoryController) CreateCategory(c *gin.Context) {
	var req vo.CreateCategoryRequest
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

	category := model.Category{
		Title:       req.Title,
		Icon:        req.Icon,
		Path:        req.Path,
		Sort:        req.Sort,
		Status:      req.Status,
		Hidden:      req.Hidden,
		ParentID:    req.ParentID,
		Description: req.Description,
	}

	err := cc.CategoryRepository.CreateCategory(&category)
	if err != nil {
		response.Fail(c, nil, "创建分类失败: "+err.Error())
		return
	}
	response.Success(c, nil, "创建分类成功")
}

// 更新分类
func (cc CategoryController) UpdateCategoryByID(c *gin.Context) {
	var req vo.UpdateCategoryRequest
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
	categoryID := c.Param("categoryID")
	// 校验父级分类ID
	category, err := cc.CategoryRepository.GetConfigByCategoryID(categoryID)
	if err != nil {
		response.Fail(c, nil, "分类不存在")
		return
	}
	if req.ParentID == &category.ID {
		response.Fail(c, nil, "不能把自己设为父分类")
		return
	}
	category.Title = req.Title
	category.Icon = req.Icon
	category.Path = req.Path
	category.Sort = req.Sort
	category.Status = req.Status
	category.Hidden = req.Hidden
	category.ParentID = req.ParentID
	category.Description = req.Description
	err = cc.CategoryRepository.UpdateCategoryByID(categoryID, &category)
	if err != nil {
		response.Fail(c, nil, "更新分类失败: "+err.Error())
		return
	}

	response.Success(c, nil, "更新分类成功")

}

// 批量删除分类
func (cc CategoryController) BatchDeleteCategoryByIds(c *gin.Context) {
	var req vo.DeleteCategoryRequest
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
	reqCategoryIds := strings.Split(req.CategoryIds, ",")
	err := cc.CategoryRepository.BatchDeleteCategoryByIds(reqCategoryIds)
	if err != nil {
		response.Fail(c, nil, "删除分类失败: "+err.Error())
		return
	}

	response.Success(c, nil, "删除分类成功")
}
