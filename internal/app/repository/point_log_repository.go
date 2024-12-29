// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"errors"
	"fmt"
	"github.com/dengmengmian/ghelper/gconvert"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"
	"strings"
	"time"
)

type IPointLogRepository interface {
	CreatePoint(userID, types, reason, eventID, ProjectID string, points float64) error // 新增积分
	GetPointLogs(req *vo.PointLogListRequest) ([]*model.PointLog, int64, error)         // 获取积分列表
}

type PointLogRepository struct {
}

// PointLogRepository构造函数
func NewPointLogRepository() IPointLogRepository {
	return PointLogRepository{}
}

// 获取推广场景列表
func (cr PointLogRepository) GetPointLogs(req *vo.PointLogListRequest) ([]*model.PointLog, int64, error) {
	var list []*model.PointLog
	db := common.DB.Model(&model.PointLog{}).Order("created_at DESC")

	projectID := strings.TrimSpace(req.ProjectID)
	if !gconvert.IsEmpty(projectID) {
		db = db.Where("project_id = ?", fmt.Sprintf("%s", projectID))
	}
	if !gconvert.IsEmpty(req.UserID) {
		db = db.Where("user_id =  ?", fmt.Sprintf("%s", req.UserID))
	}
	if !gconvert.IsEmpty(req.Nickname) {
		// 查出用户 ID。再用用户 ID 去筛选
		var user model.User
		if result := common.DB.Model(&model.User{}).Where("nickname like ?", fmt.Sprintf("%%%s%%", req.Nickname)).First(&user); result.Error != nil {
			return nil, 0, errors.New("用户不存在")
		}
		db = db.Where("user_id = ?", user.UserID)
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
	return GetPointLogOther(list), total, err
}

// 获取推广场景其他信息
func GetPointLogOther(pointLogs []*model.PointLog) []*model.PointLog {
	for _, m := range pointLogs {
		var user *model.User
		_ = common.DB.Where("user_id = ?", m.UserID).First(&user).Error
		m.User = user
	}
	return pointLogs
}

// 创建推广场景
func (cr PointLogRepository) CreatePoint(userID, types, reason, eventID, ProjectID string, points float64) error {
	pointLog := &model.PointLog{
		UserID:    userID,
		Type:      types,
		Reason:    reason,
		EventID:   eventID,
		Points:    points,
		ProjectID: ProjectID,
	}
	result := common.DB.Create(pointLog)
	if result.Error != nil {
		return result.Error
	}
	// 新增可用积分
	userPoint := &model.PointAvailable{
		ProjectID:      ProjectID,
		UserID:         userID,
		Points:         points,
		PointsLogID:    int(pointLog.ID),
		ExpirationDate: time.Now().AddDate(1, 0, 0), // 当前时间往后推一年
	}

	err := common.DB.Create(userPoint).Error
	return err
}
