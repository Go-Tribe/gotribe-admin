package common

import "github.com/gin-gonic/gin"

// MsgKey 定义统一的消息键，集中管理
type MsgKey string

const (
	MsgLoginSuccess        MsgKey = "login_success"
	MsgLogoutSuccess       MsgKey = "logout_success"
	MsgRefreshTokenSuccess MsgKey = "refresh_token_success"
	MsgListSuccess         MsgKey = "list_success"
	MsgListFail            MsgKey = "list_fail"
	MsgGetSuccess          MsgKey = "get_success"
	MsgGetFail             MsgKey = "get_fail"
	MsgCreateSuccess       MsgKey = "create_success"
	MsgCreateFail          MsgKey = "create_fail"
	MsgUpdateSuccess       MsgKey = "update_success"
	MsgUpdateFail          MsgKey = "update_fail"
	MsgDeleteSuccess       MsgKey = "delete_success"
	MsgDeleteFail          MsgKey = "delete_fail"
	MsgJWTAuthFail         MsgKey = "jwt_auth_fail"
	MsgUserSerializeFail   MsgKey = "user_serialize_fail"

	// Admin Repository 相关错误消息
	MsgUserNotFound      MsgKey = "user_not_found"
	MsgUserDisabled      MsgKey = "user_disabled"
	MsgUserRoleDisabled  MsgKey = "user_role_disabled"
	MsgPasswordIncorrect MsgKey = "password_incorrect"
	MsgUserNotLoggedIn   MsgKey = "user_not_logged_in"
	MsgUserNotFoundByID  MsgKey = "user_not_found_by_id"
	MsgNoUsersFound      MsgKey = "no_users_found"
	MsgRoleInfoFailed    MsgKey = "role_info_failed"
	MsgNoUsersWithRole   MsgKey = "no_users_with_role"

	// Role Repository 相关错误消息
	MsgGetRoleApisFailed    MsgKey = "get_role_apis_failed"
	MsgLoadRolePolicyFailed MsgKey = "load_role_policy_failed"
	MsgUpdateRoleApisFailed MsgKey = "update_role_apis_failed"
	MsgDeleteRoleApisFailed MsgKey = "delete_role_apis_failed"

	// API Repository 相关错误消息
	MsgGetApiInfoFailed    MsgKey = "get_api_info_failed"
	MsgUpdateApiFailed     MsgKey = "update_api_failed"
	MsgLoadApiPolicyFailed MsgKey = "load_api_policy_failed"
	MsgGetApiListFailed    MsgKey = "get_api_list_failed"
	MsgNoApiListFound      MsgKey = "no_api_list_found"
	MsgDeleteApiFailed     MsgKey = "delete_api_failed"
)

