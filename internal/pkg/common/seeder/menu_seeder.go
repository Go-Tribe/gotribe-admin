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
	var uint8 uint = 8
	var uint10 uint = 10
	var uint18 uint = 18
	var uint22 uint = 22
	var uint27 uint = 27

	componentStr := "component"
	systemUserStr := "/system/user"
	tableOfContents := "TableOfContents"
	briefcaseBusiness := "BriefcaseBusiness"
	storeIcon := "Store"
	shoppingBag := "LucideShoppingBag"
	shoppingBasket := "ShoppingBasket"
	alignHorizontalJustifyEnd := "AlignHorizontalJustifyEnd"
	lucideUsers := "LucideUsers"
	lucideUserRound := "LucideUserRound"
	menuSquare := "MenuSquare"
	network := "Network"
	lucideTableProperties := "LucideTableProperties"
	lucideUserCog2 := "LucideUserCog2"
	clipboardEdit := "ClipboardEdit"
	laptopMinimalCheck := "LaptopMinimalCheck"
	lucideHeadset := "LucideHeadset"
	boxes := "Boxes"
	wholeWord := "WholeWord"
	database := "Database"
	bookImage := "BookImage"
	lucideColumnsSettings := "LucideColumnsSettings"
	bookText := "BookText"
	lucideTags := "LucideTags"
	brickWall := "BrickWall"
	lucideMartini := "LucideMartini"
	lucideDatabaseZap := "LucideDatabaseZap"
	lucideMessageSquareCode := "LucideMessageSquareCode"
	lucideTableConfig := "LucideTableConfig"
	lucideChartScatter := "LucideChartScatter"

	menus := []model.Menu{
		{
			Model:     model.Model{ID: 34},
			Name:      "Dashboard",
			Title:     "Dashboard",
			Icon:      &alignHorizontalJustifyEnd,
			Path:      "/dashboard",
			Component: "Layout",
			Sort:      1,
			ParentID:  &uint0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 1},
			Name:      "System",
			Title:     "系统管理",
			Icon:      &componentStr,
			Path:      "/system",
			Component: "Layout",
			Redirect:  &systemUserStr,
			Sort:      2,
			ParentID:  &uint0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 2},
			Name:      "Admin",
			Title:     "管理员管理",
			Icon:      &lucideUsers,
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
			Icon:      &lucideUserRound,
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
			Icon:      &menuSquare,
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
			Icon:      &network,
			Path:      "api",
			Component: "/system/api/index",
			Sort:      14,
			ParentID:  &uint1,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 7},
			Name:      "OperationLog",
			Title:     "操作日志",
			Icon:      &lucideMessageSquareCode,
			Path:      "operation-log",
			Component: "/system/operation-log/index",
			Sort:      21,
			ParentID:  &uint1,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 8},
			Name:      "Business",
			Title:     "业务管理",
			Icon:      &briefcaseBusiness,
			Path:      "/business",
			Component: "Layout",
			Sort:      3,
			ParentID:  &uint0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 9},
			Name:      "Project",
			Title:     "项目管理",
			Icon:      &lucideTableProperties,
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
			Icon:      &tableOfContents,
			Path:      "/content",
			Component: "Layout",
			Sort:      4,
			ParentID:  &uint0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 11},
			Name:      "Tag",
			Title:     "标签管理",
			Icon:      &lucideTags,
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
			Icon:      &boxes,
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
			Icon:      &wholeWord,
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
			Icon:      &database,
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
			Icon:      &bookImage,
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
			Icon:      &lucideUserCog2,
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
			Icon:      &lucideColumnsSettings,
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
			Icon:      &clipboardEdit,
			Path:      "/operations",
			Component: "Layout",
			Sort:      5,
			ParentID:  &uint0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 19},
			Name:      "Scene",
			Title:     "场景管理",
			Icon:      &laptopMinimalCheck,
			Path:      "/promotion/scene",
			Component: "/promotion/scene/index",
			Sort:      1,
			ParentID:  &uint18,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 32},
			Name:      "Advertising",
			Title:     "广告管理",
			Icon:      &lucideHeadset,
			Path:      "/promotion/advertising",
			Component: "/promotion/advertising/index",
			Sort:      2,
			ParentID:  &uint18,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 20},
			Name:      "comment",
			Title:     "评论管理",
			Icon:      &lucideMartini,
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
			Icon:      &lucideDatabaseZap,
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
			Icon:      &storeIcon,
			Path:      "/store",
			Component: "Layout",
			Sort:      5,
			ParentID:  &uint0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 23},
			Name:      "ProductCategory",
			Title:     "商品分类",
			Icon:      &shoppingBasket,
			Path:      "/store/product-category",
			Component: "/store/product-category/index",
			Sort:      1,
			ParentID:  &uint22,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 33},
			Name:      "SpecValue",
			Title:     "规格属性",
			Icon:      &bookText,
			Path:      "/store/product-spec-value",
			Component: "/store/product-spec-value/index",
			Sort:      3,
			ParentID:  &uint22,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 24},
			Name:      "ProductType",
			Title:     "商品类型",
			Icon:      &shoppingBag,
			Path:      "/store/product-type",
			Component: "/store/product-type/index",
			Sort:      4,
			ParentID:  &uint22,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 25},
			Name:      "Spec",
			Title:     "规格管理",
			Icon:      &lucideChartScatter,
			Path:      "/store/product-spec",
			Component: "/store/product-spec/index",
			Sort:      2,
			ParentID:  &uint22,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 26},
			Name:      "Product",
			Title:     "商品列表",
			Icon:      &brickWall,
			Path:      "/store/product",
			Component: "/store/product/index",
			Sort:      5,
			ParentID:  &uint22,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 27},
			Name:      "Order",
			Title:     "订单管理",
			Icon:      &shoppingBag,
			Path:      "/order",
			Component: "Layout",
			Sort:      6,
			ParentID:  &uint0,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 28},
			Name:      "OrderList",
			Title:     "订单列表",
			Icon:      &shoppingBasket,
			Path:      "/order/list",
			Component: "/order/list",
			Sort:      999,
			ParentID:  &uint27,
			Roles:     []*model.Role{&adminRole},
			Creator:   "系统",
		},
		{
			Model:     model.Model{ID: 29},
			Name:      "AdminConfig",
			Title:     "后台配置",
			Icon:      &lucideTableConfig,
			Path:      "/system/config",
			Component: "/system/config/index",
			Sort:      22,
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
