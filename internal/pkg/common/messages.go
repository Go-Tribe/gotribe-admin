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
