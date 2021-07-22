package xxljob

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"sync"
	"time"
)

//实现了Logger接口
type LogEs struct {
	Contents []string    `json:"contents"` // 日志内容，主动记录
	LogId    int64       `json:"logId"`    // 日志ID
	Locker   *sync.Mutex `json:"-"`        // 加锁，防止日志覆盖
}

func NewLogEs() *LogEs {
	return &LogEs{
		Contents: make([]string, 0),
		Locker:   &sync.Mutex{},
	}
}

//追加ERROR类型的Task Log
func (l *LogEs) Error(logId int64, format string, log ...interface{}) {
	l.Locker.Lock()
	defer l.Locker.Unlock()
	l.LogId = logId
	now := time.Now().Format("2006-01-02 15:04:05")
	logInfo := fmt.Sprintf("ERROR日志 - "+now+" - "+format, log...)
	l.Contents = append(l.Contents, logInfo)
}

//追加INFO 类型的Task Log
func (l *LogEs) Info(logId int64, format string, log ...interface{}) {
	l.Locker.Lock()
	defer l.Locker.Unlock()
	l.LogId = logId
	now := time.Now().Format("2006-01-02 15:04:05")
	logInfo := fmt.Sprintf("INFO日志 - "+now+" - "+format, log...)
	l.Contents = append(l.Contents, logInfo)
}

//flush Log
func (l *LogEs) Flush() {
	l.Locker.Lock()
	defer l.Locker.Unlock()
	//发送日志到ES端
	esClient.AddLog(l)
	//重置日志
	l.Contents = make([]string, 0)
}

//初始化ES
func (l *LogEs) InitEs(server string, appId string, env string) *LogEs {
	esClient = &EsClient{
		Server: server,
		AppId:  appId,
		Env:    env,
	}
	esClient.esCli, _ = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(server))
	return l
}
