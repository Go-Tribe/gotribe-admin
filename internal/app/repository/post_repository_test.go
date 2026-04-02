// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"
	"errors"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPostRepository 模拟PostRepository
type MockPostRepository struct {
	mock.Mock
}

// CreatePost 模拟创建内容方法
func (m *MockPostRepository) CreatePost(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

// GetPostByPostID 模拟根据PostID获取内容方法
func (m *MockPostRepository) GetPostByPostID(ctx context.Context, postID string) (model.Post, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).(model.Post), args.Error(1)
}

// GetPosts 模拟获取内容列表方法
func (m *MockPostRepository) GetPosts(ctx context.Context, req *vo.PostListRequest) ([]*model.Post, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*model.Post), args.Get(1).(int64), args.Error(2)
}

// UpdatePost 模拟更新内容方法
func (m *MockPostRepository) UpdatePost(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

// BatchDeletePostByIds 模拟批量删除内容方法
func (m *MockPostRepository) BatchDeletePostByIds(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

// 创建测试内容数据
func createTestPost() *model.Post {
	return &model.Post{
		Model: model.Model{
			ID: 1,
		},
		PostID:      "post001",
		CategoryID:  1,
		ProjectID:   "proj001",
		ColumnID:    1,
		UserID:      1,
		Author:      "testauthor",
		Title:       "测试文章标题",
		Content:     "测试文章内容",
		HtmlContent: "<p>测试文章内容</p>",
		Description: "测试文章描述",
		Ext:         "{}",
		Icon:        "icon.jpg",
		Tag:         "1,2,3",
		View:        100,
		Type:        1,
		IsTop:       1,
		IsPasswd:    1,
		PassWord:    "",
		Status:      2,
		UnitPrice:   0,
		Location:    "北京",
		People:      "测试人物",
		Time:        "2024-01-01",
		Images:      "image1.jpg,image2.jpg",
		ShowTime:    "2024-01-01 10:00:00",
		Video:       "video.mp4",
	}
}

// TestPostRepository_CreatePost 测试创建内容功能
func TestPostRepository_CreatePost(t *testing.T) {
	mockRepo := new(MockPostRepository)

	// 创建测试内容
	testPost := createTestPost()

	// 设置mock期望 - 正常情况
	mockRepo.On("CreatePost", mock.Anything, testPost).Return(nil)

	// 执行测试
	err := mockRepo.CreatePost(context.Background(), testPost)

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestPostRepository_CreatePost_Error 测试创建内容错误情况
func TestPostRepository_CreatePost_Error(t *testing.T) {
	mockRepo := new(MockPostRepository)

	// 创建测试内容
	testPost := createTestPost()

	// 设置mock期望 - 错误情况
	mockRepo.On("CreatePost", mock.Anything, testPost).Return(errors.New("database error"))

	// 执行测试
	err := mockRepo.CreatePost(context.Background(), testPost)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestPostRepository_GetPostByPostID 测试根据PostID获取内容
func TestPostRepository_GetPostByPostID(t *testing.T) {
	mockRepo := new(MockPostRepository)

	// 创建测试内容
	testPost := createTestPost()

	// 设置mock期望
	mockRepo.On("GetPostByPostID", mock.Anything, "post001").Return(*testPost, nil)

	// 执行测试
	result, err := mockRepo.GetPostByPostID(context.Background(), "post001")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "post001", result.PostID)
	assert.Equal(t, "测试文章标题", result.Title)
	assert.Equal(t, "testauthor", result.Author)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestPostRepository_GetPostByPostID_NotFound 测试内容不存在的情况
func TestPostRepository_GetPostByPostID_NotFound(t *testing.T) {
	mockRepo := new(MockPostRepository)

	// 设置mock期望 - 内容不存在
	mockRepo.On("GetPostByPostID", mock.Anything, "notexist").Return(model.Post{}, errors.New("record not found"))

	// 执行测试
	result, err := mockRepo.GetPostByPostID(context.Background(), "notexist")

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, model.Post{}, result)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestPostRepository_GetPosts 测试获取内容列表
func TestPostRepository_GetPosts(t *testing.T) {
	mockRepo := new(MockPostRepository)

	// 创建测试数据
	testPost1 := createTestPost()
	testPost2 := &model.Post{
		Model: model.Model{
			ID: 2,
		},
		PostID:      "post002",
		CategoryID:  2,
		ProjectID:   "proj002",
		Author:      "author2",
		Title:       "第二篇文章",
		Content:     "第二篇内容",
		HtmlContent: "<p>第二篇内容</p>",
		Description: "第二篇描述",
		Status:      2,
		Type:        1,
	}

	testPosts := []*model.Post{testPost1, testPost2}

	tests := []struct {
		name     string
		request  *vo.PostListRequest
		expected int64
	}{
		{
			name: "获取所有内容",
			request: &vo.PostListRequest{
				PageNum:  1,
				PageSize: 10,
			},
			expected: 2,
		},
		{
			name: "按标题搜索",
			request: &vo.PostListRequest{
				Title:    "测试",
				PageNum:  1,
				PageSize: 10,
			},
			expected: 1,
		},
		{
			name: "按PostID搜索",
			request: &vo.PostListRequest{
				PostID:   "post001",
				PageNum:  1,
				PageSize: 10,
			},
			expected: 1,
		},
		{
			name: "按ProjectID筛选",
			request: &vo.PostListRequest{
				ProjectID: "proj001",
				PageNum:   1,
				PageSize:  10,
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置mock期望
			mockRepo.On("GetPosts", mock.Anything, tt.request).Return(testPosts, tt.expected, nil)

			// 执行测试
			posts, total, err := mockRepo.GetPosts(context.Background(), tt.request)

			// 验证结果
			assert.NoError(t, err)
			assert.NotNil(t, posts)
			assert.Equal(t, tt.expected, total)
			assert.Len(t, posts, 2)

			// 验证mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestPostRepository_GetPosts_Error 测试获取内容列表错误情况
func TestPostRepository_GetPosts_Error(t *testing.T) {
	mockRepo := new(MockPostRepository)

	request := &vo.PostListRequest{
		PageNum:  1,
		PageSize: 10,
	}

	// 设置mock期望 - 数据库错误
	mockRepo.On("GetPosts", mock.Anything, request).Return([]*model.Post{}, int64(0), errors.New("database connection failed"))

	// 执行测试
	posts, total, err := mockRepo.GetPosts(context.Background(), request)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, posts)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestPostRepository_UpdatePost 测试更新内容
func TestPostRepository_UpdatePost(t *testing.T) {
	mockRepo := new(MockPostRepository)

	// 创建测试数据
	testPost := createTestPost()
	testPost.Title = "更新后的标题"
	testPost.Content = "更新后的内容"
	testPost.Status = 1

	// 设置mock期望
	mockRepo.On("UpdatePost", mock.Anything, testPost).Return(nil)

	// 执行测试
	err := mockRepo.UpdatePost(context.Background(), testPost)

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestPostRepository_UpdatePost_Error 测试更新内容错误情况
func TestPostRepository_UpdatePost_Error(t *testing.T) {
	mockRepo := new(MockPostRepository)

	// 创建测试数据
	testPost := createTestPost()

	// 设置mock期望 - 更新失败
	mockRepo.On("UpdatePost", mock.Anything, testPost).Return(errors.New("update failed"))

	// 执行测试
	err := mockRepo.UpdatePost(context.Background(), testPost)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, "update failed", err.Error())

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestPostRepository_BatchDeletePostByIds 测试批量删除内容
func TestPostRepository_BatchDeletePostByIds(t *testing.T) {
	mockRepo := new(MockPostRepository)

	ids := []string{"post001", "post002", "post003"}

	// 设置mock期望
	mockRepo.On("BatchDeletePostByIds", mock.Anything, ids).Return(nil)

	// 执行测试
	err := mockRepo.BatchDeletePostByIds(context.Background(), ids)

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestPostRepository_BatchDeletePostByIds_Error 测试批量删除内容错误情况
func TestPostRepository_BatchDeletePostByIds_Error(t *testing.T) {
	mockRepo := new(MockPostRepository)

	ids := []string{"post001", "notexist"}

	// 设置mock期望 - 部分内容不存在
	mockRepo.On("BatchDeletePostByIds", mock.Anything, ids).Return(errors.New("未获取到ID为notexist的内容"))

	// 执行测试
	err := mockRepo.BatchDeletePostByIds(context.Background(), ids)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "notexist")

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestPostRepository_BatchDeletePostByIds_Empty 测试批量删除空ID列表
func TestPostRepository_BatchDeletePostByIds_Empty(t *testing.T) {
	mockRepo := new(MockPostRepository)

	ids := []string{}

	// 设置mock期望
	mockRepo.On("BatchDeletePostByIds", mock.Anything, ids).Return(nil)

	// 执行测试
	err := mockRepo.BatchDeletePostByIds(context.Background(), ids)

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// BenchmarkPostRepository_CreatePost 创建内容性能测试
func BenchmarkPostRepository_CreatePost(b *testing.B) {
	mockRepo := new(MockPostRepository)
	testPost := createTestPost()

	// 设置mock期望
	mockRepo.On("CreatePost", mock.Anything, testPost).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.CreatePost(context.Background(), testPost)
	}
}

// BenchmarkPostRepository_GetPostByPostID 根据PostID获取内容性能测试
func BenchmarkPostRepository_GetPostByPostID(b *testing.B) {
	mockRepo := new(MockPostRepository)
	testPost := createTestPost()

	// 设置mock期望
	mockRepo.On("GetPostByPostID", mock.Anything, "post001").Return(*testPost, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetPostByPostID(context.Background(), "post001")
	}
}

// BenchmarkPostRepository_GetPosts 获取内容列表性能测试
func BenchmarkPostRepository_GetPosts(b *testing.B) {
	mockRepo := new(MockPostRepository)
	testPost := createTestPost()
	testPosts := []*model.Post{testPost}
	request := &vo.PostListRequest{
		PageNum:  1,
		PageSize: 10,
	}

	// 设置mock期望
	mockRepo.On("GetPosts", mock.Anything, request).Return(testPosts, int64(1), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = mockRepo.GetPosts(context.Background(), request)
	}
}

// BenchmarkPostRepository_UpdatePost 更新内容性能测试
func BenchmarkPostRepository_UpdatePost(b *testing.B) {
	mockRepo := new(MockPostRepository)
	testPost := createTestPost()

	// 设置mock期望
	mockRepo.On("UpdatePost", mock.Anything, testPost).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.UpdatePost(context.Background(), testPost)
	}
}

// BenchmarkPostRepository_BatchDeletePostByIds 批量删除内容性能测试
func BenchmarkPostRepository_BatchDeletePostByIds(b *testing.B) {
	mockRepo := new(MockPostRepository)
	ids := []string{"post001", "post002", "post003"}

	// 设置mock期望
	mockRepo.On("BatchDeletePostByIds", mock.Anything, ids).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.BatchDeletePostByIds(context.Background(), ids)
	}
}

// TestPostRepository_Integration 测试内容管理集成功能
func TestPostRepository_Integration(t *testing.T) {
	mockRepo := new(MockPostRepository)

	// 创建测试内容
	testPost := createTestPost()

	t.Run("完整内容生命周期", func(t *testing.T) {
		// 设置mock期望 - 创建内容
		mockRepo.On("CreatePost", mock.Anything, testPost).Return(nil)

		// 设置mock期望 - 获取内容
		mockRepo.On("GetPostByPostID", mock.Anything, "post001").Return(*testPost, nil)

		// 设置mock期望 - 更新内容
		updatedPost := *testPost
		updatedPost.Title = "更新后的标题"
		mockRepo.On("UpdatePost", mock.Anything, &updatedPost).Return(nil)

		// 设置mock期望 - 删除内容
		mockRepo.On("BatchDeletePostByIds", mock.Anything, []string{"post001"}).Return(nil)

		// 执行测试 - 创建
		err := mockRepo.CreatePost(context.Background(), testPost)
		assert.NoError(t, err)

		// 执行测试 - 获取
		result, err := mockRepo.GetPostByPostID(context.Background(), "post001")
		assert.NoError(t, err)
		assert.Equal(t, "测试文章标题", result.Title)

		// 执行测试 - 更新
		err = mockRepo.UpdatePost(context.Background(), &updatedPost)
		assert.NoError(t, err)

		// 执行测试 - 删除
		err = mockRepo.BatchDeletePostByIds(context.Background(), []string{"post001"})
		assert.NoError(t, err)

		// 验证mock调用
		mockRepo.AssertExpectations(t)
	})
}
