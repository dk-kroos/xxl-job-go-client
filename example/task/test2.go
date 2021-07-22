package task

import (
	"context"
	"fmt"
	xxl "github.com/dk-laosiji/xxl-job-go-client"
	"time"
)

func Test2(cxt context.Context, param *xxl.RunReq, logger xxl.Logger) (msg string) {
	num := 1
	for {

		select {
		case <-cxt.Done():
			logger.Error(param.LogID, fmt.Sprintf("task"+param.ExecutorHandler+"被手动终止"))
			return
		default:
			num++
			time.Sleep(10 * time.Second)
			logger.Info(param.LogID, fmt.Sprintf("test one task"+param.ExecutorHandler+" param："+param.ExecutorParams+"执行行", num))

			if num > 10 {
				logger.Info(param.LogID, fmt.Sprintf("test one task"+param.ExecutorHandler+" param："+param.ExecutorParams+"执行完毕！"))
				return
			}
		}
	}

}
