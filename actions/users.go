package actions

import (
	"beer/config"
	"fmt"
	"time"
)

type alcoSums struct {
	BeerSum  float64
	VodkaSum float64
	WineSum  float64
}

func GetData(IDGroup int64, IDUser int64) (sobriety int, sums alcoSums) {
	checkUser(IDGroup, IDUser)
	db := config.GetConnection()
	query := `SELECT sobriety, beer_sum, vodka_sum, wine_sum FROM users WHERE id_group = $1 AND id_user = $2`
	params := []any{IDGroup, IDUser}
	err := db.QueryRow(query, params...).Scan(&sobriety, &sums.BeerSum, &sums.VodkaSum, &sums.WineSum)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func addUser(IDGroup int64, IDUser int64) {
	db := config.GetConnection()
	query := `INSERT INTO users (id_group, id_user, created_at) VALUES ($1, $2, $3)`
	params := []any{IDGroup, IDUser, time.Now()}
	_, err := db.Exec(query, params...)
	if err != nil {
		fmt.Println(err)
	}
}

func checkUser(IDGroup int64, IDUser int64) {
	db := config.GetConnection()
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id_group = $1 AND id_user = $2)`
	params := []any{IDGroup, IDUser}
	var exists bool
	err := db.QueryRow(query, params...).Scan(&exists)
	if err != nil {
		fmt.Println(err)
	}
	if !exists {
		addUser(IDGroup, IDUser)
	}
}

func isSobrFail(IDGroup int64, IDUser int64) bool {
	sobr, _ := GetData(IDGroup, IDUser)
	return sobr < 0
}

func isTimeoutFail(IDGroup int64, IDUser int64) (bool, int, int) {
	db := config.GetConnection()
	var t string
	query := `SELECT updated_at FROM users WHERE id_group = $1 AND id_user = $2`
	params := []any{IDGroup, IDUser}
	err := db.QueryRow(query, params...).Scan(&t)
	if err != nil {
		fmt.Println(err)
	}
	layout := "2006-01-02 15:04:05-07:00"
	date, err := time.Parse(layout, t)
	if err != nil {
		fmt.Println(err)
	}
	if time.Since(date) < 1*time.Hour {
		return true, int(time.Since(date).Minutes()), int(time.Since(date).Seconds())
	}
	return false, 0, 0
}
