// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type IFeedbackRepository interface {
	GetFeedbacks(req *vo.FeedbackListRequest) ([]*model.Feedback, int64, error) // 获取标签列表
}

type FeedbackRepository struct {
}

// FeedbackRepository构造函数
func NewFeedbackRepository() IFeedbackRepository {
	return FeedbackRepository{}
}

// 获取标签列表
func (tr FeedbackRepository) GetFeedbacks(req *vo.FeedbackListRequest) ([]*model.Feedback, int64, error) {
	var list []*model.Feedback
	db := common.DB.Model(&model.Feedback{}).Order("created_at DESC")

	projectID := strings.TrimSpace(req.ProjectID)
	if req.ProjectID != "" {
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
	return list, total, err
}
