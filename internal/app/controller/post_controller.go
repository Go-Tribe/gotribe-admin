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
	"gotribe-admin/pkg/api/known"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"gotribe-admin/pkg/util"
	"strings"

	"github.com/dengmengmian/ghelper/gconvert"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IPostController interface {
	GetPostInfo(c *gin.Context)          // 获取当前登录内容信息
	GetPosts(c *gin.Context)             // 获取内容列表
	CreatePost(c *gin.Context)           // 创建内容
	UpdatePostByID(c *gin.Context)       // 更新内容
	BatchDeletePostByIds(c *gin.Context) // 批量删除内容
	PushPostByID(c *gin.Context)         // 发布接口
}

type PostController struct {
	PostRepository    repository.IPostRepository
	ProjectRepository repository.IProjectRepository
}

// 构造函数
func NewPostController() IPostController {
	postRepository := repository.NewPostRepository()
	projectRepository := repository.NewProjectRepository()
	postController := PostController{PostRepository: postRepository, ProjectRepository: projectRepository}
	return postController
}

// 获取当前内容信息
func (pc PostController) GetPostInfo(c *gin.Context) {
	post, err := pc.PostRepository.GetPostByPostID(c.Param("postID"))
	if err != nil {
		response.Fail(c, nil, "获取当前内容信息失败: "+err.Error())
		return
	}
	postInfoDto := dto.ToPostInfoDto(&post)
	response.Success(c, gin.H{
		"post": postInfoDto,
	}, "获取当前内容信息成功")
}

// 获取内容列表
func (pc PostController) GetPosts(c *gin.Context) {
	var req vo.PostListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}

	// 获取
	post, total, err := pc.PostRepository.GetPosts(&req)
	if err != nil {
		response.Fail(c, nil, "获取内容列表失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"posts": dto.ToPostsDto(post), "total": total}, "获取内容列表成功")
}

// 创建内容
func (pc PostController) CreatePost(c *gin.Context) {
	var req vo.CreatePostRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	imageStr := strings.Join(req.Images, ",")
	post := model.Post{
		CategoryID:  req.CategoryID,
		ProjectID:   req.ProjectID,
		UserID:      req.UserID,
		Author:      req.Author,
		Title:       req.Title,
		Content:     req.Content,
		HtmlContent: req.HtmlContent,
		Description: req.Description,
		Ext:         req.Ext,
		Tag:         req.Tag,
		Icon:        req.Icon,
		Type:        req.Type,
		IsTop:       req.IsTop,
		IsPasswd:    req.IsPasswd,
		ColumnID:    req.ColumnID,
		PassWord:    req.Password,
		Time:        req.Time,
		UnitPrice:   uint(util.MoneyUtil.YuanToCents(req.UnitPrice)),
		People:      req.People,
		Location:    req.Location,
		Images:      imageStr,
		Video:       req.Video,
	}

	err := pc.PostRepository.CreatePost(&post)
	if err != nil {
		response.Fail(c, nil, "创建内容失败: "+err.Error())
		return
	}
	response.Success(c, nil, "创建内容成功")

}

// 更新内容
func (pc PostController) UpdatePostByID(c *gin.Context) {
	var req vo.UpdatePostRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}

	// 根据path中的PostID获取内容信息
	oldPost, err := pc.PostRepository.GetPostByPostID(c.Param("postID"))
	if err != nil {
		response.Fail(c, nil, "获取需要更新的内容信息失败: "+err.Error())
		return
	}
	imageStr := strings.Join(req.Images, ",")
	oldPost.Title = req.Title
	oldPost.Description = req.Description
	oldPost.IsTop = req.IsTop
	oldPost.IsPasswd = req.IsPasswd
	oldPost.ProjectID = req.ProjectID
	oldPost.PassWord = req.Password
	oldPost.Type = req.Type
	oldPost.Icon = req.Icon
	oldPost.Ext = req.Ext
	oldPost.HtmlContent = req.HtmlContent
	oldPost.Content = req.Content
	oldPost.CategoryID = req.CategoryID
	oldPost.UserID = req.UserID
	oldPost.Author = req.Author
	oldPost.Status = req.Status
	oldPost.Tag = req.Tag
	oldPost.ColumnID = req.ColumnID
	oldPost.Time = req.Time
	oldPost.UnitPrice = uint(util.MoneyUtil.YuanToCents(req.UnitPrice))
	oldPost.People = req.People
	oldPost.Location = req.Location
	oldPost.Images = imageStr
	oldPost.Video = req.Video
	// 更新内容
	err = pc.PostRepository.UpdatePost(&oldPost)
	if err != nil {
		response.Fail(c, nil, "更新内容失败: "+err.Error())
		return
	}
	response.Success(c, nil, "更新内容成功")
}

// 批量删除
func (tc PostController) BatchDeletePostByIds(c *gin.Context) {
	var req vo.DeletePostsRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}

	// 前端传来的标签ID
	reqPostIds := strings.Split(req.PostIds, ",")
	err := tc.PostRepository.BatchDeletePostByIds(reqPostIds)
	if err != nil {
		response.Fail(c, nil, "删除内容失败: "+err.Error())
		return
	}

	response.Success(c, nil, "删除内容成功")

}

// 更新内容
func (pc PostController) PushPostByID(c *gin.Context) {
	// 根据path中的PostID获取内容信息
	oldPost, err := pc.PostRepository.GetPostByPostID(c.Param("postID"))
	if err != nil {
		response.Fail(c, nil, "获取需要更新的内容信息失败: "+err.Error())
		return
	}
	oldPost.Status = known.POST_STATUS_PUBLIC
	// 更新内容
	err = pc.PostRepository.UpdatePost(&oldPost)
	if err != nil {
		response.Fail(c, nil, "更新内容失败: "+err.Error())
		return
	}
	// 同步内容至百度
	projectInfo, err := pc.ProjectRepository.GetProjectByProjectID(oldPost.ProjectID)

	if !gconvert.IsEmpty(projectInfo.PushToken) {
		// 处理 url
		postURLWithID := projectInfo.PostURL + oldPost.PostID
		go util.SEOUtil.PushBaidu(projectInfo.Domain, projectInfo.PushToken, postURLWithID)
	}

	response.Success(c, nil, "更新内容成功")
}
