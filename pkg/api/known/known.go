// Copyright 2024 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://gobase.

package known

const (
	// XRequestIDKey 用来定义 Gin 上下文中的键，代表请求的 uuid.
	XRequestIDKey = "X-Request-ID"

	// XUsernameKey 用来定义 Gin 上下文的键，代表请求的所有者.
	XUsernameKey = "X-Username"

	// 日期格式化
	TimeFormatDay = "20060102"
	TimeFormat    = "2006-01-02 15:04:05"

	// 上传资源大小
	DefUploadSize int64 = 10 * 1024 * 1024

	StatusOK = 1
	StatusNO = 2
	// 默认数据 ID
	DefulatID = 1
	// 文件类型
	IMAGE    = 1
	VIDEO    = 2
	AUDIO    = 3
	ARCHIVE  = 4
	DOCUMENT = 5
	FONT     = 6
	APP      = 7
	UNKNOWN  = 8
)
