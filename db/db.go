package db

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

// var err error
var dbConn *sql.DB
var err error

func GetDBConnection(db *sql.DB) {
	dbConn = db
}
