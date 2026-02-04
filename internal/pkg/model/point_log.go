// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

// PointLog 积分记录表
type PointLog struct {
	Model
	ProjectID string `gorm:"type:varchar(10);not null;index;comment:项目ID;" json:"projectID"`
	UserID    string `gorm:"type:varchar(10);Index;comment:用户ID" json:"userID"`
	Points    int64  `gorm:"type:bigint;NOT NULL;comment:积分数值(分)" json:"points"`
	Reason    string `gorm:"type:varchar(255);NOT NULL;comment:加减原因" json:"reason"`
	Type      string `gorm:"type:varchar(20);NOT NULL;comment:类型" json:"type"`
	EventID   string `gorm:"type:varchar(10);comment:事件ID" json:"eventID"`
	Status    uint   `gorm:"type:smallint;not null;default:1;comment:状态，1-正常；2-删除" json:"status"`
	User      *User  `gorm:"-" json:"user"`
}
