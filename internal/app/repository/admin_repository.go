// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"fmt"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
	"gotribe-admin/pkg/api/vo"
	"gotribe-admin/pkg/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/thoas/go-funk"
)

type IAdminRepository interface {
	Login(admin *model.Admin) (*model.Admin, error)    // 登录
	ChangePwd(username string, newPasswd string) error // 更新密码

	CreateAdmin(admin *model.Admin) error                              // 创建用户
	GetAdminByID(id uint) (model.Admin, error)                         // 获取单个用户
	GetAdmins(req *vo.AdminListRequest) ([]*model.Admin, int64, error) // 获取用户列表
	UpdateAdmin(admin *model.Admin) error                              // 更新用户
	BatchDeleteAdminByIds(ids []uint) error                            // 批量删除

	GetCurrentAdmin(c *gin.Context) (model.Admin, error)                  // 获取当前登录用户信息
	GetCurrentAdminMinRoleSort(c *gin.Context) (uint, model.Admin, error) // 获取当前用户角色排序最小值（最高等级角色）以及当前用户信息
	GetAdminMinRoleSortsByIds(ids []uint) ([]int, error)                  // 根据用户ID获取用户角色排序最小值

	SetAdminInfoCache(username string, admin model.Admin) // 设置用户信息缓存
	UpdateAdminInfoCacheByRoleID(roleID uint) error       // 根据角色ID更新拥有该角色的用户信息缓存
	ClearAdminInfoCache()                                 // 清理所有用户信息缓存
}

type AdminRepository struct {
}

// 当前用户信息缓存，避免频繁获取数据库
var adminInfoCache = cache.New(24*time.Hour, 48*time.Hour)

// AdminRepository构造函数
func NewAdminRepository() IAdminRepository {
	return AdminRepository{}
}

// 登录
func (ar AdminRepository) Login(admin *model.Admin) (*model.Admin, error) {
	// 根据用户名获取用户(正常状态:用户状态正常)
	var firstAdmin model.Admin
	err := common.DB.
		Where("username = ?", admin.Username).
		Preload("Roles").
		First(&firstAdmin).Error
	if err != nil {
		return nil, common.ErrUserNotFound
	}

	// 判断用户的状态
	adminStatus := firstAdmin.Status
	if adminStatus != 1 {
		return nil, common.ErrUserDisabled
	}

	// 判断用户拥有的所有角色的状态,全部角色都被禁用则不能登录
	roles := firstAdmin.Roles
	isValidate := false
	for _, role := range roles {
		// 有一个正常状态的角色就可以登录
		if role.Status == known.DEFAULT_ID {
			isValidate = true
			break
		}
	}

	if !isValidate {
		return nil, common.ErrUserRoleDisabled
	}

	// 校验密码
	err = util.PasswordUtil.ComparePasswd(firstAdmin.Password, admin.Password)
	if err != nil {
		return &firstAdmin, common.ErrPasswordIncorrect
	}
	return &firstAdmin, nil
}

// 获取当前登录用户信息
// 需要缓存，减少数据库访问
func (ar AdminRepository) GetCurrentAdmin(c *gin.Context) (model.Admin, error) {
	var newAdmin model.Admin
	ctxAdmin, exist := c.Get("user")
	if !exist {
		return newAdmin, common.ErrUserNotLoggedIn
	}
	u, _ := ctxAdmin.(model.Admin)

	// 先获取缓存
	cacheAdmin, found := adminInfoCache.Get(u.Username)
	var admin model.Admin
	var err error
	if found {
		admin = cacheAdmin.(model.Admin)
		err = nil
	} else {
		// 缓存中没有就获取数据库
		admin, err = ar.GetAdminByID(u.ID)
		// 获取成功就缓存
		if err != nil {
			adminInfoCache.Delete(u.Username)
		} else {
			adminInfoCache.Set(u.Username, admin, cache.DefaultExpiration)
		}
	}
	return admin, err
}

// 获取当前用户角色排序最小值（最高等级角色）以及当前用户信息
func (ar AdminRepository) GetCurrentAdminMinRoleSort(c *gin.Context) (uint, model.Admin, error) {
	// 获取当前用户
	ctxAdmin, err := ar.GetCurrentAdmin(c)
	if err != nil {
		return 999, ctxAdmin, err
	}
	// 获取当前用户的所有角色
	currentRoles := ctxAdmin.Roles
	// 获取当前用户角色的排序，和前端传来的角色排序做比较
	var currentRoleSorts []int
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
	}
	// 当前用户角色排序最小值（最高等级角色）
	currentRoleSortMin := uint(funk.MinInt(currentRoleSorts))

	return currentRoleSortMin, ctxAdmin, nil
}

