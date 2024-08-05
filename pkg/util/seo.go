// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package util

import (
	"github.com/dengmengmian/ghelper/ghttp"
	"log"
)

func PushBaidu(site, token string, urls string) bool {
	api := "http://data.zz.baidu.com/urls?site=" + site + "&token=" + token
	req := &ghttp.Request{
		Method:      "POST",
		Url:         api,
		Body:        urls,
		ContentType: "Content-Type: text/plain",
	}
	response, err := req.Do()
	if err != nil {
		log.Printf("推送失败：", response)
	}
	defer response.Body.Close()
	result, err := response.Body.FromToString()
	log.Printf("推送记录：", result)
	return true
}
