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

	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockAdminRepository 是 IAdminRepository 的 mock 实现
type MockAdminRepository struct {
	mock.Mock
}

func (m *MockAdminRepository) Login(ctx context.Context, admin *model.Admin) (*model.Admin, error) {
	args := m.Called(ctx, admin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Admin), args.Error(1)
}

func (m *MockAdminRepository) ChangePwd(ctx context.Context, username string, newPasswd string) error {
	args := m.Called(ctx, username, newPasswd)
	return args.Error(0)
}

func (m *MockAdminRepository) CreateAdmin(ctx context.Context, admin *model.Admin) error {
	args := m.Called(ctx, admin)
	return args.Error(0)
}

func (m *MockAdminRepository) GetAdminByID(ctx context.Context, id uint) (model.Admin, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Admin), args.Error(1)
}

func (m *MockAdminRepository) GetAdmins(ctx context.Context, req *vo.AdminListRequest) ([]*model.Admin, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Admin), args.Get(1).(int64), args.Error(2)
}

func (m *MockAdminRepository) UpdateAdmin(ctx context.Context, admin *model.Admin) error {
	args := m.Called(ctx, admin)
	return args.Error(0)
}

func (m *MockAdminRepository) BatchDeleteAdminByIds(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockAdminRepository) GetCurrentAdmin(ctx context.Context, c *gin.Context) (model.Admin, error) {
	args := m.Called(ctx, c)
	return args.Get(0).(model.Admin), args.Error(1)
}

func (m *MockAdminRepository) GetCurrentAdminMinRoleSort(ctx context.Context, c *gin.Context) (uint, model.Admin, error) {
	args := m.Called(ctx, c)
	return args.Get(0).(uint), args.Get(1).(model.Admin), args.Error(2)
}

func (m *MockAdminRepository) GetAdminMinRoleSortsByIds(ctx context.Context, ids []uint) ([]int, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockAdminRepository) SetAdminInfoCache(username string, admin model.Admin) {
	m.Called(username, admin)
}

func (m *MockAdminRepository) UpdateAdminInfoCacheByRoleID(ctx context.Context, roleID uint) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

func (m *MockAdminRepository) ClearAdminInfoCache() {
	m.Called()
}

// 创建测试用的管理员对象
func createTestAdmin() model.Admin {
	nickname := "TestAdmin"
	introduction := "Test Introduction"
	return model.Admin{
		Model: model.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username:     "testadmin",
		Password:     "$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEJF6kN7j8xQhFpO2O", // hashed "password"
		Mobile:       "13800138000",
		Avatar:       "https://example.com/avatar.png",
		Nickname:     &nickname,
		Introduction: &introduction,
		Status:       1,
		Creator:      "system",
		Roles: []*model.Role{
			{
				Model:   model.Model{ID: 1},
				Name:    "管理员",
				Keyword: "admin",
				Sort:    1,
				Status:  1,
			},
		},
	}
}

// 初始化测试用的验证器
func setupTestValidator() {
	if common.Validate == nil {
		common.Validate = validator.New()
		// 注册手机号验证函数 - 简单验证11位数字
		_ = common.Validate.RegisterValidation("checkMobile", func(fl validator.FieldLevel) bool {
			phone := fl.Field().String()
			return len(phone) == 11
		})
	}
}

