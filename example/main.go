package main

import (
	xxl "github.com/dk-laosiji/xxl-job-go-client"
	"github.com/dk-laosiji/xxl-job-go-client/example/task"
	"log"
)

func main() {
	exec := xxl.NewExecutor(
		xxl.ServerAddr("http://192.168.1.50:8080/xxl-job-admin"),
		xxl.AccessToken(""), //请求令牌(默认为空)
		//xxl.ExecutorIp("xxl-job-executor-test"),    //可自动获取
		xxl.ExecutorPort("9999"),  //默认9999（非必填）
		xxl.RegistryKey("golang-jobs"),  //执行器名称
		xxl.SetLogger(xxl.NewLogEs().InitEs("http://192.168.1.50:9200", "100058", "dev")), //自定义日志
	)
	exec.Init()
	//设置日志查看handler
	exec.LogHandler(xxl.EsLogHandler)
	//注册任务handler
	exec.RegTask("task.test", task.Test)
	exec.RegTask("task.test2", task.Test2)
	exec.RegTask("task.panic", task.Panic)
	log.Fatal(exec.Run())
}
