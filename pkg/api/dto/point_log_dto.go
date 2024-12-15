// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

type PointDto struct {
	ID        int64   `json:"id"`
	Point     float32 `json:"point"`
	UserID    string  `json:"userID"`
	Reason    string  `json:"reason"`
	Nickname  string  `json:"nickname"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

// toPointDto converts a model.Point to an PointDto.
func toPointDto(point model.PointLog) PointDto {
	var nickname string
	if point.User != nil {
		nickname = point.User.Nickname
	}
	return PointDto{
		ID:        int64(point.ID),
		Point:     point.Points,
		UserID:    point.UserID,
		Nickname:  nickname,
		Reason:    point.Reason,
		CreatedAt: point.CreatedAt.Format(known.TimeFormat),
		UpdatedAt: point.UpdatedAt.Format(known.TimeFormat),
	}
}

// ToPointInfoDto converts a model.Point to an PointDto.
func ToPointInfoDto(point model.PointLog) PointDto {
	return toPointDto(point)
}

// ToPointsDto converts a list of model.Point to a list of PointDto.
func ToPointsDto(pointList []*model.PointLog) []PointDto {
	var points []PointDto
	for _, point := range pointList {
		points = append(points, toPointDto(*point))
	}
	return points
}
