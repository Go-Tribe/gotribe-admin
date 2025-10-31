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
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
)

type IOperationLogController interface {
	GetOperationLogs(c *gin.Context)             // 获取操作日志列表
	BatchDeleteOperationLogByIds(c *gin.Context) //批量删除操作日志
}

type OperationLogController struct {
	operationLogRepository repository.IOperationLogRepository
}

func NewOperationLogController() IOperationLogController {
	operationLogRepository := repository.NewOperationLogRepository()
	operationLogController := OperationLogController{operationLogRepository: operationLogRepository}
	return operationLogController
}

// GetOperationLogs 获取操作日志列表
// @Summary      获取操作日志列表
// @Description  获取所有操作日志的列表
// @Tags         操作日志管理
// @Accept       json
// @Produce      json
// @Param        request query vo.OperationLogListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /operation-log [get]
// @Security     BearerAuth
func (oc OperationLogController) GetOperationLogs(c *gin.Context) {
	var req vo.OperationLogListRequest
	// 绑定参数
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
	logs, total, err := oc.operationLogRepository.GetOperationLogs(&req)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"logs": logs, "total": total}, common.Msg(c, common.MsgListSuccess))
}

// BatchDeleteOperationLogByIds 批量删除操作日志
// @Summary      批量删除操作日志
// @Description  根据操作日志ID列表批量删除操作日志
// @Tags         操作日志管理
// @Accept       json
// @Produce      json
// @Param        request body vo.DeleteOperationLogRequest true "操作日志ID列表"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /operation-log [delete]
// @Security     BearerAuth
func (oc OperationLogController) BatchDeleteOperationLogByIds(c *gin.Context) {
	var req vo.DeleteOperationLogRequest
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

	// 删除接口
	err := oc.operationLogRepository.BatchDeleteOperationLogByIds(req.OperationLogIds)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgDeleteFail)+": "+err.Error())
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))
}
