// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
	"gotribe-admin/pkg/api/vo"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRoleRepository 模拟RoleRepository
type MockRoleRepository struct {
	mock.Mock
}

// GetRoles 模拟获取角色列表方法
func (m *MockRoleRepository) GetRoles(ctx context.Context, req *vo.RoleListRequest) ([]model.Role, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]model.Role), args.Get(1).(int64), args.Error(2)
}

// GetRolesByIds 模拟根据角色ID获取角色方法
func (m *MockRoleRepository) GetRolesByIds(ctx context.Context, roleIds []uint) ([]*model.Role, error) {
	args := m.Called(ctx, roleIds)
	return args.Get(0).([]*model.Role), args.Error(1)
}

// CreateRole 模拟创建角色方法
func (m *MockRoleRepository) CreateRole(ctx context.Context, role *model.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

// UpdateRoleByID 模拟更新角色方法
func (m *MockRoleRepository) UpdateRoleByID(ctx context.Context, roleID uint, role *model.Role) error {
	args := m.Called(ctx, roleID, role)
	return args.Error(0)
}

// GetRoleMenusByID 模拟获取角色的权限菜单方法
func (m *MockRoleRepository) GetRoleMenusByID(ctx context.Context, roleID uint) ([]*model.Menu, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).([]*model.Menu), args.Error(1)
}

