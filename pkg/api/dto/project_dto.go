// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

type ProjectDto struct {
	ProjectID      string `json:"projectID"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	CreatedAt      string `json:"createdAt"`
	Name           string `json:"name"`
	Keywords       string `json:"keywords"`
	Domain         string `json:"domain"`
	PostURL        string `json:"postUrl"`
	ICP            string `json:"icp"`
	Author         string `json:"author"`
	BaiduAnalytics string `json:"baiduAnalytics"`
	Favicon        string `json:"favicon"`
	PublicSecurity string `json:"publicSecurity"`
	NavImage       string `json:"navImage"`
	Info           string `json:"info"`
	PushToken      string `json:"pushToken"`
}

func ToProjectInfoDto(project model.Project) ProjectDto {
	return ProjectDto{
		ProjectID:      project.ProjectID,
		Title:          project.Title,
		Description:    project.Description,
		CreatedAt:      project.CreatedAt.Format(known.TimeFormat),
		Name:           project.Name,
		Keywords:       project.Keywords,
		Domain:         project.Domain,
		PostURL:        project.PostURL,
		ICP:            project.ICP,
		Author:         project.Author,
		BaiduAnalytics: project.BaiduAnalytics,
		Favicon:        project.Favicon,
		PublicSecurity: project.PublicSecurity,
		NavImage:       project.NavImage,
		Info:           project.Info,
		PushToken:      project.PushToken,
	}
}

func ToProjectsDto(projectList []*model.Project) []ProjectDto {
	var projects []ProjectDto
	for _, project := range projectList {
		projectDto := ProjectDto{
			ProjectID:      project.ProjectID,
			Title:          project.Title,
			Description:    project.Description,
			CreatedAt:      project.CreatedAt.Format(known.TimeFormat),
			Name:           project.Name,
			Keywords:       project.Keywords,
			Domain:         project.Domain,
			PostURL:        project.PostURL,
			ICP:            project.ICP,
			Author:         project.Author,
			BaiduAnalytics: project.BaiduAnalytics,
			Favicon:        project.Favicon,
			PublicSecurity: project.PublicSecurity,
			NavImage:       project.NavImage,
			Info:           project.Info,
			PushToken:      project.PushToken,
		}

		projects = append(projects, projectDto)
	}

	return projects
}
