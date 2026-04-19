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
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IColumnController interface {
	GetColumnInfo(c *gin.Context)          // 获取当前登录专栏信息
	GetColumns(c *gin.Context)             // 获取专栏列表
	CreateColumn(c *gin.Context)           // 创建专栏
	UpdateColumnByID(c *gin.Context)       // 更新专栏
	BatchDeleteColumnByIds(c *gin.Context) // 批量删除专栏
}

type ColumnController struct {
	ColumnRepository repository.IColumnRepository
}

// 构造函数
func NewColumnController() IColumnController {
	columnRepository := repository.NewColumnRepository()
	columnController := ColumnController{ColumnRepository: columnRepository}
	return columnController
}

// GetColumnInfo 获取当前专栏信息
// @Summary 获取专栏详情
// @Description 根据专栏ID获取专栏详细信息
// @Tags 专栏管理
// @Accept json
// @Produce json
// @Param id path int true "专栏ID"
// @Success 200 {object} response.Response{data=map[string]dto.ColumnDto} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /column/{id} [get]
// @Security BearerAuth
func (pc ColumnController) GetColumnInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationFail(c, "专栏ID格式错误")
		return
	}
	ctx := c.Request.Context()
	column, err := pc.ColumnRepository.GetColumnByID(ctx, uint(id))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	columnInfoDto := dto.ToColumnInfoDto(column)
	response.Success(c, gin.H{
		"column": columnInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// GetColumns 获取专栏列表
// @Summary 获取专栏列表
// @Description 根据查询条件获取专栏列表，支持分页
// @Tags 专栏管理
// @Accept json
// @Produce json
// @Param columnID query string false "专栏ID"
// @Param projectID query string false "项目ID"
// @Param title query string false "专栏标题"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /column [get]
// @Security BearerAuth
func (pc ColumnController) GetColumns(c *gin.Context) {
	var req vo.ColumnListRequest
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

	ctx := c.Request.Context()
	// 获取
	column, total, err := pc.ColumnRepository.GetColumns(ctx, &req)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgListFail)
		return
	}
	response.Success(c, gin.H{"columns": dto.ToColumnsDto(column), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// CreateColumn 创建专栏
// @Summary 创建专栏
// @Description 创建新的专栏
// @Tags 专栏管理
// @Accept json
// @Produce json
// @Param column body vo.CreateColumnRequest true "专栏信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /column [post]
// @Security BearerAuth
func (pc ColumnController) CreateColumn(c *gin.Context) {
	var req vo.CreateColumnRequest
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

	column := model.Column{
		Title:       req.Title,
		Description: req.Description,
		Info:        req.Info,
		Icon:        req.Icon,
		ProjectId:   req.ProjectId,
	}

	ctx := c.Request.Context()
	err := pc.ColumnRepository.CreateColumn(ctx, &column)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgCreateFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgCreateSuccess))

}

// UpdateColumnByID 更新专栏
// @Summary 更新专栏
// @Description 根据专栏ID更新专栏信息
// @Tags 专栏管理
// @Accept json
// @Produce json
// @Param columnID path string true "专栏ID"
// @Param column body vo.UpdateColumnRequest true "专栏信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /column/{id} [patch]
// @Security BearerAuth
func (pc ColumnController) UpdateColumnByID(c *gin.Context) {
	var req vo.UpdateColumnRequest
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationFail(c, "专栏ID格式错误")
		return
	}

	ctx := c.Request.Context()
	// 根据path中的ID获取专栏信息
	oldColumn, err := pc.ColumnRepository.GetColumnByID(ctx, uint(id))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	oldColumn.Title = req.Title
	oldColumn.Description = req.Description
	oldColumn.Info = req.Info
	oldColumn.Icon = req.Icon
	oldColumn.ProjectId = req.ProjectId
	// 更新专栏
	err = pc.ColumnRepository.UpdateColumn(ctx, &oldColumn)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgUpdateFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// BatchDeleteColumnByIds 批量删除专栏
// @Summary 批量删除专栏
// @Description 根据专栏ID列表批量删除专栏
// @Tags 专栏管理
// @Accept json
// @Produce json
// @Param columns body vo.DeleteColumnsRequest true "专栏ID列表"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /column [delete]
// @Security BearerAuth
func (pc ColumnController) BatchDeleteColumnByIds(c *gin.Context) {
	var req vo.DeleteColumnsRequest
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

	ctx := c.Request.Context()
	err := pc.ColumnRepository.BatchDeleteColumnByIds(ctx, req.Ids)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgDeleteFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))

}
