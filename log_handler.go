package xxljob

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/**
用来日志查询，显示到xxl-job-admin后台
*/

type LogHandler func(req *LogReq) *LogRes

//默认返回
func defaultLogHandler(req *LogReq) *LogRes {
	return &LogRes{Code: 200, Msg: "", Content: LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   2,
		LogContent:  "这是日志默认返回，说明没有设置LogHandler",
		IsEnd:       true,
	}}
}

//请求错误
func reqErrLogHandler(w http.ResponseWriter, req *LogReq, err error) {
	res := &LogRes{Code: 500, Msg: err.Error(), Content: LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   0,
		LogContent:  err.Error(),
		IsEnd:       true,
	}}
	str, _ := json.Marshal(res)
	_, _ = w.Write(str)
}

//es查询日志
func EsLogHandler(req *LogReq) *LogRes {
	var res *LogRes
	if esClient == nil {
		res = &LogRes{Code: 500, Msg: "elasticsearch client is nil", Content: LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   2,
			LogContent:  "elasticsearch client is nil",
			IsEnd:       true,
		}}
		return res
	}
	//从ES查询日志
	data, err := esClient.ReadLog(req.LogID)
	if err != nil {
		res = &LogRes{Code: 500, Msg: err.Error(), Content: LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   2,
			LogContent:  err.Error(),
			IsEnd:       true,
		}}
	}
	var logContent strings.Builder
	for _, logInfo := range data {
		fmt.Fprintln(&logContent, logInfo)
	}
	res = &LogRes{Code: 200, Msg: "", Content: LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   2,
		LogContent:  logContent.String(),
		IsEnd:       true,
	}}
	return res
}
