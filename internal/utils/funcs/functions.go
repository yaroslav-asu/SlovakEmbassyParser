package funcs

import (
	"main/internal/logger"
	"main/internal/session/captcha"
	"main/internal/utils/db"
	"main/internal/utils/vars"
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
	logger.Init()
	captcha.Init()
	db.Init()
}
