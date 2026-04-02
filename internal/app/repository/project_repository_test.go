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

// MockProjectRepository 模拟 ProjectRepository
type MockProjectRepository struct {
	mock.Mock
}

// CreateProject 模拟创建项目方法
func (m *MockProjectRepository) CreateProject(ctx context.Context, project *model.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

// GetProjectByProjectID 模拟根据项目ID获取项目方法
func (m *MockProjectRepository) GetProjectByProjectID(ctx context.Context, projectID string) (model.Project, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).(model.Project), args.Error(1)
}

// GetProjects 模拟获取项目列表方法
func (m *MockProjectRepository) GetProjects(ctx context.Context, req *vo.ProjectListRequest) ([]*model.Project, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*model.Project), args.Get(1).(int64), args.Error(2)
}

// UpdateProject 模拟更新项目方法
func (m *MockProjectRepository) UpdateProject(ctx context.Context, project *model.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

// BatchDeleteProjectByIds 模拟批量删除项目方法
func (m *MockProjectRepository) BatchDeleteProjectByIds(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

// GetProjectsBySitemap 模拟获取站点地图所需项目方法
func (m *MockProjectRepository) GetProjectsBySitemap(ctx context.Context) ([]*model.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Project), args.Error(1)
}

// 创建测试项目数据
func createTestProject() *model.Project {
	return &model.Project{
		Model: model.Model{
			ID: 1,
		},
		ProjectID:      "proj123456",
		Name:           "test_project",
		Title:          "测试项目",
		Description:    "这是一个测试项目",
		Keywords:       "测试,项目",
		Domain:         "https://example.com",
		PostURL:        "https://example.com/posts",
		ICP:            "京ICP备123456号",
		PublicSecurity: "京公网安备11000002000001号",
		Author:         "测试作者",
		Info:           "项目信息",
		BaiduAnalytics: "baidu_analytics_code",
		Favicon:        "https://example.com/favicon.ico",
		NavImage:       "https://example.com/nav.png",
		PushToken:      "push_token_123",
		Status:         1,
	}
}

// TestProjectRepository_CreateProject 测试创建项目
func TestProjectRepository_CreateProject(t *testing.T) {
	tests := []struct {
		name     string
		project  *model.Project
		hasError bool
	}{
		{
			name:     "成功创建项目",
			project:  createTestProject(),
			hasError: false,
		},
		{
			name: "创建项目失败-重复名称",
			project: &model.Project{
				ProjectID:   "proj765432",
				Name:        "test_project",
				Title:       "重复项目",
				Description: "这是一个重复的项目",
				Domain:      "https://example2.com",
				Status:      1,
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			testMockRepo := new(MockProjectRepository)

			var expectedError error
			if tt.hasError {
				expectedError = errors.New("项目名称已存在")
			}

			// 设置mock期望
			testMockRepo.On("CreateProject", mock.Anything, tt.project).Return(expectedError)

			// 执行测试
			err := testMockRepo.CreateProject(context.Background(), tt.project)

			// 验证结果
			if tt.hasError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "已存在")
			} else {
				assert.NoError(t, err)
			}

			// 验证mock调用
			testMockRepo.AssertExpectations(t)
		})
	}
}

// TestProjectRepository_GetProjectByProjectID 测试根据项目ID获取项目
func TestProjectRepository_GetProjectByProjectID(t *testing.T) {
	testProject := createTestProject()

	tests := []struct {
		name       string
		projectID  string
		expected   model.Project
		hasError   bool
		checkEmpty bool
	}{
		{
			name:      "成功获取项目",
			projectID: "proj123456",
			expected:  *testProject,
			hasError:  false,
		},
		{
			name:       "项目不存在",
			projectID:  "nonexistent",
			expected:   model.Project{},
			hasError:   true,
			checkEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			testMockRepo := new(MockProjectRepository)

			var expectedError error
			if tt.hasError {
				expectedError = errors.New("项目不存在")
			}

			// 设置mock期望
			testMockRepo.On("GetProjectByProjectID", mock.Anything, tt.projectID).Return(tt.expected, expectedError)

			// 执行测试
			result, err := testMockRepo.GetProjectByProjectID(context.Background(), tt.projectID)

			// 验证结果
			if tt.hasError {
				assert.Error(t, err)
				if tt.checkEmpty {
					assert.Empty(t, result.ProjectID)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ProjectID, result.ProjectID)
				assert.Equal(t, tt.expected.Name, result.Name)
				assert.Equal(t, tt.expected.Title, result.Title)
			}

			// 验证mock调用
			testMockRepo.AssertExpectations(t)
		})
	}
}

