// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type IUserRepository interface {
	CreateUser(user *model.User) error                              // 创建用户
	GetUserByUserID(userID string) (model.User, error)              // 获取单个用户
	GetUsers(req *vo.UserListRequest) ([]*model.User, int64, error) // 获取用户列表
	UpdateUser(user *model.User) error                              // 更新用户
	BatchDeleteUserByIds(ids []string) error                        // 批量删除用户
	SearchUserByNickname(nickname string) ([]*model.User, error)
}

type UserRepository struct {
}

// UserRepository构造函数
func NewUserRepository() IUserRepository {
	return UserRepository{}
}

// 获取单个用户
func (ur UserRepository) GetUserByUserID(userID string) (model.User, error) {
	var user model.User
	err := common.DB.Where("user_id = ?", userID).First(&user).Error
	return user, err
}

// 获取用户列表
func (ur UserRepository) GetUsers(req *vo.UserListRequest) ([]*model.User, int64, error) {
	var list []*model.User
	db := common.DB.Model(&model.User{}).Order("created_at DESC")

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	userID := strings.TrimSpace(req.UserID)
	if req.UserID != "" {
		db = db.Where("user_id = ?", fmt.Sprintf("%s", userID))
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return GetUserOther(list), total, err
}

func GetUserOther(user []*model.User) []*model.User {
	for _, m := range user {
		userPoint := GetUserPoint(m.UserID)
		m.Point = userPoint
	}
	return user
}

func GetUserPoint(userID string) float64 {
	var sum sql.NullFloat64
	var pointAvailable *model.PointAvailable
	row := common.DB.Model(&pointAvailable).Select("SUM(points)").Where("user_id = ?", userID).Row()
	err := row.Scan(&sum)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果记录不存在，返回 0
			return 0
		}
		return 0
	}
	// 如果 sum 是 NULL，则返回 0
	if !sum.Valid {
		return 0
	}
	return sum.Float64
}

// 创建用户
func (ur UserRepository) CreateUser(user *model.User) error {
	err := common.DB.Create(user).Error
	return err
}

// 更新用户
func (ur UserRepository) UpdateUser(user *model.User) error {
	err := common.DB.Model(user).Updates(user).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (ur UserRepository) BatchDeleteUserByIds(ids []string) error {
	var users []model.User
	for _, id := range ids {
		// 根据ID获取用户
		user, err := ur.GetUserByUserID(id)
		if err != nil {
			return errors.New(fmt.Sprintf("未获取到ID为%s的用户", id))
		}
		users = append(users, user)
	}

	err := common.DB.Delete(&users).Error

	return err
}

// 搜索用户
func (ur UserRepository) SearchUserByNickname(nickname string) ([]*model.User, error) {
	var list []*model.User
	db := common.DB.Model(&model.User{}).Order("created_at DESC")

	if strings.TrimSpace(nickname) != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	err := db.Find(&list).Error
	return list, err
}
