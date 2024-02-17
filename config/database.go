package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DBsonek() *sql.DB{
	datasource := "root:r23password@/go_auth"
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		panic(err)
	}
	return db
}
