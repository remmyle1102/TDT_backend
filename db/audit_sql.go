package db

import (
	"time"

	"github.com/sirupsen/logrus"
)

// InsertAuditTask insert successful task to DB
func InsertAuditTask(name, location string, addBy, problemCount int) error {
	datetime := time.Now().Format("Mon Jan _2 15:04:05 2006")
	query := "INSERT INTO Report(name, problemCount, location, date, addBy) VALUES (@p1,@p2,@p3,@p4,@p5)"
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = stmt.Exec(name, problemCount, location, datetime, addBy)
	return err
}
