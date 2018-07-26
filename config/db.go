package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB, err = sql.Open("mysql", "jv:xxx@/db_apocalipse")

func TryConn() {
	if err != nil {
		panic(err.Error())
	}

	if DB.Ping() != nil {
		panic(err.Error())
	}

	fmt.Println("Banco MySQL: Ok...")

}
