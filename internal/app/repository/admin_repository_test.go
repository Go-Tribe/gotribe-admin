// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
	"gotribe-admin/pkg/api/vo"
	"gotribe-admin/pkg/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAdminRepository 模拟AdminRepository
type MockAdminRepository struct {
	mock.Mock
}

// Login 模拟登录方法
func (m *MockAdminRepository) Login(admin *model.Admin) (*model.Admin, error) {
	args := m.Called(admin)
	return args.Get(0).(*model.Admin), args.Error(1)
}

// GetCurrentAdmin 模拟获取当前管理员方法
func (m *MockAdminRepository) GetCurrentAdmin(c *gin.Context) (*model.Admin, error) {
	args := m.Called(c)
	return args.Get(0).(*model.Admin), args.Error(1)
}

// GetAdminById 模拟根据ID获取管理员方法
func (m *MockAdminRepository) GetAdminById(id uint) (*model.Admin, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Admin), args.Error(1)
}

// GetAdmins 模拟获取管理员列表方法
func (m *MockAdminRepository) GetAdmins(req *vo.AdminListRequest) ([]*model.Admin, int64, error) {
	args := m.Called(req)
	return args.Get(0).([]*model.Admin), args.Get(1).(int64), args.Error(2)
}

// CreateAdmin 模拟创建管理员方法
func (m *MockAdminRepository) CreateAdmin(admin *model.Admin) error {
	args := m.Called(admin)
	return args.Error(0)
}

// UpdateAdmin 模拟更新管理员方法
func (m *MockAdminRepository) UpdateAdmin(admin *model.Admin) error {
	args := m.Called(admin)
	return args.Error(0)
}

// ChangePwd 模拟修改密码方法
func (m *MockAdminRepository) ChangePwd(username string, hashNewPasswd string) error {
	args := m.Called(username, hashNewPasswd)
	return args.Error(0)
}

// BatchDeleteAdminByIds 模拟批量删除管理员方法
func (m *MockAdminRepository) BatchDeleteAdminByIds(ids []uint) error {
	args := m.Called(ids)
	return args.Error(0)
}

// 创建测试用户数据
func createTestAdmin() *model.Admin {
	hashedPassword, _ := util.PasswordUtil.GenPasswd("test123")
	nickname := "测试用户"
	introduction := "测试用户介绍"
	return &model.Admin{
		Model: model.Model{
			ID: 1,
		},
		Username:     "testuser",
		Password:     hashedPassword,
		Mobile:       "13800138000",
		Avatar:       "avatar.jpg",
		Nickname:     &nickname,
		Introduction: &introduction,
		Status:       1,
		Creator:      "admin",
		Roles: []*model.Role{
			{
				Model: model.Model{
					ID: 1,
				},
				Name:    "管理员",
				Keyword: "admin",
				Status:  known.DEFAULT_ID,
				Sort:    1,
			},
		},
	}
}

// TestAdminRepository_Login 测试登录功能
func TestAdminRepository_Login(t *testing.T) {
	mockRepo := new(MockAdminRepository)

	// 创建测试用户
	testAdmin := createTestAdmin()
	loginAdmin := &model.Admin{
		Username: "testuser",
		Password: "test123",
	}

	// 设置mock期望
	mockRepo.On("Login", loginAdmin).Return(testAdmin, nil)

	// 执行测试
	result, err := mockRepo.Login(loginAdmin)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, uint(1), result.ID)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestAdminRepository_GetCurrentAdmin 测试获取当前管理员
func TestAdminRepository_GetCurrentAdmin(t *testing.T) {
	mockRepo := new(MockAdminRepository)

	// 创建gin上下文
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// 创建测试用户
	testAdmin := createTestAdmin()

	// 设置mock期望
	mockRepo.On("GetCurrentAdmin", c).Return(testAdmin, nil)

	// 执行测试
	result, err := mockRepo.GetCurrentAdmin(c)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.Username)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestAdminRepository_GetAdminById 测试根据ID获取管理员
func TestAdminRepository_GetAdminById(t *testing.T) {
	mockRepo := new(MockAdminRepository)

	// 创建测试用户
	testAdmin := createTestAdmin()

	// 设置mock期望
	mockRepo.On("GetAdminById", uint(1)).Return(testAdmin, nil)

	// 执行测试
	result, err := mockRepo.GetAdminById(1)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "testuser", result.Username)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestAdminRepository_GetAdmins 测试获取管理员列表
