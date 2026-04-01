// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"fmt"
	"gotribe-admin/config"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"gotribe-admin/pkg/util/upload"
)

type IResourceRepository interface {
	CreateResource(resource *model.Resource) error                              // 创建资源
	GetResourceByID(id uint) (model.Resource, error)                            // 获取单个资源
	GetResources(req *vo.ResourceListRequest) ([]*model.Resource, int64, error) // 获取资源列表
	UpdateResource(resource *model.Resource) error                              // 更新资源
	DeleteResourceByID(id uint) error                                           // 删除资源
}

type ResourceRepository struct {
}

// ResourceRepository构造函数
func NewResourceRepository() IResourceRepository {
	return ResourceRepository{}
}

// 获取单个资源
func (rr ResourceRepository) GetResourceByID(id uint) (model.Resource, error) {
	var resource model.Resource
	err := common.DB.Where("id = ?", id).First(&resource).Error
	return resource, err
}

// 获取资源列表
func (rr ResourceRepository) GetResources(req *vo.ResourceListRequest) ([]*model.Resource, int64, error) {
	var list []*model.Resource
	db := common.DB.Model(&model.Resource{}).Order("created_at DESC")

	if int(req.Type) > 0 {
		db = db.Where("file_type = ?", req.Type)
	}

	if req.ID > 0 {
		db = db.Where("id = ?", req.ID)
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
func (rr ResourceRepository) DeleteResourceByID(id uint) error {
	project, err := rr.GetResourceByID(id)
	if err != nil {
		return fmt.Errorf("未获取到ID为%d的项目", id)
	}

	// 硬删除
	err = common.DB.Unscoped().Delete(&project).Error
	if err != nil {
		return err
	}
	// 删除 cdn 文件
	provider := config.Conf.UploadFile.GetUploadProvider(config.Conf.System.EnableOss)
	upload, err := upload.NewService(
		provider,
		config.Conf.UploadFile.Endpoint,
		config.Conf.UploadFile.Accesskey,
		config.Conf.UploadFile.Secretkey,
		config.Conf.UploadFile.Bucket,
	)
	if err != nil {
		return err
	}

	return upload.DeleteFile(project.Path)
}
