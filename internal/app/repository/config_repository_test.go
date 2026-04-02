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

// MockConfigRepository 模拟 ConfigRepository
type MockConfigRepository struct {
	mock.Mock
}

// CreateConfig 模拟创建配置方法
func (m *MockConfigRepository) CreateConfig(ctx context.Context, config *model.Config) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

// GetConfigByConfigID 模拟根据配置ID获取配置方法
func (m *MockConfigRepository) GetConfigByConfigID(ctx context.Context, configID string) (model.Config, error) {
	args := m.Called(ctx, configID)
	return args.Get(0).(model.Config), args.Error(1)
}

// GetConfigs 模拟获取配置列表方法
func (m *MockConfigRepository) GetConfigs(ctx context.Context, req *vo.ConfigListRequest) ([]*model.Config, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*model.Config), args.Get(1).(int64), args.Error(2)
}

// UpdateConfig 模拟更新配置方法
func (m *MockConfigRepository) UpdateConfig(ctx context.Context, config *model.Config) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

// BatchDeleteConfigByIds 模拟批量删除配置方法
func (m *MockConfigRepository) BatchDeleteConfigByIds(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

// 创建测试配置数据
func createTestConfig() *model.Config {
	return &model.Config{
		Model: model.Model{
			ID: 1,
		},
		ConfigID:    "cfg1234567",
		ProjectID:   "proj123456",
		Alias:       "test_config",
		Title:       "测试配置",
		Description: "这是一个测试配置",
		Type:        1,
		Info:        "配置内容",
		MDContent:   "# 配置内容",
		Status:      1,
	}
}

// 创建测试项目数据（用于关联）
func createTestProjectForConfig() *model.Project {
	return &model.Project{
		Model: model.Model{
			ID: 1,
		},
		ProjectID:   "proj123456",
		Name:        "test_project",
		Title:       "测试项目",
		Description: "这是一个测试项目",
		Domain:      "https://example.com",
		Status:      1,
	}
}

// TestConfigRepository_CreateConfig 测试创建配置
func TestConfigRepository_CreateConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   *model.Config
		hasError bool
	}{
		{
			name:     "成功创建配置",
			config:   createTestConfig(),
			hasError: false,
		},
		{
			name: "创建配置失败-重复别名",
			config: &model.Config{
				ConfigID:    "cfg7654321",
				ProjectID:   "proj123456",
				Alias:       "test_config",
				Title:       "重复配置",
				Description: "这是一个重复的配置",
				Type:        1,
				Info:        "配置内容",
				Status:      1,
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			mockRepo := new(MockConfigRepository)

			var expectedError error
			if tt.hasError {
				expectedError = errors.New("配置别名已存在")
			}

			// 设置mock期望
			mockRepo.On("CreateConfig", mock.Anything, tt.config).Return(expectedError)

			// 执行测试
			err := mockRepo.CreateConfig(context.Background(), tt.config)

			// 验证结果
			if tt.hasError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "已存在")
			} else {
				assert.NoError(t, err)
			}

			// 验证mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestConfigRepository_GetConfigByConfigID 测试根据配置ID获取配置
func TestConfigRepository_GetConfigByConfigID(t *testing.T) {
	testConfig := createTestConfig()
	testProject := createTestProjectForConfig()
	testConfig.Project = testProject

	tests := []struct {
		name       string
		configID   string
		expected   model.Config
		hasError   bool
		checkEmpty bool
	}{
		{
			name:     "成功获取配置",
			configID: "cfg1234567",
			expected: *testConfig,
			hasError: false,
		},
		{
			name:       "配置不存在",
			configID:   "nonexistent",
			expected:   model.Config{},
			hasError:   true,
			checkEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			mockRepo := new(MockConfigRepository)

			var expectedError error
			if tt.hasError {
				expectedError = errors.New("配置不存在")
			}

			// 设置mock期望
			mockRepo.On("GetConfigByConfigID", mock.Anything, tt.configID).Return(tt.expected, expectedError)

			// 执行测试
			result, err := mockRepo.GetConfigByConfigID(context.Background(), tt.configID)

			// 验证结果
			if tt.hasError {
				assert.Error(t, err)
				if tt.checkEmpty {
					assert.Empty(t, result.ConfigID)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ConfigID, result.ConfigID)
				assert.Equal(t, tt.expected.Title, result.Title)
			}

			// 验证mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestConfigRepository_GetConfigs 测试获取配置列表
func TestConfigRepository_GetConfigs(t *testing.T) {
	testConfig1 := createTestConfig()
	testConfig2 := &model.Config{
		Model: model.Model{
			ID: 2,
		},
		ConfigID:    "cfg7654321",
		ProjectID:   "proj123456",
		Alias:       "another_config",
		Title:       "另一个配置",
		Description: "这是另一个测试配置",
		Type:        2,
		Info:        "JSON内容",
		MDContent:   "",
		Status:      1,
	}

	testConfigs := []*model.Config{testConfig1, testConfig2}

	tests := []struct {
		name     string
		request  *vo.ConfigListRequest
		expected int64
	}{
		{
			name: "获取所有配置",
			request: &vo.ConfigListRequest{
				PageNum:  1,
				PageSize: 10,
			},
			expected: 2,
		},
		{
			name: "按标题搜索",
			request: &vo.ConfigListRequest{
				Title:    "测试",
				PageNum:  1,
				PageSize: 10,
			},
			expected: 1,
		},
		{
			name: "按配置ID筛选",
			request: &vo.ConfigListRequest{
				ConfigID: "cfg1234567",
				PageNum:  1,
				PageSize: 10,
			},
			expected: 1,
		},
		{
			name: "按项目ID筛选",
			request: &vo.ConfigListRequest{
				ProjectID: "proj123456",
				PageNum:   1,
				PageSize:  10,
			},
			expected: 2,
		},
		{
			name: "按类型筛选",
			request: &vo.ConfigListRequest{
				Type:     1,
				PageNum:  1,
				PageSize: 10,
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			mockRepo := new(MockConfigRepository)

			// 设置mock期望
			mockRepo.On("GetConfigs", mock.Anything, tt.request).Return(testConfigs, tt.expected, nil)

			// 执行测试
			configs, total, err := mockRepo.GetConfigs(context.Background(), tt.request)

			// 验证结果
			assert.NoError(t, err)
			assert.NotNil(t, configs)
			assert.Equal(t, tt.expected, total)
			assert.Len(t, configs, 2)

			// 验证mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestConfigRepository_UpdateConfig 测试更新配置
func TestConfigRepository_UpdateConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   *model.Config
		hasError bool
	}{
		{
			name: "成功更新配置",
			config: &model.Config{
				Model: model.Model{
					ID: 1,
				},
				ConfigID:    "cfg1234567",
				ProjectID:   "proj123456",
				Alias:       "test_config",
				Title:       "更新后的标题",
				Description: "更新后的描述",
				Type:        1,
				Info:        "更新后的内容",
				MDContent:   "# 更新后的内容",
				Status:      1,
			},
			hasError: false,
		},
		{
			name: "更新不存在的配置",
			config: &model.Config{
				ConfigID: "nonexistent",
				Title:    "不存在",
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			mockRepo := new(MockConfigRepository)

			var expectedError error
			if tt.hasError {
				expectedError = errors.New("配置不存在")
			}

			// 设置mock期望
			mockRepo.On("UpdateConfig", mock.Anything, tt.config).Return(expectedError)

			// 执行测试
			err := mockRepo.UpdateConfig(context.Background(), tt.config)

			// 验证结果
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// 验证mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestConfigRepository_BatchDeleteConfigByIds 测试批量删除配置
func TestConfigRepository_BatchDeleteConfigByIds(t *testing.T) {
	tests := []struct {
		name     string
		ids      []string
		hasError bool
	}{
		{
			name:     "成功批量删除配置",
			ids:      []string{"cfg1234567", "cfg7654321"},
			hasError: false,
		},
		{
			name:     "删除包含不存在的配置",
			ids:      []string{"cfg1234567", "nonexistent"},
			hasError: true,
		},
		{
			name:     "空ID列表",
			ids:      []string{},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			mockRepo := new(MockConfigRepository)

			var expectedError error
			if tt.hasError {
				expectedError = errors.New("未获取到ID为nonexistent的配置")
			}

			// 设置mock期望
			mockRepo.On("BatchDeleteConfigByIds", mock.Anything, tt.ids).Return(expectedError)

			// 执行测试
			err := mockRepo.BatchDeleteConfigByIds(context.Background(), tt.ids)

			// 验证结果
			if tt.hasError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "未获取到ID")
			} else {
				assert.NoError(t, err)
			}

			// 验证mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// BenchmarkConfigRepository_CreateConfig 创建配置性能测试
func BenchmarkConfigRepository_CreateConfig(b *testing.B) {
	mockRepo := new(MockConfigRepository)
	testConfig := createTestConfig()

	// 设置mock期望
	mockRepo.On("CreateConfig", mock.Anything, testConfig).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.CreateConfig(context.Background(), testConfig)
	}
}

// BenchmarkConfigRepository_GetConfigByConfigID 获取单个配置性能测试
func BenchmarkConfigRepository_GetConfigByConfigID(b *testing.B) {
	mockRepo := new(MockConfigRepository)
	testConfig := createTestConfig()

	// 设置mock期望
	mockRepo.On("GetConfigByConfigID", mock.Anything, "cfg1234567").Return(*testConfig, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetConfigByConfigID(context.Background(), "cfg1234567")
	}
}

// BenchmarkConfigRepository_GetConfigs 获取配置列表性能测试
func BenchmarkConfigRepository_GetConfigs(b *testing.B) {
	mockRepo := new(MockConfigRepository)
	testConfig1 := createTestConfig()
	testConfig2 := &model.Config{
		Model: model.Model{
			ID: 2,
		},
		ConfigID:  "cfg7654321",
		ProjectID: "proj123456",
		Alias:     "another_config",
		Title:     "另一个配置",
		Type:      2,
		Status:    1,
	}
	testConfigs := []*model.Config{testConfig1, testConfig2}

	req := &vo.ConfigListRequest{
		PageNum:  1,
		PageSize: 10,
	}

	// 设置mock期望
	mockRepo.On("GetConfigs", mock.Anything, req).Return(testConfigs, int64(2), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = mockRepo.GetConfigs(context.Background(), req)
	}
}
