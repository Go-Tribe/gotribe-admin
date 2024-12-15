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

type IPointController interface {
	GetPoints(c *gin.Context)   // 获取积分列表
	CreatePoint(c *gin.Context) // 创建积分
}

type PointController struct {
	PointRepository repository.IPointLogRepository
}

// 构造函数
func NewPointController() IPointController {
	pointLogRepository := repository.NewPointLogRepository()
	pointController := PointController{PointRepository: pointLogRepository}
	return pointController
}

// 获取积分列表
func (pc PointController) GetPoints(c *gin.Context) {
	var req vo.PointLogListRequest
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
	point, total, err := pc.PointRepository.GetPointLogs(&req)
	if err != nil {
		response.Fail(c, nil, "获取积分列表失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"points": dto.ToPointsDto(point), "total": total}, "获取积分列表成功")
}

// 创建积分
func (pc PointController) CreatePoint(c *gin.Context) {
	var req vo.CreatePointLogRequest
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

	err := pc.PointRepository.CreatePoint(req.UserID, "admin", "后台添加", "0", req.ProjectID, req.Point)
	if err != nil {
		response.Fail(c, nil, "创建积分失败: "+err.Error())
		return
	}
	response.Success(c, nil, "创建积分成功")

}
