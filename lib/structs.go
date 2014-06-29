package lib

import (
	"encoding/json"
)

type Config struct {
	MysqlUsername string
	MysqlHostname string
	MysqlPassword string
	MysqlDatabase string
	HttpPort      int
	HttpAddress   string
}

type JsonStatus struct {
	Code    int
	Message string
}

type JsonData struct {
	Code int
	Data interface{}
}

type Person struct {
	ID        int
	FirstName string
	LastName  string
}

type Task struct {
	ID   int
	Name string
}

func NewStatus(code int, message string) string {
	raw := JsonStatus{code, message}
	jsonVal, jsonErr := json.Marshal(raw)
	if jsonErr != nil {
		return jsonErr.Error()
	}
	return string(jsonVal)
}

func NewData(code int, data interface{}) string {
	raw := JsonData{code, data}
	jsonVal, jsonErr := json.Marshal(raw)
	if jsonErr != nil {
		return jsonErr.Error()
	}
	return string(jsonVal)
}
