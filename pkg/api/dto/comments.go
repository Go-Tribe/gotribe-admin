// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

type CommentDto struct {
	ID          uint   `json:"id"`
	CommentID   string `json:"commentID"`
	ProjectID   string `json:"projectID"`
	Status      uint   `json:"status"`
	UserID      string `json:"userID"`
	ObjectID    string `json:"objectID"`
	ObjectType  uint   `json:"objectType"`
	Content     string `json:"comment"`
	HtmlContent string `json:"htmlContent"`
	Nickname    string `json:"nickname"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// toCommentDto converts a model.Comment to an CommentDto.
func toCommentDto(comment model.Comment) CommentDto {
	var nickname string
	if comment.User != nil {
		nickname = comment.User.Nickname
	}
	return CommentDto{
		ID:          comment.ID,
		CommentID:   comment.CommentID,
		ProjectID:   comment.ProjectID,
		UserID:      comment.UserID,
		ObjectID:    comment.ObjectID,
		ObjectType:  comment.ObjectType,
		Content:     comment.Content,
		HtmlContent: comment.HtmlContent,
		Status:      comment.Status,
		Nickname:    nickname,
		CreatedAt:   comment.CreatedAt.Format(known.TIME_FORMAT),
		UpdatedAt:   comment.UpdatedAt.Format(known.TIME_FORMAT),
	}
}

// ToCommentInfoDto converts a model.Comment to an CommentDto.
func ToCommentInfoDto(comment model.Comment) CommentDto {
	return toCommentDto(comment)
}

// ToCommentsDto converts a list of model.Comment to a list of CommentDto.
func ToCommentsDto(commentList []*model.Comment) []CommentDto {
	var comments []CommentDto
	for _, comment := range commentList {
		comments = append(comments, toCommentDto(*comment))
	}
	return comments
}
