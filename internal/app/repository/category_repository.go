// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"errors"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"

	"gorm.io/gorm"
)

type ICategoryRepository interface {
	GetConfigByCategoryID(categoryID string) (model.Category, error)
	GetCategorys() ([]*model.Category, error)                             // 获取分类列表
	GetCategoryTree() ([]*model.Category, error)                          // 获取分类树
	CreateCategory(category *model.Category) error                        // 创建分类
	UpdateCategoryByID(categoryID string, category *model.Category) error // 更新分类
	BatchDeleteCategoryByIds(categoryIds []string) error                  // 批量删除分类
}

type CategoryRepository struct {
}

func NewCategoryRepository() ICategoryRepository {
	return CategoryRepository{}
}

// 获取单个分类详情
func (cr CategoryRepository) GetConfigByCategoryID(categoryID string) (model.Category, error) {
	var category model.Category
	err := common.DB.Where("category_id = ?", categoryID).First(&category).Error
	return category, err
}

// 获取分类列表
func (cr CategoryRepository) GetCategorys() ([]*model.Category, error) {
	var categorys []*model.Category
	err := common.DB.Order("sort").Find(&categorys).Error
	return categorys, err
}

// 获取分类树
func (cr CategoryRepository) GetCategoryTree() ([]*model.Category, error) {
	var categorys []*model.Category
	err := common.DB.Order("sort").Find(&categorys).Error
	return GenCategoryTree(0, categorys), err
}

func GenCategoryTree(parentID uint, categorys []*model.Category) []*model.Category {
	tree := make([]*model.Category, 0)

	for _, m := range categorys {
		if *m.ParentID == parentID {
			children := GenCategoryTree(m.ID, categorys)
			m.Children = children
			tree = append(tree, m)
		}
	}
	return tree
}

// 创建分类
func (cr CategoryRepository) CreateCategory(category *model.Category) error {
	err := common.DB.Create(category).Error
	return err
}

// 更新分类
func (cr CategoryRepository) UpdateCategoryByID(categoryID string, category *model.Category) error {
	err := common.DB.Model(category).Where("category_id = ?", categoryID).Updates(category).Error
	return err
}

// 批量删除分类
func (cr CategoryRepository) BatchDeleteCategoryByIds(categoryIds []string) error {
	var categorys []*model.Category

	err := common.DB.Where("category_id IN (?)", categoryIds).Find(&categorys).Error
	if err != nil {
		return err
	}
	j := 0
	for _, category := range categorys {
		if category.ID != known.DEFAULT_ID && !isPID(int64((category.ID))) {
			categorys[j] = category
			j++
		}
	}
	// Slice categorys to new size.
	categorys = categorys[:j]
	err = common.DB.Unscoped().Delete(&categorys).Error
	return err
}

// isPID 判断是否为别人的父类 ID
// 存在 true 不存在 false
func isPID(ID int64) bool {
	var category model.Category
	if err := common.DB.Where("parent_id = ?", ID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		} else {
			common.Log.Error(err.Error())
			return false
		}
	}
	return true
}
