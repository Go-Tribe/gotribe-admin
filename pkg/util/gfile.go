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
		fileType = known.IMAGE
	case filetype.IsAudio(head):
		fileType = known.AUDIO
	case filetype.IsApplication(head):
		fileType = known.APP
	case filetype.IsVideo(head):
		fileType = known.VIDEO
	case filetype.IsArchive(head):
		fileType = known.ARCHIVE
	case filetype.IsDocument(head):
		fileType = known.DOCUMENT
	case filetype.IsFont(head):
		fileType = known.FONT
	default:
		fileType = known.UNKNOWN
	}
	return fileType
}
