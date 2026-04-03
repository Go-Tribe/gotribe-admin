// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import "time"

// PointAvailable 积分记录表
type PointAvailable struct {
	Model
	ProjectID      string    `gorm:"type:varchar(10);not null;index:idx_point_available_project_user_status,priority:1;comment:项目ID;" json:"projectID"`
	UserID         uint      `gorm:"index:idx_point_available_project_user_status,priority:2;index:idx_point_available_user_status_expiration,priority:1;comment:用户ID" json:"userID"`
	Points         int64     `gorm:"type:bigint;NOT NULL;comment:积分数值(分)"`
	PointsLogID    uint      `gorm:"not null;comment:'积分记录表ID'"`
	ExpirationDate time.Time `gorm:"column:expiration_date;index:idx_point_available_user_status_expiration,priority:3;comment:'过期时间'"`
	Status         uint      `gorm:"type:smallint;not null;default:1;index:idx_point_available_project_user_status,priority:3;index:idx_point_available_user_status_expiration,priority:2;comment:状态，1-正常；2-删除" json:"status"`
}
