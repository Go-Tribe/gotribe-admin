// 修改后的完整代码示例

package dto

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/known"
)

type FeedbackDto struct {
	ID        int        `json:"id"`
	ProjectID string     `json:"projectID"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	UserID    string     `json:"userID"`
	User      UserDto    `json:"user"`
	Phone     string     `json:"phone"`
	Project   ProjectDto `json:"project"`
	CreatedAt string     `json:"createdAt"`
}

func toFeedbackDto(feedBack model.Feedback) FeedbackDto {
	dto := FeedbackDto{
		ID:        int(feedBack.ID),
		ProjectID: feedBack.ProjectID,
		Content:   feedBack.Content,
		Title:     feedBack.Title,
		UserID:    feedBack.UserID,
		Phone:     feedBack.Phone,
		CreatedAt: feedBack.CreatedAt.Format(known.TIME_FORMAT),
	}
	if feedBack.User != nil {
		dto.User = ToUserInfoDto(feedBack.User)
	}
	if feedBack.Project != nil {
		dto.Project = ToProjectInfoDto(feedBack.Project)
	}
	return dto
}

// 将单个 Feedback 转换为 FeedbackDto
func ToFeedbackInfoDto(feedBack model.Feedback) FeedbackDto {
	return toFeedbackDto(feedBack)
}

// 将多个 Feedback 转换为 []FeedbackDto
func ToFeedbacksDto(feedBackList []*model.Feedback) []FeedbackDto {
	var feedBacks = make([]FeedbackDto, len(feedBackList))
	for i, feedBack := range feedBackList {
		feedBacks[i] = toFeedbackDto(*feedBack)
	}
	return feedBacks
}
