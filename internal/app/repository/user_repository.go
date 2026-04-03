// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"strings"

	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error                              // 创建用户
	GetUserByID(ctx context.Context, id uint) (model.User, error)                        // 获取单个用户
	GetUsers(ctx context.Context, req *vo.UserListRequest) ([]*model.User, int64, error) // 获取用户列表
	UpdateUser(ctx context.Context, user *model.User) error                              // 更新用户
	BatchDeleteUserByIds(ctx context.Context, ids []uint) error                          // 批量删除用户
	SearchUserByNickname(ctx context.Context, nickname string) ([]*model.User, error)
}

type UserRepository struct {
}

// UserRepository构造函数
func NewUserRepository() IUserRepository {
	return UserRepository{}
}

// 获取单个用户
func (ur UserRepository) GetUserByID(ctx context.Context, id uint) (model.User, error) {
	var user model.User
	err := common.WithContext(ctx).DB().Where("id = ?", id).First(&user).Error
	return user, err
}

// 获取用户列表
func (ur UserRepository) GetUsers(ctx context.Context, req *vo.UserListRequest) ([]*model.User, int64, error) {
	var list []*model.User
	db := common.WithContext(ctx).DB().Model(&model.User{}).Order("created_at DESC")

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	if req.UserID > 0 {
		db = db.Where("id = ?", req.UserID)
	}
	projectID := strings.TrimSpace(req.ProjectID)
	if projectID != "" {
		db = db.Where("project_id = ?", projectID)
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
		userPoint := GetUserPoint(m.ID)
		m.Point = userPoint
	}
	return user
}

func GetUserPoint(userID uint) float64 {
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
func (ur UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	normalizeUserContactFields(user)
	err := common.WithContext(ctx).DB().Create(user).Error
	return err
}

// 更新用户
func (ur UserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	normalizeUserContactFields(user)
	err := common.WithContext(ctx).DB().Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"username":   user.Username,
		"project_id": user.ProjectID,
		"password":   user.Password,
		"nickname":   user.Nickname,
		"email":      user.Email,
		"phone":      user.Phone,
		"sex":        user.Sex,
		"status":     user.Status,
		"birthday":   user.Birthday,
		"background": user.Background,
		"ext":        user.Ext,
		"avatar_url": user.AvatarURL,
	}).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (ur UserRepository) BatchDeleteUserByIds(ctx context.Context, ids []uint) error {
	var users []model.User
	for _, id := range ids {
		// 根据ID获取用户
		user, err := ur.GetUserByID(ctx, id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%d的用户", id)
		}
		users = append(users, user)
	}

	err := common.WithContext(ctx).DB().Delete(&users).Error

	return err
}

// 搜索用户
func (ur UserRepository) SearchUserByNickname(ctx context.Context, nickname string) ([]*model.User, error) {
	var list []*model.User
	db := common.WithContext(ctx).DB().Model(&model.User{}).Order("created_at DESC")

	if strings.TrimSpace(nickname) != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	err := db.Find(&list).Error
	return list, err
}

func normalizeUserContactFields(user *model.User) {
	if user == nil {
		return
	}
	user.Email = trimOptionalString(user.Email)
	user.Phone = trimOptionalString(user.Phone)
}

func trimOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
