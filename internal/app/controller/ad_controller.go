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

type IAdController interface {
	GetAdInfo(c *gin.Context)          // 获取当前广告信息
	GetAds(c *gin.Context)             // 获取广告列表
	CreateAd(c *gin.Context)           // 创建广告
	UpdateAdByID(c *gin.Context)       // 更新广告
	BatchDeleteAdByIds(c *gin.Context) // 批量删除
}

type AdController struct {
	AdRepository repository.IAdRepository
}

// 构造函数
func NewAdController() IAdController {
	adRepository := repository.NewAdRepository()
	adController := AdController{AdRepository: adRepository}
	return adController
}

// 获取当前广告信息
func (pc AdController) GetAdInfo(c *gin.Context) {
	ad, err := pc.AdRepository.GetAdByAdID(c.Param("adID"))
	if err != nil {
		response.Fail(c, nil, "获取当前广告信息失败: "+err.Error())
		return
	}
	adInfoDto := dto.ToAdInfoDto(ad)
	response.Success(c, gin.H{
		"ad": adInfoDto,
	}, "获取当前广告信息成功")
}

// 获取广告列表
func (pc AdController) GetAds(c *gin.Context) {
	var req vo.AdListRequest
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
	ad, total, err := pc.AdRepository.GetAds(&req)
	if err != nil {
		response.Fail(c, nil, "获取广告列表失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"ads": dto.ToAdsDto(ad), "total": total}, "获取广告列表成功")
}

// 创建广告
func (pc AdController) CreateAd(c *gin.Context) {
	var req vo.CreateAdRequest
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

	ad := model.Ad{
		SceneID:     req.SceneID,
		URL:         req.URL,
		URLType:     req.URLType,
		Image:       req.Image,
		Sort:        req.Sort,
		Status:      req.Status,
		Title:       req.Title,
		Video:       req.Video,
		Ext:         req.Ext,
		Description: req.Description,
	}

	err := pc.AdRepository.CreateAd(&ad)
	if err != nil {
		response.Fail(c, nil, "创建广告失败: "+err.Error())
		return
	}
	response.Success(c, nil, "创建广告成功")

}

// 更新广告
func (pc AdController) UpdateAdByID(c *gin.Context) {
	var req vo.UpdateAdRequest
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

	// 根据path中的AdID获取广告信息
	oldAd, err := pc.AdRepository.GetAdByAdID(c.Param("adID"))
	if err != nil {
		response.Fail(c, nil, "获取需要更新的广告信息失败: "+err.Error())
		return
	}
	oldAd.Title = req.Title
	oldAd.Description = req.Description
	oldAd.Ext = req.Ext
	oldAd.Image = req.Image
	oldAd.Sort = req.Sort
	oldAd.URL = req.URL
	oldAd.URLType = req.URLType
	oldAd.Status = req.Status
	oldAd.Video = req.Video
	oldAd.SceneID = req.SceneID
	// 更新广告
	err = pc.AdRepository.UpdateAd(&oldAd)
	if err != nil {
		response.Fail(c, nil, "更新广告失败: "+err.Error())
		return
	}
	response.Success(c, nil, "更新广告成功")
}

// 批量删除广告
func (pc AdController) BatchDeleteAdByIds(c *gin.Context) {
	var req vo.DeleteAdsRequest
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

	// 前端传来的广告ID
	reqAdIds := strings.Split(req.AdIds, ",")
	err := pc.AdRepository.BatchDeleteAdByIds(reqAdIds)
	if err != nil {
		response.Fail(c, nil, "删除广告失败: "+err.Error())
		return
	}

	response.Success(c, nil, "删除广告成功")

}
