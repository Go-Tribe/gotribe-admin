// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import "time"

// PointAvailable 积分记录表
type PointAvailable struct {
	Model
	ProjectID      string    `gorm:"type:char(10);not null;index;comment:项目ID;" json:"projectID"`
	UserID         string    `gorm:"type:varchar(10);Index;comment:用户ID" json:"userID"`
	Points         float64   `gorm:"type:float(20,2);NOT NULL;comment:积分数值"`
	PointsLogID    int       `gorm:"type:int;NOT NULL;comment:'积分记录表ID'"`
	ExpirationDate time.Time `gorm:"column:expiration_date;comment:'过期时间'"`
	Status         uint      `gorm:"type:tinyint(1);not null;default:1;comment:状态，1-正常；2-删除" json:"status"`
}
