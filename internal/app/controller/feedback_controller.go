package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thoas/go-funk"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/dto"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
)

type IFeedbackController interface {
	GetFeedbacks(c *gin.Context) // 获取列表
}

type FeedbackController struct {
	FeedbackRepository repository.IFeedbackRepository
}

// 构造函数
func NewFeedbackController() IFeedbackController {
	feedbackRepository := repository.NewFeedbackRepository()
	feedbackController := FeedbackController{FeedbackRepository: feedbackRepository}
	return feedbackController
}

// GetFeedbacks 获取反馈列表
// @Summary      获取反馈列表
// @Description  获取所有反馈的列表，支持分页和筛选
// @Tags         反馈管理
// @Accept       json
// @Produce      json
// @Param        request query vo.FeedbackListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /feedback [get]
// @Security     BearerAuth
func (tc FeedbackController) GetFeedbacks(c *gin.Context) {
	var req vo.FeedbackListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.GetTransFromCtx(c))
		response.Fail(c, nil, errStr)
		return
	}

	// 获取
	feedbacks, total, err := tc.FeedbackRepository.GetFeedbacks(&req)
	if err != nil {
		response.Fail(c, nil, common.Msg(c, common.MsgListFail)+": "+err.Error())
		return
	}

	// 获取所有的用户信息和项目信息，并附加到反馈中
	feedbacks, err = getFeedbackOther(feedbacks)
	if err != nil {
		response.Fail(c, nil, "获取用户或项目信息失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"feedbacks": dto.ToFeedbacksDto(feedbacks), "total": total}, common.Msg(c, common.MsgListSuccess))
}

func getFeedbackOther(feedbacks []*model.Feedback) ([]*model.Feedback, error) {
	// 遍历feedback获取所有用户ID,查询用户信息并加进去
	userIds := make([]string, 0)
	for _, feedback := range feedbacks {
		userIds = append(userIds, feedback.UserID)
	}
	userIds = funk.UniqString(userIds)
	var users []model.User
	if len(userIds) > 0 {
		if err := common.DB.Where("user_id in (?)", userIds).Find(&users).Error; err != nil {
			return feedbacks, err
		}
	}

	// 创建用户映射以提高查找效率
	userMap := make(map[string]*model.User)
	for _, user := range users {
		userMap[user.UserID] = &user
	}

	// 将用户信息附加到反馈中
	for _, feedback := range feedbacks {
		if user, ok := userMap[feedback.UserID]; ok {
			feedback.User = user
		}
	}

	// 追加项目信息
	projectIds := funk.UniqString(funk.Map(feedbacks, func(feedback *model.Feedback) string {
		return feedback.ProjectID
	}).([]string))
	var projects []model.Project
	if len(projectIds) > 0 {
		if err := common.DB.Where("project_id in (?)", projectIds).Find(&projects).Error; err != nil {
			return feedbacks, err
		}
	}

	// 创建项目映射以提高查找效率
	projectMap := make(map[string]*model.Project)
	for _, project := range projects {
		projectMap[project.ProjectID] = &project
	}

	// 将项目信息附加到反馈中
	for _, feedback := range feedbacks {
		if project, ok := projectMap[feedback.ProjectID]; ok {
			feedback.Project = project
		}
	}

	return feedbacks, nil
}
