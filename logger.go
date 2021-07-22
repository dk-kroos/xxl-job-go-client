package xxljob

// LogFunc 应用日志
type LogFunc func(req LogReq, res *LogRes) []byte

// Logger 系统日志
type Logger interface {
	Info(logId int64, format string, a ...interface{})
	Error(logId int64, format string, a ...interface{})
	Flush()
}
