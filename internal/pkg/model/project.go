// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type Project struct {
	Model
	ProjectID      string `gorm:"type:char(10);not null;uniqueIndex:idx_project_project_id;comment:字符ID，分布式ID" json:"projectID"`
	Name           string `gorm:"type:varchar(30);not null;comment:项目名" json:"name,omitempty"`
	Title          string `gorm:"type:varchar(30);not null;comment:网站标题" json:"title,omitempty"`
	Description    string `gorm:"type:varchar(300);comment:描述" json:"description,omitempty"`
	Keywords       string `gorm:"type:varchar(30);comment:网站关键词" json:"keywords,omitempty"`
	Domain         string `gorm:"type:varchar(60);comment:项目域名" json:"domain,omitempty"`
	PostURL        string `gorm:"type:varchar(300);comment:内容链接" json:"postURL,omitempty"`
	ICP            string `gorm:"type:varchar(255);comment:icp备案信息" json:"icp,omitempty"`
	PublicSecurity string `gorm:"type:varchar(255);comment:公安备案" json:"publicSecurity,omitempty"`
	Author         string `gorm:"type:varchar(30);comment:网站版权" json:"author,omitempty"`
	Info           string `gorm:"type:longtext;comment:内容" json:"info,omitempty"`
	BaiduAnalytics string `gorm:"type:varchar(255);comment:百度统计" json:"baiduAnalytics,omitempty"`
	Favicon        string `gorm:"type:varchar(255);comment:favicon" json:"favicon,omitempty"`
	NavImage       string `gorm:"type:varchar(255);comment:导航图片" json:"navImage,omitempty"`
	PushToken      string `gorm:"type:varchar(255);comment:百度推送 API token" json:"pushToken,omitempty"`
	Status         int8   `gorm:"type:tinyint;not null;default:1;comment:状态，1-正常；2-禁用" json:"status,omitempty"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	p.ProjectID = gid.GenShortID()

	return nil
}
