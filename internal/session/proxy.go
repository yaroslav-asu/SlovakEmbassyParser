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
	if s.Proxy.Url() != "" {
		dataBase.Where("request_time != -1 and is_working = true and ip != ? and port != ?", s.Proxy.Ip, s.Proxy.Port).Order("request_time").First(&proxy)
	} else {
		dataBase.Where("request_time != -1 and is_working = true").First(&proxy)
	}

	s.Proxy = proxy
	urlInstance := url.URL{}
	urlProxy, err := urlInstance.Parse(proxy.Url())
	if err != nil {
		zap.L().Error("Failed to parse proxy url: " + proxy.Url())
	}
	s.Client.Transport = &http.Transport{Proxy: http.ProxyURL(urlProxy)}
}

func (s *Session) DisableCurrentProxy() {
	dataBase := db.Connect()
	dataBase.Model(&s.Proxy).Update("is_working", false)
}
