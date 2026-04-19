// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/dto"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IProjectController interface {
	GetProjectInfo(c *gin.Context)          // 获取当前登录项目信息
	GetProjects(c *gin.Context)             // 获取项目列表
	CreateProject(c *gin.Context)           // 创建项目
	UpdateProjectByID(c *gin.Context)       // 更新项目
	BatchDeleteProjectByIds(c *gin.Context) // 批量删除项目
}

type ProjectController struct {
	ProjectRepository repository.IProjectRepository
}

// 构造函数
func NewProjectController() IProjectController {
	projectRepository := repository.NewProjectRepository()
	projectController := ProjectController{ProjectRepository: projectRepository}
	return projectController
}

// GetProjectInfo 获取当前项目信息
// @Summary      获取项目信息
// @Description  根据项目ID获取项目详细信息
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        id path uint true "项目ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /project/{id} [get]
// @Security     BearerAuth
func (pc ProjectController) GetProjectInfo(c *gin.Context) {
	ctx := c.Request.Context()
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationFail(c, "无效的 id")
		return
	}
	project, err := pc.ProjectRepository.GetProjectByID(ctx, uint(projectID))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	projectInfoDto := dto.ToProjectInfoDto(&project)
	response.Success(c, gin.H{
		"project": projectInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// GetProjects 获取项目列表
// @Summary      获取项目列表
// @Description  获取所有项目的列表，支持分页和筛选
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        request query vo.ProjectListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /project [get]
// @Security     BearerAuth
func (pc ProjectController) GetProjects(c *gin.Context) {
	ctx := c.Request.Context()
	var req vo.ProjectListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.ValidationFail(c, errStr)
		return
	}

	// 获取
	project, total, err := pc.ProjectRepository.GetProjects(ctx, &req)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgListFail)
		return
	}
	response.Success(c, gin.H{"projects": dto.ToProjectsDto(project), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// CreateProject 创建项目
// @Summary      创建项目
// @Description  创建一个新的项目
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        request body vo.CreateProjectRequest true "创建项目请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /project [post]
// @Security     BearerAuth
func (pc ProjectController) CreateProject(c *gin.Context) {
	ctx := c.Request.Context()
	var req vo.CreateProjectRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.ValidationFail(c, errStr)
		return
	}

	project := model.Project{
		Name:           req.Name,
		Title:          req.Title,
		Description:    req.Description,
		Keywords:       req.Keywords,
		Domain:         req.Domain,
		PostURL:        req.PostURL,
		ICP:            req.ICP,
		Author:         req.Author,
		Info:           req.Info,
		PublicSecurity: req.PublicSecurity,
		Favicon:        req.Favicon,
		NavImage:       req.NavImage,
		BaiduAnalytics: req.BaiduAnalytics,
		PushToken:      req.PushToken,
	}

	err := pc.ProjectRepository.CreateProject(ctx, &project)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgCreateFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgCreateSuccess))

}

// UpdateProjectByID 更新项目
// @Summary      更新项目
// @Description  根据项目ID更新项目信息
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        id path uint true "项目ID"
// @Param        request body vo.CreateProjectRequest true "更新项目请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /project/{id} [patch]
// @Security     BearerAuth
func (pc ProjectController) UpdateProjectByID(c *gin.Context) {
	ctx := c.Request.Context()
	var req vo.CreateProjectRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.ValidationFail(c, errStr)
		return
	}

	// 根据path中的ProjectID获取项目信息
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationFail(c, "无效的 id")
		return
	}
	oldProject, err := pc.ProjectRepository.GetProjectByID(ctx, uint(projectID))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	oldProject.Title = req.Title
	oldProject.Description = req.Description
	oldProject.Name = req.Name
	oldProject.Author = req.Author
	oldProject.ICP = req.ICP
	oldProject.Keywords = req.Keywords
	oldProject.Info = req.Info
	oldProject.PostURL = req.PostURL
	oldProject.Domain = req.Domain
	oldProject.PublicSecurity = req.PublicSecurity
	oldProject.Favicon = req.Favicon
	oldProject.PostURL = req.PostURL
	oldProject.NavImage = req.NavImage
	oldProject.BaiduAnalytics = req.BaiduAnalytics
	oldProject.PushToken = req.PushToken
	// 更新项目
	err = pc.ProjectRepository.UpdateProject(ctx, &oldProject)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgUpdateFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// BatchDeleteProjectByIds 批量删除项目
// @Summary      批量删除项目
// @Description  根据项目ID列表批量删除项目
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        request body vo.DeleteProjectsRequest true "删除项目请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /project [delete]
// @Security     BearerAuth
func (tc ProjectController) BatchDeleteProjectByIds(c *gin.Context) {
	ctx := c.Request.Context()
	var req vo.DeleteProjectsRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.ValidationFail(c, errStr)
		return
	}
	// 前端传来的项目ID
	reqProjectIds := strings.Split(req.ProjectIds, ",")
	var ids []uint
	for _, idStr := range reqProjectIds {
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			response.ValidationFail(c, "项目ID格式错误")
			return
		}
		ids = append(ids, uint(id))
	}
	err := tc.ProjectRepository.BatchDeleteProjectByIds(ctx, ids)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgDeleteFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))
}
