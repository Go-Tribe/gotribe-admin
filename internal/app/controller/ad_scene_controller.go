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
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IAdSceneController interface {
	GetAdSceneInfo(c *gin.Context)          // 获取当前登录推广场景信息
	GetAdScenes(c *gin.Context)             // 获取推广场景列表
	CreateAdScene(c *gin.Context)           // 创建推广场景
	UpdateAdSceneByID(c *gin.Context)       // 更新推广场景
	BatchDeleteAdSceneByIds(c *gin.Context) // 批量删除
}

type AdSceneController struct {
	AdSceneRepository repository.IAdSceneRepository
}

// 构造函数
func NewAdSceneController() IAdSceneController {
	adSceneRepository := repository.NewAdSceneRepository()
	adSceneController := AdSceneController{AdSceneRepository: adSceneRepository}
	return adSceneController
}

// GetAdSceneInfo 获取当前推广场景信息
// @Summary 获取推广场景详情
// @Description 根据推广场景ID获取推广场景详细信息
// @Tags 推广场景管理
// @Accept json
// @Produce json
// @Param adSceneID path string true "推广场景ID"
// @Success 200 {object} response.Response{data=map[string]dto.AdSceneDto} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad/scene/{adSceneID} [get]
// @Security BearerAuth
func (pc AdSceneController) GetAdSceneInfo(c *gin.Context) {
	adScene, err := pc.AdSceneRepository.GetAdSceneByAdSceneID(c.Param("adSceneID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	adSceneInfoDto := dto.ToAdSceneInfoDto(adScene)
	response.Success(c, gin.H{
		"adScene": adSceneInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// GetAdScenes 获取推广场景列表
// @Summary 获取推广场景列表
// @Description 根据查询条件获取推广场景列表，支持分页
// @Tags 推广场景管理
// @Accept json
// @Produce json
// @Param ProjectID query string false "项目ID"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad/scene [get]
// @Security BearerAuth
func (pc AdSceneController) GetAdScenes(c *gin.Context) {
	var req vo.AdSceneListRequest
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
	adScene, total, err := pc.AdSceneRepository.GetAdScenes(&req)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"adScenes": dto.ToAdScenesDto(adScene), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// CreateAdScene 创建推广场景
// @Summary 创建推广场景
// @Description 创建新的推广场景
// @Tags 推广场景管理
// @Accept json
// @Produce json
// @Param adScene body vo.CreateAdSceneRequest true "推广场景信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad/scene [post]
// @Security BearerAuth
func (pc AdSceneController) CreateAdScene(c *gin.Context) {
	var req vo.CreateAdSceneRequest
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

	adScene := model.AdScene{
		ProjectID:   req.ProjectID,
		Title:       req.Title,
		Description: req.Description,
	}

	err := pc.AdSceneRepository.CreateAdScene(&adScene)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgCreateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgCreateSuccess))

}

// UpdateAdSceneByID 更新推广场景
// @Summary 更新推广场景
// @Description 根据推广场景ID更新推广场景信息
// @Tags 推广场景管理
// @Accept json
// @Produce json
// @Param adSceneID path string true "推广场景ID"
// @Param adScene body vo.UpdateAdSceneRequest true "推广场景信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad/scene/{adSceneID} [patch]
// @Security BearerAuth
func (pc AdSceneController) UpdateAdSceneByID(c *gin.Context) {
	var req vo.UpdateAdSceneRequest
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

	// 根据path中的AdSceneID获取推广场景信息
	oldAdScene, err := pc.AdSceneRepository.GetAdSceneByAdSceneID(c.Param("adSceneID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	oldAdScene.Title = req.Title
	oldAdScene.Description = req.Description
	// 更新推广场景
	err = pc.AdSceneRepository.UpdateAdScene(&oldAdScene)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgUpdateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// BatchDeleteAdSceneByIds 批量删除推广场景
// @Summary 批量删除推广场景
// @Description 根据推广场景ID列表批量删除推广场景
// @Tags 推广场景管理
// @Accept json
// @Produce json
// @Param adScenes body vo.DeleteAdScenesRequest true "推广场景ID列表"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad/scene [delete]
// @Security BearerAuth
func (pc AdSceneController) BatchDeleteAdSceneByIds(c *gin.Context) {
	var req vo.DeleteAdScenesRequest
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

	// 前端传来的推广场景ID
	reqAdSceneIds := strings.Split(req.AdSceneIds, ",")
	err := pc.AdSceneRepository.BatchDeleteAdSceneByIds(reqAdSceneIds)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgDeleteFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))

}
