// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package docs

import (
	"gorm.io/gorm"
	"time"
)

// SwaggerDeletedAt 为 gorm.DeletedAt 提供 Swagger 类型定义
// swagger:model DeletedAt
type SwaggerDeletedAt struct {
	// 删除时间
	// example: 2023-01-01T00:00:00Z
	Time *time.Time `json:"time,omitempty"`
	// 是否有效
	// example: true
	Valid bool `json:"valid"`
}

// SwaggerModel 为 model.Model 提供 Swagger 类型定义
// swagger:model Model
type SwaggerModel struct {
	// ID 主键
	// example: 1
	ID uint `json:"id"`
	// 创建时间
	// example: 2023-01-01T00:00:00Z
	CreatedAt time.Time `json:"createdAt"`
	// 更新时间
	// example: 2023-01-01T00:00:00Z
	UpdatedAt time.Time `json:"updatedAt"`
	// 删除时间
	DeletedAt *SwaggerDeletedAt `json:"deletedAt,omitempty"`
}

// 确保类型被正确导入
var _ gorm.DeletedAt
