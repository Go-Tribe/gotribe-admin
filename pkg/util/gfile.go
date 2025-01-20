// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package util

import (
	"github.com/h2non/filetype"
	"gotribe-admin/pkg/api/known"
	"mime/multipart"
)

func GetFileType(header *multipart.FileHeader) int {
	file, _ := header.Open()
	head := make([]byte, 261)
	file.Read(head)
	// 检查文件类型
	var fileType int
	switch {
	case filetype.IsImage(head):
		fileType = known.FILE_TYPE_IMAGE
	case filetype.IsAudio(head):
		fileType = known.FILE_TYPE_AUDIO
	case filetype.IsApplication(head):
		fileType = known.FILE_TYPE_APP
	case filetype.IsVideo(head):
		fileType = known.FILE_TYPE_VIDEO
	case filetype.IsArchive(head):
		fileType = known.FILE_TYPE_ARCHIVE
	case filetype.IsDocument(head):
		fileType = known.FILE_TYPE_DOCUMENT
	case filetype.IsFont(head):
		fileType = known.FILE_TYPE_FONT
	default:
		fileType = known.FILE_TYPE_UNKNOWN
	}
	return fileType
}
