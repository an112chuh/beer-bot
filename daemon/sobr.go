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
