// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/dto"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"strings"
)

type ITagController interface {
	GetTagInfo(c *gin.Context)          // 获取当前登录标签信息
	GetTags(c *gin.Context)             // 获取标签列表
	CreateTag(c *gin.Context)           // 创建标签
	UpdateTagByID(c *gin.Context)       // 更新标签
	BatchDeleteTagByIds(c *gin.Context) // 批量删除标签
}

type TagController struct {
	TagRepository repository.ITagRepository
}

// 构造函数
func NewTagController() ITagController {
	tagRepository := repository.NewTagRepository()
	tagController := TagController{TagRepository: tagRepository}
	return tagController
}

// 获取当前标签信息
func (tc TagController) GetTagInfo(c *gin.Context) {
	tag, err := tc.TagRepository.GetTagByTagID(c.Param("tagID"))
	if err != nil {
		response.Fail(c, nil, "获取当前标签信息失败: "+err.Error())
		return
	}
	tagInfoDto := dto.ToTagInfoDto(tag)
	response.Success(c, gin.H{
		"tag": tagInfoDto,
	}, "获取当前标签信息成功")
}

// 获取标签列表
func (tc TagController) GetTags(c *gin.Context) {
	var req vo.TagListRequest
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
	tag, total, err := tc.TagRepository.GetTags(&req)
	if err != nil {
		response.Fail(c, nil, "获取标签列表失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"tags": dto.ToTagsDto(tag), "total": total}, "获取标签列表成功")
}

// 创建标签
func (tc TagController) CreateTag(c *gin.Context) {
	var req vo.CreateTagRequest
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

	tag := model.Tag{
		Title:       req.Title,
		Description: req.Description,
		Color:       req.Color,
	}

	tagInfo, err := tc.TagRepository.CreateTag(&tag)
	if err != nil {
		response.Fail(c, nil, "创建标签失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"tag": dto.ToTagInfoDto(*tagInfo)}, "创建标签成功")
}

// 更新标签
func (tc TagController) UpdateTagByID(c *gin.Context) {
	var req vo.CreateTagRequest
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

	// 根据path中的TagID获取标签信息
	oldTag, err := tc.TagRepository.GetTagByTagID(c.Param("tagID"))
	if err != nil {
		response.Fail(c, nil, "获取需要更新的标签信息失败: "+err.Error())
		return
	}
	oldTag.Title = req.Title
	oldTag.Description = req.Description
	oldTag.Color = req.Color
	// 更新标签
	err = tc.TagRepository.UpdateTag(&oldTag)
	if err != nil {
		response.Fail(c, nil, "更新标签失败: "+err.Error())
		return
	}
	response.Success(c, nil, "更新标签成功")
}

// 批量删除标签
func (tc TagController) BatchDeleteTagByIds(c *gin.Context) {
	var req vo.DeleteTagsRequest
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
	reqTagIds := strings.Split(req.TagIds, ",")
	err := tc.TagRepository.BatchDeleteTagByIds(reqTagIds)
	if err != nil {
		response.Fail(c, nil, "删除标签失败: "+err.Error())
		return
	}

	response.Success(c, nil, "删除标签成功")

}
