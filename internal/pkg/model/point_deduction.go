// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

// PointDeduction 扣减积分表
type PointDeduction struct {
	Model
	ProjectID         string  `gorm:"type:char(10);not null;index;comment:项目ID;" json:"projectID"`
	UserID            string  `gorm:"type:varchar(10);Index;comment:用户ID" json:"userID"`
	Points            float64 `gorm:"type:float(20,2);NOT NULL;comment:积分数值"`
	PointsDetailID    int     `gorm:"type:int(10);comment:'积分明细ID'"`
	AvailablePointsID int     `gorm:"type:int(10);NOT NULL;comment:'可用积分表ID'"`
}
