// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package util

import (
	"encoding/json"
	"fmt"
)

// JSON 工具类，用于处理 JSON 序列化和反序列化
type JSON struct{}

// Struct2Json 结构体转为json
func (j *JSON) Struct2Json(obj interface{}) (string, error) {
	str, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("[Struct2Json]转换异常: %v", err)
	}
	return string(str), nil
}

// Json2Struct json转为结构体
func (j *JSON) Json2Struct(str string, obj interface{}) error {
	// 将json转为结构体
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		return fmt.Errorf("[Json2Struct]转换异常: %v", err)
	}
	return nil
}

// JsonI2Struct json interface转为结构体
func (j *JSON) JsonI2Struct(str interface{}, obj interface{}) error {
	JsonStr, ok := str.(string)
	if !ok {
		return fmt.Errorf("[JsonI2Struct]输入参数不是字符串类型")
	}
	return j.Json2Struct(JsonStr, obj)
}

// 全局实例
var JSONUtil = &JSON{}
