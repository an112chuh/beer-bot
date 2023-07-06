package actions

import (
	"beer/config"
	"fmt"
	"time"
)

var BEER = 1
var VODKA = 2
var WINE = 3

var beerSobr = 10
var vodkaSobr = 50
var wineSobr = 30

var beerVolume = 0.5
var vodkaVolume = 0.1
var wineVolume = 0.2

func Drink(IDGroup int64, IDUser int64, userName string, drinkType int) (isFail bool, reason string) {
	checkUser(IDGroup, IDUser)
	if isSobrFail(IDGroup, IDUser) {
		return true, `Вы слишком пьяны, приходите как протрезвеете`
	}
	isFail, min, _ := isTimeoutFail(IDGroup, IDUser)
	if isFail {
		return true, fmt.Sprintf(`Приходите через %d мин`, 59-min)
	}
	drinking(IDGroup, IDUser, userName, drinkType)
	return false, ``
}

func drinking(IDGroup int64, IDUser int64, userName string, drinkType int) {
	db := config.GetConnection()
	var sobr int
	var volume float64
	query := `UPDATE users SET sobriety = sobriety - $1, `
	switch drinkType {
	case BEER:
		sobr = beerSobr
		volume = beerVolume
		query += `beer_sum = beer_sum + `
	case VODKA:
		sobr = vodkaSobr
		volume = vodkaVolume
		query += `vodka_sum = vodka_sum + `
	case WINE:
		sobr = wineSobr
		volume = wineVolume
		query += `wine_sum = wine_sum + `
	default:
		sobr = 0
		volume = 0
	}
	query += `$2, updated_at = $3, user_name = $4, count = count + 1 WHERE id_group = $5 AND id_user = $6`
	params := []any{sobr, volume, time.Now(), userName, IDGroup, IDUser}
	_, err := db.Exec(query, params...)
	if err != nil {
		fmt.Println(err)
	}
}
