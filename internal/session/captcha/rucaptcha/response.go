package rucaptcha

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Response struct {
	Status    int    `json:"status"`
	Request   string `json:"request"`
	ErrorText string `json:"error_text"`
}

func (r *Response) Format() string {
	return fmt.Sprintf("Response{Status: %s, Request: %s, ErrText: %s}", r.Status, r.Request, r.ErrorText)
}

func ParseRucaptchaResponse(res *http.Response) Response {
	zap.L().Info("Parsing response")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("Failed to parse response")
	}
	zap.L().Info("Coping response to structure")
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		zap.L().Error("Failed to unmarshal text to json")
	}
	return response
}
