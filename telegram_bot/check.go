package telegrambot

import (
	"beer/actions"
	"strings"
)

func isProfile(text string) bool {
	return strings.ToLower(text) == `профиль` || checkString(text, `profile`)
}

func isDrink(text string) (bool, int) {
	if strings.ToLower(text) == `пиво` || checkString(text, `beer`) {
		return true, actions.BEER
	}
	if strings.ToLower(text) == `водка` || checkString(text, `vodka`) {
		return true, actions.VODKA
	}
	if strings.ToLower(text) == `вино` || checkString(text, `wine`) {
		return true, actions.WINE
	}
	return false, -1
}

func isHelp(text string) bool {
	return checkString(text, `help`) || checkString(text, `start`)
}

func checkString(text string, template string) bool {
	return strings.ToLower(text) == `/`+template || strings.ToLower(text) == `/`+template+`@bar_beer_bot`
}
