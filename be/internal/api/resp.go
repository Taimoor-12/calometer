package api

import (
	"calometer/internal/logger"
)

var log = logger.GetLogger()

type Response struct {
	Code map[int]string `json:"code"`
	Data interface{}    `json:"data,omitempty"`
}
