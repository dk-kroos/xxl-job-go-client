package main

import (
	"context"
	"fmt"
	xxl "github.com/dk-laosiji/xxl-job-go-client"
	"github.com/dk-laosiji/xxl-job-go-client/example/gin_example/gin_router"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"runtime"
	"time"
)

func main() {
	fmt.Println(runtime.GOOS)
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
		//平滑退出
		endLessStop()
	}()
	engine := gin.Default()
	//设置xxl-job路由
	gin_router.SetXxlRouter(engine)
	//测试正常的gin get接口
	engine.GET("testGin", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "this is a gin test")
		ctx.Abort()
		return
	})

	HttpStarter(engine)

}

func HttpStarter(engine *gin.Engine){
	// windows想支持endless stop的方式参考百度或 https://learnku.com/articles/51696
	if runtime.GOOS == "windows" {
		server, _ := initHTTPServer(context.TODO(), engine)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := endless.ListenAndServe(":"+xxl.SERVER_PORT, engine)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// InitHTTPServer is ...
func initHTTPServer(ctx context.Context, handler http.Handler) (*http.Server, func()) {
	srv := &http.Server{
		Addr:         ":" + xxl.SERVER_PORT,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return srv, func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(30))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}
}

func endLessStop()  {
	// do some thing here when application stop
}