// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

// Package controller defines the web APIs and handles HTTP requests.
package controller

// Importing necessary packages
import (
	"github.com/gin-gonic/gin"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/pkg/api/response"
)

// IIndexController is an interface defining the methods for index page data retrieval.
type IIndexController interface {
	GetIndexInfo(c *gin.Context)
	GetTimeRangeData(c *gin.Context)
}

// IndexController is a struct implementing the IIndexController interface.
type IndexController struct {
	IndexRepository repository.IIndexRepository
}

// NewIndexController is a constructor function for creating a new instance of IndexController.
// It initializes the IndexRepository and returns an IIndexController interface.
func NewIndexController() IIndexController {
	indexRepository := repository.NewIndexRepository()
	indexController := IndexController{IndexRepository: indexRepository}
	return indexController
}

// GetIndexInfo 获取首页信息
// @Summary      获取首页信息
// @Description  根据项目ID获取首页头部数据
// @Tags         首页管理
// @Accept       json
// @Produce      json
// @Param        projectID query string false "项目ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /index/info [get]
// @Security     BearerAuth
func (pc IndexController) GetIndexInfo(c *gin.Context) {
	indexInfo, err := pc.IndexRepository.GetIndexData(c.Query("projectID"))
	if err != nil {
		response.Fail(c, nil, "获取首页头部数据失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{
		"indexDate": indexInfo,
	}, "获取首页头部数据成功")
}

// GetTimeRangeData 获取时间范围数据
// @Summary      获取时间范围数据
// @Description  根据时间范围和项目ID获取首页折线图数据
// @Tags         首页管理
// @Accept       json
// @Produce      json
// @Param        projectID query string false "项目ID"
// @Param        timeRange query string false "时间范围"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /index/time-range [get]
// @Security     BearerAuth
func (pc IndexController) GetTimeRangeData(c *gin.Context) {
	timeRangeData, err := pc.IndexRepository.GetTimeRangeData(c.Query("projectID"), c.Query("timeRange"))
	if err != nil {
		response.Fail(c, nil, "获取首页折线数据失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{
		"timeRangeData": timeRangeData,
	}, "获取首页折线数据成功")
}
