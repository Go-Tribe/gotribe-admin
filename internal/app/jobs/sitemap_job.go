// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package jobs

import (
	"context"

	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"

	"github.com/douyacun/gositemap"
)

// SitemapJob 站点地图生成任务
type SitemapJob struct {
	*BaseJob
}

// NewSitemapJob 创建站点地图任务
func NewSitemapJob(config JobConfig) *SitemapJob {
	job := &SitemapJob{}
	job.BaseJob = NewBaseJob(config, job.execute)
	return job
}

// execute 执行站点地图生成
func (j *SitemapJob) execute(ctx context.Context) error {
	common.Log.Info("Starting sitemap generation job")

	// 查出 project 信息
	projects, err := repository.NewProjectRepository().GetProjectsBySitemap()
	if err != nil {
		return err
	}

	var posts []*model.Post
	st := gositemap.NewSiteMap()
	st.SetPretty(true)
	st.SetPublicPath("public")

	for _, project := range projects {
		// 检查上下文是否被取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 使用固定的文件名，确保每次都是覆盖
		st.SetFilename(project.ProjectID + ".xml")

		if err := common.DB.Model(&model.Post{}).Where("status = ? and type != ? and project_id = ?", 2, 2, project.ProjectID).Find(&posts).Error; err != nil {
			common.Log.Errorf("Failed to query posts for project %s: %v", project.ProjectID, err)
			continue
		}

		for _, post := range posts {
			url := gositemap.NewUrl()
			url.SetLoc(project.PostURL + post.PostID)
			url.SetLastmod(post.UpdatedAt)
			url.SetChangefreq(gositemap.Daily)
			url.SetPriority(1)
			st.AppendUrl(url)
		}
	}

	if _, err := st.Storage(); err != nil {
		return err
	}

	common.Log.Infof("Sitemap generation completed for %d projects", len(projects))
	return nil
}
