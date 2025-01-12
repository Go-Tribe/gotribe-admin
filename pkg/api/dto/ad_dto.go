// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

type AdDto struct {
	AdID        string `json:"adID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	SceneID     string `json:"sceneID"`
	SceneTitle  string `json:"SceneTitle"`
	Status      uint   `json:"status"`
	Image       string `json:"image"`
	Video       string `json:"video"`
	Sort        uint   `json:"sort"`
	URL         string `json:"url"`
	URLType     uint   `json:"urlType"`
	Ext         string `json:"ext"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// toAdDto converts a model.Ad to an AdDto.
func toAdDto(ad model.Ad) AdDto {
	var sceneTitle string
	if ad.Scene != nil {
		sceneTitle = ad.Scene.Title
	}
	return AdDto{
		AdID:        ad.AdID,
		Title:       ad.Title,
		Description: ad.Description,
		SceneID:     ad.SceneID,
		Status:      ad.Status,
		Image:       ad.Image,
		Video:       ad.Video,
		Sort:        ad.Sort,
		URL:         ad.URL,
		URLType:     ad.URLType,
		Ext:         ad.Ext,
		SceneTitle:  sceneTitle,
		CreatedAt:   ad.CreatedAt.Format(known.TIME_FORMAT),
		UpdatedAt:   ad.UpdatedAt.Format(known.TIME_FORMAT),
	}
}

// ToAdInfoDto converts a model.Ad to an AdDto.
func ToAdInfoDto(ad model.Ad) AdDto {
	return toAdDto(ad)
}

// ToAdsDto converts a list of model.Ad to a list of AdDto.
func ToAdsDto(adList []*model.Ad) []AdDto {
	var ads []AdDto
	for _, ad := range adList {
		ads = append(ads, toAdDto(*ad))
	}
	return ads
}
