// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type Product struct {
	Model
	ProductID     string `gorm:"type:varchar(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"productID"`
	Title         string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	ProductNumber string `gorm:"type:varchar(255);not null;comment:商品货号" json:"productNumber"`
	ProjectID     string `gorm:"type:varchar(10);Index;comment:项目 ID" json:"projectID"`
	Description   string `gorm:"not null;size:300;not null;comment:产品卖点/描述" json:"description"`
	Image         string `gorm:"type:varchar(255);not null;comment:产品主图" json:"image"`
	Video         string `gorm:"type:varchar(255);not null;comment:产品视频" json:"video"`
	BuyLimit      uint   `gorm:"type:smallint;not null;default:1;comment:购买限制" json:"buyLimit"`
	CategoryID    string `gorm:"type:varchar(10);not null;index;comment:分类ID" json:"categoryID"`
	ProductSpec   string `gorm:"type:varchar(2048);not null;comment:产品规格" json:"productSpec"`
	Content       string `gorm:"type:text;comment:内容" json:"content"`
	HtmlContent   string `gorm:"type:text;comment:html内容" json:"Htmlcontent"`
	Tag           string `gorm:"type:varchar(300);not null;comment:标签" json:"tag"`
	Tags          []*Tag `gorm:"-" json:"tags"`
	Enable        uint   `gorm:"type:smallint;not null;default:1;comment:是否启用：1-下架；2-上架" json:"enable"`
}

func (e *Product) BeforeCreate(tx *gorm.DB) error {
	e.ProductID = gid.GenShortID()
	return nil
}
