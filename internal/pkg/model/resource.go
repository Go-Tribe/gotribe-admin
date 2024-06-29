// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type Resource struct {
	Model
	ResourceID    string `gorm:"type:char(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"resourceID"`
	Title         string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Path          string `gorm:"not null;type:varchar(255);not null;comment:路径" json:"path"`
	URL           string `gorm:"not null;type:varchar(255);not null;comment:当前域名" json:"url"`
	FileExtension string `gorm:"not null;type:char(10);not null;comment:文件拓展" json:"file_extension"`
	FileType      uint   `gorm:"type:tinyint;not null;default:1;comment:资源类形，1-图片;2-文件;3-视频;4-音频" json:"fileType"`
	Description   string `gorm:"not null;size:300;not null;comment:描述" json:"description"`
	Size          int64  `gorm:"not null;type:int;not null;default:0;comment:文件大小" json:"size"`
	Status        uint   `gorm:"type:tinyint;not null;default:1;comment:状态，1-正常；2-禁用" json:"status,omitempty"`
}

func (r *Resource) BeforeCreate(tx *gorm.DB) error {
	r.ResourceID = gid.GenShortID()

	return nil
}
