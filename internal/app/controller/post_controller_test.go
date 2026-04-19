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
	"gotribe-admin/pkg/api/known"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MockPostRepository 模拟 PostRepository
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) CreatePost(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepository) GetPostByID(ctx context.Context, id uint) (model.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Post), args.Error(1)
}

func (m *MockPostRepository) GetPosts(ctx context.Context, req *vo.PostListRequest) ([]*model.Post, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*model.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostRepository) UpdatePost(ctx context.Context, post *model.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepository) BatchDeletePostByIds(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

// MockProjectRepository 模拟 ProjectRepository
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) CreateProject(ctx context.Context, project *model.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) GetProjectByID(ctx context.Context, id uint) (model.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Project), args.Error(1)
}

func (m *MockProjectRepository) GetProjects(ctx context.Context, req *vo.ProjectListRequest) ([]*model.Project, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*model.Project), args.Get(1).(int64), args.Error(2)
}

func (m *MockProjectRepository) UpdateProject(ctx context.Context, project *model.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) BatchDeleteProjectByIds(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockProjectRepository) GetProjectsBySitemap(ctx context.Context) ([]*model.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Project), args.Error(1)
}

// init 初始化验证器和日志
func init() {
	gin.SetMode(gin.TestMode)
	// 初始化一个简单的 logger 用于测试
	common.Log = zap.NewNop().Sugar()
	common.InitValidate()
}

// setupPostTest 设置测试环境
func setupPostTest() (*gin.Engine, *MockPostRepository, *MockProjectRepository, PostController) {
	router := gin.New()

	mockPostRepo := new(MockPostRepository)
	mockProjectRepo := new(MockProjectRepository)

	controller := PostController{
		PostRepository:    mockPostRepo,
		ProjectRepository: mockProjectRepo,
	}

	return router, mockPostRepo, mockProjectRepo, controller
}

// createTestPost 创建测试用的 Post 数据
func createTestPost() model.Post {
	return model.Post{
		Model: model.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Slug:        "abc123",
		Title:       "Test Post",
		Description: "Test Description",
		Content:     "Test Content",
		HtmlContent: "<p>Test Content</p>",
		Author:      "testuser",
		UserID:      1,
		CategoryID:  1,
		ProjectId:   1,
		Status:      known.POST_STATUS_DRAFT,
		Type:        known.POST_TYPE_POST,
	}
}

func TestPostController_GetPostInfo_Success(t *testing.T) {
	router, mockPostRepo, _, controller := setupPostTest()

	testPost := createTestPost()
	mockPostRepo.On("GetPostByID", mock.Anything, uint(1)).Return(testPost, nil)

	router.GET("/post/:id", controller.GetPostInfo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/post/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeSuccess, resp.Code)

	data, _ := json.Marshal(resp.Data)
	assert.Contains(t, string(data), "abc123")
	assert.Contains(t, string(data), "Test Post")

	mockPostRepo.AssertExpectations(t)
}

func TestPostController_GetPostInfo_NotFound(t *testing.T) {
	router, mockPostRepo, _, controller := setupPostTest()

	mockPostRepo.On("GetPostByID", mock.Anything, uint(999)).Return(model.Post{}, gorm.ErrRecordNotFound)

	router.GET("/post/:id", controller.GetPostInfo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/post/999", nil)
	router.ServeHTTP(w, req)

	// HandleDatabaseError 将错误映射为 CodeDatabaseError (1009)
	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeDatabaseError, resp.Code)

	mockPostRepo.AssertExpectations(t)
}

func TestPostController_GetPosts_Success(t *testing.T) {
	router, mockPostRepo, _, controller := setupPostTest()

	testPosts := []*model.Post{
		{Model: model.Model{ID: 1}, Slug: "post1", Title: "Post 1"},
		{Model: model.Model{ID: 2}, Slug: "post2", Title: "Post 2"},
	}
	mockPostRepo.On("GetPosts", mock.Anything, mock.AnythingOfType("*vo.PostListRequest")).Return(testPosts, int64(2), nil)

	router.GET("/posts", controller.GetPosts)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts?pageNum=1&pageSize=10", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeSuccess, resp.Code)

	mockPostRepo.AssertExpectations(t)
}

