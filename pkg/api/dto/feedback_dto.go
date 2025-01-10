// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

type FeedBackDto struct {
	ProjectID string  `json:"projectID"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	UserID    string  `json:"userID"`
	User      UserDto `json:"user"`
	Phone     string  `json:"phone"`
	CreatedAt string  `json:"createdAt"`
}

func toFeedBackDto(feedBack model.Feedback) FeedBackDto {
	dto := FeedBackDto{
		ProjectID: feedBack.ProjectID,
		Content:   feedBack.Content,
		Title:     feedBack.Title,
		UserID:    feedBack.UserID,
		Phone:     feedBack.Phone,
		CreatedAt: feedBack.CreatedAt.Format(known.TIME_FORMAT),
	}
	//if feedBack.User != nil {
	//	dto.User = ToUserInfoDto(feedBack.User)
	//}
	return dto
}

// 将单个 Feedback 转换为 FeedBackDto
func ToFeedBackInfoDto(feedBack model.Feedback) FeedBackDto {
	return toFeedBackDto(feedBack)
}

// 将多个 Feedback 转换为 []FeedBackDto
func ToFeedBacksDto(feedBackList []*model.Feedback) []FeedBackDto {
	var feedBacks = make([]FeedBackDto, len(feedBackList))
	for i, feedBack := range feedBackList {
		feedBacks[i] = toFeedBackDto(*feedBack)
	}
	return feedBacks
}
