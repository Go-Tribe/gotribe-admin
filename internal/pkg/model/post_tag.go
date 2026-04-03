// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import "time"

// PostTag 内容-标签关联表
type PostTag struct {
	PostID    uint      `gorm:"primaryKey;comment:内容ID" json:"postID"`
	TagID     uint      `gorm:"primaryKey;index:idx_post_tag_tag_id;comment:标签ID" json:"tagID"`
	CreatedAt time.Time `json:"createdAt"`
}
