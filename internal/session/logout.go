package session

import (
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	"net/http"
)

func Logout(client *http.Client) {
	res, err := client.Get(funcs.Linkefy("j_spring_security_logout"))
	if err != nil {
		zap.L().Warn("Cant get logout page")
	}
	switch res.StatusCode {
	case 200:
		zap.L().Info("Successfully logged out")
	default:
		zap.L().Warn("On logout got error code: " + res.Status)
	}
}
