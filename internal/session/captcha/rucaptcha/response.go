package rucaptcha

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Response struct {
	Status    int    `json:"status"`
	Request   string `json:"request"`
	ErrorText string `json:"error_text"`
}

func ParseRucaptchaResponse(res *http.Response) Response {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("Failed to parse response")
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		zap.L().Error("Failed to unmarshal text to json")
	}
	return response
}
