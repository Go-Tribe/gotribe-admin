// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package common

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// RepositoryError 定义Repository层的错误类型
type RepositoryError struct {
	Code    string        // 错误代码，对应MsgKey
	Message string        // 默认错误消息（用于日志等）
	Args    []interface{} // 格式化参数
}

// Error 实现error接口
func (e *RepositoryError) Error() string {
	if len(e.Args) > 0 {
		return fmt.Sprintf(e.Message, e.Args...)
	}
	return e.Message
}

// GetLocalizedMessage 获取本地化的错误消息
func (e *RepositoryError) GetLocalizedMessage(c *gin.Context) string {
	msgKey := MsgKey(e.Code)
	localizedMsg := Msg(c, msgKey)

	// 如果有格式化参数，进行格式化
	if len(e.Args) > 0 {
		return fmt.Sprintf(localizedMsg, e.Args...)
	}
	return localizedMsg
}

// NewRepositoryError 创建Repository错误
func NewRepositoryError(code string, message string, args ...interface{}) *RepositoryError {
	return &RepositoryError{
		Code:    code,
		Message: message,
		Args:    args,
	}
}

// 预定义的Repository错误
var (
	// Admin Repository 相关错误
	ErrUserNotFound      = NewRepositoryError("user_not_found", "用户不存在")
	ErrUserDisabled      = NewRepositoryError("user_disabled", "用户被禁用")
	ErrUserRoleDisabled  = NewRepositoryError("user_role_disabled", "用户角色被禁用")
	ErrPasswordIncorrect = NewRepositoryError("password_incorrect", "密码错误")
	ErrUserNotLoggedIn   = NewRepositoryError("user_not_logged_in", "用户未登录")
	ErrNoUsersFound      = NewRepositoryError("no_users_found", "未获取到任何用户信息")
	ErrRoleInfoFailed    = NewRepositoryError("role_info_failed", "根据角色ID获取角色信息失败")
	ErrNoUsersWithRole   = NewRepositoryError("no_users_with_role", "根据角色ID未获取到拥有该角色的用户")

	// Role Repository 相关错误
	ErrGetRoleApisFailed    = NewRepositoryError("get_role_apis_failed", "获取角色的权限接口失败")
	ErrLoadRolePolicyFailed = NewRepositoryError("load_role_policy_failed", "角色的权限接口策略加载失败")
	ErrUpdateRoleApisFailed = NewRepositoryError("update_role_apis_failed", "更新角色的权限接口失败")
	ErrDeleteRoleApisFailed = NewRepositoryError("delete_role_apis_failed", "删除角色关联权限接口失败")

	// API Repository 相关错误
	ErrGetApiInfoFailed    = NewRepositoryError("get_api_info_failed", "根据接口ID获取接口信息失败")
	ErrUpdateApiFailed     = NewRepositoryError("update_api_failed", "更新权限接口失败")
	ErrLoadApiPolicyFailed = NewRepositoryError("load_api_policy_failed", "权限接口策略加载失败")
	ErrGetApiListFailed    = NewRepositoryError("get_api_list_failed", "根据接口ID获取接口列表失败")
	ErrNoApiListFound      = NewRepositoryError("no_api_list_found", "根据接口ID未获取到接口列表")
	ErrDeleteApiFailed     = NewRepositoryError("delete_api_failed", "删除权限接口失败")
)

// NewUserNotFoundByIDError 创建指定ID用户不存在的错误
func NewUserNotFoundByIDError(id uint) *RepositoryError {
	return NewRepositoryError("user_not_found_by_id", "未获取到ID为%d的用户", id)
}
