package server

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

func ConnectDatabase() *sql.DB {
	var db *sql.DB
	configuration := mysql.NewConfig()
	(*configuration).Net = "tcp"
	(*configuration).Addr = myEnvironment["GALLERY_HOST"]
	(*configuration).User = myEnvironment["GALLERY_USER"]
	(*configuration).Passwd = myEnvironment["GALLERY_PASSWORD"]
	(*configuration).DBName = myEnvironment["GALLERY_DATABASE"]
	(*configuration).ParseTime = true

	db, err := sql.Open("mysql", configuration.FormatDSN())
	if err != nil {
		panic(err)
	}

	return db
}
