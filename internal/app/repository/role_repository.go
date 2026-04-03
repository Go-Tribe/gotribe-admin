// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"
	"fmt"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"

	"strings"
)

type IRoleRepository interface {
	GetRoles(ctx context.Context, req *vo.RoleListRequest) ([]model.Role, int64, error)       // 获取角色列表
	GetRolesByIds(ctx context.Context, roleIds []uint) ([]*model.Role, error)                 // 根据角色ID获取角色
	CreateRole(ctx context.Context, role *model.Role) error                                   // 创建角色
	UpdateRoleByID(ctx context.Context, roleID uint, role *model.Role) error                  // 更新角色
	GetRoleMenusByID(ctx context.Context, roleID uint) ([]*model.Menu, error)                 // 获取角色的权限菜单
	UpdateRoleMenus(ctx context.Context, role *model.Role) error                              // 更新角色的权限菜单
	GetRoleApisByRoleKeyword(ctx context.Context, roleKeyword string) ([]*model.Api, error)   // 根据角色关键字获取角色的权限接口
	UpdateRoleApis(ctx context.Context, roleKeyword string, reqRolePolicies [][]string) error // 更新角色的权限接口（先全部删除再新增）
	BatchDeleteRoleByIds(ctx context.Context, roleIds []uint) error                           // 删除角色
}

type RoleRepository struct {
}

func NewRoleRepository() IRoleRepository {
	return RoleRepository{}
}

func buildRoleOrder(req *vo.RoleListRequest) string {
	sortByMap := map[string]string{
		"name":       "name",
		"keyword":    "keyword",
		"sort":       "sort",
		"status":     "status",
		"creator":    "creator",
		"createdAt":  "created_at",
		"created_at": "created_at",
	}

	column, ok := sortByMap[strings.TrimSpace(req.SortBy)]
	if !ok {
		return "created_at DESC"
	}

	direction := "ASC"
	if strings.EqualFold(strings.TrimSpace(req.SortOrder), "desc") {
		direction = "DESC"
	}

	return fmt.Sprintf("%s %s", column, direction)
}

// 获取角色列表
func (r RoleRepository) GetRoles(ctx context.Context, req *vo.RoleListRequest) ([]model.Role, int64, error) {
	var list []model.Role
	db := common.WithContext(ctx).DB().Model(&model.Role{}).Order(buildRoleOrder(req))

	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		db = db.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

// 根据角色ID获取角色
func (r RoleRepository) GetRolesByIds(ctx context.Context, roleIds []uint) ([]*model.Role, error) {
	var list []*model.Role
	err := common.WithContext(ctx).DB().Where("id IN (?)", roleIds).Find(&list).Error
	return list, err
}

// 创建角色
func (r RoleRepository) CreateRole(ctx context.Context, role *model.Role) error {
	err := common.WithContext(ctx).DB().Create(role).Error
	return err
}

// 更新角色
func (r RoleRepository) UpdateRoleByID(ctx context.Context, roleID uint, role *model.Role) error {
	err := common.WithContext(ctx).DB().Model(&model.Role{}).Where("id = ?", roleID).Updates(role).Error
	return err
}

// 获取角色的权限菜单
func (r RoleRepository) GetRoleMenusByID(ctx context.Context, roleID uint) ([]*model.Menu, error) {
	var role model.Role
	err := common.WithContext(ctx).DB().Where("id = ?", roleID).Preload("Menus").First(&role).Error
	return role.Menus, err
}

// 更新角色的权限菜单
func (r RoleRepository) UpdateRoleMenus(ctx context.Context, role *model.Role) error {
	err := common.WithContext(ctx).DB().Model(role).Association("Menus").Replace(role.Menus)
	return err
}

// 根据角色关键字获取角色的权限接口
func (r RoleRepository) GetRoleApisByRoleKeyword(ctx context.Context, roleKeyword string) ([]*model.Api, error) {
	policies, _ := common.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)

	// 获取所有接口
	var apis []*model.Api
	err := common.WithContext(ctx).DB().Find(&apis).Error
	if err != nil {
		return apis, common.ErrGetRoleApisFailed
	}

	accessApis := make([]*model.Api, 0)

	for _, policy := range policies {
		path := policy[1]
		method := policy[2]
		for _, api := range apis {
			if path == api.Path && method == api.Method {
				accessApis = append(accessApis, api)
				break
			}
		}
	}

	return accessApis, err

}

// 更新角色的权限接口（先全部删除再新增）
func (r RoleRepository) UpdateRoleApis(ctx context.Context, roleKeyword string, reqRolePolicies [][]string) error {
	// 先获取path中的角色ID对应角色已有的police(需要先删除的)
	err := common.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return common.ErrLoadRolePolicyFailed
	}
	rmPolicies, _ := common.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
	if len(rmPolicies) > 0 {
		isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rmPolicies)
		if !isRemoved {
			return common.ErrUpdateRoleApisFailed
		}
	}
	isAdded, _ := common.CasbinEnforcer.AddPolicies(reqRolePolicies)
	if !isAdded {
		return common.ErrUpdateRoleApisFailed
	}
	err = common.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return common.ErrLoadRolePolicyFailed
	} else {
		return err
	}
}

// 删除角色
func (r RoleRepository) BatchDeleteRoleByIds(ctx context.Context, roleIds []uint) error {
	var roles []*model.Role
	err := common.WithContext(ctx).DB().Where("id IN (?)", roleIds).Find(&roles).Error
	if err != nil {
		return err
	}
	err = common.WithContext(ctx).DB().Select("Users", "Menus").Unscoped().Delete(&roles).Error
	// 删除成功就删除casbin policy
	if err == nil {
		for _, role := range roles {
			roleKeyword := role.Keyword
			rmPolicies, _ := common.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
			if len(rmPolicies) > 0 {
				isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rmPolicies)
				if !isRemoved {
					return common.ErrDeleteRoleApisFailed
				}
			}
		}

	}
	return err
}
