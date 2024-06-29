// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"errors"
	"fmt"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type IColumnRepository interface {
	CreateColumn(column *model.Column) error                              // 创建专栏
	GetColumnByColumnID(columnID string) (model.Column, error)            // 获取单个专栏
	GetColumns(req *vo.ColumnListRequest) ([]*model.Column, int64, error) // 获取专栏列表
	UpdateColumn(column *model.Column) error                              // 更新专栏
	BatchDeleteColumnByIds(ids []string) error                            // 批量删除专栏
}

type ColumnRepository struct {
}

// ColumnRepository构造函数
func NewColumnRepository() IColumnRepository {
	return ColumnRepository{}
}

// 获取单个专栏
func (cr ColumnRepository) GetColumnByColumnID(columnID string) (model.Column, error) {
	var column model.Column
	err := common.DB.Where("column_id = ?", columnID).First(&column).Error
	return column, err
}

// 获取专栏列表
func (cr ColumnRepository) GetColumns(req *vo.ColumnListRequest) ([]*model.Column, int64, error) {
	var list []*model.Column
	db := common.DB.Model(&model.Column{}).Order("created_at DESC")

	title := strings.TrimSpace(req.Title)
	if title != "" {
		db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title))
	}
	columnID := strings.TrimSpace(req.ColumnID)
	if req.ColumnID != "" {
		db = db.Where("column_id = ?", fmt.Sprintf("%s", columnID))
	}
	projectID := strings.TrimSpace(req.ProjectID)
	if req.ProjectID != "" {
		db = db.Where("project_id = ?", fmt.Sprintf("%s", projectID))
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

// 创建专栏
func (cr ColumnRepository) CreateColumn(column *model.Column) error {
	err := common.DB.Create(column).Error
	return err
}

// 更新专栏
func (cr ColumnRepository) UpdateColumn(column *model.Column) error {
	err := common.DB.Model(column).Updates(column).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (cr ColumnRepository) BatchDeleteColumnByIds(ids []string) error {
	var columns []model.Column
	for _, id := range ids {
		// 根据ID获取标签
		column, err := cr.GetColumnByColumnID(id)
		if err != nil {
			return errors.New(fmt.Sprintf("未获取到ID为%s的专栏", id))
		}
		columns = append(columns, column)
	}

	err := common.DB.Delete(&columns).Error

	return err
}
