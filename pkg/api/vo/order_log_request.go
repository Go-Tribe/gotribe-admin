// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 获取订单记录列表结构体
type OrderLogListRequest struct {
	OrderID string `form:"orderID" json:"orderID"`
}
