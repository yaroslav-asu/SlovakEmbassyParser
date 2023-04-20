package session

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"main/internal/utils/db"
	gorm_models "main/models/gorm"
	"net/http"
	"net/url"
	"time"
)

const proxyWaitTime = 10 * time.Second

func (s *Session) findSuitableProxy(db *gorm.DB, proxy *gorm_models.Proxy) error {
	var err error
	if s.Proxy.Url() != "" {
		err = db.Where("request_time != -1 and is_working = true and ip != ? and port != ?", s.Proxy.Ip, s.Proxy.Port).Order("request_time").First(&proxy).Error
	} else {
		err = db.Where("request_time != -1 and is_working = true").First(&proxy).Error
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) ChangeProxy() {
	var proxy gorm_models.Proxy
	dataBase := db.Connect()
	err := s.findSuitableProxy(dataBase, &proxy)
	unknownErrCounter := 0
	for err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Info("Failed to find working proxy, waiting for new one")
			time.Sleep(proxyWaitTime)
			err = s.findSuitableProxy(dataBase, &proxy)
		} else {
			zap.L().Error("Unknown error occurred while finding working proxies: " + err.Error())
			zap.L().Info("Trying to repeat")
			err = s.findSuitableProxy(dataBase, &proxy)
			unknownErrCounter++
			if unknownErrCounter == 5 {
				zap.L().Warn("Can't find proxies" + err.Error())
				break
			}
		}
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
