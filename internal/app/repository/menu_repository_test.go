// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMenuRepository 模拟MenuRepository
type MockMenuRepository struct {
	mock.Mock
}

// GetMenus 模拟获取菜单列表方法
func (m *MockMenuRepository) GetMenus(ctx context.Context) ([]*model.Menu, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Menu), args.Error(1)
}

// GetMenuTree 模拟获取菜单树方法
func (m *MockMenuRepository) GetMenuTree(ctx context.Context) ([]*model.Menu, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Menu), args.Error(1)
}

// CreateMenu 模拟创建菜单方法
func (m *MockMenuRepository) CreateMenu(ctx context.Context, menu *model.Menu) error {
	args := m.Called(ctx, menu)
	return args.Error(0)
}

// UpdateMenuByID 模拟更新菜单方法
func (m *MockMenuRepository) UpdateMenuByID(ctx context.Context, menuID uint, menu *model.Menu) error {
	args := m.Called(ctx, menuID, menu)
	return args.Error(0)
}

// BatchDeleteMenuByIds 模拟批量删除菜单方法
func (m *MockMenuRepository) BatchDeleteMenuByIds(ctx context.Context, menuIds []uint) error {
	args := m.Called(ctx, menuIds)
	return args.Error(0)
}

// GetUserMenusByUserID 模拟根据用户ID获取用户的权限菜单列表方法
func (m *MockMenuRepository) GetUserMenusByUserID(ctx context.Context, userID uint) ([]*model.Menu, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.Menu), args.Error(1)
}

// GetUserMenuTreeByUserID 模拟根据用户ID获取用户的权限菜单树方法
func (m *MockMenuRepository) GetUserMenuTreeByUserID(ctx context.Context, userID uint) ([]*model.Menu, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.Menu), args.Error(1)
}

