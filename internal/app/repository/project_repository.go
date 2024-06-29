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

type IProjectRepository interface {
	CreateProject(project *model.Project) error                              // 创建项目
	GetProjectByProjectID(projectID string) (model.Project, error)           // 获取单个项目
	GetProjects(req *vo.ProjectListRequest) ([]*model.Project, int64, error) // 获取项目列表
	UpdateProject(project *model.Project) error                              // 更新项目
	BatchDeleteProjectByIds(ids []string) error                              // 批量删除项目
	GetProjectsBySitemap() ([]*model.Project, error)
}

type ProjectRepository struct {
}

// ProjectRepository构造函数
func NewProjectRepository() IProjectRepository {
	return ProjectRepository{}
}

// 获取单个项目
func (pr ProjectRepository) GetProjectByProjectID(projectID string) (model.Project, error) {
	var project model.Project
	err := common.DB.Where("project_id = ?", projectID).First(&project).Error
	return project, err
}

// 获取项目列表
func (pr ProjectRepository) GetProjects(req *vo.ProjectListRequest) ([]*model.Project, int64, error) {
	var list []*model.Project
	db := common.DB.Model(&model.Project{}).Order("created_at DESC")

	title := strings.TrimSpace(req.Title)
	if title != "" {
		db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title))
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

// 创建项目
func (pr ProjectRepository) CreateProject(project *model.Project) error {
	err := common.DB.Create(project).Error
	return err
}

// 更新项目
func (pr ProjectRepository) UpdateProject(project *model.Project) error {
	err := common.DB.Model(project).Updates(project).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (pr ProjectRepository) BatchDeleteProjectByIds(ids []string) error {
	var projects []model.Project
	for _, id := range ids {
		// 根据ID获取标签
		project, err := pr.GetProjectByProjectID(id)
		if err != nil {
			return errors.New(fmt.Sprintf("未获取到ID为%s的项目", id))
		}
		projects = append(projects, project)
	}

	err := common.DB.Delete(&projects).Error

	return err
}

// 获取sitmap所需要的 projects 信息
func (pr ProjectRepository) GetProjectsBySitemap() ([]*model.Project, error) {
	var list []*model.Project
	err := common.DB.Model(&model.Project{}).Order("created_at DESC").Find(&list).Error
	return list, err
}