func TestPostController_GetPosts_DatabaseError(t *testing.T) {
	router, mockPostRepo, _, controller := setupPostTest()

	mockPostRepo.On("GetPosts", mock.Anything, mock.AnythingOfType("*vo.PostListRequest")).Return([]*model.Post{}, int64(0), errors.New("database error"))

	router.GET("/posts", controller.GetPosts)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts", nil)
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeDatabaseError, resp.Code)

	mockPostRepo.AssertExpectations(t)
}

func TestPostController_CreatePost_Success(t *testing.T) {
	router, mockPostRepo, _, controller := setupPostTest()

	mockPostRepo.On("CreatePost", mock.Anything, mock.AnythingOfType("*model.Post")).Return(nil)

	router.POST("/post", controller.CreatePost)

	createReq := vo.CreatePostRequest{
		Title:       "New Post",
		Description: "New Description",
		Content:     "New Content",
		HtmlContent: "<p>New Content</p>",
		Author:      "testuser",
		UserID:      1,
		CategoryID:  1,
		ProjectId:   1,
		Type:        known.POST_TYPE_POST,
	}

	jsonData, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeSuccess, resp.Code)

	mockPostRepo.AssertExpectations(t)
}

func TestPostController_CreatePost_ValidationFail(t *testing.T) {
	router, _, _, controller := setupPostTest()

	router.POST("/post", controller.CreatePost)

	// 缺少必填字段
	invalidReq := map[string]interface{}{
		"title": "New Post",
		// 缺少其他必填字段
	}

	jsonData, _ := json.Marshal(invalidReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeValidationFailed, resp.Code)
}

func TestPostController_CreatePost_InvalidJSON(t *testing.T) {
	router, _, _, controller := setupPostTest()

	router.POST("/post", controller.CreatePost)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// HandleBindError 返回 422 (CodeValidationFailed)
	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeValidationFailed, resp.Code)
}

func TestPostController_UpdatePostByID_Success(t *testing.T) {
	router, mockPostRepo, _, controller := setupPostTest()

	testPost := createTestPost()
	mockPostRepo.On("GetPostByID", mock.Anything, uint(1)).Return(testPost, nil)
	mockPostRepo.On("UpdatePost", mock.Anything, mock.AnythingOfType("*model.Post")).Return(nil)

	router.PATCH("/post/:id", controller.UpdatePostByID)

	updateReq := vo.UpdatePostRequest{
		Title:       "Updated Post",
		Description: "Updated Description",
		Content:     "Updated Content",
		HtmlContent: "<p>Updated Content</p>",
		Author:      "testuser",
		UserID:      1,
		CategoryID:  1,
		ProjectId:   1,
		Type:        known.POST_TYPE_POST,
	}

	jsonData, _ := json.Marshal(updateReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/post/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeSuccess, resp.Code)

	mockPostRepo.AssertExpectations(t)
}

func TestPostController_UpdatePostByID_NotFound(t *testing.T) {
	router, mockPostRepo, _, controller := setupPostTest()

	mockPostRepo.On("GetPostByID", mock.Anything, uint(999)).Return(model.Post{}, gorm.ErrRecordNotFound)

	router.PATCH("/post/:id", controller.UpdatePostByID)

	updateReq := vo.UpdatePostRequest{
		Title:       "Updated Post",
		Description: "Updated Description",
		Content:     "Updated Content",
		HtmlContent: "<p>Updated Content</p>",
		Author:      "testuser",
		UserID:      1,
		CategoryID:  1,
		ProjectId:   1,
		Type:        known.POST_TYPE_POST,
	}

	jsonData, _ := json.Marshal(updateReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/post/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// HandleDatabaseError 返回 CodeDatabaseError (1009)
	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeDatabaseError, resp.Code)

	mockPostRepo.AssertExpectations(t)
}

