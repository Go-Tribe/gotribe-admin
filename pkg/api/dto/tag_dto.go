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
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Sort        uint   `json:"sort"`
	Count       uint   `json:"count"`
	Status      uint8  `json:"status"`
	CreatedAt   string `json:"createdAt"`
}

func toTagDto(tag *model.Tag) TagDto {
	if tag == nil {
		return TagDto{}
	}
	return TagDto{
		ID:          tag.ID,
		Title:       tag.Title,
		Slug:        tag.Slug,
		Color:       tag.Color,
		Description: tag.Description,
		Sort:        tag.Sort,
		Count:       tag.Count,
		Status:      tag.Status,
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
