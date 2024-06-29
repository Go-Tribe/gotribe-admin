package jobs

import (
	"github.com/dengmengmian/ghelper/gconvert"
	"github.com/douyacun/gositemap"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
)

func sitemapJob() {
	// 查出 porject 信息
	projects, err := repository.NewProjectRepository().GetProjectsBySitemap()
	if err != nil {
		common.Log.Error("sitemapJob:", err.Error())
		return
	}
	var posts []*model.Post
	st := gositemap.NewSiteMap()
	st.SetPretty(true)
	st.SetPublicPath("public")
	for idx, project := range projects {
		st.SetFilename(project.ProjectID + gconvert.String(idx) + ".xml")

		if err := common.DB.Model(&model.Post{}).Where("status = ? and type != ? and project_id = ?", 2, 2, project.ProjectID).Find(&posts).Error; err != nil {
			common.Log.Error("post:", err.Error())
			return
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
	st.Storage()
	return
}
