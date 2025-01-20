// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type Config struct {
	Model
	ConfigID    string   `gorm:"type:char(10);not null;uniqueIndex;comment:字符ID，分布式 ID;" json:"configID"`
	ProjectID   string   `gorm:"type:char(10);not null;index;comment:项目ID;" json:"projectID"`
	Alias       string   `gorm:"type:varchar(20);not null;uniqueIndex;comment:别名" json:"alias"`
	Title       string   `gorm:"type:varchar(30);not null;comment:标题" json:"title"`
	Description string   `gorm:"type:varchar(300);not null;comment:描述" json:"description"`
	Type        uint     `gorm:"type:tinyint;not null;default:1;comment:类型，1表示普通配置2:json类型" json:"type"`
	Info        string   `gorm:"type:longtext;not null;comment:内容" json:"info"`
	MDContent   string   `gorm:"type:longtext;not null;comment:MD内容" json:"mdContent"`
	Status      uint     `gorm:"type:tinyint;not null;default:1;comment:状态，1-正常；2-禁用" json:"status"`
	Project     *Project `gorm:"-" json:"project"`
}

func (c *Config) BeforeCreate(tx *gorm.DB) error {
	c.ConfigID = gid.GenShortID()

	return nil
}
