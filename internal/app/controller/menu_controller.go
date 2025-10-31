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

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

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

// GetMenus 获取菜单列表
// @Summary      获取菜单列表
// @Description  获取所有菜单的列表
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /menu/list [get]
// @Security     BearerAuth
func (mc MenuController) GetMenus(c *gin.Context) {
	menus, err := mc.MenuRepository.GetMenus()
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"menus": menus}, common.Msg(c, common.MsgListSuccess))
}

// GetMenuTree 获取菜单树
// @Summary      获取菜单树
// @Description  获取树形结构的菜单列表
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /menu/tree [get]
// @Security     BearerAuth
func (mc MenuController) GetMenuTree(c *gin.Context) {
	menuTree, err := mc.MenuRepository.GetMenuTree()
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"menuTree": menuTree}, common.Msg(c, common.MsgGetSuccess))
}

// CreateMenu 创建菜单
// @Summary      创建菜单
// @Description  创建一个新的菜单
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        request body vo.CreateMenuRequest true "创建菜单请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /menu/create [post]
// @Security     BearerAuth
func (mc MenuController) CreateMenu(c *gin.Context) {
	var req vo.CreateMenuRequest
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
		response.Fail(c, nil, common.Msg(c, common.MsgCreateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgCreateSuccess))
}

// UpdateMenuByID 更新菜单
// @Summary      更新菜单
// @Description  根据菜单ID更新菜单信息
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        menuID path string true "菜单ID"
// @Param        request body vo.UpdateMenuRequest true "更新菜单请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /menu/update/{menuID} [patch]
// @Security     BearerAuth
func (mc MenuController) UpdateMenuByID(c *gin.Context) {
	var req vo.UpdateMenuRequest
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
		response.Fail(c, nil, common.Msg(c, common.MsgUpdateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))

}

// BatchDeleteMenuByIds 批量删除菜单
// @Summary      批量删除菜单
// @Description  根据菜单ID列表批量删除菜单
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        request body vo.DeleteMenuRequest true "删除菜单请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /menu/delete/batch [delete]
// @Security     BearerAuth
func (mc MenuController) BatchDeleteMenuByIds(c *gin.Context) {
	var req vo.DeleteMenuRequest
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
	err := mc.MenuRepository.BatchDeleteMenuByIds(req.MenuIds)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgDeleteFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))
}

// GetUserMenusByUserID 获取用户的可访问菜单列表
// @Summary      获取用户的可访问菜单列表
// @Description  根据用户ID获取用户的可访问菜单列表
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        userID path string true "用户ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /menu/access/list/{userID} [get]
// @Security     BearerAuth
func (mc MenuController) GetUserMenusByUserID(c *gin.Context) {
	// 获取路径中的userID
	userID, _ := strconv.Atoi(c.Param("userID"))
	if userID <= 0 {
		response.Fail(c, nil, "用户ID不正确")
		return
	}

	menus, err := mc.MenuRepository.GetUserMenusByUserID(uint(userID))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"menus": menus}, common.Msg(c, common.MsgListSuccess))
}

// GetUserMenuTreeByUserID 获取用户的可访问菜单树
// @Summary      获取用户的可访问菜单树
// @Description  根据用户ID获取用户的可访问菜单树
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        userID path string true "用户ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /menu/access/tree/{userID} [get]
// @Security     BearerAuth
func (mc MenuController) GetUserMenuTreeByUserID(c *gin.Context) {
	// 获取路径中的userID
	userID, _ := strconv.Atoi(c.Param("userID"))
	if userID <= 0 {
		response.Fail(c, nil, "用户ID不正确")
		return
	}

	menuTree, err := mc.MenuRepository.GetUserMenuTreeByUserID(uint(userID))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"menuTree": menuTree}, "获取用户的可访问菜单树成功")
}
