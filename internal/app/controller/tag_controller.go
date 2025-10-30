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
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	tagInfoDto := dto.ToTagInfoDto(tag)
	response.Success(c, gin.H{
		"tag": tagInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
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
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.Fail(c, nil, errStr)
		return
	}

	// 获取
	tag, total, err := tc.TagRepository.GetTags(&req)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"tags": dto.ToTagsDto(tag), "total": total}, common.Msg(c, common.MsgListSuccess))
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
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
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
		response.Fail(c, nil, common.Msg(c, common.MsgCreateFail)+": "+err.Error())
		return
	}
	response.Success(c, gin.H{"tag": dto.ToTagInfoDto(*tagInfo)}, common.Msg(c, common.MsgCreateSuccess))
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
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.Fail(c, nil, errStr)
		return
	}

	// 根据path中的TagID获取标签信息
	oldTag, err := tc.TagRepository.GetTagByTagID(c.Param("tagID"))
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgGetFail)+": "+err.Error())
		return
	}
	oldTag.Title = req.Title
	oldTag.Description = req.Description
	oldTag.Color = req.Color
	// 更新标签
	err = tc.TagRepository.UpdateTag(&oldTag)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgUpdateFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
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
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.Fail(c, nil, errStr)
		return
	}

	// 前端传来的标签ID
	reqTagIds := strings.Split(req.TagIds, ",")
	err := tc.TagRepository.BatchDeleteTagByIds(reqTagIds)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgDeleteFail)+": "+err.Error())
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))

}
