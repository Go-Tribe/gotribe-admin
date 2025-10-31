// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/pkg/api/dto"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

// GetSystemConfigInfo 获取系统配置信息
// @Summary      获取系统配置信息
// @Description  获取当前系统配置信息
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /base/config [get]
func (tc SystemConfigController) GetSystemConfigInfo(c *gin.Context) {
	systemConfig, err := tc.SystemConfigRepository.GetSystemConfig()
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	systemConfigInfoDto := dto.ToSystemConfigInfoDto(&systemConfig)
	response.Success(c, gin.H{
		"systemConfig": systemConfigInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// UpdateSystemConfigByID 更新系统配置
// @Summary      更新系统配置
// @Description  更新系统配置信息
// @Tags         系统配置
// @Accept       json
// @Produce      json
// @Param        request body vo.CreateSystemConfigRequest true "更新系统配置请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /systemConfig/update [patch]
// @Security     BearerAuth
func (tc SystemConfigController) UpdateSystemConfigByID(c *gin.Context) {
	var req vo.CreateSystemConfigRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.ValidationFail(c, errStr)
		return
	}

	// 根据path中的SystemConfigID获取系统配置信息
	oldSystemConfig, err := tc.SystemConfigRepository.GetSystemConfig()
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
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
		response.HandleDatabaseError(c, err, common.MsgUpdateFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}
