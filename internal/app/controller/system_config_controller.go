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
	"gotribe-admin/pkg/api/dto"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
)

type ISystemConfigController interface {
	GetSystemConfigInfo(c *gin.Context)    // 获取当前系统配置信息
	UpdateSystemConfigByID(c *gin.Context) // 更新系统配置
}

type SystemConfigController struct {
	SystemConfigRepository repository.ISystemConfigRepository
}

// 构造函数
func NewSystemConfigController() ISystemConfigController {
	systemConfigRepository := repository.NewSystemConfigRepository()
	systemConfigController := SystemConfigController{SystemConfigRepository: systemConfigRepository}
	return systemConfigController
}

// 获取当前系统配置信息
func (tc SystemConfigController) GetSystemConfigInfo(c *gin.Context) {
	systemConfig, err := tc.SystemConfigRepository.GetSystemConfig()
	if err != nil {
		response.Fail(c, nil, "获取当前系统配置信息失败: "+err.Error())
		return
	}
	systemConfigInfoDto := dto.ToSystemConfigInfoDto(&systemConfig)
	response.Success(c, gin.H{
		"systemConfig": systemConfigInfoDto,
	}, "获取当前系统配置信息成功")
}

// 更新系统配置
func (tc SystemConfigController) UpdateSystemConfigByID(c *gin.Context) {
	var req vo.CreateSystemConfigRequest
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

	// 根据path中的SystemConfigID获取系统配置信息
	oldSystemConfig, err := tc.SystemConfigRepository.GetSystemConfig()
	if err != nil {
		response.Fail(c, nil, "获取需要更新的系统配置信息失败: "+err.Error())
		return
	}
	oldSystemConfig.Title = req.Title
	oldSystemConfig.Content = req.Content
	oldSystemConfig.Footer = req.Footer
	oldSystemConfig.Icon = req.Icon
	oldSystemConfig.Logo = req.Logo
	// 更新系统配置
	err = tc.SystemConfigRepository.UpdateSystemConfig(&oldSystemConfig)
	if err != nil {
		response.Fail(c, nil, "更新系统配置失败: "+err.Error())
		return
	}
	response.Success(c, nil, "更新系统配置成功")
}
