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

// GetIndexInfo retrieves the index page information based on the project ID.
// It calls the repository to get data and returns a success or fail response.
func (pc IndexController) GetIndexInfo(c *gin.Context) {
	indexInfo, err := pc.IndexRepository.GetIndexData(c.Param("projectID"))
	if err != nil {
		response.Fail(c, nil, "获取首页头部数据失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{
		"indexDate": indexInfo,
	}, "获取首页头部数据成功")
}

// GetTimeRangeData retrieves the time range data for the index page based on the specified time range and project ID.
// It calls the repository to get data and returns a success or fail response.
func (pc IndexController) GetTimeRangeData(c *gin.Context) {
	timeRangeData, err := pc.IndexRepository.GetTimeRangeData(c.Param("timeRange"), c.Param("projectID"))
	if err != nil {
		response.Fail(c, nil, "获取首页折线数据失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{
		"timeRangeData": timeRangeData,
	}, "获取首页折线数据成功")
}
