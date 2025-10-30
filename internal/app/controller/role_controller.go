// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"fmt"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thoas/go-funk"
)

type IRoleController interface {
	GetRoles(c *gin.Context)             // 获取角色列表
	CreateRole(c *gin.Context)           // 创建角色
	UpdateRoleByID(c *gin.Context)       // 更新角色
	GetRoleMenusByID(c *gin.Context)     // 获取角色的权限菜单
	UpdateRoleMenusByID(c *gin.Context)  // 更新角色的权限菜单
	GetRoleApisByID(c *gin.Context)      // 获取角色的权限接口
	UpdateRoleApisByID(c *gin.Context)   // 更新角色的权限接口
	BatchDeleteRoleByIds(c *gin.Context) // 批量删除角色
}

type RoleController struct {
	RoleRepository repository.IRoleRepository
}

func NewRoleController() IRoleController {
	roleRepository := repository.NewRoleRepository()
	roleController := RoleController{RoleRepository: roleRepository}
	return roleController
}

// 获取角色列表
func (rc RoleController) GetRoles(c *gin.Context) {
	var req vo.RoleListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.Fail(c, nil, errStr)
		return
	}

	// 获取角色列表
	roles, total, err := rc.RoleRepository.GetRoles(&req)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"roles": roles, "total": total}, common.Msg(c, common.MsgListSuccess))
}

// 创建角色
func (rc RoleController) CreateRole(c *gin.Context) {
	var req vo.CreateRoleRequest
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

	// 获取当前用户最高角色等级
	uc := repository.NewAdminRepository()
	sort, ctxUser, err := uc.GetCurrentAdminMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, "获取当前用户最高角色等级失败: "+err.Error())
		return
	}

	// 用户不能创建比自己等级高或相同等级的角色
	if sort >= req.Sort {
		response.Fail(c, nil, "不能创建比自己等级高或相同等级的角色")
		return
	}

	role := model.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    &req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.Username,
	}

	// 创建角色
	err = rc.RoleRepository.CreateRole(&role)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgCreateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgCreateSuccess))

}

// 更新角色
func (rc RoleController) UpdateRoleByID(c *gin.Context) {
	var req vo.CreateRoleRequest
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
	// 获取path中的roleID
	roleID, _ := strconv.Atoi(c.Param("roleID"))
	if roleID <= 0 {
		response.Fail(c, nil, "角色ID不正确")
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	ur := repository.NewAdminRepository()
	minSort, ctxUser, err := ur.GetCurrentAdminMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// 不能更新比自己角色等级高或相等的角色
	// 根据path中的角色ID获取该角色信息
	roles, err := rc.RoleRepository.GetRolesByIds([]uint{uint(roleID)})
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "未获取到角色信息")
		return
	}
	if minSort >= roles[0].Sort {
		response.Fail(c, nil, "不能更新比自己角色等级高或相等的角色")
		return
	}

	// 不能把角色等级更新得比当前用户的等级高
	if minSort >= req.Sort {
		response.Fail(c, nil, "不能把角色等级更新得比当前用户的等级高或相同")
		return
	}

	role := model.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    &req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.Username,
	}

	// 更新角色
	err = rc.RoleRepository.UpdateRoleByID(uint(roleID), &role)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgUpdateFail)+": "+err.Error())
		return
	}

	// 如果更新成功，且更新了角色的keyword, 则更新casbin中policy
	if req.Keyword != roles[0].Keyword {
		// 获取policy
		rolePolicies := common.CasbinEnforcer.GetFilteredPolicy(0, roles[0].Keyword)
		if len(rolePolicies) == 0 {
			response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
			return
		}
		rolePoliciesCopy := make([][]string, 0)
		// 替换keyword
		for _, policy := range rolePolicies {
			policyCopy := make([]string, len(policy))
			copy(policyCopy, policy)
			rolePoliciesCopy = append(rolePoliciesCopy, policyCopy)
			policy[0] = req.Keyword
		}

		//gormadapter未实现UpdatePolicies方法，等gorm更新---
		//isUpdated, _ := common.CasbinEnforcer.UpdatePolicies(rolePoliciesCopy, rolePolicies)
		//if !isUpdated {
		//	response.Fail(c, nil, "更新角色成功，但角色关键字关联的权限接口更新失败！")
		//	return
		//}

		// 这里需要先新增再删除（先删除再增加会出错）
		isAdded, _ := common.CasbinEnforcer.AddPolicies(rolePolicies)
		if !isAdded {
			response.Fail(c, nil, "更新角色成功，但角色关键字关联的权限接口更新失败")
			return
		}
		isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rolePoliciesCopy)
		if !isRemoved {
			response.Fail(c, nil, "更新角色成功，但角色关键字关联的权限接口更新失败")
			return
		}
		err := common.CasbinEnforcer.LoadPolicy()
		if err != nil {
			response.Fail(c, nil, "更新角色成功，但角色关键字关联角色的权限接口策略加载失败")
			return
		}

	}

	// 更新角色成功处理用户信息缓存有两种做法:（这里使用第二种方法，因为一个角色下用户数量可能很多，第二种方法可以分散数据库压力）
	// 1.可以帮助用户更新拥有该角色的用户信息缓存,使用下面方法
	// err = ur.UpdateUserInfoCacheByRoleID(uint(roleID))
	// 2.直接清理缓存，让活跃的用户自己重新缓存最新用户信息
	ur.ClearAdminInfoCache()

	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// 获取角色的权限菜单
