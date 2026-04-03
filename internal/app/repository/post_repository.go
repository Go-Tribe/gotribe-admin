// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"
	"fmt"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"sort"
	"strconv"
	"strings"

	"github.com/dengmengmian/ghelper/gconvert"
	"gorm.io/gorm"
)

type IPostRepository interface {
	CreatePost(ctx context.Context, post *model.Post) error                              // 创建内容
	GetPostByPostID(ctx context.Context, postID string) (model.Post, error)              // 获取单个内容
	GetPosts(ctx context.Context, req *vo.PostListRequest) ([]*model.Post, int64, error) // 获取内容列表
	UpdatePost(ctx context.Context, post *model.Post) error                              // 更新内容
	BatchDeletePostByIds(ctx context.Context, ids []string) error                        // 批量删除内容
}

type PostRepository struct {
}

// PostRepository构造函数
func NewPostRepository() IPostRepository {
	return PostRepository{}
}

// 获取单个内容
func (pr PostRepository) GetPostByPostID(ctx context.Context, postID string) (model.Post, error) {
	var post model.Post
	err := common.WithContext(ctx).DB().Where("post_id = ?", postID).First(&post).Error
	if err != nil {
		return post, err
	}

	posts, err := GetPostOther(ctx, []*model.Post{&post})
	if err != nil {
		return post, err
	}
	if len(posts) == 0 {
		return post, nil
	}

	return *posts[0], nil
}

// 获取内容列表
func (pr PostRepository) GetPosts(ctx context.Context, req *vo.PostListRequest) ([]*model.Post, int64, error) {
	var list []*model.Post
	db := common.WithContext(ctx).DB().Model(&model.Post{}).Order("created_at DESC")

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
	if req.Status > 0 {
		db = db.Where("status = ?", req.Status)
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
	list, err = GetPostOther(ctx, list)
	return list, total, err
}

func GetPostOther(ctx context.Context, posts []*model.Post) ([]*model.Post, error) {
	// 收集所有需要查询的 CategoryID, Tag, ProjectID
	categoryIDSet := make(map[uint]struct{})
	projectIDSet := make(map[string]struct{})
	postIDs := make([]uint, 0, len(posts))

	for _, m := range posts {
		if m.CategoryID > 0 {
			categoryIDSet[m.CategoryID] = struct{}{}
		}
		if m.ProjectID != "" {
			projectIDSet[m.ProjectID] = struct{}{}
		}
		if m.ID > 0 {
			postIDs = append(postIDs, m.ID)
		}
	}

	// 转换为切片
	categoryIDs := make([]uint, 0, len(categoryIDSet))
	for id := range categoryIDSet {
		categoryIDs = append(categoryIDs, id)
	}

	projectIDs := make([]string, 0, len(projectIDSet))
	for id := range projectIDSet {
		projectIDs = append(projectIDs, id)
	}

	// 批量查询 Category
	var categories []*model.Category
	if len(categoryIDs) > 0 {
		if err := common.WithContext(ctx).DB().Where("id IN (?)", categoryIDs).Find(&categories).Error; err != nil {
			return nil, err
		}
	}
	// 批量查询 PostTag 关联
	postTagMap := make(map[uint][]uint)
	tagIDSet := make(map[uint]struct{})
	if len(postIDs) > 0 {
		var postTags []model.PostTag
		if err := common.WithContext(ctx).DB().Where("post_id IN (?)", postIDs).Find(&postTags).Error; err != nil {
			return nil, err
		}
		for _, postTag := range postTags {
			postTagMap[postTag.PostID] = append(postTagMap[postTag.PostID], postTag.TagID)
			tagIDSet[postTag.TagID] = struct{}{}
		}
	}

	// 批量查询 Tag
	var allTags []*model.Tag
	tagIDs := make([]uint, 0, len(tagIDSet))
	for tagID := range tagIDSet {
		tagIDs = append(tagIDs, tagID)
	}
	if len(tagIDs) > 0 {
		if err := common.WithContext(ctx).DB().Where("id IN (?)", tagIDs).Find(&allTags).Error; err != nil {
			return nil, err
		}
	}

	// 批量查询 Project
	var projects []*model.Project
	if len(projectIDs) > 0 {
		if err := common.WithContext(ctx).DB().Where("project_id IN (?)", projectIDs).Find(&projects).Error; err != nil {
			return nil, err
		}
	}

	// 将查询结果赋值给 posts
	categoryMap := make(map[uint]*model.Category)
	for _, category := range categories {
		categoryMap[category.ID] = category
	}

	tagMap := make(map[uint]*model.Tag)
	for _, tag := range allTags {
		tagMap[tag.ID] = tag
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
		for _, tagID := range postTagMap[m.ID] {
			if tag, ok := tagMap[tagID]; ok {
				tags = append(tags, tag)
			}
		}
		m.Tags = tags
		m.Tag = formatPostTagIDs(tags)
		if project, ok := projectMap[m.ProjectID]; ok {
			m.Project = project
		}
	}
	return posts, nil
}

// 创建内容
func (pr PostRepository) CreatePost(ctx context.Context, post *model.Post) error {
	db := common.WithContext(ctx).DB()
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(post).Error; err != nil {
			return err
		}
		return syncPostTags(tx, post.ID, post.Tag)
	})
}

