package response

import "net/http"

// Envelope is the unified API response shape.
type Envelope struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// Success builds a successful response payload.
func Success(data interface{}, msg string) Envelope {
	return Envelope{Code: http.StatusOK, Data: data, Msg: msg}
}

// Error builds an error response payload with a custom status code.
func Error(status int, msg string, data interface{}) Envelope {
	return Envelope{Code: status, Data: data, Msg: msg}
}
