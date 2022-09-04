package httpclient

type LogType string

const (
	Info  LogType = "info"
	Error LogType = "error"
)

type LogData struct {
	ContextId string  `json:"context_id"`
	Country   string  `json:"country"`
	Level     LogType `json:"level"`
	Message   string  `json:"msg"`
}