// 更新内容
func (pr PostRepository) UpdatePost(ctx context.Context, post *model.Post) error {
	db := common.WithContext(ctx).DB()
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Post{}).Where("id = ?", post.ID).Updates(buildPostUpdateMap(post)).Error; err != nil {
			return err
		}
		return syncPostTags(tx, post.ID, post.Tag)
	})
}

// 批量删除
func (pr PostRepository) BatchDeletePostByIds(ctx context.Context, ids []string) error {
	var posts []model.Post
	for _, id := range ids {
		// 根据ID获取标签
		post, err := pr.GetPostByPostID(ctx, id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%s的内容", id)
		}
		posts = append(posts, post)
	}

	db := common.WithContext(ctx).DB()
	return db.Transaction(func(tx *gorm.DB) error {
		postIDs := make([]uint, 0, len(posts))
		for _, post := range posts {
			postIDs = append(postIDs, post.ID)
		}
		if len(postIDs) > 0 {
			if err := tx.Where("post_id IN (?)", postIDs).Delete(&model.PostTag{}).Error; err != nil {
				return err
			}
		}
		return tx.Delete(&posts).Error
	})
}

func syncPostTags(tx *gorm.DB, postID uint, rawTagIDs string) error {
	tagIDs, err := parsePostTagIDs(rawTagIDs)
	if err != nil {
		return err
	}

	if err := tx.Where("post_id = ?", postID).Delete(&model.PostTag{}).Error; err != nil {
		return err
	}
	if len(tagIDs) == 0 {
		return nil
	}

	postTags := make([]model.PostTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		postTags = append(postTags, model.PostTag{
			PostID: postID,
			TagID:  tagID,
		})
	}

	return tx.Create(&postTags).Error
}

func parsePostTagIDs(raw string) ([]uint, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	values := strings.Split(raw, ",")
	seen := make(map[uint]struct{}, len(values))
	tagIDs := make([]uint, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		tagID, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("无效的标签ID: %s", value)
		}
		id := uint(tagID)
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		tagIDs = append(tagIDs, id)
	}

	sort.Slice(tagIDs, func(i, j int) bool {
		return tagIDs[i] < tagIDs[j]
	})
	return tagIDs, nil
}

func formatPostTagIDs(tags []*model.Tag) string {
	if len(tags) == 0 {
		return ""
	}

	tagIDs := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		tagIDs = append(tagIDs, strconv.FormatUint(uint64(tag.ID), 10))
	}
	return strings.Join(tagIDs, ",")
}

func buildPostUpdateMap(post *model.Post) map[string]interface{} {
	return map[string]interface{}{
		"post_id":      post.PostID,
		"category_id":  post.CategoryID,
		"project_id":   post.ProjectID,
		"column_id":    post.ColumnID,
		"user_id":      post.UserID,
		"author":       post.Author,
		"title":        post.Title,
		"content":      post.Content,
		"html_content": post.HtmlContent,
		"description":  post.Description,
		"ext":          post.Ext,
		"icon":         post.Icon,
		"view":         post.View,
		"type":         post.Type,
		"is_top":       post.IsTop,
		"is_passwd":    post.IsPasswd,
		"pass_word":    post.PassWord,
		"status":       post.Status,
		"unit_price":   post.UnitPrice,
		"location":     post.Location,
		"people":       post.People,
		"time":         post.Time,
		"images":       post.Images,
		"show_time":    post.ShowTime,
		"video":        post.Video,
	}
}
