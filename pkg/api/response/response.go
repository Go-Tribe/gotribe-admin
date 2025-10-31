// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
// @Description 统一API响应格式
type Response struct {
	Code    int         `json:"code" example:"200"`     // 状态码
	Message string      `json:"message" example:"操作成功"` // 响应消息
	Data    interface{} `json:"data"`                   // 响应数据
}

// 返回前端
func ResponseFunc(c *gin.Context, httpStatus int, code int, data gin.H, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// ResponseWithCode 根据错误码自动设置HTTP状态码
func ResponseWithCode(c *gin.Context, code int, data gin.H, message string) {
	httpStatus := GetHTTPStatus(code)
	ResponseFunc(c, httpStatus, code, data, message)
}

// 返回前端-成功
func Success(c *gin.Context, data gin.H, message string) {
	ResponseFunc(c, http.StatusOK, CodeSuccess, data, message)
}

// 返回前端-失败
func Fail(c *gin.Context, data gin.H, message string) {
	ResponseFunc(c, http.StatusBadRequest, CodeBadRequest, data, message)
}

// 参数验证失败
func ValidationFail(c *gin.Context, message string) {
	ResponseWithCode(c, CodeValidationFailed, nil, message)
}

// 未认证
func Unauthorized(c *gin.Context, message string) {
	ResponseWithCode(c, CodeUnauthorized, nil, message)
}

// 权限不足
func Forbidden(c *gin.Context, message string) {
	ResponseWithCode(c, CodeForbidden, nil, message)
}

// 资源不存在
func NotFound(c *gin.Context, message string) {
	ResponseWithCode(c, CodeNotFound, nil, message)
}

// 内部服务器错误
func InternalServerError(c *gin.Context, message string) {
	ResponseWithCode(c, CodeInternalServerError, nil, message)
}

// 数据库错误
func DatabaseError(c *gin.Context, message string) {
	ResponseWithCode(c, CodeDatabaseError, nil, message)
}

// 用户不存在
func UserNotFound(c *gin.Context, message string) {
	ResponseWithCode(c, CodeUserNotFound, nil, message)
}

// 用户已被禁用
func UserDisabled(c *gin.Context, message string) {
	ResponseWithCode(c, CodeUserDisabled, nil, message)
}

// 密码错误
func PasswordIncorrect(c *gin.Context, message string) {
	ResponseWithCode(c, CodePasswordIncorrect, nil, message)
}

// Token无效
func TokenInvalid(c *gin.Context, message string) {
	ResponseWithCode(c, CodeTokenInvalid, nil, message)
}

// 权限不足（业务级别）
func PermissionDenied(c *gin.Context, message string) {
	ResponseWithCode(c, CodePermissionDenied, nil, message)
}
