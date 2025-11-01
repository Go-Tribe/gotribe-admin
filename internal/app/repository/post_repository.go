// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"fmt"
	"github.com/dengmengmian/ghelper/gconvert"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type IPostRepository interface {
	CreatePost(post *model.Post) error                              // 创建内容
	GetPostByPostID(postID string) (model.Post, error)              // 获取单个内容
	GetPosts(req *vo.PostListRequest) ([]*model.Post, int64, error) // 获取内容列表
	UpdatePost(post *model.Post) error                              // 更新内容
	BatchDeletePostByIds(ids []string) error                        // 批量删除内容
}

type PostRepository struct {
}

// PostRepository构造函数
func NewPostRepository() IPostRepository {
	return PostRepository{}
}

// 获取单个内容
func (pr PostRepository) GetPostByPostID(postID string) (model.Post, error) {
	var post model.Post
	err := common.DB.Where("post_id = ?", postID).First(&post).Error
	//var category model.Category
	//err = common.DB.Where("category_id = ?", post.CategoryID).First(&category).Error
	//post.Category = &category
	return post, err
}

// 获取内容列表
func (pr PostRepository) GetPosts(req *vo.PostListRequest) ([]*model.Post, int64, error) {
	var list []*model.Post
	db := common.DB.Model(&model.Post{}).Order("created_at DESC")

	title := strings.TrimSpace(req.Title)
	if !gconvert.IsEmpty(title) {
		db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title))
	}
	postID := strings.TrimSpace(req.PostID)
	if !gconvert.IsEmpty(postID) {
		db = db.Where("post_id = ?", postID)
	}
	projectID := strings.TrimSpace(req.ProjectID)
	if !gconvert.IsEmpty(projectID) {
		db = db.Where("project_id = ?", projectID)
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
	if err != nil {
		return nil, 0, err
	}
	// 调用 GetPostOther 并处理返回值
	list, err = GetPostOther(list)
	return list, total, err
}

func GetPostOther(posts []*model.Post) ([]*model.Post, error) {
	// 收集所有需要查询的 CategoryID, Tag, ProjectID
	categoryIDs := make([]string, 0, len(posts))
	projectIDs := make([]string, 0, len(posts))
	tagsMap := make(map[uint][]string)
	allTagsSet := make(map[string]bool)

	for _, m := range posts {
		if m.CategoryID != "" {
			categoryIDs = append(categoryIDs, m.CategoryID)
		}
		if m.ProjectID != "" {
			projectIDs = append(projectIDs, m.ProjectID)
		}
		if m.Tag != "" {
			strs := strings.Split(m.Tag, ",")
			for _, str := range strs {
				tag := strings.TrimSpace(str)
				if tag != "" {
					tagsMap[m.ID] = append(tagsMap[m.ID], tag)
					allTagsSet[tag] = true
				}
			}
		}
	}

	// 批量查询 Category
	var categories []*model.Category
	if err := common.DB.Where("category_id IN (?)", categoryIDs).Find(&categories).Error; err != nil {
		return nil, err
	}
	// 批量查询 Tag
	var allTags []*model.Tag
	tagIDs := make([]string, 0, len(allTagsSet))
	for tag := range allTagsSet {
		tagIDs = append(tagIDs, tag)
	}
	if err := common.DB.Where("tag_id IN (?)", tagIDs).Find(&allTags).Error; err != nil {
		return nil, err
	}

	// 批量查询 Project
	var projects []*model.Project
	if err := common.DB.Where("project_id IN (?)", projectIDs).Find(&projects).Error; err != nil {
		return nil, err
	}

	// 将查询结果赋值给 posts
	categoryMap := make(map[string]*model.Category)
	for _, category := range categories {
		common.Log.Info("category", "category", category.CategoryID)
		categoryMap[category.CategoryID] = category
	}

	tagMap := make(map[string]*model.Tag)
	for _, tag := range allTags {
		tagMap[tag.TagID] = tag
	}

	projectMap := make(map[string]*model.Project)
	for _, project := range projects {
		projectMap[project.ProjectID] = project
	}

	for _, m := range posts {
		if category, ok := categoryMap[m.CategoryID]; ok {
			m.Category = category
		}
		var tags []*model.Tag
		for _, tagID := range tagsMap[m.ID] {
			if tag, ok := tagMap[tagID]; ok {
				tags = append(tags, tag)
			}
		}
		m.Tags = tags
		if project, ok := projectMap[m.ProjectID]; ok {
			m.Project = project
		}
	}
	return posts, nil
}

// 创建内容
func (pr PostRepository) CreatePost(post *model.Post) error {
	err := common.DB.Create(post).Error
	return err
}

// 更新内容
func (pr PostRepository) UpdatePost(post *model.Post) error {
	err := common.DB.Model(post).Updates(post).Error
	if err != nil {
		return err
	}

	return err
}

// 批量删除
func (pr PostRepository) BatchDeletePostByIds(ids []string) error {
	var posts []model.Post
	for _, id := range ids {
		// 根据ID获取标签
		post, err := pr.GetPostByPostID(id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%s的内容", id)
		}
		posts = append(posts, post)
	}

	err := common.DB.Delete(&posts).Error

	return err
}
