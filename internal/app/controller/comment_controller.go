// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/pkg/api/dto"
	"gotribe-admin/pkg/api/known"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
)

type ICommentController interface {
	GetComments(c *gin.Context)       // 获取评论列表
	UpdateCommentByID(c *gin.Context) // 更新评论
}

type CommentController struct {
	CommentRepository repository.ICommentRepository
}

// 构造函数
func NewCommentController() ICommentController {
	commentRepository := repository.NewCommentRepository()
	commentController := CommentController{CommentRepository: commentRepository}
	return commentController
}

// GetComments 获取评论列表
// @Summary      获取评论列表
// @Description  获取所有评论的列表，支持分页和筛选
// @Tags         评论管理
// @Accept       json
// @Produce      json
// @Param        request query vo.CommentListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /comment [get]
// @Security     BearerAuth
func (pc CommentController) GetComments(c *gin.Context) {
	var req vo.CommentListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.ValidationFail(c, errStr)
		return
	}

	// 获取
	comment, total, err := pc.CommentRepository.GetComments(&req)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgListFail)
		return
	}
	response.Success(c, gin.H{"comments": dto.ToCommentsDto(comment), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// UpdateCommentByID 更新评论
// @Summary      更新评论状态
// @Description  根据评论ID更新评论的审核状态
// @Tags         评论管理
// @Accept       json
// @Produce      json
// @Param        commentID path string true "评论ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /comment/{commentID} [patch]
// @Security     BearerAuth
func (pc CommentController) UpdateCommentByID(c *gin.Context) {
	// 根据path中的CommentID获取评论信息
	oldComment, err := pc.CommentRepository.GetCommentByComentID(c.Param("commentID"))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}
	var reqStatus uint
	if oldComment.Status == known.AUDIT_STATUS_PENDING {
		reqStatus = known.AUDIT_STATUS_PASS
	} else {
		reqStatus = known.AUDIT_STATUS_PENDING
	}
	oldComment.Status = reqStatus
	// 更新评论
	err = pc.CommentRepository.UpdateComment(&oldComment)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgUpdateFail)
		return
	}
	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}
