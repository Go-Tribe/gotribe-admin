// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package repository

import (
	"context"
	"fmt"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/vo"

	"strings"
	"time"
)

type IOperationLogRepository interface {
	GetOperationLogs(ctx context.Context, req *vo.OperationLogListRequest) ([]model.OperationLog, int64, error)
	BatchDeleteOperationLogByIds(ctx context.Context, ids []uint) error
	SaveOperationLogChannel(olc <-chan *model.OperationLog) //处理OperationLogChan将日志记录到数据库
}

type OperationLogRepository struct {
}

func NewOperationLogRepository() IOperationLogRepository {
	return OperationLogRepository{}
}

func buildOperationLogOrder(req *vo.OperationLogListRequest) string {
	sortByMap := map[string]string{
		"username":   "username",
		"ip":         "ip",
		"path":       "path",
		"status":     "status",
		"startTime":  "start_time",
		"start_time": "start_time",
		"timeCost":   "time_cost",
		"time_cost":  "time_cost",
		"desc":       "\"desc\"",
	}

	column, ok := sortByMap[strings.TrimSpace(req.SortBy)]
	if !ok {
		return "start_time DESC"
	}

	direction := "ASC"
	if strings.EqualFold(strings.TrimSpace(req.SortOrder), "desc") {
		direction = "DESC"
	}

	return fmt.Sprintf("%s %s", column, direction)
}

func (o OperationLogRepository) GetOperationLogs(ctx context.Context, req *vo.OperationLogListRequest) ([]model.OperationLog, int64, error) {
	var list []model.OperationLog
	db := common.WithContext(ctx).DB().Model(&model.OperationLog{}).Order(buildOperationLogOrder(req))

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	ip := strings.TrimSpace(req.Ip)
	if ip != "" {
		db = db.Where("ip LIKE ?", fmt.Sprintf("%%%s%%", ip))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}

	// 分页
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	return list, total, err

}

func (o OperationLogRepository) BatchDeleteOperationLogByIds(ctx context.Context, ids []uint) error {
	err := common.WithContext(ctx).DB().Where("id IN (?)", ids).Unscoped().Delete(&model.OperationLog{}).Error
	return err
}

// var Logs []model.OperationLog //全局变量多个线程需要加锁，所以每个线程自己维护一个
// 处理OperationLogChan将日志记录到数据库
func (o OperationLogRepository) SaveOperationLogChannel(olc <-chan *model.OperationLog) {
	logs := make([]model.OperationLog, 0, 10)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	flush := func() {
		if len(logs) == 0 {
			return
		}
		if err := common.DB.Create(&logs).Error; err != nil {
			common.Log.Errorf("批量写入操作日志失败: %v", err)
		}
		logs = logs[:0]
	}

	for {
		select {
		case log, ok := <-olc:
			if !ok {
				flush()
				return
			}
			logs = append(logs, *log)
			if len(logs) >= 10 {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}
