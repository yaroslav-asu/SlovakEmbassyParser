package funcs

import (
	"bytes"
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"io"
	"main/internal/utils/vars"
	"net/http"
	"strings"
)

func Linkify(linkParts ...string) string {
	link := vars.SiteUrl
	for _, linkPart := range linkParts {
		link += strings.Replace(linkPart, "/", "", -1)
	}
	return link
}

func ResponseToSoup(res *http.Response) soup.Root {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("Failed to read response")
		return soup.Root{}
	}
	html := string(body)
	res.Body = io.NopCloser(bytes.NewBuffer(body))
	return soup.HTMLParse(html)
}
