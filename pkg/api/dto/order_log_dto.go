// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/util"
)

type OrderLogDto struct {
	OrderLogID string `json:"orderLogID"`
	OrderID    string `json:"orderID"`
	Remark     string `json:"remark"`
	CreatedAt  string `json:"createdAt"`
}

func ToOrderLogInfoDto(orderLog *model.OrderLog) OrderLogDto {
	if orderLog == nil {
		return OrderLogDto{}
	}

	return OrderLogDto{
		OrderLogID: orderLog.OrderLogID,
		Remark:     orderLog.Remark,
		CreatedAt:  util.FormatTime(orderLog.CreatedAt),
	}
}

func ToOrderLogsDto(orderLogList []*model.OrderLog) []OrderLogDto {
	if orderLogList == nil {
		return nil
	}
	orderLogs := make([]OrderLogDto, 0, len(orderLogList))
	for _, orderLog := range orderLogList {
		if orderLog == nil {
			continue
		}
		orderLogDto := ToOrderLogInfoDto(orderLog)
		orderLogs = append(orderLogs, orderLogDto)
	}

	return orderLogs
}
