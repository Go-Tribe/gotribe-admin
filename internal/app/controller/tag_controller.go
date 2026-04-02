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

// GetTagInfo 获取标签信息
// @Summary      获取标签信息
// @Description  根据标签ID获取标签详细信息
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        id path int true "标签ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /tag/{id} [get]
// @Security     BearerAuth
func (tc TagController) GetTagInfo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationFail(c, "标签ID格式错误")
		return
	}
	ctx := c.Request.Context()
	tag, err := tc.TagRepository.GetTagByID(ctx, uint(id))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	tagInfoDto := dto.ToTagInfoDto(tag)
	response.Success(c, gin.H{
		"tag": tagInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// GetTags 获取标签列表
// @Summary      获取标签列表
// @Description  获取所有标签的列表
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        request query vo.TagListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /tag [get]
// @Security     BearerAuth
func (tc TagController) GetTags(c *gin.Context) {
	var req vo.TagListRequest
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

	ctx := c.Request.Context()
	// 获取
	tag, total, err := tc.TagRepository.GetTags(ctx, &req)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgListFail)
		return
	}
	response.Success(c, gin.H{"tags": dto.ToTagsDto(tag), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// CreateTag 创建标签
// @Summary      创建标签
// @Description  创建一个新的标签
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        request body vo.CreateTagRequest true "创建标签请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /tag [post]
// @Security     BearerAuth
func (tc TagController) CreateTag(c *gin.Context) {
	var req vo.CreateTagRequest
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

	tag := model.Tag{
		Title:       req.Title,
		Slug:        req.Slug,
		Description: req.Description,
		Color:       req.Color,
	}

	ctx := c.Request.Context()
	tagInfo, err := tc.TagRepository.CreateTag(ctx, &tag)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgCreateFail)
		return
	}
	response.Success(c, gin.H{"tag": dto.ToTagInfoDto(*tagInfo)}, common.Msg(c, common.MsgCreateSuccess))
}

// UpdateTagByID 更新标签
// @Summary      更新标签
// @Description  根据标签ID更新标签信息
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        id path int true "标签ID"
// @Param        request body vo.CreateTagRequest true "标签信息"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /tag/{id} [patch]
// @Security     BearerAuth
func (tc TagController) UpdateTagByID(c *gin.Context) {
	var req vo.CreateTagRequest
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

	// 根据path中的ID获取标签信息
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationFail(c, "标签ID格式错误")
		return
	}
	ctx := c.Request.Context()
	oldTag, err := tc.TagRepository.GetTagByID(ctx, uint(id))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	oldTag.Title = req.Title
	oldTag.Slug = req.Slug
	oldTag.Description = req.Description
	oldTag.Color = req.Color
	oldTag.Sort = req.Sort
	oldTag.Status = req.Status
	// 更新标签
	err = tc.TagRepository.UpdateTag(ctx, &oldTag)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgUpdateFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// BatchDeleteTagByIds 批量删除标签
// @Summary      批量删除标签
// @Description  根据标签ID列表批量删除标签
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        request body vo.DeleteTagsRequest true "标签ID列表"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /tag [delete]
// @Security     BearerAuth
func (tc TagController) BatchDeleteTagByIds(c *gin.Context) {
	var req vo.DeleteTagsRequest
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

	ctx := c.Request.Context()
	err := tc.TagRepository.BatchDeleteTagByIds(ctx, req.Ids)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgDeleteFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))

}