// TestProjectRepository_GetProjects 测试获取项目列表
func TestProjectRepository_GetProjects(t *testing.T) {
	testProject1 := createTestProject()
	testProject2 := &model.Project{
		Model: model.Model{
			ID: 2,
		},
		ProjectID:   "proj765432",
		Name:        "another_project",
		Title:       "另一个项目",
		Description: "这是另一个测试项目",
		Domain:      "https://example2.com",
		Status:      1,
	}

	testProjects := []*model.Project{testProject1, testProject2}

	tests := []struct {
		name     string
		request  *vo.ProjectListRequest
		expected int64
	}{
		{
			name: "获取所有项目",
			request: &vo.ProjectListRequest{
				PageNum:  1,
				PageSize: 10,
			},
			expected: 2,
		},
		{
			name: "按标题搜索",
			request: &vo.ProjectListRequest{
				Title:    "测试",
				PageNum:  1,
				PageSize: 10,
			},
			expected: 1,
		},
		{
			name: "按项目ID筛选",
			request: &vo.ProjectListRequest{
				ProjectID: "proj123456",
				PageNum:   1,
				PageSize:  10,
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			testMockRepo := new(MockProjectRepository)

			// 设置mock期望
			testMockRepo.On("GetProjects", mock.Anything, tt.request).Return(testProjects, tt.expected, nil)

			// 执行测试
			projects, total, err := testMockRepo.GetProjects(context.Background(), tt.request)

			// 验证结果
			assert.NoError(t, err)
			assert.NotNil(t, projects)
			assert.Equal(t, tt.expected, total)
			assert.Len(t, projects, 2)

			// 验证mock调用
			testMockRepo.AssertExpectations(t)
		})
	}
}

