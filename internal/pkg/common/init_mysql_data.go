// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package common

import (
	"errors"
	"github.com/dengmengmian/ghelper/gid"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"gotribe-admin/config"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/util"
)

// 初始化mysql数据
func InitData() {
	// 是否初始化数据
	if !config.Conf.System.InitData {
		return
	}

	// 1.写入角色数据
	newRoles := make([]*model.Role, 0)
	roles := []*model.Role{
		{
			Model:   model.Model{ID: 1},
			Name:    "管理员",
			Keyword: "admin",
			Desc:    new(string),
			Sort:    1,
			Status:  1,
			Creator: "系统",
		},
		{
			Model:   model.Model{ID: 2},
			Name:    "普通管理员",
			Keyword: "user",
			Desc:    new(string),
			Sort:    3,
			Status:  1,
			Creator: "系统",
		},
		{
			Model:   model.Model{ID: 3},
			Name:    "访客",
			Keyword: "guest",
			Desc:    new(string),
			Sort:    5,
			Status:  1,
			Creator: "系统",
		},
	}

	for _, role := range roles {
		err := DB.First(&role, role.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newRoles = append(newRoles, role)
		}
	}

	if len(newRoles) > 0 {
		err := DB.Create(&newRoles).Error
		if err != nil {
			Log.Errorf("写入系统角色数据失败：%v", err)
		}
	}

	// 2写入菜单
	newMenus := make([]model.Menu, 0)
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:2],
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
			Roles:     roles[:2],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
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
			Roles:     roles[:1],
			Creator:   "系统",
		},
	}
	for _, menu := range menus {
		err := DB.First(&menu, menu.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newMenus = append(newMenus, menu)
		}
	}
	if len(newMenus) > 0 {
		err := DB.Create(&newMenus).Error
		if err != nil {
			Log.Errorf("写入系统菜单数据失败：%v", err)
		}
	}

	// 3.写入管理员
	newAdmins := make([]model.Admin, 0)
	admins := []model.Admin{
		{
			Model:        model.Model{ID: 1},
			Username:     "admin",
			Password:     util.GenPasswd("123456"),
			Mobile:       "18888888888",
			Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
			Nickname:     new(string),
			Introduction: new(string),
			Status:       1,
			Creator:      "系统",
			Roles:        roles[:1],
		},
	}

	for _, admin := range admins {
		err := DB.First(&admin, admin.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newAdmins = append(newAdmins, admin)
		}
	}

	if len(newAdmins) > 0 {
		err := DB.Create(&newAdmins).Error
		if err != nil {
			Log.Errorf("写入管理员数据失败：%v", err)
		}
	}

	// 4.写入api
	apis := []model.Api{
		{
			Method:   "POST",
			Path:     "/base/login",
			Category: "base",
			Desc:     "管理员登录",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/base/logout",
			Category: "base",
			Desc:     "管理员登出",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/base/refreshToken",
			Category: "base",
			Desc:     "刷新JWT令牌",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/admin/info",
			Category: "admin",
			Desc:     "获取当前登录管理员信息",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/admin/list",
			Category: "admin",
			Desc:     "获取管理员列表",
			Creator:  "系统",
		},
		{
			Method:   "PUT",
			Path:     "/admin/changePwd",
			Category: "admin",
			Desc:     "更新管理员登录密码",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/admin/create",
			Category: "admin",
			Desc:     "创建管理员",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/admin/update/:userID",
			Category: "admin",
			Desc:     "更新管理员",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/admin/delete/batch",
			Category: "admin",
			Desc:     "批量删除管理员",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/role/list",
			Category: "role",
			Desc:     "获取角色列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/role/create",
			Category: "role",
			Desc:     "创建角色",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/role/update/:roleID",
			Category: "role",
			Desc:     "更新角色",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/role/menus/get/:roleID",
			Category: "role",
			Desc:     "获取角色的权限菜单",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/role/menus/update/:roleID",
			Category: "role",
			Desc:     "更新角色的权限菜单",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/role/apis/get/:roleID",
			Category: "role",
			Desc:     "获取角色的权限接口",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/role/apis/update/:roleID",
			Category: "role",
			Desc:     "更新角色的权限接口",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/role/delete/batch",
			Category: "role",
			Desc:     "批量删除角色",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/menu/list",
			Category: "menu",
			Desc:     "获取菜单列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/menu/tree",
			Category: "menu",
			Desc:     "获取菜单树",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/menu/create",
			Category: "menu",
			Desc:     "创建菜单",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/menu/update/:menuID",
			Category: "menu",
			Desc:     "更新菜单",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/menu/delete/batch",
			Category: "menu",
			Desc:     "批量删除菜单",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/menu/access/list/:userID",
			Category: "menu",
			Desc:     "获取管理员的可访问菜单列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/menu/access/tree/:userID",
			Category: "menu",
			Desc:     "获取管理员的可访问菜单树",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/api/list",
			Category: "api",
			Desc:     "获取接口列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/api/tree",
			Category: "api",
			Desc:     "获取接口树",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/api/create",
			Category: "api",
			Desc:     "创建接口",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/api/update/:roleID",
			Category: "api",
			Desc:     "更新接口",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/api/delete/batch",
			Category: "api",
			Desc:     "批量删除接口",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/log/operation/list",
			Category: "log",
			Desc:     "获取操作日志列表",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/log/operation/delete/batch",
			Category: "log",
			Desc:     "批量删除操作日志",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/project/:projectID",
			Category: "project",
			Desc:     "获取单条项目详情",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/project",
			Category: "project",
			Desc:     "获取项目列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/project",
			Category: "project",
			Desc:     "创建项目",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/project/:projectID",
			Category: "project",
			Desc:     "更新项目",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/project",
			Category: "project",
			Desc:     "批量删除项目",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/config/:configID",
			Category: "config",
			Desc:     "获取单条配置详情",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/config",
			Category: "config",
			Desc:     "获取配置列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/config",
			Category: "config",
			Desc:     "创建配置",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/config/:configID",
			Category: "config",
			Desc:     "更新配置",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/config",
			Category: "config",
			Desc:     "批量删除配置",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/tag/:tagID",
			Category: "tag",
			Desc:     "获取单条标签详情",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/tag",
			Category: "tag",
			Desc:     "获取标签列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/tag",
			Category: "tag",
			Desc:     "创建标签",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/tag/:tagID",
			Category: "tag",
			Desc:     "更新标签",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/tag",
			Category: "tag",
			Desc:     "批量删除标签",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/category/:categoryID",
			Category: "category",
			Desc:     "获取分类信息",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/category/tree",
			Category: "category",
			Desc:     "获取分类树",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/category",
			Category: "category",
			Desc:     "获取分类列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/category",
			Category: "category",
			Desc:     "创建分类",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/category/:categoryID",
			Category: "category",
			Desc:     "更新分类",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/category",
			Category: "category",
			Desc:     "批量删除分类",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/post/:postID",
			Category: "post",
			Desc:     "获取单条内容详情",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/post",
			Category: "post",
			Desc:     "获取内容列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/post",
			Category: "post",
			Desc:     "创建内容",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/post/:postID",
			Category: "post",
			Desc:     "更新内容",
			Creator:  "系统",
		},
		{
			Method:   "PUT",
			Path:     "/post/:postID",
			Category: "post",
			Desc:     "发布内容",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/post",
			Category: "post",
			Desc:     "批量删除内容",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/user/:userID",
			Category: "user",
			Desc:     "获取单个用户详情",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/user",
			Category: "user",
			Desc:     "获取用户列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/user/search",
			Category: "user",
			Desc:     "搜索用户列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/user",
			Category: "user",
			Desc:     "创建用户",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/user/:userID",
			Category: "user",
			Desc:     "更新用户",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/user",
			Category: "user",
			Desc:     "批量删除用户",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/resource/upload",
			Category: "resource",
			Desc:     "上传资源",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/resource",
			Category: "resource",
			Desc:     "获取资源列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/resource/:resourceID",
			Category: "resource",
			Desc:     "获取资源详情",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/resource/:resourceID",
			Category: "resource",
			Desc:     "更新资源信息",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/resource",
			Category: "resource",
			Desc:     "删除资源",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/column",
			Category: "column",
			Desc:     "新增专栏",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/column",
			Category: "column",
			Desc:     "获取专栏列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/column/:columnID",
			Category: "column",
			Desc:     "获取专栏详情",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/column/:columnID",
			Category: "column",
			Desc:     "更新专栏信息",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/column",
			Category: "column",
			Desc:     "删除专栏",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/ad/scene/:adSceneID",
			Category: "ad_scene",
			Desc:     "获取单条推广场景",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/ad/scene",
			Category: "ad_scene",
			Desc:     "获取所有推广场景",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/ad/scene",
			Category: "ad_scene",
			Desc:     "创建推广场景",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/ad/scene",
			Category: "ad_scene",
			Desc:     "删除推广场景",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/ad/scene/:adSceneID",
			Category: "ad_scene",
			Desc:     "更新推广场景信息",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/ad/:adID",
			Category: "ad",
			Desc:     "获取单个广告位",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/ad",
			Category: "ad",
			Desc:     "获取广告列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/ad",
			Category: "ad",
			Desc:     "创建广告",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/ad/:adID",
			Category: "ad",
			Desc:     "更新广告",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/ad",
			Category: "ad",
			Desc:     "删除广告",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/comment",
			Category: "comment",
			Desc:     "获取评论列表",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/comment/:commentID",
			Category: "comment",
			Desc:     "审核评论",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/point",
			Category: "point",
			Desc:     "获取评论列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/point",
			Category: "point",
			Desc:     "后台增加积分",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/product/category/tree",
			Category: "product_category",
			Desc:     "获取分类树",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/product/category",
			Category: "product_category",
			Desc:     "获取分类列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/product/category",
			Category: "product_category",
			Desc:     "创建分类",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/product/category/:productCategoryID",
			Category: "product_category",
			Desc:     "更新分类信息",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/product/category",
			Category: "product_category",
			Desc:     "删除分类",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/product/category/:productCategoryID",
			Category: "product_category",
			Desc:     "获取分类详情",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/product/type/:productTypeID",
			Category: "product_type",
			Desc:     "获取商品类型详情",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/product/type",
			Category: "product_type",
			Desc:     "获取商品类型列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/product/type",
			Category: "product_type",
			Desc:     "创建商品类型",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/product/type/:productTypeID",
			Category: "product_type",
			Desc:     "更新商品类型信息",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/product/type",
			Category: "product_type",
			Desc:     "删除商品类型",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/product/spec/:productSpecID",
			Category: "product_spec",
			Desc:     "获取商品规格",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/product/spec",
			Category: "product_spec",
			Desc:     "获取商品规格列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/product/spec",
			Category: "product_spec",
			Desc:     "创建商品规格",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/product/spec/:productSpecID",
			Category: "product_spec",
			Desc:     "更新商品规格",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/product/spec",
			Category: "product_spec",
			Desc:     "删除商品规格",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/product/spec/item/:productSpecItemID",
			Category: "product_spec_item",
			Desc:     "获取商品规格值",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/api/product/spec/item",
			Category: "product_spec_item",
			Desc:     "获取商品规格值列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/product/spec/item",
			Category: "product_spec_item",
			Desc:     "创建商品规格值",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/product/spec/item/:productSpecItemID",
			Category: "product_spec_item",
			Desc:     "更新商品规格值",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/product/spec/item",
			Category: "product_spec_item",
			Desc:     "删除商品规格值",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/product/:productID",
			Category: "product",
			Desc:     "获取商品信息",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/product",
			Category: "product",
			Desc:     "获取商品列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/product",
			Category: "product",
			Desc:     "创建商品",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/product/:productID",
			Category: "product",
			Desc:     "更新商品",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/product",
			Category: "product",
			Desc:     "删除商品",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/product/spec/info/:categoryID",
			Category: "product_spec",
			Desc:     "获取规格和规格项信息",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/order/:orderID",
			Category: "order",
			Desc:     "获取订单详情",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/order/log/:orderID",
			Category: "order",
			Desc:     "获取订单记录",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/order",
			Category: "order",
			Desc:     "获取订单列表",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/order/:orderID",
			Category: "order",
			Desc:     "修改订单信息",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/order",
			Category: "order",
			Desc:     "删除订单",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/base/config",
			Category: "base",
			Desc:     "获取后台配置",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/system",
			Category: "system",
			Desc:     "更新配置",
			Creator:  "系统",
		},
	}
	newApi := make([]model.Api, 0)
	newRoleCasbin := make([]model.RoleCasbin, 0)
	for i, api := range apis {
		api.ID = uint(i + 1)
		err := DB.First(&api, api.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newApi = append(newApi, api)

			// 管理员拥有所有API权限
			newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{
				Keyword: roles[0].Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})

			// 非管理员拥有基础权限
			basePaths := []string{
				"/base/login",
				"/base/logout",
				"/base/refreshToken",
				"/user/info",
				"/base/config",
				"/menu/access/tree/:userID",
			}

			if funk.ContainsString(basePaths, api.Path) {
				newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{
					Keyword: roles[1].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
				newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{
					Keyword: roles[2].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
			}
		}
	}

	if len(newApi) > 0 {
		if err := DB.Create(&newApi).Error; err != nil {
			Log.Errorf("写入api数据失败：%v", err)
		}
	}

	if len(newRoleCasbin) > 0 {
		rules := make([][]string, 0)
		for _, c := range newRoleCasbin {
			rules = append(rules, []string{
				c.Keyword, c.Path, c.Method,
			})
		}
		isAdd, err := CasbinEnforcer.AddPolicies(rules)
		if !isAdd {
			Log.Errorf("写入casbin数据失败：%v", err)
		}
	}
	// 3.写入分类
	newCategory := make([]model.Category, 0)
	categorys := []model.Category{
		{
			Model:       model.Model{ID: 1},
			Title:       "默认分类",
			Description: "默认分类",
			CategoryID:  "24ejga",
			Status:      1,
		},
	}

	for _, category := range categorys {
		err := DB.First(&category, category.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newCategory = append(newCategory, category)
		}
	}

	if len(newCategory) > 0 {
		err := DB.Create(&newCategory).Error
		if err != nil {
			Log.Errorf("写入分类数据失败：%v", err)
		}
	}

	// 6.写入TAG数据
	newTags := make([]model.Tag, 0)
	tags := []model.Tag{
		{
			Model:       model.Model{ID: 1},
			TagID:       gid.GenShortID(),
			Title:       "默认标签",
			Description: "默认标签",
		},
	}

	for _, tag := range tags {
		err := DB.First(&tag, tag.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newTags = append(newTags, tag)
		}
	}

	if len(newTags) > 0 {
		err := DB.Create(&newTags).Error
		if err != nil {
			Log.Errorf("写入tag数据失败：%v", err)
		}
	}

	// 7.写入项目数据
	newProjects := make([]model.Project, 0)
	projects := []model.Project{
		{
			Model:       model.Model{ID: 1},
			Title:       "默认项目",
			ProjectID:   "245eko",
			Description: "默认项目",
		},
	}

	for _, project := range projects {
		err := DB.First(&project, project.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newProjects = append(newProjects, project)
		}
	}

	if len(newProjects) > 0 {
		err := DB.Create(&newProjects).Error
		if err != nil {
			Log.Errorf("写入项目数据失败：%v", err)
		}
	}
	// 8.初始化nav
	newNavs := make([]model.Config, 0)
	navs := []model.Config{
		{
			Model:       model.Model{ID: 1},
			Type:        2,
			ProjectID:   "245eko",
			Alias:       "nav",
			Description: "导航配置",
			Title:       "默认导航",
			Info:        `[{"name":"首页","url":"/","child":[]},{"name":"精选栏目","url":"/series","child":[]},{"name":"示例页面","url":"/p/24g1i6","child":[]}]`,
		},
	}

	for _, nav := range navs {
		err := DB.First(&nav, nav.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newNavs = append(newNavs, nav)
		}
	}

	if len(newNavs) > 0 {
		err := DB.Create(&newNavs).Error
		if err != nil {
			Log.Errorf("写入nav数据失败：%v", err)
		}
	}

	// 9.初始化 user
	newUsers := make([]model.User, 0)
	users := []model.User{
		{
			Model:     model.Model{ID: 1},
			UserID:    gid.GenShortID(),
			Username:  "gotribe",
			Nickname:  "gotribe",
			ProjectID: "245eko",
		},
	}

	for _, user := range users {
		err := DB.First(&user, user.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUsers = append(newUsers, user)
		}
	}

	if len(newUsers) > 0 {
		err := DB.Create(&newUsers).Error
		if err != nil {
			Log.Errorf("写入用户数据失败：%v", err)
		}
	}

	// 10.写入文章数据
	newPosts := make([]model.Post, 0)
	posts := []model.Post{
		{
			Model:       model.Model{ID: 1},
			PostID:      "243x9g",
			Title:       "欢迎使用GoTribe",
			Description: "这是一篇示例文章",
			Content:     "# 这是一篇示例文章",
			Icon:        "https://cdn.dengmengmian.com/20240528/1716909013037462047.jpg",
			HtmlContent: "<h1>这是一篇示例文章</h1>",
			UserID:      "245eko",
			CategoryID:  "24ejga",
			Author:      "GoTribe",
			ProjectID:   "245eko",
		},
	}

	for _, post := range posts {
		err := DB.First(&post, post.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newPosts = append(newPosts, post)
		}
	}

	if len(newPosts) > 0 {
		err := DB.Create(&newPosts).Error
		if err != nil {
			Log.Errorf("写入文章数据失败：%v", err)
		}
	}

}
