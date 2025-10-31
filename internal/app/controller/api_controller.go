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

	"strconv"
)

type IApiController interface {
	GetApis(c *gin.Context)             // 获取接口列表
	GetApiTree(c *gin.Context)          // 获取接口树(按接口Category字段分类)
	CreateApi(c *gin.Context)           // 创建接口
	UpdateApiByID(c *gin.Context)       // 更新接口
	BatchDeleteApiByIds(c *gin.Context) // 批量删除接口
}

type ApiController struct {
	ApiRepository repository.IApiRepository
}

func NewApiController() IApiController {
	apiRepository := repository.NewApiRepository()
	apiController := ApiController{ApiRepository: apiRepository}
	return apiController
}

// GetApis 获取接口列表
// @Summary      获取接口列表
// @Description  获取所有接口的列表
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Param        request query vo.ApiListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /api/list [get]
// @Security     BearerAuth
func (ac ApiController) GetApis(c *gin.Context) {
	var req vo.ApiListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	// 获取
	apis, total, err := ac.ApiRepository.GetApis(&req)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgListFail)
		return
	}
	response.Success(c, gin.H{
		"apis": apis, "total": total,
	}, common.Msg(c, common.MsgListSuccess))
}

// GetApiTree 获取接口树
// @Summary      获取接口树
// @Description  获取树形结构的接口列表(按接口Category字段分类)
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /api/tree [get]
// @Security     BearerAuth
func (ac ApiController) GetApiTree(c *gin.Context) {
	tree, err := ac.ApiRepository.GetApiTree()
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	response.Success(c, gin.H{
		"apis": tree,
	}, common.Msg(c, common.MsgGetSuccess))
}

// CreateApi 创建接口
// @Summary      创建接口
// @Description  创建一个新的接口
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Param        request body vo.CreateApiRequest true "创建接口请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /api/create [post]
// @Security     BearerAuth
func (ac ApiController) CreateApi(c *gin.Context) {
	var req vo.CreateApiRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	// 获取当前用户
	ur := repository.NewAdminRepository()
	ctxUser, err := ur.GetCurrentAdmin(c)
	if err != nil {
		response.InternalServerError(c, "获取当前用户信息失败")
		return
	}

	api := model.Api{
		Method:   req.Method,
		Path:     req.Path,
		Category: req.Category,
		Desc:     req.Desc,
		Creator:  ctxUser.Username,
	}

	// 创建接口
	err = ac.ApiRepository.CreateApi(&api)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgCreateFail)
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgCreateSuccess))
	return
}

// UpdateApiByID 更新接口
// @Summary      更新接口
// @Description  根据接口ID更新接口信息
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Param        apiID path string true "接口ID"
// @Param        request body vo.UpdateApiRequest true "更新接口请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /api/update/{apiID} [patch]
// @Security     BearerAuth
func (ac ApiController) UpdateApiByID(c *gin.Context) {
	var req vo.UpdateApiRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	// 获取路径中的apiID
	apiID, _ := strconv.Atoi(c.Param("apiID"))
	if apiID <= 0 {
		response.ValidationFail(c, "接口ID不正确")
		return
	}

	// 获取当前用户
	ur := repository.NewAdminRepository()
	ctxUser, err := ur.GetCurrentAdmin(c)
	if err != nil {
		response.InternalServerError(c, "获取当前用户信息失败")
		return
	}

	api := model.Api{
		Method:   req.Method,
		Path:     req.Path,
		Category: req.Category,
		Desc:     req.Desc,
		Creator:  ctxUser.Username,
	}

	err = ac.ApiRepository.UpdateApiByID(uint(apiID), &api)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgUpdateFail)
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// BatchDeleteApiByIds 批量删除接口
// @Summary      批量删除接口
// @Description  根据接口ID列表批量删除接口
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Param        request body vo.DeleteApiRequest true "删除接口请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /api/delete/batch [delete]
// @Security     BearerAuth
func (ac ApiController) BatchDeleteApiByIds(c *gin.Context) {
	var req vo.DeleteApiRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	// 删除接口
	err := ac.ApiRepository.BatchDeleteApiByIds(req.ApiIds)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgDeleteFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))
}
