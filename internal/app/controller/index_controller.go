// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"github.com/gin-gonic/gin"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/pkg/api/response"
)

type IIndexController interface {
	GetIndexInfo(c *gin.Context)
	GetTimeRangeData(c *gin.Context)
}

type IndexController struct {
	IndexRepository repository.IIndexRepository
}

// 构造函数
func NewIndexController() IIndexController {
	indexRepository := repository.NewIndexRepository()
	indexController := IndexController{IndexRepository: indexRepository}
	return indexController
}

func (pc IndexController) GetIndexInfo(c *gin.Context) {
	// 获取当前广告信息
	indexInfo, err := pc.IndexRepository.GetIndexData(c.Param("porjectID"))
	if err != nil {
		response.Fail(c, nil, "获取首页头部数据失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{
		"indexDate": indexInfo,
	}, "获取首页头部数据成功")
}

func (pc IndexController) GetTimeRangeData(c *gin.Context) {
	// 获取当前广告信息
	timeRangeData, err := pc.IndexRepository.GetTimeRangeData(c.Param("timeRange"), c.Param("porjectID"))
	if err != nil {
		response.Fail(c, nil, "获取首页折线数据失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{
		"timeRangeData": timeRangeData,
	}, "获取首页折线数据成功")
}
