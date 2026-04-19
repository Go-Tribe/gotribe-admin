// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository 模拟UserRepository
type MockUserRepository struct {
	mock.Mock
}

// CreateUser 模拟创建用户方法
func (m *MockUserRepository) CreateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// GetUserByID 模拟根据ID获取用户方法
func (m *MockUserRepository) GetUserByID(ctx context.Context, id uint) (model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.User), args.Error(1)
}

// GetUsers 模拟获取用户列表方法
func (m *MockUserRepository) GetUsers(ctx context.Context, req *vo.UserListRequest) ([]*model.User, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

// UpdateUser 模拟更新用户方法
func (m *MockUserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// BatchDeleteUserByIds 模拟批量删除用户方法
func (m *MockUserRepository) BatchDeleteUserByIds(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

// SearchUserByNickname 模拟根据昵称搜索用户方法
func (m *MockUserRepository) SearchUserByNickname(ctx context.Context, nickname string) ([]*model.User, error) {
	args := m.Called(ctx, nickname)
	return args.Get(0).([]*model.User), args.Error(1)
}

// createTestUser 创建测试用户数据
func createTestUser() *model.User {
	phone := "13800138000"
	email := "test@example.com"
	return &model.User{
		Model: model.Model{
			ID: 1,
		},
		Username:  "testuser",
		Nickname:  "测试用户",
		AvatarURL: "avatar.jpg",
		Phone:     &phone,
		Email:     &email,
		Status:    1,
		Point:     100.0,
		ProjectId: 1,
	}
}

// TestUserRepository_CreateUser 测试创建用户
func TestUserRepository_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)

	testUser := createTestUser()

	// 设置mock期望 - 成功创建
	mockRepo.On("CreateUser", mock.Anything, testUser).Return(nil)

	// 执行测试
	err := mockRepo.CreateUser(context.Background(), testUser)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserRepository_CreateUser_Error 测试创建用户失败
func TestUserRepository_CreateUser_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)

	testUser := createTestUser()

	// 设置mock期望 - 创建失败
	mockRepo.On("CreateUser", mock.Anything, testUser).Return(assert.AnError)

	// 执行测试
	err := mockRepo.CreateUser(context.Background(), testUser)

	// 验证结果
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserRepository_GetUserByID 测试根据ID获取用户
func TestUserRepository_GetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)

	testUser := createTestUser()

	// 设置mock期望
	mockRepo.On("GetUserByID", mock.Anything, uint(1)).Return(*testUser, nil)

	// 执行测试
	result, err := mockRepo.GetUserByID(context.Background(), 1)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "testuser", result.Username)
	mockRepo.AssertExpectations(t)
}

// TestUserRepository_GetUserByID_NotFound 测试用户不存在
func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	// 设置mock期望 - 用户不存在
	mockRepo.On("GetUserByID", mock.Anything, uint(999)).Return(model.User{}, assert.AnError)

	// 执行测试
	result, err := mockRepo.GetUserByID(context.Background(), 999)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, model.User{}, result)
	mockRepo.AssertExpectations(t)
}

// TestUserRepository_GetUsers 测试获取用户列表
func TestUserRepository_GetUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)

	testUser1 := createTestUser()
	testUser2 := &model.User{
		Model: model.Model{
			ID: 2,
		},
		Username:  "testuser2",
		Nickname:  "测试用户2",
		Status:    1,
		Point:     50.0,
		ProjectId: 1,
	}

	testUsers := []*model.User{testUser1, testUser2}

	tests := []struct {
		name     string
		request  *vo.UserListRequest
		expected int64
	}{
		{
			name: "获取所有用户",
			request: &vo.UserListRequest{
				PageNum:  1,
				PageSize: 10,
			},
			expected: 2,
		},
		{
			name: "按用户名搜索",
			request: &vo.UserListRequest{
				Username: "testuser",
				PageNum:  1,
				PageSize: 10,
			},
			expected: 1,
		},
		{
			name: "按昵称搜索",
			request: &vo.UserListRequest{
				Nickname: "测试",
				PageNum:  1,
				PageSize: 10,
			},
			expected: 2,
		},
		{
			name: "按用户ID搜索",
			request: &vo.UserListRequest{
				UserID:   1,
				PageNum:  1,
				PageSize: 10,
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置mock期望
			mockRepo.On("GetUsers", mock.Anything, tt.request).Return(testUsers, tt.expected, nil)

			// 执行测试
			users, total, err := mockRepo.GetUsers(context.Background(), tt.request)

			// 验证结果
			assert.NoError(t, err)
			assert.NotNil(t, users)
			assert.Equal(t, tt.expected, total)
			assert.Len(t, users, 2)

			// 验证mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserRepository_GetUsers_Error 测试获取用户列表失败
func TestUserRepository_GetUsers_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)

	request := &vo.UserListRequest{
		PageNum:  1,
		PageSize: 10,
	}

	// 设置mock期望
	mockRepo.On("GetUsers", mock.Anything, request).Return([]*model.User{}, int64(0), assert.AnError)

	// 执行测试
	users, total, err := mockRepo.GetUsers(context.Background(), request)

	// 验证结果
	assert.Error(t, err)
	assert.Empty(t, users)
	assert.Equal(t, int64(0), total)
	mockRepo.AssertExpectations(t)
}

// TestUserRepository_UpdateUser 测试更新用户
func TestUserRepository_UpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)

	testUser := createTestUser()
	testUser.Nickname = "更新后的昵称"
	newPhone := "13800138999"
	testUser.Phone = &newPhone

	// 设置mock期望
	mockRepo.On("UpdateUser", mock.Anything, testUser).Return(nil)

	// 执行测试
	err := mockRepo.UpdateUser(context.Background(), testUser)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserRepository_UpdateUser_Error 测试更新用户失败
