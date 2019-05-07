package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time
}

func main() {
	var (
		expr         *cronexpr.Expression
		now          time.Time
		cronJob      *CronJob
		scheduleTale map[string]*CronJob
	)

	scheduleTale = make(map[string]*CronJob)

	//当前时间
	now = time.Now()

	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	scheduleTale["job1"] = cronJob

	expr = cronexpr.MustParse("*/4 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	scheduleTale["job2"] = cronJob

	//启动一个调度协程
	go func() {
		var (
			now     time.Time
			jobName string
			cronJob *CronJob
		)

		//定时检查一下任务调度表
		for {
			now = time.Now()

			for jobName, cronJob = range scheduleTale {

				//判断是否过期
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					//启动一个协程，执行这个任务
					go func(jobName string) {
						fmt.Println("执行：", jobName)
					}(jobName)

					//计算下一次调度时间
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println(jobName, "下次执行时间：", cronJob.nextTime)
				}

			}

			//睡眠100毫秒
			select {
			case <-time.NewTimer(100 * time.Millisecond).C:

			}
		}

	}()

	time.Sleep(100 * time.Second)
}
