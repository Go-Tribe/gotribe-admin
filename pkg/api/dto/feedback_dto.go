// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

type FeedbackDto struct {
	ProjectID string     `json:"projectID"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	UserID    string     `json:"userID"`
	User      UserDto    `json:"user"`
	Phone     string     `json:"phone"`
	Project   ProjectDto `json:"project"`
	CreatedAt string     `json:"createdAt"`
}

func toFeedbackDto(feedBack model.Feedback) FeedbackDto {
	dto := FeedbackDto{
		ProjectID: feedBack.ProjectID,
		Content:   feedBack.Content,
		Title:     feedBack.Title,
		UserID:    feedBack.UserID,
		Phone:     feedBack.Phone,
		CreatedAt: feedBack.CreatedAt.Format(known.TIME_FORMAT),
	}
	if feedBack.User != nil {
		dto.User = ToUserInfoDto(feedBack.User)
	}
	if feedBack.Project != nil {
		dto.Project = ToProjectInfoDto(*feedBack.Project)
	}
	return dto
}

// 将单个 Feedback 转换为 FeedbackDto
func ToFeedbackInfoDto(feedBack model.Feedback) FeedbackDto {
	return toFeedbackDto(feedBack)
}

// 将多个 Feedback 转换为 []FeedbackDto
func ToFeedbacksDto(feedBackList []*model.Feedback) []FeedbackDto {
	var feedBacks = make([]FeedbackDto, len(feedBackList))
	for i, feedBack := range feedBackList {
		feedBacks[i] = toFeedbackDto(*feedBack)
	}
	return feedBacks
}
