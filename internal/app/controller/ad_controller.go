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

// GetAdInfo 获取当前广告信息
// @Summary 获取广告详情
// @Description 根据广告ID获取广告详细信息
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param adID path string true "广告ID"
// @Success 200 {object} response.Response{data=map[string]dto.AdDto} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad/{adID} [get]
// @Security BearerAuth
func (pc AdController) GetAdInfo(c *gin.Context) {
	ad, err := pc.AdRepository.GetAdByAdID(c.Param("adID"))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	adInfoDto := dto.ToAdInfoDto(ad)
	response.Success(c, gin.H{
		"ad": adInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// GetAds 获取广告列表
// @Summary 获取广告列表
// @Description 根据查询条件获取广告列表，支持分页
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param sceneID query string false "场景ID"
// @Param title query string false "广告标题"
// @Param status query int false "状态(1:启用 2:禁用)"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad [get]
// @Security BearerAuth
func (pc AdController) GetAds(c *gin.Context) {
	var req vo.AdListRequest
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
	ad, total, err := pc.AdRepository.GetAds(&req)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgListFail)
		return
	}
	response.Success(c, gin.H{"ads": dto.ToAdsDto(ad), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// CreateAd 创建广告
// @Summary 创建广告
// @Description 创建新的广告
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param ad body vo.CreateAdRequest true "广告信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad [post]
// @Security BearerAuth
func (pc AdController) CreateAd(c *gin.Context) {
	var req vo.CreateAdRequest
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
		response.HandleDatabaseError(c, err, common.MsgCreateFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgCreateSuccess))

}

// UpdateAdByID 更新广告
// @Summary 更新广告
// @Description 根据广告ID更新广告信息
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param adID path string true "广告ID"
// @Param ad body vo.UpdateAdRequest true "广告信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad/{adID} [patch]
// @Security BearerAuth
func (pc AdController) UpdateAdByID(c *gin.Context) {
	var req vo.UpdateAdRequest
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

	// 根据path中的AdID获取广告信息
	oldAd, err := pc.AdRepository.GetAdByAdID(c.Param("adID"))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
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
		response.HandleDatabaseError(c, err, common.MsgUpdateFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// BatchDeleteAdByIds 批量删除广告
// @Summary 批量删除广告
// @Description 根据广告ID列表批量删除广告
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param ads body vo.DeleteAdsRequest true "广告ID列表"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad [delete]
// @Security BearerAuth
func (pc AdController) BatchDeleteAdByIds(c *gin.Context) {
	var req vo.DeleteAdsRequest
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

	// 前端传来的广告ID
	reqAdIds := strings.Split(req.AdIds, ",")
	err := pc.AdRepository.BatchDeleteAdByIds(reqAdIds)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgDeleteFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))

}
