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
	componentStr := "component"
	systemUserStr := "/system/user"
	userStr := "user"
	peoplesStr := "peoples"
	treeTableStr := "tree-table"
	treeStr := "tree"
	exampleStr := "example"
	logOperationStr := "/log/operation-log"
	documentationStr := "documentation"
	excel := "excel"
	//tab := "tab"
	documentation := "documentation"
	education := "education"
	bug := "bug"
	form := "form"
	language := "language"
	password := "password"
	chart := "chart"

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
			Icon:      &exampleStr,
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
			Icon:      &documentationStr,
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
			Icon:      &excel,
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
			Icon:      &documentation,
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
			Icon:      &bug,
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
			Icon:      &form,
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
			Icon:      &password,
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
			Icon:      &chart,
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
			ParentID:  &uint8,
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
	newUsers := make([]model.Admin, 0)
	users := []model.Admin{
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

	for _, user := range users {
		err := DB.First(&user, user.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUsers = append(newUsers, user)
		}
	}

	if len(newUsers) > 0 {
		err := DB.Create(&newUsers).Error
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

	// 5.写入分类数据
	newCategory := model.Category{}
	newCategory.ID = 1
	newCategory.Title = "默认分类"
	newCategory.Description = "默认分类"
	if err := DB.Create(&newCategory).Error; err != nil {
		Log.Errorf("写入默认分类数据失败：%v", err)
	}

	// 6.写入TAG数据
	newTag := model.Tag{}
	newTag.ID = 1
	newTag.TagID = gid.GenShortID()
	newTag.Title = "默认标签"
	newTag.Description = "默认标签"
	if err := DB.Create(&newTag).Error; err != nil {
		Log.Errorf("写入默认tag数据失败：%v", err)
	}

	// 7.写入项目数据
	newProject := model.Project{}
	newProject.ID = 1
	newProject.Title = "默认项目"
	newProject.ProjectID = "245eko"
	newProject.Description = "默认项目"
	if err := DB.Create(&newProject).Error; err != nil {
		Log.Errorf("写入项目数据失败：%v", err)
	}

	// 8.初始化nav
	nav := model.Config{}
	nav.Type = 2
	nav.ProjectID = "245eko"
	nav.Alias = "nav"
	nav.Description = "导航配置"
	nav.Title = "博客导航"
	nav.Info = `[{"name":"首页","url":"/","child":[]},{"name":"精选栏目","url":"/series","child":[]},{"name":"示例页面","url":"/p/24g1i6","child":[]}]`
	if err := DB.Create(&nav).Error; err != nil {
		Log.Errorf("写入nav数据失败：%v", err)
	}

	// 9.初始化 user
	user := model.User{}
	user.UserID = gid.GenShortID()
	user.Username = "author"
	user.Nickname = "作者"
	user.ProjectID = "245eko"
	if err := DB.Create(&user).Error; err != nil {
		Log.Errorf("写入用户数据失败：%v", err)
	}
}
