// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package response

// 错误码定义
const (
	// 成功
	CodeSuccess = 200

	// 客户端错误 4xx
	CodeBadRequest       = 400 // 请求参数错误
	CodeUnauthorized     = 401 // 未认证
	CodeForbidden        = 403 // 权限不足
	CodeNotFound         = 404 // 资源不存在
	CodeMethodNotAllowed = 405 // 方法不允许
	CodeConflict         = 409 // 资源冲突
	CodeValidationFailed = 422 // 参数验证失败
	CodeTooManyRequests  = 429 // 请求过于频繁

	// 服务器错误 5xx
	CodeInternalServerError = 500 // 内部服务器错误
	CodeBadGateway          = 502 // 网关错误
	CodeServiceUnavailable  = 503 // 服务不可用
	CodeGatewayTimeout      = 504 // 网关超时

	// 业务错误码 1xxx
	CodeUserNotFound      = 1001 // 用户不存在
	CodeUserDisabled      = 1002 // 用户已被禁用
	CodePasswordIncorrect = 1003 // 密码错误
	CodeTokenInvalid      = 1004 // Token无效
	CodeTokenExpired      = 1005 // Token已过期
	CodePermissionDenied  = 1006 // 权限不足
	CodeResourceNotFound  = 1007 // 资源不存在
	CodeResourceConflict  = 1008 // 资源冲突
	CodeDatabaseError     = 1009 // 数据库错误
	CodeExternalAPIError  = 1010 // 外部API错误
)

// 错误码对应的HTTP状态码映射
var CodeToHTTPStatus = map[int]int{
	CodeSuccess:             200,
	CodeBadRequest:          400,
	CodeUnauthorized:        401,
	CodeForbidden:           403,
	CodeNotFound:            404,
	CodeMethodNotAllowed:    405,
	CodeConflict:            409,
	CodeValidationFailed:    422,
	CodeTooManyRequests:     429,
	CodeInternalServerError: 500,
	CodeBadGateway:          502,
	CodeServiceUnavailable:  503,
	CodeGatewayTimeout:      504,

	// 业务错误码映射到HTTP状态码
	CodeUserNotFound:      404,
	CodeUserDisabled:      403,
	CodePasswordIncorrect: 401,
	CodeTokenInvalid:      401,
	CodeTokenExpired:      401,
	CodePermissionDenied:  403,
	CodeResourceNotFound:  404,
	CodeResourceConflict:  409,
	CodeDatabaseError:     500,
	CodeExternalAPIError:  502,
}

// GetHTTPStatus 根据业务错误码获取对应的HTTP状态码
func GetHTTPStatus(code int) int {
	if status, exists := CodeToHTTPStatus[code]; exists {
		return status
	}
	return 500 // 默认返回500
}
