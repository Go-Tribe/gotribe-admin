// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
	"gotribe-admin/pkg/util"
	"strings"
)

// 返回给前端的内容列表
type PostsDto struct {
	ColumnID    string          `json:"columnID,omitempty"`
	PostID      string          `json:"postID"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	CategoryID  string          `json:"categoryID"`
	ProjectID   string          `json:"projectID"`
	UserID      string          `json:"userID" `
	Author      string          `json:"author" `
	Content     string          `json:"content" `
	HtmlContent string          `json:"htmlContent"`
	Ext         string          `json:"ext"`
	Icon        string          `json:"icon"`
	Tag         string          `json:"tag"`
	Type        uint            `json:"type" `
	IsTop       uint            `json:"isTop" `
	IsPasswd    uint            `json:"isPasswd"`
	Category    *model.Category `json:"category"`
	Tags        []*model.Tag    `json:"tags"`
	Project     *model.Project  `json:"project"`
	CreatedAt   string          `json:"createdAt"`
	Status      uint            `json:"status"`
	Location    string          `json:"location"`
	People      string          `json:"people"`
	Time        string          `json:"time"`
	Images      []string        `json:"images"`
	UnitPrice   float64         `json:"unitPrice"`
	Video       string          `json:"video"`
}

func ToPostInfoDto(post *model.Post) PostsDto {
	if post == nil {
		return PostsDto{}
	}
	var imageList []string
	if len(post.Images) > 0 {
		// 用,分割成数组
		imageList = strings.Split(post.Images, ",")
	}
	return PostsDto{
		ColumnID:    post.ColumnID,
		PostID:      post.PostID,
		Title:       post.Title,
		Description: post.Description,
		CategoryID:  post.CategoryID,
		ProjectID:   post.ProjectID,
		UserID:      post.UserID,
		Author:      post.Author,
		Content:     post.Content,
		HtmlContent: post.HtmlContent,
		Ext:         post.Ext,
		Icon:        post.Icon,
		Tag:         post.Tag,
		Type:        post.Type,
		IsTop:       post.IsTop,
		IsPasswd:    post.IsPasswd,
		Category:    post.Category,
		CreatedAt:   post.CreatedAt.Format(known.TIME_FORMAT),
		Tags:        post.Tags,
		Project:     post.Project,
		Status:      post.Status,
		Location:    post.Location,
		People:      post.People,
		Time:        post.Time,
		Images:      imageList,
		UnitPrice:   util.FenToYuan(int(post.UnitPrice)),
		Video:       post.Video,
	}
}

func ToPostsDto(postList []*model.Post) []PostsDto {
	var posts []PostsDto
	for _, post := range postList {
		if post == nil {
			continue
		}
		postDto := ToPostInfoDto(post)
		posts = append(posts, postDto)
	}

	return posts
}
