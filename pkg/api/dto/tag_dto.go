// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

type TagDto struct {
	TagID       string `json:"tagID"`
	Title       string `json:"title"`
	Color       string `json:"color"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

func toTagDto(tag *model.Tag) TagDto {
	if tag == nil {
		return TagDto{}
	}
	return TagDto{
		TagID:       tag.TagID,
		Title:       tag.Title,
		Color:       tag.Color,
		Description: tag.Description,
		CreatedAt:   tag.CreatedAt.Format(known.TIME_FORMAT),
	}
}

func ToTagInfoDto(tag model.Tag) TagDto {
	return toTagDto(&tag)
}

func ToTagsDto(tagList []*model.Tag) []TagDto {
	if tagList == nil {
		return []TagDto{}
	}

	tags := make([]TagDto, 0, len(tagList))
	for _, tag := range tagList {
		tags = append(tags, toTagDto(tag))
	}

	return tags
}