func TestUserRepository_UpdateUser_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)

	testUser := createTestUser()

	// 设置mock期望 - 更新失败
	mockRepo.On("UpdateUser", mock.Anything, testUser).Return(assert.AnError)

	// 执行测试
	err := mockRepo.UpdateUser(context.Background(), testUser)

	// 验证结果
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserRepository_BatchDeleteUserByIds 测试批量删除用户
func TestUserRepository_BatchDeleteUserByIds(t *testing.T) {
	mockRepo := new(MockUserRepository)

	ids := []uint{1, 2, 3}

	// 设置mock期望
	mockRepo.On("BatchDeleteUserByIds", mock.Anything, ids).Return(nil)

	// 执行测试
	err := mockRepo.BatchDeleteUserByIds(context.Background(), ids)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserRepository_BatchDeleteUserByIds_Error 测试批量删除用户失败
func TestUserRepository_BatchDeleteUserByIds_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)

	ids := []uint{999}

	// 设置mock期望 - 删除失败（用户不存在）
	mockRepo.On("BatchDeleteUserByIds", mock.Anything, ids).Return(assert.AnError)

	// 执行测试
	err := mockRepo.BatchDeleteUserByIds(context.Background(), ids)

	// 验证结果
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserRepository_BatchDeleteUserByIds_Empty 测试批量删除空ID列表
func TestUserRepository_BatchDeleteUserByIds_Empty(t *testing.T) {
	mockRepo := new(MockUserRepository)

	ids := []uint{}

	// 设置mock期望
	mockRepo.On("BatchDeleteUserByIds", mock.Anything, ids).Return(nil)

	// 执行测试
	err := mockRepo.BatchDeleteUserByIds(context.Background(), ids)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserRepository_SearchUserByNickname 测试根据昵称搜索用户
func TestUserRepository_SearchUserByNickname(t *testing.T) {
	mockRepo := new(MockUserRepository)

	testUser1 := createTestUser()
	testUser2 := &model.User{
		Model: model.Model{
			ID: 2,
		},
		Username:  "testuser2",
		Nickname:  "测试用户2",
		Status:    1,
		ProjectId: 1,
	}

	testUsers := []*model.User{testUser1, testUser2}

	tests := []struct {
		name     string
		nickname string
		expected []*model.User
	}{
		{
			name:     "搜索测试用户",
			nickname: "测试",
			expected: testUsers,
		},
		{
			name:     "搜索特定用户",
			nickname: "测试用户2",
			expected: []*model.User{testUser2},
		},
		{
			name:     "搜索不存在用户",
			nickname: "不存在的用户",
			expected: []*model.User{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置mock期望
			mockRepo.On("SearchUserByNickname", mock.Anything, tt.nickname).Return(tt.expected, nil)

			// 执行测试
			users, err := mockRepo.SearchUserByNickname(context.Background(), tt.nickname)

			// 验证结果
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, users)

			// 验证mock调用
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUserRepository_SearchUserByNickname_Error 测试搜索用户失败
func TestUserRepository_SearchUserByNickname_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)

	// 设置mock期望
	mockRepo.On("SearchUserByNickname", mock.Anything, "test").Return([]*model.User{}, assert.AnError)

	// 执行测试
	users, err := mockRepo.SearchUserByNickname(context.Background(), "test")

	// 验证结果
	assert.Error(t, err)
	assert.Empty(t, users)
	mockRepo.AssertExpectations(t)
}

// BenchmarkUserRepository_CreateUser 创建用户性能测试
func BenchmarkUserRepository_CreateUser(b *testing.B) {
	mockRepo := new(MockUserRepository)
	testUser := createTestUser()

	// 设置mock期望
	mockRepo.On("CreateUser", mock.Anything, testUser).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mockRepo.CreateUser(context.Background(), testUser)
	}
}

// BenchmarkUserRepository_GetUserByID 根据ID获取用户性能测试
func BenchmarkUserRepository_GetUserByID(b *testing.B) {
	mockRepo := new(MockUserRepository)
	testUser := createTestUser()

	// 设置mock期望
	mockRepo.On("GetUserByID", mock.Anything, uint(1)).Return(*testUser, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockRepo.GetUserByID(context.Background(), 1)
	}
}

// BenchmarkUserRepository_GetUsers 获取用户列表性能测试
func BenchmarkUserRepository_GetUsers(b *testing.B) {
	mockRepo := new(MockUserRepository)
	testUser := createTestUser()
	testUsers := []*model.User{testUser}
	request := &vo.UserListRequest{
		PageNum:  1,
		PageSize: 10,
	}

	// 设置mock期望
	mockRepo.On("GetUsers", mock.Anything, request).Return(testUsers, int64(1), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = mockRepo.GetUsers(context.Background(), request)
	}
}
