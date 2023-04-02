package session

import (
	"go.uber.org/zap"
	"main/internal/utils/db"
	gorm_models "main/models/gorm"
	"net/http"
	"net/url"
)

func (s *Session) ChangeProxy() {
	dataBase := db.Connect()
	var proxy gorm_models.Proxy
	if s.Proxy.Url != "" {
		dataBase.Order("request_time").Not("url = ? and working = false", s.Proxy.Url).First(&proxy)
	} else {
		dataBase.Order("request_time").Not("working = false").First(&proxy)
	}

	s.Proxy = proxy
	urlInstance := url.URL{}
	urlProxy, err := urlInstance.Parse("http://" + proxy.Url)
	if err != nil {
		zap.L().Error("Failed to parse proxy url: http://" + proxy.Url)
	}
	s.Client.Transport = &http.Transport{Proxy: http.ProxyURL(urlProxy)}
}

func (s *Session) DisableCurrentProxy() {
	dataBase := db.Connect()
	dataBase.Model(&s.Proxy).Update("working", false)
}
