package task

import (
	"context"
	xxl "github.com/dk-laosiji/xxl-job-go-client"
)

func Panic(cxt context.Context, param *xxl.RunReq,logger xxl.Logger) (msg string) {
	panic("test panic")
	return
}
