// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

type ColumnDto struct {
	ColumnID    string `json:"columnID"`
	Title       string `json:"title"`
	Info        string `json:"info"`
	ProjectID   string `json:"projectID"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

func ToColumnInfoDto(column model.Column) ColumnDto {
	return ColumnDto{
		ColumnID:    column.ColumnID,
		Title:       column.Title,
		Description: column.Description,
		Info:        column.Info,
		Icon:        column.Icon,
		ProjectID:   column.ProjectID,
		CreatedAt:   column.CreatedAt.Format(known.TIME_FORMAT),
	}
}

func ToColumnsDto(columnList []*model.Column) []ColumnDto {
	var columns []ColumnDto
	for _, column := range columnList {
		columnDto := ColumnDto{
			ColumnID:    column.ColumnID,
			Title:       column.Title,
			ProjectID:   column.ProjectID,
			Description: column.Description,
			Info:        column.Info,
			Icon:        column.Icon,
			CreatedAt:   column.CreatedAt.Format(known.TIME_FORMAT),
		}

		columns = append(columns, columnDto)
	}

	return columns
}
