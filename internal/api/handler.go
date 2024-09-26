package api

import (
	"calometer/internal/logger"
)

var log = logger.GetLogger()

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}
