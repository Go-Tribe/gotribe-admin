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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// init 初始化验证器和日志（如果尚未初始化）
func init() {
	gin.SetMode(gin.TestMode)
	if common.Log == nil {
		common.Log = zap.NewNop().Sugar()
	}
	if common.Validate == nil {
		common.InitValidate()
	}
}

// MockCategoryRepository 模拟 CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) GetCategoryByID(ctx context.Context, id uint) (model.Category, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetCategorys(ctx context.Context) ([]*model.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetCategoryTree(ctx context.Context) ([]*model.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Category), args.Error(1)
}

func (m *MockCategoryRepository) CreateCategory(ctx context.Context, category *model.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) UpdateCategoryByID(ctx context.Context, id uint, category *model.Category) error {
	args := m.Called(ctx, id, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) BatchDeleteCategoryByIds(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

// init 初始化验证器
func init() {
	gin.SetMode(gin.TestMode)
	// 验证器在 post_controller_test.go 中已经初始化
	// 这里不需要重复初始化
}

// setupCategoryTest 设置测试环境
func setupCategoryTest() (*gin.Engine, *MockCategoryRepository, CategoryController) {
	router := gin.New()

	mockCategoryRepo := new(MockCategoryRepository)

	controller := CategoryController{
		CategoryRepository: mockCategoryRepo,
	}

	return router, mockCategoryRepo, controller
}

// createTestCategory 创建测试用的 Category 数据
func createTestCategory() model.Category {
	return model.Category{
		Model: model.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Title:       "Test Category",
		Slug:        "test-category",
		Description: "Test Description",
		Icon:        "icon.png",
		Path:        "/test",
		Sort:        1,
		Status:      1,
		Hidden:      1,
		ParentID:    0,
	}
}

func TestCategoryController_GetCategoryInfo_Success(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	testCategory := createTestCategory()
	mockCategoryRepo.On("GetCategoryByID", mock.Anything, uint(1)).Return(testCategory, nil)

	router.GET("/category/:id", controller.GetCategoryInfo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/category/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeSuccess, resp.Code)

	data, _ := json.Marshal(resp.Data)
	assert.Contains(t, string(data), "Test Category")
	assert.Contains(t, string(data), "test-category")

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryController_GetCategoryInfo_NotFound(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	mockCategoryRepo.On("GetCategoryByID", mock.Anything, uint(999)).Return(model.Category{}, gorm.ErrRecordNotFound)

	router.GET("/category/:id", controller.GetCategoryInfo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/category/999", nil)
	router.ServeHTTP(w, req)

	// HandleDatabaseError 返回 CodeDatabaseError (1009)
	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeDatabaseError, resp.Code)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryController_GetCategoryInfo_InvalidID(t *testing.T) {
	router, _, controller := setupCategoryTest()

	router.GET("/category/:id", controller.GetCategoryInfo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/category/invalid", nil)
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeValidationFailed, resp.Code)
}

func TestCategoryController_GetCategoryInfo_DatabaseError(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	mockCategoryRepo.On("GetCategoryByID", mock.Anything, uint(1)).Return(model.Category{}, errors.New("database error"))

	router.GET("/category/:id", controller.GetCategoryInfo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/category/1", nil)
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeDatabaseError, resp.Code)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryController_GetCategorys_Success(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	testCategories := []*model.Category{
		{Model: model.Model{ID: 1}, Title: "Category 1", Slug: "category-1"},
		{Model: model.Model{ID: 2}, Title: "Category 2", Slug: "category-2"},
	}
	mockCategoryRepo.On("GetCategorys", mock.Anything).Return(testCategories, nil)

	router.GET("/categories", controller.GetCategorys)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeSuccess, resp.Code)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryController_GetCategorys_DatabaseError(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	mockCategoryRepo.On("GetCategorys", mock.Anything).Return([]*model.Category{}, errors.New("database error"))

	router.GET("/categories", controller.GetCategorys)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories", nil)
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeDatabaseError, resp.Code)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryController_GetCategoryTree_Success(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	testTree := []*model.Category{
		{
			Model:    model.Model{ID: 1},
			Title:    "Parent Category",
			Slug:     "parent",
			ParentID: 0,
			Children: []*model.Category{
				{Model: model.Model{ID: 2}, Title: "Child Category", Slug: "child", ParentID: 1},
			},
		},
	}
	mockCategoryRepo.On("GetCategoryTree", mock.Anything).Return(testTree, nil)

	router.GET("/category/tree", controller.GetCategoryTree)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/category/tree", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeSuccess, resp.Code)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryController_GetCategoryTree_DatabaseError(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	mockCategoryRepo.On("GetCategoryTree", mock.Anything).Return([]*model.Category{}, errors.New("database error"))

	router.GET("/category/tree", controller.GetCategoryTree)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/category/tree", nil)
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeDatabaseError, resp.Code)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryController_CreateCategory_Success(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	mockCategoryRepo.On("CreateCategory", mock.Anything, mock.AnythingOfType("*model.Category")).Return(nil)

	router.POST("/category", controller.CreateCategory)

	createReq := vo.CreateCategoryRequest{
		Title:       "New Category",
		Slug:        "new-category",
		Icon:        "icon.png",
		Path:        "/new",
		Sort:        1,
		Status:      1,
		Hidden:      1,
		Description: "New Description",
		ParentID:    0,
	}

	jsonData, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/category", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeSuccess, resp.Code)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryController_CreateCategory_ValidationFail(t *testing.T) {
	router, _, controller := setupCategoryTest()

	router.POST("/category", controller.CreateCategory)

	// 缺少必填字段 title 和 slug
	invalidReq := map[string]interface{}{
		"icon": "icon.png",
	}

	jsonData, _ := json.Marshal(invalidReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/category", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeValidationFailed, resp.Code)
}

func TestCategoryController_CreateCategory_InvalidJSON(t *testing.T) {
	router, _, controller := setupCategoryTest()

	router.POST("/category", controller.CreateCategory)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/category", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// HandleBindError 返回 422 (CodeValidationFailed)
	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeValidationFailed, resp.Code)
}

func TestCategoryController_CreateCategory_InvalidSort(t *testing.T) {
	router, _, controller := setupCategoryTest()

	router.POST("/category", controller.CreateCategory)

	// Sort 超出范围 (1-999)
	invalidReq := vo.CreateCategoryRequest{
		Title: "Test",
		Slug:  "test",
		Sort:  1000,
	}

	jsonData, _ := json.Marshal(invalidReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/category", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeValidationFailed, resp.Code)
}

func TestCategoryController_CreateCategory_DatabaseError(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	mockCategoryRepo.On("CreateCategory", mock.Anything, mock.AnythingOfType("*model.Category")).Return(errors.New("database error"))

	router.POST("/category", controller.CreateCategory)

	// 需要提供所有必填字段
	createReq := vo.CreateCategoryRequest{
		Title:  "New Category",
		Slug:   "new-category",
		Sort:   1,
		Hidden: 1,
	}

	jsonData, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/category", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeDatabaseError, resp.Code)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryController_UpdateCategoryByID_Success(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	testCategory := createTestCategory()
	mockCategoryRepo.On("GetCategoryByID", mock.Anything, uint(1)).Return(testCategory, nil)
	mockCategoryRepo.On("UpdateCategoryByID", mock.Anything, uint(1), mock.AnythingOfType("*model.Category")).Return(nil)

	router.PATCH("/category/:id", controller.UpdateCategoryByID)

	updateReq := vo.UpdateCategoryRequest{
		Title:       "Updated Category",
		Slug:        "updated-category",
		Icon:        "updated-icon.png",
		Path:        "/updated",
		Sort:        2,
		Status:      1,
		Hidden:      1,
		Description: "Updated Description",
		ParentID:    0,
	}

	jsonData, _ := json.Marshal(updateReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/category/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeSuccess, resp.Code)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryController_UpdateCategoryByID_NotFound(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	mockCategoryRepo.On("GetCategoryByID", mock.Anything, uint(999)).Return(model.Category{}, gorm.ErrRecordNotFound)

	router.PATCH("/category/:id", controller.UpdateCategoryByID)

	// 需要提供所有必填字段，包括 Hidden (oneof=1 2)
	updateReq := vo.UpdateCategoryRequest{
		Title:    "Updated Category",
		Slug:     "updated-category",
		Sort:     1,
		Hidden:   1,
		ParentID: 0,
	}

	jsonData, _ := json.Marshal(updateReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/category/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// HandleDatabaseError 返回 CodeDatabaseError (1009)
	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeDatabaseError, resp.Code)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryController_UpdateCategoryByID_InvalidID(t *testing.T) {
	router, _, controller := setupCategoryTest()

	router.PATCH("/category/:id", controller.UpdateCategoryByID)

	// 需要提供所有必填字段，包括 Hidden (oneof=1 2)
	updateReq := vo.UpdateCategoryRequest{
		Title:    "Updated Category",
		Slug:     "updated-category",
		Sort:     1,
		Hidden:   1,
		ParentID: 0,
	}

	jsonData, _ := json.Marshal(updateReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/category/invalid", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeValidationFailed, resp.Code)
}

func TestCategoryController_UpdateCategoryByID_SelfParent(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	testCategory := createTestCategory()
	testCategory.ID = 5
	mockCategoryRepo.On("GetCategoryByID", mock.Anything, uint(5)).Return(testCategory, nil)

	router.PATCH("/category/:id", controller.UpdateCategoryByID)

	// 尝试将自己设为父分类
	updateReq := vo.UpdateCategoryRequest{
		Title:    "Updated Category",
		Slug:     "updated-category",
		Sort:     1,
		ParentID: 5, // 将自己设为父分类
	}

	jsonData, _ := json.Marshal(updateReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/category/5", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeValidationFailed, resp.Code)
}

func TestCategoryController_UpdateCategoryByID_ValidationFail(t *testing.T) {
	router, _, controller := setupCategoryTest()

	router.PATCH("/category/:id", controller.UpdateCategoryByID)

	// 缺少必填字段
	invalidReq := map[string]interface{}{
		"title": "Updated Category",
		// 缺少 slug
	}

	jsonData, _ := json.Marshal(invalidReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/category/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeValidationFailed, resp.Code)
}

func TestCategoryController_UpdateCategoryByID_DatabaseError(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	testCategory := createTestCategory()
	mockCategoryRepo.On("GetCategoryByID", mock.Anything, uint(1)).Return(testCategory, nil)
	mockCategoryRepo.On("UpdateCategoryByID", mock.Anything, uint(1), mock.AnythingOfType("*model.Category")).Return(errors.New("database error"))

	router.PATCH("/category/:id", controller.UpdateCategoryByID)

	// 需要提供所有必填字段，包括 Hidden (oneof=1 2)
	updateReq := vo.UpdateCategoryRequest{
		Title:    "Updated Category",
		Slug:     "updated-category",
		Sort:     1,
		Hidden:   1,
		ParentID: 0,
	}

	jsonData, _ := json.Marshal(updateReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/category/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeDatabaseError, resp.Code)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryController_BatchDeleteCategoryByIds_Success(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	mockCategoryRepo.On("BatchDeleteCategoryByIds", mock.Anything, []uint{1, 2, 3}).Return(nil)

	router.DELETE("/categories", controller.BatchDeleteCategoryByIds)

	deleteReq := vo.DeleteCategoryRequest{
		Ids: []uint{1, 2, 3},
	}

	jsonData, _ := json.Marshal(deleteReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/categories", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeSuccess, resp.Code)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryController_BatchDeleteCategoryByIds_DatabaseError(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	mockCategoryRepo.On("BatchDeleteCategoryByIds", mock.Anything, []uint{1}).Return(errors.New("database error"))

	router.DELETE("/categories", controller.BatchDeleteCategoryByIds)

	deleteReq := vo.DeleteCategoryRequest{
		Ids: []uint{1},
	}

	jsonData, _ := json.Marshal(deleteReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/categories", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeDatabaseError, resp.Code)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryController_BatchDeleteCategoryByIds_EmptyIds(t *testing.T) {
	router, mockCategoryRepo, controller := setupCategoryTest()

	mockCategoryRepo.On("BatchDeleteCategoryByIds", mock.Anything, []uint{}).Return(nil)

	router.DELETE("/categories", controller.BatchDeleteCategoryByIds)

	deleteReq := vo.DeleteCategoryRequest{
		Ids: []uint{},
	}

	jsonData, _ := json.Marshal(deleteReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/categories", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockCategoryRepo.AssertExpectations(t)
}

func TestCategoryController_BatchDeleteCategoryByIds_InvalidJSON(t *testing.T) {
	router, _, controller := setupCategoryTest()

	router.DELETE("/categories", controller.BatchDeleteCategoryByIds)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/categories", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// HandleBindError 返回 422 (CodeValidationFailed)
	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeValidationFailed, resp.Code)
}

// 测试 NewCategoryController 构造函数
func TestNewCategoryController(t *testing.T) {
	// 由于 NewCategoryController 使用真实的 Repository，
	// 我们只需要验证它返回了非 nil 的接口
	controller := NewCategoryController()
	assert.NotNil(t, controller)
}
