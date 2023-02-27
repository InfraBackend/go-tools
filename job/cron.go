package job

import "github.com/robfig/cron/v3"

// 使用定时任务库
var Cron = cron.New()

func init() {
	Cron.Run()
}
