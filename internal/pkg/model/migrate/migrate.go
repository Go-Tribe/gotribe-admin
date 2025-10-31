// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package migrate

import (
	"fmt"

	"gotribe-admin/internal/pkg/model"

	"gorm.io/gorm"
)

// DBAutoMigrate 自动迁移表结构
func DBAutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.Admin{},
		&model.UserEvent{},
		&model.Role{},
		&model.Menu{},
		&model.Api{},
		&model.Post{},
		&model.Tag{},
		&model.Category{},
		&model.Column{},
		&model.Comment{},
		&model.Project{},
		&model.Config{},
		&model.Resource{},
		&model.AdScene{},
		&model.Ad{},
		&model.PointLog{},
		&model.PointDeduction{},
		&model.PointAvailable{},
		&model.ProductCategory{},
		&model.ProductType{},
		&model.ProductSpec{},
		&model.ProductSpecItem{},
		&model.Product{},
		&model.ProductSku{},
		&model.Order{},
		&model.OrderLog{},
		&model.SystemConfig{},
		&model.OperationLog{},
		&model.ThirdPartyAccounts{},
		&model.Feedback{},
	)
	if err != nil {
		panic(fmt.Sprintf("database migration failed: %v", err))
	}
}
