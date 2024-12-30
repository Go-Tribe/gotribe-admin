// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
	"gotribe-admin/pkg/util/upload"
)

type ResourceDto struct {
	ResourceID    string `json:"resourceID"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	URL           string `json:"url"`
	Path          string `json:"path"`
	FileType      uint   `json:"fileType"`
	FileExtension string `json:"file_extension"`
	Size          int64  `json:"size"`
	CreatedAt     string `json:"createdAt"`
}

func ToResourceInfoDto(resource model.Resource) ResourceDto {
	return ResourceDto{
		ResourceID:    resource.ResourceID,
		Title:         resource.Title,
		Description:   resource.Description,
		URL:           resource.URL,
		Path:          resource.Path,
		FileType:      resource.FileType,
		FileExtension: resource.FileExtension,
		Size:          resource.Size,
		CreatedAt:     resource.CreatedAt.Format(known.TIME_FORMAT),
	}
}

func ToResourcesDto(resourceList []*model.Resource) []ResourceDto {
	var resources []ResourceDto
	for _, resource := range resourceList {
		resourceDto := ResourceDto{
			ResourceID:    resource.ResourceID,
			Title:         resource.Title,
			Description:   resource.Description,
			URL:           resource.URL,
			Path:          resource.Path,
			FileType:      resource.FileType,
			FileExtension: resource.FileExtension,
			Size:          resource.Size,
			CreatedAt:     resource.CreatedAt.Format(known.TIME_FORMAT),
		}

		resources = append(resources, resourceDto)
	}

	return resources
}

type UploadResourceDto struct {
	FileExt  string `json:"fileExt"`
	Key      string `json:"key"`
	Domain   string `json:"domain"`
	FileType int    `json:"fileType"`
}

func ToUploadResourceDto(resource *upload.FileRet) UploadResourceDto {
	return UploadResourceDto{
		FileExt:  resource.FileExt,
		Key:      resource.Key,
		Domain:   "",
		FileType: 0,
	}
}
