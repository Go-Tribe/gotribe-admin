// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type ITagRepository interface {
	CreateTag(ctx context.Context, tag *model.Tag) (*model.Tag, error)                // 创建标签
	GetTagByID(ctx context.Context, id uint) (model.Tag, error)                       // 获取单个标签
	GetTags(ctx context.Context, req *vo.TagListRequest) ([]*model.Tag, int64, error) // 获取标签列表
	UpdateTag(ctx context.Context, tag *model.Tag) error                              // 更新标签
	BatchDeleteTagByIds(ctx context.Context, ids []uint) error                        // 批量删除
}

type TagRepository struct {
}

// TagRepository构造函数
func NewTagRepository() ITagRepository {
	return TagRepository{}
}

// 获取单个标签
func (tr TagRepository) GetTagByID(ctx context.Context, id uint) (model.Tag, error) {
	var tag model.Tag
	err := common.WithContext(ctx).DB().Where("id = ?", id).First(&tag).Error
	return tag, err
}

// 获取标签列表
func (tr TagRepository) GetTags(ctx context.Context, req *vo.TagListRequest) ([]*model.Tag, int64, error) {
	var list []*model.Tag
	db := common.WithContext(ctx).DB().Model(&model.Tag{}).Order("created_at DESC")

	title := strings.TrimSpace(req.Title)
	if title != "" {
		db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title))
	}
	id := req.ID
	if id > 0 {
		db = db.Where("id = ?", id)
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

// 创建标签
func (tr TagRepository) CreateTag(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	if isTagExist(ctx, tag.Title) {
		return nil, fmt.Errorf("%s标签已存在", tag.Title)
	}
	result := common.WithContext(ctx).DB().Create(tag)
	if result.Error != nil {
		return nil, result.Error
	}
	return tag, nil
}

// 更新标签
func (tr TagRepository) UpdateTag(ctx context.Context, tag *model.Tag) error {
	err := common.WithContext(ctx).DB().Model(tag).Updates(tag).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (tr TagRepository) BatchDeleteTagByIds(ctx context.Context, ids []uint) error {
	var tags []model.Tag
	for _, id := range ids {
		// 根据ID获取标签
		tag, err := tr.GetTagByID(ctx, id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%d的标签", id)
		}
		tags = append(tags, tag)
	}

	err := common.WithContext(ctx).DB().Unscoped().Delete(&tags).Error

	return err
}

func isTagExist(ctx context.Context, title string) bool {
	var tag model.Tag
	result := common.WithContext(ctx).DB().Where("title = ?", title).First(&tag)
	return !errors.Is(result.Error, gorm.ErrRecordNotFound)
}