func (rc RoleController) GetRoleMenusByID(c *gin.Context) {
	// 获取path中的roleID
	roleID, _ := strconv.Atoi(c.Param("roleID"))
	if roleID <= 0 {
		response.Fail(c, nil, "角色ID不正确")
		return
	}
	menus, err := rc.RoleRepository.GetRoleMenusByID(uint(roleID))
	if err != nil {
		response.Fail(c, nil, "获取角色的权限菜单失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"menus": menus}, common.Msg(c, common.MsgGetSuccess))
}

// 更新角色的权限菜单
func (rc RoleController) UpdateRoleMenusByID(c *gin.Context) {
	var req vo.UpdateRoleMenusRequest
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
	// 获取path中的roleID
	roleID, _ := strconv.Atoi(c.Param("roleID"))
	if roleID <= 0 {
		response.Fail(c, nil, "角色ID不正确")
		return
	}
	// 根据path中的角色ID获取该角色信息
	roles, err := rc.RoleRepository.GetRolesByIds([]uint{uint(roleID)})
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "未获取到角色信息")
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	ur := repository.NewAdminRepository()
	minSort, ctxUser, err := ur.GetCurrentAdminMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// (非管理员)不能更新比自己角色等级高或相等角色的权限菜单
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			response.Fail(c, nil, "不能更新比自己角色等级高或相等角色的权限菜单")
			return
		}
	}

	// 获取当前用户所拥有的权限菜单
	mr := repository.NewMenuRepository()
	ctxUserMenus, err := mr.GetUserMenusByUserID(ctxUser.ID)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}

	// 获取当前用户所拥有的权限菜单ID
	ctxUserMenusIds := make([]uint, 0)
	for _, menu := range ctxUserMenus {
		ctxUserMenusIds = append(ctxUserMenusIds, menu.ID)
	}

	// 前端传来最新的MenuIds集合
	menuIds := req.MenuIds

	// 用户需要修改的菜单集合
	reqMenus := make([]*model.Menu, 0)

	// (非管理员)不能把角色的权限菜单设置的比当前用户所拥有的权限菜单多
	if minSort != 1 {
		for _, id := range menuIds {
			if !funk.Contains(ctxUserMenusIds, id) {
				response.Fail(c, nil, fmt.Sprintf("无权设置ID为%d的菜单", id))
				return
			}
		}

		for _, id := range menuIds {
			for _, menu := range ctxUserMenus {
				if id == menu.ID {
					reqMenus = append(reqMenus, menu)
					break
				}
			}
		}
	} else {
		// 管理员随意设置
		// 根据menuIds查询查询菜单
		menus, err := mr.GetMenus()
		if err != nil {
			response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
			return
		}
		for _, menuID := range menuIds {
			for _, menu := range menus {
				if menuID == menu.ID {
					reqMenus = append(reqMenus, menu)
				}
			}
		}
	}

	roles[0].Menus = reqMenus

	err = rc.RoleRepository.UpdateRoleMenus(roles[0])
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgUpdateFail)+": "+err.Error())
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))

}

