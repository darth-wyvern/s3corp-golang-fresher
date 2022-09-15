package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type SuccessResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

// WriteJSONResponse writes the response to client
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	var v, ok = data.(string) // if datatype of data is string
	if ok {
		w.WriteHeader(statusCode)
		w.Write([]byte(v))
		return
	}
	result, err := json.Marshal(data)
	if err != nil {
		log.Printf("marsharing json error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	if data != nil {
		w.Write(result)
	}
}

// WriteErrorResponse writes the error response to client
//
// Deprecated: use WriteJSONResponse instead
func WriteErrorResponse(w http.ResponseWriter, errResp error) {
	result, err := json.Marshal(errResp)
	if err != nil {
		log.Printf("marsharing json error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var v, ok = errResp.(ErrorResponse)
	// Display status code if error is of type ErrorResponse else 500
	if ok {
		w.WriteHeader(v.Status)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(result)
}
