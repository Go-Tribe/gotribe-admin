// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建用户结构体
type CreateUserRequest struct {
	Username  string `form:"username" json:"username" validate:"required,min=2,max=20,alphanum"`
	Nickname  string `form:"nickname" json:"nickname" validate:"required,min=2,max=20"`
	Email     string `form:"email" json:"email"`
	Phone     string `form:"phone" json:"phone"`
	ProjectID string `form:"projectID" json:"projectID" validate:"required,min=2,max=20"`
	Password  string `form:"password" json:"password" validate:"required,min=6,max=20"`
}

// 获取用户列表结构体
type UserListRequest struct {
	UserID   string `form:"userID" json:"userID"`
	Nickname string `form:"nickname" json:"nickname"`
	Username string `form:"username" json:"username"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

// 批量删除用户结构体
type DeleteUsersRequest struct {
	UserIds string `json:"userIds" form:"userIds"`
}

// 更新用户结构体
type UpdateUserRequest struct {
	Nickname string `form:"nickname" json:"nickname"`
	Email    string `form:"email" json:"email"`
	Phone    string `form:"phone" json:"phone"`
	Password string `form:"password" json:"password"`
}
