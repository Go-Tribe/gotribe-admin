// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

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