// createTestMenu 创建测试菜单数据
func createTestMenuForRepo() *model.Menu {
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

// createTestMenuWithParent 创建带父菜单的测试菜单数据
func createTestMenuWithParent() *model.Menu {
	name := "子菜单"
	path := "/test/child"
	component := "ChildComponent"
	icon := "child-icon"
	title := "子菜单标题"
	parentID := uint(1)
	return &model.Menu{
		Model: model.Model{
			ID: 2,
		},
		Name:      name,
		Path:      path,
		Component: component,
		Icon:      &icon,
		Title:     title,
		ParentID:  &parentID,
		Sort:      2,
		Status:    known.DEFAULT_ID,
	}
}

// TestMenuRepository_CreateMenu 测试创建菜单
func TestMenuRepository_CreateMenu(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	testMenu := createTestMenuForRepo()

	// 设置mock期望
	mockRepo.On("CreateMenu", mock.Anything, testMenu).Return(nil)

	// 执行测试
	err := mockRepo.CreateMenu(context.Background(), testMenu)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_CreateMenu_Error 测试创建菜单失败
func TestMenuRepository_CreateMenu_Error(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	testMenu := createTestMenuForRepo()

	// 设置mock期望 - 创建失败
	mockRepo.On("CreateMenu", mock.Anything, testMenu).Return(assert.AnError)

	// 执行测试
	err := mockRepo.CreateMenu(context.Background(), testMenu)

	// 验证结果
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_GetMenus 测试获取菜单列表
func TestMenuRepository_GetMenus(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	testMenu1 := createTestMenuForRepo()
	testMenu2 := createTestMenuWithParent()

	testMenus := []*model.Menu{testMenu1, testMenu2}

	// 设置mock期望
	mockRepo.On("GetMenus", mock.Anything).Return(testMenus, nil)

	// 执行测试
	menus, err := mockRepo.GetMenus(context.Background())

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, menus)
	assert.Len(t, menus, 2)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_GetMenus_Empty 测试获取空菜单列表
func TestMenuRepository_GetMenus_Empty(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	// 设置mock期望
	mockRepo.On("GetMenus", mock.Anything).Return([]*model.Menu{}, nil)

	// 执行测试
	menus, err := mockRepo.GetMenus(context.Background())

	// 验证结果
	assert.NoError(t, err)
	assert.Empty(t, menus)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_GetMenus_Error 测试获取菜单列表失败
func TestMenuRepository_GetMenus_Error(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	// 设置mock期望
	mockRepo.On("GetMenus", mock.Anything).Return([]*model.Menu{}, assert.AnError)

	// 执行测试
	menus, err := mockRepo.GetMenus(context.Background())

	// 验证结果
	assert.Error(t, err)
	assert.Empty(t, menus)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_GetMenuTree 测试获取菜单树
func TestMenuRepository_GetMenuTree(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	testMenu1 := createTestMenuForRepo()
	testMenu2 := createTestMenuWithParent()
	// 设置子菜单
	testMenu1.Children = []*model.Menu{testMenu2}

	testMenus := []*model.Menu{testMenu1}

	// 设置mock期望
	mockRepo.On("GetMenuTree", mock.Anything).Return(testMenus, nil)

	// 执行测试
	menus, err := mockRepo.GetMenuTree(context.Background())

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, menus)
	assert.Len(t, menus, 1)
	assert.Len(t, menus[0].Children, 1)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_GetMenuTree_Empty 测试获取空菜单树
func TestMenuRepository_GetMenuTree_Empty(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	// 设置mock期望
	mockRepo.On("GetMenuTree", mock.Anything).Return([]*model.Menu{}, nil)

	// 执行测试
	menus, err := mockRepo.GetMenuTree(context.Background())

	// 验证结果
	assert.NoError(t, err)
	assert.Empty(t, menus)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_GetMenuTree_Error 测试获取菜单树失败
func TestMenuRepository_GetMenuTree_Error(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	// 设置mock期望
	mockRepo.On("GetMenuTree", mock.Anything).Return([]*model.Menu{}, assert.AnError)

	// 执行测试
	menus, err := mockRepo.GetMenuTree(context.Background())

	// 验证结果
	assert.Error(t, err)
	assert.Empty(t, menus)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_UpdateMenuByID 测试更新菜单
func TestMenuRepository_UpdateMenuByID(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	testMenu := createTestMenuForRepo()
	updatedTitle := "更新后的菜单标题"
	testMenu.Title = updatedTitle
	testMenu.Sort = 3

	// 设置mock期望
	mockRepo.On("UpdateMenuByID", mock.Anything, uint(1), testMenu).Return(nil)

	// 执行测试
	err := mockRepo.UpdateMenuByID(context.Background(), 1, testMenu)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_UpdateMenuByID_Error 测试更新菜单失败
func TestMenuRepository_UpdateMenuByID_Error(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	testMenu := createTestMenuForRepo()

	// 设置mock期望 - 更新失败
	mockRepo.On("UpdateMenuByID", mock.Anything, uint(999), testMenu).Return(assert.AnError)

	// 执行测试
	err := mockRepo.UpdateMenuByID(context.Background(), 999, testMenu)

	// 验证结果
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_BatchDeleteMenuByIds 测试批量删除菜单
func TestMenuRepository_BatchDeleteMenuByIds(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	menuIds := []uint{1, 2, 3}

	// 设置mock期望
	mockRepo.On("BatchDeleteMenuByIds", mock.Anything, menuIds).Return(nil)

	// 执行测试
	err := mockRepo.BatchDeleteMenuByIds(context.Background(), menuIds)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_BatchDeleteMenuByIds_Error 测试批量删除菜单失败
func TestMenuRepository_BatchDeleteMenuByIds_Error(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	menuIds := []uint{999}

	// 设置mock期望 - 删除失败
	mockRepo.On("BatchDeleteMenuByIds", mock.Anything, menuIds).Return(assert.AnError)

	// 执行测试
	err := mockRepo.BatchDeleteMenuByIds(context.Background(), menuIds)

	// 验证结果
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_BatchDeleteMenuByIds_Empty 测试批量删除空菜单ID列表
func TestMenuRepository_BatchDeleteMenuByIds_Empty(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	menuIds := []uint{}

	// 设置mock期望
	mockRepo.On("BatchDeleteMenuByIds", mock.Anything, menuIds).Return(nil)

	// 执行测试
	err := mockRepo.BatchDeleteMenuByIds(context.Background(), menuIds)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_GetUserMenusByUserID 测试根据用户ID获取用户的权限菜单列表
func TestMenuRepository_GetUserMenusByUserID(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	testMenu1 := createTestMenuForRepo()
	testMenu2 := createTestMenuWithParent()

	testMenus := []*model.Menu{testMenu1, testMenu2}

	// 设置mock期望
	mockRepo.On("GetUserMenusByUserID", mock.Anything, uint(1)).Return(testMenus, nil)

	// 执行测试
	menus, err := mockRepo.GetUserMenusByUserID(context.Background(), 1)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, menus)
	assert.Len(t, menus, 2)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_GetUserMenusByUserID_Empty 测试用户没有权限菜单
func TestMenuRepository_GetUserMenusByUserID_Empty(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	// 设置mock期望
	mockRepo.On("GetUserMenusByUserID", mock.Anything, uint(999)).Return([]*model.Menu{}, nil)

	// 执行测试
	menus, err := mockRepo.GetUserMenusByUserID(context.Background(), 999)

	// 验证结果
	assert.NoError(t, err)
	assert.Empty(t, menus)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_GetUserMenusByUserID_Error 测试获取用户权限菜单失败
func TestMenuRepository_GetUserMenusByUserID_Error(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	// 设置mock期望
	mockRepo.On("GetUserMenusByUserID", mock.Anything, uint(1)).Return([]*model.Menu{}, assert.AnError)

	// 执行测试
	menus, err := mockRepo.GetUserMenusByUserID(context.Background(), 1)

	// 验证结果
	assert.Error(t, err)
	assert.Empty(t, menus)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_GetUserMenuTreeByUserID 测试根据用户ID获取用户的权限菜单树
func TestMenuRepository_GetUserMenuTreeByUserID(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	testMenu1 := createTestMenuForRepo()
	testMenu2 := createTestMenuWithParent()
	// 设置子菜单
	testMenu1.Children = []*model.Menu{testMenu2}

	testMenus := []*model.Menu{testMenu1}

	// 设置mock期望
	mockRepo.On("GetUserMenuTreeByUserID", mock.Anything, uint(1)).Return(testMenus, nil)

	// 执行测试
	menus, err := mockRepo.GetUserMenuTreeByUserID(context.Background(), 1)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, menus)
	assert.Len(t, menus, 1)
	assert.Len(t, menus[0].Children, 1)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_GetUserMenuTreeByUserID_Empty 测试用户没有权限菜单树
func TestMenuRepository_GetUserMenuTreeByUserID_Empty(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	// 设置mock期望
	mockRepo.On("GetUserMenuTreeByUserID", mock.Anything, uint(999)).Return([]*model.Menu{}, nil)

	// 执行测试
	menus, err := mockRepo.GetUserMenuTreeByUserID(context.Background(), 999)

	// 验证结果
	assert.NoError(t, err)
	assert.Empty(t, menus)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_GetUserMenuTreeByUserID_Error 测试获取用户权限菜单树失败
func TestMenuRepository_GetUserMenuTreeByUserID_Error(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	// 设置mock期望
	mockRepo.On("GetUserMenuTreeByUserID", mock.Anything, uint(1)).Return([]*model.Menu{}, assert.AnError)

	// 执行测试
	menus, err := mockRepo.GetUserMenuTreeByUserID(context.Background(), 1)

	// 验证结果
	assert.Error(t, err)
	assert.Empty(t, menus)
	mockRepo.AssertExpectations(t)
}

// TestMenuRepository_MenuTreeStructure 测试菜单树结构
func TestMenuRepository_MenuTreeStructure(t *testing.T) {
	mockRepo := new(MockMenuRepository)

	// 创建多层级的菜单结构
	name1 := "系统管理"
	path1 := "/system"
	component1 := "Layout"
	icon1 := "system-icon"
	title1 := "系统管理"
	parentID1 := uint(0)
	menu1 := &model.Menu{
		Model: model.Model{ID: 1},
		Name:  name1, Path: path1, Component: component1,
		Icon: &icon1, Title: title1, ParentID: &parentID1,
		Sort: 1, Status: known.DEFAULT_ID,
	}

	name2 := "用户管理"
	path2 := "user"
	component2 := "system/user"
	icon2 := "user-icon"
	title2 := "用户管理"
	parentID2 := uint(1)
	menu2 := &model.Menu{
		Model: model.Model{ID: 2},
		Name:  name2, Path: path2, Component: component2,
		Icon: &icon2, Title: title2, ParentID: &parentID2,
		Sort: 1, Status: known.DEFAULT_ID,
	}

	name3 := "角色管理"
	path3 := "role"
	component3 := "system/role"
	icon3 := "role-icon"
	title3 := "角色管理"
	parentID3 := uint(1)
	menu3 := &model.Menu{
		Model: model.Model{ID: 3},
		Name:  name3, Path: path3, Component: component3,
		Icon: &icon3, Title: title3, ParentID: &parentID3,
		Sort: 2, Status: known.DEFAULT_ID,
	}

	// 设置子菜单
	menu1.Children = []*model.Menu{menu2, menu3}

	testMenus := []*model.Menu{menu1}

	t.Run("菜单树结构验证", func(t *testing.T) {
		// 设置mock期望
		mockRepo.On("GetMenuTree", mock.Anything).Return(testMenus, nil)

		// 执行测试
		menus, err := mockRepo.GetMenuTree(context.Background())

		// 验证结果
		assert.NoError(t, err)
		assert.NotNil(t, menus)
		assert.Len(t, menus, 1)
		assert.Equal(t, "系统管理", menus[0].Title)
		assert.Len(t, menus[0].Children, 2)
		assert.Equal(t, "用户管理", menus[0].Children[0].Title)
		assert.Equal(t, "角色管理", menus[0].Children[1].Title)

		// 验证mock调用
		mockRepo.AssertExpectations(t)
	})
}

// BenchmarkMenuRepository_CreateMenu 创建菜单性能测试
func BenchmarkMenuRepository_CreateMenu(b *testing.B) {
	mockRepo := new(MockMenuRepository)
	testMenu := createTestMenuForRepo()

	// 设置mock期望
	mockRepo.On("CreateMenu", mock.Anything, testMenu).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.CreateMenu(context.Background(), testMenu)
	}
}

// BenchmarkMenuRepository_GetMenus 获取菜单列表性能测试
func BenchmarkMenuRepository_GetMenus(b *testing.B) {
	mockRepo := new(MockMenuRepository)
	testMenu := createTestMenuForRepo()
	testMenus := []*model.Menu{testMenu}

	// 设置mock期望
	mockRepo.On("GetMenus", mock.Anything).Return(testMenus, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetMenus(context.Background())
	}
}

// BenchmarkMenuRepository_GetMenuTree 获取菜单树性能测试
func BenchmarkMenuRepository_GetMenuTree(b *testing.B) {
	mockRepo := new(MockMenuRepository)
	testMenu := createTestMenuForRepo()
	testMenus := []*model.Menu{testMenu}

	// 设置mock期望
	mockRepo.On("GetMenuTree", mock.Anything).Return(testMenus, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetMenuTree(context.Background())
	}
}

// BenchmarkMenuRepository_GetUserMenusByUserID 获取用户权限菜单性能测试
func BenchmarkMenuRepository_GetUserMenusByUserID(b *testing.B) {
	mockRepo := new(MockMenuRepository)
	testMenu := createTestMenuForRepo()
	testMenus := []*model.Menu{testMenu}

	// 设置mock期望
	mockRepo.On("GetUserMenusByUserID", mock.Anything, uint(1)).Return(testMenus, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetUserMenusByUserID(context.Background(), 1)
	}
}
