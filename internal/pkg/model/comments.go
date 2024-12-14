// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

type Comment struct {
	Model
	ProjectID   string `gorm:"type:char(10);not null;index;comment:项目ID;" json:"projectID"`
	Content     string `gorm:"not null;type:longtext;not null;comment:内容" json:"content"`
	HtmlContent string `gorm:"not null;type:longtext;not null;comment:HTML内容" json:"htmlContent"`
	Status      uint   `gorm:"type:tinyint;not null;default:1;comment:状态，1-待审核；2-审核通过" json:"status,omitempty"`
	ObjectID    string `gorm:"type:char(10);not null;index;comment:评论对象ID" json:"objectID"`
	ObjectType  uint   `gorm:"type:tinyint;not null;default:1;comment:评论对象类型，1-文章；2-商品" json:"objectType"`
	Type        uint   `gorm:"type:tinyint;not null;default:1;comment:评论类型，1-评论；2-回复" json:"type"`
	UserID      string `gorm:"type:char(10);not null;index;comment:用户ID" json:"userID"`
	ToUserID    string `gorm:"type:char(10);not null;index;comment:被评论用户ID" json:"toUserID"`
	PID         int    `gorm:"type:int;not null;default:0;comment:父评论ID" json:"pid"`
	Hot         int    `gorm:"type:int;default:0;comment:热度" json:"hot"`
	Like        int    `gorm:"type:int;default:0;comment:点赞数" json:"like"`
	Dislike     int    `gorm:"type:int;default:0;comment:踩数" json:"dislike"`
}
