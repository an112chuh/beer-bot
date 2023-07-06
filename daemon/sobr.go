package daemon

import (
	"beer/config"
	"fmt"
	"time"
)

func ReduceSobriety() {
	for {
		db := config.GetConnection()
		query := `UPDATE users SET sobriety = IIF(sobriety = 100, sobriety, sobriety + 1)`
		_, err := db.Exec(query)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(10 * time.Minute)
	}
}

func FillCount() {
	db := config.GetConnection()
	type data struct {
		IDGroup  int
		IDUser   int
		BeerSum  float64
		VodkaSum float64
		WineSum  float64
	}
	var d []data
	query := `SELECT id_group, id_user, beer_sum, vodka_sum, wine_sum FROM users`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var r data
		err := rows.Scan(&r.IDGroup, &r.IDUser, &r.BeerSum, &r.VodkaSum, &r.WineSum)
		if err != nil {
			fmt.Println(err)
		}
		d = append(d, r)
	}
	for _, item := range d {
		count := 0.
		count += item.BeerSum / 0.5
		count += item.VodkaSum / 0.1
		count += item.WineSum / 0.2
		query := `UPDATE users SET count = $1 WHERE id_group = $2 AND id_user = $3`
		params := []any{count, item.IDGroup, item.IDUser}
		_, err := db.Exec(query, params...)
		if err != nil {
			fmt.Println(err)
		}
	}
}
