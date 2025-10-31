# 国际化错误处理示例

本文档展示如何在Controller层处理Repository返回的国际化错误。

## 在Controller中处理Repository错误

### 方法1：直接处理RepositoryError

```go
func (ac *AuthController) Login(c *gin.Context) {
    var u model.Admin
    if err := c.ShouldBindJSON(&u); err != nil {
        response.ParamError(c, err.Error())
        return
    }

    admin, err := ac.AdminRepository.Login(u)
    if err != nil {
        // 检查是否为RepositoryError
        if repoErr, ok := err.(*common.RepositoryError); ok {
            // 获取本地化错误消息
            localizedMsg := repoErr.GetLocalizedMessage(c)
            
            // 根据错误类型返回不同的HTTP状态码
            switch repoErr.Code {
            case "user_not_found", "password_incorrect":
                response.PasswordIncorrect(c, localizedMsg)
            case "user_disabled", "user_role_disabled":
                response.Forbidden(c, localizedMsg)
            case "user_not_logged_in":
                response.Unauthorized(c, localizedMsg)
            default:
                response.InternalServerError(c, localizedMsg)
            }
            return
        }
        
        // 处理其他类型的错误
        response.InternalServerError(c, err.Error())
        return
    }

    // 成功处理...
}
```

### 方法2：创建通用错误处理函数

```go
// 在controller包中创建通用错误处理函数
func HandleRepositoryError(c *gin.Context, err error) bool {
    if err == nil {
        return false
    }
    
    // 检查是否为RepositoryError
    if repoErr, ok := err.(*common.RepositoryError); ok {
        localizedMsg := repoErr.GetLocalizedMessage(c)
        
        switch repoErr.Code {
        case "user_not_found", "password_incorrect":
            response.PasswordIncorrect(c, localizedMsg)
        case "user_disabled", "user_role_disabled":
            response.Forbidden(c, localizedMsg)
        case "user_not_logged_in":
            response.Unauthorized(c, localizedMsg)
        case "no_users_found", "role_info_failed", "no_users_with_role":
            response.NotFound(c, localizedMsg)
        default:
            response.InternalServerError(c, localizedMsg)
        }
        return true
    }
    
    // 处理其他类型的错误
    response.InternalServerError(c, err.Error())
    return true
}

// 在Controller中使用
func (ac *AuthController) Login(c *gin.Context) {
    var u model.Admin
    if err := c.ShouldBindJSON(&u); err != nil {
        response.ParamError(c, err.Error())
        return
    }

    admin, err := ac.AdminRepository.Login(u)
    if HandleRepositoryError(c, err) {
        return
    }

    // 成功处理...
}
```

## 支持的错误类型

| 错误代码 | 中文消息 | 英文消息 | 建议HTTP状态码 |
|---------|---------|---------|---------------|
| `user_not_found` | 用户不存在 | User not found | 401 |
| `user_disabled` | 用户被禁用 | User is disabled | 403 |
| `user_role_disabled` | 用户角色被禁用 | User role is disabled | 403 |
| `password_incorrect` | 密码错误 | Incorrect password | 401 |
| `user_not_logged_in` | 用户未登录 | User not logged in | 401 |
| `user_not_found_by_id` | 未获取到ID为%d的用户 | User not found by ID %d | 404 |
| `no_users_found` | 未获取到任何用户信息 | No user information found | 404 |
| `role_info_failed` | 根据角色ID获取角色信息失败 | Failed to get role information by role ID | 500 |
| `no_users_with_role` | 根据角色ID未获取到拥有该角色的用户 | No users with this role found by role ID | 404 |

## 语言检测

系统通过以下方式检测用户语言偏好：

1. 检查Gin Context中的`lang`字段
2. 如果未设置，默认使用中文(`zh`)
3. 支持的语言：`zh`(中文)、`en`(英文)

## 性能考虑

- 错误消息本地化性能：约20ns/op
- 错误对象创建性能：约0.2ns/op
- 建议在高频调用场景中缓存错误对象