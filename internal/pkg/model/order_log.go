// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type OrderLog struct {
	Model
	OrderLogID string `gorm:"type:char(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"OrderLogID"`
	OrderID    string `gorm:"type:varchar(255);not null;comment:订单号" json:"orderID"`
	Remark     string `gorm:"type:varchar(255);not null;comment:操作记录" json:"remark"`
}

func (t *OrderLog) BeforeCreate(tx *gorm.DB) error {
	t.OrderLogID = gid.GenShortID()
	return nil
}
