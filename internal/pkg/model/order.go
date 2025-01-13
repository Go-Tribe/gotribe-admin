// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	Model
	OrderID      string     `gorm:"type:char(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"OrderID"`
	OrderNumber  string     `gorm:"type:varchar(255);uniqueIndex;not null;comment:订单号" json:"orderNumber"`
	OrderType    uint       `gorm:"type:tinyint(4);not null;Index;comment:订单类型：1-普通订单；2-积分订单" json:"orderType"`
	UserID       string     `gorm:"type:char(10);not null;Index;comment:用户ID" json:"userID"`
	Username     string     `gorm:"type:varchar(255);not null;Index;comment:用户名" json:"username"`
	ProductID    string     `gorm:"type:char(10);not null;Index;comment:产品ID" json:"productID"`
	ProductSku   string     `gorm:"type:char(10);not null;comment:产品SKU" json:"productSku"`
	ProductName  string     `gorm:"type:varchar(255);not null;comment:产品名称" json:"productName"`
	Status       uint       `gorm:"type:tinyint(4);not null;Index;comment:状态1-待支付；2-已支付；3-已发货；4-已收货；5-已取消；6-待退款；7.已退款" json:"status"`
	PayNumber    string     `gorm:"type:varchar(255);not null;comment:支付单号" json:"payNumber"`
	PayTime      *time.Time `gorm:"type:datetime;default:null;comment:支付时间" json:"payTime"`
	PayMethod    uint       `gorm:"type:tinyint(4);not null;comment:支付方式：1-微信支付；2-支付宝支付；3-积分支付；4-余额支付" json:"payMethod"`
	RefundTime   *time.Time `gorm:"type:datetime;default:null;comment:退款时间" json:"refundTime"`
	PayStatus    uint       `gorm:"type:tinyint(4);not null;comment:支付状态：1-待支付；2-已支付；3-已退款" json:"payStatus"`
	RefundStatus uint       `gorm:"type:tinyint(4);not null;comment:退款状态：1-待退款；2-已退款" json:"refundStatus"`
	ProjectID    string     `gorm:"type:varchar(10);Index;comment:项目 ID" json:"projectID"`
	ProductImage string     `gorm:"type:varchar(255);not null;comment:产品主图" json:"productImage"`
	Amount       int        `gorm:"type:int(10);not null;comment:总金额" json:"amount"`
	AmountPay    int        `gorm:"type:int(10);not null;comment:实际支付金额" json:"amountPay"`
	Quantity     uint       `gorm:"type:int(10);not null;comment:购买数量" json:"quantity"`
	UnitPrice    int        `gorm:"type:int(10);not null;comment:商品价格" json:"unitPrice"`
	UnitPoint    int        `gorm:"type:float(20,2);NOT NULL;comment:积分数值" json:"unitPoint"`
	// 收货人信息
	ConsigneeName     string `gorm:"type:varchar(255);not null;comment:收货人姓名" json:"consigneeName"`
	ConsigneePhone    string `gorm:"type:varchar(20);not null;comment:收货人电话" json:"consigneePhone"`
	ConsigneeProvince string `gorm:"type:varchar(100);not null;comment:省" json:"consigneeProvince"`
	ConsigneeCity     string `gorm:"type:varchar(100);not null;comment:市" json:"consigneeCity"`
	ConsigneeDistrict string `gorm:"type:varchar(100);not null;comment:区/县" json:"consigneeDistrict"`
	ConsigneeStreet   string `gorm:"type:varchar(255);not null;comment:街道" json:"consigneeStreet"`
	ConsigneeAddress  string `gorm:"type:varchar(255);not null;comment:详细地址" json:"consigneeAddress"`
	Remark            string `gorm:"type:varchar(255);not null;comment:买家留言" json:"remark"`
	RemarkAdmin       string `gorm:"type:varchar(255);not null;comment:订单备注" json:"remarkAdmin"`
	User              *User  `gorm:"-" json:"user"`
}

func (t *Order) BeforeCreate(tx *gorm.DB) error {
	t.OrderID = gid.GenShortID()
	return nil
}
