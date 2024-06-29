// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

type ConfigDto struct {
	Alias       string `json:"alias"`
	Type        uint   `json:"type"`
	ConfigID    string `json:"configID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Info        string `json:"info"`
	MDContent   string `json:"mdContent"`
	ProjectID   string `json:"projectID"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func ToConfigInfoDto(config model.Config) ConfigDto {
	return ConfigDto{
		Alias:       config.Alias,
		Type:        config.Type,
		Info:        config.Info,
		ConfigID:    config.ConfigID,
		Title:       config.Title,
		ProjectID:   config.ProjectID,
		Description: config.Description,
		MDContent:   config.MDContent,
		CreatedAt:   config.CreatedAt.Format(known.TimeFormat),
		UpdatedAt:   config.UpdatedAt.Format(known.TimeFormat),
	}
}

func ToConfigsDto(configList []*model.Config) []ConfigDto {
	var configs []ConfigDto
	for _, config := range configList {
		configDto := ConfigDto{
			Alias:       config.Alias,
			Type:        config.Type,
			ConfigID:    config.ConfigID,
			Title:       config.Title,
			Description: config.Description,
			ProjectID:   config.ProjectID,
			MDContent:   config.MDContent,
			CreatedAt:   config.CreatedAt.Format(known.TimeFormat),
			UpdatedAt:   config.UpdatedAt.Format(known.TimeFormat),
		}

		configs = append(configs, configDto)
	}

	return configs
}
