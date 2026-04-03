// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"

	"github.com/thoas/go-funk"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

type IMenuRepository interface {
	GetMenus(ctx context.Context) ([]*model.Menu, error)                     // 获取菜单列表
	GetMenuTree(ctx context.Context) ([]*model.Menu, error)                  // 获取菜单树
	CreateMenu(ctx context.Context, menu *model.Menu) error                  // 创建菜单
	UpdateMenuByID(ctx context.Context, menuID uint, menu *model.Menu) error // 更新菜单
	BatchDeleteMenuByIds(ctx context.Context, menuIds []uint) error          // 批量删除菜单

	GetUserMenusByUserID(ctx context.Context, userID uint) ([]*model.Menu, error)    // 根据用户ID获取用户的权限(可访问)菜单列表
	GetUserMenuTreeByUserID(ctx context.Context, userID uint) ([]*model.Menu, error) // 根据用户ID获取用户的权限(可访问)菜单树
}

type MenuRepository struct {
}

func NewMenuRepository() IMenuRepository {
	return MenuRepository{}
}

// 获取菜单列表
func (m MenuRepository) GetMenus(ctx context.Context) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := common.WithContext(ctx).DB().Order("sort").Find(&menus).Error
	return menus, err
}

// 获取菜单树
func (m MenuRepository) GetMenuTree(ctx context.Context) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := common.WithContext(ctx).DB().Order("sort").Find(&menus).Error
	// parentID为0的是根菜单
	return GenMenuTree(0, menus), err
}

func GenMenuTree(parentID uint, menus []*model.Menu) []*model.Menu {
	tree := make([]*model.Menu, 0)

	for _, m := range menus {
		if *m.ParentID == parentID {
			children := GenMenuTree(m.ID, menus)
			m.Children = children
			tree = append(tree, m)
		}
	}
	return tree
}

// 创建菜单
func (m MenuRepository) CreateMenu(ctx context.Context, menu *model.Menu) error {
	err := common.WithContext(ctx).DB().Create(menu).Error
	return err
}

// 更新菜单
func (m MenuRepository) UpdateMenuByID(ctx context.Context, menuID uint, menu *model.Menu) error {
	err := common.WithContext(ctx).DB().Model(menu).Where("id = ?", menuID).Updates(menu).Error
	return err
}

// 批量删除菜单
func (m MenuRepository) BatchDeleteMenuByIds(ctx context.Context, menuIds []uint) error {
	var menus []*model.Menu
	err := common.WithContext(ctx).DB().Where("id IN (?)", menuIds).Find(&menus).Error
	if err != nil {
		return err
	}
	err = common.WithContext(ctx).DB().Select("Roles").Unscoped().Delete(&menus).Error
	return err
}

// 根据用户ID获取用户的权限(可访问)菜单列表
func (m MenuRepository) GetUserMenusByUserID(ctx context.Context, userID uint) ([]*model.Menu, error) {
	// 获取用户
	var user model.Admin
	err := common.WithContext(ctx).DB().Where("id = ?", userID).Preload("Roles").First(&user).Error
	if err != nil {
		return nil, err
	}
	// 获取角色
	roles := user.Roles
	// 所有角色的菜单集合
	allRoleMenus := make([]*model.Menu, 0)
	for _, role := range roles {
		var userRole model.Role
		err := common.WithContext(ctx).DB().Where("id = ?", role.ID).Preload("Menus").First(&userRole).Error
		if err != nil {
			return nil, err
		}
		// 获取角色的菜单
		menus := userRole.Menus
		allRoleMenus = append(allRoleMenus, menus...)
	}

	// 所有角色的菜单集合去重
	allRoleMenusID := make([]int, 0)
	for _, menu := range allRoleMenus {
		allRoleMenusID = append(allRoleMenusID, int(menu.ID))
	}
	allRoleMenusIDUniq := funk.UniqInt(allRoleMenusID)
	allRoleMenusUniq := make([]*model.Menu, 0)
	for _, id := range allRoleMenusIDUniq {
		for _, menu := range allRoleMenus {
			if id == int(menu.ID) {
				allRoleMenusUniq = append(allRoleMenusUniq, menu)
				break
			}
		}
	}

	// 获取状态status为1的菜单
	accessMenus := make([]*model.Menu, 0)
	for _, menu := range allRoleMenusUniq {
		if menu.Status == known.DEFAULT_ID {
			accessMenus = append(accessMenus, menu)
		}
	}

	return accessMenus, err
}

// 根据用户ID获取用户的权限(可访问)菜单树
func (m MenuRepository) GetUserMenuTreeByUserID(ctx context.Context, userID uint) ([]*model.Menu, error) {
	menus, err := m.GetUserMenusByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	tree := GenMenuTree(0, menus)
	return tree, err
}
