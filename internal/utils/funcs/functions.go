package funcs

import (
	"fmt"
	"go.uber.org/zap"
	"main/internal/logger"
	"main/internal/utils/random"
	"main/internal/utils/vars"
	"math/rand"
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
	random.InitRandom()
	logger.InitLogger()
}

func Linkefy(linkParts ...string) string {
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
func RandomSleep() {
	sleepingTime := rand.Float64()*2 + 1
	zap.L().Info("Started random sleeping with time: " + fmt.Sprintf("%.2f", sleepingTime))
	time.Sleep(time.Duration(sleepingTime) * time.Second)
	zap.L().Info("Finished random sleeping")
}
