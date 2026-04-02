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

// MockAdRepository 模拟 AdRepository
type MockAdRepository struct {
	mock.Mock
}

// CreateAd 模拟创建广告方法
func (m *MockAdRepository) CreateAd(ctx context.Context, ad *model.Ad) error {
	args := m.Called(ctx, ad)
	return args.Error(0)
}

// GetAdByID 模拟根据ID获取广告方法
func (m *MockAdRepository) GetAdByID(ctx context.Context, id uint) (model.Ad, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Ad), args.Error(1)
}

// GetAds 模拟获取广告列表方法
func (m *MockAdRepository) GetAds(ctx context.Context, req *vo.AdListRequest) ([]*model.Ad, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*model.Ad), args.Get(1).(int64), args.Error(2)
}

// UpdateAd 模拟更新广告方法
func (m *MockAdRepository) UpdateAd(ctx context.Context, ad *model.Ad) error {
	args := m.Called(ctx, ad)
	return args.Error(0)
}

// BatchDeleteAdByIds 模拟批量删除广告方法
func (m *MockAdRepository) BatchDeleteAdByIds(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

// 创建测试广告数据
func createTestAd() *model.Ad {
	return &model.Ad{
		Model: model.Model{
			ID: 1,
		},
		Title:       "测试广告",
		Description: "这是一个测试广告",
		URL:         "https://example.com/ad",
		URLType:     1,
		Sort:        1,
		Status:      2, // 已发布
		SceneID:     1,
		Ext:         `{"key": "value"}`,
		Image:       "https://example.com/ad.jpg",
		Video:       "",
	}
}

// 创建测试广告场景数据（用于关联）
func createTestAdScene() *model.AdScene {
	return &model.AdScene{
		Model: model.Model{
			ID: 1,
		},
		Title:       "首页广告位",
		Description: "首页顶部广告位",
		ProjectID:   "proj123456",
	}
}

// TestAdRepository_CreateAd 测试创建广告
func TestAdRepository_CreateAd(t *testing.T) {
	tests := []struct {
		name     string
		ad       *model.Ad
		hasError bool
	}{
		{
			name:     "成功创建广告",
			ad:       createTestAd(),
			hasError: false,
		},
		{
			name: "创建广告失败-场景不存在",
			ad: &model.Ad{
				Title:       "无效广告",
				Description: "场景不存在的广告",
				URL:         "https://example.com/ad",
				URLType:     1,
				Sort:        1,
				Status:      1,
				SceneID:     999, // 不存在的场景
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			testMockRepo := new(MockAdRepository)

			var expectedError error
			if tt.hasError {
				expectedError = errors.New("广告场景不存在")
			}

			// 设置mock期望
			testMockRepo.On("CreateAd", mock.Anything, tt.ad).Return(expectedError)

			// 执行测试
			err := testMockRepo.CreateAd(context.Background(), tt.ad)

			// 验证结果
			if tt.hasError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "不存在")
			} else {
				assert.NoError(t, err)
			}

			// 验证mock调用
			testMockRepo.AssertExpectations(t)
		})
	}
}

// TestAdRepository_GetAdByID 测试根据ID获取广告
func TestAdRepository_GetAdByID(t *testing.T) {
	testAd := createTestAd()
	testScene := createTestAdScene()
	testAd.Scene = testScene

	tests := []struct {
		name       string
		id         uint
		expected   model.Ad
		hasError   bool
		checkEmpty bool
	}{
		{
			name:     "成功获取广告",
			id:       1,
			expected: *testAd,
			hasError: false,
		},
		{
			name:       "广告不存在",
			id:         999,
			expected:   model.Ad{},
			hasError:   true,
			checkEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			testMockRepo := new(MockAdRepository)

			var expectedError error
			if tt.hasError {
				expectedError = errors.New("广告不存在")
			}

			// 设置mock期望
			testMockRepo.On("GetAdByID", mock.Anything, tt.id).Return(tt.expected, expectedError)

			// 执行测试
			result, err := testMockRepo.GetAdByID(context.Background(), tt.id)

			// 验证结果
			if tt.hasError {
				assert.Error(t, err)
				if tt.checkEmpty {
					assert.Empty(t, result.ID)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.Title, result.Title)
				assert.Equal(t, tt.expected.URL, result.URL)
			}

			// 验证mock调用
			testMockRepo.AssertExpectations(t)
		})
	}
}