func TestPostController_UpdatePostByID_ValidationFail(t *testing.T) {
	router, _, _, controller := setupPostTest()

	router.PATCH("/post/:id", controller.UpdatePostByID)

	// 缺少必填字段
	invalidReq := map[string]interface{}{
		"title": "Updated Post",
		// 缺少其他必填字段
	}

	jsonData, _ := json.Marshal(invalidReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/post/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeValidationFailed, resp.Code)
}

func TestPostController_BatchDeletePostByIds_Success(t *testing.T) {
	router, mockPostRepo, _, controller := setupPostTest()

	mockPostRepo.On("BatchDeletePostByIds", mock.Anything, []uint{1, 2}).Return(nil)

	router.DELETE("/posts", controller.BatchDeletePostByIds)

	deleteReq := vo.DeletePostsRequest{
		PostIds: "1,2",
	}

	jsonData, _ := json.Marshal(deleteReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeSuccess, resp.Code)

	mockPostRepo.AssertExpectations(t)
}

func TestPostController_BatchDeletePostByIds_EmptyIds(t *testing.T) {
	router, mockPostRepo, _, controller := setupPostTest()

	mockPostRepo.On("BatchDeletePostByIds", mock.Anything, []uint(nil)).Return(nil)

	router.DELETE("/posts", controller.BatchDeletePostByIds)

	deleteReq := vo.DeletePostsRequest{
		PostIds: "",
	}

	jsonData, _ := json.Marshal(deleteReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostController_BatchDeletePostByIds_DatabaseError(t *testing.T) {
	router, mockPostRepo, _, controller := setupPostTest()

	mockPostRepo.On("BatchDeletePostByIds", mock.Anything, []uint{1}).Return(errors.New("database error"))

	router.DELETE("/posts", controller.BatchDeletePostByIds)

	deleteReq := vo.DeletePostsRequest{
		PostIds: "1",
	}

	jsonData, _ := json.Marshal(deleteReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeDatabaseError, resp.Code)

	mockPostRepo.AssertExpectations(t)
}

func TestPostController_PushPostByID_Success(t *testing.T) {
	router, mockPostRepo, mockProjectRepo, controller := setupPostTest()

	testPost := createTestPost()
	mockPostRepo.On("GetPostByID", mock.Anything, uint(1)).Return(testPost, nil)
	mockPostRepo.On("UpdatePost", mock.Anything, mock.AnythingOfType("*model.Post")).Return(nil)
	mockProjectRepo.On("GetProjectByID", mock.Anything, uint(1)).Return(model.Project{
		Name:  "proj123",
		Title: "Test Project",
	}, nil)

	router.PUT("/post/:id", controller.PushPostByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/post/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeSuccess, resp.Code)

	mockPostRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
}

func TestPostController_PushPostByID_NotFound(t *testing.T) {
	router, mockPostRepo, _, controller := setupPostTest()

	mockPostRepo.On("GetPostByID", mock.Anything, uint(999)).Return(model.Post{}, gorm.ErrRecordNotFound)

	router.PUT("/post/:id", controller.PushPostByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/post/999", nil)
	router.ServeHTTP(w, req)

	// HandleDatabaseError 返回 CodeDatabaseError (1009)
	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeDatabaseError, resp.Code)

	mockPostRepo.AssertExpectations(t)
}

func TestPostController_PushPostByID_UpdateError(t *testing.T) {
	router, mockPostRepo, _, controller := setupPostTest()

	testPost := createTestPost()
	mockPostRepo.On("GetPostByID", mock.Anything, uint(1)).Return(testPost, nil)
	mockPostRepo.On("UpdatePost", mock.Anything, mock.AnythingOfType("*model.Post")).Return(errors.New("update error"))

	router.PUT("/post/:id", controller.PushPostByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/post/1", nil)
	router.ServeHTTP(w, req)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, response.CodeDatabaseError, resp.Code)

	mockPostRepo.AssertExpectations(t)
}

// 测试 NewPostController 构造函数
func TestNewPostController(t *testing.T) {
	// 由于 NewPostController 使用真实的 Repository，
	// 我们只需要验证它返回了非 nil 的接口
	controller := NewPostController()
	assert.NotNil(t, controller)
}
