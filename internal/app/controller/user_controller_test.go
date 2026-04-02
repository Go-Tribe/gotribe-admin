// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gotribe-admin/config"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func init() {
	// 初始化测试配置
	if config.Conf.System == nil {
		config.Conf.System = &config.SystemConfig{
			CDNDomain: "https://cdn.example.com",
		}
	}
}

// MockUserRepository 是 IUserRepository 的 mock 实现
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id uint) (model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) GetUsers(ctx context.Context, req *vo.UserListRequest) ([]*model.User, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) BatchDeleteUserByIds(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockUserRepository) SearchUserByNickname(ctx context.Context, nickname string) ([]*model.User, error) {
	args := m.Called(ctx, nickname)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.User), args.Error(1)
}

// 创建测试用的用户对象
func createTestUser() model.User {
	return model.User{
		Model: model.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username:  "testuser",
		Password:  "$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEJF6kN7j8xQhFpO2O",
		Nickname:  "Test User",
		Email:     "test@example.com",
		Phone:     "13800138000",
		ProjectID: "project1",
		Sex:       "M",
		Status:    1,
		Point:     100.5,
	}
}

func TestUserController_GetUserInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		userID       string
		mockSetup    func(*MockUserRepository)
		wantHttpCode int
		wantBizCode  int
	}{
		{
			name:   "成功获取用户信息",
			userID: "1",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByID", mock.Anything, uint(1)).Return(createTestUser(), nil)
			},
			wantHttpCode: http.StatusOK,
			wantBizCode:  response.CodeSuccess,
		},
		{
			name:         "用户ID格式错误",
			userID:       "invalid",
			mockSetup:    func(m *MockUserRepository) {},
			wantHttpCode: http.StatusUnprocessableEntity,
			wantBizCode:  response.CodeValidationFailed,
		},
		{
			name:   "用户不存在",
			userID: "999",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByID", mock.Anything, uint(999)).Return(model.User{}, gorm.ErrRecordNotFound)
			},
			wantHttpCode: http.StatusInternalServerError,
			wantBizCode:  response.CodeDatabaseError,
		},
		{
			name:   "数据库错误",
			userID: "1",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByID", mock.Anything, uint(1)).Return(model.User{}, errors.New("database error"))
			},
			wantHttpCode: http.StatusInternalServerError,
			wantBizCode:  response.CodeDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			controller := UserController{UserRepository: mockRepo}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: tt.userID}}
			c.Request = httptest.NewRequest("GET", "/api/user/"+tt.userID, nil)

			controller.GetUserInfo(c)

			assert.Equal(t, tt.wantHttpCode, w.Code)
			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBizCode, resp.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserController_GetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		query        string
		mockSetup    func(*MockUserRepository)
		wantHttpCode int
		wantBizCode  int
	}{
		{
			name:  "成功获取用户列表",
			query: "?pageNum=1&pageSize=10",
			mockSetup: func(m *MockUserRepository) {
				users := []*model.User{
					func() *model.User {
						u := createTestUser()
						return &u
					}(),
				}
				m.On("GetUsers", mock.Anything, mock.Anything).Return(users, int64(1), nil)
			},
			wantHttpCode: http.StatusOK,
			wantBizCode:  response.CodeSuccess,
		},
		{
			name:  "带筛选条件获取用户列表",
			query: "?username=test&nickname=Test&userID=1",
			mockSetup: func(m *MockUserRepository) {
				users := []*model.User{
					func() *model.User {
						u := createTestUser()
						return &u
					}(),
				}
				m.On("GetUsers", mock.Anything, mock.Anything).Return(users, int64(1), nil)
			},
			wantHttpCode: http.StatusOK,
			wantBizCode:  response.CodeSuccess,
		},
		{
			name:  "获取列表失败-数据库错误",
			query: "?pageNum=1&pageSize=10",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUsers", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("database error"))
			},
			wantHttpCode: http.StatusInternalServerError,
			wantBizCode:  response.CodeDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			controller := UserController{UserRepository: mockRepo}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/user"+tt.query, nil)

			controller.GetUsers(c)

			assert.Equal(t, tt.wantHttpCode, w.Code)
			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBizCode, resp.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserController_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		body         interface{}
		mockSetup    func(*MockUserRepository)
		wantHttpCode int
		wantBizCode  int
	}{
		{
			name: "成功创建用户",
			body: vo.CreateUserRequest{
				Username:  "newuser",
				Password:  "password123",
				Nickname:  "New User",
				Email:     "new@example.com",
				Phone:     "13900139000",
				ProjectID: "project1",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("CreateUser", mock.Anything, mock.Anything).Return(nil)
			},
			wantHttpCode: http.StatusOK,
			wantBizCode:  response.CodeSuccess,
		},
		{
			name: "参数绑定失败",
			body: map[string]interface{}{
				"username": 123,
			},
			mockSetup:    func(m *MockUserRepository) {},
			wantHttpCode: http.StatusUnprocessableEntity,
			wantBizCode:  response.CodeValidationFailed,
		},
		{
			name: "创建用户失败-数据库错误",
			body: vo.CreateUserRequest{
				Username:  "newuser",
				Password:  "password123",
				Nickname:  "New User",
				Email:     "new@example.com",
				Phone:     "13900139000",
				ProjectID: "project1",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("CreateUser", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			wantHttpCode: http.StatusInternalServerError,
			wantBizCode:  response.CodeDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			controller := UserController{UserRepository: mockRepo}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonBody, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest("POST", "/api/user", bytes.NewBuffer(jsonBody))
			c.Request.Header.Set("Content-Type", "application/json")

			controller.CreateUser(c)

			assert.Equal(t, tt.wantHttpCode, w.Code)
			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBizCode, resp.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserController_UpdateUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		userID       string
		body         interface{}
		mockSetup    func(*MockUserRepository)
		wantHttpCode int
		wantBizCode  int
	}{
		{
			name:   "成功更新用户",
			userID: "1",
			body: vo.UpdateUserRequest{
				Nickname: "Updated User",
				Email:    "updated@example.com",
				Phone:    "13700137000",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByID", mock.Anything, uint(1)).Return(createTestUser(), nil)
				m.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)
			},
			wantHttpCode: http.StatusOK,
			wantBizCode:  response.CodeSuccess,
		},
		{
			name:   "成功更新用户（带密码）",
			userID: "1",
			body: vo.UpdateUserRequest{
				Nickname: "Updated User",
				Password: "newpassword123",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByID", mock.Anything, uint(1)).Return(createTestUser(), nil)
				m.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)
			},
			wantHttpCode: http.StatusOK,
			wantBizCode:  response.CodeSuccess,
		},
		{
			name:   "用户ID格式错误",
			userID: "invalid",
			body: vo.UpdateUserRequest{
				Nickname: "Updated User",
			},
			mockSetup:    func(m *MockUserRepository) {},
			wantHttpCode: http.StatusUnprocessableEntity,
			wantBizCode:  response.CodeValidationFailed,
		},
		{
			name:   "用户不存在",
			userID: "999",
			body: vo.UpdateUserRequest{
				Nickname: "Updated User",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByID", mock.Anything, uint(999)).Return(model.User{}, gorm.ErrRecordNotFound)
			},
			wantHttpCode: http.StatusInternalServerError,
			wantBizCode:  response.CodeDatabaseError,
		},
		{
			name:   "更新失败-数据库错误",
			userID: "1",
			body: vo.UpdateUserRequest{
				Nickname: "Updated User",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByID", mock.Anything, uint(1)).Return(createTestUser(), nil)
				m.On("UpdateUser", mock.Anything, mock.Anything).Return(errors.New("update failed"))
			},
			wantHttpCode: http.StatusInternalServerError,
			wantBizCode:  response.CodeDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			controller := UserController{UserRepository: mockRepo}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: tt.userID}}

			jsonBody, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest("PATCH", "/api/user/"+tt.userID, bytes.NewBuffer(jsonBody))
			c.Request.Header.Set("Content-Type", "application/json")

			controller.UpdateUserByID(c)

			assert.Equal(t, tt.wantHttpCode, w.Code)
			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBizCode, resp.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserController_BatchDeleteUserByIds(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		body         interface{}
		mockSetup    func(*MockUserRepository)
		wantHttpCode int
		wantBizCode  int
	}{
		{
			name: "成功批量删除用户",
			body: vo.DeleteUsersRequest{
				Ids: []uint{1, 2, 3},
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("BatchDeleteUserByIds", mock.Anything, []uint{uint(1), uint(2), uint(3)}).Return(nil)
			},
			wantHttpCode: http.StatusOK,
			wantBizCode:  response.CodeSuccess,
		},
		{
			name: "参数绑定失败-错误的数据类型",
			body: map[string]interface{}{
				"ids": "invalid",
			},
			mockSetup:    func(m *MockUserRepository) {},
			wantHttpCode: http.StatusUnprocessableEntity,
			wantBizCode:  response.CodeValidationFailed,
		},
		{
			name: "删除失败-数据库错误",
			body: vo.DeleteUsersRequest{
				Ids: []uint{1, 2, 3},
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("BatchDeleteUserByIds", mock.Anything, []uint{uint(1), uint(2), uint(3)}).Return(errors.New("delete failed"))
			},
			wantHttpCode: http.StatusInternalServerError,
			wantBizCode:  response.CodeDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			controller := UserController{UserRepository: mockRepo}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonBody, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest("DELETE", "/api/user", bytes.NewBuffer(jsonBody))
			c.Request.Header.Set("Content-Type", "application/json")

			controller.BatchDeleteUserByIds(c)

			assert.Equal(t, tt.wantHttpCode, w.Code)
			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBizCode, resp.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserController_SearchUserByUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		query        string
		mockSetup    func(*MockUserRepository)
		wantHttpCode int
		wantBizCode  int
	}{
		{
			name:  "成功搜索用户",
			query: "?nickname=test",
			mockSetup: func(m *MockUserRepository) {
				users := []*model.User{
					func() *model.User {
						u := createTestUser()
						return &u
					}(),
				}
				m.On("SearchUserByNickname", mock.Anything, "test").Return(users, nil)
			},
			wantHttpCode: http.StatusOK,
			wantBizCode:  response.CodeSuccess,
		},
		{
			name:  "搜索用户失败-数据库错误",
			query: "?nickname=test",
			mockSetup: func(m *MockUserRepository) {
				m.On("SearchUserByNickname", mock.Anything, "test").Return(nil, errors.New("search failed"))
			},
			wantHttpCode: http.StatusInternalServerError,
			wantBizCode:  response.CodeInternalServerError,
		},
		{
			name:  "空昵称搜索",
			query: "?nickname=",
			mockSetup: func(m *MockUserRepository) {
				m.On("SearchUserByNickname", mock.Anything, "").Return([]*model.User{}, nil)
			},
			wantHttpCode: http.StatusOK,
			wantBizCode:  response.CodeSuccess,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			controller := UserController{UserRepository: mockRepo}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/user/search"+tt.query, nil)

			controller.SearchUserByUsername(c)

			assert.Equal(t, tt.wantHttpCode, w.Code)
			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBizCode, resp.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}
