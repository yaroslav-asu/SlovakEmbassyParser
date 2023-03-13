package session

import (
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

func (s Session) getSoup(link string) (string, error) {
	return soup.GetWithClient(link, &s.client)
}

func (s Session) GetParsedSoup(link string) soup.Root {
	doc, err := s.getSoup(link)
	if err != nil {
		zap.L().Error("Cant get: " + link)
		return soup.Root{}
	}
	return soup.HTMLParse(doc)
}
func (s Session) Get(url string) *http.Response {
	res, err := s.client.Get(url)
	if err != nil {
		zap.L().Error("Cant get: " + url)
	}
	return res
}

func (s Session) PostForm(url string, data url.Values) *http.Response {
	res, err := s.client.PostForm(url, data)
	if err != nil {
		zap.L().Error("Cant post form to: " + url)
	}
	return res
}
