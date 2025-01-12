// Copyright 2024 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://gobase.

package known

const (
	// XRequestIDKey 用来定义 Gin 上下文中的键，代表请求的 UUID.
	X_REQUEST_ID_KEY = "X-Request-ID"

	// XUsernameKey 用来定义 Gin 上下文的键，代表请求的所有者.
	X_USERNAME_KEY = "X-Username"

	// 日期格式化
	TIME_FORMAT_DAY   = "20060102"
	TIME_FORMAT       = "2006-01-02 15:04:05"
	TIME_FORMAT_SHORT = "2006-01-02"

	// 默认数据 ID
	DEFAULT_ID = 1

	// 文件类型
	FILE_TYPE_IMAGE    = 1
	FILE_TYPE_VIDEO    = 2
	FILE_TYPE_AUDIO    = 3
	FILE_TYPE_ARCHIVE  = 4
	FILE_TYPE_DOCUMENT = 5
	FILE_TYPE_FONT     = 6
	FILE_TYPE_APP      = 7
	FILE_TYPE_UNKNOWN  = 8

	// 审核状态
	AUDIT_STATUS_PENDING = 1
	AUDIT_STATUS_PASS    = 2

	// 上传资源大小 (默认 10MB)
	DEFAULT_UPLOAD_SIZE int64 = 100 * 1024 * 1024
)
