// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dengmengmian/ghelper/gconvert"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
)

type IPointLogRepository interface {
	CreatePoint(ctx context.Context, userID uint, types, reason, eventID, ProjectID string, points float64) error // 新增积分
	GetPointLogs(ctx context.Context, req *vo.PointLogListRequest) ([]*model.PointLog, int64, error)              // 获取积分列表
}

type PointLogRepository struct {
}

// PointLogRepository构造函数
func NewPointLogRepository() IPointLogRepository {
	return PointLogRepository{}
}

// 获取推广场景列表
func (cr PointLogRepository) GetPointLogs(ctx context.Context, req *vo.PointLogListRequest) ([]*model.PointLog, int64, error) {
	var list []*model.PointLog
	db := common.WithContext(ctx).DB().Model(&model.PointLog{}).Order("created_at DESC")

	projectID := strings.TrimSpace(req.ProjectID)
	if !gconvert.IsEmpty(projectID) {
		db = db.Where("project_id = ?", projectID)
	}
	if req.UserID > 0 {
		db = db.Where("user_id =  ?", req.UserID)
	}
	if !gconvert.IsEmpty(req.Nickname) {
		// 查出用户 ID。再用用户 ID 去筛选
		var user model.User
		if result := common.WithContext(ctx).DB().Model(&model.User{}).Where("nickname like ?", fmt.Sprintf("%%%s%%", req.Nickname)).First(&user); result.Error != nil {
			return nil, 0, common.ErrUserNotFound
		}
		db = db.Where("user_id = ?", user.ID)
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
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	list, err = GetPointLogOther(ctx, list)
	return list, total, err
}

// 获取其他信息
func GetPointLogOther(ctx context.Context, pointLogs []*model.PointLog) ([]*model.PointLog, error) {
	if len(pointLogs) == 0 {
		return pointLogs, nil
	}

	// 收集所有 UserID
	userIDs := make([]uint, 0, len(pointLogs))
	for _, m := range pointLogs {
		if m.UserID > 0 {
			userIDs = append(userIDs, m.UserID)
		}
	}

	if len(userIDs) == 0 {
		return pointLogs, nil
	}

	// 批量查询
	var users []*model.User
	if err := common.WithContext(ctx).DB().Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return pointLogs, err
	}

	// 建立映射
	userMap := make(map[uint]*model.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// 赋值
	for _, m := range pointLogs {
		if user, ok := userMap[m.UserID]; ok {
			m.User = user
		}
	}
	return pointLogs, nil
}

// 创建推广场景
func (cr PointLogRepository) CreatePoint(ctx context.Context, userID uint, types, reason, eventID, ProjectID string, points float64) error {
	// 将元转换为分
	pointsCents := int64(points * 100)

	pointLog := &model.PointLog{
		UserID:    userID,
		Type:      types,
		Reason:    reason,
		EventID:   eventID,
		Points:    pointsCents,
		ProjectID: ProjectID,
	}
	result := common.WithContext(ctx).DB().Create(pointLog)
	if result.Error != nil {
		return result.Error
	}
	// 新增可用积分
	userPoint := &model.PointAvailable{
		ProjectID:      ProjectID,
		UserID:         userID,
		Points:         pointsCents,
		PointsLogID:    pointLog.ID,
		ExpirationDate: time.Now().AddDate(1, 0, 0), // 当前时间往后推一年
	}

	err := common.WithContext(ctx).DB().Create(userPoint).Error
	return err
}
