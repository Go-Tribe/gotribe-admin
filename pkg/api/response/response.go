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

// 返回前端-成功
func Success(c *gin.Context, data gin.H, message string) {
	ResponseFunc(c, http.StatusOK, 200, data, message)
}

// 返回前端-失败
func Fail(c *gin.Context, data gin.H, message string) {
	ResponseFunc(c, http.StatusBadRequest, 400, data, message)
}
