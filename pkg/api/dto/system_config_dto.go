// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
)

type SystemConfigDto struct {
	SystemConfigID string `json:"systemConfigID"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	Logo           string `json:"logo"`
	Icon           string `json:"icon"`
	Footer         string `json:"footer"`
}

func ToSystemConfigInfoDto(systemConfig *model.SystemConfig) SystemConfigDto {
	return SystemConfigDto{
		SystemConfigID: systemConfig.SystemConfigID,
		Title:          systemConfig.Title,
		Content:        systemConfig.Content,
		Logo:           systemConfig.Logo,
		Icon:           systemConfig.Icon,
		Footer:         systemConfig.Footer,
	}
}
