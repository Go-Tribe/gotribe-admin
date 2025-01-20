// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建项目结构体
type CreateProjectRequest struct {
	Name           string `form:"name" json:"name" validate:"required,min=2,max=20"`
	Title          string `form:"title" json:"title"`
	Description    string `form:"description" json:"description"`
	Keywords       string `form:"keywords" json:"keywords"`
	Domain         string `form:"domain" json:"domain"`
	PostURL        string `form:"postUrl" json:"postUrl"`
	ICP            string `form:"icp" json:"icp"`
	BaiduAnalytics string `form:"baiduAnalytics" json:"baiduAnalytics"`
	Favicon        string `form:"favicon" json:"favicon"`
	PublicSecurity string `form:"publicSecurity" json:"publicSecurity"`
	Author         string `form:"author" json:"author"`
	NavImage       string `form:"navImage" json:"navImage"`
	Info           string `form:"info" json:"info"`
	PushToken      string `form:"pushToken" json:"pushToken"`
}

// 获取项目列表结构体
type ProjectListRequest struct {
	ProjectID string `form:"projectID" json:"projectID"`
	Title     string `form:"title" json:"title"`
	PageNum   uint   `json:"pageNum" form:"pageNum"`
	PageSize  uint   `json:"pageSize" form:"pageSize"`
}

// 批量删除项目结构体
type DeleteProjectsRequest struct {
	ProjectIds string `json:"projectIds" form:"projectIds"`
}
