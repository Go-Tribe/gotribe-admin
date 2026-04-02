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

// MockTagRepository 模拟TagRepository
type MockTagRepository struct {
	mock.Mock
}

// CreateTag 模拟创建标签方法
func (m *MockTagRepository) CreateTag(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	args := m.Called(ctx, tag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}

// GetTagByID 模拟根据ID获取标签方法
func (m *MockTagRepository) GetTagByID(ctx context.Context, id uint) (model.Tag, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Tag), args.Error(1)
}

// GetTags 模拟获取标签列表方法
func (m *MockTagRepository) GetTags(ctx context.Context, req *vo.TagListRequest) ([]*model.Tag, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*model.Tag), args.Get(1).(int64), args.Error(2)
}

// UpdateTag 模拟更新标签方法
func (m *MockTagRepository) UpdateTag(ctx context.Context, tag *model.Tag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

// BatchDeleteTagByIds 模拟批量删除标签方法
func (m *MockTagRepository) BatchDeleteTagByIds(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

// 创建测试标签数据
func createTestTag() *model.Tag {
	return &model.Tag{
		Model: model.Model{
			ID: 1,
		},
		Title:       "测试标签",
		Slug:        "test-tag",
		Description: "这是一个测试标签",
		Color:       "#FF5733",
		Sort:        1,
		Count:       10,
		Status:      1,
	}
}

// TestTagRepository_CreateTag 测试创建标签功能
func TestTagRepository_CreateTag(t *testing.T) {
	mockRepo := new(MockTagRepository)

	// 创建测试标签
	testTag := createTestTag()

	// 设置mock期望 - 正常情况
	mockRepo.On("CreateTag", mock.Anything, testTag).Return(testTag, nil)

	// 执行测试
	result, err := mockRepo.CreateTag(context.Background(), testTag)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testTag.Title, result.Title)
	assert.Equal(t, testTag.Slug, result.Slug)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestTagRepository_CreateTag_Exist 测试创建已存在的标签
func TestTagRepository_CreateTag_Exist(t *testing.T) {
	mockRepo := new(MockTagRepository)

	// 创建测试标签
	testTag := createTestTag()

	// 设置mock期望 - 标签已存在
	mockRepo.On("CreateTag", mock.Anything, testTag).Return(nil, errors.New("测试标签标签已存在"))

	// 执行测试
	result, err := mockRepo.CreateTag(context.Background(), testTag)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "已存在")

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestTagRepository_CreateTag_Error 测试创建标签数据库错误
func TestTagRepository_CreateTag_Error(t *testing.T) {
	mockRepo := new(MockTagRepository)

	// 创建测试标签
	testTag := createTestTag()

	// 设置mock期望 - 数据库错误
	mockRepo.On("CreateTag", mock.Anything, testTag).Return(nil, errors.New("database error"))

	// 执行测试
	result, err := mockRepo.CreateTag(context.Background(), testTag)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestTagRepository_GetTagByID 测试根据ID获取标签
func TestTagRepository_GetTagByID(t *testing.T) {
	mockRepo := new(MockTagRepository)

	// 创建测试标签
	testTag := createTestTag()

	// 设置mock期望
	mockRepo.On("GetTagByID", mock.Anything, uint(1)).Return(*testTag, nil)

	// 执行测试
	result, err := mockRepo.GetTagByID(context.Background(), 1)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "测试标签", result.Title)
	assert.Equal(t, "test-tag", result.Slug)
	assert.Equal(t, "#FF5733", result.Color)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestTagRepository_GetTagByID_NotFound 测试标签不存在的情况
func TestTagRepository_GetTagByID_NotFound(t *testing.T) {
	mockRepo := new(MockTagRepository)

	// 设置mock期望 - 标签不存在
	mockRepo.On("GetTagByID", mock.Anything, uint(999)).Return(model.Tag{}, errors.New("record not found"))

	// 执行测试
	result, err := mockRepo.GetTagByID(context.Background(), 999)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, model.Tag{}, result)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestTagRepository_GetTags 测试获取标签列表
func TestTagRepository_GetTags(t *testing.T) {
	mockRepo := new(MockTagRepository)

	// 创建测试数据
	testTag1 := createTestTag()
	testTag2 := &model.Tag{
		Model: model.Model{
			ID: 2,
		},
		Title:       "第二个标签",
		Slug:        "second-tag",
		Description: "第二个测试标签",
		Color:       "#33FF57",
		Sort:        2,
		Count:       5,
		Status:      1,
	}

	testTags := []*model.Tag{testTag1, testTag2}

	tests := []struct {
		name     string
		request  *vo.TagListRequest
		expected int64
	}{
		{
			name: "获取所有标签",
			request: &vo.TagListRequest{
				PageNum:  1,
				PageSize: 10,
			},
			expected: 2,
		},
		{
			name: "按标题搜索",
			request: &vo.TagListRequest{
				Title:    "测试",
				PageNum:  1,
				PageSize: 10,
			},
			expected: 1,
		},
		{
			name: "按ID搜索",
			request: &vo.TagListRequest{
				ID:       1,
				PageNum:  1,
				PageSize: 10,
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置mock期望
			mockRepo.On("GetTags", mock.Anything, tt.request).Return(testTags, tt.expected, nil)

			// 执行测试
			tags, total, err := mockRepo.GetTags(context.Background(), tt.request)

			// 验证结果
			assert.NoError(t, err)
			assert.NotNil(t, tags)
			assert.Equal(t, tt.expected, total)
			assert.Len(t, tags, 2)

			// 验证mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestTagRepository_GetTags_Error 测试获取标签列表错误情况
func TestTagRepository_GetTags_Error(t *testing.T) {
	mockRepo := new(MockTagRepository)

	request := &vo.TagListRequest{
		PageNum:  1,
		PageSize: 10,
	}

	// 设置mock期望 - 数据库错误
	mockRepo.On("GetTags", mock.Anything, request).Return([]*model.Tag{}, int64(0), errors.New("database connection failed"))

	// 执行测试
	tags, total, err := mockRepo.GetTags(context.Background(), request)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, tags)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestTagRepository_UpdateTag 测试更新标签
func TestTagRepository_UpdateTag(t *testing.T) {
	mockRepo := new(MockTagRepository)

	// 创建测试数据
	testTag := createTestTag()
	testTag.Title = "更新后的标签名称"
	testTag.Color = "#3357FF"
	testTag.Sort = 10

	// 设置mock期望
	mockRepo.On("UpdateTag", mock.Anything, testTag).Return(nil)

	// 执行测试
	err := mockRepo.UpdateTag(context.Background(), testTag)

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestTagRepository_UpdateTag_Error 测试更新标签错误情况
func TestTagRepository_UpdateTag_Error(t *testing.T) {
	mockRepo := new(MockTagRepository)

	// 创建测试数据
	testTag := createTestTag()

	// 设置mock期望 - 更新失败
	mockRepo.On("UpdateTag", mock.Anything, testTag).Return(errors.New("update failed"))

	// 执行测试
	err := mockRepo.UpdateTag(context.Background(), testTag)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, "update failed", err.Error())

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestTagRepository_BatchDeleteTagByIds 测试批量删除标签
func TestTagRepository_BatchDeleteTagByIds(t *testing.T) {
	mockRepo := new(MockTagRepository)

	ids := []uint{1, 2, 3}

	// 设置mock期望
	mockRepo.On("BatchDeleteTagByIds", mock.Anything, ids).Return(nil)

	// 执行测试
	err := mockRepo.BatchDeleteTagByIds(context.Background(), ids)

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestTagRepository_BatchDeleteTagByIds_NotFound 测试批量删除标签不存在
func TestTagRepository_BatchDeleteTagByIds_NotFound(t *testing.T) {
	mockRepo := new(MockTagRepository)

	ids := []uint{1, 999}

	// 设置mock期望 - 部分标签不存在
	mockRepo.On("BatchDeleteTagByIds", mock.Anything, ids).Return(errors.New("未获取到ID为999的标签"))

	// 执行测试
	err := mockRepo.BatchDeleteTagByIds(context.Background(), ids)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "999")

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestTagRepository_BatchDeleteTagByIds_Empty 测试批量删除空ID列表
func TestTagRepository_BatchDeleteTagByIds_Empty(t *testing.T) {
	mockRepo := new(MockTagRepository)

	ids := []uint{}

	// 设置mock期望
	mockRepo.On("BatchDeleteTagByIds", mock.Anything, ids).Return(nil)

	// 执行测试
	err := mockRepo.BatchDeleteTagByIds(context.Background(), ids)

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// BenchmarkTagRepository_CreateTag 创建标签性能测试
func BenchmarkTagRepository_CreateTag(b *testing.B) {
	mockRepo := new(MockTagRepository)
	testTag := createTestTag()

	// 设置mock期望
	mockRepo.On("CreateTag", mock.Anything, testTag).Return(testTag, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.CreateTag(context.Background(), testTag)
	}
}

// BenchmarkTagRepository_GetTagByID 根据ID获取标签性能测试
func BenchmarkTagRepository_GetTagByID(b *testing.B) {
	mockRepo := new(MockTagRepository)
	testTag := createTestTag()

	// 设置mock期望
	mockRepo.On("GetTagByID", mock.Anything, uint(1)).Return(*testTag, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetTagByID(context.Background(), 1)
	}
}

// BenchmarkTagRepository_GetTags 获取标签列表性能测试
func BenchmarkTagRepository_GetTags(b *testing.B) {
	mockRepo := new(MockTagRepository)
	testTag := createTestTag()
	testTags := []*model.Tag{testTag}
	request := &vo.TagListRequest{
		PageNum:  1,
		PageSize: 10,
	}

	// 设置mock期望
	mockRepo.On("GetTags", mock.Anything, request).Return(testTags, int64(1), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = mockRepo.GetTags(context.Background(), request)
	}
}

// BenchmarkTagRepository_UpdateTag 更新标签性能测试
func BenchmarkTagRepository_UpdateTag(b *testing.B) {
	mockRepo := new(MockTagRepository)
	testTag := createTestTag()

	// 设置mock期望
	mockRepo.On("UpdateTag", mock.Anything, testTag).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.UpdateTag(context.Background(), testTag)
	}
}

// BenchmarkTagRepository_BatchDeleteTagByIds 批量删除标签性能测试
func BenchmarkTagRepository_BatchDeleteTagByIds(b *testing.B) {
	mockRepo := new(MockTagRepository)
	ids := []uint{1, 2, 3}

	// 设置mock期望
	mockRepo.On("BatchDeleteTagByIds", mock.Anything, ids).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.BatchDeleteTagByIds(context.Background(), ids)
	}
}

// TestTagRepository_Integration 测试标签管理集成功能
func TestTagRepository_Integration(t *testing.T) {
	mockRepo := new(MockTagRepository)

	// 创建测试标签
	testTag := createTestTag()

	t.Run("完整标签生命周期", func(t *testing.T) {
		// 设置mock期望 - 创建标签
		mockRepo.On("CreateTag", mock.Anything, testTag).Return(testTag, nil)

		// 设置mock期望 - 获取标签
		mockRepo.On("GetTagByID", mock.Anything, uint(1)).Return(*testTag, nil)

		// 设置mock期望 - 更新标签
		updatedTag := *testTag
		updatedTag.Title = "更新后的标签"
		mockRepo.On("UpdateTag", mock.Anything, &updatedTag).Return(nil)

		// 设置mock期望 - 删除标签
		mockRepo.On("BatchDeleteTagByIds", mock.Anything, []uint{1}).Return(nil)

		// 执行测试 - 创建
		result, err := mockRepo.CreateTag(context.Background(), testTag)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// 执行测试 - 获取
		retrievedTag, err := mockRepo.GetTagByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, "测试标签", retrievedTag.Title)

		// 执行测试 - 更新
		err = mockRepo.UpdateTag(context.Background(), &updatedTag)
		assert.NoError(t, err)

		// 执行测试 - 删除
		err = mockRepo.BatchDeleteTagByIds(context.Background(), []uint{1})
		assert.NoError(t, err)

		// 验证mock调用
		mockRepo.AssertExpectations(t)
	})

	t.Run("标签列表查询", func(t *testing.T) {
		// 创建新的mock实例
		listMockRepo := new(MockTagRepository)

		// 创建多个测试标签
		tag1 := createTestTag()
		tag2 := &model.Tag{
			Model: model.Model{
				ID: 2,
			},
			Title: "Go",
			Slug:  "go",
			Color: "#00ADD8",
		}
		tag3 := &model.Tag{
			Model: model.Model{
				ID: 3,
			},
			Title: "Vue",
			Slug:  "vue",
			Color: "#4FC08D",
		}

		testTags := []*model.Tag{tag1, tag2, tag3}

		request := &vo.TagListRequest{
			PageNum:  1,
			PageSize: 10,
		}

		// 设置mock期望 - 获取标签列表
		listMockRepo.On("GetTags", mock.Anything, request).Return(testTags, int64(3), nil)

		// 执行测试
		tags, total, err := listMockRepo.GetTags(context.Background(), request)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, tags, 3)

		// 验证标签属性
		assert.Equal(t, "测试标签", tags[0].Title)
		assert.Equal(t, "Go", tags[1].Title)
		assert.Equal(t, "Vue", tags[2].Title)

		// 验证mock调用
		listMockRepo.AssertExpectations(t)
	})

	t.Run("创建重复标签验证", func(t *testing.T) {
		// 创建新的mock实例
		existMockRepo := new(MockTagRepository)

		// 创建测试标签
		existTag := createTestTag()

		// 设置mock期望 - 标签已存在
		existMockRepo.On("CreateTag", mock.Anything, existTag).Return(nil, errors.New("测试标签标签已存在"))

		// 执行测试
		result, err := existMockRepo.CreateTag(context.Background(), existTag)

		// 验证结果
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "已存在")

		// 验证mock调用
		existMockRepo.AssertExpectations(t)
	})
}
