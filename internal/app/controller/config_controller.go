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
	"gotribe-admin/pkg/api/dto"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"strings"
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

// 获取当前配置信息
func (pc ConfigController) GetConfigInfo(c *gin.Context) {
	config, err := pc.ConfigRepository.GetConfigByConfigID(c.Param("configID"))
	if err != nil {
		response.Fail(c, nil, "获取当前配置信息失败: "+err.Error())
		return
	}
	configInfoDto := dto.ToConfigInfoDto(config)
	response.Success(c, gin.H{
		"config": configInfoDto,
	}, "获取当前配置信息成功")
}

// 获取配置列表
func (pc ConfigController) GetConfigs(c *gin.Context) {
	var req vo.ConfigListRequest
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

	// 获取
	config, total, err := pc.ConfigRepository.GetConfigs(&req)
	if err != nil {
		response.Fail(c, nil, "获取配置列表失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"configs": dto.ToConfigsDto(config), "total": total}, "获取配置列表成功")
}

// 创建配置
func (pc ConfigController) CreateConfig(c *gin.Context) {
	var req vo.CreateConfigRequest
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
		response.Fail(c, nil, "创建配置失败: "+err.Error())
		return
	}
	response.Success(c, nil, "创建配置成功")

}

// 更新配置
func (pc ConfigController) UpdateConfigByID(c *gin.Context) {
	var req vo.UpdateConfigRequest
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

	// 根据path中的ConfigID获取配置信息
	oldConfig, err := pc.ConfigRepository.GetConfigByConfigID(c.Param("configID"))
	if err != nil {
		response.Fail(c, nil, "获取需要更新的配置信息失败: "+err.Error())
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
		response.Fail(c, nil, "更新配置失败: "+err.Error())
		return
	}
	response.Success(c, nil, "更新配置成功")
}

// 批量删除配置
func (pc ConfigController) BatchDeleteConfigByIds(c *gin.Context) {
	var req vo.DeleteConfigsRequest
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

	// 前端传来的配置ID
	reqConfigIds := strings.Split(req.ConfigIds, ",")
	err := pc.ConfigRepository.BatchDeleteConfigByIds(reqConfigIds)
	if err != nil {
		response.Fail(c, nil, "删除配置失败: "+err.Error())
		return
	}

	response.Success(c, nil, "删除配置成功")

}
