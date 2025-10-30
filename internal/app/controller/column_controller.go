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

// 获取当前专栏信息
func (pc ColumnController) GetColumnInfo(c *gin.Context) {
	column, err := pc.ColumnRepository.GetColumnByColumnID(c.Param("columnID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	columnInfoDto := dto.ToColumnInfoDto(column)
	response.Success(c, gin.H{
		"column": columnInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// 获取专栏列表
func (pc ColumnController) GetColumns(c *gin.Context) {
	var req vo.ColumnListRequest
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
	column, total, err := pc.ColumnRepository.GetColumns(&req)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"columns": dto.ToColumnsDto(column), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// 创建专栏
func (pc ColumnController) CreateColumn(c *gin.Context) {
	var req vo.CreateColumnRequest
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

	column := model.Column{
		Title:       req.Title,
		Description: req.Description,
		Info:        req.Info,
		Icon:        req.Icon,
		ProjectID:   req.ProjectID,
	}

	err := pc.ColumnRepository.CreateColumn(&column)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgCreateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgCreateSuccess))

}

// 更新专栏
func (pc ColumnController) UpdateColumnByID(c *gin.Context) {
	var req vo.UpdateColumnRequest
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

	// 根据path中的ColumnID获取专栏信息
	oldColumn, err := pc.ColumnRepository.GetColumnByColumnID(c.Param("columnID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	oldColumn.Title = req.Title
	oldColumn.Description = req.Description
	oldColumn.Info = req.Info
	oldColumn.Icon = req.Icon
	// 更新专栏
	err = pc.ColumnRepository.UpdateColumn(&oldColumn)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgUpdateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// 批量删除
func (tc ColumnController) BatchDeleteColumnByIds(c *gin.Context) {
	var req vo.DeleteColumnsRequest
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
	reqColumnIds := strings.Split(req.ColumnIds, ",")
	err := tc.ColumnRepository.BatchDeleteColumnByIds(reqColumnIds)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgDeleteFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))

}
