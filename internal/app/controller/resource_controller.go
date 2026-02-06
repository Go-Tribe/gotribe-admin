// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"gotribe-admin/config"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/dto"
	"gotribe-admin/pkg/api/known"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"gotribe-admin/pkg/util"
	"gotribe-admin/pkg/util/upload"

	"github.com/dengmengmian/ghelper/gconvert"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IResourceController interface {
	GetResourceInfo(c *gin.Context)    // 获取当前资源信息
	GetResources(c *gin.Context)       // 获取资源列表
	UpdateResourceByID(c *gin.Context) // 更新资源
	UploadResources(c *gin.Context)    //上传资源
	DeleteResourceByID(c *gin.Context) // 删除资源
}

type ResourceController struct {
	ResourceRepository repository.IResourceRepository
}

// 构造函数
func NewResourceController() IResourceController {
	resourceRepository := repository.NewResourceRepository()
	resourceController := ResourceController{ResourceRepository: resourceRepository}
	return resourceController
}

// GetResourceInfo 获取当前资源信息
// @Summary 获取资源信息
// @Description 根据资源ID获取资源详细信息
// @Tags 资源管理
// @Accept json
// @Produce json
// @Param resourceID path string true "资源ID"
// @Success 200 {object} response.Response{data=object{resource=dto.ResourceDto}} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/resources/{resourceID} [get]
// @Security BearerAuth
func (pc ResourceController) GetResourceInfo(c *gin.Context) {
	resource, err := pc.ResourceRepository.GetResourceByResourceID(c.Param("resourceID"))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	resourceInfoDto := dto.ToResourceInfoDto(resource)
	response.Success(c, gin.H{
		"resource": resourceInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// GetResources 获取资源列表
// @Summary 获取资源列表
// @Description 获取资源列表，支持分页和筛选
// @Tags 资源管理
// @Accept json
// @Produce json
// @Param request query vo.ResourceListRequest false "查询参数"
// @Success 200 {object} response.Response{data=object{resources=[]dto.ResourceDto,total=int}} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/resources [get]
// @Security BearerAuth
func (pc ResourceController) GetResources(c *gin.Context) {
	var req vo.ResourceListRequest
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
	resource, total, err := pc.ResourceRepository.GetResources(&req)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgListFail)
		return
	}
	response.Success(c, gin.H{"resources": dto.ToResourcesDto(resource), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// UpdateResourceByID 更新资源
// @Summary 更新资源
// @Description 根据资源ID更新资源信息
// @Tags 资源管理
// @Accept json
// @Produce json
// @Param resourceID path string true "资源ID"
// @Param request body vo.CreateResourceRequest true "更新资源请求参数"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/resources/{resourceID} [put]
// @Security BearerAuth
func (pc ResourceController) UpdateResourceByID(c *gin.Context) {
	var req vo.CreateResourceRequest
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

	// 根据path中的ResourceID获取资源信息
	oldResource, err := pc.ResourceRepository.GetResourceByResourceID(c.Param("resourceID"))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	oldResource.Title = req.Title
	oldResource.Description = req.Description
	// 更新资源
	err = pc.ResourceRepository.UpdateResource(&oldResource)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgUpdateFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// UploadResources 上传资源
// @Summary 上传资源
// @Description 上传文件资源到服务器
// @Tags 资源管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "上传的文件"
// @Success 200 {object} response.Response{data=object{upload=dto.UploadResourceDto}} "上传成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/resources/upload [post]
// @Security BearerAuth
func (pc ResourceController) UploadResources(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.HandleBindError(c, err)
		return
	}
	if fileHeader.Size > known.DEFAULT_UPLOAD_SIZE {
		response.ValidationFail(c, "上传资源过大")
		return
	}
	provider := config.Conf.UploadFile.GetUploadProvider(config.Conf.System.EnableOss)
	upload, err := upload.NewService(
		provider,
		config.Conf.UploadFile.Endpoint,
		config.Conf.UploadFile.Accesskey,
		config.Conf.UploadFile.Secretkey,
		config.Conf.UploadFile.Bucket,
	)
	if err != nil {
		response.InternalServerError(c, "初始化上传服务失败："+err.Error())
		return
	}
	fileRes, err := upload.UploadFile(fileHeader)
	if err != nil {
		response.InternalServerError(c, "上传 CDN 失败："+err.Error())
		return
	}
	uploadRes := dto.ToUploadResourceDto(&fileRes)
	uploadRes.Domain = config.Conf.System.CDNDomain
	fileType, err := util.FileUtil.GetFileType(fileHeader)
	if err != nil {
		response.InternalServerError(c, "获取文件类型失败: "+err.Error())
		return
	}
	uploadRes.FileType = fileType

	// 资源入库
	resource := model.Resource{
		Title:         fileHeader.Filename,
		Path:          uploadRes.Key,
		URL:           uploadRes.Domain,
		FileExtension: uploadRes.FileExt,
		Size:          fileHeader.Size,
		FileType:      gconvert.Uint(uploadRes.FileType),
	}

	if err = pc.ResourceRepository.CreateResource(&resource); err != nil {
		response.HandleDatabaseError(c, err, common.MsgCreateFail)
		return
	}
	response.Success(c, gin.H{"upload": uploadRes}, "上传资源成功")
}

// DeleteResourceByID 删除资源
// @Summary 删除资源
// @Description 根据资源ID删除资源
// @Tags 资源管理
// @Accept json
// @Produce json
// @Param request body vo.DeleteResourcesRequest true "删除资源请求参数"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/resources [delete]
// @Security BearerAuth
func (pc ResourceController) DeleteResourceByID(c *gin.Context) {
	var req vo.DeleteResourcesRequest
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

	err := pc.ResourceRepository.DeleteResourceByID(req.ResourceID)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgDeleteFail)
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))
}
