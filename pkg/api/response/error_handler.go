// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package response

import (
	"gotribe-admin/internal/pkg/common"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// HandleValidationError 处理参数验证错误
func HandleValidationError(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errStr := validationErrors[0].Translate(common.GetTransFromCtx(c))
		ValidationFail(c, errStr)
	} else {
		ValidationFail(c, err.Error())
	}
}

// HandleBindError 处理参数绑定错误
func HandleBindError(c *gin.Context, err error) {
	ValidationFail(c, err.Error())
}

// HandleDatabaseError 处理数据库错误
func HandleDatabaseError(c *gin.Context, err error, operation common.MsgKey) {
	message := common.Msg(c, operation) + ": " + err.Error()
	DatabaseError(c, message)
}

// HandleRepositoryError 处理Repository层错误
func HandleRepositoryError(c *gin.Context, err error, operation common.MsgKey) {
	message := common.Msg(c, operation) + ": " + err.Error()
	InternalServerError(c, message)
}

// HandleNotFoundError 处理资源不存在错误
func HandleNotFoundError(c *gin.Context, resource string) {
	message := resource + "不存在"
	NotFound(c, message)
}

// HandleConflictError 处理资源冲突错误
func HandleConflictError(c *gin.Context, message string) {
	ResponseWithCode(c, CodeConflict, nil, message)
}