// TestAdRepository_GetAds 测试获取广告列表
func TestAdRepository_GetAds(t *testing.T) {
	testAd1 := createTestAd()
	testAd2 := &model.Ad{
		Model: model.Model{
			ID: 2,
		},
		Title:       "另一个广告",
		Description: "这是另一个测试广告",
		URL:         "https://example.com/ad2",
		URLType:     2,
		Sort:        2,
		Status:      1, // 未发布
		SceneID:     1,
		Image:       "https://example.com/ad2.jpg",
	}

	testAds := []*model.Ad{testAd1, testAd2}

	tests := []struct {
		name     string
		request  *vo.AdListRequest
		expected int64
	}{
		{
			name: "获取所有广告",
			request: &vo.AdListRequest{
				PageNum:  1,
				PageSize: 10,
			},
			expected: 2,
		},
		{
			name: "按标题搜索",
			request: &vo.AdListRequest{
				Title:    "测试",
				PageNum:  1,
				PageSize: 10,
			},
			expected: 1,
		},
		{
			name: "按场景ID筛选",
			request: &vo.AdListRequest{
				SceneID:  1,
				PageNum:  1,
				PageSize: 10,
			},
			expected: 2,
		},
		{
			name: "按状态筛选",
			request: &vo.AdListRequest{
				Status:   2,
				PageNum:  1,
				PageSize: 10,
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			testMockRepo := new(MockAdRepository)

			// 设置mock期望
			testMockRepo.On("GetAds", mock.Anything, tt.request).Return(testAds, tt.expected, nil)

			// 执行测试
			ads, total, err := testMockRepo.GetAds(context.Background(), tt.request)

			// 验证结果
			assert.NoError(t, err)
			assert.NotNil(t, ads)
			assert.Equal(t, tt.expected, total)
			assert.Len(t, ads, 2)

			// 验证mock调用
			testMockRepo.AssertExpectations(t)
		})
	}
}

// TestAdRepository_UpdateAd 测试更新广告
func TestAdRepository_UpdateAd(t *testing.T) {
	tests := []struct {
		name     string
		ad       *model.Ad
		hasError bool
	}{
		{
			name: "成功更新广告",
			ad: &model.Ad{
				Model: model.Model{
					ID: 1,
				},
				Title:       "更新后的广告",
				Description: "更新后的描述",
				URL:         "https://example.com/updated",
				URLType:     1,
				Sort:        1,
				Status:      2,
				SceneID:     1,
				Image:       "https://example.com/updated.jpg",
			},
			hasError: false,
		},
		{
			name: "更新不存在的广告",
			ad: &model.Ad{
				Model: model.Model{
					ID: 999,
				},
				Title: "不存在",
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			testMockRepo := new(MockAdRepository)

			var expectedError error
			if tt.hasError {
				expectedError = errors.New("广告不存在")
			}

			// 设置mock期望
			testMockRepo.On("UpdateAd", mock.Anything, tt.ad).Return(expectedError)

			// 执行测试
			err := testMockRepo.UpdateAd(context.Background(), tt.ad)

			// 验证结果
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// 验证mock调用
			testMockRepo.AssertExpectations(t)
		})
	}
}

// TestAdRepository_BatchDeleteAdByIds 测试批量删除广告
func TestAdRepository_BatchDeleteAdByIds(t *testing.T) {
	tests := []struct {
		name     string
		ids      []uint
		hasError bool
	}{
		{
			name:     "成功批量删除广告",
			ids:      []uint{1, 2},
			hasError: false,
		},
		{
			name:     "删除包含不存在的广告",
			ids:      []uint{1, 999},
			hasError: true,
		},
		{
			name:     "空ID列表",
			ids:      []uint{},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			testMockRepo := new(MockAdRepository)

			var expectedError error
			if tt.hasError {
				expectedError = errors.New("未获取到ID为999的推广场景")
			}

			// 设置mock期望
			testMockRepo.On("BatchDeleteAdByIds", mock.Anything, tt.ids).Return(expectedError)

			// 执行测试
			err := testMockRepo.BatchDeleteAdByIds(context.Background(), tt.ids)

			// 验证结果
			if tt.hasError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "未获取到ID")
			} else {
				assert.NoError(t, err)
			}

			// 验证mock调用
			testMockRepo.AssertExpectations(t)
		})
	}
}

