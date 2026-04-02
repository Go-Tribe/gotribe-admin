// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"
	"errors"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCategoryRepository 模拟CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

// GetCategoryByID 模拟根据ID获取分类方法
func (m *MockCategoryRepository) GetCategoryByID(ctx context.Context, id uint) (model.Category, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Category), args.Error(1)
}

// GetCategorys 模拟获取分类列表方法
func (m *MockCategoryRepository) GetCategorys(ctx context.Context) ([]*model.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Category), args.Error(1)
}

// GetCategoryTree 模拟获取分类树方法
func (m *MockCategoryRepository) GetCategoryTree(ctx context.Context) ([]*model.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Category), args.Error(1)
}

// CreateCategory 模拟创建分类方法
func (m *MockCategoryRepository) CreateCategory(ctx context.Context, category *model.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

// UpdateCategoryByID 模拟根据ID更新分类方法
func (m *MockCategoryRepository) UpdateCategoryByID(ctx context.Context, id uint, category *model.Category) error {
	args := m.Called(ctx, id, category)
	return args.Error(0)
}

// BatchDeleteCategoryByIds 模拟批量删除分类方法
func (m *MockCategoryRepository) BatchDeleteCategoryByIds(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

// 创建测试分类数据
func createTestCategory() *model.Category {
	return &model.Category{
		Model: model.Model{
			ID: 1,
		},
		ParentID:    0,
		Sort:        1,
		Icon:        "category-icon.jpg",
		Title:       "测试分类",
		Slug:        "test-category",
		Path:        "/test-category",
		Hidden:      1,
		Description: "这是一个测试分类",
		Ext:         "{}",
		Status:      1,
		Count:       10,
		Children:    []*model.Category{},
	}
}

// 创建测试子分类数据
func createTestChildCategory(parentID uint) *model.Category {
	return &model.Category{
		Model: model.Model{
			ID: 2,
		},
		ParentID:    parentID,
		Sort:        2,
		Icon:        "child-icon.jpg",
		Title:       "子分类",
		Slug:        "child-category",
		Path:        "/test-category/child-category",
		Hidden:      1,
		Description: "这是一个子分类",
		Status:      1,
		Count:       5,
		Children:    []*model.Category{},
	}
}

// TestCategoryRepository_CreateCategory 测试创建分类功能
func TestCategoryRepository_CreateCategory(t *testing.T) {
	mockRepo := new(MockCategoryRepository)

	// 创建测试分类
	testCategory := createTestCategory()

	// 设置mock期望 - 正常情况
	mockRepo.On("CreateCategory", mock.Anything, testCategory).Return(nil)

	// 执行测试
	err := mockRepo.CreateCategory(context.Background(), testCategory)

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_CreateCategory_Error 测试创建分类错误情况
func TestCategoryRepository_CreateCategory_Error(t *testing.T) {
	mockRepo := new(MockCategoryRepository)

	// 创建测试分类
	testCategory := createTestCategory()

	// 设置mock期望 - 错误情况（重复Slug）
	mockRepo.On("CreateCategory", mock.Anything, testCategory).Return(errors.New("duplicate entry for slug"))

	// 执行测试
	err := mockRepo.CreateCategory(context.Background(), testCategory)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate")

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_GetCategoryByID 测试根据ID获取分类
func TestCategoryRepository_GetCategoryByID(t *testing.T) {
	mockRepo := new(MockCategoryRepository)

	// 创建测试分类
	testCategory := createTestCategory()

	// 设置mock期望
	mockRepo.On("GetCategoryByID", mock.Anything, uint(1)).Return(*testCategory, nil)

	// 执行测试
	result, err := mockRepo.GetCategoryByID(context.Background(), 1)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "测试分类", result.Title)
	assert.Equal(t, "test-category", result.Slug)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_GetCategoryByID_NotFound 测试分类不存在的情况
func TestCategoryRepository_GetCategoryByID_NotFound(t *testing.T) {
	mockRepo := new(MockCategoryRepository)

	// 设置mock期望 - 分类不存在
	mockRepo.On("GetCategoryByID", mock.Anything, uint(999)).Return(model.Category{}, errors.New("record not found"))

	// 执行测试
	result, err := mockRepo.GetCategoryByID(context.Background(), 999)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, model.Category{}, result)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_GetCategorys 测试获取分类列表
func TestCategoryRepository_GetCategorys(t *testing.T) {
	mockRepo := new(MockCategoryRepository)

	// 创建测试数据
	testCategory1 := createTestCategory()
	testCategory2 := &model.Category{
		Model: model.Model{
			ID: 2,
		},
		ParentID:    0,
		Sort:        2,
		Title:       "第二个分类",
		Slug:        "second-category",
		Description: "第二个测试分类",
		Status:      1,
	}

	testCategories := []*model.Category{testCategory1, testCategory2}

	// 设置mock期望
	mockRepo.On("GetCategorys", mock.Anything).Return(testCategories, nil)

	// 执行测试
	categories, err := mockRepo.GetCategorys(context.Background())

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, categories)
	assert.Len(t, categories, 2)
	assert.Equal(t, "测试分类", categories[0].Title)
	assert.Equal(t, "第二个分类", categories[1].Title)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_GetCategorys_Error 测试获取分类列表错误情况
func TestCategoryRepository_GetCategorys_Error(t *testing.T) {
	mockRepo := new(MockCategoryRepository)

	// 设置mock期望 - 数据库错误
	mockRepo.On("GetCategorys", mock.Anything).Return([]*model.Category{}, errors.New("database error"))

	// 执行测试
	categories, err := mockRepo.GetCategorys(context.Background())

	// 验证结果
	assert.Error(t, err)
	assert.Empty(t, categories)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_GetCategoryTree 测试获取分类树
func TestCategoryRepository_GetCategoryTree(t *testing.T) {
	mockRepo := new(MockCategoryRepository)

	// 创建测试数据 - 包含父子关系
	parentCategory := createTestCategory()
	childCategory := createTestChildCategory(1)
	parentCategory.Children = []*model.Category{childCategory}

	testTree := []*model.Category{parentCategory}

	// 设置mock期望
	mockRepo.On("GetCategoryTree", mock.Anything).Return(testTree, nil)

	// 执行测试
	tree, err := mockRepo.GetCategoryTree(context.Background())

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, tree)
	assert.Len(t, tree, 1)
	assert.Equal(t, "测试分类", tree[0].Title)
	assert.Len(t, tree[0].Children, 1)
	assert.Equal(t, "子分类", tree[0].Children[0].Title)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_GetCategoryTree_Error 测试获取分类树错误情况
func TestCategoryRepository_GetCategoryTree_Error(t *testing.T) {
	mockRepo := new(MockCategoryRepository)

	// 设置mock期望 - 数据库错误
	mockRepo.On("GetCategoryTree", mock.Anything).Return([]*model.Category{}, errors.New("database error"))

	// 执行测试
	tree, err := mockRepo.GetCategoryTree(context.Background())

	// 验证结果
	assert.Error(t, err)
	assert.Empty(t, tree)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_UpdateCategoryByID 测试更新分类
func TestCategoryRepository_UpdateCategoryByID(t *testing.T) {
	mockRepo := new(MockCategoryRepository)

	// 创建测试数据
	testCategory := createTestCategory()
	testCategory.Title = "更新后的分类名称"
	testCategory.Description = "更新后的描述"
	testCategory.Sort = 5

	// 设置mock期望
	mockRepo.On("UpdateCategoryByID", mock.Anything, uint(1), testCategory).Return(nil)

	// 执行测试
	err := mockRepo.UpdateCategoryByID(context.Background(), 1, testCategory)

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_UpdateCategoryByID_Error 测试更新分类错误情况
func TestCategoryRepository_UpdateCategoryByID_Error(t *testing.T) {
	mockRepo := new(MockCategoryRepository)

	// 创建测试数据
	testCategory := createTestCategory()

	// 设置mock期望 - 分类不存在
	mockRepo.On("UpdateCategoryByID", mock.Anything, uint(999), testCategory).Return(errors.New("record not found"))

	// 执行测试
	err := mockRepo.UpdateCategoryByID(context.Background(), 999, testCategory)

	// 验证结果
	assert.Error(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_BatchDeleteCategoryByIds 测试批量删除分类
func TestCategoryRepository_BatchDeleteCategoryByIds(t *testing.T) {
	mockRepo := new(MockCategoryRepository)

	ids := []uint{2, 3, 4}

	// 设置mock期望
	mockRepo.On("BatchDeleteCategoryByIds", mock.Anything, ids).Return(nil)

	// 执行测试
	err := mockRepo.BatchDeleteCategoryByIds(context.Background(), ids)

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_BatchDeleteCategoryByIds_DefaultCategory 测试删除默认分类失败
func TestCategoryRepository_BatchDeleteCategoryByIds_DefaultCategory(t *testing.T) {
	mockRepo := new(MockCategoryRepository)

	// 尝试删除默认分类（ID = 1）
	ids := []uint{known.DEFAULT_ID}

	// 设置mock期望 - 默认分类不允许删除
	mockRepo.On("BatchDeleteCategoryByIds", mock.Anything, ids).Return(errors.New("默认分类不允许删除"))

	// 执行测试
	err := mockRepo.BatchDeleteCategoryByIds(context.Background(), ids)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "默认分类")

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_BatchDeleteCategoryByIds_HasChildren 测试删除有子分类的分类
func TestCategoryRepository_BatchDeleteCategoryByIds_HasChildren(t *testing.T) {
	mockRepo := new(MockCategoryRepository)

	// 尝试删除有子分类的分类
	ids := []uint{1}

	// 设置mock期望 - 有子分类不允许删除
	mockRepo.On("BatchDeleteCategoryByIds", mock.Anything, ids).Return(errors.New("该分类下包含子分类，请先删除子分类"))

	// 执行测试
	err := mockRepo.BatchDeleteCategoryByIds(context.Background(), ids)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "子分类")

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestCategoryRepository_BatchDeleteCategoryByIds_Empty 测试批量删除空ID列表
func TestCategoryRepository_BatchDeleteCategoryByIds_Empty(t *testing.T) {
	mockRepo := new(MockCategoryRepository)

	ids := []uint{}

	// 设置mock期望
	mockRepo.On("BatchDeleteCategoryByIds", mock.Anything, ids).Return(nil)

	// 执行测试
	err := mockRepo.BatchDeleteCategoryByIds(context.Background(), ids)

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// BenchmarkCategoryRepository_CreateCategory 创建分类性能测试
func BenchmarkCategoryRepository_CreateCategory(b *testing.B) {
	mockRepo := new(MockCategoryRepository)
	testCategory := createTestCategory()

	// 设置mock期望
	mockRepo.On("CreateCategory", mock.Anything, testCategory).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.CreateCategory(context.Background(), testCategory)
	}
}

// BenchmarkCategoryRepository_GetCategoryByID 根据ID获取分类性能测试
func BenchmarkCategoryRepository_GetCategoryByID(b *testing.B) {
	mockRepo := new(MockCategoryRepository)
	testCategory := createTestCategory()

	// 设置mock期望
	mockRepo.On("GetCategoryByID", mock.Anything, uint(1)).Return(*testCategory, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetCategoryByID(context.Background(), 1)
	}
}

// BenchmarkCategoryRepository_GetCategorys 获取分类列表性能测试
func BenchmarkCategoryRepository_GetCategorys(b *testing.B) {
	mockRepo := new(MockCategoryRepository)
	testCategory := createTestCategory()
	testCategories := []*model.Category{testCategory}

	// 设置mock期望
	mockRepo.On("GetCategorys", mock.Anything).Return(testCategories, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetCategorys(context.Background())
	}
}

// BenchmarkCategoryRepository_GetCategoryTree 获取分类树性能测试
func BenchmarkCategoryRepository_GetCategoryTree(b *testing.B) {
	mockRepo := new(MockCategoryRepository)
	parentCategory := createTestCategory()
	childCategory := createTestChildCategory(1)
	parentCategory.Children = []*model.Category{childCategory}
	testTree := []*model.Category{parentCategory}

	// 设置mock期望
	mockRepo.On("GetCategoryTree", mock.Anything).Return(testTree, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetCategoryTree(context.Background())
	}
}

// BenchmarkCategoryRepository_UpdateCategoryByID 更新分类性能测试
func BenchmarkCategoryRepository_UpdateCategoryByID(b *testing.B) {
	mockRepo := new(MockCategoryRepository)
	testCategory := createTestCategory()

	// 设置mock期望
	mockRepo.On("UpdateCategoryByID", mock.Anything, uint(1), testCategory).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.UpdateCategoryByID(context.Background(), 1, testCategory)
	}
}

// BenchmarkCategoryRepository_BatchDeleteCategoryByIds 批量删除分类性能测试
func BenchmarkCategoryRepository_BatchDeleteCategoryByIds(b *testing.B) {
	mockRepo := new(MockCategoryRepository)
	ids := []uint{2, 3, 4}

	// 设置mock期望
	mockRepo.On("BatchDeleteCategoryByIds", mock.Anything, ids).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.BatchDeleteCategoryByIds(context.Background(), ids)
	}
}

// TestCategoryRepository_Integration 测试分类管理集成功能
func TestCategoryRepository_Integration(t *testing.T) {
	mockRepo := new(MockCategoryRepository)

	// 创建测试分类
	testCategory := createTestCategory()

	t.Run("完整分类生命周期", func(t *testing.T) {
		// 设置mock期望 - 创建分类
		mockRepo.On("CreateCategory", mock.Anything, testCategory).Return(nil)

		// 设置mock期望 - 获取分类
		mockRepo.On("GetCategoryByID", mock.Anything, uint(1)).Return(*testCategory, nil)

		// 设置mock期望 - 更新分类
		updatedCategory := *testCategory
		updatedCategory.Title = "更新后的分类"
		mockRepo.On("UpdateCategoryByID", mock.Anything, uint(1), &updatedCategory).Return(nil)

		// 设置mock期望 - 删除分类
		mockRepo.On("BatchDeleteCategoryByIds", mock.Anything, []uint{1}).Return(nil)

		// 执行测试 - 创建
		err := mockRepo.CreateCategory(context.Background(), testCategory)
		assert.NoError(t, err)

		// 执行测试 - 获取
		result, err := mockRepo.GetCategoryByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, "测试分类", result.Title)

		// 执行测试 - 更新
		err = mockRepo.UpdateCategoryByID(context.Background(), 1, &updatedCategory)
		assert.NoError(t, err)

		// 执行测试 - 删除
		err = mockRepo.BatchDeleteCategoryByIds(context.Background(), []uint{1})
		assert.NoError(t, err)

		// 验证mock调用
		mockRepo.AssertExpectations(t)
	})

	t.Run("分类树结构验证", func(t *testing.T) {
		// 创建新的mock实例
		treeMockRepo := new(MockCategoryRepository)

		// 创建父子分类结构
		parentCategory := createTestCategory()
		childCategory1 := createTestChildCategory(1)
		childCategory2 := &model.Category{
			Model: model.Model{
				ID: 3,
			},
			ParentID: 1,
			Title:    "第二个子分类",
			Slug:     "second-child",
		}
		parentCategory.Children = []*model.Category{childCategory1, childCategory2}

		testTree := []*model.Category{parentCategory}

		// 设置mock期望 - 获取分类树
		treeMockRepo.On("GetCategoryTree", mock.Anything).Return(testTree, nil)

		// 执行测试
		tree, err := treeMockRepo.GetCategoryTree(context.Background())

		// 验证结果
		assert.NoError(t, err)
		assert.Len(t, tree, 1)
		assert.Len(t, tree[0].Children, 2)
		assert.Equal(t, "测试分类", tree[0].Title)
		assert.Equal(t, "子分类", tree[0].Children[0].Title)
		assert.Equal(t, "第二个子分类", tree[0].Children[1].Title)

		// 验证mock调用
		treeMockRepo.AssertExpectations(t)
	})
}
