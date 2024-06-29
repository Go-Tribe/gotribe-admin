// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"errors"
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
		db = db.Where("post_id = ?", fmt.Sprintf("%s", postID))
	}
	projectID := strings.TrimSpace(req.ProjectID)
	if !gconvert.IsEmpty(projectID) {
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
	return GetPostOther(list), total, err
}

func GetPostOther(posts []*model.Post) []*model.Post {
	for _, m := range posts {
		var category *model.Category
		_ = common.DB.Where("category_id = ?", m.CategoryID).First(&category).Error
		m.Category = category
		var tags []*model.Tag
		// 处理tag
		strs := strings.Split(m.Tag, ",")
		// 去除每个字符串首部和尾部的空白字符（如果有的话）
		for i, str := range strs {
			strs[i] = strings.TrimSpace(str)
		}
		_ = common.DB.Where("tag_id in (?)", strs).Find(&tags).Error
		m.Tags = tags
		var project *model.Project
		_ = common.DB.Where("project_id = ?", m.ProjectID).First(&project).Error
		m.Project = project
	}
	return posts
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
			return errors.New(fmt.Sprintf("未获取到ID为%s的内容", id))
		}
		posts = append(posts, post)
	}

	err := common.DB.Delete(&posts).Error

	return err
}