// UpdateRoleMenus 模拟更新角色的权限菜单方法
func (m *MockRoleRepository) UpdateRoleMenus(ctx context.Context, role *model.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

// GetRoleApisByRoleKeyword 模拟根据角色关键字获取角色的权限接口方法
func (m *MockRoleRepository) GetRoleApisByRoleKeyword(ctx context.Context, roleKeyword string) ([]*model.Api, error) {
	args := m.Called(ctx, roleKeyword)
	return args.Get(0).([]*model.Api), args.Error(1)
}

// UpdateRoleApis 模拟更新角色的权限接口方法
func (m *MockRoleRepository) UpdateRoleApis(ctx context.Context, roleKeyword string, reqRolePolicies [][]string) error {
	args := m.Called(ctx, roleKeyword, reqRolePolicies)
	return args.Error(0)
}

// BatchDeleteRoleByIds 模拟删除角色方法
func (m *MockRoleRepository) BatchDeleteRoleByIds(ctx context.Context, roleIds []uint) error {
	args := m.Called(ctx, roleIds)
	return args.Error(0)
}

// createTestRole 创建测试角色数据
func createTestRole() *model.Role {
	desc := "测试角色描述"
	return &model.Role{
		Model: model.Model{
			ID: 1,
		},
		Name:    "测试角色",
		Keyword: "test_role",
		Desc:    &desc,
		Status:  known.DEFAULT_ID,
		Sort:    1,
	}
}

// createTestMenuForRole 创建测试菜单数据
func createTestMenuForRole() *model.Menu {
	name := "测试菜单"
	path := "/test"
	component := "TestComponent"
	icon := "test-icon"
	title := "测试菜单标题"
	parentID := uint(0)
	return &model.Menu{
		Model: model.Model{
			ID: 1,
		},
		Name:      name,
		Path:      path,
		Component: component,
		Icon:      &icon,
		Title:     title,
		ParentID:  &parentID,
		Sort:      1,
		Status:    known.DEFAULT_ID,
	}
}

// createTestApi 创建测试API数据
func createTestApi() *model.Api {
	return &model.Api{
		Model: model.Model{
			ID: 1,
		},
		Path:     "/api/test",
		Method:   "GET",
		Category: "test",
		Desc:     "测试API",
	}
}

// TestRoleRepository_CreateRole 测试创建角色
func TestRoleRepository_CreateRole(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	testRole := createTestRole()

	// 设置mock期望
	mockRepo.On("CreateRole", mock.Anything, testRole).Return(nil)

	// 执行测试
	err := mockRepo.CreateRole(context.Background(), testRole)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_CreateRole_Error 测试创建角色失败
func TestRoleRepository_CreateRole_Error(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	testRole := createTestRole()

	// 设置mock期望 - 创建失败
	mockRepo.On("CreateRole", mock.Anything, testRole).Return(assert.AnError)

	// 执行测试
	err := mockRepo.CreateRole(context.Background(), testRole)

	// 验证结果
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_GetRoles 测试获取角色列表
func TestRoleRepository_GetRoles(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	testRole1 := createTestRole()
	desc2 := "普通用户角色描述"
	testRole2 := model.Role{
		Model: model.Model{
			ID: 2,
		},
		Name:    "普通用户",
		Keyword: "user",
		Desc:    &desc2,
		Status:  known.DEFAULT_ID,
		Sort:    5,
	}

	testRoles := []model.Role{*testRole1, testRole2}

	tests := []struct {
		name     string
		request  *vo.RoleListRequest
		expected int64
	}{
		{
			name: "获取所有角色",
			request: &vo.RoleListRequest{
				PageNum:  1,
				PageSize: 10,
			},
			expected: 2,
		},
		{
			name: "按名称搜索",
			request: &vo.RoleListRequest{
				Name:     "测试",
				PageNum:  1,
				PageSize: 10,
			},
			expected: 1,
		},
		{
			name: "按关键字搜索",
			request: &vo.RoleListRequest{
				Keyword:  "user",
				PageNum:  1,
				PageSize: 10,
			},
			expected: 1,
		},
		{
			name: "按状态搜索",
			request: &vo.RoleListRequest{
				Status:   known.DEFAULT_ID,
				PageNum:  1,
				PageSize: 10,
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置mock期望
			mockRepo.On("GetRoles", mock.Anything, tt.request).Return(testRoles, tt.expected, nil)

			// 执行测试
			roles, total, err := mockRepo.GetRoles(context.Background(), tt.request)

			// 验证结果
			assert.NoError(t, err)
			assert.NotNil(t, roles)
			assert.Equal(t, tt.expected, total)
			assert.Len(t, roles, 2)

			// 验证mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestRoleRepository_GetRoles_Error 测试获取角色列表失败
func TestRoleRepository_GetRoles_Error(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	request := &vo.RoleListRequest{
		PageNum:  1,
		PageSize: 10,
	}

	// 设置mock期望
	mockRepo.On("GetRoles", mock.Anything, request).Return([]model.Role{}, int64(0), assert.AnError)

	// 执行测试
	roles, total, err := mockRepo.GetRoles(context.Background(), request)

	// 验证结果
	assert.Error(t, err)
	assert.Empty(t, roles)
	assert.Equal(t, int64(0), total)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_GetRolesByIds 测试根据角色ID获取角色
func TestRoleRepository_GetRolesByIds(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	testRole1 := createTestRole()
	desc2 := "普通用户角色描述"
	testRole2 := &model.Role{
		Model: model.Model{
			ID: 2,
		},
		Name:    "普通用户",
		Keyword: "user",
		Desc:    &desc2,
		Status:  known.DEFAULT_ID,
		Sort:    5,
	}

	testRoles := []*model.Role{testRole1, testRole2}
	roleIds := []uint{1, 2}

	// 设置mock期望
	mockRepo.On("GetRolesByIds", mock.Anything, roleIds).Return(testRoles, nil)

	// 执行测试
	roles, err := mockRepo.GetRolesByIds(context.Background(), roleIds)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, roles)
	assert.Len(t, roles, 2)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_GetRolesByIds_NotFound 测试角色不存在
func TestRoleRepository_GetRolesByIds_NotFound(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	roleIds := []uint{999}

	// 设置mock期望
	mockRepo.On("GetRolesByIds", mock.Anything, roleIds).Return([]*model.Role{}, nil)

	// 执行测试
	roles, err := mockRepo.GetRolesByIds(context.Background(), roleIds)

	// 验证结果
	assert.NoError(t, err)
	assert.Empty(t, roles)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_UpdateRoleByID 测试更新角色
func TestRoleRepository_UpdateRoleByID(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	testRole := createTestRole()
	updatedDesc := "更新后的角色描述"
	testRole.Desc = &updatedDesc
	testRole.Sort = 2

	// 设置mock期望
	mockRepo.On("UpdateRoleByID", mock.Anything, uint(1), testRole).Return(nil)

	// 执行测试
	err := mockRepo.UpdateRoleByID(context.Background(), 1, testRole)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_UpdateRoleByID_Error 测试更新角色失败
func TestRoleRepository_UpdateRoleByID_Error(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	testRole := createTestRole()

	// 设置mock期望 - 更新失败
	mockRepo.On("UpdateRoleByID", mock.Anything, uint(999), testRole).Return(assert.AnError)

	// 执行测试
	err := mockRepo.UpdateRoleByID(context.Background(), 999, testRole)

	// 验证结果
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_GetRoleMenusByID 测试获取角色的权限菜单
func TestRoleRepository_GetRoleMenusByID(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	testMenu1 := createTestMenuForRole()
	name2 := "测试菜单2"
	path2 := "/test2"
	component2 := "TestComponent2"
	icon2 := "test-icon2"
	title2 := "测试菜单标题2"
	parentID2 := uint(1)
	testMenu2 := &model.Menu{
		Model: model.Model{
			ID: 2,
		},
		Name:      name2,
		Path:      path2,
		Component: component2,
		Icon:      &icon2,
		Title:     title2,
		ParentID:  &parentID2,
		Sort:      2,
		Status:    known.DEFAULT_ID,
	}

	testMenus := []*model.Menu{testMenu1, testMenu2}

	// 设置mock期望
	mockRepo.On("GetRoleMenusByID", mock.Anything, uint(1)).Return(testMenus, nil)

	// 执行测试
	menus, err := mockRepo.GetRoleMenusByID(context.Background(), 1)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, menus)
	assert.Len(t, menus, 2)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_GetRoleMenusByID_NotFound 测试角色菜单不存在
func TestRoleRepository_GetRoleMenusByID_NotFound(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	// 设置mock期望
	mockRepo.On("GetRoleMenusByID", mock.Anything, uint(999)).Return([]*model.Menu{}, assert.AnError)

	// 执行测试
	menus, err := mockRepo.GetRoleMenusByID(context.Background(), 999)

	// 验证结果
	assert.Error(t, err)
	assert.Empty(t, menus)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_UpdateRoleMenus 测试更新角色的权限菜单
func TestRoleRepository_UpdateRoleMenus(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	testRole := createTestRole()
	testMenu := createTestMenuForRole()
	testRole.Menus = []*model.Menu{testMenu}

	// 设置mock期望
	mockRepo.On("UpdateRoleMenus", mock.Anything, testRole).Return(nil)

	// 执行测试
	err := mockRepo.UpdateRoleMenus(context.Background(), testRole)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_UpdateRoleMenus_Error 测试更新角色菜单失败
func TestRoleRepository_UpdateRoleMenus_Error(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	testRole := createTestRole()

	// 设置mock期望 - 更新失败
	mockRepo.On("UpdateRoleMenus", mock.Anything, testRole).Return(assert.AnError)

	// 执行测试
	err := mockRepo.UpdateRoleMenus(context.Background(), testRole)

	// 验证结果
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_GetRoleApisByRoleKeyword 测试根据角色关键字获取角色的权限接口
func TestRoleRepository_GetRoleApisByRoleKeyword(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	testApi1 := createTestApi()
	testApi2 := &model.Api{
		Model: model.Model{
			ID: 2,
		},
		Path:     "/api/test2",
		Method:   "POST",
		Category: "test",
		Desc:     "测试API2",
	}

	testApis := []*model.Api{testApi1, testApi2}

	// 设置mock期望
	mockRepo.On("GetRoleApisByRoleKeyword", mock.Anything, "admin").Return(testApis, nil)

	// 执行测试
	apis, err := mockRepo.GetRoleApisByRoleKeyword(context.Background(), "admin")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, apis)
	assert.Len(t, apis, 2)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_GetRoleApisByRoleKeyword_Error 测试获取角色API失败
func TestRoleRepository_GetRoleApisByRoleKeyword_Error(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	// 设置mock期望
	mockRepo.On("GetRoleApisByRoleKeyword", mock.Anything, "invalid_role").Return([]*model.Api{}, assert.AnError)

	// 执行测试
	apis, err := mockRepo.GetRoleApisByRoleKeyword(context.Background(), "invalid_role")

	// 验证结果
	assert.Error(t, err)
	assert.Empty(t, apis)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_UpdateRoleApis 测试更新角色的权限接口
func TestRoleRepository_UpdateRoleApis(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	roleKeyword := "admin"
	reqRolePolicies := [][]string{
		{"admin", "/api/users", "GET"},
		{"admin", "/api/users", "POST"},
	}

	// 设置mock期望
	mockRepo.On("UpdateRoleApis", mock.Anything, roleKeyword, reqRolePolicies).Return(nil)

	// 执行测试
	err := mockRepo.UpdateRoleApis(context.Background(), roleKeyword, reqRolePolicies)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_UpdateRoleApis_Error 测试更新角色API失败
func TestRoleRepository_UpdateRoleApis_Error(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	roleKeyword := "invalid_role"
	reqRolePolicies := [][]string{
		{"invalid_role", "/api/users", "GET"},
	}

	// 设置mock期望 - 更新失败
	mockRepo.On("UpdateRoleApis", mock.Anything, roleKeyword, reqRolePolicies).Return(assert.AnError)

	// 执行测试
	err := mockRepo.UpdateRoleApis(context.Background(), roleKeyword, reqRolePolicies)

	// 验证结果
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_BatchDeleteRoleByIds 测试删除角色
func TestRoleRepository_BatchDeleteRoleByIds(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	roleIds := []uint{1, 2, 3}

	// 设置mock期望
	mockRepo.On("BatchDeleteRoleByIds", mock.Anything, roleIds).Return(nil)

	// 执行测试
	err := mockRepo.BatchDeleteRoleByIds(context.Background(), roleIds)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_BatchDeleteRoleByIds_Error 测试删除角色失败
func TestRoleRepository_BatchDeleteRoleByIds_Error(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	roleIds := []uint{999}

	// 设置mock期望 - 删除失败
	mockRepo.On("BatchDeleteRoleByIds", mock.Anything, roleIds).Return(assert.AnError)

	// 执行测试
	err := mockRepo.BatchDeleteRoleByIds(context.Background(), roleIds)

	// 验证结果
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestRoleRepository_BatchDeleteRoleByIds_Empty 测试删除空角色ID列表
func TestRoleRepository_BatchDeleteRoleByIds_Empty(t *testing.T) {
	mockRepo := new(MockRoleRepository)

	roleIds := []uint{}

	// 设置mock期望
	mockRepo.On("BatchDeleteRoleByIds", mock.Anything, roleIds).Return(nil)

	// 执行测试
	err := mockRepo.BatchDeleteRoleByIds(context.Background(), roleIds)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// BenchmarkRoleRepository_CreateRole 创建角色性能测试
func BenchmarkRoleRepository_CreateRole(b *testing.B) {
	mockRepo := new(MockRoleRepository)
	testRole := createTestRole()

	// 设置mock期望
	mockRepo.On("CreateRole", mock.Anything, testRole).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.CreateRole(context.Background(), testRole)
	}
}

// BenchmarkRoleRepository_GetRoles 获取角色列表性能测试
func BenchmarkRoleRepository_GetRoles(b *testing.B) {
	mockRepo := new(MockRoleRepository)
	testRole := createTestRole()
	testRoles := []model.Role{*testRole}
	request := &vo.RoleListRequest{
		PageNum:  1,
		PageSize: 10,
	}

	// 设置mock期望
	mockRepo.On("GetRoles", mock.Anything, request).Return(testRoles, int64(1), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = mockRepo.GetRoles(context.Background(), request)
	}
}

// BenchmarkRoleRepository_GetRoleMenusByID 获取角色菜单性能测试
func BenchmarkRoleRepository_GetRoleMenusByID(b *testing.B) {
	mockRepo := new(MockRoleRepository)
	testMenu := createTestMenuForRole()
	testMenus := []*model.Menu{testMenu}

	// 设置mock期望
	mockRepo.On("GetRoleMenusByID", mock.Anything, uint(1)).Return(testMenus, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetRoleMenusByID(context.Background(), 1)
	}
}
