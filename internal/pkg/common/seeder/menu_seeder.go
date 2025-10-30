// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package seeder

import (
	"gotribe-admin/internal/pkg/model"

	"gorm.io/gorm"
)

// MenuSeeder 菜单种子
type MenuSeeder struct {
	*BaseSeeder
}

// NewMenuSeeder 创建菜单种子
func NewMenuSeeder() *MenuSeeder {
	return &MenuSeeder{
		BaseSeeder: NewBaseSeeder("menu"),
	}
}

// Run 执行菜单数据种子
func (s *MenuSeeder) Run(db *gorm.DB) error {
	// 获取角色
	var adminRole model.Role
	if err := db.First(&adminRole, 1).Error; err != nil {
		return err
	}

	var uint0 uint = 0
	var uint1 uint = 1
	var uint6 uint = 6
	var uint8 uint = 8
	var uint10 uint = 10
	var uint18 uint = 18
	var uint22 uint = 22
	var uint27 uint = 27

	componentStr := "component"
	systemUserStr := "/system/user"
	userStr := "user"
	peoplesStr := "peoples"
	treeTableStr := "tree-table"
	treeStr := "tree"
	logOperationStr := "/log/operation-log"
	documentationStr := "documentation"
	education := "education"
	tagIcon := "24gf-tags2"
	language := "language"
	xitongrizhi := "xitongrizhi"
	skill := "skill"
	yewu := "yewu"
	xiangmu := "xiangmu"
	nested := "nested"
	shuju := "shuju"
	ziyuan := "ziyuan"
	yunyingzhongxin := "yunyingzhongxin"
	eye := "eye"
	message := "message"
	jifen := "jifen"
	shopping := "shopping"
	list := "list"
	theme := "theme"
	guige := "guige"
	shangpinliebiao := "shangpinliebiao"
	shangpin := "shangpin-"
	dingdanliebiao := "dingdanliebiao"
	peizhishezhi := "peizhishezhi"

	menus := []model.Menu{
		{
			Model:     model.Model{ID: 1},
			Name:      "System",
			Title:     "系统管理",
			Icon:      &componentStr,
			Path:      "/system",
			Component: "Layout",
			Redirect:  &systemUserStr,
			Sort:      99,
			ParentID:  &uint0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 2},
			Name:      "Admin",
			Title:     "管理员管理",
			Icon:      &userStr,
			Path:      "admin",
			Component: "/system/admin/index",
			Sort:      11,
			ParentID:  &uint1,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 3},
			Name:      "Role",
			Title:     "角色管理",
			Icon:      &peoplesStr,
			Path:      "role",
			Component: "/system/role/index",
			Sort:      12,
			ParentID:  &uint1,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 4},
			Name:      "Menu",
			Title:     "菜单管理",
			Icon:      &treeTableStr,
			Path:      "menu",
			Component: "/system/menu/index",
			Sort:      13,
			ParentID:  &uint1,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 5},
			Name:      "Api",
			Title:     "接口管理",
			Icon:      &treeStr,
			Path:      "api",
			Component: "/system/api/index",
			Sort:      14,
			ParentID:  &uint1,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 6},
			Name:      "Log",
			Title:     "日志管理",
			Icon:      &xitongrizhi,
			Path:      "/log",
			Component: "Layout",
			Redirect:  &logOperationStr,
			Sort:      98,
			ParentID:  &uint0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 7},
			Name:      "OperationLog",
			Title:     "操作日志",
			Icon:      &skill,
			Path:      "operation-log",
			Component: "/log/operation-log/index",
			Sort:      21,
			ParentID:  &uint6,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 8},
			Name:      "Business",
			Title:     "业务管理",
			Icon:      &yewu,
			Path:      "/business",
			Component: "Layout",
			Sort:      1,
			ParentID:  &uint0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 9},
			Name:      "Project",
			Title:     "项目管理",
			Icon:      &xiangmu,
			Path:      "/business/project",
			Component: "/business/project/index",
			Sort:      33,
			ParentID:  &uint8,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 10},
			Name:      "Content",
			Title:     "内容管理",
			Icon:      &education,
			Path:      "/content",
			Component: "Layout",
			Sort:      2,
			ParentID:  &uint0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 11},
			Name:      "Tag",
			Title:     "标签管理",
			Icon:      &tagIcon,
			Path:      "/content/tag",
			Component: "/content/tag/index",
			Sort:      33,
			ParentID:  &uint10,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 12},
			Name:      "Category",
			Title:     "分类管理",
			Icon:      &nested,
			Path:      "/content/category",
			Component: "/content/category/index",
			Sort:      33,
			ParentID:  &uint10,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 13},
			Name:      "Article",
			Title:     "文章管理",
			Icon:      &language,
			Path:      "/content/article",
			Component: "/content/article/index",
			Sort:      33,
			ParentID:  &uint10,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 14},
			Name:      "Config",
			Title:     "数据管理",
			Icon:      &shuju,
			Path:      "/content/config",
			Component: "/content/config/index",
			Sort:      33,
			ParentID:  &uint10,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 15},
			Name:      "Resource",
			Title:     "资源管理",
			Icon:      &ziyuan,
			Path:      "/content/resource",
			Component: "/content/resource/index",
			Sort:      33,
			ParentID:  &uint10,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 16},
			Name:      "User",
			Title:     "用户管理",
			Icon:      &userStr,
			Path:      "/business/user",
			Component: "/business/user/index",
			Sort:      33,
			ParentID:  &uint8,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 17},
			Name:      "Column",
			Title:     "专栏管理",
			Icon:      &documentationStr,
			Path:      "/content/column",
			Component: "/content/column/index",
			Sort:      34,
			ParentID:  &uint10,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 18},
			Name:      "Operations",
			Title:     "运营管理",
			Icon:      &yunyingzhongxin,
			Path:      "/operations",
			Component: "Layout",
			Sort:      999,
			ParentID:  &uint0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 19},
			Name:      "promotion",
			Title:     "广告位管理",
			Icon:      &eye,
			Path:      "/operation/promotion",
			Component: "/operation/promotion/index",
			Sort:      35,
			ParentID:  &uint18,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 20},
			Name:      "comment",
			Title:     "评论管理",
			Icon:      &message,
			Path:      "/operation/comment",
			Component: "/operation/comment/index",
			Sort:      999,
			ParentID:  &uint18,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 21},
			Name:      "point",
			Title:     "积分管理",
			Icon:      &jifen,
			Path:      "/operation/point",
			Component: "/operation/point/index",
			Sort:      999,
			ParentID:  &uint18,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 22},
			Name:      "Store",
			Title:     "商城管理",
			Icon:      &shopping,
			Path:      "/store",
			Component: "Layout",
			Sort:      999,
			ParentID:  &uint0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 23},
			Name:      "ProductCategory",
			Title:     "商品分类",
			Icon:      &list,
			Path:      "/store/product-category",
			Component: "/store/product-category/index",
			Sort:      999,
			ParentID:  &uint22,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 24},
			Name:      "ProductType",
			Title:     "商品类型",
			Icon:      &theme,
			Path:      "/store/product-type",
			Component: "/store/product-type/index",
			Sort:      1,
			ParentID:  &uint22,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 25},
			Name:      "Spec",
			Title:     "规格管理",
			Icon:      &guige,
			Path:      "/store/product-spec",
			Component: "/store/product-spec/index",
			Sort:      1,
			ParentID:  &uint22,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 26},
			Name:      "Product",
			Title:     "商品列表",
			Icon:      &shangpinliebiao,
			Path:      "/store/product",
			Component: "/store/product/index",
			Sort:      1,
			ParentID:  &uint22,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 27},
			Name:      "Order",
			Title:     "订单管理",
			Icon:      &shangpin,
			Path:      "/order",
			Component: "Layout",
			Sort:      1,
			ParentID:  &uint0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 28},
			Name:      "OrderList",
			Title:     "订单列表",
			Icon:      &dingdanliebiao,
			Path:      "/store/order",
			Component: "/store/order/index",
			Sort:      999,
			ParentID:  &uint27,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 29},
			Name:      "AdminConfig",
			Title:     "后台配置",
			Icon:      &peizhishezhi,
			Path:      "/system/config",
			Component: "/system/config/index",
			Sort:      1,
			ParentID:  &uint1,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
	}

	for _, menu := range menus {
		if err := createIfNotExists(db, &menu, menu.ID); err != nil {
			return err
		}
	}

	return nil
}
