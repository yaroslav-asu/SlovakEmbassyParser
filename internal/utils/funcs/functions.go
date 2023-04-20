package funcs

import (
	"bytes"
	"fmt"
	"github.com/anaskhan96/soup"
	"go.uber.org/zap"
	"io"
	"main/internal/logger"
	"main/internal/session/captcha"
	"main/internal/utils/vars"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func StripString(s string) string {
	return stripStringRunes(s, ' ', '\n', '\t')
}

func stripStringRunes(s string, runes ...rune) string {
	for i, letter := range s {
		if !isRuneInList(letter, runes) {
			s = s[i:]
			break
		}
	}
	for i := len(s) - 1; i >= 0; i-- {
		if !isRuneInList(rune(s[i]), runes) {
			s = s[:i+1]
			break
		}
	}
	return s
}
func isRuneInList(checkingRune rune, runes []rune) bool {
	for _, r := range runes {
		if checkingRune == r {
			return true
		}
	}
	return false
}

func Init() {
	vars.InitEnv()
	logger.InitLogger()
	captcha.InitCaptcha()
}

func Linkify(linkParts ...string) string {
	link := vars.SiteUrl
	for _, linkPart := range linkParts {
		link += strings.Replace(linkPart, "/", "", -1)
	}
	return link
}

func StringsToIntArray(stringArr []string) []int {
	intArr := make([]int, len(stringArr))
	for i, s := range stringArr {
		intEl, err := strconv.Atoi(s)
		if err != nil {
			zap.L().Error("Can't convert string to int: " + s)
			return []int{}
		}
		intArr[i] = intEl
	}
	return intArr
}
func Sleep() {
	SleepTime(2, 5)
}

func SleepTime(from, to float64) {
	sleepingTime := rand.Float64()*(to-from) + from
	zap.L().Info("Started random sleeping with time: " + fmt.Sprintf("%.2f", sleepingTime))
	time.Sleep(time.Duration(sleepingTime) * time.Second)
	zap.L().Info("Finished random sleeping with time: " + fmt.Sprintf("%.2f", sleepingTime))
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