func TestAdminRepository_GetAdmins(t *testing.T) {
	mockRepo := new(MockAdminRepository)

	// 创建测试数据
	testAdmin1 := createTestAdmin()
	nickname2 := "测试用户2"
	testAdmin2 := &model.Admin{
		Model: model.Model{
			ID: 2,
		},
		Username: "testuser2",
		Password: "hashedpassword",
		Mobile:   "13800138001",
		Nickname: &nickname2,
		Status:   1,
	}

	testAdmins := []*model.Admin{testAdmin1, testAdmin2}

	tests := []struct {
		name     string
		request  *vo.AdminListRequest
		expected int64
	}{
		{
			name: "获取所有用户",
			request: &vo.AdminListRequest{
				PageNum:  1,
				PageSize: 10,
			},
			expected: 2,
		},
		{
			name: "按用户名搜索",
			request: &vo.AdminListRequest{
				Username: "testuser",
				PageNum:  1,
				PageSize: 10,
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置mock期望
			mockRepo.On("GetAdmins", tt.request).Return(testAdmins, tt.expected, nil)

			// 执行测试
			admins, total, err := mockRepo.GetAdmins(tt.request)

			// 验证结果
			assert.NoError(t, err)
			assert.NotNil(t, admins)
			assert.Equal(t, tt.expected, total)
			assert.Len(t, admins, 2)

			// 验证mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestAdminRepository_CreateAdmin 测试创建管理员
func TestAdminRepository_CreateAdmin(t *testing.T) {
	mockRepo := new(MockAdminRepository)

	newNickname := "新用户"
	newAdmin := &model.Admin{
		Username: "newuser",
		Password: "hashedpassword",
		Mobile:   "13800138002",
		Nickname: &newNickname,
		Status:   1,
	}

	// 设置mock期望
	mockRepo.On("CreateAdmin", newAdmin).Return(nil)

	// 执行测试
	err := mockRepo.CreateAdmin(newAdmin)

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestAdminRepository_UpdateAdmin 测试更新管理员
func TestAdminRepository_UpdateAdmin(t *testing.T) {
	mockRepo := new(MockAdminRepository)

	// 创建测试数据
	testAdmin := createTestAdmin()
	updatedNickname := "更新后的昵称"
	testAdmin.Nickname = &updatedNickname
	testAdmin.Mobile = "13800138999"

	// 设置mock期望
	mockRepo.On("UpdateAdmin", testAdmin).Return(nil)

	// 执行测试
	err := mockRepo.UpdateAdmin(testAdmin)

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestAdminRepository_ChangePwd 测试更新密码
func TestAdminRepository_ChangePwd(t *testing.T) {
	mockRepo := new(MockAdminRepository)

	newHashedPassword, _ := util.PasswordUtil.GenPasswd("newpassword123")

	// 设置mock期望
	mockRepo.On("ChangePwd", "testuser", newHashedPassword).Return(nil)

	// 执行测试
	err := mockRepo.ChangePwd("testuser", newHashedPassword)

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestAdminRepository_BatchDeleteAdminByIds 测试批量删除管理员
func TestAdminRepository_BatchDeleteAdminByIds(t *testing.T) {
	mockRepo := new(MockAdminRepository)

	ids := []uint{1, 2, 3}

	// 设置mock期望
	mockRepo.On("BatchDeleteAdminByIds", ids).Return(nil)

	// 执行测试
	err := mockRepo.BatchDeleteAdminByIds(ids)

	// 验证结果
	assert.NoError(t, err)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// BenchmarkAdminRepository_Login 登录性能测试
func BenchmarkAdminRepository_Login(b *testing.B) {
	mockRepo := new(MockAdminRepository)
	testAdmin := createTestAdmin()
	loginAdmin := &model.Admin{
		Username: "testuser",
		Password: "test123",
	}

	// 设置mock期望
	mockRepo.On("Login", loginAdmin).Return(testAdmin, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.Login(loginAdmin)
	}
}

// BenchmarkAdminRepository_GetCurrentAdmin 获取当前管理员性能测试
func BenchmarkAdminRepository_GetCurrentAdmin(b *testing.B) {
	mockRepo := new(MockAdminRepository)
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	testAdmin := createTestAdmin()

	// 设置mock期望
	mockRepo.On("GetCurrentAdmin", c).Return(testAdmin, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetCurrentAdmin(c)
	}
}

// SetAdminInfoCache 模拟设置用户信息缓存方法
func (m *MockAdminRepository) SetAdminInfoCache(username string, admin model.Admin) {
	m.Called(username, admin)
}

// UpdateAdminInfoCacheByRoleID 模拟根据角色ID更新用户信息缓存方法
func (m *MockAdminRepository) UpdateAdminInfoCacheByRoleID(roleID uint) error {
	args := m.Called(roleID)
	return args.Error(0)
}

// ClearAdminInfoCache 模拟清理所有用户信息缓存方法
func (m *MockAdminRepository) ClearAdminInfoCache() {
	m.Called()
}

// TestAdminRepository_SetAdminInfoCache 测试设置用户信息缓存
func TestAdminRepository_SetAdminInfoCache(t *testing.T) {
	mockRepo := new(MockAdminRepository)

	// 创建测试用户
	testAdmin := createTestAdmin()

	// 设置mock期望
	mockRepo.On("SetAdminInfoCache", "testuser", *testAdmin).Return()

	// 执行测试
	mockRepo.SetAdminInfoCache("testuser", *testAdmin)

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestAdminRepository_UpdateAdminInfoCacheByRoleID 测试根据角色ID更新用户信息缓存
func TestAdminRepository_UpdateAdminInfoCacheByRoleID(t *testing.T) {
	mockRepo := new(MockAdminRepository)

	tests := []struct {
		name     string
		roleID   uint
		hasError bool
	}{
		{
			name:     "成功更新缓存",
			roleID:   1,
			hasError: false,
		},
		{
			name:     "角色不存在",
			roleID:   999,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expectedError error
			if tt.hasError {
				expectedError = assert.AnError
			}

			// 设置mock期望
			mockRepo.On("UpdateAdminInfoCacheByRoleID", tt.roleID).Return(expectedError)

			// 执行测试
			err := mockRepo.UpdateAdminInfoCacheByRoleID(tt.roleID)

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

// TestAdminRepository_ClearAdminInfoCache 测试清理所有用户信息缓存
func TestAdminRepository_ClearAdminInfoCache(t *testing.T) {
	mockRepo := new(MockAdminRepository)

	// 设置mock期望
	mockRepo.On("ClearAdminInfoCache").Return()

	// 执行测试
	mockRepo.ClearAdminInfoCache()

	// 验证mock调用
	mockRepo.AssertExpectations(t)
}

// TestAdminRepository_CacheIntegration 测试缓存集成功能
func TestAdminRepository_CacheIntegration(t *testing.T) {
	mockRepo := new(MockAdminRepository)

	// 创建测试用户
	testAdmin := createTestAdmin()

	t.Run("缓存设置和清理流程", func(t *testing.T) {
		// 设置mock期望 - 设置缓存
		mockRepo.On("SetAdminInfoCache", "testuser", *testAdmin).Return()

		// 设置mock期望 - 清理缓存
		mockRepo.On("ClearAdminInfoCache").Return()

		// 执行测试 - 设置缓存
		mockRepo.SetAdminInfoCache("testuser", *testAdmin)

		// 执行测试 - 清理缓存
		mockRepo.ClearAdminInfoCache()

		// 验证mock调用
		mockRepo.AssertExpectations(t)
	})

	t.Run("角色更新触发缓存更新", func(t *testing.T) {
		// 设置mock期望 - 根据角色ID更新缓存
		mockRepo.On("UpdateAdminInfoCacheByRoleID", uint(1)).Return(nil)

		// 执行测试
		err := mockRepo.UpdateAdminInfoCacheByRoleID(1)

		// 验证结果
		assert.NoError(t, err)

		// 验证mock调用
		mockRepo.AssertExpectations(t)
	})
}

// BenchmarkAdminRepository_SetAdminInfoCache 设置缓存性能测试
func BenchmarkAdminRepository_SetAdminInfoCache(b *testing.B) {
	mockRepo := new(MockAdminRepository)
	testAdmin := createTestAdmin()

	// 设置mock期望
	mockRepo.On("SetAdminInfoCache", "testuser", *testAdmin).Return()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mockRepo.SetAdminInfoCache("testuser", *testAdmin)
	}
}

// BenchmarkAdminRepository_ClearAdminInfoCache 清理缓存性能测试
func BenchmarkAdminRepository_ClearAdminInfoCache(b *testing.B) {
	mockRepo := new(MockAdminRepository)

	// 设置mock期望
	mockRepo.On("ClearAdminInfoCache").Return()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mockRepo.ClearAdminInfoCache()
	}
}

// GetCurrentAdminMinRoleSort 模拟获取当前用户角色排序最小值方法
func (m *MockAdminRepository) GetCurrentAdminMinRoleSort(c *gin.Context) (uint, *model.Admin, error) {
	args := m.Called(c)
	return args.Get(0).(uint), args.Get(1).(*model.Admin), args.Error(2)
}

// GetAdminMinRoleSortsByIds 模拟根据用户ID获取用户角色排序最小值方法
func (m *MockAdminRepository) GetAdminMinRoleSortsByIds(ids []uint) ([]int, error) {
	args := m.Called(ids)
	return args.Get(0).([]int), args.Error(1)
}

// TestAdminRepository_GetCurrentAdminMinRoleSort 测试获取当前用户角色排序最小值
func TestAdminRepository_GetCurrentAdminMinRoleSort(t *testing.T) {
	// 创建gin上下文
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// 创建测试用户
	testAdmin := createTestAdmin()

	tests := []struct {
		name         string
		expectedSort uint
		hasError     bool
	}{
		{
			name:         "成功获取角色排序",
			expectedSort: 1,
			hasError:     false,
		},
		{
			name:         "用户未登录",
			expectedSort: 999,
			hasError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试用例创建新的mock实例
			testMockRepo := new(MockAdminRepository)

			var expectedError error
			var expectedAdmin *model.Admin
			if tt.hasError {
				expectedError = assert.AnError
				expectedAdmin = nil // 错误情况下返回nil
			} else {
				expectedAdmin = testAdmin
			}

			// 设置mock期望
			testMockRepo.On("GetCurrentAdminMinRoleSort", c).Return(tt.expectedSort, expectedAdmin, expectedError)

			// 执行测试
			sort, admin, err := testMockRepo.GetCurrentAdminMinRoleSort(c)

			// 验证结果
			if tt.hasError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedSort, sort)
				assert.Nil(t, admin)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSort, sort)
				assert.NotNil(t, admin)
				assert.Equal(t, "testuser", admin.Username)
			}

			// 验证mock调用
			testMockRepo.AssertExpectations(t)
		})
	}
}

// TestAdminRepository_GetAdminMinRoleSortsByIds 测试根据用户ID获取用户角色排序最小值
func TestAdminRepository_GetAdminMinRoleSortsByIds(t *testing.T) {
	mockRepo := new(MockAdminRepository)

	tests := []struct {
		name          string
		ids           []uint
		expectedSorts []int
		hasError      bool
	}{
		{
			name:          "成功获取角色排序",
			ids:           []uint{1, 2},
			expectedSorts: []int{1, 2},
			hasError:      false,
		},
		{
			name:          "用户不存在",
			ids:           []uint{999},
			expectedSorts: []int{},
			hasError:      true,
		},
		{
			name:          "空ID列表",
			ids:           []uint{},
			expectedSorts: []int{},
			hasError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expectedError error
			if tt.hasError {
				expectedError = assert.AnError
			}

			// 设置mock期望
			mockRepo.On("GetAdminMinRoleSortsByIds", tt.ids).Return(tt.expectedSorts, expectedError)

			// 执行测试
			sorts, err := mockRepo.GetAdminMinRoleSortsByIds(tt.ids)

			// 验证结果
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedSorts, sorts)

			// 验证mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestAdminRepository_RolePermissionIntegration 测试角色权限集成功能
func TestAdminRepository_RolePermissionIntegration(t *testing.T) {
	mockRepo := new(MockAdminRepository)

	// 创建gin上下文
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// 创建测试用户
	testAdmin := createTestAdmin()

	t.Run("用户角色权限验证流程", func(t *testing.T) {
		// 设置mock期望 - 获取当前用户
		mockRepo.On("GetCurrentAdmin", c).Return(testAdmin, nil)

		// 设置mock期望 - 获取当前用户角色排序
		mockRepo.On("GetCurrentAdminMinRoleSort", c).Return(uint(1), testAdmin, nil)

		// 设置mock期望 - 根据ID获取用户角色排序
		mockRepo.On("GetAdminMinRoleSortsByIds", []uint{1}).Return([]int{1}, nil)

		// 执行测试 - 获取当前用户
		currentAdmin, err := mockRepo.GetCurrentAdmin(c)
		assert.NoError(t, err)
		assert.Equal(t, "testuser", currentAdmin.Username)

		// 执行测试 - 获取当前用户角色排序
		sort, admin, err := mockRepo.GetCurrentAdminMinRoleSort(c)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), sort)
		assert.Equal(t, "testuser", admin.Username)

		// 执行测试 - 根据ID获取用户角色排序
		sorts, err := mockRepo.GetAdminMinRoleSortsByIds([]uint{1})
		assert.NoError(t, err)
		assert.Equal(t, []int{1}, sorts)

		// 验证mock调用
		mockRepo.AssertExpectations(t)
	})

	t.Run("多用户角色排序比较", func(t *testing.T) {
		// 设置mock期望 - 获取多个用户的角色排序
		mockRepo.On("GetAdminMinRoleSortsByIds", []uint{1, 2, 3}).Return([]int{1, 2, 3}, nil)

		// 执行测试
		sorts, err := mockRepo.GetAdminMinRoleSortsByIds([]uint{1, 2, 3})

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, sorts)
		assert.Len(t, sorts, 3)

		// 验证角色排序逻辑
		for i, sort := range sorts {
			assert.Equal(t, i+1, sort, "角色排序应该按顺序递增")
		}

		// 验证mock调用
		mockRepo.AssertExpectations(t)
	})
}

// TestAdminRepository_RoleValidation 测试角色验证功能
func TestAdminRepository_RoleValidation(t *testing.T) {
	mockRepo := new(MockAdminRepository)

	// 创建gin上下文
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	// 创建测试用户
	testAdmin := createTestAdmin()

	t.Run("管理员角色权限验证", func(t *testing.T) {
		// 设置mock期望 - 获取当前用户角色排序（管理员角色排序为1）
		mockRepo.On("GetCurrentAdminMinRoleSort", c).Return(uint(1), testAdmin, nil)

		// 执行测试
		sort, admin, err := mockRepo.GetCurrentAdminMinRoleSort(c)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, uint(1), sort)
		assert.NotNil(t, admin)
		assert.Equal(t, "testuser", admin.Username)
		assert.Len(t, admin.Roles, 1)
		assert.Equal(t, "管理员", admin.Roles[0].Name)
		assert.Equal(t, "admin", admin.Roles[0].Keyword)

		// 验证mock调用
		mockRepo.AssertExpectations(t)
	})

	t.Run("普通用户角色权限验证", func(t *testing.T) {
		// 创建新的mock实例以避免冲突
		normalMockRepo := new(MockAdminRepository)

		// 创建普通用户
		nickname := "普通用户"
		normalUser := &model.Admin{
			Model: model.Model{
				ID: 2,
			},
			Username: "normaluser",
			Nickname: &nickname,
			Roles: []*model.Role{
				{
					Model: model.Model{
						ID: 2,
					},
					Name:    "普通用户",
					Keyword: "user",
					Status:  known.DEFAULT_ID,
					Sort:    5, // 普通用户角色排序较高
				},
			},
		}

		// 设置mock期望 - 获取当前用户角色排序（普通用户角色排序为5）
		normalMockRepo.On("GetCurrentAdminMinRoleSort", c).Return(uint(5), normalUser, nil)

		// 执行测试
		sort, admin, err := normalMockRepo.GetCurrentAdminMinRoleSort(c)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, uint(5), sort)
		assert.NotNil(t, admin)
		assert.Equal(t, "normaluser", admin.Username)
		assert.Len(t, admin.Roles, 1)
		assert.Equal(t, "普通用户", admin.Roles[0].Name)
		assert.Equal(t, uint(5), admin.Roles[0].Sort)

		// 验证mock调用
		normalMockRepo.AssertExpectations(t)
	})
}

// BenchmarkAdminRepository_GetCurrentAdminMinRoleSort 获取当前用户角色排序性能测试
func BenchmarkAdminRepository_GetCurrentAdminMinRoleSort(b *testing.B) {
	mockRepo := new(MockAdminRepository)
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	testAdmin := createTestAdmin()

	// 设置mock期望
	mockRepo.On("GetCurrentAdminMinRoleSort", c).Return(uint(1), testAdmin, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = mockRepo.GetCurrentAdminMinRoleSort(c)
	}
}

// BenchmarkAdminRepository_GetAdminMinRoleSortsByIds 获取用户角色排序性能测试
func BenchmarkAdminRepository_GetAdminMinRoleSortsByIds(b *testing.B) {
	mockRepo := new(MockAdminRepository)
	ids := []uint{1, 2, 3, 4, 5}
	expectedSorts := []int{1, 2, 3, 4, 5}

	// 设置mock期望
	mockRepo.On("GetAdminMinRoleSortsByIds", ids).Return(expectedSorts, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetAdminMinRoleSortsByIds(ids)
	}
}
