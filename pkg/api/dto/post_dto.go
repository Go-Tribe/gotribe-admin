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
	"time"
)

// 返回给前端的内容列表
type PostsDto struct {
	ID          uint            `json:"id"`
	Slug        string          `json:"slug"`
	ColumnID    uint            `json:"columnID,omitempty"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	CategoryID  uint            `json:"categoryID"`
	ProjectId   uint            `json:"projectId"`
	UserID      uint            `json:"userID"`
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
	ShowTime    string          `json:"showTime"`
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
		ID:          post.ID,
		Slug:        post.Slug,
		ColumnID:    post.ColumnID,
		Title:       post.Title,
		Description: post.Description,
		CategoryID:  post.CategoryID,
		ProjectId:   post.ProjectId,
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
		Time:        formatPostTime(post.Time, known.TIME_FORMAT_SHORT),
		Images:      imageList,
		UnitPrice:   util.MoneyUtil.CentsToYuan(int64(post.UnitPrice)),
		Video:       post.Video,
		ShowTime:    formatPostTime(post.ShowTime, known.TIME_FORMAT),
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

func formatPostTime(value *time.Time, layout string) string {
	if value == nil {
		return ""
	}
	return value.Format(layout)
}