func TestAdminController_GetAdminInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestValidator()

	tests := []struct {
		name         string
		mockSetup    func(*MockAdminRepository)
		wantHttpCode int
		wantBizCode  int
	}{
		{
			name: "成功获取管理员信息",
			mockSetup: func(m *MockAdminRepository) {
				testAdmin := createTestAdmin()
				m.On("GetCurrentAdmin", mock.Anything, mock.Anything).Return(testAdmin, nil)
			},
			wantHttpCode: http.StatusOK,
			wantBizCode:  response.CodeSuccess,
		},
		{
			name: "获取管理员信息失败",
			mockSetup: func(m *MockAdminRepository) {
				m.On("GetCurrentAdmin", mock.Anything, mock.Anything).Return(model.Admin{}, errors.New("database error"))
			},
			wantHttpCode: http.StatusInternalServerError,
			wantBizCode:  response.CodeDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAdminRepository)
			tt.mockSetup(mockRepo)

			controller := AdminController{AdminRepository: mockRepo}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/admin/info", nil)

			controller.GetAdminInfo(c)

			assert.Equal(t, tt.wantHttpCode, w.Code)
			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBizCode, resp.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAdminController_GetAdmins(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestValidator()

	tests := []struct {
		name         string
		query        string
		mockSetup    func(*MockAdminRepository)
		wantHttpCode int
		wantBizCode  int
	}{
		{
			name:  "成功获取管理员列表",
			query: "?pageNum=1&pageSize=10",
			mockSetup: func(m *MockAdminRepository) {
				admins := []*model.Admin{
					func() *model.Admin {
						a := createTestAdmin()
						return &a
					}(),
				}
				m.On("GetAdmins", mock.Anything, mock.Anything).Return(admins, int64(1), nil)
			},
			wantHttpCode: http.StatusOK,
			wantBizCode:  response.CodeSuccess,
		},
		{
			name:         "参数绑定失败-无效页码",
			query:        "?pageNum=invalid",
			mockSetup:    func(m *MockAdminRepository) {},
			wantHttpCode: http.StatusUnprocessableEntity,
			wantBizCode:  response.CodeValidationFailed,
		},
		{
			name:  "获取列表失败-数据库错误",
			query: "?pageNum=1&pageSize=10",
			mockSetup: func(m *MockAdminRepository) {
				m.On("GetAdmins", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("database error"))
			},
			wantHttpCode: http.StatusInternalServerError,
			wantBizCode:  response.CodeDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAdminRepository)
			tt.mockSetup(mockRepo)

			controller := AdminController{AdminRepository: mockRepo}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/admin/list"+tt.query, nil)

			controller.GetAdmins(c)

			assert.Equal(t, tt.wantHttpCode, w.Code)
			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBizCode, resp.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAdminController_CreateAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestValidator()

	tests := []struct {
		name         string
		body         interface{}
		mockSetup    func(*MockAdminRepository)
		wantHttpCode int
		wantBizCode  int
	}{
		{
			name: "参数绑定失败-错误的数据类型",
			body: map[string]interface{}{
				"username": 123,
			},
			mockSetup:    func(m *MockAdminRepository) {},
			wantHttpCode: http.StatusUnprocessableEntity,
			wantBizCode:  response.CodeValidationFailed,
		},
		{
			name: "参数校验失败-必填字段缺失",
			body: map[string]interface{}{
				"username": "test",
			},
			mockSetup:    func(m *MockAdminRepository) {},
			wantHttpCode: http.StatusUnprocessableEntity,
			wantBizCode:  response.CodeValidationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAdminRepository)
			tt.mockSetup(mockRepo)

			controller := AdminController{AdminRepository: mockRepo}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonBody, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest("POST", "/api/admin/create", bytes.NewBuffer(jsonBody))
			c.Request.Header.Set("Content-Type", "application/json")

			controller.CreateAdmin(c)

			assert.Equal(t, tt.wantHttpCode, w.Code)
			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBizCode, resp.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAdminController_UpdateAdminByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestValidator()

	tests := []struct {
		name         string
		userID       string
		body         interface{}
		mockSetup    func(*MockAdminRepository)
		wantHttpCode int
		wantBizCode  int
	}{
		{
			name:   "用户ID格式错误",
			userID: "invalid",
			body: map[string]interface{}{
				"username": "test",
			},
			mockSetup:    func(m *MockAdminRepository) {},
			wantHttpCode: http.StatusUnprocessableEntity,
			wantBizCode:  response.CodeValidationFailed,
		},
		{
			name:   "用户不存在",
			userID: "999",
			body: map[string]interface{}{
				"username":     "updatedadmin",
				"mobile":       "13900139000",
				"nickname":     "Updated Admin",
				"introduction": "Updated Introduction",
				"status":       1,
				"roleIds":      []uint{2},
			},
			mockSetup: func(m *MockAdminRepository) {
				m.On("GetAdminByID", mock.Anything, uint(999)).Return(model.Admin{}, gorm.ErrRecordNotFound)
			},
			wantHttpCode: http.StatusInternalServerError,
			wantBizCode:  response.CodeDatabaseError,
		},
		// 注意：成功更新管理员和获取当前用户失败的测试需要RoleRepository
		// 由于Controller内部直接实例化RoleRepository，这些场景难以在单元测试中覆盖
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAdminRepository)
			tt.mockSetup(mockRepo)

			controller := AdminController{AdminRepository: mockRepo}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "userID", Value: tt.userID}}

			jsonBody, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest("PATCH", "/api/admin/update/"+tt.userID, bytes.NewBuffer(jsonBody))
			c.Request.Header.Set("Content-Type", "application/json")

			controller.UpdateAdminByID(c)

			assert.Equal(t, tt.wantHttpCode, w.Code)
			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBizCode, resp.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAdminController_BatchDeleteAdminByIds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestValidator()

	tests := []struct {
		name         string
		body         interface{}
		mockSetup    func(*MockAdminRepository)
		wantHttpCode int
		wantBizCode  int
	}{
		{
			name: "成功批量删除管理员",
			body: vo.DeleteAdminRequest{
				UserIds: []uint{2, 3},
			},
			mockSetup: func(m *MockAdminRepository) {
				m.On("GetAdminMinRoleSortsByIds", mock.Anything, []uint{uint(2), uint(3)}).Return([]int{2, 3}, nil)
				m.On("GetCurrentAdminMinRoleSort", mock.Anything, mock.Anything).Return(uint(1), createTestAdmin(), nil)
				m.On("BatchDeleteAdminByIds", mock.Anything, []uint{uint(2), uint(3)}).Return(nil)
			},
			wantHttpCode: http.StatusOK,
			wantBizCode:  response.CodeSuccess,
		},
		{
			name: "参数绑定失败-错误的数据类型",
			body: map[string]interface{}{
				"userIds": "invalid",
			},
			mockSetup:    func(m *MockAdminRepository) {},
			wantHttpCode: http.StatusUnprocessableEntity,
			wantBizCode:  response.CodeValidationFailed,
		},
		{
			name: "尝试删除自己-权限不足",
			body: vo.DeleteAdminRequest{
				UserIds: []uint{1}, // 当前用户的ID
			},
			mockSetup: func(m *MockAdminRepository) {
				m.On("GetAdminMinRoleSortsByIds", mock.Anything, []uint{uint(1)}).Return([]int{1}, nil)
				m.On("GetCurrentAdminMinRoleSort", mock.Anything, mock.Anything).Return(uint(1), createTestAdmin(), nil)
			},
			wantHttpCode: http.StatusForbidden,
			wantBizCode:  response.CodeForbidden,
		},
		{
			name: "删除失败-数据库错误",
			body: vo.DeleteAdminRequest{
				UserIds: []uint{2, 3},
			},
			mockSetup: func(m *MockAdminRepository) {
				m.On("GetAdminMinRoleSortsByIds", mock.Anything, []uint{uint(2), uint(3)}).Return([]int{2, 3}, nil)
				m.On("GetCurrentAdminMinRoleSort", mock.Anything, mock.Anything).Return(uint(1), createTestAdmin(), nil)
				m.On("BatchDeleteAdminByIds", mock.Anything, []uint{uint(2), uint(3)}).Return(errors.New("delete failed"))
			},
			wantHttpCode: http.StatusInternalServerError,
			wantBizCode:  response.CodeDatabaseError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAdminRepository)
			tt.mockSetup(mockRepo)

			controller := AdminController{AdminRepository: mockRepo}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonBody, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest("DELETE", "/api/admin/delete", bytes.NewBuffer(jsonBody))
			c.Request.Header.Set("Content-Type", "application/json")

			controller.BatchDeleteAdminByIds(c)

			assert.Equal(t, tt.wantHttpCode, w.Code)
			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBizCode, resp.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAdminController_ChangePwd(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestValidator()

	tests := []struct {
		name         string
		body         interface{}
		mockSetup    func(*MockAdminRepository)
		wantHttpCode int
		wantBizCode  int
	}{
		// 注意：成功修改密码的测试需要使用真实的bcrypt哈希值
		// 由于密码验证涉及加密操作，建议在集成测试中覆盖此场景
		{
			name: "参数绑定失败-错误的数据类型",
			body: map[string]interface{}{
				"oldPassword": 123,
			},
			mockSetup:    func(m *MockAdminRepository) {},
			wantHttpCode: http.StatusUnprocessableEntity,
			wantBizCode:  response.CodeValidationFailed,
		},
		{
			name: "原密码错误",
			body: vo.ChangePwdRequest{
				OldPassword: "wrongpassword",
				NewPassword: "newpassword123",
			},
			mockSetup: func(m *MockAdminRepository) {
				admin := createTestAdmin()
				admin.Password = "$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqQzBZN0UfGNEJF6kN7j8xQhFpO2O"
				m.On("GetCurrentAdmin", mock.Anything, mock.Anything).Return(admin, nil)
			},
			wantHttpCode: http.StatusUnauthorized,
			wantBizCode:  response.CodePasswordIncorrect,
		},
		{
			name: "获取当前用户失败",
			body: vo.ChangePwdRequest{
				OldPassword: "password",
				NewPassword: "newpassword123",
			},
			mockSetup: func(m *MockAdminRepository) {
				m.On("GetCurrentAdmin", mock.Anything, mock.Anything).Return(model.Admin{}, errors.New("not logged in"))
			},
			wantHttpCode: http.StatusUnprocessableEntity,
			wantBizCode:  response.CodeValidationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAdminRepository)
			tt.mockSetup(mockRepo)

			controller := AdminController{AdminRepository: mockRepo}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonBody, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest("PUT", "/api/admin/changePwd", bytes.NewBuffer(jsonBody))
			c.Request.Header.Set("Content-Type", "application/json")

			controller.ChangePwd(c)

			assert.Equal(t, tt.wantHttpCode, w.Code)
			var resp response.Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBizCode, resp.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}
