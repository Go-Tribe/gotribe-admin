// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"errors"
	"fmt"
	"gotribe-admin/config"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"gotribe-admin/pkg/util/upload"
	"strings"
)

type IResourceRepository interface {
	CreateResource(resource *model.Resource) error                              // 创建资源
	GetResourceByResourceID(resourceID string) (model.Resource, error)          // 获取单个资源
	GetResources(req *vo.ResourceListRequest) ([]*model.Resource, int64, error) // 获取资源列表
	UpdateResource(resource *model.Resource) error                              // 更新资源
	DeleteResourceByID(id string) error                                         // 删除资源
}

type ResourceRepository struct {
}

// ResourceRepository构造函数
func NewResourceRepository() IResourceRepository {
	return ResourceRepository{}
}

// 获取单个资源
func (rr ResourceRepository) GetResourceByResourceID(resourceID string) (model.Resource, error) {
	var resource model.Resource
	err := common.DB.Where("resource_id = ?", resourceID).First(&resource).Error
	return resource, err
}

// 获取资源列表
func (rr ResourceRepository) GetResources(req *vo.ResourceListRequest) ([]*model.Resource, int64, error) {
	var list []*model.Resource
	db := common.DB.Model(&model.Resource{}).Order("created_at DESC")

	if int(req.Type) > 0 {
		db = db.Where("file_type = ?", req.Type)
	}

	resourceID := strings.TrimSpace(req.ResourceID)
	if req.ResourceID != "" {
		db = db.Where("resource_id = ?", fmt.Sprintf("%s", resourceID))
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

// 创建资源
func (rr ResourceRepository) CreateResource(resource *model.Resource) error {
	err := common.DB.Create(resource).Error
	return err
}

// 更新资源
func (rr ResourceRepository) UpdateResource(resource *model.Resource) error {
	err := common.DB.Model(resource).Updates(resource).Error
	if err != nil {
		return err
	}

	return err
}

// 删除文件
func (rr ResourceRepository) DeleteResourceByID(id string) error {
	project, err := rr.GetResourceByResourceID(id)
	if err != nil {
		return errors.New(fmt.Sprintf("未获取到ID为%s的项目", id))
	}

	// 硬删除
	err = common.DB.Unscoped().Delete(&project).Error
	// 删除 cdn 文件
	upload, err := upload.NewUploadFile(
		config.Conf.UploadFile.Endpoint,
		config.Conf.UploadFile.Accesskey,
		config.Conf.UploadFile.Secretkey,
		config.Conf.UploadFile.Bucket,
		config.Conf.System.EnableOss,
	)

	//qiniuUpload := upload.NewQiniu(config.Conf.QiniuConfig.Accesskey, config.Conf.QiniuConfig.Secretkey, config.Conf.QiniuConfig.Bucket)
	//qiniuUpload.DeletdFile(project.Path)
	return upload.DeleteFile(project.Path)
}
