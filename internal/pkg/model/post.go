// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type Post struct {
	Model
	PostID      string    `gorm:"type:varchar(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"postID"`
	CategoryID  string    `gorm:"type:varchar(10);Index;comment:分类 ID" json:"categoryID"`
	ProjectID   string    `gorm:"type:varchar(10);Index;comment:项目 ID" json:"projectID"`
	ColumnID    string    `gorm:"type:varchar(10);Index;comment:专栏ID" json:"columnID"`
	UserID      string    `gorm:"type:varchar(10);Index;comment:用户ID" json:"userID"`
	Author      string    `gorm:"type:varchar(30);not null;index:idx_username;comment:作者" json:"author"`
	Title       string    `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Content     string    `gorm:"not null;type:text;comment:内容" json:"content"`
	HtmlContent string    `gorm:"not null;type:text;comment:html内容" json:"htmlContent"`
	Description string    `gorm:"not null;size:300;comment:描述" json:"description"`
	Ext         string    `gorm:"type:text;comment:'扩展字段'" json:"ext"`
	Icon        string    `gorm:"type:varchar(255);comment:图标" json:"icon"`
	Tag         string    `gorm:"type:varchar(30);comment:tag" json:"tag"`
	View        uint      `gorm:"default:1;comment:'阅读量'" json:"view"`
	Type        uint      `gorm:"type:smallint;default:1;comment:类型，1.文章 2.page 3.短文" json:"type"`
	IsTop       uint      `gorm:"type:smallint;default:1;comment:是否置顶：1-禁用;2-启用" json:"isTop"`
	IsPasswd    uint      `gorm:"type:smallint;default:1;comment:是否加密：1-禁用;2-启用" json:"isPasswd"`
	PassWord    string    `gorm:"type:varchar(255);not null;comment:密码" json:"password"`
	Status      uint      `gorm:"type:smallint;not null;default:1;comment:状态，1-草稿；2-发布" json:"status"`
	UnitPrice   uint      `gorm:"type:integer;not null;comment:商品价格(分)" json:"unitPrice"`
	Location    string    `gorm:"type:varchar(255);comment:地点" json:"location"`
	People      string    `gorm:"type:varchar(255);comment:人物" json:"people"`
	Time        string    `gorm:"type:varchar(255);comment:时间" json:"time"`
	Images      string    `gorm:"type:varchar(1000);comment:图片" json:"images"`
	Video       string    `gorm:"type:varchar(255);not null;comment:产品视频" json:"video"`
	Category    *Category `gorm:"-" json:"category"`
	Tags        []*Tag    `gorm:"-" json:"tags"`
	Project     *Project  `gorm:"-" json:"project"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	p.PostID = gid.GenShortID()

	return nil
}
