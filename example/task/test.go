package task

import (
	"context"
	xxl "github.com/dk-laosiji/xxl-job-go-client"
	"time"
)

func Test(cxt context.Context, param *xxl.RunReq, logger xxl.Logger) (msg string) {
	logger.Info(param.LogID, "我开始干活了！！！！")
	time.Sleep(time.Second * 1)
	logger.Info(param.LogID, "我执行结束了！！！！")
	return "test done"
}
