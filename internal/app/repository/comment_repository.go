// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"
	"fmt"
	"github.com/dengmengmian/ghelper/gconvert"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type ICommentRepository interface {
	GetCommentByID(ctx context.Context, id uint) (model.Comment, error)                           //获取单条评论
	GetComments(ctx context.Context, req *vo.CommentListRequest) ([]*model.Comment, int64, error) // 获取评论列表
	UpdateComment(ctx context.Context, comment *model.Comment) error                              // 更新评论
}

type CommentRepository struct {
}

// CommentRepository构造函数
func NewCommentRepository() ICommentRepository {
	return CommentRepository{}
}

func buildCommentOrder(req *vo.CommentListRequest) string {
	sortByMap := map[string]string{
		"id":         "id",
		"comment":    "content",
		"userID":     "user_id",
		"user_id":    "user_id",
		"status":     "status",
		"createdAt":  "created_at",
		"created_at": "created_at",
	}

	column, ok := sortByMap[strings.TrimSpace(req.SortBy)]
	if !ok {
		return "created_at DESC"
	}

	direction := "ASC"
	if strings.EqualFold(strings.TrimSpace(req.SortOrder), "desc") {
		direction = "DESC"
	}

	return fmt.Sprintf("%s %s", column, direction)
}

func (cr CommentRepository) GetCommentByID(ctx context.Context, id uint) (model.Comment, error) {
	var comment model.Comment
	err := common.WithContext(ctx).DB().Where("id = ?", id).First(&comment).Error
	return comment, err
}

// 获取评论列表
func (cr CommentRepository) GetComments(ctx context.Context, req *vo.CommentListRequest) ([]*model.Comment, int64, error) {
	var list []*model.Comment
	db := common.WithContext(ctx).DB().Model(&model.Comment{})

	objectID := strings.TrimSpace(req.ObjectID)
	if !gconvert.IsEmpty(objectID) {
		db = db.Where("object_id = ?", objectID)
	}
	if !gconvert.IsEmpty(req.ObjectType) {
		db = db.Where("object_type = ?", req.ObjectType)
	}
	if !gconvert.IsEmpty(req.Status) {
		db = db.Where("status = ?", req.Status)
	}
	if req.ProjectId > 0 {
		db = db.Where("project_id = ?", req.ProjectId)
	}
	if !gconvert.IsEmpty(req.Nickname) {
		nicknameIDs := common.WithContext(ctx).DB().
			Model(&model.User{}).
			Select("id").
			Where("nickname like ?", fmt.Sprintf("%%%s%%", req.Nickname))
		db = db.Where("user_id IN (?)", nicknameIDs)
	}

	db = db.Order(buildCommentOrder(req))

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
	list, err = GetCommentOther(ctx, list)
	return list, total, err
}

// 获取评论其他信息
func GetCommentOther(ctx context.Context, comments []*model.Comment) ([]*model.Comment, error) {
	if len(comments) == 0 {
		return comments, nil
	}

	// 收集所有 UserID
	userIDs := make([]uint, 0, len(comments))
	for _, m := range comments {
		if m.UserID > 0 {
			userIDs = append(userIDs, m.UserID)
		}
	}

	if len(userIDs) == 0 {
		return comments, nil
	}

	// 批量查询
	var users []*model.User
	if err := common.WithContext(ctx).DB().Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return comments, err
	}

	// 建立映射
	userMap := make(map[uint]*model.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// 赋值
	for _, m := range comments {
		if user, ok := userMap[m.UserID]; ok {
			m.User = user
		}
	}
	return comments, nil
}

// 更新评论
func (cr CommentRepository) UpdateComment(ctx context.Context, comment *model.Comment) error {
	err := common.WithContext(ctx).DB().Model(comment).Updates(comment).Error
	if err != nil {
		return err
	}
	return err
}
