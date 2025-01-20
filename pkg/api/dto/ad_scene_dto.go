// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

type AdSceneDto struct {
	AdSceneID    string `json:"adSceneID"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ProjectID    string `json:"projectID"`
	ProjectTitle string `json:"projectTitle"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// toAdSceneDto converts a model.AdScene to an AdSceneDto.
func toAdSceneDto(adScene model.AdScene) AdSceneDto {
	var projectTitle string
	if adScene.Project != nil {
		projectTitle = adScene.Project.Title
	}
	return AdSceneDto{
		AdSceneID:    adScene.AdSceneID,
		Title:        adScene.Title,
		Description:  adScene.Description,
		ProjectID:    adScene.ProjectID,
		CreatedAt:    adScene.CreatedAt.Format(known.TIME_FORMAT),
		UpdatedAt:    adScene.UpdatedAt.Format(known.TIME_FORMAT),
		ProjectTitle: projectTitle,
	}
}

// ToAdSceneInfoDto converts a model.AdScene to an AdSceneDto.
func ToAdSceneInfoDto(adScene model.AdScene) AdSceneDto {
	return toAdSceneDto(adScene)
}

// ToAdScenesDto converts a list of model.AdScene to a list of AdSceneDto.
func ToAdScenesDto(adSceneList []*model.AdScene) []AdSceneDto {
	var adScenes []AdSceneDto
	for _, adScene := range adSceneList {
		adScenes = append(adScenes, toAdSceneDto(*adScene))
	}
	return adScenes
}
