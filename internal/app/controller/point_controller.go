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

// GetPoints 获取积分列表
// @Summary      获取积分列表
// @Description  根据查询条件获取积分列表，支持分页
// @Tags         积分管理
// @Accept       json
// @Produce      json
// @Param        userID query string false "用户ID"
// @Param        projectID query string false "项目ID"
// @Param        pageNum query int false "页码"
// @Param        pageSize query int false "每页数量"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /point [get]
// @Security     BearerAuth
func (pc PointController) GetPoints(c *gin.Context) {
	var req vo.PointLogListRequest
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
	point, total, err := pc.PointRepository.GetPointLogs(&req)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"points": dto.ToPointsDto(point), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// CreatePoint 创建积分
// @Summary      创建积分
// @Description  为用户创建积分记录
// @Tags         积分管理
// @Accept       json
// @Produce      json
// @Param        request body vo.CreatePointLogRequest true "积分信息"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /point [post]
// @Security     BearerAuth
func (pc PointController) CreatePoint(c *gin.Context) {
	var req vo.CreatePointLogRequest
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

	err := pc.PointRepository.CreatePoint(req.UserID, "admin", "后台添加", "0", req.ProjectID, req.Point)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgCreateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgCreateSuccess))

}
