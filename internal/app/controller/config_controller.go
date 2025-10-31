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

type IConfigController interface {
	GetConfigInfo(c *gin.Context)          // 获取当前登录配置信息
	GetConfigs(c *gin.Context)             // 获取配置列表
	CreateConfig(c *gin.Context)           // 创建配置
	UpdateConfigByID(c *gin.Context)       // 更新配置
	BatchDeleteConfigByIds(c *gin.Context) // 批量删除
}

type ConfigController struct {
	ConfigRepository repository.IConfigRepository
}

// 构造函数
func NewConfigController() IConfigController {
	configRepository := repository.NewConfigRepository()
	cnfigController := ConfigController{ConfigRepository: configRepository}
	return cnfigController
}

// GetConfigInfo 获取配置信息
// @Summary      获取配置信息
// @Description  根据配置ID获取配置详细信息
// @Tags         配置管理
// @Accept       json
// @Produce      json
// @Param        configID path string true "配置ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /config/{configID} [get]
// @Security     BearerAuth
func (pc ConfigController) GetConfigInfo(c *gin.Context) {
	config, err := pc.ConfigRepository.GetConfigByConfigID(c.Param("configID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	configInfoDto := dto.ToConfigInfoDto(config)
	response.Success(c, gin.H{
		"config": configInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// GetConfigs 获取配置列表
// @Summary      获取配置列表
// @Description  获取所有配置的列表
// @Tags         配置管理
// @Accept       json
// @Produce      json
// @Param        request query vo.ConfigListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /config [get]
// @Security     BearerAuth
func (pc ConfigController) GetConfigs(c *gin.Context) {
	var req vo.ConfigListRequest
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
	config, total, err := pc.ConfigRepository.GetConfigs(&req)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"configs": dto.ToConfigsDto(config), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// CreateConfig 创建配置
// @Summary      创建配置
// @Description  创建一个新的配置
// @Tags         配置管理
// @Accept       json
// @Produce      json
// @Param        request body vo.CreateConfigRequest true "创建配置请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /config [post]
// @Security     BearerAuth
func (pc ConfigController) CreateConfig(c *gin.Context) {
	var req vo.CreateConfigRequest
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

	config := model.Config{
		ProjectID:   req.ProjectID,
		Alias:       req.Alias,
		Title:       req.Title,
		MDContent:   req.MDContent,
		Description: req.Description,
		Type:        req.Type,
		Info:        req.Info,
	}

	err := pc.ConfigRepository.CreateConfig(&config)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgCreateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgCreateSuccess))

}

// UpdateConfigByID 更新配置
// @Summary      更新配置
// @Description  根据配置ID更新配置信息
// @Tags         配置管理
// @Accept       json
// @Produce      json
// @Param        configID path string true "配置ID"
// @Param        request body vo.UpdateConfigRequest true "配置信息"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /config/{configID} [patch]
// @Security     BearerAuth
func (pc ConfigController) UpdateConfigByID(c *gin.Context) {
	var req vo.UpdateConfigRequest
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

	// 根据path中的ConfigID获取配置信息
	oldConfig, err := pc.ConfigRepository.GetConfigByConfigID(c.Param("configID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	oldConfig.Title = req.Title
	oldConfig.Description = req.Description
	oldConfig.Info = req.Info
	oldConfig.ProjectID = req.ProjectID
	oldConfig.MDContent = req.MDContent
	// 更新配置
	err = pc.ConfigRepository.UpdateConfig(&oldConfig)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgUpdateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// BatchDeleteConfigByIds 批量删除配置
// @Summary      批量删除配置
// @Description  根据配置ID列表批量删除配置
// @Tags         配置管理
// @Accept       json
// @Produce      json
// @Param        request body vo.DeleteConfigsRequest true "配置ID列表"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /config [delete]
// @Security     BearerAuth
func (pc ConfigController) BatchDeleteConfigByIds(c *gin.Context) {
	var req vo.DeleteConfigsRequest
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

	// 前端传来的配置ID
	reqConfigIds := strings.Split(req.ConfigIds, ",")
	err := pc.ConfigRepository.BatchDeleteConfigByIds(reqConfigIds)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgDeleteFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))

}
