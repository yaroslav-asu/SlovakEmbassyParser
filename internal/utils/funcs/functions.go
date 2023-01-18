package funcs

import (
	"main/internal/logger"
	"main/internal/utils/random"
	"main/internal/utils/vars"
	"strings"
)

func StripString(s string) string {
	return stripStringRunes(s, ' ', '\n', '\t')
}

func stripStringRunes(s string, runes ...rune) string {
	for i, letter := range s {
		if !runeInList(letter, runes) {
			s = s[i:]
			break
		}
	}
	for i := len(s) - 1; i >= 0; i-- {
		if !runeInList(rune(s[i]), runes) {
			s = s[:i+1]
			break
		}
	}
	return s
}
func runeInList(checkingRune rune, runes []rune) bool {
	fl := false
	for _, r := range runes {
		if checkingRune == r {
			fl = true
			break
		}
	}
	return fl
}

func Init() {
	random.InitRandom()
	logger.InitLogger()
	vars.InitEnv()
}

func Linkefy(linkParts ...string) string {
	link := vars.SiteUrl
	for _, linkPart := range linkParts {
		link += strings.Replace(linkPart, "/", "", -1)
	}
	return link
}