// 获取角色的权限接口
func (rc RoleController) GetRoleApisByID(c *gin.Context) {
	// 获取path中的roleID
	roleID, _ := strconv.Atoi(c.Param("roleID"))
	if roleID <= 0 {
		response.Fail(c, nil, "角色ID不正确")
		return
	}
	// 根据path中的角色ID获取该角色信息
	roles, err := rc.RoleRepository.GetRolesByIds([]uint{uint(roleID)})
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "未获取到角色信息")
		return
	}
	// 根据角色keyword获取casbin中policy
	keyword := roles[0].Keyword
	apis, err := rc.RoleRepository.GetRoleApisByRoleKeyword(keyword)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"apis": apis}, "获取角色的权限接口成功")
}

// 更新角色的权限接口
func (rc RoleController) UpdateRoleApisByID(c *gin.Context) {
	var req vo.UpdateRoleApisRequest
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

	// 获取path中的roleID
	roleID, _ := strconv.Atoi(c.Param("roleID"))
	if roleID <= 0 {
		response.Fail(c, nil, "角色ID不正确")
		return
	}
	// 根据path中的角色ID获取该角色信息
	roles, err := rc.RoleRepository.GetRolesByIds([]uint{uint(roleID)})
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "未获取到角色信息")
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	ur := repository.NewAdminRepository()
	minSort, ctxUser, err := ur.GetCurrentAdminMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// (非管理员)不能更新比自己角色等级高或相等角色的权限接口
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			response.Fail(c, nil, "不能更新比自己角色等级高或相等角色的权限接口")
			return
		}
	}

	// 获取当前用户所拥有的权限接口
	ctxRoles := ctxUser.Roles
	ctxRolesPolicies := make([][]string, 0)
	for _, role := range ctxRoles {
		policy := common.CasbinEnforcer.GetFilteredPolicy(0, role.Keyword)
		ctxRolesPolicies = append(ctxRolesPolicies, policy...)
	}
	// 得到path中的角色ID对应角色能够设置的权限接口集合
	for _, policy := range ctxRolesPolicies {
		policy[0] = roles[0].Keyword
	}

	// 前端传来最新的ApiID集合
	apiIds := req.ApiIds
	// 根据apiID获取接口详情
	ar := repository.NewApiRepository()
	apis, err := ar.GetApisByID(apiIds)
	if err != nil {
		response.Fail(c, nil, "根据接口ID获取接口信息失败")
		return
	}
	// 生成前端想要设置的角色policies
	reqRolePolicies := make([][]string, 0)
	for _, api := range apis {
		reqRolePolicies = append(reqRolePolicies, []string{
			roles[0].Keyword, api.Path, api.Method,
		})
	}

	// (非管理员)不能把角色的权限接口设置的比当前用户所拥有的权限接口多
	if minSort != 1 {
		for _, reqPolicy := range reqRolePolicies {
			if !funk.Contains(ctxRolesPolicies, reqPolicy) {
				response.Fail(c, nil, fmt.Sprintf("无权设置路径为%s,请求方式为%s的接口", reqPolicy[1], reqPolicy[2]))
				return
			}
		}
	}

	// 更新角色的权限接口
	err = rc.RoleRepository.UpdateRoleApis(roles[0].Keyword, reqRolePolicies)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))

}

// 批量删除角色
func (rc RoleController) BatchDeleteRoleByIds(c *gin.Context) {
	var req vo.DeleteRoleRequest
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

	// 获取当前用户最高等级角色
	ur := repository.NewAdminRepository()
	minSort, _, err := ur.GetCurrentAdminMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// 前端传来需要删除的角色ID
	roleIds := req.RoleIds
	// 获取角色信息
	roles, err := rc.RoleRepository.GetRolesByIds(roleIds)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "未获取到角色信息")
		return
	}

	// 不能删除比自己角色等级高或相等的角色
	for _, role := range roles {
		if minSort >= role.Sort {
			response.Fail(c, nil, "不能删除比自己角色等级高或相等的角色")
			return
		}
	}

	// 删除角色
	err = rc.RoleRepository.BatchDeleteRoleByIds(roleIds)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgDeleteFail))
		return
	}

	// 删除角色成功直接清理缓存，让活跃的用户自己重新缓存最新用户信息
	ur.ClearAdminInfoCache()
	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))

}