// i18nMessages 双语消息表（可按需扩展）
var i18nMessages = map[string]map[MsgKey]string{
	"zh": {
		MsgLoginSuccess:        "登录成功",
		MsgLogoutSuccess:       "退出成功",
		MsgRefreshTokenSuccess: "刷新token成功",
		MsgListSuccess:         "获取列表成功",
		MsgListFail:            "获取列表失败",
		MsgGetSuccess:          "获取成功",
		MsgGetFail:             "获取失败",
		MsgCreateSuccess:       "创建成功",
		MsgCreateFail:          "创建失败",
		MsgUpdateSuccess:       "更新成功",
		MsgUpdateFail:          "更新失败",
		MsgDeleteSuccess:       "删除成功",
		MsgDeleteFail:          "删除失败",
		MsgJWTAuthFail:         "JWT认证失败",
		MsgUserSerializeFail:   "用户信息序列化失败",

		// Admin Repository 相关错误消息
		MsgUserNotFound:      "用户不存在",
		MsgUserDisabled:      "用户被禁用",
		MsgUserRoleDisabled:  "用户角色被禁用",
		MsgPasswordIncorrect: "密码错误",
		MsgUserNotLoggedIn:   "用户未登录",
		MsgUserNotFoundByID:  "未获取到ID为%d的用户",
		MsgNoUsersFound:      "未获取到任何用户信息",
		MsgRoleInfoFailed:    "根据角色ID获取角色信息失败",
		MsgNoUsersWithRole:   "根据角色ID未获取到拥有该角色的用户",

		// Role Repository 相关错误消息
		MsgGetRoleApisFailed:    "获取角色的权限接口失败",
		MsgLoadRolePolicyFailed: "角色的权限接口策略加载失败",
		MsgUpdateRoleApisFailed: "更新角色的权限接口失败",
		MsgDeleteRoleApisFailed: "删除角色关联权限接口失败",

		// API Repository 相关错误消息
		MsgGetApiInfoFailed:    "根据接口ID获取接口信息失败",
		MsgUpdateApiFailed:     "更新权限接口失败",
		MsgLoadApiPolicyFailed: "权限接口策略加载失败",
		MsgGetApiListFailed:    "根据接口ID获取接口列表失败",
		MsgNoApiListFound:      "根据接口ID未获取到接口列表",
		MsgDeleteApiFailed:     "删除权限接口失败",
	},
	"en": {
		MsgLoginSuccess:        "Login successful",
		MsgLogoutSuccess:       "Logout successful",
		MsgRefreshTokenSuccess: "Refresh token successful",
		MsgListSuccess:         "Get list successfully",
		MsgListFail:            "Failed to get list",
		MsgGetSuccess:          "Get successfully",
		MsgGetFail:             "Failed to get",
		MsgCreateSuccess:       "Create successful",
		MsgCreateFail:          "Failed to create",
		MsgUpdateSuccess:       "Update successful",
		MsgUpdateFail:          "Failed to update",
		MsgDeleteSuccess:       "Delete successful",
		MsgDeleteFail:          "Failed to delete",
		MsgJWTAuthFail:         "JWT authentication failed",
		MsgUserSerializeFail:   "User information serialization failed",

		// Admin Repository 相关错误消息
		MsgUserNotFound:      "User not found",
		MsgUserDisabled:      "User is disabled",
		MsgUserRoleDisabled:  "User role is disabled",
		MsgPasswordIncorrect: "Incorrect password",
		MsgUserNotLoggedIn:   "User not logged in",
		MsgUserNotFoundByID:  "User not found by ID %d",
		MsgNoUsersFound:      "No users found",
		MsgRoleInfoFailed:    "Failed to get role information by role ID",
		MsgNoUsersWithRole:   "No users found with this role",

		// Role Repository 相关错误消息
		MsgGetRoleApisFailed:    "Failed to get role APIs",
		MsgLoadRolePolicyFailed: "Failed to load role policy",
		MsgUpdateRoleApisFailed: "Failed to update role APIs",
		MsgDeleteRoleApisFailed: "Failed to delete role APIs",

		// API Repository 相关错误消息
		MsgGetApiInfoFailed:    "Failed to get API information by ID",
		MsgUpdateApiFailed:     "Failed to update API",
		MsgLoadApiPolicyFailed: "Failed to load API policy",
		MsgGetApiListFailed:    "Failed to get API list by ID",
		MsgNoApiListFound:      "No API list found by ID",
		MsgDeleteApiFailed:     "Failed to delete API",
	},
}

// getLang 从上下文中获取语言，默认 zh
func getLang(c *gin.Context) string {
	if c == nil {
		return "zh"
	}
	if v, ok := c.Get("lang"); ok {
		if s, sok := v.(string); sok && s == "en" {
			return "en"
		}
	}
	return "zh"
}

// Msg 返回当前请求语言的消息内容
func Msg(c *gin.Context, key MsgKey) string {
	lang := getLang(c)
	if m, ok := i18nMessages[lang]; ok {
		if s, ok2 := m[key]; ok2 {
			return s
		}
	}
	// 回退到中文
	return i18nMessages["zh"][key]
}
