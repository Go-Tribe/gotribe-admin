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
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

func ToTagInfoDto(tag model.Tag) TagDto {
	return TagDto{
		TagID:       tag.TagID,
		Title:       tag.Title,
		Description: tag.Description,
		CreatedAt:   tag.CreatedAt.Format(known.TIME_FORMAT),
	}
}

func ToTagsDto(tagList []*model.Tag) []TagDto {
	var tags []TagDto
	for _, tag := range tagList {
		tagDto := TagDto{
			TagID:       tag.TagID,
			Title:       tag.Title,
			Description: tag.Description,
			CreatedAt:   tag.CreatedAt.Format(known.TIME_FORMAT),
		}

		tags = append(tags, tagDto)
	}

	return tags
}
