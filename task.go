package xxljob

import (
	"context"
	"fmt"
)

// TaskFunc 任务执行函数
type TaskFunc func(cxt context.Context, param *RunReq, logger Logger) string

// Task 任务
type Task struct {
	Id        int64
	Name      string
	Ext       context.Context
	Param     *RunReq
	fn        TaskFunc
	Cancel    context.CancelFunc
	StartTime int64
	EndTime   int64
	//日志
	log   Logger
	LogId int64
}

// Run 运行任务
func (t *Task) Run(callback func(code int64, msg string)) {
	defer t.log.Flush()
	defer func(cancel func()) {
		if err := recover(); err != nil {
			t.log.Error(t.LogId, t.Info()+" panic: ", fmt.Sprintf("%s", err))
			callback(500, "task panic:"+fmt.Sprintf("%s", err))
			cancel()
		}
	}(t.Cancel)
	msg := t.fn(t.Ext, t.Param, t.log)
	t.log.Info(t.LogId, "任务信息:", "任务ID["+Int64ToStr(t.Id)+"]；任务名称["+t.Name+"]；参数："+t.Param.ExecutorParams+" 执行完毕")
	callback(200, msg)
	return
}

// Info 任务信息
func (t *Task) Info() string {
	return "任务ID[" + Int64ToStr(t.Id) + "]；任务名称[" + t.Name + "]；参数：" + t.Param.ExecutorParams
}
