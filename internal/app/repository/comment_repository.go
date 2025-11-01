// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"fmt"
	"github.com/dengmengmian/ghelper/gconvert"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type ICommentRepository interface {
	GetCommentByComentID(commentID string) (model.Comment, error)            //获取单条评论
	GetComments(req *vo.CommentListRequest) ([]*model.Comment, int64, error) // 获取评论列表
	UpdateComment(comment *model.Comment) error                              // 更新评论
}

type CommentRepository struct {
}

// CommentRepository构造函数
func NewCommentRepository() ICommentRepository {
	return CommentRepository{}
}

func (cr CommentRepository) GetCommentByComentID(commentID string) (model.Comment, error) {
	var comment model.Comment
	err := common.DB.Where("comment_id = ?", commentID).First(&comment).Error
	return comment, err
}

// 获取评论列表
func (cr CommentRepository) GetComments(req *vo.CommentListRequest) ([]*model.Comment, int64, error) {
	var list []*model.Comment
	db := common.DB.Model(&model.Comment{}).Order("created_at DESC")

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
	if !gconvert.IsEmpty(req.ProjectID) {
		db = db.Where("project_id = ?", req.ProjectID)
	}
	if !gconvert.IsEmpty(req.Nickname) {
		// 查出用户 ID。再用用户 ID 去筛选
		var user model.User
		if result := common.DB.Model(&model.User{}).Where("nickname like ?", fmt.Sprintf("%%%s%%", req.Nickname)).First(&user); result.Error != nil {
			return nil, 0, result.Error
		}
		db = db.Where("user_id = ?", user.UserID)
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
	return GetCommentOther(list), total, err
}

// 获取评论其他信息
func GetCommentOther(comments []*model.Comment) []*model.Comment {
	for _, m := range comments {
		var user *model.User
		_ = common.DB.Where("user_id = ?", m.UserID).First(&user).Error
		m.User = user
	}
	return comments
}

// 更新评论
func (cr CommentRepository) UpdateComment(comment *model.Comment) error {
	err := common.DB.Model(comment).Updates(comment).Error
	if err != nil {
		return err
	}
	return err
}
