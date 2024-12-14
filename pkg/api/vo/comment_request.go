// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 获取评论
type CommentListRequest struct {
	ProjectID  string `form:"projectID" json:"projectID"`
	ObjectID   string `form:"objectID" json:"objectID"`
	ObjectType uint   `form:"objectType" json:"objectType"`
	Status     uint   `form:"status" json:"status"`
	Username   string `form:"username" json:"username"`
	PageNum    uint   `json:"pageNum" form:"pageNum"`
	PageSize   uint   `json:"pageSize" form:"pageSize"`
}

// 更新评论
type UpdateCommentRequest struct {
	Status uint `form:"status" json:"status" validate:"oneof=1 2"`
}
