// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"github.com/dengmengmian/ghelper/gconvert"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

// 获取当前资源信息
func (pc ResourceController) GetResourceInfo(c *gin.Context) {
	resource, err := pc.ResourceRepository.GetResourceByResourceID(c.Param("resourceID"))
	if err != nil {
		response.Fail(c, nil, "获取当前资源信息失败: "+err.Error())
		return
	}
	resourceInfoDto := dto.ToResourceInfoDto(resource)
	response.Success(c, gin.H{
		"resource": resourceInfoDto,
	}, "获取当前资源信息成功")
}

// 获取资源列表
func (pc ResourceController) GetResources(c *gin.Context) {
	var req vo.ResourceListRequest
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
	resource, total, err := pc.ResourceRepository.GetResources(&req)
	if err != nil {
		response.Fail(c, nil, "获取资源列表失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"resources": dto.ToResourcesDto(resource), "total": total}, "获取资源列表成功")
}

// 更新资源
func (pc ResourceController) UpdateResourceByID(c *gin.Context) {
	var req vo.CreateResourceRequest
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

	// 根据path中的ResourceID获取资源信息
	oldResource, err := pc.ResourceRepository.GetResourceByResourceID(c.Param("resourceID"))
	if err != nil {
		response.Fail(c, nil, "获取需要更新的资源信息失败: "+err.Error())
		return
	}
	oldResource.Title = req.Title
	oldResource.Description = req.Description
	// 更新资源
	err = pc.ResourceRepository.UpdateResource(&oldResource)
	if err != nil {
		response.Fail(c, nil, "更新资源失败: "+err.Error())
		return
	}
	response.Success(c, nil, "更新资源成功")
}

// 上传资源
func (pc ResourceController) UploadResources(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, nil, "上传资源失败: "+err.Error())
		return
	}
	if fileHeader.Size > known.DEFAULT_UPLOAD_SIZE {
		response.Fail(c, nil, "上传资源过大")
		return
	}
	qiniuUpload := upload.NewQiniu(config.Conf.QiniuConfig.Accesskey, config.Conf.QiniuConfig.Secretkey, config.Conf.QiniuConfig.Bucket)
	fileRes, err := qiniuUpload.UploadFile(fileHeader)
	if err != nil {
		response.Fail(c, nil, "上传 CDN 失败："+err.Error())
		return
	}
	uploadRes := dto.ToUploadResourceDto(fileRes)
	uploadRes.Domain = config.Conf.QiniuConfig.Domain
	uploadRes.FileType = util.GetFileType(fileHeader)

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
		response.Fail(c, nil, "创建资源失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"upload": uploadRes}, "上传资源成功")
}

// 批量删除
func (pc ResourceController) DeleteResourceByID(c *gin.Context) {
	var req vo.DeleteResourcesRequest
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

	err := pc.ResourceRepository.DeleteResourceByID(req.ResourceID)
	if err != nil {
		response.Fail(c, nil, "删除资源失败: "+err.Error())
		return
	}

	response.Success(c, nil, "删除资源成功")
}
