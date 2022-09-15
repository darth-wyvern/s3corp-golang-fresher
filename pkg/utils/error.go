package utils

import (
	"fmt"
)

type ErrorResponse struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Desc   string `json:"desc"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("{\"status\":%d,\"code\":\"%s\",\"desc\":\"%s\"}", e.Status, e.Code, e.Desc)
}
