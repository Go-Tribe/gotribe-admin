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

	"strconv"
)

type IMenuController interface {
	GetMenus(c *gin.Context)             // 获取菜单列表
	GetMenuTree(c *gin.Context)          // 获取菜单树
	CreateMenu(c *gin.Context)           // 创建菜单
	UpdateMenuByID(c *gin.Context)       // 更新菜单
	BatchDeleteMenuByIds(c *gin.Context) // 批量删除菜单

	GetUserMenusByUserID(c *gin.Context)    // 获取用户的可访问菜单列表
	GetUserMenuTreeByUserID(c *gin.Context) // 获取用户的可访问菜单树
}

type MenuController struct {
	MenuRepository repository.IMenuRepository
}

func NewMenuController() IMenuController {
	menuRepository := repository.NewMenuRepository()
	menuController := MenuController{MenuRepository: menuRepository}
	return menuController
}

// 获取菜单列表
func (mc MenuController) GetMenus(c *gin.Context) {
	menus, err := mc.MenuRepository.GetMenus()
	if err != nil {
		response.Fail(c, nil, "获取菜单列表失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"menus": menus}, "获取菜单列表成功")
}

// 获取菜单树
func (mc MenuController) GetMenuTree(c *gin.Context) {
	menuTree, err := mc.MenuRepository.GetMenuTree()
	if err != nil {
		response.Fail(c, nil, "获取菜单树失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"menuTree": menuTree}, "获取菜单树成功")
}

// 创建菜单
func (mc MenuController) CreateMenu(c *gin.Context) {
	var req vo.CreateMenuRequest
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

	// 获取当前用户
	ur := repository.NewAdminRepository()
	ctxUser, err := ur.GetCurrentAdmin(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户信息失败")
		return
	}

	menu := model.Menu{
		Name:       req.Name,
		Title:      req.Title,
		Icon:       &req.Icon,
		Path:       req.Path,
		Redirect:   &req.Redirect,
		Component:  req.Component,
		Sort:       req.Sort,
		Status:     req.Status,
		Hidden:     req.Hidden,
		NoCache:    req.NoCache,
		AlwaysShow: req.AlwaysShow,
		Breadcrumb: req.Breadcrumb,
		ActiveMenu: &req.ActiveMenu,
		ParentID:   &req.ParentID,
		Creator:    ctxUser.Username,
	}

	err = mc.MenuRepository.CreateMenu(&menu)
	if err != nil {
		response.Fail(c, nil, "创建菜单失败: "+err.Error())
		return
	}
	response.Success(c, nil, "创建菜单成功")
}

// 更新菜单
func (mc MenuController) UpdateMenuByID(c *gin.Context) {
	var req vo.UpdateMenuRequest
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

	// 获取路径中的menuID
	menuID, _ := strconv.Atoi(c.Param("menuID"))
	if menuID <= 0 {
		response.Fail(c, nil, "菜单ID不正确")
		return
	}

	// 获取当前用户
	ur := repository.NewAdminRepository()
	ctxUser, err := ur.GetCurrentAdmin(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户信息失败")
		return
	}

	menu := model.Menu{
		Name:       req.Name,
		Title:      req.Title,
		Icon:       &req.Icon,
		Path:       req.Path,
		Redirect:   &req.Redirect,
		Component:  req.Component,
		Sort:       req.Sort,
		Status:     req.Status,
		Hidden:     req.Hidden,
		NoCache:    req.NoCache,
		AlwaysShow: req.AlwaysShow,
		Breadcrumb: req.Breadcrumb,
		ActiveMenu: &req.ActiveMenu,
		ParentID:   &req.ParentID,
		Creator:    ctxUser.Username,
	}

	err = mc.MenuRepository.UpdateMenuByID(uint(menuID), &menu)
	if err != nil {
		response.Fail(c, nil, "更新菜单失败: "+err.Error())
		return
	}

	response.Success(c, nil, "更新菜单成功")

}

// 批量删除菜单
func (mc MenuController) BatchDeleteMenuByIds(c *gin.Context) {
	var req vo.DeleteMenuRequest
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
	err := mc.MenuRepository.BatchDeleteMenuByIds(req.MenuIds)
	if err != nil {
		response.Fail(c, nil, "删除菜单失败: "+err.Error())
		return
	}

	response.Success(c, nil, "删除菜单成功")
}

// 根据用户ID获取用户的可访问菜单列表
func (mc MenuController) GetUserMenusByUserID(c *gin.Context) {
	// 获取路径中的userID
	userID, _ := strconv.Atoi(c.Param("userID"))
	if userID <= 0 {
		response.Fail(c, nil, "用户ID不正确")
		return
	}

	menus, err := mc.MenuRepository.GetUserMenusByUserID(uint(userID))
	if err != nil {
		response.Fail(c, nil, "获取用户的可访问菜单列表失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"menus": menus}, "获取用户的可访问菜单列表成功")
}

// 根据用户ID获取用户的可访问菜单树
func (mc MenuController) GetUserMenuTreeByUserID(c *gin.Context) {
	// 获取路径中的userID
	userID, _ := strconv.Atoi(c.Param("userID"))
	if userID <= 0 {
		response.Fail(c, nil, "用户ID不正确")
		return
	}

	menuTree, err := mc.MenuRepository.GetUserMenuTreeByUserID(uint(userID))
	if err != nil {
		response.Fail(c, nil, "获取用户的可访问菜单树失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"menuTree": menuTree}, "获取用户的可访问菜单树成功")
}