// TestAdRepository_Integration 测试广告相关功能的集成
func TestAdRepository_Integration(t *testing.T) {
	t.Run("广告创建和查询流程", func(t *testing.T) {
		mockRepo := new(MockAdRepository)
		testAd := createTestAd()

		// 创建广告
		mockRepo.On("CreateAd", mock.Anything, testAd).Return(nil)
		err := mockRepo.CreateAd(context.Background(), testAd)
		assert.NoError(t, err)

		// 查询广告
		mockRepo.On("GetAdByID", mock.Anything, uint(1)).Return(*testAd, nil)
		result, err := mockRepo.GetAdByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, testAd.Title, result.Title)

		// 更新广告
		testAd.Title = "更新后的标题"
		mockRepo.On("UpdateAd", mock.Anything, testAd).Return(nil)
		err = mockRepo.UpdateAd(context.Background(), testAd)
		assert.NoError(t, err)

		// 删除广告
		mockRepo.On("BatchDeleteAdByIds", mock.Anything, []uint{1}).Return(nil)
		err = mockRepo.BatchDeleteAdByIds(context.Background(), []uint{1})
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("广告列表查询与过滤", func(t *testing.T) {
		mockRepo := new(MockAdRepository)
		testAds := []*model.Ad{
			createTestAd(),
			{
				Model: model.Model{
					ID: 2,
				},
				Title:   "草稿广告",
				Status:  1,
				SceneID: 1,
			},
		}

		// 查询已发布广告
		req := &vo.AdListRequest{
			Status:   2,
			PageNum:  1,
			PageSize: 10,
		}
		mockRepo.On("GetAds", mock.Anything, req).Return(testAds[:1], int64(1), nil)
		ads, total, err := mockRepo.GetAds(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, ads, 1)
		assert.Equal(t, uint(2), ads[0].Status)

		mockRepo.AssertExpectations(t)
	})
}

// BenchmarkAdRepository_CreateAd 创建广告性能测试
func BenchmarkAdRepository_CreateAd(b *testing.B) {
	mockRepo := new(MockAdRepository)
	testAd := createTestAd()

	// 设置mock期望
	mockRepo.On("CreateAd", mock.Anything, testAd).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.CreateAd(context.Background(), testAd)
	}
}

// BenchmarkAdRepository_GetAdByID 获取单个广告性能测试
func BenchmarkAdRepository_GetAdByID(b *testing.B) {
	mockRepo := new(MockAdRepository)
	testAd := createTestAd()

	// 设置mock期望
	mockRepo.On("GetAdByID", mock.Anything, uint(1)).Return(*testAd, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetAdByID(context.Background(), 1)
	}
}

// BenchmarkAdRepository_GetAds 获取广告列表性能测试
func BenchmarkAdRepository_GetAds(b *testing.B) {
	mockRepo := new(MockAdRepository)
	testAd1 := createTestAd()
	testAd2 := &model.Ad{
		Model: model.Model{
			ID: 2,
		},
		Title:   "另一个广告",
		Status:  1,
		SceneID: 1,
	}
	testAds := []*model.Ad{testAd1, testAd2}

	req := &vo.AdListRequest{
		PageNum:  1,
		PageSize: 10,
	}

	// 设置mock期望
	mockRepo.On("GetAds", mock.Anything, req).Return(testAds, int64(2), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = mockRepo.GetAds(context.Background(), req)
	}
}

// BenchmarkAdRepository_UpdateAd 更新广告性能测试
func BenchmarkAdRepository_UpdateAd(b *testing.B) {
	mockRepo := new(MockAdRepository)
	testAd := createTestAd()

	// 设置mock期望
	mockRepo.On("UpdateAd", mock.Anything, testAd).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.UpdateAd(context.Background(), testAd)
	}
}

// BenchmarkAdRepository_BatchDeleteAdByIds 批量删除广告性能测试
func BenchmarkAdRepository_BatchDeleteAdByIds(b *testing.B) {
	mockRepo := new(MockAdRepository)
	ids := []uint{1, 2, 3, 4, 5}

	// 设置mock期望
	mockRepo.On("BatchDeleteAdByIds", mock.Anything, ids).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.BatchDeleteAdByIds(context.Background(), ids)
	}
}
