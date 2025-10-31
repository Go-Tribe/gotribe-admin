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
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IUserController interface {
	GetUserInfo(c *gin.Context)          // 获取当前登录用户信息
	GetUsers(c *gin.Context)             // 获取用户列表
	CreateUser(c *gin.Context)           // 创建用户
	UpdateUserByID(c *gin.Context)       // 更新用户
	BatchDeleteUserByIds(c *gin.Context) // 批量删除用户
	SearchUserByUsername(c *gin.Context)
}

type UserController struct {
	UserRepository repository.IUserRepository
}

// 构造函数
func NewUserController() IUserController {
	userRepository := repository.NewUserRepository()
	userController := UserController{UserRepository: userRepository}
	return userController
}

// GetUserInfo 获取用户信息
// @Summary      获取用户信息
// @Description  根据用户ID获取用户详细信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        userID path string true "用户ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /user/{userID} [get]
// @Security     BearerAuth
func (pc UserController) GetUserInfo(c *gin.Context) {
	user, err := pc.UserRepository.GetUserByUserID(c.Param("userID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	userInfoDto := dto.ToUserInfoDto(&user)
	response.Success(c, gin.H{
		"user": userInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// GetUsers 获取用户列表
// @Summary      获取用户列表
// @Description  获取所有用户的列表
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        request query vo.UserListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /user [get]
// @Security     BearerAuth
func (pc UserController) GetUsers(c *gin.Context) {
	var req vo.UserListRequest
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
	user, total, err := pc.UserRepository.GetUsers(&req)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"users": dto.ToUsersDto(user), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// CreateUser 创建用户
// @Summary      创建用户
// @Description  创建一个新的用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        request body vo.CreateUserRequest true "创建用户请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /user [post]
// @Security     BearerAuth
func (pc UserController) CreateUser(c *gin.Context) {
	var req vo.CreateUserRequest
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

	user := model.User{
		Username:  req.Username,
		Nickname:  req.Nickname,
		Phone:     req.Phone,
		Email:     req.Email,
		ProjectID: req.ProjectID,
	}

	err := pc.UserRepository.CreateUser(&user)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgCreateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgCreateSuccess))

}

// UpdateUserByID 更新用户
// @Summary      更新用户
// @Description  根据用户ID更新用户信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        userID path string true "用户ID"
// @Param        request body vo.UpdateUserRequest true "更新用户请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /user/{userID} [patch]
// @Security     BearerAuth
func (pc UserController) UpdateUserByID(c *gin.Context) {
	var req vo.UpdateUserRequest
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

	// 根据path中的UserID获取用户信息
	oldUser, err := pc.UserRepository.GetUserByUserID(c.Param("userID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	oldUser.Nickname = req.Nickname
	oldUser.Phone = req.Phone
	oldUser.Email = req.Email
	if len(req.Password) > 0 {
		newPassword, _ := util.PasswordUtil.Encrypt(req.Password)
		oldUser.Password = newPassword
	}

	// 更新用户
	err = pc.UserRepository.UpdateUser(&oldUser)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgUpdateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// BatchDeleteUserByIds 批量删除用户
// @Summary      批量删除用户
// @Description  根据用户ID列表批量删除用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        request body vo.DeleteUsersRequest true "删除用户请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /user [delete]
// @Security     BearerAuth
func (tc UserController) BatchDeleteUserByIds(c *gin.Context) {
	var req vo.DeleteUsersRequest
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

	// 前端传来的标签ID
	reqUserIds := strings.Split(req.UserIds, ",")
	err := tc.UserRepository.BatchDeleteUserByIds(reqUserIds)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgDeleteFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))

}

// SearchUserByUsername 搜索用户
// @Summary      搜索用户
// @Description  根据昵称搜索用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        nickname query string true "用户昵称"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /user/search [get]
// @Security     BearerAuth
func (tc UserController) SearchUserByUsername(c *gin.Context) {
	user, err := tc.UserRepository.SearchUserByNickname(c.Param("nickname"))
	if err != nil {
		response.Fail(c, nil, "获取需要更新的用户信息失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"users": dto.ToUsersDto(user)}, common.Msg(c, common.MsgListSuccess))
}
