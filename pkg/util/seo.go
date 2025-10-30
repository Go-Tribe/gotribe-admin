// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package util

import (
	"fmt"
	"log"

	"github.com/dengmengmian/ghelper/ghttp"
)

// SEO SEO工具类，用于处理SEO相关操作
type SEO struct{}

// PushBaidu 推送URL到百度
func (s *SEO) PushBaidu(site, token string, urls string) (bool, error) {
	api := "http://data.zz.baidu.com/urls?site=" + site + "&token=" + token
	req := &ghttp.Request{
		Method:      "POST",
		Url:         api,
		Body:        urls,
		ContentType: "Content-Type: text/plain",
	}
	response, err := req.Do()
	if err != nil {
		return false, fmt.Errorf("推送失败: %v", err)
	}
	defer response.Body.Close()

	result, err := response.Body.FromToString()
	if err != nil {
		return false, fmt.Errorf("解析响应失败: %v", err)
	}

	log.Printf("推送记录: %s", result)
	return true, nil
}

// 全局实例
var SEOUtil = &SEO{}
