package jobs

import "github.com/robfig/cron/v3"

var Cron *cron.Cron

func InitCron() {
	secondParser := cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour |
			cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	job := cron.New(cron.WithParser(secondParser), cron.WithChain())
	Cron = job
	JobRun(job)
}

func JobRun(job *cron.Cron) {
	// 示例定时任务
	//job.AddFunc("@every 5s", func() {
	//	exampleJob()
	//})
	job.AddFunc("@every 1m", func() {
		sitemapJob()
	})
}
