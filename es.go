package xxljob

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
	"time"
)

var esClient *EsClient

const INDEX_KEY_PREFIX = "xxl-job-log-"

type EsClient struct {
	Servers []string //服务地址
	AppId   string   `json:"appId"` // 应用ID
	Env     string   `json:"env"`   // 运行环境
	esCli   *elastic.Client
}

//往es写日志
func (e *EsClient) AddLog(logEs *LogEs) error {
	if len(logEs.Contents) < 1 {
		return errors.New("log is empty")
	}
	//如果ES没有初始化成功，再次尝试初始化ES
	if e.esCli == nil {
		errConn := e.Connect()
		if errConn != nil {
			return errConn
		}
	}

	//组装index
	key := INDEX_KEY_PREFIX + time.Now().Format("2006-01-02")
	//组装新的map
	logData := make(map[string]interface{}, 4)
	logData["appId"] = e.AppId
	logData["env"] = e.Env
	logData["logId"] = logEs.LogId
	logData["content"] = logEs.Contents
	logData["@timestamp"] = time.Now()
	//map转json字符串
	jsonLog, _ := json.Marshal(logData)

	//写入es
	_, err := e.esCli.Index().
		Index(key).
		BodyJson(string(jsonLog)).
		Do(context.Background())

	return err
}

//从es读日志
func (e *EsClient) ReadLog(logId int64) ([]interface{}, error) {
	if logId < 1 {
		return nil, errors.New("logId is empty")
	}

	//如果ES没有初始化成功，再次尝试初始化ES
	if e.esCli == nil {
		errConn := e.Connect()
		if errConn != nil {
			return nil, errConn
		}
	}

	//组装index
	key := INDEX_KEY_PREFIX + time.Now().Format("2006-01-02")
	//读取es
	esq := elastic.NewTermQuery("logId", logId)
	searchRes, err := e.esCli.Search().
		Index(key).
		Query(esq).
		Pretty(true).
		Do(context.Background())

	if err != nil {
		return nil, err
	}
	if searchRes == nil || searchRes.Hits.TotalHits.Value < 1 {
		return nil, errors.New("log is empty")
	}

	//遍历，组装数据
	data := make([]interface{}, 0)
	for _, hit := range searchRes.Hits.Hits {
		logInfo := make(map[string]interface{})
		err := json.Unmarshal(hit.Source, &logInfo)
		if err != nil || len(logInfo) < 1 {
			continue
		}
		//合并日志内容
		if content, ok := logInfo["content"]; ok {
			if logContent, ok := content.([]interface{}); ok {
				data = append(data, logContent...)
			}

		}
	}

	if len(data) < 1 {
		return nil, errors.New("log is empty")
	}
	return data, nil
}

//连接ES
func (e *EsClient) Connect() error {
	if len(e.Servers) < 1 {
		return errors.New("es server is empty")
	}
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(e.Servers...))
	if client == nil || err != nil {
		return err
	}
	e.esCli = client
	return nil
}
