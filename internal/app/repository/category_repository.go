// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"
	"errors"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"

	"gorm.io/gorm"
)

type ICategoryRepository interface {
	GetCategoryByID(ctx context.Context, id uint) (model.Category, error)
	GetCategorys(ctx context.Context) ([]*model.Category, error)                     // 获取分类列表
	GetCategoryTree(ctx context.Context) ([]*model.Category, error)                  // 获取分类树
	CreateCategory(ctx context.Context, category *model.Category) error              // 创建分类
	UpdateCategoryByID(ctx context.Context, id uint, category *model.Category) error // 更新分类
	BatchDeleteCategoryByIds(ctx context.Context, ids []uint) error                  // 批量删除分类
}

type CategoryRepository struct {
}

func NewCategoryRepository() ICategoryRepository {
	return CategoryRepository{}
}

// 获取单个分类详情
func (cr CategoryRepository) GetCategoryByID(ctx context.Context, id uint) (model.Category, error) {
	var category model.Category
	err := common.WithContext(ctx).DB().Where("id = ?", id).First(&category).Error
	return category, err
}

// 获取分类列表
func (cr CategoryRepository) GetCategorys(ctx context.Context) ([]*model.Category, error) {
	var categorys []*model.Category
	err := common.WithContext(ctx).DB().Order("sort").Find(&categorys).Error
	return categorys, err
}

// 获取分类树
func (cr CategoryRepository) GetCategoryTree(ctx context.Context) ([]*model.Category, error) {
	var categorys []*model.Category
	err := common.WithContext(ctx).DB().Order("sort").Find(&categorys).Error
	return GenCategoryTree(0, categorys), err
}

func GenCategoryTree(parentID uint, categorys []*model.Category) []*model.Category {
	tree := make([]*model.Category, 0)

	for _, m := range categorys {
		if m.ParentID == parentID {
			children := GenCategoryTree(m.ID, categorys)
			m.Children = children
			tree = append(tree, m)
		}
	}
	return tree
}

// 创建分类
func (cr CategoryRepository) CreateCategory(ctx context.Context, category *model.Category) error {
	err := common.WithContext(ctx).DB().Create(category).Error
	return err
}

// 更新分类
func (cr CategoryRepository) UpdateCategoryByID(ctx context.Context, id uint, category *model.Category) error {
	err := common.WithContext(ctx).DB().Model(category).Where("id = ?", id).Updates(category).Error
	return err
}

// 批量删除分类
func (cr CategoryRepository) BatchDeleteCategoryByIds(ctx context.Context, ids []uint) error {
	var categorys []*model.Category

	err := common.WithContext(ctx).DB().Where("id IN (?)", ids).Find(&categorys).Error
	if err != nil {
		return err
	}
	// 校验分类是否可以删除
	for _, category := range categorys {
		if category.ID == known.DEFAULT_ID {
			return errors.New("默认分类不允许删除")
		}
		if isPID(ctx, int64(category.ID)) {
			return errors.New("该分类下包含子分类，请先删除子分类")
		}
	}

	err = common.WithContext(ctx).DB().Unscoped().Delete(&categorys).Error
	return err
}

// isPID 判断是否为别人的父类 ID
// 存在 true 不存在 false
func isPID(ctx context.Context, ID int64) bool {
	var category model.Category
	if err := common.WithContext(ctx).DB().Where("parent_id = ?", ID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		} else {
			common.Log.Error(err.Error())
			return false
		}
	}
	return true
}
