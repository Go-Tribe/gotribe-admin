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

// 获取当前项目信息
func (pc ProjectController) GetProjectInfo(c *gin.Context) {
	project, err := pc.ProjectRepository.GetProjectByProjectID(c.Param("projectID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	projectInfoDto := dto.ToProjectInfoDto(&project)
	response.Success(c, gin.H{
		"project": projectInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// 获取项目列表
func (pc ProjectController) GetProjects(c *gin.Context) {
	var req vo.ProjectListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.Fail(c, nil, errStr)
		return
	}

	// 获取
	project, total, err := pc.ProjectRepository.GetProjects(&req)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"projects": dto.ToProjectsDto(project), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// 创建项目
func (pc ProjectController) CreateProject(c *gin.Context) {
	var req vo.CreateProjectRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.Fail(c, nil, errStr)
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

	err := pc.ProjectRepository.CreateProject(&project)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgCreateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgCreateSuccess))

}

// 更新项目
func (pc ProjectController) UpdateProjectByID(c *gin.Context) {
	var req vo.CreateProjectRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.Fail(c, nil, errStr)
		return
	}

	// 根据path中的ProjectID获取项目信息
	oldProject, err := pc.ProjectRepository.GetProjectByProjectID(c.Param("projectID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	oldProject.Title = req.Title
	oldProject.Description = req.Description
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
	err = pc.ProjectRepository.UpdateProject(&oldProject)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgUpdateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// 批量删除
func (tc ProjectController) BatchDeleteProjectByIds(c *gin.Context) {
	var req vo.DeleteProjectsRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.Fail(c, nil, errStr)
		return
	}

	// 前端传来的标签ID
	reqProjectIds := strings.Split(req.ProjectIds, ",")
	err := tc.ProjectRepository.BatchDeleteProjectByIds(reqProjectIds)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgDeleteFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))

}
