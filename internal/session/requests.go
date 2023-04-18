package session

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	"net"
	"net/http"
	"net/url"
)

func (s *Session) GetParsedSoup(url string) soup.Root {
	doc, err := soup.GetWithClient(url, s.Client)
	if err != nil {
		s.handleRequestError(url, err)
		return s.GetParsedSoup(url)
	}
	parsedSoup := soup.HTMLParse(doc)
	if !sessionWorking(parsedSoup) {
		zap.L().Warn("Session cookies aren't valid, starting to log in")
		s.LogIn()
		return s.GetParsedSoup(url)
	}
	return parsedSoup
}

func (s *Session) Get(url string) *http.Response {
	res, err := s.Client.Get(url)
	parsedSoup := funcs.ResponseToSoup(res)
	if err != nil {
		s.handleRequestError(url, err)
		return s.Get(url)
	}
	if !sessionWorking(parsedSoup) {
		zap.L().Warn("Session cookies aren't valid, starting to log in")
		s.LogIn()
		return s.Get(url)
	}
	return res
}

func (s *Session) PostForm(url string, data url.Values) *http.Response {
	res, err := s.Client.PostForm(url, data)
	if err != nil {
		s.handleRequestError(url, err)
		return s.PostForm(url, data)
	}
	return res
}

func (s *Session) handleRequestError(url string, err error) {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		zap.L().Warn("Proxy timeout: " + s.Proxy.Url())
	} else if err != nil {
		zap.L().Warn("Cant access to:" + url + " with proxy: " + s.Proxy.Url())
	}
	zap.L().Info("Trying to change proxy")
	funcs.SleepTime(5, 10)
	s.DisableCurrentProxy()
	s.ChangeProxy()
}