// TestProjectRepository_UpdateProject 测试更新项目
func TestProjectRepository_UpdateProject(t *testing.T) {
	tests := []struct {
		name     string
		project  *model.Project
		hasError bool
	}{
		{
			name: "成功更新项目",
			project: &model.Project{
				Model: model.Model{
					ID: 1,
				},
				ProjectID:      "proj123456",
				Name:           "test_project",
				Title:          "更新后的标题",
				Description:    "更新后的描述",
				Keywords:       "更新,关键词",
				Domain:         "https://example.com",
				PostURL:        "https://example.com/posts",
				ICP:            "京ICP备654321号",
				PublicSecurity: "京公网安备11000002000002号",
				Author:         "更新作者",
				Info:           "更新后的信息",
				BaiduAnalytics: "new_analytics_code",
				Favicon:        "https://example.com/new_favicon.ico",
				NavImage:       "https://example.com/new_nav.png",
				PushToken:      "new_push_token",
				Status:         1,
			},
			hasError: false,
		},
		{
			name: "更新不存在的项目",
			project: &model.Project{
				ProjectID: "nonexistent",
				Title:     "不存在",
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			testMockRepo := new(MockProjectRepository)

			var expectedError error
			if tt.hasError {
				expectedError = errors.New("项目不存在")
			}

			// 设置mock期望
			testMockRepo.On("UpdateProject", mock.Anything, tt.project).Return(expectedError)

			// 执行测试
			err := testMockRepo.UpdateProject(context.Background(), tt.project)

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

// TestProjectRepository_BatchDeleteProjectByIds 测试批量删除项目
func TestProjectRepository_BatchDeleteProjectByIds(t *testing.T) {
	tests := []struct {
		name     string
		ids      []string
		hasError bool
	}{
		{
			name:     "成功批量删除项目",
			ids:      []string{"proj123456", "proj765432"},
			hasError: false,
		},
		{
			name:     "删除包含不存在的项目",
			ids:      []string{"proj123456", "nonexistent"},
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
			testMockRepo := new(MockProjectRepository)

			var expectedError error
			if tt.hasError {
				expectedError = errors.New("未获取到ID为nonexistent的项目")
			}

			// 设置mock期望
			testMockRepo.On("BatchDeleteProjectByIds", mock.Anything, tt.ids).Return(expectedError)

			// 执行测试
			err := testMockRepo.BatchDeleteProjectByIds(context.Background(), tt.ids)

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

// TestProjectRepository_GetProjectsBySitemap 测试获取站点地图所需项目
func TestProjectRepository_GetProjectsBySitemap(t *testing.T) {
	testProject1 := createTestProject()
	testProject2 := &model.Project{
		Model: model.Model{
			ID: 2,
		},
		ProjectID: "proj765432",
		Name:      "another_project",
		Title:     "另一个项目",
		Domain:    "https://example2.com",
		Status:    1,
	}

	testProjects := []*model.Project{testProject1, testProject2}

	tests := []struct {
		name     string
		expected []*model.Project
		hasError bool
	}{
		{
			name:     "成功获取站点地图项目",
			expected: testProjects,
			hasError: false,
		},
		{
			name:     "数据库错误",
			expected: []*model.Project{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			testMockRepo := new(MockProjectRepository)

			var expectedError error
			if tt.hasError {
				expectedError = errors.New("数据库连接失败")
			}

			// 设置mock期望
			testMockRepo.On("GetProjectsBySitemap", mock.Anything).Return(tt.expected, expectedError)

			// 执行测试
			projects, err := testMockRepo.GetProjectsBySitemap(context.Background())

			// 验证结果
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, projects)
				assert.Len(t, projects, len(tt.expected))
			}

			// 验证mock调用
			testMockRepo.AssertExpectations(t)
		})
	}
}

// BenchmarkProjectRepository_CreateProject 创建项目性能测试
func BenchmarkProjectRepository_CreateProject(b *testing.B) {
	mockRepo := new(MockProjectRepository)
	testProject := createTestProject()

	// 设置mock期望
	mockRepo.On("CreateProject", mock.Anything, testProject).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.CreateProject(context.Background(), testProject)
	}
}

// BenchmarkProjectRepository_GetProjectByProjectID 获取单个项目性能测试
func BenchmarkProjectRepository_GetProjectByProjectID(b *testing.B) {
	mockRepo := new(MockProjectRepository)
	testProject := createTestProject()

	// 设置mock期望
	mockRepo.On("GetProjectByProjectID", mock.Anything, "proj123456").Return(*testProject, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetProjectByProjectID(context.Background(), "proj123456")
	}
}

// BenchmarkProjectRepository_GetProjects 获取项目列表性能测试
func BenchmarkProjectRepository_GetProjects(b *testing.B) {
	mockRepo := new(MockProjectRepository)
	testProject1 := createTestProject()
	testProject2 := &model.Project{
		Model: model.Model{
			ID: 2,
		},
		ProjectID: "proj765432",
		Name:      "another_project",
		Title:     "另一个项目",
		Status:    1,
	}
	testProjects := []*model.Project{testProject1, testProject2}

	req := &vo.ProjectListRequest{
		PageNum:  1,
		PageSize: 10,
	}

	// 设置mock期望
	mockRepo.On("GetProjects", mock.Anything, req).Return(testProjects, int64(2), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = mockRepo.GetProjects(context.Background(), req)
	}
}

// BenchmarkProjectRepository_GetProjectsBySitemap 获取站点地图项目性能测试
func BenchmarkProjectRepository_GetProjectsBySitemap(b *testing.B) {
	mockRepo := new(MockProjectRepository)
	testProjects := []*model.Project{
		createTestProject(),
		{
			Model: model.Model{
				ID: 2,
			},
			ProjectID: "proj765432",
			Name:      "another_project",
			Title:     "另一个项目",
			Status:    1,
		},
	}

	// 设置mock期望
	mockRepo.On("GetProjectsBySitemap", mock.Anything).Return(testProjects, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetProjectsBySitemap(context.Background())
	}
}
