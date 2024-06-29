// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

// 角色权限规则
type RoleCasbin struct {
	Keyword string `json:"keyword"` // 角色关键字
	Path    string `json:"path"`    // 访问路径
	Method  string `json:"method"`  // 请求方式
}