// 获取单个用户
func (ar AdminRepository) GetAdminByID(id uint) (model.Admin, error) {
	var admin model.Admin
	err := common.DB.Where("id = ?", id).Preload("Roles").First(&admin).Error
	return admin, err
}

// 获取用户列表
func (ar AdminRepository) GetAdmins(req *vo.AdminListRequest) ([]*model.Admin, int64, error) {
	var list []*model.Admin
	db := common.DB.Model(&model.Admin{}).Order("created_at DESC")

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Preload("Roles").Find(&list).Error
	} else {
		err = db.Preload("Roles").Find(&list).Error
	}
	return list, total, err
}

// 更新密码
func (ar AdminRepository) ChangePwd(username string, hashNewPasswd string) error {
	err := common.DB.Model(&model.Admin{}).Where("username = ?", username).Update("password", hashNewPasswd).Error
	// 如果更新密码成功，则更新当前用户信息缓存
	// 先获取缓存
	cacheAdmin, found := adminInfoCache.Get(username)
	if err == nil {
		if found {
			admin := cacheAdmin.(model.Admin)
			admin.Password = hashNewPasswd
			adminInfoCache.Set(username, admin, cache.DefaultExpiration)
		} else {
			// 没有缓存就获取用户信息缓存
			var admin model.Admin
			common.DB.Where("username = ?", username).First(&admin)
			adminInfoCache.Set(username, admin, cache.DefaultExpiration)
		}
	}

	return err
}

// 创建用户
func (ar AdminRepository) CreateAdmin(admin *model.Admin) error {
	err := common.DB.Create(admin).Error
	return err
}

// 更新用户
func (ar AdminRepository) UpdateAdmin(admin *model.Admin) error {
	err := common.DB.Model(admin).Updates(admin).Error
	if err != nil {
		return err
	}
	err = common.DB.Model(admin).Association("Roles").Replace(admin.Roles)

	//err := common.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&admin).Error

	// 如果更新成功就更新用户信息缓存
	if err == nil {
		adminInfoCache.Set(admin.Username, *admin, cache.DefaultExpiration)
	}
	return err
}

// 批量删除
func (ar AdminRepository) BatchDeleteAdminByIds(ids []uint) error {
	// 用户和角色存在多对多关联关系
	var admins []model.Admin
	for _, id := range ids {
		// 根据ID获取用户
		admin, err := ar.GetAdminByID(id)
		if err != nil {
			return common.NewUserNotFoundByIDError(id)
		}
		admins = append(admins, admin)
	}

	err := common.DB.Select("Roles").Unscoped().Delete(&admins).Error
	// 删除用户成功，则删除用户信息缓存
	if err == nil {
		for _, admin := range admins {
			adminInfoCache.Delete(admin.Username)
		}
	}
	return err
}

// 根据用户ID获取用户角色排序最小值
func (ar AdminRepository) GetAdminMinRoleSortsByIds(ids []uint) ([]int, error) {
	// 根据用户ID获取用户信息
	var adminList []model.Admin
	err := common.DB.Where("id IN (?)", ids).Preload("Roles").Find(&adminList).Error
	if err != nil {
		return []int{}, err
	}
	if len(adminList) == 0 {
		return []int{}, common.ErrNoUsersFound
	}
	var roleMinSortList []int
	for _, admin := range adminList {
		roles := admin.Roles
		var roleSortList []int
		for _, role := range roles {
			roleSortList = append(roleSortList, int(role.Sort))
		}
		roleMinSort := funk.MinInt(roleSortList)
		roleMinSortList = append(roleMinSortList, roleMinSort)
	}
	return roleMinSortList, nil
}

// 设置用户信息缓存
func (ar AdminRepository) SetAdminInfoCache(username string, admin model.Admin) {
	adminInfoCache.Set(username, admin, cache.DefaultExpiration)
}

// 根据角色ID更新拥有该角色的用户信息缓存
func (ar AdminRepository) UpdateAdminInfoCacheByRoleID(roleID uint) error {

	var role model.Role
	err := common.DB.Where("id = ?", roleID).Preload("Admins").First(&role).Error
	if err != nil {
		return common.ErrRoleInfoFailed
	}

	admins := role.Admin
	if len(admins) == 0 {
		return common.ErrNoUsersWithRole
	}

	// 更新用户信息缓存
	for _, admin := range admins {
		_, found := adminInfoCache.Get(admin.Username)
		if found {
			adminInfoCache.Set(admin.Username, *admin, cache.DefaultExpiration)
		}
	}

	return err
}

// 清理所有用户信息缓存
func (ar AdminRepository) ClearAdminInfoCache() {
	adminInfoCache.Flush()
}
