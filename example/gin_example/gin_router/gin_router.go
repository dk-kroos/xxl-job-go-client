package gin_router

import (
	xxl "github.com/dk-laosiji/xxl-job-go-client"
	"github.com/dk-laosiji/xxl-job-go-client/example/task"
	"github.com/gin-gonic/gin"
)

// 把gin引擎传过来
func SetXxlRouter(engine *gin.Engine) {
	//初始化执行器
	exec := xxl.NewExecutor(
		xxl.ServerAddr("http://192.168.1.50:8080/xxl-job-admin"),                             //admin 地址
		xxl.AccessToken(""),                                                                  //请求令牌(默认为空)
		xxl.ExecutorPort(xxl.SERVER_PORT),                                                    //默认9999（此处要与gin服务启动port必需一至）
		xxl.RegistryKey("golang-jobs"),                                                       //执行器名称
		xxl.SetLogger(xxl.NewLogEs().InitEs([]string{"192.168.1.50:9200"}, "100058", "dev")), //自定义日志组件
	)
	exec.Init()
	//设置日志handler为ES
	exec.LogHandler(xxl.EsLogHandler)
	//注册admin需要访问的路由方法
	engine.POST("run", gin.WrapF(exec.RunTask))
	engine.POST("kill", gin.WrapF(exec.KillTask))
	engine.POST("log", gin.WrapF(exec.TaskLog))
	engine.POST("beat", gin.WrapF(exec.Beat))
	engine.POST("idleBeat", gin.WrapF(exec.IdleBeat))

	//注册任务handler
	exec.RegTask("task.test", task.Test)
	exec.RegTask("task.test2", task.Test2)
	exec.RegTask("task.panic", task.Panic)
}
