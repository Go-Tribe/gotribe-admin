// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

type UserDto struct {
	UserID    string `json:"userID"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	CreatedAt string `json:"createdAt"`
}

func ToUserInfoDto(user model.User) UserDto {
	return UserDto{
		UserID:    user.UserID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		CreatedAt: user.CreatedAt.Format(known.TimeFormat),
	}
}

func ToUsersDto(userList []*model.User) []UserDto {
	var users []UserDto
	for _, user := range userList {
		userDto := UserDto{
			UserID:    user.UserID,
			Username:  user.Username,
			Nickname:  user.Nickname,
			CreatedAt: user.CreatedAt.Format(known.TimeFormat),
		}

		users = append(users, userDto)
	}

	return users
}
