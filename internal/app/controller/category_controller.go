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

// GetCategoryInfo 获取当前分类信息
// @Summary      获取分类详情
// @Description  根据分类ID获取分类详细信息
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Param        categoryID path string true "分类ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /category/{categoryID} [get]
// @Security     BearerAuth
func (cc CategoryController) GetCategoryInfo(c *gin.Context) {
	category, err := cc.CategoryRepository.GetConfigByCategoryID(c.Param("categoryID"))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	response.Success(c, gin.H{
		"category": category,
	}, common.Msg(c, common.MsgGetSuccess))
}

// GetCategorys 获取分类列表
// @Summary      获取分类列表
// @Description  获取所有分类的列表
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /category [get]
// @Security     BearerAuth
func (cc CategoryController) GetCategorys(c *gin.Context) {
	categorys, err := cc.CategoryRepository.GetCategorys()
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgListFail)
		return
	}
	response.Success(c, gin.H{"categorys": categorys}, common.Msg(c, common.MsgListSuccess))
}

// GetCategoryTree 获取分类树
// @Summary      获取分类树
// @Description  获取树形结构的分类列表
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /category/tree [get]
// @Security     BearerAuth
func (cc CategoryController) GetCategoryTree(c *gin.Context) {
	categoryTree, err := cc.CategoryRepository.GetCategoryTree()
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	response.Success(c, gin.H{"categoryTree": categoryTree}, "获取分类树成功")
}

// CreateCategory 创建分类
// @Summary      创建分类
// @Description  创建一个新的分类
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Param        request body vo.CreateCategoryRequest true "创建分类请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /category [post]
// @Security     BearerAuth
func (cc CategoryController) CreateCategory(c *gin.Context) {
	var req vo.CreateCategoryRequest
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
		response.HandleDatabaseError(c, err, common.MsgCreateFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgCreateSuccess))
}

// UpdateCategoryByID 更新分类
// @Summary      更新分类
// @Description  根据分类ID更新分类信息
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Param        categoryID path string true "分类ID"
// @Param        request body vo.UpdateCategoryRequest true "更新分类请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /category/{categoryID} [patch]
// @Security     BearerAuth
func (cc CategoryController) UpdateCategoryByID(c *gin.Context) {
	var req vo.UpdateCategoryRequest
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
	categoryID := c.Param("categoryID")
	// 校验父级分类ID
	category, err := cc.CategoryRepository.GetConfigByCategoryID(categoryID)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	if req.ParentID == &category.ID {
		response.ValidationFail(c, "不能把自己设为父分类")
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
		response.HandleDatabaseError(c, err, common.MsgUpdateFail)
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))

}

// BatchDeleteCategoryByIds 批量删除分类
// @Summary      批量删除分类
// @Description  根据分类ID列表批量删除分类
// @Tags         分类管理
// @Accept       json
// @Produce      json
// @Param        request body vo.DeleteCategoryRequest true "删除分类请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /category [delete]
// @Security     BearerAuth
func (cc CategoryController) BatchDeleteCategoryByIds(c *gin.Context) {
	var req vo.DeleteCategoryRequest
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
	reqCategoryIds := strings.Split(req.CategoryIds, ",")
	err := cc.CategoryRepository.BatchDeleteCategoryByIds(reqCategoryIds)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgDeleteFail)
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))
}
