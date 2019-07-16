package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

// var err error
var dbConn *sql.DB
var err error

func GetDBConnection(db *sql.DB) {
	dbConn = db
}

func SelectFilePath(db *sql.DB) {
	var logFileName, logFilePath string
	rows, err := db.Query("SELECT log_file_path, log_file_name  from sys.server_file_audits")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&logFilePath, &logFileName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(logFilePath, logFileName)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
