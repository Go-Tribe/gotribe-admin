// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/dto"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"gotribe-admin/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thoas/go-funk"

	"strconv"
)

type IAdminController interface {
	GetAdminInfo(c *gin.Context)          // 获取当前登录用户信息
	GetAdmins(c *gin.Context)             // 获取用户列表
	ChangePwd(c *gin.Context)             // 更新用户登录密码
	CreateAdmin(c *gin.Context)           // 创建用户
	UpdateAdminByID(c *gin.Context)       // 更新用户
	BatchDeleteAdminByIds(c *gin.Context) // 批量删除用户
}

type AdminController struct {
	AdminRepository repository.IAdminRepository
}

// 构造函数
func NewAdminController() IAdminController {
	userRepository := repository.NewAdminRepository()
	userController := AdminController{AdminRepository: userRepository}
	return userController
}

// GetAdminInfo 获取当前登录管理员信息
// @Summary      获取当前登录管理员信息
// @Description  获取当前登录管理员详细信息
// @Tags         管理员管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /admin/info [post]
// @Security     BearerAuth
func (uc AdminController) GetAdminInfo(c *gin.Context) {
	user, err := uc.AdminRepository.GetCurrentAdmin(c)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	userInfoDto := dto.ToAdminInfoDto(user)
	response.Success(c, gin.H{
		"admin": userInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// GetAdmins 获取管理员列表
// @Summary      获取管理员列表
// @Description  获取所有管理员的列表
// @Tags         管理员管理
// @Accept       json
// @Produce      json
// @Param        request query vo.AdminListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /admin/list [get]
// @Security     BearerAuth
func (uc AdminController) GetAdmins(c *gin.Context) {
	var req vo.AdminListRequest
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

	// 获取
	users, total, err := uc.AdminRepository.GetAdmins(&req)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"admins": dto.ToAdminsDto(users), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// ChangePwd 更新管理员登录密码
// @Summary      更新管理员登录密码
// @Description  更新当前管理员的登录密码
// @Tags         管理员管理
// @Accept       json
// @Produce      json
// @Param        request body vo.ChangePwdRequest true "更新密码请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /admin/changePwd [put]
// @Security     BearerAuth
func (uc AdminController) ChangePwd(c *gin.Context) {
	var req vo.ChangePwdRequest

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

	// 直接使用明文密码进行校验与更新

	// 获取当前用户
	user, err := uc.AdminRepository.GetCurrentAdmin(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 获取用户的真实正确密码
	correctPasswd := user.Password
	// 判断前端请求的密码是否等于真实密码
	err = util.PasswordUtil.ComparePasswd(correctPasswd, req.OldPassword)
	if err != nil {
		response.Fail(c, nil, "原密码有误")
		return
	}
	// 更新密码
	hashedPassword, err := util.PasswordUtil.GenPasswd(req.NewPassword)
	if err != nil {
		response.Fail(c, nil, "密码加密失败: "+err.Error())
		return
	}
	err = uc.AdminRepository.ChangePwd(user.Username, hashedPassword)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgUpdateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// CreateAdmin 创建管理员
// @Summary      创建管理员
// @Description  创建一个新的管理员
// @Tags         管理员管理
// @Accept       json
// @Produce      json
// @Param        request body vo.CreateAdminRequest true "创建管理员请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /admin/create [post]
// @Security     BearerAuth
func (uc AdminController) CreateAdmin(c *gin.Context) {
	var req vo.CreateAdminRequest
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

	// 直接使用明文密码；如提供则校验长度
	if req.Password != "" {
		if len(req.Password) < 6 {
			response.Fail(c, nil, "密码长度至少为6位")
			return
		}
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	currentRoleSortMin, ctxAdmin, err := uc.AdminRepository.GetCurrentAdminMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// 获取前端传来的用户角色id
	reqRoleIds := req.RoleIds
	// 根据角色id获取角色
	rr := repository.NewRoleRepository()
	roles, err := rr.GetRolesByIds(reqRoleIds)
	if err != nil {
		response.Fail(c, nil, "根据角色ID获取角色信息失败: "+err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "未获取到角色信息")
		return
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// 前端传来用户角色排序最小值（最高等级角色）
	reqRoleSortMin := uint(funk.MinInt(reqRoleSorts))

	// 当前用户的角色排序最小值 需要小于 前端传来的角色排序最小值（用户不能创建比自己等级高的或者相同等级的用户）
	if currentRoleSortMin >= reqRoleSortMin {
		response.Fail(c, nil, "用户不能创建比自己等级高的或者相同等级的用户")
		return
	}

	// 密码为空就默认123456
	if req.Password == "" {
		req.Password = "123456"
	}
	hashedPassword, err := util.PasswordUtil.GenPasswd(req.Password)
	if err != nil {
		response.Fail(c, nil, "密码加密失败: "+err.Error())
		return
	}
	user := model.Admin{
		Username:     req.Username,
		Password:     hashedPassword,
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		Nickname:     &req.Nickname,
		Introduction: &req.Introduction,
		Status:       req.Status,
		Creator:      ctxAdmin.Username,
		Roles:        roles,
	}

	err = uc.AdminRepository.CreateAdmin(&user)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgCreateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgCreateSuccess))

}

// UpdateAdminByID 更新管理员
// @Summary      更新管理员
// @Description  根据管理员ID更新管理员信息
// @Tags         管理员管理
// @Accept       json
// @Produce      json
// @Param        userID path string true "管理员ID"
// @Param        request body vo.CreateAdminRequest true "更新管理员请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /admin/update/{userID} [patch]
// @Security     BearerAuth
func (uc AdminController) UpdateAdminByID(c *gin.Context) {
	var req vo.CreateAdminRequest
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

	//获取path中的userID
	userID, _ := strconv.Atoi(c.Param("userID"))
	if userID <= 0 {
		response.Fail(c, nil, "用户ID不正确")
		return
	}

	// 根据path中的userID获取用户信息
	oldAdmin, err := uc.AdminRepository.GetAdminByID(uint(userID))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}

	// 获取当前用户
	ctxAdmin, err := uc.AdminRepository.GetCurrentAdmin(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 获取当前用户的所有角色
	currentRoles := ctxAdmin.Roles
	// 获取当前用户角色的排序，和前端传来的角色排序做比较
	var currentRoleSorts []int
	// 当前用户角色ID集合
	var currentRoleIds []uint
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
		currentRoleIds = append(currentRoleIds, role.ID)
	}
	// 当前用户角色排序最小值（最高等级角色）
	currentRoleSortMin := funk.MinInt(currentRoleSorts)

	// 获取前端传来的用户角色id
	reqRoleIds := req.RoleIds
	// 根据角色id获取角色
	rr := repository.NewRoleRepository()
	roles, err := rr.GetRolesByIds(reqRoleIds)
	if err != nil {
		response.Fail(c, nil, "根据角色ID获取角色信息失败: "+err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "未获取到角色信息")
		return
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// 前端传来用户角色排序最小值（最高等级角色）
	reqRoleSortMin := funk.MinInt(reqRoleSorts)

	user := model.Admin{
		Model:        oldAdmin.Model,
		Username:     req.Username,
		Password:     oldAdmin.Password,
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		Nickname:     &req.Nickname,
		Introduction: &req.Introduction,
		Status:       req.Status,
		Creator:      ctxAdmin.Username,
		Roles:        roles,
	}
	// 判断是更新自己还是更新别人
	if userID == int(ctxAdmin.ID) {
		// 如果是更新自己
		// 不能禁用自己
		if req.Status == 2 {
			response.Fail(c, nil, "不能禁用自己")
			return
		}
		// 不能更改自己的角色
		reqDiff, currentDiff := funk.Difference(req.RoleIds, currentRoleIds)
		if len(reqDiff.([]uint)) > 0 || len(currentDiff.([]uint)) > 0 {
			response.Fail(c, nil, "不能更改自己的角色")
			return
		}

		// 不能更新自己的密码，只能在个人中心更新
		if req.Password != "" {
			response.Fail(c, nil, "请到个人中心更新自身密码")
			return
		}

		// 密码赋值
		user.Password = ctxAdmin.Password

	} else {
		// 如果是更新别人
		// 用户不能更新比自己角色等级高的或者相同等级的用户
		// 根据path中的userID获取用户角色排序最小值
		minRoleSorts, err := uc.AdminRepository.GetAdminMinRoleSortsByIds([]uint{uint(userID)})
		if err != nil || len(minRoleSorts) == 0 {
			response.Fail(c, nil, "根据用户ID获取用户角色排序最小值失败")
			return
		}
		if currentRoleSortMin >= minRoleSorts[0] {
			response.Fail(c, nil, "用户不能更新比自己角色等级高的或者相同等级的用户")
			return
		}

		// 用户不能把别的用户角色等级更新得比自己高或相等
		if currentRoleSortMin >= reqRoleSortMin {
			response.Fail(c, nil, "用户不能把别的用户角色等级更新得比自己高或相等")
			return
		}

		// 密码赋值
		if req.Password != "" {
			hashedPassword, err := util.PasswordUtil.GenPasswd(req.Password)
			if err != nil {
				response.Fail(c, nil, "密码加密失败: "+err.Error())
				return
			}
			user.Password = hashedPassword
		}

	}

	// 更新用户
	err = uc.AdminRepository.UpdateAdmin(&user)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgUpdateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))

}

// BatchDeleteAdminByIds 批量删除管理员
// @Summary      批量删除管理员
// @Description  根据管理员ID列表批量删除管理员
// @Tags         管理员管理
// @Accept       json
// @Produce      json
// @Param        request body vo.DeleteAdminRequest true "删除管理员请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /admin/delete/batch [delete]
// @Security     BearerAuth
func (uc AdminController) BatchDeleteAdminByIds(c *gin.Context) {
	var req vo.DeleteAdminRequest
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

	// 前端传来的用户ID
	reqAdminIds := req.UserIds
	// 根据用户ID获取用户角色排序最小值
	roleMinSortList, err := uc.AdminRepository.GetAdminMinRoleSortsByIds(reqAdminIds)
	if err != nil || len(roleMinSortList) == 0 {
		response.Fail(c, nil, "根据用户ID获取用户角色排序最小值失败")
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	minSort, ctxAdmin, err := uc.AdminRepository.GetCurrentAdminMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	currentRoleSortMin := int(minSort)

	// 不能删除自己
	if funk.Contains(reqAdminIds, ctxAdmin.ID) {
		response.Fail(c, nil, "用户不能删除自己")
		return
	}

	// 不能删除比自己角色排序低(等级高)的用户
	for _, sort := range roleMinSortList {
		if currentRoleSortMin >= sort {
			response.Fail(c, nil, "用户不能删除比自己角色等级高的用户")
			return
		}
	}

	err = uc.AdminRepository.BatchDeleteAdminByIds(reqAdminIds)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgDeleteFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))

}
